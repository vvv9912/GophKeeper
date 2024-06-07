package service

import (
	"GophKeeper/internal/Agent/server"
	"context"
	"fmt"
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
func (s *Service) CreateFile(ctx context.Context, path string) error {
	r := NewReader(path)
	n, err := r.NumChunk()
	if err != nil {
		return err
	}
	var uuidChunk string
	for i := 1; i <= n; i++ {
		data, err := r.ReadFile(i)
		if err != nil {
			return err
		}

		//encrypt, err := EncryptData(data, s.GetPublicKey())
		//if err != nil {
		//	return err
		//}
		fmt.Println("Old uuid-chunk:", uuidChunk, "Num chunk:", n)
		ennd := r.SizeChunk * (i)
		if i == n {
			ennd = int(r.Size())
		}
		uuidChunk, err = s.PostCrateFileStartChunks(ctx, data, r.Name(), uuidChunk, int(r.SizeChunk)*(i-1), ennd, int(r.Size()))
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
