package main

import (
	"fmt"
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
	// ----- Postgres -----
	db, err := sqlx.Connect("postgres", "host=localhost port=5432 user=postgres password=your_password dbname=delivery sslmode=disable")
	if err != nil {
		panic(err)
	}

	// ----- Repositories and Cache -----
	rRepo := repo.NewDeliveryRepo(db)
	rCache := cache.NewRedisCache("localhost:6379")
	deliveryHandler := handlers.NewDeliveryHandler(rRepo, rCache)
	authHandler := handlers.NewAuthHandler(db) // Handles signup/login

	// ----- Router -----
	router := chi.NewRouter()

	// Serve frontend static files
	fs := http.FileServer(http.Dir("./frontend"))
	router.Handle("/*", fs)

	// Register API routes
	handlers.RegisterRoutes(router, deliveryHandler, authHandler)

	fmt.Println("Server running on :8080")
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
}
