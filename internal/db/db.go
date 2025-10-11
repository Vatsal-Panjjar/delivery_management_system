package db

import (
    "github.com/jmoiron/sqlx"
    _ "github.com/lib/pq"
)

var DB *sqlx.DB

func ConnectWithURL(url string) (*sqlx.DB, error) {
    var err error
    DB, err = sqlx.Connect("postgres", url)
    return DB, err
}
