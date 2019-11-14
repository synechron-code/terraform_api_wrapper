package handler

import (
	"fmt"
	"github.com/google/uuid"
	"os/exec"
)

type JobContext struct {
	ContextID   uuid.UUID
	Statefiles  map[string]string
	Vendor      int
	Credentials map[string]string
}

var JobContexts map[uuid.UUID]JobContext

var planLocation_g string
var contextLocation_g string

func ContextsHandlerInit(planLocation string, contextLocation string) {
	JobContexts = make(map[uuid.UUID]JobContext)
	planLocation_g = planLocation
	contextLocation_g = contextLocation
}

func CreateNewJobContext() uuid.UUID {
	contextID := uuid.New()
	JobContexts[contextID] = JobContext{ContextID: contextID}

	//mkdir context-location/contextID
	mkContextDir := exec.Command("mkdir", fmt.Sprintf("%s/%s", contextLocation_g, contextID))
	cpPlansToContext := exec.Command("sh", "-c", fmt.Sprintf("cp -r %s/* %s/%s", planLocation_g, contextLocation_g, contextID.String()))
	//cp plansLocation context-location/contextID

	if out, err := mkContextDir.CombinedOutput(); err != nil {
		fmt.Printf("Error making directory %v\n", err)
	} else {
		fmt.Printf("workspace mkdir: %s", out)
	}

	if cpout, err := cpPlansToContext.CombinedOutput(); err != nil {
		fmt.Printf("Error making directory %v\n", err)
	} else {
		fmt.Printf("workspace cp: %v", cpout)
	}
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

/*
func CreatePlanContext(contextID uuid.UUID, plansLocation string) {

}
*/

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
