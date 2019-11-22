package handler

import (
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"

	"github.com/google/uuid"
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

func CreateNewJobContext() (uuid.UUID, error) {
	contextID := uuid.New()
	JobContexts[contextID] = JobContext{ContextID: contextID}

	//mkdir context-location/contextID
	mkContextDir := exec.Command("sh", "-c", fmt.Sprintf("mkdir -p %s/%s", contextLocation_g, contextID))
	cpPlansToContext := exec.Command("sh", "-c", fmt.Sprintf("cp -r %s/* %s/%s", planLocation_g, contextLocation_g, contextID.String()))
	//cp plansLocation context-location/contextID

	if out, err := mkContextDir.CombinedOutput(); err != nil {
		fmt.Printf("Error making directory %s, error is %s\n", fmt.Sprintf("%s/%s", contextLocation_g, contextID), out)
		return contextID, err
	} else {
		fmt.Printf("workspace mkdir: %s", out)
	}

	if cpout, err := cpPlansToContext.CombinedOutput(); err != nil {
		fmt.Printf("Error copying plans to workspace, error is %s\n", cpout)
		return contextID, err
	} else {
		fmt.Printf("workspace cp: %v", cpout)
	}

	return contextID, nil
}

func SetCredentials(contextID uuid.UUID, credentials map[string]string) {

	jobContext := JobContexts[contextID]

	jobContext.Credentials = credentials

	JobContexts[contextID] = jobContext
}

func SetCertificates(contextID uuid.UUID, certificateData map[string]string) error {
	jobContext := JobContexts[contextID]

	for k, v := range certificateData {
		certLoc := fmt.Sprintf("%s/%s/%s", contextLocation_g, contextID, k)
		data, err := base64.StdEncoding.DecodeString(v)
		if err != nil {
			fmt.Printf("Error base64 decoding certificate data for %s", k)
			return err
		}

		f, err := os.Create(certLoc)
		if err != nil {
			fmt.Printf("Error %v creating certificate file in location %s", err, certLoc)
			return err
		}
		defer f.Close()

		if _, err := f.Write(data); err != nil {
			fmt.Printf("Error %v writing certificate file for %s to location %s", err, k, certLoc)
			return err
		}
		jobContext.Credentials[k] = certLoc
	}

	JobContexts[contextID] = jobContext
	return nil
}

func SetVendor(contextID uuid.UUID, vendor int) {
	jobContext := JobContexts[contextID]

	jobContext.Vendor = vendor
	//JobContexts[contextID] = jobContext
}

/*
func CreatePlanContext(contextID uuid.UUID, plansLocation string) {

}
*/

func SetStateFiles(contextID uuid.UUID, statefiles map[string]string) {
	jobContext := JobContexts[contextID]

	jobContext.Statefiles = statefiles
	//JobContexts[contextID] = jobContext
}

func CreateJobContext(vendor int, credentials map[string]string, certificateData map[string]string, statefiles map[string]string) (map[string]string, error) {
	contextID, err := CreateNewJobContext()
	if err != nil {
		return nil, err
	}

	SetVendor(contextID, vendor)
	SetCredentials(contextID, credentials)

	if err := SetCertificates(contextID, certificateData); err != nil {
		return nil, err
	}
	SetStateFiles(contextID, statefiles)

	//we return a map to facilitate the API json response body
	var returnVal = make(map[string]string)

	returnVal["context_id"] = contextID.String()
	return returnVal, nil
}
