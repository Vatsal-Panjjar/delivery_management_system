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

type AuthResponse struct {
    Token string `json:"token"`
    Role  string `json:"role"`
}

// Dummy users for demo
var users = map[string]struct {
    Password string
    Role     string
}{
    "admin": {Password: "admin123", Role: "admin"},
    "user":  {Password: "user123", Role: "user"},
}

// RegisterAuthRoutes registers login endpoints
func RegisterAuthRoutes(r *http.ServeMux, db *sqlx.DB) {
    r.HandleFunc("/login", loginHandler)
}

// loginHandler handles user login
func loginHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var req AuthRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    user, ok := users[req.Username]
    if !ok || user.Password != req.Password {
        http.Error(w, "Invalid credentials", http.StatusUnauthorized)
        return
    }

    // Generate token (2 arguments: username, role)
    token, err := auth.GenerateToken(req.Username, user.Role)
    if err != nil {
        http.Error(w, "Failed to generate token", http.StatusInternalServerError)
        return
    }

    resp := AuthResponse{
        Token: token,
        Role:  user.Role,
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(resp)
}
