package handler

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/gruntwork-io/terratest/modules/terraform"
)

const (
	APPLY = iota
	PLAN
	DESTROY
)

const (
	CREATED = iota
	QUEUED
	RUNNING
	COMPLETE
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

func CreateJob(jobInstructions JobInstructions, action int, stage string) *Job {
	//create the Job object and add it to the Jobs map.
	//TFOptions struct for this job and put it into TF

	tfOptions := terraform.Options{
		Vars:         jobInstructions.TfVars,
		TerraformDir: "/terraform/" + stage,
		//BackendConfig: JobContexts[jobInstructions.ContextID].Statefiles[jobInstructions.Stage],
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
	fmt.Println("Entering JobHandler")
	job.Status = RUNNING

	fmt.Println(job.JobID)
	fmt.Println("Set status of job to running")
	
	time.Sleep(20 * time.Second)
	job.Response.TfOutput = "hello"
	//job.Response.TfOutput, job.Response.TfError := terraform.InitAndApplyE(blah)
	
	fmt.Println(job.JobID)
	fmt.Println("Set status of job to complete")
	
	job.Status = COMPLETE
}

func QueryJobStatus(jobId uuid.UUID) int {
	return Jobs[jobId].Status
}

func GetJobResponse(jobId uuid.UUID) *JobResponse {
	return &Jobs[jobId].Response
}
