package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/Vatsal-Panjiar/delivery_management_system/internal/cache"
	"github.com/Vatsal-Panjiar/delivery_management_system/internal/models"
	"github.com/Vatsal-Panjiar/delivery_management_system/internal/repo"
)

type DeliveryHandler struct {
	Repo  *repo.DeliveryRepo
	Cache *cache.RedisCache
}

// NewDeliveryHandler creates a new DeliveryHandler with repo and Redis cache
func NewDeliveryHandler(r *repo.DeliveryRepo, c *cache.RedisCache) *DeliveryHandler {
	return &DeliveryHandler{
		Repo:  r,
		Cache: c,
	}
}

// Create a new delivery
func (h *DeliveryHandler) Create(w http.ResponseWriter, r *http.Request) {
	var d models.Delivery
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if d.ID == "" {
		d.ID = uuid.New().String()
	}

	if err := h.Repo.Create(&d); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Optionally, cache the delivery in Redis
	if h.Cache != nil {
		h.Cache.Set(d.ID, d)
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(d)
}

// Get delivery by ID
func (h *DeliveryHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	// Check Redis cache first
	if h.Cache != nil {
		if cached, found := h.Cache.Get(id); found {
			json.NewEncoder(w).Encode(cached)
			return
		}
	}

	d, err := h.Repo.GetByID(id)
	if err != nil {
		http.Error(w, "delivery not found", http.StatusNotFound)
		return
	}

	// Store in cache
	if h.Cache != nil {
		h.Cache.Set(d.ID, d)
	}

	json.NewEncoder(w).Encode(d)
}

// List deliveries by status
func (h *DeliveryHandler) ListByStatus(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")

	deliveries, err := h.Repo.ListByStatus(status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(deliveries)
}

// RegisterDeliveryRoutes sets up the routes
func RegisterDeliveryRoutes(r chi.Router, h *DeliveryHandler) {
	r.Post("/deliveries", h.Create)
	r.Get("/deliveries/{id}", h.Get)
	r.Get("/deliveries", h.ListByStatus)
}
