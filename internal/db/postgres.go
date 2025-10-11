package db

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// NewPostgres connects to the PostgreSQL database using the given URL.
func NewPostgres(url string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", url)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// ExecSchema executes a SQL schema string (used for migrations).
func ExecSchema(db *sqlx.DB, schemaSQL string) error {
	_, err := db.Exec(schemaSQL)
	return err
}
