package service

import (
	"context"
)

func (s *Service) SignIn(ctx context.Context, username, password string) (string, error) {
	user, err := s.AuthService.SignIn(ctx, username, password)
	if err != nil {
		return "", err
	}

	// Выставляем в структуре AuthServer JWT token
	s.AuthService.SetJWTToken(user.JWT)

	// Добавляем в БД sqlite
	err = s.StorageData.SetJWTToken(ctx, user.JWT)
	if err != nil {
		return "", err
	}

	return user.JWT, err
}

func (s *Service) SignUp(ctx context.Context, username, password string) (string, error) {
	user, err := s.AuthService.SignUp(ctx, username, password)
	if err != nil {
		return "", err
	}

	// Выставляем в структуре AuthServer JWT token
	s.AuthService.SetJWTToken(user.JWT)

	// Добавляем в БД sqlite
	err = s.StorageData.SetJWTToken(ctx, user.JWT)
	if err != nil {
		return "", err
	}

	return user.JWT, err
}
