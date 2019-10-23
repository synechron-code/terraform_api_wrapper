package main

import (
	"citihub.com/terraform_api_wrapper/runner"
	"flag"
)

func main() {

	var (
		port         = flag.Int("port", 8080, "Port to serve the API on")
		planLocation = flag.String("plan-location", "~/terraform", "Top level directory containing your terraform plans")
	)
	flag.Parse()

	fmt.Fprintf("Serving API on port %d", port)
	runner.API_runner(port, *planLocation)
}
