package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/Vatsal-Panjjar/delivery_management_system/internal/auth"
)

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func RegisterRoutes(mux *http.ServeMux, db *sqlx.DB) {
	mux.HandleFunc("/auth/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var req AuthRequest
		json.NewDecoder(r.Body).Decode(&req)
		hash, _ := auth.HashPassword(req.Password)
		_, err := db.Exec("INSERT INTO users(username, password) VALUES($1,$2)", req.Username, hash)
		if err != nil {
			http.Error(w, "Registration failed", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc("/auth/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var req AuthRequest
		json.NewDecoder(r.Body).Decode(&req)

		var hash string
		err := db.QueryRow("SELECT password FROM users WHERE username=$1", req.Username).Scan(&hash)
		if err != nil || !auth.CheckPassword(hash, req.Password) {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		token, _ := auth.GenerateToken(req.Username)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"token": token})
	})
}
