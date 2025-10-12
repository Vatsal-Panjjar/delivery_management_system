package models

// Status represents an order's status
type Status struct {
	OrderID string `json:"order_id"`
	Status  string `json:"status"`
}
