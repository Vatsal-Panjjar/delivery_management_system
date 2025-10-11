package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/go-redis/redis/v8"
	"github.com/Vatsal-Panjjar/delivery_management_system/internal/auth"
)

type AuthHandler struct {
	DB  *sqlx.DB
	RDB *redis.Client
}

func RegisterAuthRoutes(r chi.Router, db *sqlx.DB, rdb *redis.Client) {
	h := &AuthHandler{DB: db, RDB: rdb}
	r.Post("/register", h.Register)
	r.Post("/login", h.Login)
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	type Req struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}
	var req Req
	json.NewDecoder(r.Body).Decode(&req)

	_, err := h.DB.Exec(`INSERT INTO users (username, password, role, created_at) VALUES ($1, $2, $3, $4)`,
		req.Username, req.Password, req.Role, time.Now())
	if err != nil {
		http.Error(w, "User registration failed", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("User registered successfully"))
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	type Req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	var req Req
	json.NewDecoder(r.Body).Decode(&req)

	var id string
	var role string
	var dbPass string
	err := h.DB.QueryRow(`SELECT id, role, password FROM users WHERE username=$1`, req.Username).Scan(&id, &role, &dbPass)
	if err != nil || dbPass != req.Password {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, _ := auth.GenerateToken(id, role)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
