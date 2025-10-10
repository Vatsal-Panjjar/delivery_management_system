package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

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

	if err := h.Repo.Create(&d); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Cache delivery for 5 minutes
	h.Cache.Set(d.ID, d, 5*time.Minute)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(d)
}

// Get delivery by ID
func (h *DeliveryHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if cached, found := h.Cache.Get(id); found {
		json.NewEncoder(w).Encode(cached)
		return
	}

	d, err := h.Repo.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Cache delivery for 5 minutes
	h.Cache.Set(d.ID, d, 5*time.Minute)
	json.NewEncoder(w).Encode(d)
}

// List deliveries by status with pagination
func (h *DeliveryHandler) ListByStatus(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 10
	offset := 0

	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
		limit = l
	}
	if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
		offset = o
	}

	deliveries, err := h.Repo.ListByStatus(status, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(deliveries)
}

// Register delivery routes
func RegisterDeliveryRoutes(r chi.Router, h *DeliveryHandler) {
	r.Post("/deliveries", h.Create)
	r.Get("/deliveries/{id}", h.Get)
	r.Get("/deliveries", h.ListByStatus)
}
