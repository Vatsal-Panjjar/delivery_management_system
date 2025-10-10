package repo

import (
	"database/sql"
	"errors"
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
    query := `INSERT INTO deliveries (id, customer_id, courier_id, pickup_address, dropoff_address, status, price_cents, created_at, updated_at)
              VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`
    _, err := r.DB.Exec(query, d.ID, d.CustomerID, d.CourierID, d.PickupAddress, d.DropoffAddress, d.Status, d.PriceCents, time.Now(), time.Now())
    return err
}


func (r *DeliveryRepo) GetByID(id string) (*models.Delivery, error) {
	var d models.Delivery
	err := r.DB.Get(&d, "SELECT * FROM deliveries WHERE id=$1", id)
	if err == sql.ErrNoRows {
		return nil, errors.New("not found")
	}
	return &d, err
}

func (r *DeliveryRepo) UpdateStatus(id, status string, courierID *int) error {
	_, err := r.DB.Exec("UPDATE deliveries SET status=$1, courier_id=$2, updated_at=$3 WHERE id=$4", status, courierID, time.Now(), id)
	return err
}

func (r *DeliveryRepo) ListByStatus(status string, limit, offset int) ([]models.Delivery, error) {
	var ds []models.Delivery
	err := r.DB.Select(&ds, "SELECT * FROM deliveries WHERE status=$1 ORDER BY created_at DESC LIMIT $2 OFFSET $3", status, limit, offset)
	return ds, err
}
