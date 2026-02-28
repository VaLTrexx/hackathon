package main

import (
	"log"
	"net/http"
	"sync"

	"yourproject/internal/api"
	"yourproject/internal/models"
	"yourproject/internal/scheduler"
	"yourproject/internal/worker"
)

func main() {

	// Channels
	jobs := make(chan models.Job, 10)
	results := make(chan models.Result, 10)

	// In-memory storage
	store := make(map[string]models.Result)
	var mu sync.RWMutex

	// Start worker pool
	workerCount := 3
	for i := 0; i < workerCount; i++ {
		go worker.StartWorker(i, jobs, results)
	}

	// Start results aggregator
	go func() {
		for res := range results {
			mu.Lock()
			store[res.Country] = res
			mu.Unlock()
		}
	}()

	// Start scheduler
	go scheduler.StartScheduler(jobs)

	// Setup HTTP handlers
	http.HandleFunc("/risk", api.GetAllRisk(&store, &mu))
	http.HandleFunc("/risk/", api.GetRiskByCountry(&store, &mu))

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
