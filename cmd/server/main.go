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

    // Hardcoded DB & Redis credentials
    dbURL := "postgres://postgres:rupupuru@01@localhost:5432/delivery_db?sslmode=disable"
    redisAddr := "localhost:6379"

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

    // Serve static files (CSS, JS)
    fs := http.FileServer(http.Dir("./web"))
    r.Handle("/static/*", http.StripPrefix("/static/", fs))

    // Serve HTML pages
    r.Get("/", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "./web/index.html")
    })
    r.Get("/login", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "./web/login.html")
    })
    r.Get("/register", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "./web/register.html")
    })
    r.Get("/user_dashboard", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "./web/user_dashboard.html")
    })
    r.Get("/admin_dashboard", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "./web/admin_dashboard.html")
    })

    // Mount backend API routes
    handlers.RegisterAuthRoutes(r, db, rdb)
    handlers.RegisterOrderRoutes(r, db, rdb)
    handlers.RegisterAdminRoutes(r, db, rdb)

    // Start server
    port := "8080"
    fmt.Printf("Server running on port %s\n", port)
    log.Fatal(http.ListenAndServe(":"+port, r))
}
