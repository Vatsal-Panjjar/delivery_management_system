package repo

import (
	"github.com/jmoiron/sqlx"
	"github.com/Vatsal-Panjiar/delivery_management_system/internal/models"
)

type DeliveryRepo struct {
	DB *sqlx.DB
}

func NewDeliveryRepo(db *sqlx.DB) *DeliveryRepo {
	return &DeliveryRepo{
		DB: db,
	}
}

func (r *DeliveryRepo) Create(d *models.Delivery) error {
	_, err := r.DB.Exec(
		`INSERT INTO deliveries (id, customer_id, courier_id, pickup_address, dropoff_address, status, price_cents)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		d.ID, d.CustomerID, d.CourierID, d.PickupAddress, d.DropoffAddress, d.Status, d.PriceCents,
	)
	return err
}

func (r *DeliveryRepo) GetByID(id string) (*models.Delivery, error) {
	var d models.Delivery
	err := r.DB.Get(&d, "SELECT * FROM deliveries WHERE id=$1", id)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func (r *DeliveryRepo) ListByStatus(status string, limit, offset int) ([]*models.Delivery, error) {
	var deliveries []*models.Delivery
	err := r.DB.Select(&deliveries,
		"SELECT * FROM deliveries WHERE status=$1 ORDER BY created_at DESC LIMIT $2 OFFSET $3",
		status, limit, offset,
	)
	if err != nil {
		return nil, err
	}
	return deliveries, nil
}
