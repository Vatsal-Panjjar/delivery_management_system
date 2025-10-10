package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/Vatsal-Panjiar/delivery_management_system/internal/cache"
	"github.com/Vatsal-Panjiar/delivery_management_system/internal/handlers"
	"github.com/Vatsal-Panjiar/delivery_management_system/internal/repo"
)

func main() {
	// Use environment variables, fallback to defaults
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "rupupuru@01")
	dbName := getEnv("DB_NAME", "delivery")

	// Connect to Postgres
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		panic(fmt.Errorf("failed to connect to Postgres: %w", err))
	}

	// Connect to Redis
	redisAddr := getEnv("REDIS_ADDR", "localhost:6379")
	rCache := cache.NewRedisCache(redisAddr)
	ctx := context.Background()
	if err := rCache.Ping(ctx); err != nil {
		panic(fmt.Errorf("failed to connect to Redis: %w", err))
	}

	// Initialize repository and handlers
	repo := repo.NewDeliveryRepo(db)
	handler := handlers.NewDeliveryHandler(repo, rCache)

	// Setup router and register routes
	router := chi.NewRouter()
	RegisterRoutes(router, handler) // We'll define this below

	fmt.Println("Server running on :8080")
	http.ListenAndServe(":8080", router)
}

// getEnv fetches environment variable or returns fallback
func getEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}

// RegisterRoutes sets up the HTTP routes for deliveries
func RegisterRoutes(r *chi.Mux, h *handlers.DeliveryHandler) {
	r.Route("/deliveries", func(r chi.Router) {
		r.Post("/", h.Create)
		r.Get("/", h.ListByStatus)
		r.Get("/{id}", h.Get)
	})
}
