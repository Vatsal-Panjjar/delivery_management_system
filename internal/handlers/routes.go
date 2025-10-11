package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"github.com/go-redis/redis/v8"
)

func RegisterRoutes(r chi.Router, db *sqlx.DB, rdb *redis.Client) {
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Delivery Management System API is running"))
	})
	RegisterAuthRoutes(r, db, rdb)
	RegisterOrderRoutes(r, db, rdb)
	RegisterAdminRoutes(r, db, rdb)
}
