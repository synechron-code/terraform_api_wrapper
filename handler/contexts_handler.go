package handler

import (
	"github.com/google/uuid"
)

const (
	AWS = iota
	AZURE
	GCP
)

type JobContext struct (
	statefiles []string
	vendor int
	credentials interface{}
)

var JobContexts map[uuid.UUID]JobContext

func createNewJobContext() (uuid.UUID) {
	contextID := uuid.New()
	JobContexts[contextID] = new JobContext{}

	return contextID
}

func setCredentials(contextID uuid.UUID, credentials interface{}) {
	JobContexts[contextID].credentials = credentials
}

func setVendor(contextID uuid.UUID, vendor int) {
	JobContexts[contextID].vendor = vendor
}

func setStateFiles(contextID uuid.UUID, statefiles []string) {
	JobContexts[contextID].statefiles = statefiles
}