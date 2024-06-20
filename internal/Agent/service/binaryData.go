package service

import (
	"GophKeeper/internal/Agent/server"
	"GophKeeper/pkg/logger"
	"GophKeeper/pkg/store"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"io"
	"os"
	path2 "path"
)

// CreateBinaryFile - создание файла бинарного
func (s *Service) CreateBinaryFile(ctx context.Context, path string, name, description string, ch chan<- string) error {
	// Получение jwt токена
	if err := s.setJwtToken(ctx); err != nil {
		return err
	}

	// Создаем новый шифрованный файл в tmp папке
	pathTmpFile, err := s.createEncryptedFile(path)
	if err != nil {
		return err
	}

	// Очистка из tmp папки
	defer func() {
		os.Remove(pathTmpFile)
	}()

	// Получаем оригинальное название файла пользователя
	_, originalFileName := path2.Split(path)

	// Распределение на чанки
	r := NewReader(pathTmpFile)

	// Подготовка данных для запроса на сервер
	reqData, reqDataJson, err := s.prepareReqBinaryFile(originalFileName, name, description)
	if err != nil {
		return err
	}

	// Передача на сервер
	resp, err := s.transferCreateDataBinaryFile(ctx, r, reqDataJson, ch)
	if err != nil {
		return err
	}
	// сохранение в локальный репозиторий
	if err = s.saveLocalFile(ctx, r, pathTmpFile, name, description, reqData, resp); err != nil {
		return err
	}

	return nil
}

// UpdateBinaryFile - обновление данных бинарного формата
func (s *Service) UpdateBinaryFile(ctx context.Context, path string, userDataId int64, ch chan<- string) error {
	if err := s.setJwtToken(ctx); err != nil {
		return err
	}

	// Создаем новый шифрованный файл в tmp папке
	pathTmpFile, err := s.createEncryptedFile(path)
	if err != nil {
		return err
	}
	defer func() {
		os.Remove(pathTmpFile)
	}()

	r := NewReader(pathTmpFile)
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
	//todo удаление предыдущего файла

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

func (s *Service) createEncryptedFile(path string) (string, error) {
	// Создаем новый шифрованный файл в tmp папке
	// Имя нового файла
	newNameFile := uuid.NewString()
	//Полный путь к новому файлу
	pathTmp := path2.Join(PathTmp, newNameFile)

	err := s.Encrypter.EncryptFile(path, pathTmp)
	if err != nil {
		return "", err
	}

	return pathTmp, nil
}
func (s *Service) prepareReqBinaryFile(originalFileName string, name, description string) (*server.ReqData, []byte, error) {

	// Данные о файле
	infoOriginalFile := server.DataFileInfo{OriginalFileName: originalFileName}

	// Метаданные оригинального файла
	dataOriginalFile, err := json.Marshal(infoOriginalFile)
	if err != nil {
		logger.Log.Error("Marshal json failed", zap.Error(err))
		return nil, nil, err
	}

	// Структура запроса с данными о файле на сервере
	reqData := &server.ReqData{
		Name:        name,
		Description: description,
		Data:        dataOriginalFile,
	}

	// Шифруем данные о файле
	if err := s.encryptData(reqData); err != nil {
		return nil, nil, err
	}

	reqDataJson, err := json.Marshal(reqData)
	if err != nil {
		logger.Log.Error("Marshal json failed", zap.Error(err))
		return nil, nil, err
	}
	return reqData, reqDataJson, nil

}

func (s *Service) transferCreateDataBinaryFile(ctx context.Context, r *Reader, reqDataJson []byte, ch chan<- string) (*server.RespData, error) {
	var resp *server.RespData

	// Количество чанков в файле
	n, err := r.NumChunk()
	if err != nil {
		return nil, err
	}

	var uuidChunk string
	// Передача файла на сервер
	for i := 1; i <= n; i++ {
		// Считываем файл по чанкам
		data, err := r.ReadFile(i)
		if err != nil {
			return nil, err
		}

		maxChunk := r.SizeChunk * (i)
		if i == n {
			maxChunk = int(r.Size())
		}

		uuidChunk, resp, err = s.PostCrateFileStartChunks(ctx, data, r.NameFile, uuidChunk, int(r.SizeChunk)*(i-1), maxChunk, int(r.Size()), reqDataJson)
		if err != nil {
			logger.Log.Error("PostCrateFileStartChunks failed", zap.Error(err))
			return nil, err
		}

		// Вывод полезной информации пользователю
		ch <- fmt.Sprintf("Передано кБайт %0.1f из %0.1f", float64(maxChunk)/1024.0, float64(r.Size())/1024)

	}
	if resp == nil {
		err := fmt.Errorf("resp is nil")
		logger.Log.Error("resp is nil", zap.Error(err))
		return nil, err
	}

	return resp, nil
}

func (s *Service) saveLocalFile(ctx context.Context, r *Reader, pathTmpFile, name, description string, reqData *server.ReqData, resp *server.RespData) error {
	// Копируем файл в локальное хранилище Агента
	NewNameFile := uuid.NewString()
	if err := copyFile(pathTmpFile, PathStorage, NewNameFile); err != nil {
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
	if err := s.StorageData.CreateBinaryFile(ctx, reqData.Data, resp.UserDataId, name, description, resp.Hash, resp.CreatedAt, resp.UpdateAt, metaData); err != nil {
		logger.Log.Error("CreateBinaryFile failed", zap.Error(err))
		return err
	}
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

func (s *Service) decryptFile(ctx context.Context, meta *store.MetaData, originalFileName string) (string, error) {

	saveOrigFile := path2.Join(PathUserData, originalFileName)

	err := s.DecryptFile(meta.PathSave, saveOrigFile)
	if err != nil {
		return "", err
	}

	return saveOrigFile, nil
}
