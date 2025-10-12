package service

import (
	"delivery_management_system/internal/db"
	"delivery_management_system/internal/redis"
	"log"
	"sync"
	"time"
)

// TrackOrder simulates tracking an order and updating its status
func TrackOrder(orderID int, wg *sync.WaitGroup) {
	defer wg.Done()

	// Simulate order status progression
	statuses := []string{"pending", "shipped", "out_for_delivery", "delivered"}
	for _, status := range statuses {
		time.Sleep(2 * time.Second) // Simulate delay in processing status

		// Update PostgreSQL
		err := db.UpdateOrderStatus(orderID, status)
		if err != nil {
			log.Printf("Error updating order %d status: %v", orderID, err)
			return
		}

		// Cache in Redis for real-time tracking
		redis.SetOrderTracking(string(orderID), status)
	}
}

func StartTrackingOrder(orderID int) {
	var wg sync.WaitGroup
	wg.Add(1)

	go TrackOrder(orderID, &wg)

	wg.Wait() // Wait for the order tracking to complete
}
