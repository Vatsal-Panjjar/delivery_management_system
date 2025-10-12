package db

import (
    "fmt"
    "log"

    "github.com/jmoiron/sqlx"
    _ "github.com/lib/pq"
)

type Store struct {
    DB *sqlx.DB
}

// NewStore creates a new database connection
func NewStore(dbURL string) *Store {
    db, err := sqlx.Connect("postgres", dbURL)
    if err != nil {
        log.Fatalf("Failed to connect to Postgres: %v", err)
    }

    fmt.Println("Connected to Postgres")
    return &Store{DB: db}
}

// Close closes the database connection
func (s *Store) Close() {
    if s.DB != nil {
        s.DB.Close()
    }
}
