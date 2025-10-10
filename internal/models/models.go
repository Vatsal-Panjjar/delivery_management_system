package models


import "time"


type User struct {
ID int `db:"id" json:"id"`
Name string `db:"name" json:"name"`
Email string `db:"email" json:"email"`
PasswordHash string `db:"password_hash" json:"-"`
Role string `db:"role" json:"role"`
CreatedAt time.Time `db:"created_at" json:"created_at"`
}


type Delivery struct{
ID string `db:"id" json:"id"`
CustomerID int `db:"customer_id" json:"customer_id"`
CourierID *int `db:"courier_id" json:"courier_id,omitempty"`
PickupAddress string `db:"pickup_address" json:"pickup_address"`
DropoffAddress string `db:"dropoff_address" json:"dropoff_address"`
Status string `db:"status" json:"status"`
PriceCents int `db:"price_cents" json:"price_cents"`
CreatedAt time.Time `db:"created_at" json:"created_at"`
UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
