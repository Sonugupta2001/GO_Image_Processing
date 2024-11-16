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

var StoreMasterData = make(map[string]StoreMaster)

func LoadStoreMaster(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open store master file %s: %w", filePath, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read store master file %s: %w", filePath, err)
	}

	if len(records) < 2 {
		return errors.New("store master file is empty or has no data")
	}

	for _, record := range records[1:] {
		if len(record) < 3 {
			return errors.New("invalid record in store master file")
		}

		StoreMasterData[record[2]] = StoreMaster{
			StoreID:   record[2],
			StoreName: record[1],
			AreaCode:  record[0],
		}
	}

	return nil
}
