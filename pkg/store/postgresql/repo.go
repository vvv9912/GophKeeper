package postgresql

import (
	"GophKeeper/pkg/logger"
	"context"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type Database struct {
	db *sqlx.DB
}

func NewDatabase(db *sqlx.DB) *Database {
	return &Database{db: db}
}
func (db *Database) createUser(ctx context.Context, login, password string) (int64, error) {
	query := "INSERT INTO users (login, password) VALUES ($1, $2) RETURNING id"
	var id int64
	err := db.db.QueryRowxContext(ctx, query, login, password).Scan(&id)
	if err != nil {
		logger.Log.Error("Error while creating user", zap.String("login", login), zap.String("password", password))
		return 0, err
	}
	return id, nil
}
