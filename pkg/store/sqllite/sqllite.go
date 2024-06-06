package sqllite

import (
	"GophKeeper/pkg/logger"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

var (
	TypeCredentials    = 1
	TypeCreditCardData = 2
	TypeFile           = 3
)

type Database struct {
	db *sqlx.DB
}

func NewDatabase(db *sqlx.DB) *Database {
	return &Database{db: db}
}

func (db *Database) CreateFileData(ctx context.Context, data []byte, userDataId int64, name, description, hash string) error {
	// Получаем userDataId
	tx, err := db.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			newErr := tx.Rollback()
			if newErr != nil {
				err = errors.Join(err, newErr)
			}
		} else {
			newErr := tx.Commit()
			if newErr != nil {
				err = errors.Join(err, newErr)
			}
		}
	}()
	dataId, err := db.createData(ctx, tx, data)
	if err != nil {
		return err
	}
	err = db.createUserData(ctx, tx, userDataId, dataId, TypeFile, name, description, hash)
	if err != nil {
		return err
	}
	return nil
}
func (db *Database) CreateCredentials(ctx context.Context, data []byte, userDataId int64, name, description, hash string) error {
	// Получаем userDataId
	tx, err := db.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			newErr := tx.Rollback()
			if newErr != nil {
				err = errors.Join(err, newErr)
			}
		} else {
			newErr := tx.Commit()
			if newErr != nil {
				err = errors.Join(err, newErr)
			}
		}
	}()
	dataId, err := db.createData(ctx, tx, data)
	if err != nil {
		return err
	}
	err = db.createUserData(ctx, tx, userDataId, dataId, TypeCredentials, name, description, hash)
	if err != nil {
		return err
	}
	return nil
}

// createData - добавление пользовательских данных.
func (db *Database) createData(ctx context.Context, tx *sql.Tx, data []byte) (int64, error) {

	query := "INSERT INTO data (encrypt_data) VALUES (?) RETURNING data_id"
	var id int64
	err := tx.QueryRowContext(ctx, query, data).Scan(&id)
	if err != nil {
		logger.Log.Error("Add data")
		return 0, err
	}

	return id, nil
}
func (db *Database) createUserData(ctx context.Context, tx *sql.Tx, userDataId int64, dataId int64, dataType int, name, description, hash string) error {

	query := "INSERT INTO users_data (user_data_id,data_id, data_type, name, description,hash) VALUES (?, ?,?,?,?,?)"

	res, err := tx.Exec(query, userDataId, dataId, dataType, name, description, hash)
	if err != nil {
		logger.Log.Error("Add credentials error", zap.String("name", name), zap.String("description", description), zap.String("hash", hash), zap.Int("data_type", dataType))
		return err
	}
	r, err := res.RowsAffected()
	if err != nil {
		logger.Log.Error("Error getting rows affected", zap.Error(err))
		return err
	}
	if r == 0 {
		err = fmt.Errorf("Insert 0")
		return err
	}

	return err
}
