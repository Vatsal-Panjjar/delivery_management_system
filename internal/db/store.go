package db

import (
	"database/sql"
	"log"
	"fmt"

	_ "github.com/lib/pq"
)

var db *sql.DB

// Initialize DB connection
func Initialize(connStr string) {
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Unable to connect to database:", err)
	}
}

// Close DB connection
func Close() {
	if db != nil {
		db.Close()
	}
}

// CreateOrder creates a new order in the database
func CreateOrder(userID int, status string) (int, error) {
	var orderID int
	query := `INSERT INTO orders (user_id, order_status) VALUES ($1, $2) RETURNING id`
	err := db.QueryRow(query, userID, status).Scan(&orderID)
	if err != nil {
		return 0, fmt.Errorf("could not create order: %v", err)
	}
	return orderID, nil
}
