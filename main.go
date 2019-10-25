package main

import (
	"flag"
	"fmt"
	"github.com/citihub/terraform_api_wrapper/runner"
)

func main() {

	var (
		port         = flag.Int("port", 8080, "Port to serve the API on")
		planLocation = flag.String("plan-location", "~/terraform", "Top level directory containing your terraform plans")
	)
	flag.Parse()

	fmt.Printf("Serving API on port %v", port)
	runner.API_runner(port, *planLocation)
}
