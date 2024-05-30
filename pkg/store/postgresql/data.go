package postgresql

import (
	"GophKeeper/pkg/customErrors"
	"context"
	"errors"
	"net/http"
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
