package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Todo struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
}

func main() {
	ctx := context.Background()

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Println("Not found Database Url");
	}

	var pool *pgxpool.Pool
	var err error

	// App-level Connection Retry Loop (pgxpool ဖြင့် Database ချိတ်ဆက်ရန် စောင့်ဆိုင်းခြင်း)
	for {
		pool, err = pgxpool.New(ctx, databaseURL)
		if err == nil {
			err = pool.Ping(ctx)
			if err == nil {
				log.Println("Connected to PostgreSQL successfully using pgxpool!")
				break
			}
			pool.Close()
		}

		log.Println("waiting for PostgreSQL...")
		time.Sleep(2 * time.Second)
	}
	defer pool.Close()

	// 'todos' table မရှိသေးပါက အလိုအလျောက် ဆောက်ပေးမည့် Query
	_, err = pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS todos (
			id SERIAL PRIMARY KEY,
			text TEXT NOT NULL
		)
	`)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	http.HandleFunc("/todos", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// GET Method: Database မှ Todo များကို SELECT လုပ်၍ ပြန်ထုတ်ပေးခြင်း
		if r.Method == http.MethodGet {
			rows, err := pool.Query(ctx, "SELECT id, text FROM todos ORDER BY id ASC")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer rows.Close()

			todos := []Todo{}
			for rows.Next() {
				var t Todo
				if err := rows.Scan(&t.ID, &t.Text); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				todos = append(todos, t)
			}

			json.NewEncoder(w).Encode(todos)
			return
		}

		// POST Method: Todo အသစ်ကို Database ထဲသို့ INSERT လုပ်၍ သိမ်းဆည်းခြင်း
		if r.Method == http.MethodPost {
			var newTodo Todo
			if err := json.NewDecoder(r.Body).Decode(&newTodo); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			query := `INSERT INTO todos (text) VALUES ($1) RETURNING id`
			err := pool.QueryRow(ctx, query, newTodo.Text).Scan(&newTodo.ID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(newTodo)
			return
		}
	})

	log.Println("Todo Backend listening on port 5000...")
	if err := http.ListenAndServe(":5000", nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}