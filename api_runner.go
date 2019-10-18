package api_wrapper

import (
	"handler"
)

func api_runner(maxQueueSize int, maxWorkers int, port int) {

	/*var (
		maxQueueSize = flag.Int("max_queue_size", 5, "The maximum size of the terraform execution jobs queue")
		maxWorkers = flag.Int("max_workers", 1, "The maximum number of concurrent terraform executions. Best keep this to 1 unless you're confident")
		port = flag.Int("port", 8080, "Port to serve the API on")
	)
	flag.Parse()*/

	jobsChan := make(chan handler.Job, *maxQueueSize)

	for i := 1; i <= *maxWorkers; i++ {
		go handler.JobHandler(jobsChan)
	}

	// handler for adding jobs
	http.HandleFunc("/v1/plan/{action}/{stage}", func(w http.ResponseWriter, r *http.Request) {
		jobsChan <- handler.createJob(w, r)
	})

	http.HandleFunc("/v1/context/statefiles", func(w http.ResponseWriter, r *http.Request) {
		handler.setStatefiles(w, r)
	})

	http.HandleFunc("/v1/context/vendor", func(w http.ResponseWriter, r *http.Request) {
		jobsChan <- handler.setVendor(w, r)
	})

	http.HandleFunc("/v1/context/credentials", func(w http.ResponseWriter, r *http.Request) {
		jobsChan <- handler.setCredentials(w, r)
	})

	http.HandleFunc("/v1/query/{job}/status", func(w http.ResponseWriter, r *http.Request) {
		jobsChan <- handler.queryJobStatus(w, r)
	})

	http.HandleFunc("/v1/query/{job}/response", func(w http.ResponseWriter, r *http.Request) {
		jobsChan <- handler.getResponse(w, r)
	})

	log.Fatal(http.ListenAndServe(":"+*port, nil))	
}