package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/Vatsal-Panjjar/delivery_management_system/internal/cache"
	"github.com/Vatsal-Panjjar/delivery_management_system/internal/db"
	"github.com/Vatsal-Panjjar/delivery_management_system/internal/handlers"
	"github.com/Vatsal-Panjjar/delivery_management_system/internal/workers"
)

func main() {
	pgURL := os.Getenv("POSTGRES_URL")
	if pgURL == "" {
		pgURL = "postgres://postgres:rupupuru@01@localhost:5432/delivery_db?sslmode=disable"
	}
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	sqlxDB, err := sqlx.Connect("postgres", pgURL)
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}
	defer sqlxDB.Close()
	fmt.Println("Connected to Postgres")

	rdb := redis.NewClient(&redis.Options{Addr: redisAddr})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("failed to connect redis: %v", err)
	}
	fmt.Println("Connected to Redis")

	store := db.New(sqlxDB)
	cacheClient := cache.New(rdb)

	// start async tracker
	tracker := workers.NewOrderTracker(store, cacheClient)
	tracker.Start()
	defer tracker.Stop()

	r := chi.NewRouter()

	// Serve static web files (html/css/js)
	fs := http.FileServer(http.Dir("./web"))
	r.Handle("/*", http.StripPrefix("/", fs))

	// API mount
	api := chi.NewRouter()
	handlers.RegisterAuthHandlers(api, store, cacheClient)
	handlers.RegisterOrderHandlers(api, store, cacheClient, tracker)
	handlers.RegisterAdminHandlers(api, store, cacheClient, tracker)
	r.Mount("/api", api)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("Server running on http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
