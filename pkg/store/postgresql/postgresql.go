package postgresql

import (
	"GophKeeper/pkg/logger"
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
)

func (db *Database) CreateUser(login, password string) (int64, error) {
	id, err := db.createUser(context.Background(), login, password)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				logger.Log.Error("Error creating user, duplicate", zap.String("error", pgErr.Message))
				return 0, nil
			}
			if pgErr.Code == "22001" {
				logger.Log.Error("Error creating user, long", zap.String("error", pgErr.Message))
				return 0, nil
			}
		}
		logger.Log.Error("Error while creating user", zap.String("login", login), zap.String("password", password))
		return 0, err
	}

	return id, nil
}
