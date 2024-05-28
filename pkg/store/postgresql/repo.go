package postgresql

import (
	"github.com/jackc/pgx/v5"
)

type Database struct {
	pgx *pgx.Conn
}

func NewDatabase(db *pgx.Conn) *Database {
	return &Database{pgx: db}
}
