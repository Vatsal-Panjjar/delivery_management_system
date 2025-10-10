package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/Vatsal-Panjiar/delivery_management_system/internal/models"
	"github.com/Vatsal-Panjiar/delivery_management_system/internal/repo"
)

type DeliveryHandler struct {
	Repo *repo.DeliveryRepo
}

func NewDeliveryHandler(r *repo.DeliveryRepo) *DeliveryHandler {
	return &DeliveryHandler{Repo: r}
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

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(d)
}

// Get delivery by ID
func (h *DeliveryHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	d, err := h.Repo.GetByID(id)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(d)
}

// List deliveries by status
func (h *DeliveryHandler) ListByStatus(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	limit := 100 // default
	offset := 0  // default

	ds, err := h.Repo.ListByStatus(status, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(ds)
}
