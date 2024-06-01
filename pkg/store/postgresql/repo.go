package postgresql

import (
	"GophKeeper/pkg/logger"
	"GophKeeper/pkg/store"
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"time"
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

func (db *Database) createUserData(ctx context.Context, tx *sql.Tx, userId int64, dataId int64, dataType int, name, description, hash string) error {

	query := "INSERT INTO users_data (data_id,user_id, data_type, name, description,hash) VALUES ($1, $2,$3,$4,$5,$6)"

	res, err := db.db.Exec(query, dataId, userId, dataType, name, description, hash)
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

func (db *Database) createData(ctx context.Context, tx *sql.Tx, data []byte) (int64, error) {

	query := "INSERT INTO data (encrypt_data) VALUES ($1) RETURNING data_id"
	var id int64
	err := tx.QueryRowContext(ctx, query, data).Scan(&id)
	if err != nil {
		logger.Log.Error("Add data")
		return 0, err
	}

	return id, nil
}

func (db *Database) getDataByUserId(ctx context.Context, userId int64) ([]store.UsersData, error) {
	query := "SELECT user_data_id, data_id,users_id,data_type,name, description, hash, created_at,update_at,isDeleted FROM data WHERE user_id = $1 and isDeleted = false"
	row, err := db.db.QueryContext(ctx, query, userId)
	if err != nil {
		logger.Log.Error("Error while querying data", zap.Error(err))
		return nil, err
	}
	defer row.Close()

	var data []store.UsersData

	for row.Next() {
		var userDataId int64
		var DataId int64
		var dataType int
		var name string
		var description string
		var hash string
		var createdAt time.Time
		var updateAt time.Time
		var isDeleted bool

		err = row.Scan(&userDataId, &DataId, &userId, &dataType, &name, &description, &hash, &createdAt, &updateAt, &isDeleted)
		if err != nil {
			logger.Log.Error("Error getting data", zap.Error(err))
			return nil, err
		}
		data = append(data, store.UsersData{
			UserDataId:  userDataId,
			UserId:      userId,
			DataType:    dataType,
			Name:        name,
			Description: description,
			Hash:        hash,
			CreatedAt:   createdAt,
			UpdateAt:    updateAt,
			IsDeleted:   isDeleted,
		})
	}
	return data, err
}

func (db *Database) changeData(ctx context.Context, userId int64, lastTimeUpdate time.Time) ([]store.UsersData, error) {

	query := "SELECT user_data_id, name, description,data_type, hash, update_at,isDeleted FROM data WHERE user_id = $1 and update_at > $2"

	row, err := db.db.QueryContext(ctx, query, userId, lastTimeUpdate)
	if err != nil {
		logger.Log.Error("Error while querying data", zap.Error(err))
		return nil, err
	}
	defer row.Close()

	var data []store.UsersData

	for row.Next() {
		var userDataId int64
		var name string
		var description string
		var dataType int
		var hash string
		var updateAt time.Time
		var isDeleted bool

		err = row.Scan(&userDataId, &name, &description, &dataType, &hash, &updateAt, &isDeleted)
		if err != nil {
			logger.Log.Error("Error getting data", zap.Error(err))
			return nil, err
		}
		data = append(data, store.UsersData{
			UserDataId:  userDataId,
			Name:        name,
			Description: description,
			DataType:    dataType,
			Hash:        hash,
			UpdateAt:    updateAt,
			IsDeleted:   isDeleted,
		})
	}

	return data, nil
}
