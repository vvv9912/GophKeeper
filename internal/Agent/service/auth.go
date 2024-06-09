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

// setJwtToken - выставляем токен для будущих запросов
func (s *Service) setJwtToken(ctx context.Context) error {
	// Проверка на пустой токен
	if s.AuthService.GetJWTToken() == "" {

		// Получаем токен из локального хранилища
		jwt, err := s.StorageData.GetJWTToken(ctx)
		if err != nil {
			return err
		}

		if jwt == "" {
			fmt.Println("jwt is empty")
			return fmt.Errorf("jwt is empty")
		}

		// Выставляем в структуре AuthServer JWT token
		s.AuthService.SetJWTToken(jwt)
		fmt.Println("jwt", jwt)
	}

	return nil
}
