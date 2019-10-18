package handler

import (
	"github.com/google/uuid"
)

const (
	AWS = iota
	AZURE
	GCP
)

type JobContext struct {
	ContextID   uuid.UUID
	Statefiles  []string
	Vendor      int
	Credentials interface{}
}

var JobContexts map[uuid.UUID]JobContext

func CreateNewJobContext() uuid.UUID {
	contextID := uuid.New()
	JobContexts[contextID] = JobContext{ContextID: contextID}

	return contextID
}

func SetCredentials(contextID uuid.UUID, credentials interface{}) {
	var jobContext = JobContexts[contextID]

	jobContext.Credentials = credentials
	JobContexts[contextID] = jobContext
}

func SetVendor(contextID uuid.UUID, vendor int) {
	var jobContext = JobContexts[contextID]

	jobContext.Vendor = vendor
	JobContexts[contextID] = jobContext
}

func SetStateFiles(contextID uuid.UUID, statefiles []string) {
	var jobContext = JobContexts[contextID]

	jobContext.Statefiles = statefiles
	JobContexts[contextID] = jobContext
}
