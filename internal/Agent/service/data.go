package service

import (
	"GophKeeper/internal/Agent/server"
	"GophKeeper/pkg/store/sqllite"
	"context"
	"fmt"
)

// CreateCredentials - Создание данных логин/пароль
func (s *Service) CreateCredentials(ctx context.Context, data *server.ReqData) error {
	// Получение jwt токена
	if err := s.setJwtToken(ctx); err != nil {
		return err
	}

	if err := s.encryptData(data); err != nil {
		return err
	}

	resp, err := s.DataInterface.PostCredentials(ctx, data)
	if err != nil {
		return err
	}

	return s.StorageData.CreateCredentials(ctx, data.Data, resp.UserDataId, data.Name, data.Description, resp.Hash, resp.CreatedAt, resp.UpdateAt)
}

// todo :text
func (s *Service) CreateFileData(ctx context.Context, data *server.ReqData) error {

	if err := s.setJwtToken(ctx); err != nil {
		return err
	}
	//todo шифруем data
	resp, err := s.DataInterface.PostCrateFile(ctx, data)
	if err != nil {
		return err
	}
	return s.StorageData.CreateFileData(ctx, data.Data, resp.UserDataId, data.Name, data.Description, resp.Hash, resp.CreatedAt, resp.UpdateAt)
}

// CreateCreditCard - создание данных о кредитной карте
func (s *Service) CreateCreditCard(ctx context.Context, data *server.ReqData) error {

	if err := s.setJwtToken(ctx); err != nil {
		return err
	}

	if err := s.encryptData(data); err != nil {
		return err
	}

	resp, err := s.DataInterface.PostCreditCard(ctx, data)
	if err != nil {
		return err
	}
	return s.StorageData.CreateCreditCard(ctx, data.Data, resp.UserDataId, data.Name, data.Description, resp.Hash, resp.CreatedAt, resp.UpdateAt)
}

// PingServer - пинг сервера
func (s *Service) PingServer(ctx context.Context) bool {
	if err := s.DataInterface.Ping(ctx); err != nil {
		return false
	}
	return true
}

// GetData - получение данных любого формата
func (s *Service) GetData(ctx context.Context, userDataId int64) ([]byte, error) {
	// Проверяем доступен ли сервер
	if !s.PingServer(ctx) {
		fmt.Println("Сервер недоступен")
		resp, err := s.GetDataFromAgentStorage(ctx, userDataId)
		if err != nil {
			return nil, err
		}
		return resp, err
	}

	// Выставляем токен для будущих запросов
	if err := s.setJwtToken(ctx); err != nil {
		return nil, err
	}

	// Проверяем Новые ли данные на сервере
	ok, err := s.CheckNewData(ctx, userDataId)

	if !ok {
		// Если не новые скачиваем из локального хранилища
		resp, err := s.GetDataFromAgentStorage(ctx, userDataId)
		if err != nil {
			return nil, err
		}
		// сохраняем файл
		return resp, err
	}

	// Получение файла из сервера
	resp, err := s.DataInterface.GetData(ctx, userDataId)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(resp))

	decrypt, err := s.Encrypter.Decrypt(resp)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(decrypt))
	return decrypt, nil
}

// CheckNewData - проверка на новые данные
func (s *Service) CheckNewData(ctx context.Context, userDataId int64) (bool, error) {
	// Получаем инофрмацию о обновление текущих данных
	data, err := s.StorageData.GetInfoData(ctx, userDataId)
	if err != nil {
		return false, err
	}

	// Выставляем токен для будущих запросов
	if err := s.setJwtToken(ctx); err != nil {
		return false, err
	}

	ok, err := s.CheckUpdate(ctx, userDataId, data.UpdateAt)
	if err != nil {
		return false, err
	}
	return ok, nil
}

// GetDataFromAgentStorage - получение данных из хранилища агента
func (s *Service) GetDataFromAgentStorage(ctx context.Context, userDataId int64) ([]byte, error) {

	fmt.Println("Скачиваем данные из локального хранилища")

	// Получение файла из хранилища
	usersData, dataFile, err := s.StorageData.GetData(ctx, userDataId)
	if err != nil {
		return nil, err
	}

	if usersData.DataType == sqllite.TypeFile {
		metaData, err := s.GetMetaData(ctx, userDataId)
		if err != nil {
			return nil, err
		}
		//todo
		// расшифровываем бинарные файлы

		// Сохраняем файл

		resp := fmt.Sprintf("Файл сохранен %s/%s; Название оригинальное %s", metaData.PathSave, metaData.FileName, string(dataFile.EncryptData))

		return []byte(resp), nil
	}

	decrypt, err := s.Encrypter.Decrypt(dataFile.EncryptData)
	if err != nil {
		return nil, err
	}

	resp := fmt.Sprintf("Данные %s", string(decrypt))

	return []byte(resp), err
}

// GetListData - получение списка актуальных данных пользователя
func (s *Service) GetListData(ctx context.Context) ([]byte, error) {

	if err := s.setJwtToken(ctx); err != nil {
		return nil, err
	}
	// todo Получение из локального хралищиа

	resp, err := s.DataInterface.GetListData(ctx)
	if err != nil {
		return nil, err
	}
	return resp, nil

}

// UpdateData - обновление данных пользователя (кроме бинарного файла)
func (s *Service) UpdateData(ctx context.Context, userDataId int64, data []byte) ([]byte, error) {
	if err := s.setJwtToken(ctx); err != nil {
		return nil, err
	}
	//todo шифруем data
	data, err := s.Encrypter.Encrypt(data)
	if err != nil {
		return nil, err
	}

	resp, err := s.DataInterface.PostUpdateData(ctx, userDataId, data)
	if err != nil {
		return nil, err
	}

	err = s.StorageData.UpdateData(ctx, userDataId, data, resp.Hash, resp.UpdateAt)
	if err != nil {
		return nil, err
	}

	return []byte("Data updated"), nil
}
func (s *Service) encryptData(reqData *server.ReqData) error {

	// Шифруем данные о файле
	DataEncrypt, err := s.Encrypter.Encrypt((reqData.Data))
	if err != nil {
		return err
	}

	reqData.Data = DataEncrypt

	return nil
}
func (s *Service) decryptData(reqData *server.ReqData) error {
	DataDecrypt, err := s.Encrypter.Decrypt(reqData.Data)
	if err != nil {
		return err
	}

	reqData.Data = DataDecrypt
	return nil
}
