package models

import "time"

type Delivery struct {
	ID          string    `db:"id" json:"id"`
	CustomerID  string    `db:"customer_id" json:"customer_id"`
	CourierID   *string   `db:"courier_id" json:"courier_id,omitempty"`
	PickupAddr  string    `db:"pickup_address" json:"pickup_address"`
	DropoffAddr string    `db:"dropoff_address" json:"dropoff_address"`
	Status      string    `db:"status" json:"status"`
	PriceCents  int       `db:"price_cents" json:"price_cents"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}
