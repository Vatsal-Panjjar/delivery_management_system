package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/Vatsal-Panjjar/delivery_management_system/internal/handlers"
	"github.com/Vatsal-Panjjar/delivery_management_system/internal/middleware"
)

func main() {
	fmt.Println("Starting Delivery Management Server...")

	// --- DATABASE & REDIS SETUP ---
	// Replace these with your actual credentials
	dbURL := "postgres://postgres:rupupuru@01@localhost:5432/delivery_db?sslmode=disable"
	redisAddr := "localhost:6379"

	// Connect to Postgres
	db, err := sqlx.Connect("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to Postgres: %v", err)
	}
	defer db.Close()
	fmt.Println("Connected to Postgres")

	// Connect to Redis
	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	_, err = rdb.Ping(rdb.Context()).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	fmt.Println("Connected to Redis")

	// --- ROUTER ---
	r := chi.NewRouter()

	// Root status page
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Delivery Management System API is running"))
	})

	// --- API ROUTES ---
	handlers.RegisterOrderRoutes(r, db, rdb)
	// Add admin or user routes here as needed

	// --- SERVE STATIC HTML ---
	fs := http.FileServer(http.Dir("./web"))
	r.Handle("/*", http.StripPrefix("/", fs))

	// --- START SERVER ---
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("Server running on http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
