package runner

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"citihub.com/terraform_api_wrapper/handler"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func queryJobStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	status := make(map[string]string)
	status["jobid"] = vars["jobid"]

	jobID, _ := uuid.Parse(vars["jobid"])

	switch handler.QueryJobStatus(jobID) {
	case handler.CREATED:
		status["status"] = "created"
	case handler.QUEUED:
		status["status"] = "queued"
	case handler.RUNNING:
		status["status"] = "running"
	case handler.COMPLETE:
		status["status"] = "complete"
	default:
		status["status"] = "unknown"
	}

	json.NewEncoder(w).Encode(status)
}

func createContext(w http.ResponseWriter, r *http.Request) {
	fmt.Println("createContext hit")

	var jobContext handler.JsonJobContext
	var vendor int

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}

	jsonerr := json.Unmarshal(body, &jobContext)
	if jsonerr != nil {
		panic(jsonerr)
	}

	switch jobContext.Vendor {
	case "aws":
		vendor = handler.AWS
	case "azure":
		vendor = handler.AZURE
	case "gcp":
		vendor = handler.GCP
	default:
		//panic
		vendor = 99
	}

	json.NewEncoder(w).Encode(handler.CreateJobContext(vendor, jobContext.Credentials, jobContext.Statefiles))
}

func runJob(w http.ResponseWriter, r *http.Request) {
	fmt.Println("runJob hit")

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
		action = handler.APPLY
	case "plan":
		action = handler.PLAN
	case "destroy":
		action = handler.DESTROY
	default:
		//panic
	}

	job := handler.CreateJob(jobInstructions, handler.JobContexts[jobInstructions.ContextID], action, vars["stage"])
	go handler.JobHandler(job)

	json.NewEncoder(w).Encode(job)
}

func getJobResponse(w http.ResponseWriter, r *http.Request) {
	fmt.Println("getJobResponse hit")

	vars := mux.Vars(r)

	jobId, _ := uuid.Parse(vars["jobid"])

	json.NewEncoder(w).Encode(handler.GetJobResponse(jobId))
}

func Api_runner(port *int, plan_location string) {
	handler.JobHandlerInit(plan_location)
	handler.ContextsHandlerInit()

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/v1/context/create", createContext)

	router.HandleFunc("/v1/job/create/{stage}/{action}", runJob)
	router.HandleFunc("/v1/job/response/{jobid}", getJobResponse)
	router.HandleFunc("/v1/job/status/{jobid}", queryJobStatus)

	http.Handle("/", router)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
