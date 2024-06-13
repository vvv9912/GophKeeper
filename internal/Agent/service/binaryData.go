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

// UpdateBinaryFile - обновление данных бинарного формата
func (s *Service) UpdateBinaryFile(ctx context.Context, path string, userDataId int64, ch chan<- string) error {
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
