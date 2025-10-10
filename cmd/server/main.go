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
	// Connect to PostgreSQL
	db, err := sqlx.Connect("postgres", "host=localhost port=5432 user=postgres password=rupupuru@01 dbname=delivery sslmode=disable")
	if err != nil {
		log.Fatalf("DB connection failed: %v", err)
	}

	// Initialize repositories
	userRepo := repo.NewUserRepo(db)
	deliveryRepo := repo.NewDeliveryRepo(db)

	// Initialize Redis cache
	rCache := cache.NewRedisCache("localhost:6379")

	// Initialize handlers
	jwtSecret := []byte("your_super_secret_key")
	authHandler := handlers.NewAuthHandler(userRepo, jwtSecret)
	deliveryHandler := handlers.NewDeliveryHandler(deliveryRepo, rCache)

	// Create router
	router := chi.NewRouter()

	// Auth routes
	router.Post("/signup", authHandler.Signup)
	router.Post("/login", authHandler.Login)

	// Delivery routes
	router.Post("/deliveries", deliveryHandler.Create)
	router.Get("/deliveries/{id}", deliveryHandler.Get)
	router.Get("/deliveries", deliveryHandler.ListByStatus)

	// Serve static frontend files
	router.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("frontend/static"))))
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "frontend/index.html")
	})

	fmt.Println("Server running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
