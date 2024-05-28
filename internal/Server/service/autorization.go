package service

import (
	"GophKeeper/internal/Server/authorization"
	"GophKeeper/pkg/customErrors"
	"context"
	"net/http"
)

func (s *Service) SignUp(ctx context.Context, login, password string) (string, error) {
	// Получаем хэш пароль передает User

	hashPassword := authorization.Sha256Hash(password)

	userId, err := s.StoreAuth.CreateUser(ctx, login, hashPassword)
	if err != nil {
		return "", customErrors.NewCustomError(err, http.StatusInternalServerError, "Error create user")
	}

	jwt, err := s.Auth.BuildJWTString(userId)
	if err != nil {
		return "", customErrors.NewCustomError(err, http.StatusInternalServerError, "Error create JWT")
	}

	return jwt, nil
}

func (s *Service) SignIn(ctx context.Context, login, password string) (string, error) {
	// Получаем хэш пароль передает User

	hashPassword := authorization.Sha256Hash(password)

	// TODO:Добавить ошибки кастомные
	userId, err := s.StoreAuth.GetUserId(ctx, login, hashPassword)
	if err != nil {
		return "", customErrors.NewCustomError(err, http.StatusInternalServerError, "Error get user id")
	}

	jwt, err := s.Auth.BuildJWTString(userId)
	if err != nil {
		return "", customErrors.NewCustomError(err, http.StatusInternalServerError, "Error create JWT")
	}

	return jwt, err
}

//func (s *Service) GetUserIdFromJwt(jwt string) (int64, error) {
//	return s.Auth.GetUserId(jwt)
//}
