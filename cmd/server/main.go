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

    dbURL := "postgres://postgres:rupupuru@01@localhost:5432/delivery_db?sslmode=disable"
    redisAddr := "localhost:6379"

    db, err := sqlx.Connect("postgres", dbURL)
    if err != nil {
        log.Fatalf("Failed to connect to Postgres: %v", err)
    }
    defer db.Close()
    fmt.Println("Connected to Postgres")

    rdb := redis.NewClient(&redis.Options{Addr: redisAddr})
    _, err = rdb.Ping(rdb.Context()).Result()
    if err != nil {
        log.Fatalf("Failed to connect to Redis: %v", err)
    }
    fmt.Println("Connected to Redis")

    r := chi.NewRouter()

    // Root
    r.Get("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Delivery Management System API is running"))
    })

    // Auth routes (updated to accept *chi.Mux)
    handlers.RegisterRoutes(r, db)

    // Order routes
    handlers.RegisterOrderRoutes(r, db, rdb)

    // Serve static HTML
    fs := http.FileServer(http.Dir("./web"))
    r.Handle("/*", http.StripPrefix("/", fs))

    port := "8080"
    fmt.Printf("Server running on http://localhost:%s\n", port)
    log.Fatal(http.ListenAndServe(":"+port, r))
}
