package main

import (
    "fmt"
    "log"
    "net/http"

    "github.com/go-chi/chi/v5"
    "github.com/go-redis/redis/v8"

    "github.com/Vatsal-Panjjar/delivery_management_system/internal/db"
    "github.com/Vatsal-Panjjar/delivery_management_system/internal/handlers"
)

func main() {
    fmt.Println("Starting Delivery Management Server...")

    // Hardcoded PostgreSQL connection string
    dbURL := "postgres://postgres:rupupuru@01@localhost:5432/delivery_db?sslmode=disable"

    // Connect to Postgres
    var err error
    db.DB, err = db.ConnectWithURL(dbURL)
    if err != nil {
        log.Fatalf("Failed to connect to Postgres: %v", err)
    }
    defer db.DB.Close()
    fmt.Println("Connected to Postgres")

    // Hardcoded Redis address
    rdb := redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
    })
    _, err = rdb.Ping(rdb.Context()).Result()
    if err != nil {
        log.Fatalf("Failed to connect to Redis: %v", err)
    }
    fmt.Println("Connected to Redis")

    // Create router
    r := chi.NewRouter()

    // Auth routes
    r.Post("/signup", handlers.SignupHandler)
    r.Post("/login", handlers.LoginHandler)

    // Order routes
    handlers.RegisterOrderRoutes(r, rdb)

    // Start server
    port := "8080"
    fmt.Printf("Server running on port %s\n", port)
    log.Fatal(http.ListenAndServe(":"+port, r))
}
