package service

import (
	"delivery_management_system/internal/db"
	"delivery_management_system/internal/redis"
	"strconv"
)

// UpdateOrderStatus updates the status of an order
func UpdateOrderStatus(orderID int, status string) error {
	// Convert orderID (int) to string
	orderIDStr := strconv.Itoa(orderID)

	// Update order status in DB
	err := db.UpdateOrderStatus(orderIDStr, status)
	if err != nil {
		return err
	}

	// Set order tracking in Redis (simulated)
	trackingInfo := "Order status updated to: " + status
	err = redis.SetOrderTracking(orderIDStr, trackingInfo)
	if err != nil {
		return err
	}

	return nil
}
