package sqllite

import (
	"GophKeeper/pkg/customErrors"
	"GophKeeper/pkg/logger"
	"GophKeeper/pkg/store"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
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

type Database struct {
	db *sqlx.DB
}

func NewDatabase(db *sqlx.DB) *Database {
	return &Database{db: db}
}

func (db *Database) CreateBinaryFile(ctx context.Context, data []byte, userDataId int64, name, description, hash string, createdAt *time.Time, UpdateAt *time.Time, metaData *store.MetaData) error {
	tx, err := db.db.Begin()

	if err != nil {
		return err
	}
	defer func() {
		err = handleTransaction(tx, err)
	}()

	m, err := json.Marshal(metaData)
	if err != nil {
		err = customErrors.NewCustomError(err, http.StatusInternalServerError, "err json marshal metadata")
		return err
	}

	// возвращаем user_data_id
	dataId, err := db.createDataWithMeta(ctx, tx, data, m)
	if err != nil {
		err = customErrors.NewCustomError(err, http.StatusInternalServerError, "add file failed")
		return err
	}
	err = db.createUserData(ctx, tx, userDataId, dataId, TypeBinaryFile, name, description, hash, createdAt, UpdateAt)
	if err != nil {
		err = customErrors.NewCustomError(err, http.StatusInternalServerError, "add file failed")
		return err
	}

	return nil
}

func (db *Database) CreateFileData(ctx context.Context, data []byte, userDataId int64, name, description, hash string, createdAt *time.Time, UpdateAt *time.Time) error {
	// Получаем userDataId
	tx, err := db.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		err = handleTransaction(tx, err)
	}()
	dataId, err := db.createData(ctx, tx, data)
	if err != nil {
		return err
	}
	err = db.createUserData(ctx, tx, userDataId, dataId, TypeBinaryFile, name, description, hash, createdAt, UpdateAt)
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) CreateCredentials(ctx context.Context, data []byte, userDataId int64, name, description, hash string, createdAt *time.Time, UpdateAt *time.Time) error {
	// Получаем userDataId
	tx, err := db.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		err = handleTransaction(tx, err)
	}()
	dataId, err := db.createData(ctx, tx, data)
	if err != nil {
		return err
	}
	err = db.createUserData(ctx, tx, userDataId, dataId, TypeCredentials, name, description, hash, createdAt, UpdateAt)
	if err != nil {
		return err
	}
	return nil
}
func (db *Database) CreateCreditCard(ctx context.Context, data []byte, userDataId int64, name, description, hash string, createdAt *time.Time, UpdateAt *time.Time) error {
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
	err = db.createUserData(ctx, tx, userDataId, dataId, TypeCreditCardData, name, description, hash, createdAt, UpdateAt)
	if err != nil {
		return err
	}
	return nil
}
func (db *Database) GetMetaData(ctx context.Context, userDataId int64) (*store.MetaData, error) {
	query := `SELECT d.meta_data FROM users_data as u
		JOIN data as d on u.data_id = d.data_id
		WHERE user_data_id = ?`

	var metaData []byte
	if err := db.db.QueryRowContext(ctx, query, userDataId).Scan(&metaData); err != nil {
		err = customErrors.NewCustomError(err, http.StatusInternalServerError, "get file size failed")
		logger.Log.Error("get file size failed", zap.Error(err))
		return nil, err
	}

	var MetaData store.MetaData
	if err := json.Unmarshal(metaData, &MetaData); err != nil {
		err = customErrors.NewCustomError(err, http.StatusInternalServerError, "get file size failed")
		return nil, err
	}
	return &MetaData, nil
}

func (db *Database) GetInfoData(ctx context.Context, userDataId int64) (*store.UsersData, error) {
	usersData, err := db.getDataUserByUserId(ctx, userDataId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = customErrors.NewCustomError(err, http.StatusNotFound, "get data failed")
			return nil, err
		}
		err = customErrors.NewCustomError(err, http.StatusInternalServerError, "get data failed")
		return nil, err
	}
	return usersData, nil
}

func (db *Database) GetData(ctx context.Context, usersDataId int64) (*store.UsersData, *store.DataFile, error) {
	usersData, err := db.getDataUserByUserId(ctx, usersDataId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = customErrors.NewCustomError(err, http.StatusNotFound, "get data failed")
			return nil, nil, err
		}
		err = customErrors.NewCustomError(err, http.StatusInternalServerError, "get data failed")
		return nil, nil, err
	}

	data, err := db.getDataByDataId(ctx, usersData.DataId)
	if err != nil {
		err = customErrors.NewCustomError(err, http.StatusInternalServerError, "get data failed")
		return nil, nil, err
	}

	return usersData, data, nil
}
func (db *Database) getDataUserByUserId(ctx context.Context, userDataId int64) (*store.UsersData, error) {
	query := "SELECT user_data_id, data_id,data_type,name, description, hash, created_at,update_at,is_deleted FROM users_data WHERE user_data_id = $1 and is_deleted = false "
	row := db.db.QueryRowContext(ctx, query, userDataId)

	var data store.UsersData

	err := row.Scan(
		&data.UserDataId,
		&data.DataId,
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

func (db *Database) getDataByDataId(ctx context.Context, dataId int64) (*store.DataFile, error) {
	query := "SELECT data_id, encrypt_data FROM data WHERE data_id = $1"
	var data store.DataFile
	err := db.db.QueryRow(query, dataId).Scan(&data.DataId, &data.EncryptData)
	if err != nil {
		//if err == sql.ErrNoRows {
		//	err
		//}
		//todo
		logger.Log.Error("Get data by id", zap.Error(err))
		return nil, err
	}
	return &data, nil
}

// createData - добавление пользовательских данных.
func (db *Database) createDataWithMeta(ctx context.Context, tx *sql.Tx, data []byte, metaData []byte) (int64, error) {

	query := "INSERT INTO data (encrypt_data, meta_data) VALUES (?,?) RETURNING data_id"
	var id int64
	err := tx.QueryRowContext(ctx, query, data, metaData).Scan(&id)
	if err != nil {
		logger.Log.Error("Add data")
		return 0, err
	}

	return id, nil
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
func (db *Database) createUserData(ctx context.Context, tx *sql.Tx, userDataId int64, dataId int64, dataType int, name, description, hash string, createdAt, updateAt *time.Time) error {

	query := "INSERT INTO users_data (user_data_id,data_id, data_type, name, description,hash, created_at, update_at) VALUES (?, ?,?,?,?,?,?,?)"

	res, err := tx.Exec(query, userDataId, dataId, dataType, name, description, hash, createdAt, updateAt)
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

func (db *Database) createBinaryFile(ctx context.Context, tx *sql.Tx, data []byte, userDataId int64, name, description, hash string, createdAt, updateAt *time.Time, metaData *store.MetaData) error {
	m, err := json.Marshal(metaData)
	if err != nil {
		err = customErrors.NewCustomError(err, http.StatusInternalServerError, "err json marshal metadata")
		return err
	}
	// возвращаем user_data_id
	dataId, err := db.createDataWithMeta(ctx, tx, data, m)
	if err != nil {
		return err
	}
	err = db.createUserData(ctx, tx, userDataId, dataId, TypeBinaryFile, name, description, hash, createdAt, updateAt)
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) UpdateData(ctx context.Context, userDataId int64, data []byte, hash string, updateAt *time.Time) error {
	tx, err := db.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		err = handleTransaction(tx, err)
	}()

	err = db.updateData(ctx, tx, userDataId, data, hash, updateAt)
	if err != nil {
		return err
	}
	return nil
}
func (db *Database) UpdateDataBinary(ctx context.Context, userDataId int64, data []byte, hash string, updateAt *time.Time, metaData []byte) error {
	tx, err := db.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		err = handleTransaction(tx, err)
	}()

	err = db.updateData(ctx, tx, userDataId, data, hash, updateAt)
	if err != nil {
		return err
	}

	err = db.updateMetaData(ctx, tx, userDataId, metaData)
	if err != nil {
		return err
	}
	return nil
}

// updateData - обновление пользовательских данных
func (db *Database) updateData(ctx context.Context, tx *sql.Tx, userDataId int64, data []byte, hash string, updateAt *time.Time) error {

	queryBlock1 := "SELECT data_id FROM users_data WHERE user_data_id = ? "

	// Изменение данных в таблице users_data
	query1 := "UPDATE users_data SET  hash=?, update_at=? WHERE user_data_id = ?"
	// Изменение данных в таблице data
	query2 := "UPDATE data SET encrypt_data = ? where data_id=?"

	// Получим dataId и заблокируем на изменение табилцу
	var dataId int64
	err := tx.QueryRowContext(ctx, queryBlock1, userDataId).Scan(&dataId)
	if err != nil {
		logger.Log.Error("Error while querying data", zap.Error(err))
		return err
	}

	//// Заблокируем таблицу data
	//_, err = tx.QueryContext(ctx, queryBlock2, dataId)
	//if err != nil {
	//	logger.Log.Error("Error while querying data", zap.Error(err))
	//	return err
	//}

	// Вставляем новые данные
	res, err := tx.ExecContext(ctx, query2, data, dataId)
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

	res, err = tx.ExecContext(ctx, query1, hash, updateAt, userDataId)
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

	// Возвращаем updateAt
	return nil
}
func (db *Database) updateMetaData(ctx context.Context, tx *sql.Tx, userDataId int64, metaData []byte) error {

	q := "UPDATE data SET meta_data = ? where data_id = ?"

	_, err := tx.ExecContext(ctx, q, metaData, userDataId)
	if err != nil {
		logger.Log.Error("Error updating meta data", zap.Error(err))
		return err
	}
	return nil
}
func handleTransaction(tx *sql.Tx, err error) error {
	if err != nil {
		newErr := tx.Rollback()
		if newErr != nil {
			logger.Log.Error("Error while rollback", zap.Error(newErr))
			err = errors.Join(err, newErr)
		}
	} else {
		newErr := tx.Commit()
		if newErr != nil {
			logger.Log.Error("Error while commit", zap.Error(newErr))
			err = errors.Join(err, newErr)
		}
	}
	return err
}
