package service

import (
	"GophKeeper/internal/Agent/server"
	"context"
)

func (s *Service) CreateCredentials(ctx context.Context, data *server.ReqData) error {
	//todo encrypt
	resp, err := s.DataInterface.PostCredentials(ctx, data)
	if err != nil {
		return err
	}
	return s.StorageData.CreateCredentials(ctx, data.Data, resp.UserDataId, data.Name, data.Description, resp.Hash)
}
