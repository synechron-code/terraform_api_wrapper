package handler

import (
	"github.com/google/uuid"
)

type JobContext struct {
	ContextID   uuid.UUID
	Statefiles  map[string]string
	Vendor      int
	Credentials map[string]string
}

var JobContexts map[uuid.UUID]JobContext

func ContextsHandlerInit() {
	JobContexts = make(map[uuid.UUID]JobContext)
}

func CreateNewJobContext() uuid.UUID {
	contextID := uuid.New()
	JobContexts[contextID] = JobContext{ContextID: contextID}

	return contextID
}

func SetCredentials(contextID uuid.UUID, credentials map[string]string) {
	var jobContext = JobContexts[contextID]

	jobContext.Credentials = credentials
	JobContexts[contextID] = jobContext
}

func SetVendor(contextID uuid.UUID, vendor int) {
	var jobContext = JobContexts[contextID]

	jobContext.Vendor = vendor
	JobContexts[contextID] = jobContext
}

func SetStateFiles(contextID uuid.UUID, statefiles map[string]string) {
	var jobContext = JobContexts[contextID]

	jobContext.Statefiles = statefiles
	JobContexts[contextID] = jobContext
}

func CreateJobContext(vendor int, credentials map[string]string, statefiles map[string]string) map[string]string {
	contextID := CreateNewJobContext()
	SetVendor(contextID, vendor)
	SetCredentials(contextID, credentials)
	SetStateFiles(contextID, statefiles)

	var returnVal = make(map[string]string)

	returnVal["context_id"] = contextID.String()
	return returnVal
}
