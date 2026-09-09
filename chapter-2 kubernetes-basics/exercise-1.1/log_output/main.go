package main

import (
	// "crypto/rand"
	// "encoding/hex"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
)
var randomString string

func statusHandler(w http.ResponseWriter, r *http.Request) {
	
	currentTimestamp := time.Now().Format(time.RFC3339)
	
	// Format: <timestamp>: <randomString>
	response := fmt.Sprintf("%s: %s\n", currentTimestamp, randomString)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}

func main (){

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	randomString = uuid.New().String()

	// Function to print the timestamp and the stored string
	logWithTimestamp := func() {
		timestamp := time.Now().UTC().Format(time.RFC3339)
		fmt.Printf("%s: %s\n", timestamp, randomString)
	}

	// Output immediately on startup
	logWithTimestamp()

	// Create a ticker that ticks every 5 seconds
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			logWithTimestamp()
		}
	}()
	
	http.HandleFunc("/", statusHandler)

	log.Printf("Server is listening on port %s...", port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}