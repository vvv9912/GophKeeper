package service

import "context"

func (s *Service) SignIn(ctx context.Context, username, password string) (string, error) {
	user, err := s.AuthService.SignIn(ctx, username, password)
	if err != nil {
		return "", err
	}
	//todo в бд
	s.AuthService.SetJWTToken(user.JWT)

	return user.JWT, err
}

func (s *Service) SignUp(ctx context.Context, username, password string) (string, error) {
	user, err := s.AuthService.SignUp(ctx, username, password)
	if err != nil {
		return "", err
	}

	return user.JWT, err
}
