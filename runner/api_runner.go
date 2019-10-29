package runner

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/citihub/terraform_api_wrapper/handler"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func queryJobStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	status := make(map[string]string)
	status["jobid"] = vars["jobid"]

	jobID, err := uuid.Parse(vars["jobid"])
	if err != nil {
		json.NewEncoder(w).Encode(fmt.Sprintf("{error: %v}", err))
		return
	}

	switch handler.QueryJobStatus(jobID) {
	case handler.CREATED:
		status["status"] = "created"
	case handler.QUEUED:
		status["status"] = "queued"
	case handler.RUNNING:
		status["status"] = "running"
	case handler.COMPLETE:
		status["status"] = "complete"
	case handler.ASSERTING:
		status["status"] = "asserting"
	case handler.JOBERROR:
		status["status"] = "error"
	default:
		status["status"] = "unknown"
	}

	json.NewEncoder(w).Encode(status)
}

func createContext(w http.ResponseWriter, r *http.Request) {
	var jobContext handler.JsonJobContext
	var vendor int

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		json.NewEncoder(w).Encode(fmt.Sprintf("{error: %v}", err))
		return
	}

	jsonerr := json.Unmarshal(body, &jobContext)
	if jsonerr != nil {
		json.NewEncoder(w).Encode(fmt.Sprintf("{error: %v}", err))
		return
	}

	switch jobContext.Vendor {
	case "aws":
		vendor = handler.AWS
	case "azure":
		vendor = handler.AZURE
	case "gcp":
		vendor = handler.GCP
	default:
		vendor = 99
	}

	json.NewEncoder(w).Encode(handler.CreateJobContext(vendor, jobContext.Credentials, jobContext.Statefiles))
}

func runJob(w http.ResponseWriter, r *http.Request) {
	var jobInstructions handler.JobInstructions
	var action int

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		json.NewEncoder(w).Encode(fmt.Sprintf("{error: %v}", err))
		//panic(err)
		return
	}

	jsonerr := json.Unmarshal(body, &jobInstructions)
	if jsonerr != nil {
		json.NewEncoder(w).Encode(fmt.Sprintf("{error: %v}", err))
		//panic(jsonerr)
		return
	}

	vars := mux.Vars(r)

	switch vars["action"] {
	case "apply":
		action = handler.APPLY
	case "plan":
		action = handler.PLAN
	case "destroy":
		action = handler.DESTROY
	default:
		json.NewEncoder(w).Encode(fmt.Sprintf("{error: action not recognised}"))
		return
	}

	job := handler.CreateJob(jobInstructions, handler.JobContexts[jobInstructions.ContextID], action, vars["stage"])

	go handler.JobHandler(job)

	json.NewEncoder(w).Encode(job)
}

func getJobResponse(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	jobId, err := uuid.Parse(vars["jobid"])
	if err != nil {
		json.NewEncoder(w).Encode(fmt.Sprintf("{error: %v}", err))
	} else {
		json.NewEncoder(w).Encode(handler.GetJobResponse(jobId))
	}
}

func API_runner(port *int, plan_location string, context_location string) {
	handler.JobHandlerInit(plan_location)
	handler.ContextsHandlerInit(plan_location, context_location)

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/v1/context/create", createContext)

	router.HandleFunc("/v1/job/create/{stage}/{action}", runJob)
	router.HandleFunc("/v1/job/response/{jobid}", getJobResponse)
	router.HandleFunc("/v1/job/status/{jobid}", queryJobStatus)

	http.Handle("/", router)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}
