# terraform_api_wrapper
An HTTP API wrapper around Terraform. 

Why might you want to use this?
- you've created an approved infrastructure "build" using Terraform and want to memorialise it in a container, but keep the config external
- you have a complex multi-stage pipeline and want to orchestrate it from something else
- ...