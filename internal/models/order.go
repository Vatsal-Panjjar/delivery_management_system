package models

// Order represents an order placed by a user
type Order struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
	Status string `json:"status"`
}
