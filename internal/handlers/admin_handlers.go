package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Vatsal-Panjjar/delivery_management_system/internal/db"
)

type AdminHandler struct {
	store *db.Store
}

func NewAdminHandler(store *db.Store) *AdminHandler {
	return &AdminHandler{store: store}
}

func (h *AdminHandler) GetAllOrders(w http.ResponseWriter, r *http.Request) {
	rows, _ := h.store.DB.Query(`SELECT id, user_id, pickup_address, delivery_address, status FROM orders`)
	defer rows.Close()

	var orders []map[string]interface{}
	for rows.Next() {
		var o struct {
			ID              int
			UserID          int
			PickupAddress   string
			DeliveryAddress string
			Status          string
		}
		rows.Scan(&o.ID, &o.UserID, &o.PickupAddress, &o.DeliveryAddress, &o.Status)
		orders = append(orders, map[string]interface{}{
			"id":       o.ID,
			"user_id":  o.UserID,
			"pickup":   o.PickupAddress,
			"delivery": o.DeliveryAddress,
			"status":   o.Status,
		})
	}

	json.NewEncoder(w).Encode(orders)
}
