package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"citihub.com/terraform_api_wrapper/handler"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

const (
	APPLY = iota
	PLAN
	DESTROY
)

const (
	CREATED = iota
	QUEUED
	RUNNING
	COMPLETE
)

var jobChan chan *handler.Job

func queryJobStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	status := make(map[string]string)
	status["jobid"] = vars["jobid"]

	jobID, _ := uuid.Parse(vars["jobid"])

	switch handler.QueryJobStatus(jobID) {
	case CREATED:
		status["status"] = "created"
	case QUEUED:
		status["status"] = "queued"
	case RUNNING:
		status["status"] = "running"
	case COMPLETE:
		status["status"] = "complete"
	default:
		status["status"] = "unknown"
	}

	json.NewEncoder(w).Encode(status)
}

/*
to do functions

func getJobResponse(w http.ResponseWriter, r *http.Request) {

}


func setCredentials(w http.ResponseWriter, r *http.Request) {

}

func setVendor(w http.ResponseWriter, r *http.Request) {

}

func setStatefiles(w http.ResponseWriter, r *http.Request) {

}


*/

func runJob(w http.ResponseWriter, r *http.Request) {
	fmt.Println("createJob hit")

	var jobInstructions handler.JobInstructions
	var action int

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}

	jsonerr := json.Unmarshal(body, &jobInstructions)
	if jsonerr != nil {
		panic(jsonerr)
	}

	vars := mux.Vars(r)

	switch vars["action"] {
	case "apply":
		action = APPLY
	case "plan":
		action = PLAN
	case "destroy":
		action = DESTROY
	default:
		//panic
	}

	job := handler.CreateJob(jobInstructions, action, vars["stage"])
	go handler.JobHandler(job)

	json.NewEncoder(w).Encode(job)
}

func api_runner(port *int) {

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/v1/plan/{action}/{stage}", runJob)
	router.HandleFunc("/v1/query/status/{jobid}", queryJobStatus)

	/*
		//to do functions

		router.HandleFunc("/v1/context/statefiles", setStatefiles)
		router.HandleFunc("/v1/context/vendor", setVendor)
		router.HandleFunc("/v1/context/credentials", setCredentials)
		router.HandleFunc("/v1/query/response/{jobid}", getJobResponse)
	*/

	http.Handle("/", router)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {
	handler.JobHandlerInit()

	var (
		port = flag.Int("port", 8080, "Port to serve the API on")
	)
	flag.Parse()

	api_runner(port)
}
