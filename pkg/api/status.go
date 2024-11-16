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

	jobsMutex.Lock()
	j, exists := jobs[jobID]
	jobsMutex.Unlock()

	if !exists {
		http.Error(w, "Job ID not found", http.StatusNotFound)
		return
	}

	response := map[string]interface{}{
		"status": j.Status,
		"job_id": strconv.FormatInt(jobID, 10),
	}

	if j.Status == "failed" {
		var errorDetails []map[string]string
		for _, errMsg := range j.Errors {
			errorDetails = append(errorDetails, map[string]string{
				"store_id": errMsg["store_id"],
				"error":    "",
			})
		}
		response["error"] = errorDetails
	}


	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}