package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Vatsal-Panjjar/delivery_management_system/internal/db"
)

type OrderHandler struct {
	store *db.Store
}

func NewOrderHandler(store *db.Store) *OrderHandler {
	return &OrderHandler{store: store}
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID          int    `json:"user_id"`
		PickupAddress   string `json:"pickup"`
		DeliveryAddress string `json:"delivery"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	_, err := h.store.DB.Exec(`INSERT INTO orders (user_id, pickup_address, delivery_address) VALUES ($1, $2, $3)`,
		req.UserID, req.PickupAddress, req.DeliveryAddress)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write([]byte("✅ Order created"))
}
