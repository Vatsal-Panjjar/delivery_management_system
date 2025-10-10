package repo

import (
	"github.com/jmoiron/sqlx"
	"github.com/Vatsal-Panjiar/delivery_management_system/internal/models"
)

type DeliveryRepo struct {
	DB *sqlx.DB
}

func NewDeliveryRepo(db *sqlx.DB) *DeliveryRepo {
	return &DeliveryRepo{DB: db}
}

// Create a new delivery
func (r *DeliveryRepo) Create(d *models.Delivery) error {
	_, err := r.DB.Exec(
		`INSERT INTO deliveries (id, customer_id, courier_id, pickup_address, dropoff_address, stat
