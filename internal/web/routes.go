
package web

import (
    "github.com/gorilla/mux"
    "net/http"
)

func RegisterRoutes(r *mux.Router) {
    // Define your routes here
    r.HandleFunc("/api/your-endpoint", YourHandler).Methods("GET")
}
