package repo

import (
	"database/sql"
	"fmt"

	"github.com/Vatsal-Panjiar/delivery_management_system/internal/models"
)

type DeliveryRepo struct {
	DB *sql.DB
}

func NewDeliveryRepo(db *sql.DB) *DeliveryRepo {
	return &DeliveryRepo{DB: db}
}

// Create a delivery
func (r *DeliveryRepo) Create(d *models.Delivery) error {
	query := `
		INSERT INTO deliveries (id, customer_id, courier_id, pickup_address, dropoff_address, status, price_cents, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,NOW(),NOW())
	`
	_, err := r.DB.Exec(query, d.ID, d.CustomerID, d.CourierID, d.PickupAddress, d.DropoffAddress, d.Status, d.PriceCents)
	return err
}

// Get delivery by ID
func (r *DeliveryRepo) GetByID(id string) (*models.Delivery, error) {
	d := &models.Delivery{}
	query := `
		SELECT id, customer_id, courier_id, pickup_address, dropoff_address, status, price_cents, created_at, updated_at
		FROM deliveries WHERE id=$1
	`
	err := r.DB.QueryRow(query, id).Scan(
		&d.ID, &d.CustomerID, &d.CourierID, &d.PickupAddress, &d.DropoffAddress,
		&d.Status, &d.PriceCents, &d.CreatedAt, &d.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return d, nil
}

// List deliveries by status with pagination
func (r *DeliveryRepo) ListByStatus(status string, offset, limit int) ([]*models.Delivery, error) {
	query := `
		SELECT id, customer_id, courier_id, pickup_address, dropoff_address, status, price_cents, created_at, updated_at
		FROM deliveries WHERE status=$1
		ORDER BY created_at DESC
		OFFSET $2 LIMIT $3
	`
	rows, err := r.DB.Query(query, status, offset, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	deliveries := []*models.Delivery{}
	for rows.Next() {
		d := &models.Delivery{}
		if err := rows.Scan(
			&d.ID, &d.CustomerID, &d.CourierID, &d.PickupAddress, &d.DropoffAddress,
			&d.Status, &d.PriceCents, &d.CreatedAt, &d.UpdatedAt,
		); err != nil {
			return nil, err
		}
		deliveries = append(deliveries, d)
	}
	return deliveries, nil
}
