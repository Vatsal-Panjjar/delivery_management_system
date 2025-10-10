package handlers

import "github.com/go-chi/chi/v5"

func RegisterDeliveryRoutes(r chi.Router, h *DeliveryHandler) {
	r.Route("/deliveries", func(r chi.Router) {
		r.Get("/", h.ListByStatus)
		r.Post("/", h.Create)
		r.Get("/{id}", h.Get)
	})
}
