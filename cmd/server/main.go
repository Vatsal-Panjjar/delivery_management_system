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

    // PostgreSQL connection
    dbURL := "postgres://postgres:rupupuru@01@localhost:5432/delivery_db?sslmode=disable"
    db, err := sqlx.Connect("postgres", dbURL)
    if err != nil {
        log.Fatalf("Failed to connect to Postgres: %v", err)
    }
    defer db.Close()
    fmt.Println("Connected to Postgres")

    // Redis connection
    rdb := redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
    })
    _, err = rdb.Ping(rdb.Context()).Result()
    if err != nil {
        log.Fatalf("Failed to connect to Redis: %v", err)
    }
    fmt.Println("Connected to Redis")

    // Router
    r := chi.NewRouter()

    // Serve frontend files
    r.Handle("/web/*", http.StripPrefix("/web/", http.FileServer(http.Dir("./web"))))

    // API routes
    apiMux := http.NewServeMux()
    handlers.RegisterAuthRoutes(apiMux, db)      // login/register
    handlers.RegisterOrderRoutes(r, db, rdb)    // orders
    r.Mount("/api", apiMux)

    port := "8080"
    fmt.Printf("Server running on port %s\n", port)
    log.Fatal(http.ListenAndServe(":"+port, r))
}
