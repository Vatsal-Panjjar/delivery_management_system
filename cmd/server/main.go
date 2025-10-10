package main

import (
    "fmt"
    "net/http"

    "github.com/go-chi/chi/v5"
    "github.com/Vatsal-Panjiar/delivery_management_system/internal/handlers"
)

func main() {
    r := chi.NewRouter()

    // Register your routes
    handlers.RegisterRoutes(r)

    fmt.Println("Server running on :8080")
    http.ListenAndServe(":8080", r)
}
