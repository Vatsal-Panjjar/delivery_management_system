package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/Vatsal-Panjjar/delivery_management_system/internal/handlers"
)

func main() {
	fmt.Println("Starting Delivery Management Server...")

	// Hardcoded PostgreSQL URL (replace username/password/dbname if needed)
	dbURL := "postgres://postgres:rurupuru@01@localhost:5432/delivery_db?sslmode=disable"
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

	// Create router
	r := chi.NewRouter()

	// Simple root route
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Delivery Management System API is running"))
	})

	// Register routes
	handlers.RegisterAuthRoutes(r, db)             // registration & login
	handlers.RegisterUserRoutes(r, db, rdb)       // order management

	// Start server
	port := "8080"
	fmt.Printf("Server running on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
