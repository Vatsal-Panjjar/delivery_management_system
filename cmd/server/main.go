package main

import (
    "fmt"
    "log"
    "net/http"
    "os"

    "github.com/go-chi/chi/v5"
    "github.com/go-redis/redis/v8"
    "github.com/jmoiron/sqlx"
    _ "github.com/lib/pq"

    "github.com/Vatsal-Panjjar/delivery_management_system/internal/handlers"
)

func main() {
    fmt.Println("Starting Delivery Management Server...")

    // Database connection (hardcode your password if you want)
    dbURL := "postgres://postgres:rupupuru@01@localhost:5432/delivery_db?sslmode=disable"
    db, err := sqlx.Connect("postgres", dbURL)
    if err != nil {
        log.Fatalf("Failed to connect to Postgres: %v", err)
    }
    defer db.Close()
    fmt.Println("Connected to Postgres")

    // Redis
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

    // Serve frontend
    r.Handle("/web/*", http.StripPrefix("/web/", http.FileServer(http.Dir("./web"))))

    // API Routes
    handlers.RegisterAuthRoutes(r, db, rdb)
    handlers.RegisterOrderRoutes(r, db, rdb)
    handlers.RegisterAdminRoutes(r, db, rdb)

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    fmt.Printf("Server running on port %s\n", port)
    log.Fatal(http.ListenAndServe(":"+port, r))
}
