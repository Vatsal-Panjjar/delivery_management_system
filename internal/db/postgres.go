package db

import (
    "github.com/jmoiron/sqlx"
    _ "github.com/lib/pq"
)

func ConnectPostgres(url string) (*sqlx.DB, error) {
    db, err := sqlx.Connect("postgres", url)
    if err != nil {
        return nil, err
    }
    return db, nil
}
