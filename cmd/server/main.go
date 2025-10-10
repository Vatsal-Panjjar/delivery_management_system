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
	db, _ := sqlx.Connect("postgres", "host=localhost port=5432 user=postgres password=YOUR_PASSWORD dbname=delivery sslmode=disable")

	userRepo := repo.NewUserRepo(db)
	authHandler := handlers.NewAuthHandler(userRepo, []byte("secretkey"))

	deliveryRepo := repo.NewDeliveryRepo(db)
	rCache := cache.NewRedisCache("localhost:6379")
	dHandler := handlers.NewDeliveryHandler(deliveryRepo, rCache)

	router := chi.NewRouter()
	router.Post("/auth/signup", authHandler.Signup)
	router.Post("/auth/login
