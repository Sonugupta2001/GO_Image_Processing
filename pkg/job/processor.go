package job

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"sync"
	"time"
)

type Job struct {
	ID        int64
	StoreJobs []StoreJobRequest
	Status    string
	Errors    []string
}

type StoreJobRequest struct {
	StoreID   string   `json:"store_id"`
	ImageURLs []string `json:"image_url"`
	VisitTime string   `json:"visit_time"`
}



func ProcessJob(jobID int64, storeJobs []StoreJobRequest) error {
	var wg sync.WaitGroup
	errChan := make(chan map[string]string, len(storeJobs))

	// Process each store job
	for _, storeJob := range storeJobs {
		// Check if the store exists in the StoreMaster data
		if _, exists := StoreMasterData[storeJob.StoreID]; !exists {
			// Send error if store does not exist
			errChan <- map[string]string{
				"store_id": storeJob.StoreID,
				"error":    fmt.Sprintf("store_id %s does not exist", storeJob.StoreID),
			}
			continue
		}

		// Process images for the store
		for _, url := range storeJob.ImageURLs {
			wg.Add(1)
			go func(imgURL string) {
				defer wg.Done()
				// Process the image, if it fails, send an error
				if err := processImage(imgURL); err != nil {
					errChan <- map[string]string{
						"store_id": storeJob.StoreID,
						"error":    fmt.Sprintf("error processing image %s: %v", imgURL, err),
					}
				}
			}(url)
		}
	}

	wg.Wait()
	close(errChan)

	// Collect errors related to stores
	var errors []map[string]string
	for err := range errChan {
		errors = append(errors, err)
	}

	

	// If there were errors, return them
	if len(errors) > 0 {
		return fmt.Errorf("Job completed with errors: %v", errors)
	}

	return nil
}

func processImage(url string) error {
	client := &http.Client{Timeout: 1000 * time.Second}

	resp, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download image: %w", err)
	}
	defer resp.Body.Close()

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to decode image: %w", err)
	}

	width := img.Bounds().Dx()
	height := img.Bounds().Dy()
	perimeter := 2 * (width + height)

	fmt.Printf("Image processed: Perimeter = %d\n", perimeter)
	return nil
}
