package handlers

import (
	"github.com/go-chi/chi/v5"
)

// RegisterRoutes sets up all routes for deliveries and auth
func RegisterRoutes(r chi.Router, delivery *DeliveryHandler, auth *AuthHandler) {
	r.Post("/deliveries", delivery.Create)
	r.Get("/deliveries/{id}", delivery.Get)
	r.Get("/deliveries", delivery.ListByStatus)

	r.Post("/signup", auth.Signup)
	r.Post("/login", auth.Login)
}
