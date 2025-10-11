package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"github.com/go-redis/redis/v8"
)

type AdminHandler struct {
	DB  *sqlx.DB
	RDB *redis.Client
}

func RegisterAdminRoutes(r chi.Router, db *sqlx.DB, rdb *redis.Client) {
	h := &AdminHandler{DB: db, RDB: rdb}
	r.Get("/admin/overview", h.AdminOverview)
}

func (h *AdminHandler) AdminOverview(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to Admin Dashboard"))
}
