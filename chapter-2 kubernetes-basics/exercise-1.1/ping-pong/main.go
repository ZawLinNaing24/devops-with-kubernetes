package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

var counter int = 0

func pingpongHandler(w http.ResponseWriter, r *http.Request) {
	// Request လာတိုင်း pong <counter> ဟု တုံ့ပြန်ပြီး counter ကို ၁ တိုးမည်
	response := fmt.Sprintf("pong %d\n", counter)
	counter++

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	// /pingpong path Endpoint
	http.HandleFunc("/pingpong", pingpongHandler)

	log.Printf("Ping-pong server listening on port %s...", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}