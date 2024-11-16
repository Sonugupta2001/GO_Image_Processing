package job

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
)

type StoreMaster struct {
	StoreID   string
	StoreName string
	AreaCode  string
}

// Global variable to hold store data.
var StoreMasterData = make(map[string]StoreMaster)

// LoadStoreMaster loads store master data from a CSV file.
func LoadStoreMaster(filePath string) error {
	// Open the file
	file, err := os.Open(filePath) // Correct file path
	if err != nil {
		return fmt.Errorf("failed to open store master file %s: %w", filePath, err)
	}
	defer file.Close()

	// Read the CSV data
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read store master file %s: %w", filePath, err)
	}

	// Ensure there is at least one record (the header)
	if len(records) < 2 {
		return errors.New("store master file is empty or has no data")
	}

	// Skip the header and load the rest of the records
	for _, record := range records[1:] { // Skipping header
		if len(record) < 3 {
			return errors.New("invalid record in store master file")
		}

		// Adjusted order to match AreaCode, StoreName, StoreID
		StoreMasterData[record[2]] = StoreMaster{
			StoreID:   record[2], // Third column is StoreID
			StoreName: record[1], // Second column is StoreName
			AreaCode:  record[0], // First column is AreaCode
		}
	}

	return nil
}
