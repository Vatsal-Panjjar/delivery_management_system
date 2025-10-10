package handlers

import (
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r *chi.Mux, auth *AuthHandler, delivery *DeliveryHandler) {
	r.Post("/signup", auth.Signup)
	r.Post("/login", auth.Login)

	r.Route("/deliveries", func(r chi.Router) {
		r.Post("/", delivery.Create)
		r.Get("/", delivery.ListByStatus)
		r.Get("/{id}", delivery.Get)
	})
}
