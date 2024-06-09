package service

import (
	"context"
	"fmt"
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
