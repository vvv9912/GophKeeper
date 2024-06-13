package service

import (
	"GophKeeper/internal/Agent/server"
	"GophKeeper/pkg/logger"
	"GophKeeper/pkg/store"
	"GophKeeper/pkg/store/sqllite"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"io"
	"os"
	path2 "path"
)

var PathStorage = "FileAgent/storage"
var PathTmp = "FileAgent/tmp"

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

// CreateFile - создание файла бинарного
func (s *Service) CreateFile(ctx context.Context, path string, name, description string, ch chan<- string) error {
	// Получение jwt токена
	if err := s.setJwtToken(ctx); err != nil {
		return err
	}
	//todo шифруем data
	// Считывание файла по чанкам

	// Создаеим новый шифрованный файл в tmp папке
	newNameFile := uuid.NewString()
	err := s.Encrypter.EncryptFile(path, path2.Join(PathTmp, newNameFile))
	if err != nil {
		return err
	}

	defer func() {
		os.Remove(path2.Join(PathTmp, newNameFile))
	}()

	r := NewReader(path2.Join(PathTmp, newNameFile))
	n, err := r.NumChunk()
	if err != nil {
		return err
	}

	// Оригинальное название файла
	_, originalFileName := path2.Split(path)

	// id - чанка при передачи
	var uuidChunk string

	// Данные о файле
	dataInfo := server.DataFileInfo{OriginalFileName: originalFileName}

	data, err := json.Marshal(dataInfo)
	if err != nil {
		logger.Log.Error("Marshal json failed", zap.Error(err))
		return err
	}

	// Структура запроса с данными о файле на сервер
	reqData := &server.ReqData{
		Name:        name,
		Description: description,
		Data:        data,
	}

	// Шифруем данные о файле
	if err := s.encryptData(reqData); err != nil {
		return err
	}

	reqDataJson, err := json.Marshal(reqData)
	if err != nil {
		logger.Log.Error("Marshal json failed", zap.Error(err))
		return err
	}

	var resp *server.RespData
	fmt.Println(r.NameFile)

	// Передача файла на сервер
	for i := 1; i <= n; i++ {
		// Считываем файл по чанкам
		data, err := r.ReadFile(i)
		if err != nil {
			return err
		}

		maxChunk := r.SizeChunk * (i)
		if i == n {
			maxChunk = int(r.Size())
		}

		uuidChunk, resp, err = s.PostCrateFileStartChunks(ctx, data, r.NameFile, uuidChunk, int(r.SizeChunk)*(i-1), maxChunk, int(r.Size()), reqDataJson)
		if err != nil {
			logger.Log.Error("PostCrateFileStartChunks failed", zap.Error(err))
			return err
		}

		// Вывод полезной информации пользователю
		ch <- fmt.Sprintf("Передано кБайт %0.1f из %0.1f", float64(maxChunk)/1024.0, float64(r.Size())/1024)

	}
	if resp == nil {
		err := fmt.Errorf("resp is nil")
		logger.Log.Error("resp is nil", zap.Error(err))
		return err
	}

	// Копируем файл в локальное хранилище Агента
	NewNameFile := uuid.NewString()
	if err := copyFile(path, PathStorage, NewNameFile); err != nil {
		logger.Log.Error("copyFile failed", zap.Error(err))
		return err
	}
	// Новое поле с мета данными, тк путь и названия файла поменялись
	metaData := &store.MetaData{
		FileName: NewNameFile,
		Size:     r.size,
		PathSave: PathStorage,
	}

	// Сохранение в локальное хранилище
	if err := s.StorageData.CreateBinaryFile(ctx, data, resp.UserDataId, name, description, resp.Hash, resp.CreatedAt, resp.UpdateAt, metaData); err != nil {
		logger.Log.Error("CreateBinaryFile failed", zap.Error(err))
		return err
	}

	return nil
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

// UpdateBinaryFile - обновление данных бинарного формата
func (s *Service) UpdateBinaryFile(ctx context.Context, path string, userDataId int64, ch chan<- string) error {
	ctx = context.Background()
	if err := s.setJwtToken(ctx); err != nil {
		return err
	}
	// Создаеим новый шифрованный файл
	newNameFile := uuid.NewString()
	err := s.Encrypter.EncryptFile(path, path2.Join(PathTmp, newNameFile))
	if err != nil {
		return err
	}

	defer func() {
		os.Remove(path2.Join(PathTmp, newNameFile))
	}()

	r := NewReader(path2.Join(PathTmp, newNameFile))
	n, err := r.NumChunk()
	if err != nil {
		return err
	}

	// Оригинальное название файла
	_, originalFileName := path2.Split(path)

	var uuidChunk string

	dataInfo := server.DataFileInfo{OriginalFileName: originalFileName}

	data, err := json.Marshal(dataInfo)
	if err != nil {
		logger.Log.Error("Marshal json failed", zap.Error(err))
		return err
	}

	reqData := &server.ReqData{
		Data: data,
	}

	// Шифруем данные о файле
	if err := s.encryptData(reqData); err != nil {
		return err
	}

	reqDataJson, err := json.Marshal(reqData)
	if err != nil {
		logger.Log.Error("Marshal json failed", zap.Error(err))
		return err
	}
	var resp *server.RespData

	for i := 1; i <= n; i++ {

		data, err := r.ReadFile(i)
		if err != nil {
			return err
		}

		maxChunk := r.SizeChunk * (i)
		if i == n {
			maxChunk = int(r.Size())
		}
		uuidChunk, resp, err = s.PostUpdateBinaryFile(ctx, data, r.NameFile, uuidChunk, int(r.SizeChunk)*(i-1), maxChunk, int(r.Size()), reqDataJson, userDataId)
		if err != nil {
			logger.Log.Error("PostCrateFileStartChunks failed", zap.Error(err))
			return err
		}

		ch <- fmt.Sprintf("Передано кБайт %0.1f из %0.1f", float64(maxChunk)/1024.0, float64(r.Size())/1024)

	}

	// Копирование к себе в хранилище
	NewNameFile := uuid.NewString()
	if err := copyFile(path, PathStorage, NewNameFile); err != nil {
		logger.Log.Error("copyFile failed", zap.Error(err))
		return err
	}

	metaData := &store.MetaData{
		FileName: NewNameFile,
		Size:     r.size,
		PathSave: PathStorage,
	}

	if resp == nil {
		err := fmt.Errorf("resp is nil")
		logger.Log.Error("resp is nil", zap.Error(err))
		return err
	}

	meta, err := json.Marshal(metaData)
	if err != nil {
		logger.Log.Error("Marshal json failed", zap.Error(err))
		return err
	}
	if err := s.StorageData.UpdateDataBinary(ctx, userDataId, data, resp.Hash, resp.UpdateAt, meta); err != nil {
		logger.Log.Error("CreateBinaryFile failed", zap.Error(err))
		return err
	}

	return nil
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

// copyFile - копирование файла по новому пути и новым именим
func copyFile(src, newPath string, newNameFile string) error {

	if err := os.MkdirAll(newPath, os.ModePerm); err != nil {
		return err
	}

	dst := path2.Join(newPath, newNameFile)

	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}

	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}

	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	return nil
}
