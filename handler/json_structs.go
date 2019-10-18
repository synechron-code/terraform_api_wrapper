package handler

import (
	"github.com/google/uuid"
)

type JobInstructions struct {
	ContextID uuid.UUID              `json:"context_id"`
	TfVars    map[string]interface{} `json:"tfvars"`
	RemoteStates []string			`json:"remote_states"`
}

type AzureCredentials struct {
	ARM_Tenant_ID       string `json:"arm_tenant_id"`
	ARM_Subscription_ID string `json:"arm_subscription_id"`
	ARM_Client_ID       string `json:"arm_client_id"`
	ARM_Client_Secret   string `json:"arm_client_secret`
}

type StatefileLocations struct {
	Statefiles map[string]string `json:"statefiles"`
}

type Vendor struct {
	Vendor string `json:"vendor"`
}

type JsonJobContext struct {
	Vendor string `json:"vendor"`
	Statefiles map[string]string	`json:"statefiles"`	
	Credentials map[string]string `json:"credentials"`
}