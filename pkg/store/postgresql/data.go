package postgresql

import (
	"GophKeeper/pkg/customErrors"
	"GophKeeper/pkg/logger"
	"GophKeeper/pkg/store"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"go.uber.org/zap"
	"net/http"
	"time"
)

// CreateCredentials - Создание пары логин/пароль.
func (db *Database) CreateCredentials(ctx context.Context, userId int64, data []byte, name, description, hash string) (*store.UsersData, error) {
	tx, err := db.db.Begin()
	if err != nil {
		return nil, err
	}

	defer func() {
		err = handleTransaction(tx, err)
	}()

	dataId, err := db.createData(ctx, tx, data)
	if err != nil {
		err = customErrors.NewCustomError(err, http.StatusInternalServerError, "add credentials failed")
		return nil, err
	}
	userData, err := db.createUserData(ctx, tx, userId, dataId, TypeCredentials, name, description, hash)
	if err != nil {
		err = customErrors.NewCustomError(err, http.StatusInternalServerError, "add credentials failed")
		return nil, err
	}

	return userData, nil
}

// CreateCreditCard - Создание пары данные банковских карт.
func (db *Database) CreateCreditCard(ctx context.Context, userId int64, data []byte, name, description, hash string) (*store.UsersData, error) {
	tx, err := db.db.Begin()
	if err != nil {
		return nil, err
	}

	defer func() {
		err = handleTransaction(tx, err)
	}()

	dataId, err := db.createData(ctx, tx, data)
	if err != nil {
		err = customErrors.NewCustomError(err, http.StatusInternalServerError, "add credit card failed")
		return nil, err
	}
	userData, err := db.createUserData(ctx, tx, userId, dataId, TypeCreditCardData, name, description, hash)
	if err != nil {
		err = customErrors.NewCustomError(err, http.StatusInternalServerError, "add credit card failed")
		return nil, err
	}

	return userData, nil
}

// CreateFileData - Создание произвольных данных.
func (db *Database) CreateFileDataChunks(ctx context.Context, userId int64, data []byte, name, description, hash string, metaData *store.MetaData) (*store.UsersData, error) {
	tx, err := db.db.Begin()
	if err != nil {
		return nil, err
	}

	defer func() {
		err = handleTransaction(tx, err)
	}()

	m, err := json.Marshal(metaData)
	if err != nil {
		err = customErrors.NewCustomError(err, http.StatusInternalServerError, "err json marshal metadata")
		return nil, err
	}

	// возвращаем user_data_id
	dataId, err := db.createDataWithMeta(ctx, tx, data, m)
	if err != nil {
		err = customErrors.NewCustomError(err, http.StatusInternalServerError, "add file failed")
		return nil, err
	}
	userData, err := db.createUserData(ctx, tx, userId, dataId, TypeBinaryFile, name, description, hash)
	if err != nil {
		err = customErrors.NewCustomError(err, http.StatusInternalServerError, "add file failed")
		return nil, err
	}

	return userData, nil
}

// CreateFileData - Создание произвольных данных.
func (db *Database) CreateFileData(ctx context.Context, userId int64, data []byte, name, description, hash string) (*store.UsersData, error) {
	tx, err := db.db.Begin()
	if err != nil {
		return nil, err
	}

	defer func() {
		err = handleTransaction(tx, err)
	}()

	// возвращаем user_data_id
	dataId, err := db.createData(ctx, tx, data)
	if err != nil {
		err = customErrors.NewCustomError(err, http.StatusInternalServerError, "add file failed")
		return nil, err
	}
	userData, err := db.createUserData(ctx, tx, userId, dataId, TypeBinaryFile, name, description, hash)
	if err != nil {
		err = customErrors.NewCustomError(err, http.StatusInternalServerError, "add file failed")
		return nil, err
	}

	return userData, nil
}

func (db *Database) ChangeData(ctx context.Context, userId int64, userDataId int64, lastTimeUpdate time.Time) (bool, error) {

	return db.changeData(ctx, userId, userDataId, lastTimeUpdate)
}

func (db *Database) ChangeAllData(ctx context.Context, userId int64, lastTimeUpdate time.Time) ([]store.UsersData, error) {
	data, err := db.changeAllData(ctx, userId, lastTimeUpdate)
	if err != nil {
		err = customErrors.NewCustomError(err, http.StatusInternalServerError, "get change data failed")
		return nil, err
	}
	return data, nil
}

func (db *Database) GetFileSize(ctx context.Context, userId int64, userDataId int64) (int64, error) {

	metaData, err := db.GetMetaData(ctx, userId, userDataId)
	if err != nil {
		return 0, err
	}
	return metaData.Size, nil
}

func (db *Database) GetMetaData(ctx context.Context, userId, userDataId int64) (*store.MetaData, error) {
	query := `SELECT d.meta_data FROM users_data as u
		JOIN data as d on u.data_id = d.data_id
		WHERE u.user_data_id = $1 and u.user_id = $2`

	var metaData []byte
	if err := db.db.QueryRowContext(ctx, query, userDataId, userId).Scan(&metaData); err != nil {
		err = customErrors.NewCustomError(err, http.StatusInternalServerError, "get file size failed")
		return nil, err
	}

	var MetaData store.MetaData
	if err := json.Unmarshal(metaData, &MetaData); err != nil {
		err = customErrors.NewCustomError(err, http.StatusInternalServerError, "get file size failed")
		return nil, err
	}
	return &MetaData, nil
}

func (db *Database) GetData(ctx context.Context, userId int64, usersDataId int64) (*store.UsersData, *store.DataFile, error) {

	tx, err := db.db.Begin()
	if err != nil {
		logger.Log.Error("Error while begin transaction", zap.Error(err))
		return nil, nil, err
	}
	defer func() {
		err = handleTransaction(tx, err)
	}()

	usersData, err := db.getDataUserByUserId(ctx, tx, userId, usersDataId)
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

func (db *Database) UpdateData(ctx context.Context, userId, userDataId int64, data []byte, hash string) (*store.UsersData, error) {
	tx, err := db.db.Begin()
	if err != nil {
		logger.Log.Error("Error while begin transaction", zap.Error(err))
		return nil, err
	}
	defer func() {
		err = handleTransaction(tx, err)
	}()

	err = db.updateData(ctx, tx, userId, userDataId, data, hash)
	if err != nil {
		err = customErrors.NewCustomError(err, http.StatusInternalServerError, "update data failed")
		return nil, err
	}
	usersData, err := db.getDataUserByUserId(ctx, tx, userId, userDataId)
	if err != nil {
		err = customErrors.NewCustomError(err, http.StatusInternalServerError, "update data failed")
		return nil, err
	}

	return usersData, nil
}

func (db *Database) RemoveData(ctx context.Context, userId, usersDataId int64) error {
	err := db.removeData(ctx, userId, usersDataId)
	if err != nil {
		err = customErrors.NewCustomError(err, http.StatusInternalServerError, "remove data failed")
		return err
	}
	return nil
}

func (db *Database) GetListData(ctx context.Context, userId int64) ([]store.UsersData, error) {
	data, err := db.getListData(ctx, userId)
	if err != nil {
		err = customErrors.NewCustomError(err, http.StatusInternalServerError, "remove all data failed")
		return nil, err
	}
	return data, nil
}

// UpdateBinaryFile - Создание произвольных данных.
func (db *Database) UpdateBinaryFile(ctx context.Context, userId int64, userDataId int64, data []byte, hash string, metaData []byte) (*store.UsersData, error) {
	tx, err := db.db.Begin()
	if err != nil {
		return nil, err
	}

	defer func() {
		err = handleTransaction(tx, err)
	}()

	err = db.updateData(ctx, tx, userId, userDataId, data, hash)
	if err != nil {
		err = customErrors.NewCustomError(err, http.StatusInternalServerError, "update data failed")
		return nil, err
	}
	usersData, err := db.getDataUserByUserId(ctx, tx, userId, userDataId)
	if err != nil {
		err = customErrors.NewCustomError(err, http.StatusInternalServerError, "update data failed")
		return nil, err
	}
	if err := db.updateMetaData(ctx, tx, usersData.DataId, metaData); err != nil {
		err = customErrors.NewCustomError(err, http.StatusInternalServerError, "update data failed")
		return nil, err
	}

	return usersData, nil
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
