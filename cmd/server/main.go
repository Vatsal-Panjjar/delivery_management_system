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
	db, err := sqlx.Connect("postgres", "host=localhost port=5432 user=postgres password=rupupuru@01 dbname=delivery sslmode=disable")
	if err != nil {
		panic(err)
	}

	userRepo := repo.NewUserRepo(db)
	deliveryRepo := repo.NewDeliveryRepo(db)
	rCache := cache.NewRedisCache("localhost:6379")

	authHandler := handlers.NewAuthHandler(userRepo, []byte("supersecretkey"))
	deliveryHandler := handlers.NewDeliveryHandler(deliveryRepo, rCache)

	router := chi.NewRouter()
	handlers.RegisterRoutes(router, authHandler, deliveryHandler)

	fmt.Println("Server running on :8080")
	http.ListenAndServe(":8080", router)
}
