package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/Vatsal-Panjiar/delivery_management_system/internal/auth"
	"github.com/Vatsal-Panjiar/delivery_management_system/internal/cache"
	"github.com/Vatsal-Panjiar/delivery_management_system/internal/handlers"
	"github.com/Vatsal-Panjiar/delivery_management_system/internal/repo"
)

func main() {
	db, err := sqlx.Connect("postgres", "host=localhost port=5432 user=postgres password=YOUR_PASSWORD dbname=delivery sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	rRepo := repo.NewDeliveryRepo(db)
	rCache := cache.NewRedisCache("localhost:6379")

	authHandler := auth.NewAuthHandler(rRepo, []byte("supersecretkey"))
	deliveryHandler := handlers.NewDeliveryHandler(rRepo, rCache)

	r := chi.NewRouter()
	auth.RegisterAuthRoutes(r, authHandler)
	handlers.RegisterDeliveryRoutes(r, deliveryHandler)

	fmt.Println("Server running on :8080")
	http.ListenAndServe(":8080", r)
}
