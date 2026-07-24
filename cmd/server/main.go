package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joetomjob/user-api/internal/user"
	"github.com/joho/godotenv"
)

type LogEvent struct {
	Method   string
	Path     string
	Duration int
}

func LogWorker(ch <-chan LogEvent, wg *sync.WaitGroup) {
	defer wg.Done()

	for e := range ch {
		fmt.Printf("Method: %s, Path: %s, Duration: %d", e.Method, e.Path, e.Duration)
		fmt.Println()
	}
}

func middleware(next http.Handler, ch chan<- LogEvent) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		select {
		case ch <- LogEvent{
			Method:   r.Method,
			Path:     r.URL.Path,
			Duration: int(time.Since(start).Milliseconds()),
		}:
		default:
		}
	})
}

func health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]string{"status": "ok"}

	json.NewEncoder(w).Encode(response)
}

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	ch := make(chan LogEvent, 100)

	go LogWorker(ch, &wg)

	// load env file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("env variable fetch failed with error: %v", err)
	}

	// db connection
	url := os.Getenv("DATABASE_URL")
	if strings.TrimSpace(url) == "" {
		log.Fatal("connection string cannot be empty")
	}

	config, err := pgxpool.ParseConfig(url)
	if err != nil {
		log.Fatalf("cannot parse database URL: %v", err)
	}

	config.MaxConnIdleTime = time.Minute * 30
	config.MaxConns = 10
	config.MinConns = 2

	ctx := context.Background()
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Fatalf("cannot create new connection: %v", err)
	}
	defer pool.Close()

	if err = pool.Ping(ctx); err != nil {
		log.Fatalf("cannot ping the db: %v", err)
	}

	fmt.Println("succesfully connected to db")

	repo := user.NewRepo(pool)
	service := user.NewService(repo)
	handler := user.NewHandler(service)
	// router
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", health)
	mux.HandleFunc("POST /users", handler.Create)
	mux.HandleFunc("GET /users/{id}", handler.GetById)
	mux.HandleFunc("PUT /users/{id}", handler.Update)
	mux.HandleFunc("DELETE /users/{id}", handler.Delete)

	// server
	server := &http.Server{
		Addr:         ":8080",
		Handler:      middleware(mux, ch),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed %s", err)
	}

	close(ch)
	wg.Wait()
}
