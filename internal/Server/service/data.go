package service

import (
	"GophKeeper/pkg/ShaHash"
	"GophKeeper/pkg/customErrors"
	"GophKeeper/pkg/logger"
	"context"
	"errors"
	"net/http"
)

type ServiceData struct {
	StoreData
}

func NewServiceData(storeData StoreData) *ServiceData {
	return &ServiceData{StoreData: storeData}
}

func (s *ServiceData) CreateCredentials(ctx context.Context, userId int64, data []byte, name, description string) error {
	var err error
	if data == nil {
		logger.Log.Error("data is empty")
		err = errors.Join(err, customErrors.NewCustomError(nil, http.StatusBadRequest, "data is empty"))
	}
	if name == "" {
		logger.Log.Error("name is empty")
		err = errors.Join(err, customErrors.NewCustomError(nil, http.StatusBadRequest, "name is empty"))
	}
	if description == "" {
		logger.Log.Error("description is empty")
		err = errors.Join(err, customErrors.NewCustomError(nil, http.StatusBadRequest, "description is empty"))
	}
	//todo в вывод добавить unwrap
	if err != nil {
		return err
	}
	// Считаем хэш полученных данных
	hash := ShaHash.Sha256Hash(data)
	err = s.StoreData.CreateCredentials(ctx, userId, data, name, description, hash)
	return nil
}
