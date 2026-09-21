package main

import (
	
	"log"
	"net/http"
	"os"
	"fmt"
	"strings"
)

const pingpongFilePath = "/usr/src/app/files/pingpong.txt";
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
	pingData, err := os.ReadFile(pingpongFilePath)
	if err == nil {
		pingCount = strings.TrimSpace(string(pingData))
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