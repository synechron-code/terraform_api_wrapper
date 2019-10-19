package handler

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"testing"
	//"time"
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

func JobHandlerInit() {
	Jobs = make(map[uuid.UUID]*Job)
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
		TerraformDir:  "/home/ian/test/" + stage,
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

func JobHandler(job *Job) {
	var tfOutput string
	var tfError error

	fmt.Println("Entering JobHandler")
	job.Status = RUNNING

	fmt.Println(job.JobID)
	fmt.Println("Set status of job to running")

	//time.Sleep(20 * time.Second)
	//job.Response.TfOutput = "hello"

	t := new(testing.T)

	switch job.Request.Action {
	case APPLY:
		fmt.Println("running apply")
		tfOutput, tfError = terraform.InitAndApplyE(t, &job.Request.tfOptions)
	case PLAN:
		fmt.Println("running plan")
		tfOutput, tfError = terraform.InitAndPlanE(t, &job.Request.tfOptions)
	case DESTROY:
		fmt.Println("running destroy")
		tfOutput, tfError = terraform.DestroyE(t, &job.Request.tfOptions)
	default:
		tfOutput = "none"
		//panic
	}

	job.Response.TfOutput = tfOutput
	job.Response.TfError = tfError

	fmt.Println(job.JobID)
	fmt.Println("Set status of job to complete")

	job.Status = COMPLETE
}

func QueryJobStatus(jobId uuid.UUID) int {
	return Jobs[jobId].Status
}

func GetJobResponse(jobId uuid.UUID) JobResponse {
	return Jobs[jobId].Response
}
