package main

import (
	
	"log"
	"net/http"
	"os"
)

func statusHandler(w http.ResponseWriter, r *http.Request) {
	// Writer ရေးနေတဲ့ Shared File ကို လှမ်းဖတ်မယ်
	filePath := "/usr/src/app/files/log.txt"
	data, err := os.ReadFile(filePath)
	
	if err != nil {
		// ဖိုင် မရှိသေးရင် (သို့) ဖတ်လို့ မရသေးရင် ပြမယ့်စာ
		http.Error(w, "Log file not generated yet...", http.StatusNotFound)
		return
	}

	// ဖတ်လို့ရတဲ့ စာကို Browser / curl ဆီ ပြန်ပို့မယ်
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
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