package main

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

func main (){
	randomString := uuid.New().String()

	// Function to print the timestamp and the stored string
	logWithTimestamp := func() {
		timestamp := time.Now().UTC().Format(time.RFC3339)
		fmt.Printf("%s: %s\n", timestamp, randomString)
	}

	// Output immediately on startup
	logWithTimestamp()

	// Create a ticker that ticks every 5 seconds
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	// Loop indefinitely, waiting for the ticker
	for range ticker.C {
		logWithTimestamp()
	}
}