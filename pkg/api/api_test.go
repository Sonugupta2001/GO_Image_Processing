package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestSubmitJobHandler_Success(t *testing.T) {
	payload := `{
		"count": 1,
		"visits": [
			{
				"store_id": "S00339218",
				"image_url": ["https://www.gstatic.com/webp/gallery/2.jpg"],
				"visit_time": "2024-11-15T12:00:00Z"
			}
		]
	}`

	req, err := http.NewRequest("POST", "/api/submit/", bytes.NewReader([]byte(payload)))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(SubmitJobHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, status)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}

	if _, ok := response["job_id"]; !ok {
		t.Error("Response missing job_id")
	}
}

func TestSubmitJobHandler_InvalidPayload(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/submit/", bytes.NewReader([]byte(`{}`)))
	w := httptest.NewRecorder()

	SubmitJobHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestSubmitJobHandler_EmptyPayload(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/submit/", bytes.NewReader([]byte(`{}`)))
	w := httptest.NewRecorder()

	SubmitJobHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}

	if response["error"] != "Empty payload" {
		t.Errorf("Expected error message 'Empty payload', got %v", response["error"])
	}
}

func TestGetJobStatusHandler_Success(t *testing.T) {
	payload := `{
		"count": 1,
		"visits": [
			{
				"store_id": "S00339218",
				"image_url": ["https://www.gstatic.com/webp/gallery/2.jpg"],
				"visit_time": "2024-11-15T12:00:00Z"
			}
		]
	}`

	req, err := http.NewRequest("POST", "/api/submit/", bytes.NewReader([]byte(payload)))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(SubmitJobHandler)
	handler.ServeHTTP(rr, req)

	var response map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}

	jobID := response["job_id"].(float64)
	jobIDStr := strconv.Itoa(int(jobID))

	reqStatus, err := http.NewRequest("GET", "/api/status?jobid="+jobIDStr, nil)
	if err != nil {
		t.Fatal(err)
	}

	rrStatus := httptest.NewRecorder()
	statusHandler := http.HandlerFunc(GetJobStatusHandler)

	statusHandler.ServeHTTP(rrStatus, reqStatus)

	if status := rrStatus.Code; status != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, status)
	}

	var statusResponse map[string]interface{}
	if err := json.NewDecoder(rrStatus.Body).Decode(&statusResponse); err != nil {
		t.Fatal(err)
	}

	if statusResponse["job_id"] != jobID {
		t.Errorf("Expected job_id %v, got %v", jobID, statusResponse["job_id"])
	}
}

func TestGetJobStatusHandler_Failure(t *testing.T) {
	reqStatus, err := http.NewRequest("GET", "/api/status?jobid=999", nil)
	if err != nil {
		t.Fatal(err)
	}

	rrStatus := httptest.NewRecorder()
	statusHandler := http.HandlerFunc(GetJobStatusHandler)

	statusHandler.ServeHTTP(rrStatus, reqStatus)

	if status := rrStatus.Code; status != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, status)
	}

	var statusResponse map[string]interface{}
	if err := json.NewDecoder(rrStatus.Body).Decode(&statusResponse); err != nil {
		t.Fatal(err)
	}

	if len(statusResponse) != 0 {
		t.Errorf("Expected empty response, got %v", statusResponse)
	}
}
