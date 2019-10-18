package handler

import (
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/google/uuid"
)

const (
	APPLY = iota
	PLAN
	DESTROY
)

const (
	QUEUED = iota
	RUNNING
	COMPLETE
)

type JobRequest struct {
	JobContextUUID uuid.UUID
	tfOptions terraform.Options
	Stage string
}

type JobResponse struct {
	TfOutput string
	TfError error 
}

type Job struct (
	JobID uuid.UUID
	Request JobRequest
	Response JobResponse
	Status int
)

var Jobs map[uuid.UUID]Job

func createJob(contextID uuid.UUID, action int, tfVars interface{}, stage string) (&Job) {
	//create the Job object and add it to the Jobs map.
	//TFOptions struct for this job and put it into TF
}

func jobHandler(jobsChan <-chan *Job) {
		for job := range jobsChan {
			job.Status = RUNNING
			//job.Response.TfOutput, job.Response.TfError := terraform.InitAndApplyE(blah)	
			job.Status = COMPLETE
		}
}

func queryJobStatus(jobId uuid.UUID) int {
	return Jobs[jobId].Status
}

func getJobResponse(jobId uuid.UUID) (JobResponse) {
	return Jobs[jobId].Response
}
