package sqllite

import "github.com/jmoiron/sqlx"

type Database struct {
	db *sqlx.DB
}

func NewDatabase(db *sqlx.DB) *Database {
	return &Database{db: db}
}
