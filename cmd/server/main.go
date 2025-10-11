package main

import (
    "fmt"
    "log"
    "net/http"

    "github.com/go-chi/chi/v5"
    "github.com/go-redis/redis/v8"
    "github.com/jmoiron/sqlx"
    _ "github.com/lib/pq"

    "github.com/Vatsal-Panjjar/delivery_management_system/internal/handlers"
)

func main() {
    fmt.Println("Starting Delivery Management Server...")

    // === Hardcoded database and Redis info ===
    dbURL := "postgres://postgres:MySecretPass123@localhost:5432/delivery_db?sslmode=disable"
    redisAddr := "localhost:6379"
    // ========================================

    // Connect to Postgres
    db, err := sqlx.Connect("postgres", dbURL)
    if err != nil {
        log.Fatalf("Failed to connect to Postgres: %v", err)
    }
    defer db.Close()
    fmt.Println("Connected to Postgres")

    // Connect to Redis
    rdb := redis.NewClient(&redis.Options{
        Addr: redisAddr,
    })
    _, err = rdb.Ping(rdb.Context()).Result()
    if err != nil {
        log.Fatalf("Failed to connect to Redis: %v", err)
    }
    fmt.Println("Connected to Redis")

    // Create router
    r := chi.NewRouter()

    // Redirect root to login page
    r.Get("/", func(w http.ResponseWriter, r *http.Request) {
        http.Redirect(w, r, "/web/index.html", http.StatusSeeOther)
    })

    // Serve static web files
    r.Handle("/web/*", http.StripPrefix("/web/", http.FileServer(http.Dir("./web"))))

    // API base route
    r.Get("/api", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Delivery Management API is running"))
    })

    // Mount order routes
    handlers.RegisterOrderRoutes(r, db, rdb)

    // Start server on port 8080
    port := "8080"
    fmt.Printf("Server running on port %s\n", port)
    log.Fatal(http.ListenAndServe(":"+port, r))
}
