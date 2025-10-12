package db

import (
	"database/sql"
	"errors"
	"log"
	"delivery_management_system/internal/models"
	"github.com/lib/pq"
)

// Initialize connects to the PostgreSQL database
func Initialize(connStr string) {
	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
}

// CreateOrder creates a new order in the database
func CreateOrder(order *models.Order) (string, error) {
	var orderID string
	err := DB.QueryRow("INSERT INTO orders (user_id, status) VALUES ($1, $2) RETURNING id", order.UserID, order.Status).Scan(&orderID)
	if err != nil {
		log.Println("Error inserting order:", err)
		return "", err
	}
	return orderID, nil
}

// UpdateOrderStatus updates the status of an existing order
func UpdateOrderStatus(orderID, status string) error {
	_, err := DB.Exec("UPDATE orders SET status = $1 WHERE id = $2", status, orderID)
	if err != nil {
		log.Println("Error updating order status:", err)
		return err
	}
	return nil
}

// GetOrderStatus retrieves the status of an order by its ID
func GetOrderStatus(orderID string) (string, error) {
	var status string
	err := DB.QueryRow("SELECT status FROM orders WHERE id = $1", orderID).Scan(&status)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New("order not found")
		}
		log.Println("Error fetching order status:", err)
		return "", err
	}
	return status, nil
}
