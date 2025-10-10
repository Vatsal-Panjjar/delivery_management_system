package handlers

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/Vatsal-Panjiar/delivery_management_system/internal/models"
	"github.com/Vatsal-Panjiar/delivery_management_system/internal/repo"
	"github.com/Vatsal-Panjiar/delivery_management_system/internal/cache"
)

type DeliveryHandler struct {
	Repo  *repo.DeliveryRepo
	Cache *cache.RedisCache
	mu    sync.Mutex
}

func NewDeliveryHandler(r *repo.DeliveryRepo, c *cache.RedisCache) *DeliveryHandler {
	return &DeliveryHandler{Repo: r, Cache: c}
}

// Create a new delivery
func (h *DeliveryHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.Delivery
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	req.ID = uuid.NewString()
	req.Status = "pending"
	req.CreatedAt = time.Now()
	req.UpdatedAt = time.Now()

	if err := h.Repo.Create(&req); err != nil {
		http.Error(w, "failed to create delivery", http.StatusInternalServerError)
		return
	}

	_ = h.Cache.Del("deliveries:pending")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(req)
}

// Get delivery by ID
func (h *DeliveryHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	cacheKey := "delivery:" + id

	if cached, err := h.Cache.Get(cacheKey); err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(cached))
		return
	}

	d, err := h.Repo.GetByID(id)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	b, _ := json.Marshal(d)
	_ = h.Cache.Set(cacheKey, b, 5*time.Minute)

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

// List deliveries by status
func (h *DeliveryHandler) ListByStatus(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	if status == "" {
		status = "pending"
	}

	cacheKey := "deliveries:" + status
	if cached, err := h.Cache.Get(cacheKey); err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(cached))
		return
	}

	ds, err := h.Repo.ListByStatus(status, 100, 0)
	if err != nil {
		http.Error(w, "failed to fetch deliveries", http.StatusInternalServerError)
		return
	}

	b, _ := json.Marshal(ds)
	_ = h.Cache.Set(cacheKey, b, 5*time.Minute)

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

// RegisterRoutes registers all delivery endpoints
func RegisterRoutes(r *chi.Mux, h *DeliveryHandler) {
	r.Post("/deliveries", h.Create)
	r.Get("/deliveries/{id}", h.Get)
	r.Get("/deliveries", h.ListByStatus)
}
