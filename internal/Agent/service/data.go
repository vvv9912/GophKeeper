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
func (s *Service) CreateFile(ctx context.Context, path string, name, description string) error {
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
		g, _ := generateKey()
		encrypt, err := encryptData(data, g)
		if err != nil {
			return err
		}
		_ = encrypt
		fmt.Println("Old uuid-chunk:", uuidChunk, "Num chunk:", n)
		ennd := r.SizeChunk * (i)
		if i == n {
			ennd = int(r.Size())
		}
		uuidChunk, err = s.PostCrateFileStartChunks(ctx, data, r.NameFile, uuidChunk, int(r.SizeChunk)*(i-1), ennd, int(r.Size()), reqDataJson)
		if err != nil {
			return err
		}
		fmt.Println("New uuid-chunk:", uuidChunk, "Num chunk:", n)
	}
	return nil
}
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
		s.AuthService.SetJWTToken(jwt)
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
