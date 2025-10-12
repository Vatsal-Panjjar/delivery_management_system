package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Vatsal-Panjjar/delivery_management_system/internal/auth"
	"github.com/Vatsal-Panjjar/delivery_management_system/internal/db"
)

type UserHandler struct {
	store *db.Store
}

func NewUserHandler(store *db.Store) *UserHandler {
	return &UserHandler{store: store}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	hash, _ := auth.HashPassword(req.Password)
	_, err := h.store.DB.Exec(`INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3)`,
		req.Username, req.Email, hash)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write([]byte("✅ User registered successfully"))
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	var storedHash, username, role string
	err := h.store.DB.QueryRow(`SELECT password_hash, username, role FROM users WHERE email=$1`, req.Email).
		Scan(&storedHash, &username, &role)
	if err != nil {
		http.Error(w, "invalid credentials", 401)
		return
	}

	if !auth.CheckPassword(storedHash, req.Password) {
		http.Error(w, "invalid password", 401)
		return
	}

	token, _ := auth.GenerateToken(username, role)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
