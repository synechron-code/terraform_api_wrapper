package main

import (
	"flag"
	"fmt"
	"github.com/citihub/terraform_api_wrapper/runner"
	"os"
)

func main() {

	dollarhome, _ := os.UserHomeDir()
	var (
		port              = flag.Int("port", 8080, "Port to serve the API on")
		planLocation      = flag.String("plan-location", fmt.Sprintf("%s/terraform/plans", dollarhome), "Top level directory containing your terraform plans")
		workspaceLocation = flag.String("workspace-location", fmt.Sprintf("%s/terraform/workspace", dollarhome), "Top level directory where context workspace is set up (avoid conflicts caused by running multiple contexts")
	)
	flag.Parse()

	fmt.Printf("Serving API on port %d\n", port)
	runner.API_runner(port, *planLocation, *workspaceLocation)
}
