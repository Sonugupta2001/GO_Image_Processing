package job

import (
	"testing"
)

func TestProcessJob_Success(t *testing.T) {
	StoreMasterData["S00339218"] = StoreMaster{StoreID: "S00339218", StoreName: "Test Store", AreaCode: "123"}

	storeJobs := []StoreJobRequest{
		{
			StoreID:   "S00339218",
			ImageURLs: []string{"https://www.gstatic.com/webp/gallery/2.jpg"},
			VisitTime: "2024-11-15T12:00:00Z",
		},
	}

	err := ProcessJob(1, storeJobs)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestProcessJob_InvalidStoreID(t *testing.T) {
	storeJobs := []StoreJobRequest{
		{
			StoreID:   "INVALID_STORE_ID",
			ImageURLs: []string{"https://www.gstatic.com/webp/gallery/2.jpg"},
			VisitTime: "2024-11-15T12:00:00Z",
		},
	}

	err := ProcessJob(1, storeJobs)
	if err == nil {
		t.Errorf("Expected error for invalid store_id, got nil")
	}
}
