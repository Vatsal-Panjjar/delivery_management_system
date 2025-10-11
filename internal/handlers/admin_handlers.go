package handlers

import (
    "net/http"
    "encoding/json"
    "github.com/Vatsal-Panjjar/delivery_management_system/internal/auth"
    "github.com/Vatsal-Panjjar/delivery_management_system/internal/db"
)

type User struct {
    ID       int    `db:"id"`
    Username string `db:"username"`
    Password string `db:"password"`
    Role     string `db:"role"`
}

func SignupHandler(w http.ResponseWriter, r *http.Request) {
    var user User
    json.NewDecoder(r.Body).Decode(&user)
    hashedPassword, _ := auth.HashPassword(user.Password)
    _, err := db.DB.Exec("INSERT INTO users (username, password, role) VALUES ($1, $2, $3)", user.Username, hashedPassword, "user")
    if err != nil {
        http.Error(w, "Signup failed", http.StatusInternalServerError)
        return
    }
    w.Write([]byte("User registered successfully"))
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
    var creds User
    json.NewDecoder(r.Body).Decode(&creds)

    var user User
    err := db.DB.Get(&user, "SELECT * FROM users WHERE username=$1", creds.Username)
    if err != nil || !auth.CheckPassword(creds.Password, user.Password) {
        http.Error(w, "Invalid credentials", http.StatusUnauthorized)
        return
    }

    token, _ := auth.GenerateToken(string(rune(user.ID)), user.Role)
    json.NewEncoder(w).Encode(map[string]string{"token": token})
}
