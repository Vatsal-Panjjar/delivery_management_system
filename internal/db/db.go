package db

import (
    "github.com/jmoiron/sqlx"
    _ "github.com/lib/pq"
    "log"
    "os"
)

var DB *sqlx.DB

func Connect() {
    dbURL := os.Getenv("POSTGRES_URL") // e.g. postgres://user:pass@localhost:5432/delivery_db?sslmode=disable
    var err error
    DB, err = sqlx.Connect("postgres", dbURL)
    if err != nil {
        log.Fatalf("Failed to connect to Postgres: %v", err)
    }
}
