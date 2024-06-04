package postgresql

import (
	"GophKeeper/pkg/customErrors"
	"GophKeeper/pkg/store"
	"context"
	"database/sql"
	"errors"
	"net/http"
	"time"
)

// CreateCredentials - Создание пары логин/пароль.
func (db *Database) CreateCredentials(ctx context.Context, userId int64, data []byte, name, description, hash string) (int64, error) {
	tx, err := db.db.Begin()
	if err != nil {
		return 0, err
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
		err = customErrors.NewCustomError(err, http.StatusInternalServerError, "add credentials failed")
		return 0, err
	}
	userDataId, err := db.createUserData(ctx, tx, userId, dataId, TypeCredentials, name, description, hash)
	if err != nil {
		err = customErrors.NewCustomError(err, http.StatusInternalServerError, "add credentials failed")
		return 0, err
	}

	return userDataId, nil
}

// CreateCreditCard - Создание пары данные банковских карт.
func (db *Database) CreateCreditCard(ctx context.Context, userId int64, data []byte, name, description, hash string) error {
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
		err = customErrors.NewCustomError(err, http.StatusInternalServerError, "add credit card failed")
		return err
	}
	_, err = db.createUserData(ctx, tx, userId, dataId, TypeCreditCardData, name, description, hash)
	if err != nil {
		err = customErrors.NewCustomError(err, http.StatusInternalServerError, "add credit card failed")
		return err
	}

	return nil
}

// CreateFileData - Создание произвольных данных.
func (db *Database) CreateFileData(ctx context.Context, userId int64, data []byte, name, description, hash string) error {
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

	// возвращаем user_data_id

	dataId, err := db.createData(ctx, tx, data)
	if err != nil {
		err = customErrors.NewCustomError(err, http.StatusInternalServerError, "add file failed")
		return err
	}
	_, err = db.createUserData(ctx, tx, userId, dataId, TypeFile, name, description, hash)
	if err != nil {
		err = customErrors.NewCustomError(err, http.StatusInternalServerError, "add file failed")
		return err
	}

	return nil
}

func (db *Database) ChangeData(ctx context.Context, userId int64, lastTimeUpdate time.Time) ([]store.UsersData, error) {
	data, err := db.changeData(ctx, userId, lastTimeUpdate)
	if err != nil {
		err = customErrors.NewCustomError(err, http.StatusInternalServerError, "get change data failed")
		return nil, err
	}
	return data, nil
}

func (db *Database) GetData(ctx context.Context, userId int64, usersDataId int64) (*store.UsersData, *store.DataFile, error) {
	usersData, err := db.getDataUserByUserId(ctx, userId, usersDataId)
	if err != nil {
		if err == sql.ErrNoRows {
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

func (db *Database) UpdateData(ctx context.Context, updateData *store.UpdateUsersData, data []byte) error {
	err := db.updateData(ctx, updateData, data)
	if err != nil {
		err = customErrors.NewCustomError(err, http.StatusInternalServerError, "update data failed")
		return err
	}
	return nil
}

func (db *Database) RemoveData(ctx context.Context, userId, usersDataId int64) error {
	err := db.removeData(ctx, userId, usersDataId)
	if err != nil {
		err = customErrors.NewCustomError(err, http.StatusInternalServerError, "remove data failed")
		return err
	}
	return nil
}
