package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/Vatsal-Panjiar/delivery_management_system/internal/cache"
	"github.com/Vatsal-Panjiar/delivery_management_system/internal/handlers"
	"github.com/Vatsal-Panjiar/delivery_management_system/internal/repo"
)

func main() {
	// --- PostgreSQL Connection ---
	db, err := sqlx.Connect("postgres", "host=localhost port=5432 user=postgres password=YOUR_PASSWORD dbname=delivery sslmode=disable")
	if err != nil {
		log.Fatalf("Postgres connection error: %v", err)
	}

	// --- Redis Connection ---
	rCache := cache.NewRedisCache("localhost:6379")

	// --- Repos ---
	deliveryRepo := repo.NewDeliveryRepo(db)
	userRepo := repo.NewUserRepo(db)

	// --- Handlers ---
	deliveryHandler := handlers.NewDeliveryHandler(deliveryRepo, rCache)
	authHandler := handlers.NewAuthHandler(userRepo, []byte("mysecretkey123")) // JWT secret

	// --- Router ---
	router := chi.NewRouter()

	// Auth Routes
	router.Post("/register", authHandler.Register)
	router.Post("/login", authHandler.Login)

	// Delivery Routes
	router.Post("/deliveries", deliveryHandler.Create)
	router.Get("/deliveries", deliver
