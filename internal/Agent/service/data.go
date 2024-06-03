package service

import "GophKeeper/internal/Agent/server"

func (s *Service) CreateCredentials(ctx context.Context, data *server.ReqData) error {
	//todo encrypt
	err := s.DataInterface.PostCredentials(ctx, data)
	if err != nil {
		return err
	}
	return nil
}
