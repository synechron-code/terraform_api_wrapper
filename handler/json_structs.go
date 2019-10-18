//JSON structures for TF stuff
type JobInstructions struct {
	ContextID			uuid.UUID   	`json:"context_id"`
	Action				string			`json:"action"`
	Stage				string			`json:"stage"`
	TfVars				interface{}		`json:"tfvars"`
}

type AzureCredentials struct {
	ARM_Tenant_ID		string			`json:"arm_tenant_id"`
	ARM_Subscription_ID string 			`json:"arm_subscription_id"`
	ARM_Client_ID 		string 			`json:"arm_client_id"`
	ARM_Client_Secret 	string 			`json:"arm_client_secret`
}

type StatefileLocations struct {
	Statefiles 			map[string]string `json:"statefiles"`
}

type Vendor struct {
	Vendor				string			`json:"vendor"`
}