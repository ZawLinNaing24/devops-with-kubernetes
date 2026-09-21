package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"path/filepath"
)

// ၁။ Package-level မှာ Variable တွေ ကြေညာခြင်း
var (
	counter int  = 0
	dir string  = "/usr/src/app/files"
	filePath string
)

// ၂။ init() function ထဲမှာ ဖိုလ်ဒါဆောက်ခြင်းနဲ့ Path သတ်မှတ်ခြင်း
func init() {
	os.MkdirAll(dir, os.ModePerm)
	filePath = filepath.Join(dir, "pingpong.txt")
}

func pingpongHandler(w http.ResponseWriter, r *http.Request) {
	data, err := os.ReadFile(filePath)
	if err == nil {
		// ဖိုင်ထဲမှ စာသားကို Integer ဂဏန်းသို့ ပြောင်းခြင်း
		savedCount, convErr := strconv.Atoi(strings.TrimSpace(string(data)))
		if convErr == nil {
			counter = savedCount
		}
	}
	// Request လာတိုင်း pong <counter> ဟု တုံ့ပြန်ပြီး counter ကို ၁ တိုးမည်
	response := fmt.Sprintf("pong %d\n", counter)
	counter++
	err = os.WriteFile(filePath, []byte(strconv.Itoa(counter)), 0644)
	if err != nil {
		fmt.Printf("Error writing file: %v\n", err)
	}
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