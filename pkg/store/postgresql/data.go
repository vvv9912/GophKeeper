package postgresql

import (
	"GophKeeper/pkg/customErrors"
	"GophKeeper/pkg/store"
	"context"
	"errors"
	"net/http"
	"time"
)

// CreateCredentials - Создание пары логин/пароль.
func (db *Database) CreateCredentials(ctx context.Context, userId int64, data []byte, name, description, hash string) error {
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
		err = customErrors.NewCustomError(err, http.StatusInternalServerError, "add credentials failed")
		return err
	}
	err = db.createUserData(ctx, tx, userId, dataId, TypeCredentials, name, description, hash)
	if err != nil {
		err = customErrors.NewCustomError(err, http.StatusInternalServerError, "add credentials failed")
		return err
	}

	return nil
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
	err = db.createUserData(ctx, tx, userId, dataId, TypeCreditCardData, name, description, hash)
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

	dataId, err := db.createData(ctx, tx, data)
	if err != nil {
		err = customErrors.NewCustomError(err, http.StatusInternalServerError, "add file failed")
		return err
	}
	err = db.createUserData(ctx, tx, userId, dataId, TypeFile, name, description, hash)
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
