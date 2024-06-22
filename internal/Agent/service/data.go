package service

import (
	"GophKeeper/internal/Agent/Encrypt"
	"GophKeeper/internal/Agent/model"
	"GophKeeper/internal/Agent/server"
	"GophKeeper/pkg/logger"
	"GophKeeper/pkg/store/sqllite"
	"context"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type UseCase struct {
	AuthService
	DataInterface
	StorageData
	Encrypter
	JWTToken string
}

func NewUseCase(db *sqlx.DB, key []byte, certFile, keyFile string, serverDns string) *UseCase {

	serv := server.NewAgentServer(certFile, keyFile, serverDns)

	encrypt, err := Encrypt.NewEncrypt(key)
	if err != nil {
		panic(err)
	}

	return &UseCase{
		AuthService:   serv,
		DataInterface: serv,
		StorageData:   sqllite.NewDatabase(db),
		Encrypter:     encrypt,
	}
}

// CreateCredentials - Создание данных логин/пароль
func (s *UseCase) CreateCredentials(ctx context.Context, data *server.ReqData) error {
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

// CreateFileData - создание файла
func (s *UseCase) CreateFileData(ctx context.Context, data *server.ReqData) error {

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
func (s *UseCase) CreateCreditCard(ctx context.Context, data *server.ReqData) error {

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
func (s *UseCase) PingServer(ctx context.Context) bool {
	if err := s.DataInterface.Ping(ctx); err != nil {
		return false
	}
	return true
}

// GetData - получение данных любого формата
func (s *UseCase) GetData(ctx context.Context, userDataId int64) ([]byte, error) {
	// Проверяем доступен ли сервер
	if !s.PingServer(ctx) {
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
	if err != nil {
		return nil, err
	}

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

	logger.Log.Debug("data from server", zap.String("resp", string(resp)))

	decrypt, err := s.Encrypter.Decrypt(resp)
	if err != nil {
		return nil, err
	}

	logger.Log.Debug("decrypt from server", zap.String("decrypt", string(decrypt)))
	return decrypt, nil
}

// CheckNewData - проверка на новые данные
func (s *UseCase) CheckNewData(ctx context.Context, userDataId int64) (bool, error) {
	// Получаем инофрмацию о обновление текущих данных
	data, err := s.StorageData.GetInfoData(ctx, userDataId)
	if err != nil {
		return false, err
	}

	// Выставляем токен для будущих запросов
	if err = s.setJwtToken(ctx); err != nil {
		return false, err
	}

	ok, err := s.CheckUpdate(ctx, userDataId, data.UpdateAt)
	if err != nil {
		return false, err
	}
	return ok, nil
}

// GetDataFromAgentStorage - получение данных из хранилища агента
func (s *UseCase) GetDataFromAgentStorage(ctx context.Context, userDataId int64) ([]byte, error) {
	logger.Log.Info("Получаем данные из локального хранилища")

	// Получение файла из хранилища
	usersData, dataFile, err := s.StorageData.GetData(ctx, userDataId)
	if err != nil {
		return nil, err
	}

	decrypt, err := s.Encrypter.Decrypt(dataFile.EncryptData)
	if err != nil {
		return nil, err
	}

	if usersData.DataType == sqllite.TypeBinaryFile {
		metaData, err := s.StorageData.GetMetaData(ctx, userDataId)
		if err != nil {
			return nil, err
		}

		var fileData model.Data
		err = json.Unmarshal(decrypt, &fileData)
		if err != nil {
			logger.Log.Error("Unmarshal", zap.Error(err))
			return nil, err
		}

		origFileName := string(fileData.DecryptData)
		origPathSave, err := s.decryptFile(ctx, metaData, origFileName)
		if err != nil {
			return nil, err
		}

		resp := fmt.Sprintf("Файл сохранен %s/%s; Название оригинальное %s", origPathSave, metaData.FileName, string(dataFile.EncryptData))

		return []byte(resp), nil
	}

	resp := fmt.Sprintf("Данные %s", string(decrypt))

	return []byte(resp), err
}

// GetListData - получение списка актуальных данных пользователя
func (s *UseCase) GetListData(ctx context.Context) ([]byte, error) {

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
func (s *UseCase) UpdateData(ctx context.Context, userDataId int64, data []byte) ([]byte, error) {
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
func (s *UseCase) encryptData(reqData *server.ReqData) error {

	// Шифруем данные о файле
	DataEncrypt, err := s.Encrypter.Encrypt((reqData.Data))
	if err != nil {
		return err
	}

	reqData.Data = DataEncrypt

	return nil
}
func (s *UseCase) decryptData(reqData *server.ReqData) error {
	DataDecrypt, err := s.Encrypter.Decrypt(reqData.Data)
	if err != nil {
		return err
	}

	reqData.Data = DataDecrypt
	return nil
}
