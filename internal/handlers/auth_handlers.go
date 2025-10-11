package handlers

import (
    "encoding/json"
    "net/http"

    "github.com/Vatsal-Panjjar/delivery_management_system/internal/auth"
    "github.com/jmoiron/sqlx"
)

type RegisterRequest struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

type LoginRequest struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

func RegisterHandler(db *sqlx.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var req RegisterRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            http.Error(w, "Invalid request", http.StatusBadRequest)
            return
        }
        _, err := db.Exec("INSERT INTO users (username, password, role) VALUES ($1,$2,'customer')",
            req.Username, req.Password)
        if err != nil {
            http.Error(w, "Failed to register", http.StatusInternalServerError)
            return
        }
        w.WriteHeader(http.StatusCreated)
    }
}

func LoginHandler(db *sqlx.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var req LoginRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            http.Error(w, "Invalid request", http.StatusBadRequest)
            return
        }

        var dbPassword, role string
        err := db.QueryRow("SELECT password, role FROM users WHERE username=$1", req.Username).
            Scan(&dbPassword, &role)
        if err != nil || dbPassword != req.Password {
            http.Error(w, "Invalid credentials", http.StatusUnauthorized)
            return
        }

        token, err := auth.GenerateToken(req.Username, role)
        if err != nil {
            http.Error(w, "Failed to generate token", http.StatusInternalServerError)
            return
        }

        json.NewEncoder(w).Encode(map[string]string{
            "token": token,
            "role":  role,
        })
    }
}
