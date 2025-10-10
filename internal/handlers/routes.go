package handlers

import (
	"github.com/go-chi/chi/v5"
)

func RegisterAuthRoutes(r chi.Router, h *AuthHandler) {
	r.Post("/signup", h.Signup)
	r.Post("/login", h.Login)
}
