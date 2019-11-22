package handler

import (
	"github.com/google/uuid"
)

type JobInstructions struct {
	ContextID    uuid.UUID              `json:"context_id"`
	TfVars       map[string]interface{} `json:"tfvars"`
	RemoteStates []string               `json:"remote_states"`
}

type StatefileLocations struct {
	Statefiles map[string]string `json:"statefiles"`
}

type Vendor struct {
	Vendor string `json:"vendor"`
}

type JsonJobContext struct {
	Vendor           string            `json:"vendor"`
	Statefiles       map[string]string `json:"statefiles"`
	Credentials      map[string]string `json:"credentials"`
	Certificate_Data []CertificateData `json:"certificate_data"`
}

type CertificateData struct {
	CredentialName string `json:"credential_name"`
	Data           string `json:"data"`
	Type           string `json:"encoding"`
}
