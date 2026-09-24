package main

import (
	
	"log"
	"net/http"
	"os"
	"fmt"
	"strings"
	"io"
)

const logFilePath = "/usr/src/app/files/log.txt";

func statusHandler(w http.ResponseWriter, r *http.Request) {
	// Writer ရေးနေတဲ့ Shared File ကို လှမ်းဖတ်မယ်
	
	// ၁။ Log File (Writer က ရေးထားသည်) ကို ဖတ်ပါ
	logData, err := os.ReadFile(logFilePath)
	if err != nil {
		// Log ဖိုင် မရှိသေးပါက Error ပြမည်
		http.Error(w, "Log file not generated yet...", http.StatusNotFound)
		return
	}
	// ၂။ Ping-pong Count ဖိုင်ကို ဖတ်
	pingCount := "0" // Ping-pong ဖိုင် မရှိသေးပါက (သို့) မဖတ်နိုင်ပါက Default '0' ဟု သတ်မှတ်မည်
	
	// Environment variable မရှိပါက Default 'http://pingpong-svc:8080/count' ကို သုံးမည်
	pingPongURL := os.Getenv("PINGPONG_URL")
	if pingPongURL == "" {
		pingPongURL = "http://pingpong-svc:8080/count"
	}
	resp, err := http.Get(pingPongURL)
	if err == nil {
		defer resp.Body.Close()
		body, readErr := io.ReadAll(resp.Body)
		if readErr == nil {
			pingCount = strings.TrimSpace(string(body))
		}
	} else {
		// Ping Pong service ကို လှမ်းခေါ်လို့ မရပါက Log ထဲ အမှားရိုက်ပြမည် (Default '0' ပဲ ပြထားမည်)
		log.Printf("Error fetching ping count from %s: %v", pingPongURL, err)
	}

	responseText := fmt.Sprintf("%s.\nPing / Pongs: %s\n", strings.TrimSpace(string(logData)), pingCount)
	// ဖတ်လို့ရတဲ့ စာကို Browser / curl ဆီ ပြန်ပို့မယ်
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(responseText))
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	http.HandleFunc("/", statusHandler)

	log.Printf("Reader Server is listening on port %s...", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}