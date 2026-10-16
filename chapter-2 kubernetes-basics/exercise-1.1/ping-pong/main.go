package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx := context.Background()

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Println("Not Found Database URL")
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

	// ဇယားမရှိပါက ဖန်တီးခြင်း
	_, err = pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS ping_pongs (
			id INTEGER PRIMARY KEY,
			count INTEGER NOT NULL
		)
	`)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	// Initial row (id = 1, count = 0) ထည့်သွင်းခြင်း
	_, err = pool.Exec(ctx, `
		INSERT INTO ping_pongs (id, count)
		VALUES (1, 0)
		ON CONFLICT (id) DO NOTHING
	`)
	if err != nil {
		log.Fatalf("Failed to initialize row: %v", err)
	}

	// 1. /pingpong Endpoint: Request ရောက်တိုင်း count ကို Database ထဲတွင် တိုးမည်
	http.HandleFunc("/pingpong", func(w http.ResponseWriter, r *http.Request) {
		var count int
		err := pool.QueryRow(ctx, `
			UPDATE ping_pongs
			SET count = count + 1
			WHERE id = 1
			RETURNING count
		`).Scan(&count)
		if err != nil {
			log.Printf("Database error on /pingpong: %v", err)
			http.Error(w, "database error", http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "pong %d", count)
	})

	// 2. /count Endpoint: Database ထဲရှိ လက်ရှိ count တန်ဖိုးကို ထုတ်ပေးမည်
	http.HandleFunc("/count", func(w http.ResponseWriter, r *http.Request) {
		var count int
		err := pool.QueryRow(ctx, `
			SELECT count FROM ping_pongs WHERE id = 1
		`).Scan(&count)
		if err != nil {
			log.Printf("Database error on /count: %v", err)
			http.Error(w, "database error", http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "%d", count)
	})

	log.Println("Ping-pong app listening on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}