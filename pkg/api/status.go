package api

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func GetJobStatusHandler(w http.ResponseWriter, r *http.Request) {
	jobIDStr := r.URL.Query().Get("jobid")
	var jobID int64
	var err error
	jobID, err = strconv.ParseInt(jobIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid job ID", http.StatusBadRequest)
		return
	}

	// Locking to check the job status
	jobsMutex.Lock()
	j, exists := jobs[jobID]
	jobsMutex.Unlock()

	if !exists {
		http.Error(w, "Job ID not found", http.StatusNotFound)
		return
	}

	// Prepare base response
	response := map[string]interface{}{
		"status": j.Status,
		"job_id": strconv.FormatInt(jobID, 10),
	}

	// Add errors only if the job failed
	if j.Status == "failed" {
		var errorDetails []map[string]string
		for _, errMsg := range j.Errors {
			errorDetails = append(errorDetails, map[string]string{
				"store_id": errMsg["store_id"],
				"error":    "", // Keep error message empty as per requirements
			})
		}
		response["error"] = errorDetails
	}

	// Send JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}