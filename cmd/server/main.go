package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/Vatsal-Panjiar/delivery_management_system/internal/cache"
	"github.com/Vatsal-Panjiar/delivery_management_system/internal/handlers"
	"github.com/Vatsal-Panjiar/delivery_management_system/internal/repo"
)

func main() {
	db, err := sqlx.Connect("postgres", "host=localhost port=5432 user=postgres password=YOUR_PASSWORD dbname=delivery sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	// Repositories
	userRepo := repo.NewUserRepo(db)
	deliveryRepo := repo.NewDeliveryRepo(db)

	// Redis cache
	rCache := cache.NewRedisCache("localhost:6379")

	// Handlers
	authHandler := handlers.NewAuthHandler(userRepo)
	deliveryHandler := handlers.NewDeliveryHandler(deliveryRepo, rCache)

	// Router
	router := chi.NewRouter()
	handlers.RegisterAuthRoutes(router, authHandler)
	handlers.RegisterDeliveryRoutes(router, deliveryHandler)

	fmt.Println("Server running on :8080")
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
