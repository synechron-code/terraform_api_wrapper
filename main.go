package main

import (
	"citihub.com/terraform_api_wrapper/runner"
	"flag"
)

func main() {

	var (
		port = flag.Int("port", 8080, "Port to serve the API on")
	)
	flag.Parse()

	runner.Api_runner(port)
}
