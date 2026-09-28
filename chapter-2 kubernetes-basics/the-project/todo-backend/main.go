package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
)

type Todo struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
}

var (
	todos = []Todo{{ID: 1, Text: "Todo Item အဟောင်း"}}
	mu    sync.Mutex
)

func main() {
	http.HandleFunc("/todos", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// CORS (Frontend က တိုက်ရိုက် ခေါ်ပါက လိုအပ်သည်)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		if r.Method == http.MethodGet {
			mu.Lock()
			defer mu.Unlock()
			json.NewEncoder(w).Encode(todos)
			return
		}

		if r.Method == http.MethodPost {
			var newTodo Todo
			if err := json.NewDecoder(r.Body).Decode(&newTodo); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			mu.Lock()
			newTodo.ID = len(todos) + 1
			todos = append(todos, newTodo)
			mu.Unlock()

			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(newTodo)
			return
		}
	})

	log.Println("Todo Backend listening on port 5000...")
	http.ListenAndServe(":5000", nil)
}