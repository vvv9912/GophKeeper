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
	hash, err := s.createData(ctx, userId, data, name, description)
	if err != nil {
		return err
	}
	return s.StoreData.CreateCredentials(ctx, userId, data, name, description, hash)
}
func (s *ServiceData) CreateCreditCard(ctx context.Context, userId int64, data []byte, name, description string) error {
	hash, err := s.createData(ctx, userId, data, name, description)
	if err != nil {
		return err
	}
	return s.StoreData.CreateCreditCard(ctx, userId, data, name, description, hash)
}
func (s *ServiceData) CreateFile(ctx context.Context, userId int64, data []byte, name, description string) error {
	hash, err := s.createData(ctx, userId, data, name, description)
	if err != nil {
		return err
	}
	return s.StoreData.CreateFileData(ctx, userId, data, name, description, hash)
}
func (s *ServiceData) createData(ctx context.Context, userId int64, data []byte, name, description string) (string, error) {
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
	if userId == 0 {
		logger.Log.Error("userId is empty")
		err = errors.Join(err, customErrors.NewCustomError(nil, http.StatusBadRequest, "userId is empty"))
	}
	//todo в вывод добавить unwrap
	if err != nil {
		return "", err
	}
	// Считаем хэш полученных данных
	hash := ShaHash.Sha256Hash(data)
	return hash, err
}
