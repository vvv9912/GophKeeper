package postgresql

import (
	"GophKeeper/pkg/logger"
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
)

func (db *Database) CreateUser(ctx context.Context, login, password string) (int64, error) {
	id, err := db.createUser(ctx, login, password)
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

func (db *Database) GetUserId(ctx context.Context, login, password string) (int64, error) {
	row := db.db.QueryRowxContext(ctx, "SELECT user_id FROM users WHERE login = $1 AND password = $2", login, password)
	var id int64
	err := row.Scan(&id)
	if err != nil {
		logger.Log.Error("Error while getting user id", zap.String("login", login), zap.String("password", password), zap.Error(err))
		return 0, err
	}
	return id, nil
}
