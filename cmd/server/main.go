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

    // --- PostgreSQL connection ---
    dbURL := "postgres://postgres:MySecretPass123@localhost:5432/delivery_db?sslmode=disable"
    db, err := sqlx.Connect("postgres", dbURL)
    if err != nil {
        log.Fatalf("Failed to connect to Postgres: %v", err)
    }
    defer db.Close()
    fmt.Println("Connected to Postgres")

    // --- Redis connection ---
    rdb := redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
    })
    _, err = rdb.Ping(rdb.Context()).Result()
    if err != nil {
        log.Fatalf("Failed to connect to Redis: %v", err)
    }
    fmt.Println("Connected to Redis")

    // --- Router ---
    r := chi.NewRouter()

    // --- Routes ---
    r.Get("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Delivery Management System API is running"))
    })

    // Register auth routes
    handlers.RegisterAuthRoutes(r, db, rdb)

    // Register order routes
    handlers.RegisterOrderRoutes(r, db, rdb)

    // Register admin routes
    handlers.RegisterAdminRoutes(r, db)

    // --- Start server ---
    port := "8080"
    fmt.Printf("Server running on port %s\n", port)
    log.Fatal(http.ListenAndServe(":"+port, r))
}
