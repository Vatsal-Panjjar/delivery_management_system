import (
    "encoding/json"
    "net/http"

    "github.com/jmoiron/sqlx"
    "github.com/go-chi/chi/v5"
    "github.com/Vatsal-Panjjar/delivery_management_system/internal/auth"
)

func RegisterRoutes(r *chi.Mux, db *sqlx.DB) {
    r.Post("/auth/register", func(w http.ResponseWriter, r *http.Request) {
        var req struct {
            Username string `json:"username"`
            Password string `json:"password"`
        }
        json.NewDecoder(r.Body).Decode(&req)
        hash, _ := auth.HashPassword(req.Password)
        _, err := db.Exec("INSERT INTO users(username, password) VALUES($1,$2)", req.Username, hash)
        if err != nil {
            http.Error(w, "Registration failed", http.StatusBadRequest)
            return
        }
        w.WriteHeader(http.StatusOK)
    })

    r.Post("/auth/login", func(w http.ResponseWriter, r *http.Request) {
        var req struct {
            Username string `json:"username"`
            Password string `json:"password"`
        }
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
