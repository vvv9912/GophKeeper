package service

import (
	"GophKeeper/pkg/ShaHash"
	"GophKeeper/pkg/customErrors"
	"GophKeeper/pkg/logger"
	"GophKeeper/pkg/store"
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
func (s *ServiceData) CreateCredentials(ctx context.Context, userId int64, data []byte, name, description string) (*RespData, error) {
	hash, err := s.createData(ctx, userId, data, name, description)
	if err != nil {
		return nil, err
	}
	userDataId, err := s.StoreData.CreateCredentials(ctx, userId, data, name, description, hash)
	if err != nil {
		return nil, err
	}
	resp := &RespData{
		UserDataId: userDataId,
		Hash:       hash,
	}

	return resp, nil
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

func (s *ServiceData) GetData(ctx context.Context, userId int64, userDataId int64) ([]byte, error) {
	if userId == 0 {
		logger.Log.Error("userId is empty")
		return nil, customErrors.NewCustomError(nil, http.StatusBadRequest, "userId is empty")
	}
	usersData, data, err := s.StoreData.GetData(ctx, userId, userDataId)
	if err != nil {
		return nil, err
	}
	type Data struct {
		InfoUsersData *store.UsersData `json:"infoUsersData"`
		EncryptData   *store.DataFile  `json:"encryptData"`
	}
	resp := Data{
		InfoUsersData: usersData,
		EncryptData:   data,
	}
	response, err := json.Marshal(resp)
	if err != nil {
		err = customErrors.NewCustomError(err, http.StatusInternalServerError, "Marshal data error")
		return nil, err
	}
	return response, nil
}

func (s *ServiceData) UpdateData(ctx context.Context, userId int64, usersData *store.UpdateUsersData, data []byte) error {
	if userId == 0 {
		logger.Log.Error("userId is empty")
		return customErrors.NewCustomError(nil, http.StatusBadRequest, "userId is empty")
	}
	usersData.UserId = userId

	// todo проверка, если данные уже обновлены с другого устр-ва
	err := s.StoreData.UpdateData(ctx, usersData, data)
	if err != nil {
		return err
	}
	return nil

}
func (s *ServiceData) RemoveData(ctx context.Context, userId, userDataId int64) error {
	if userId == 0 {
		logger.Log.Error("userId is empty")
		return customErrors.NewCustomError(nil, http.StatusBadRequest, "userId is empty")
	}
	if userDataId == 0 {
		logger.Log.Error("userDataId is empty")
		return customErrors.NewCustomError(nil, http.StatusBadRequest, "userDataId is empty")
	}

	err := s.StoreData.RemoveData(ctx, userId, userDataId)
	if err != nil {
		return err
	}
	return nil

}
