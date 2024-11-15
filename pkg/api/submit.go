package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"retailpulse-image-service/pkg/job"
	"sync"
)

// Job state
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
		// Validate store ID directly from the global `storeMaster` map
		if _, exists := job.StoreMasterData[storeJob.StoreID]; !exists {
			errors = append(errors, map[string]string{
				"store_id": storeJob.StoreID,
				"error":    fmt.Sprintf("store_id %s does not exist", storeJob.StoreID),
			})
			continue
		}
	}

	// Process the job after store validation
	err := job.ProcessJob(j.ID, j.StoreJobs)

	jobsMutex.Lock()
	defer jobsMutex.Unlock()

	if len(errors) > 0 || err != nil {
		j.Status = "failed"

		if err != nil {
			errors = append(errors, map[string]string{
				"store_id": "Unknown",
				"error":    err.Error(),
			})
		}

		// Add the errors to the job
		for _, e := range errors {
			j.Errors = append(j.Errors, fmt.Sprintf("store_id %s: %s", e["store_id"], e["error"]))
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

	// Check if the payload is empty
	if jobRequest.Count == 0 && len(jobRequest.Visits) == 0 {
		w.Header().Set("Content-Type", "application/json") // Ensure JSON response
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Empty payload"})
		return
	}

	// Validate top-level fields
	if jobRequest.Count != len(jobRequest.Visits) {
		http.Error(w, "`count` does not match number of `visits`", http.StatusBadRequest)
		return
	}

	// Create the job
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

	// Process the job asynchronously
	go jobProcessor(jobObj)

	w.WriteHeader(http.StatusCreated)
	response := map[string]interface{}{
		"job_id": jobID,
		"status": jobObj.Status,
	}

	// If the job failed, include errors in the response
	if jobObj.Status == "failed" {
		response["error"] = jobObj.Errors
	}

	json.NewEncoder(w).Encode(response)
}
