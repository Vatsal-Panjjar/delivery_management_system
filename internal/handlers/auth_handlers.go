package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/go-chi/chi/v5"
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
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	hashed, err := auth.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "failed to hash password", http.StatusInternalServerError)
		return
	}

	_, err = h.DB.Exec("INSERT INTO users (username, password, role) VALUES ($1,$2,'user')", req.Username, hashed)
	if err != nil {
		http.Error(w, "failed to create user", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]any{"success": true})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

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

	token, err := auth.GenerateToken(req.Username, role) // 2 args only
	if err != nil {
		http.Error(w, "failed to generate token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]any{"token": token, "role": role})
}
