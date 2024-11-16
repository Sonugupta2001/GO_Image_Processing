package job
import (
	"os"
	"testing"
)

func TestLoadStoreMaster_Success(t *testing.T) {
	mockCSV := "store_id,store_name,area_code\nS00339218,Test Store,123\n"
	file, _ := os.CreateTemp("", "store_master.csv")
	defer os.Remove(file.Name())
	file.WriteString(mockCSV)
	file.Close()

	err := LoadStoreMaster(file.Name())
	if err != nil {
		t.Fatalf("Failed to load store master: %v", err)
	}

	if len(StoreMasterData) != 1 {
		t.Errorf("Expected 1 store, got %d", len(StoreMasterData))
	}

	if StoreMasterData["S00339218"].StoreName != "Test Store" {
		t.Errorf("Store name mismatch: expected 'Test Store', got '%s'", StoreMasterData["S00339218"].StoreName)
	}
}

func TestLoadStoreMaster_InvalidFile(t *testing.T) {
	err := LoadStoreMaster("non_existent_file.csv")
	if err == nil {
		t.Errorf("Expected error for invalid file, got nil")
	}
}
