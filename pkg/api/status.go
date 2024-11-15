package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

/*
func GetJobStatusHandler(w http.ResponseWriter, r *http.Request) {
	jobIDStr := r.URL.Query().Get("jobid")
	var jobID int64
	var err error
	jobID, err = strconv.ParseInt(jobIDStr, 10, 64)
	if err != nil {
		fmt.Println("Invalid job ID", err)
		http.Error(w, "Invalid job ID", http.StatusBadRequest)
		return
	}

	jobsMutex.Lock()
	j, exists := jobs[jobID]
	jobsMutex.Unlock()

	if !exists {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{})
		return
	}

	// Prepare the response map
	response := map[string]interface{}{
		"status": j.Status,
		"job_id": jobID,
	}

	// If job failed, return errors
	if j.Status == "failed" {
		// Create an error array with detailed store errors
		var errorDetails []map[string]string
		for _, errMsg := range j.Errors {
			// Assuming the error contains store_id details, adjust as needed
			errorDetails = append(errorDetails, map[string]string{
				"store_id": "Unknown", // Retrieve store ID details as per your logic
				"error":    errMsg,
			})
		}
		response["error"] = errorDetails
	}

	// Return the final response
	json.NewEncoder(w).Encode(response)
} */

func GetJobStatusHandler(w http.ResponseWriter, r *http.Request) {
	jobIDStr := r.URL.Query().Get("jobid")
	var jobID int64
	var err error
	jobID, err = strconv.ParseInt(jobIDStr, 10, 64)
	if err != nil {
		fmt.Println("Invalid job ID", err)
		http.Error(w, "Invalid job ID", http.StatusBadRequest)
		return
	}

	// Locking to check the job status
	jobsMutex.Lock()
	j, exists := jobs[jobID]
	jobsMutex.Unlock()

	// Check if the job exists
	if !exists {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{})
		return
	}

	// Prepare response
	response := map[string]interface{}{
		"status": j.Status,
		"job_id": jobID,
	}

	// If job failed, return errors specific to each store
	if j.Status == "failed" {
		var errorDetails []map[string]string
		// Collect each individual error
		for _, errMsg := range j.Errors {
			// Assuming that the error message is already in the correct format
			errorDetails = append(errorDetails, map[string]string{
				"store_id": "Unknown", // You can also extract the store_id from the error message if needed
				"error":    errMsg,
			})
		}
		response["error"] = errorDetails
	}

	// Final JSON response
	json.NewEncoder(w).Encode(response)
}
