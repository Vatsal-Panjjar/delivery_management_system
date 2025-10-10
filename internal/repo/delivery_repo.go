package repo

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/Vatsal-Panjiar/delivery_management_system/internal/models"
)

type DeliveryRepo struct {
	DB *sqlx.DB
}

func NewDeliveryRepo(db *sqlx.DB) *DeliveryRepo {
	return &DeliveryRepo{DB: db}
}

func (r *DeliveryRepo) Create(d *models.Delivery) error {
	_, err := r.DB.Exec(
		`INSERT INTO deliveries (id, customer_id, courier_id, pickup_address, dropoff_address, status, price_cents, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
		d.ID, d.CustomerID, d.CourierID, d.PickupAddr, d.DropoffAddr, d.Status,
		d.PriceCents, time.Now(), time.Now(),
	)
	return err
}

func (r *DeliveryRepo) ListByStatus(status string) ([]models.Delivery, error) {
	var ds []models.Delivery
	err := r.DB.Select(&ds, "SELECT * FROM deliveries WHERE status=$1 ORDER BY created_at DESC", status)
	return ds, err
}
