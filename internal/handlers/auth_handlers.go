package handlers

import (
    "encoding/json"
    "net/http"
    "time"

    "github.com/go-chi/chi/v5"
    "github.com/jmoiron/sqlx"
    "github.com/Vatsal-Panjjar/delivery_management_system/internal/auth"
)

type AuthHandler struct {
    DB *sqlx.DB
}

func RegisterAuthRoutes(r *chi.Mux, db *sqlx.DB) {
    h := &AuthHandler{DB: db}
    r.Post("/register", h.Register)
    r.Post("/login", h.Login)
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
    var req struct{ Username, Password string }
    json.NewDecoder(r.Body).Decode(&req)
    hashed := auth.HashPassword(req.Password)
    _, err := h.DB.Exec("INSERT INTO users (username, password, role) VALUES ($1,$2,'user')", req.Username, hashed)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(map[string]any{"success": true})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
    var req struct{ Username, Password string }
    json.NewDecoder(r.Body).Decode(&req)
    var id int
    var hashed, role string
    err := h.DB.QueryRow("SELECT id,password,role FROM users WHERE username=$1", req.Username).Scan(&id, &hashed, &role)
    if err != nil {
        http.Error(w, "invalid credentials", http.StatusUnauthorized)
        return
    }
    if !auth.CheckPassword(req.Password, hashed) {
        http.Error(w, "invalid credentials", http.StatusUnauthorized)
        return
    }
    token, _ := auth.GenerateToken(req.Username, role) // 2 arguments
    json.NewEncoder(w).Encode(map[string]any{"token": token, "role": role})
}
