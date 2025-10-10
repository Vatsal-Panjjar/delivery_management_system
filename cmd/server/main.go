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
	// Connect to Postgres
	db, err := sqlx.Connect("postgres", "postgres://user:pass@localhost:5432/delivery?sslmode=disable")
	if err != nil {
		panic(err)
	}

	r := repo.NewDeliveryRepo(db)
	c := cache.NewRedisCache("localhost:6379")
	h := handlers.NewDeliveryHandler(r, c)

	router := chi.NewRouter()
	handlers.RegisterRoutes(router, h)

	fmt.Println("Server running on :8080")
	http.ListenAndServe(":8080", router)
}
