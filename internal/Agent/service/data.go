package service

import (
	"GophKeeper/internal/Agent/server"
	"GophKeeper/pkg/logger"
	"context"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
)

func (s *Service) CreateCredentials(ctx context.Context, data *server.ReqData) error {

	if err := s.setJwtToken(ctx); err != nil {
		return err
	}

	resp, err := s.DataInterface.PostCredentials(ctx, data)
	if err != nil {
		return err
	}

	return s.StorageData.CreateCredentials(ctx, data.Data, resp.UserDataId, data.Name, data.Description, resp.Hash)
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
		uuidChunk, err = s.PostCrateFileStartChunks(ctx, data, r.NameFile, uuidChunk, int(r.SizeChunk)*(i-1), maxChunk, int(r.Size()), reqDataJson)
		if err != nil {
			logger.Log.Error("PostCrateFileStartChunks failed", zap.Error(err))
			return err
		}

		ch <- fmt.Sprintf("Передано кБайт %0.1f из %0.1f", float64(maxChunk)/1024.0, float64(r.Size())/1024)

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
	return s.StorageData.CreateFileData(ctx, data.Data, resp.UserDataId, data.Name, data.Description, resp.Hash)
}

func (s *Service) CreateCreditCard(ctx context.Context, data *server.ReqData) error {

	if err := s.setJwtToken(ctx); err != nil {
		return err
	}

	resp, err := s.DataInterface.PostCreditCard(ctx, data)
	if err != nil {
		return err
	}
	return s.StorageData.CreateCreditCard(ctx, data.Data, resp.UserDataId, data.Name, data.Description, resp.Hash)
}
func (s *Service) setJwtToken(ctx context.Context) error {

	if s.AuthService.GetJWTToken() == "" {
		jwt, err := s.StorageData.GetJWTToken(ctx)
		if err != nil {
			return err
		}
		if jwt == "" {
			fmt.Println("jwt is empty")
			return fmt.Errorf("jwt is empty")
		}
		s.AuthService.SetJWTToken(jwt)
		fmt.Println("jwt", jwt)
	}

	return nil
}

func (s *Service) GetData(ctx context.Context, userDataId int64) ([]byte, error) {
	if err := s.setJwtToken(ctx); err != nil {
		return nil, err
	}

	resp, err := s.DataInterface.GetData(ctx, userDataId)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(resp))
	return resp, nil
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

//// todo обновление только файла
//func (s *Service) UpdateDataFile(ctx context.Context, userDataId int64, name, description string) error {
//	if err := s.setJwtToken(ctx); err != nil {
//		return err
//	}
//	return s.StorageData.UpdateDataFile(ctx, userDataId, name, description)
//}
