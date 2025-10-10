package handlers

import (
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r *chi.Mux, dh *DeliveryHandler, ah *AuthHandler) {
	r.Post("/register", ah.Register)
	r.Post("/login", ah.Login)

	r.Post("/deliveries", dh.Create)
	r.Get("/deliveries", dh.ListByStatus)
	r.Get("/deliveries/{id}", dh.Get)
}
