package handlers

import (
    "database/sql"
    "encoding/json"
    "net/http"

    "github.com/jmoiron/sqlx"
    "github.com/Vatsal-Panjjar/delivery_management_system/internal/auth"
)

type User struct {
    ID       int    `db:"id" json:"id"`
    Username string `db:"username" json:"username"`
    Role     string `db:"role" json:"role"`
}

func RegisterUserHandler(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var input struct {
        Username string `json:"username"`
        Password string `json:"password"`
    }
    if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        http.Error(w, "Invalid input", http.StatusBadRequest)
        return
    }

    hash, _ := auth.HashPassword(input.Password)

    _, err := db.Exec("INSERT INTO users (username, password_hash) VALUES ($1, $2)", input.Username, hash)
    if err != nil {
        http.Error(w, "Username already exists", http.StatusBadRequest)
        return
    }

    w.WriteHeader(http.StatusCreated)
    w.Write([]byte("User registered successfully"))
}

func LoginHandler(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var input struct {
        Username string `json:"username"`
        Password string `json:"password"`
    }
    if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        http.Error(w, "Invalid input", http.StatusBadRequest)
        return
    }

    var user struct {
        ID           int    `db:"id"`
        Username     string `db:"username"`
        PasswordHash string `db:"password_hash"`
        Role         string `db:"role"`
    }

    err := db.Get(&user, "SELECT id, username, password_hash, role FROM users WHERE username=$1", input.Username)
    if err != nil {
        if err == sql.ErrNoRows {
            http.Error(w, "Invalid credentials", http.StatusUnauthorized)
            return
        }
        http.Error(w, "Server error", http.StatusInternalServerError)
        return
    }

    if err := auth.CheckPassword(user.PasswordHash, input.Password); err != nil {
        http.Error(w, "Invalid credentials", http.StatusUnauthorized)
        return
    }

    token, _ := auth.GenerateToken(user.Username, user.Role)
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
        "token": token,
        "role":  user.Role,
    })
}

func RegisterAuthRoutes(r *http.ServeMux, db *sqlx.DB) {
    r.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
        RegisterUserHandler(w, r, db)
    })
    r.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
        LoginHandler(w, r, db)
    })
}
