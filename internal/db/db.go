package db

import (
    "log"

    "github.com/jmoiron/sqlx"
    _ "github.com/lib/pq"
)

var DB *sqlx.DB

func Init() {
    // Hardcoded Postgres connection string
    dbURL := "postgres://postgres:MySecretPass123@localhost:5432/delivery_db?sslmode=disable"

    var err error
    DB, err = sqlx.Connect("postgres", dbURL)
    if err != nil {
        log.Fatalf("Failed to connect to Postgres: %v", err)
    }

    log.Println("Connected to Postgres")
}
