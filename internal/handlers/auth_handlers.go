package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/Vatsal-Panjjar/delivery_management_system/internal/auth"
	"github.com/Vatsal-Panjjar/delivery_management_system/internal/cache"
	"github.com/Vatsal-Panjjar/delivery_management_system/internal/models"
	"github.com/Vatsal-Panjjar/delivery_management_system/internal/db"
)

type AuthHandler struct {
	store *db.Store
	cache *cache.Cache
}

func RegisterAuthHandlers(r chi.Router, store *db.Store, cacheClient *cache.Cache) {
	h := &AuthHandler{store: store, cache: cacheClient}
	r.Post("/auth/register", h.Register)
	r.Post("/auth/login", h.Login)
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
	_, err = h.store.DB.Exec(`INSERT INTO users (username, email, password_hash, role, created_at) VALUES ($1,$2,$3,$4,now())`,
		req.Username, req.Email, hash, req.Role)
	if err != nil {
		http.Error(w, "cannot create user", http.StatusBadRequest)
		return
	}
	_ = h.cache.Set(context.Background(), "user:"+req.Email, req.Username, time.Hour)
	w.WriteHeader(http.StatusCreated)
}
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	var u models.User
	if err := h.store.DB.Get(&u, "SELECT id, username, password_hash, role FROM users WHERE email=$1", req.Email); err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}
	if !auth.CheckPassword(u.PasswordHash, req.Password) {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}
	tok, err := auth.GenerateToken(u.ID, u.Username, u.Role, 24*time.Hour)
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"token": tok, "username": u.Username, "role": u.Role})
}
