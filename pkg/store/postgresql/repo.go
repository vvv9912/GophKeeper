package postgresql

import (
	"GophKeeper/pkg/customErrors"
	"GophKeeper/pkg/logger"
	"GophKeeper/pkg/store"
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"net/http"
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

// createData - добавление пользовательских данных.
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

// getDataByUserId - Получение информации о данных пользователя, которые не удалены
func (db *Database) getDataByUserId(ctx context.Context, userId int64) ([]store.UsersData, error) {
	query := "SELECT user_data_id, data_id,users_id,data_type,name, description, hash, created_at,update_at,is_deleted FROM users_data WHERE user_id = $1 and is_deleted = false FOR UPDATE "
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

// changeData - получение информации об изменненных данных
func (db *Database) changeData(ctx context.Context, userId int64, lastTimeUpdate time.Time) ([]store.UsersData, error) {

	query := "SELECT user_data_id, name, description,data_type, hash, update_at,is_deleted FROM users_data WHERE user_id = $1 and update_at > $2 FOR UPDATE "

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

// updateData - обновление пользовательских данных
func (db *Database) updateData(ctx context.Context, newData store.UsersData, data []byte) error {

	// Блокирующая транзацкция SELECT * FROM table_name WHERE condition FOR UPDATE;
	tx, err := db.db.Begin()
	if err != nil {
		logger.Log.Error("Error while begin transaction", zap.Error(err))
		return err
	}

	queryBlock1 := "SELECT data_id FROM users_data WHERE user_id = $1 and user_data_id = $2 for update"
	queryBlock2 := "SELECT data_id FROM data WHERE data_id = $1 for update"

	query1 := "UPDATE users_data SET name=$1, description=$2, hash=$3, update_at=$4 WHERE user_data_id = $6"
	query2 := "UPDATE data SET encrypt_data = $1 where data_id=$2"

	// Получим dataId и заблокируем на изменение табилцу
	var dataId int64
	err = tx.QueryRowContext(ctx, queryBlock1, newData.DataId, newData.DataId).Scan(&dataId)
	if err != nil {
		logger.Log.Error("Error while querying data", zap.Error(err))
		return err
	}

	if dataId != newData.DataId {
		//todo add err in repo
		err = customErrors.NewCustomError(nil, http.StatusBadRequest, "dataId != newData.DataId")
		return err
	}

	// Заблокируем таблицу data
	_, err = tx.QueryContext(ctx, queryBlock2, newData.DataId)
	if err != nil {
		logger.Log.Error("Error while querying data", zap.Error(err))
		return err
	}

	res, err := tx.ExecContext(ctx, query2, data, newData.DataId)
	if err != nil {
		logger.Log.Error("Error while querying data", zap.Error(err))
		return err
	}

	nr, err := res.RowsAffected()
	if err != nil {
		logger.Log.Error("Error while querying data", zap.Error(err))
		return err
	}
	if nr == 0 {
		//todo add err in repo
		return customErrors.NewCustomError(nil, http.StatusNotFound, "data not found")
	}

	res, err = tx.ExecContext(ctx, query1, newData.Name, newData.Description, newData.Hash, newData.UpdateAt, newData.UserDataId)
	if err != nil {
		logger.Log.Error("Error while querying data", zap.Error(err))
		return err
	}

	nr, err = res.RowsAffected()
	if err != nil {
		logger.Log.Error("Error while querying data", zap.Error(err))
		return err
	}
	if nr == 0 {
		//todo add err in repo
		return customErrors.NewCustomError(nil, http.StatusNotFound, "data not found")
	}

	return nil
}
