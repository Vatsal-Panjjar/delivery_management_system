package handlers

import (
    "net/http"

    "github.com/go-chi/chi/v5"
    "github.com/jmoiron/sqlx"
)

func RegisterAdminRoutes(r *chi.Mux, db *sqlx.DB) {
    r.Route("/admin", func(r chi.Router) {
        r.Get("/dashboard", AdminDashboardHandler(db))
        r.Put("/orders/{id}/status", UpdateOrderStatusHandler(db))
    })
}

func AdminDashboardHandler(db *sqlx.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Admin Dashboard"))
    }
}

func UpdateOrderStatusHandler(db *sqlx.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Update Order Status"))
    }
}
