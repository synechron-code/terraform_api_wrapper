package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"testing"
)

type JobRequest struct {
	Action    int
	tfOptions terraform.Options
	Stage     string
}

type JobResponse struct {
	TfOutput string
	TfError  error
}

type Job struct {
	JobID     uuid.UUID
	ContextID uuid.UUID
	Request   JobRequest
	Response  JobResponse
	Status    int
}

var Jobs map[uuid.UUID]*Job
var planLocation string
var contextLocation string

func JobHandlerInit(plan_location string) {
	Jobs = make(map[uuid.UUID]*Job)
	planLocation = plan_location
}

func CreateJob(jobInstructions JobInstructions, jobContext JobContext, action int, stage string) *Job {
	//create the Job object and add it to the Jobs map.
	//TFOptions struct for this job and put it into TF

	backendConfig := map[string]interface{}{
		"path": jobContext.Statefiles[stage],
	}

	credentials := JobContexts[jobInstructions.ContextID].Credentials

	vars := jobInstructions.TfVars

	for _, remoteState := range jobInstructions.RemoteStates {
		vars["remote_state_"+remoteState] = JobContexts[jobInstructions.ContextID].Statefiles[remoteState]
	}

	tfOptions := terraform.Options{
		Vars:          vars,
		TerraformDir:  fmt.Sprintf("%v/%v/%v", contextLocation, jobContext.ContextID, stage),
		BackendConfig: backendConfig,
		EnvVars:       credentials,
	}

	request := JobRequest{
		Action:    action,
		Stage:     stage,
		tfOptions: tfOptions,
	}

	newJob := Job{
		JobID:     uuid.New(),
		ContextID: jobInstructions.ContextID,
		Request:   request,
		Status:    CREATED,
	}

	Jobs[newJob.JobID] = &newJob

	fmt.Println(request.tfOptions)

	return &newJob
}

func AssertJobStatus(job *Job) {
	job.Status = COMPLETE
}

func JobHandler(job *Job) {
	var tfOutput string
	var tfError error

	fmt.Println("Entering JobHandler")
	job.Status = RUNNING

	fmt.Println(fmt.Sprintf("Set status of job %s to running", job.JobID))

	t := new(testing.T)

	switch job.Request.Action {
	case APPLY:
		fmt.Println("running apply")
		tfOutput, tfError = terraform.InitAndApplyE(t, &job.Request.tfOptions)
		//TODO: improve job Status based on Terratest assertion
	case PLAN:
		fmt.Println("running plan")
		tfOutput, tfError = terraform.InitAndPlanE(t, &job.Request.tfOptions)
		//TODO: improve job Status based on Terratest assertion
	case DESTROY:
		fmt.Println("running destroy")
		tfOutput, tfError = terraform.InitE(t, &job.Request.tfOptions)
		if tfError != nil {
			break
		}
		tfOutput, tfError = terraform.DestroyE(t, &job.Request.tfOptions)
		//TODO: improve job Status based on Terratest assertion
	default:
		tfOutput = ""
		errorMessage, _ := json.RawMessage("\"JobHandler Error\": \"Action not recognised\"").MarshalJSON()
		tfError = errors.New(fmt.Sprintf("%v", errorMessage))
		job.Status = JOBERROR
		//panic
	}

	job.Response.TfOutput = tfOutput
	job.Response.TfError = tfError

	AssertJobStatus(job)

	fmt.Println(fmt.Sprintf("Set status of job %s to running", job.JobID))
}

func QueryJobStatus(jobId uuid.UUID) int {
	return Jobs[jobId].Status
}

func GetJobResponse(jobId uuid.UUID) JobResponse {
	return Jobs[jobId].Response
}
