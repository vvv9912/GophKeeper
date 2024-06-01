package service

import (
	"GophKeeper/pkg/ShaHash"
	"GophKeeper/pkg/customErrors"
	"GophKeeper/pkg/logger"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

// ServiceData - структура для работы с данными пользователя.
type ServiceData struct {
	StoreData
}

// NewServiceData - конструктор структуры для работы с данными пользователя.
func NewServiceData(storeData StoreData) *ServiceData {
	return &ServiceData{StoreData: storeData}
}

// CreateCredentials - Создание пары логин/пароль.
func (s *ServiceData) CreateCredentials(ctx context.Context, userId int64, data []byte, name, description string) error {
	hash, err := s.createData(ctx, userId, data, name, description)
	if err != nil {
		return err
	}
	return s.StoreData.CreateCredentials(ctx, userId, data, name, description, hash)
}

// CreateCreditCard - Создание пары данные банковских карт.
func (s *ServiceData) CreateCreditCard(ctx context.Context, userId int64, data []byte, name, description string) error {
	hash, err := s.createData(ctx, userId, data, name, description)
	if err != nil {
		return err
	}
	return s.StoreData.CreateCreditCard(ctx, userId, data, name, description, hash)
}

// CreateFile - Создание произвольных данных.
func (s *ServiceData) CreateFile(ctx context.Context, userId int64, data []byte, name, description string) error {
	hash, err := s.createData(ctx, userId, data, name, description)
	if err != nil {
		return err
	}
	return s.StoreData.CreateFileData(ctx, userId, data, name, description, hash)
}

// createData - проверка правильности данных и расчет хэша.
func (s *ServiceData) createData(ctx context.Context, userId int64, data []byte, name, description string) (string, error) {
	var err error
	if data == nil || len(data) == 0 {
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
func (s *ServiceData) ChangeData(ctx context.Context, userId int64, lastTimeUpdate time.Time) ([]byte, error) {
	if userId == 0 {
		logger.Log.Error("userId is empty")
		return nil, customErrors.NewCustomError(nil, http.StatusBadRequest, "userId is empty")
	}
	data, err := s.StoreData.ChangeData(ctx, userId, lastTimeUpdate)
	if err != nil {
		return nil, err
	}
	resp, err := json.Marshal(data)
	if err != nil {
		err = customErrors.NewCustomError(err, http.StatusInternalServerError, "Marshal data error")
		return nil, err
	}
	return resp, nil
}
