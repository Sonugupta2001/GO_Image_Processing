package api

import (
	"encoding/json"
	"net/http"
	"retailpulse-image-service/pkg/job"
	"sync"
)


var jobs = make(map[int64]*job.Job)
var jobsMutex sync.Mutex
var jobCounter int64

type JobRequest struct {
	Count  int
	Visits []job.StoreJobRequest
}


func jobProcessor(j *job.Job) {
	jobsMutex.Lock()
	j.Status = "ongoing"
	jobsMutex.Unlock()

	var errors []map[string]string

	for _, storeJob := range j.StoreJobs {
		if _, exists := job.StoreMasterData[storeJob.StoreID]; !exists {
			errors = append(errors, map[string]string{
				"store_id": storeJob.StoreID,
				"error":    "",
			})
		}
	}

	err := job.ProcessJob(j.ID, j.StoreJobs)

	jobsMutex.Lock()
	defer jobsMutex.Unlock()

	if len(errors) > 0 || err != nil {
		j.Status = "failed"

		if len(errors) > 0 {
			j.Errors = append(j.Errors, errors...)
		}

	} else {
		j.Status = "completed"
	}
}



func SubmitJobHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var jobRequest JobRequest
	if err := json.NewDecoder(r.Body).Decode(&jobRequest); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	if jobRequest.Count == 0 && len(jobRequest.Visits) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Empty payload"})
		return
	}

	if jobRequest.Count != len(jobRequest.Visits) {
		http.Error(w, "`count` does not match number of `visits`", http.StatusBadRequest)
		return
	}


	jobsMutex.Lock()
	jobCounter++
	jobID := jobCounter
	jobObj := &job.Job{
		ID:        jobID,
		StoreJobs: jobRequest.Visits,
		Status:    "pending",
	}
	jobs[jobID] = jobObj
	jobsMutex.Unlock()

	go jobProcessor(jobObj)

	w.WriteHeader(http.StatusCreated)
	response := map[string]interface{}{
		"job_id": jobID,
		"status": jobObj.Status,
	}

	json.NewEncoder(w).Encode(response)
}
