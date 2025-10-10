package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/Vatsal-Panjiar/delivery_management_system/internal/cache"
	"github.com/Vatsal-Panjiar/delivery_management_system/internal/models"
	"github.com/Vatsal-Panjiar/delivery_management_system/internal/repo"
)

type DeliveryHandler struct {
	Repo  *repo.DeliveryRepo
	Cache *cache.RedisCache
}

func NewDeliveryHandler(r *repo.DeliveryRepo, c *cache.RedisCache) *DeliveryHandler {
	return &DeliveryHandler{Repo: r, Cache: c}
}

func (h *DeliveryHandler) ListByStatus(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	deliveries, _ := h.Repo.ListByStatus(status)
	json.NewEncoder(w).Encode(deliveries)
}

func RegisterDeliveryRoutes(r chi.Router, h *DeliveryHandler) {
	r.Get("/deliveries", h.ListByStatus)
}
