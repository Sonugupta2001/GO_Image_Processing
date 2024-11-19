package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"retailpulse-image-service/pkg/job"
	"net/http"
	"retailpulse-image-service/pkg/api"
	"time"
)

func main() {
	err := job.LoadStoreMaster("./StoreMaster.csv")
	if err != nil {
		log.Fatalf("Error loading store master file: %v", err)
	}

	http.HandleFunc("/api/submit/", api.SubmitJobHandler)
	http.HandleFunc("/api/status", api.GetJobStatusHandler)
	http.HandleFunc("/health", healthCheckHandler)

	server := &http.Server{Addr: ":8080"}

	go func() {
		log.Println("Server running on http://localhost:8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	log.Println("Shutting down gracefully...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped")
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}