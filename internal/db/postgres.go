package db


import (
"fmt"
"github.com/jmoiron/sqlx"
_ "github.com/lib/pq"
"time"
)


type Postgres struct{ *sqlx.DB }


func NewPostgres(user, pass, dbname, host string, port int) (*Postgres, error) {
dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", user, pass, host, port, dbname)
db, err := sqlx.Connect("postgres", dsn)
if err != nil {
return nil, err
}
db.SetMaxOpenConns(25)
db.SetMaxIdleConns(5)
db.SetConnMaxLifetime(5 * time.Minute)
return &Postgres{DB: db}, nil
}
