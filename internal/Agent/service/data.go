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

var PathStorage = "Agent/storage"

func (s *Service) CreateCredentials(ctx context.Context, data *server.ReqData) error {

	if err := s.setJwtToken(ctx); err != nil {
		return err
	}

	resp, err := s.DataInterface.PostCredentials(ctx, data)
	if err != nil {
		return err
	}

	return s.StorageData.CreateCredentials(ctx, data.Data, resp.UserDataId, data.Name, data.Description, resp.Hash, resp.CreatedAt, resp.UpdateAt)
}

// binary file
func (s *Service) CreateFile(ctx context.Context, path string, name, description string, ch chan<- string) error {
	if err := s.setJwtToken(ctx); err != nil {
		return err
	}

	r := NewReader(path)
	n, err := r.NumChunk()
	if err != nil {
		return err
	}
	var uuidChunk string

	dataInfo := server.DataFileInfo{OriginalFileName: r.NameFile}
	data, err := json.Marshal(dataInfo)
	if err != nil {
		logger.Log.Error("Marshal json failed", zap.Error(err))
		return err
	}

	reqData := &server.ReqData{
		Name:        name,
		Description: description,
		Data:        data, //todo шифруем data
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

		//g, _ := generateKey()
		//
		//encrypt, err := encryptData(data, g)
		//if err != nil {
		//	return err
		//}
		//_ = encrypt

		maxChunk := r.SizeChunk * (i)
		if i == n {
			maxChunk = int(r.Size())
		}
		uuidChunk, resp, err = s.PostCrateFileStartChunks(ctx, data, r.NameFile, uuidChunk, int(r.SizeChunk)*(i-1), maxChunk, int(r.Size()), reqDataJson)
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

	resp, err := s.DataInterface.PostCrateFile(ctx, data)
	if err != nil {
		return err
	}
	return s.StorageData.CreateFileData(ctx, data.Data, resp.UserDataId, data.Name, data.Description, resp.Hash, resp.CreatedAt, resp.UpdateAt)
}

func (s *Service) CreateCreditCard(ctx context.Context, data *server.ReqData) error {

	if err := s.setJwtToken(ctx); err != nil {
		return err
	}

	resp, err := s.DataInterface.PostCreditCard(ctx, data)
	if err != nil {
		return err
	}
	return s.StorageData.CreateCreditCard(ctx, data.Data, resp.UserDataId, data.Name, data.Description, resp.Hash, resp.CreatedAt, resp.UpdateAt)
}

func (s *Service) PingServer(ctx context.Context) bool {
	if err := s.DataInterface.Ping(ctx); err != nil {
		return false
	}
	return true
}

func (s *Service) GetData(ctx context.Context, userDataId int64) ([]byte, error) {
	if !s.PingServer(ctx) {
		fmt.Println("Сервер недоступен")
		usersData, dataFile, err := s.GetDataFromAgentStorage(ctx, userDataId)
		if err != nil {
			return nil, err
		}

		if usersData.DataType == sqllite.TypeFile {
			metaData, err := s.GetMetaData(ctx, userDataId)
			if err != nil {
				return nil, err
			}
			resp := fmt.Sprintf("Файл сохранен %s/%s; Название оригинальное %s", metaData.PathSave, metaData.FileName, string(dataFile.EncryptData))
			return []byte(resp), nil
		}
		resp := fmt.Sprintf("Данные %s", string(dataFile.EncryptData))
		return []byte(resp), err
	}

	if err := s.setJwtToken(ctx); err != nil {
		return nil, err
	}

	// Получение файла из сервера
	resp, err := s.DataInterface.GetData(ctx, userDataId)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(resp))
	return resp, nil
}

func (s *Service) CheckNewData(ctx context.Context, userDataId int64) (bool, error) {
	data, err := s.StorageData.GetInfoData(ctx, userDataId)
	if err != nil {
		return false, err
	}

	if err := s.setJwtToken(ctx); err != nil {
		return false, err
	}

	ok, err := s.CheckUpdate(ctx, userDataId, data.UpdateAt)
	if err != nil {
		return false, err
	}
	return ok, nil
}

func (s *Service) GetDataFromAgentStorage(ctx context.Context, userDataId int64) (*store.UsersData, *store.DataFile, error) {

	// Получение файла из хранилища
	usersData, data, err := s.StorageData.GetData(ctx, userDataId)
	if err != nil {
		return nil, nil, err
	}

	return usersData, data, nil
}
func (s *Service) GetListData(ctx context.Context) ([]byte, error) {
	if err := s.setJwtToken(ctx); err != nil {
		return nil, err
	}

	resp, err := s.DataInterface.GetListData(ctx)
	if err != nil {
		return nil, err
	}
	return resp, nil

}
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

// todo обновление только файла
//func (s *Service) UpdateBinaryFile(ctx context.Context, userDataId int64, name, description string) error {
//	//if err := s.setJwtToken(ctx); err != nil {
//	//	return err
//	//}
//
//	return s.StorageData.UpdateDataFile(ctx, userDataId, name, description)
//}

//func (s *Service) UpdateDataCreditCard(ctx context.Context, userDataId int64, name, description string) error {
//	if err := s.setJwtToken(ctx); err != nil {
//		return err
//	}
//	return s.StorageData.UpdateDataCreditCard(ctx, userDataId, name, description)
//}

//func (s *Service) UpdateDataCredentials(ctx context.Context, userDataId int64, name, description string) error {
//	if err := s.setJwtToken(ctx); err != nil {
//		return err
//	}
//	return s.StorageData.UpdateDataCredentials(ctx, userDataId, name, description)
//}

//func (s *Service) UpdateDataFile(ctx context.Context, userDataId int64, name, description string) error {
//	if err := s.setJwtToken(ctx); err != nil {
//		return err
//	}
//	return s.StorageData.UpdateDataFile(ctx, userDataId, name, description)
//}
