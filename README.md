# terraform_api_wrapper
This is a lightweight HTTP API wrapper around Terraform. It uses terratest from gruntworks under the hood to run Terraform.

## Why might you want to use this?
- you've created an approved infrastructure "build" using Terraform and want to memorialise it in a container, but externalise the config and need a way to push config at runtime
- you have a complex multi-stage pipeline and want to orchestrate it 

## How to use it
There are two steps to using this API
1. Create an execution context. The execution context contains your provider credentials and statefile locations.
2. Create jobs. The job contains the execution context, variables and remote state names (matching the names configured in your context) needed to run your plan.

## Example
Example of creating a resource group and a managed disk in Azure

1. Create an execution context
```
    {
        "vendor": "azure",
        "credentials": {
            "ARM_CLIENT_ID":"clientid_goes_here",
            "ARM_CLIENT_SECRET":"secret_goes_here",
            "ARM_SUBSCRIPTION_ID":"subscription_id_goes_here",
            "ARM_TENANT_ID":"tenant_id_goes_here"
        },
        "statefiles": {
            "resource_group": "/home/ian/test/resource_group/terraform.tfstate",
            "disk": "/home/ian/test/disk/terraform.tfstate"
        }
    }
```
2. POST it
```
curl -X POST -d @context.json localhost:8080/v1/context/create
{"context_id":"2b4b3fae-fb8e-4750-b06b-7378343cfaf1"}
```
3. Create the resource group json
```
    {
        "context_id": "2b4b3fae-fb8e-4750-b06b-7378343cfaf1",
        "tfvars": {
            "name": "test_rg"
        }
    }
```
4. Create the resource_group YAML and put it in a folder called "resource_group" (the directory name has to match the value in "statefiles" in the context, and the plan name in the job/create/{plan}/{action} URI)

5. POST it 
```
    curl -X POST -d @rg_instructions.json localhost:8080/v1/job/create/resource_group/apply
    {"JobID":"2b4b3fae-fb8e-4750-b06b-7378343cfaf1"}
```
6. Create the disk json
```
    {
            "context_id": "2b4b3fae-fb8e-4750-b06b-7378343cfaf1",
            "tfvars": {
                    "name": "testdisk",
                    "disk_size": "1"
            },
            "remote_states": [
                    "resource_group"
            ]
    }
```
7. Create the disk YAML and put it in a folder called "disk"

8. POST it
```
curl -X POST -d @disk_instructions.json localhost:8080/v1/job/create/disk/apply
{"JobID":"0de42156-3ce8-4c30-a4e1-60e74e8a8b37"}
```
## Checking the job status
Because creating resources in public cloud can involve long-running jobs, the instructions POSTed to the endpoint will immediately return a JobID while Terraform runs asynchronously.
You can check the status of a running job by calling /job/status/{jobId}

curl -X POST localhost:8080/v1/job/status/2b4b3fae-fb8e-4750-b06b-7378343cfaf1
{"jobid":"2b4b3fae-fb8e-4750-b06b-7378343cfaf1","status":"running"}

## Getting the job output
Once a job status has transitioned to "complete" the output will be available to view.  It will need some parsing to pretty print it. You can use terratest's log analyser to view it.

curl -X POST localhost:8080/v1/job/response/2b4b3fae-fb8e-4750-b06b-7378343cfaf1
{"TfOutput":"\u001b[0m\u001b[1mdata.terraform_remote_state.rg_group: Refreshing state...\u001b[0m\n\u001b[0m\u001b[1mazurerm_managed_disk.test: Creating...\u001b[0m\n  create_option:        \"\" =\u003e \"Empty\"\n  disk_iops_read_write: \"\" =\u003e \"\u003ccomputed\u003e\"\n  disk_mbps_read_write: \"\" =\u003e \"\u003ccomputed\u003e\"\n  disk_size_gb:         \"\" =\u003e \"1\"\n  location:             \"\" =\u003e \"westus\"\n  name:                 \"\" =\u003e \"testdisk\"\n  resource_group_name:  \"\" =\u003e \"test_rg\"\n  source_uri:           \"\" =\u003e \"\u003ccomputed\u003e\"\n  storage_account_type: \"\" =\u003e \"Standard_LRS\"\n  tags.%:               \"\" =\u003e \"\u003ccomputed\u003e\"\u001b[0m\n\u001b[0m\u001b[1mazurerm_managed_disk.test: Creation complete after 2s (ID: /subscriptions/ea9a7509-82bb-430a-857f-...iders/Microsoft.Compute/disks/testdisk)\u001b[0m\u001b[0m\n\u001b[0m\u001b[1m\u001b[32m\nApply complete! Resources: 1 added, 0 changed, 0 destroyed.\u001b[0m","TfError":null}