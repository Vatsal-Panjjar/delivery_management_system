package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/Vatsal-Panjiar/delivery_management_system/internal/cache"
	"github.com/Vatsal-Panjiar/delivery_management_system/internal/handlers"
	"github.com/Vatsal-Panjiar/delivery_management_system/internal/repo"
)

func main() {
	// ----- Connect to PostgreSQL -----
	db, err := sqlx.Connect("postgres", "host=localhost port=5432 user=postgres password=rupupuru@01 dbname=delivery sslmode=disable")
	if err != nil {
		panic(err)
	}

	// ----- Initialize Repositories and Cache -----
	deliveryRepo := repo.NewDeliveryRepo(db)
	userRepo := repo.NewUserRepo(db) // for authentication
	rCache := cache.NewRedisCache("localhost:6379")

	// ----- Initialize Handlers -----
	deliveryHandler := handlers.NewDeliveryHandler(deliveryRepo, rCache)

	jwtSecret := []byte("your_super_secret_key")
	authHandler := handlers.NewAuthHandler(userRepo, jwtSecret)

	// ----- Initialize Router -----
	router := chi.NewRouter()
	handlers.RegisterRoutes(router, deliveryHandler, authHandler)

	// ----- Start Server -----
	fmt.Println("Server running on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		panic(err)
	}
}
