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

const (
	TypeCredentials    = 1
	TypeCreditCardData = 2
	TypeBinaryFile     = 3
	TypeTxt            = 4
)

var DataType = map[int]string{
	1: "Credentials",
	2: "CreditCardData",
	3: "File",
}

type Database struct {
	db *sqlx.DB
}

func NewDatabase(db *sqlx.DB) *Database {
	return &Database{db: db}
}
func (db *Database) createUser(ctx context.Context, login, password string) (int64, error) {
	q := "INSERT INTO users (login, password) VALUES ($1, $2) RETURNING user_id"
	var id int64
	err := db.db.QueryRowxContext(ctx, q, login, password).Scan(&id)
	if err != nil {
		logger.Log.Error("Error while creating user", zap.String("login", login), zap.String("password", password), zap.Error(err))
		return 0, err
	}
	return id, nil
}

func (db *Database) createUserData(ctx context.Context, tx *sql.Tx, userId int64, dataId int64, dataType int, name, description, hash string) (*store.UsersData, error) {

	q := "INSERT INTO users_data (data_id,user_id, data_type, name, description,hash) VALUES ($1, $2,$3,$4,$5,$6) RETURNING user_data_id,created_at,update_at"

	var usersData store.UsersData

	err := db.db.QueryRowContext(ctx, q, dataId, userId, dataType, name, description, hash).Scan(&usersData.UserDataId, &usersData.CreatedAt, &usersData.UpdateAt)
	if err != nil {
		logger.Log.Error("Add credentials error", zap.String("name", name), zap.String("description", description), zap.String("hash", hash), zap.Int("data_type", dataType))
		return nil, err
	}
	usersData.UserId = userId
	usersData.DataId = dataId
	usersData.DataType = dataType
	usersData.Name = name
	usersData.Description = description
	usersData.Hash = hash

	return &usersData, err
}

// createData - добавление пользовательских данных.
func (db *Database) createDataWithMeta(ctx context.Context, tx *sql.Tx, data []byte, metaData []byte) (int64, error) {

	q := "INSERT INTO data (encrypt_data, meta_data) VALUES ($1,$2) RETURNING data_id"
	var id int64
	err := tx.QueryRowContext(ctx, q, data, metaData).Scan(&id)
	if err != nil {
		logger.Log.Error("Add data")
		return 0, err
	}

	return id, nil
}

// createData - добавление пользовательских данных.
func (db *Database) createData(ctx context.Context, tx *sql.Tx, data []byte) (int64, error) {

	q := "INSERT INTO data (encrypt_data) VALUES ($1) RETURNING data_id"
	var id int64
	err := tx.QueryRowContext(ctx, q, data).Scan(&id)
	if err != nil {
		logger.Log.Error("Add data")
		return 0, err
	}

	return id, nil
}

func (db *Database) getDataByDataId(ctx context.Context, dataId int64) (*store.DataFile, error) {
	q := "SELECT data_id, encrypt_data FROM data WHERE data_id = $1"
	var data store.DataFile
	err := db.db.QueryRow(q, dataId).Scan(&data.DataId, &data.EncryptData)
	if err != nil {
		logger.Log.Error("Get data by id", zap.Error(err))
		return nil, err
	}
	return &data, nil
}

// getDataByUserId - Получение информации о данных пользователя, которые не удалены
func (db *Database) getDataUserByUserId(ctx context.Context, tx *sql.Tx, userId int64, userDataId int64) (*store.UsersData, error) {
	q := "SELECT user_data_id, data_id,user_id,data_type,name, description, hash, created_at,update_at,is_deleted FROM users_data WHERE user_data_id = $1 and is_deleted = false and user_id = $2 FOR UPDATE "
	row := tx.QueryRowContext(ctx, q, userDataId, userId)

	var data store.UsersData

	err := row.Scan(
		&data.UserDataId,
		&data.DataId,
		&data.UserId,
		&data.DataType,
		&data.Name,
		&data.Description,
		&data.Hash,
		&data.CreatedAt,
		&data.UpdateAt,
		&data.IsDeleted,
	)

	if err != nil {
		logger.Log.Error("Error getting data", zap.Error(err))
		return nil, err
	}

	return &data, err
}

// getDataByUserId - Получение информации о данных пользователя
func (db *Database) getListData(ctx context.Context, userId int64) ([]store.UsersData, error) {
	q := "SELECT user_data_id, data_id,user_id,data_type,name, description, hash, created_at,update_at,is_deleted FROM users_data WHERE is_deleted = false and user_id = $1 FOR UPDATE "
	rows, err := db.db.QueryContext(ctx, q, userId)
	if err != nil {
		logger.Log.Error("Error getting data", zap.Error(err))
		return nil, err
	}
	var data []store.UsersData

	for rows.Next() {
		var d store.UsersData
		err = rows.Scan(
			&d.UserDataId,
			&d.DataId,
			&d.UserId,
			&d.DataType,
			&d.Name,
			&d.Description,
			&d.Hash,
			&d.CreatedAt,
			&d.UpdateAt,
			&d.IsDeleted,
		)
		if err != nil {
			logger.Log.Error("Error getting data", zap.Error(err))
			return nil, err
		}
		data = append(data, d)
	}

	return data, err
}

// changeData - получение информации об изменненнu данных
func (db *Database) changeData(ctx context.Context, userId int64, userDataId int64, lastTimeUpdate time.Time) (bool, error) {
	q := "SELECT EXISTS(SELECT 1 FROM users_data WHERE user_id = $1 AND user_data_id = $2 AND update_at > $3)"
	var exist bool
	err := db.db.QueryRowContext(ctx, q, userId, userDataId, lastTimeUpdate).Scan(&exist)
	if err != nil {
		logger.Log.Error("Error while querying data", zap.Error(err))
		return false, err
	}

	return exist, nil
}

// changeData - получение информации об изменненных данных
func (db *Database) changeAllData(ctx context.Context, userId int64, lastTimeUpdate time.Time) ([]store.UsersData, error) {

	q := "SELECT user_data_id, name, description,data_type, hash, update_at,is_deleted FROM users_data WHERE user_id = $1 and update_at > $2 FOR UPDATE "

	rows, err := db.db.QueryContext(ctx, q, userId, lastTimeUpdate)
	if err != nil {
		logger.Log.Error("Error while querying data", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var data []store.UsersData

	for rows.Next() {
		var userData store.UsersData
		err = rows.Scan(&userData.UserDataId, &userData.Name, &userData.Description, &userData.DataType, &userData.Hash, &userData.UpdateAt, &userData.IsDeleted)
		if err != nil {
			logger.Log.Error("Error getting data", zap.Error(err))
			return nil, err
		}
		data = append(data, userData)
	}

	return data, nil
}
func (db *Database) updateMetaData(ctx context.Context, tx *sql.Tx, dataId int64, meta []byte) error {
	q := "UPDATE data SET meta_data = $1 where data_id = $2"

	_, err := tx.ExecContext(ctx, q, meta, dataId)
	if err != nil {
		logger.Log.Error("Error updating meta data", zap.Error(err))
		return err
	}
	return nil
}

// updateData - обновление пользовательских данных
func (db *Database) updateData(ctx context.Context, tx *sql.Tx, userId, userDataId int64, data []byte, hash string) error {

	// Блокировка таблицы users_data и получение dataid

	queryBlock1 := "SELECT data_id FROM users_data WHERE user_id = $1 and user_data_id = $2 FOR UPDATE "

	// Изменение данных в таблице users_data
	queryChangeUsersData := "UPDATE users_data SET  hash=$1, update_at=$2 WHERE user_data_id = $3"
	// Изменение данных в таблице data
	queryChangeData := "UPDATE data SET encrypt_data = $1 where data_id=$2"

	// Получим dataId и заблокируем на изменение табилцу
	var dataId int64
	err := tx.QueryRowContext(ctx, queryBlock1, userId, userDataId).Scan(&dataId)
	if err != nil {
		logger.Log.Error("Error while querying data", zap.Error(err))
		return err
	}

	res, err := tx.ExecContext(ctx, queryChangeData, data, dataId)
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
		return customErrors.NewCustomError(nil, http.StatusNotFound, "data not found")
	}
	updateAt := time.Now().UTC()

	res, err = tx.ExecContext(ctx, queryChangeUsersData, hash, updateAt, userDataId)
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
		return fmt.Errorf("data not found")
	}

	// Возвращаем updateAt
	return nil
}

// removeData - удаление пользовательских данных
func (db *Database) removeData(ctx context.Context, userId int64, usersDataId int64) error {

	// Удаление данных в таблице users_data
	queryRemoveUsersData := "UPDATE users_data SET is_deleted=$1 WHERE user_data_id = $2 and user_id = $3"

	res, err := db.db.ExecContext(ctx, queryRemoveUsersData, true, usersDataId, userId)
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
		return customErrors.NewCustomError(nil, http.StatusNotFound, "data not found")
	}

	return nil
}
