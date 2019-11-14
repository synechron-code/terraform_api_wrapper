package handler

const (
	APPLY = iota
	PLAN
	DESTROY
)

const (
	CREATED   = iota
	QUEUED    //Job is queued for execution by Terraform
	RUNNING   //Terraform is currently running this Job
	ASSERTING //working out whether there was an error in Terraform execution
	COMPLETE  //Job completed successfully
	JOBERROR  //There was an error outside of Terraform execution
	TFERROR   //Terrform returned errors
	TFWARNING //Terraform returned warnings but no errors
)

const (
	AWS = iota
	AZURE
	GCP
)
