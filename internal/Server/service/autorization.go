package service

import (
	"GophKeeper/internal/Server/authorization"
	"GophKeeper/pkg/customErrors"
	"context"
	"net/http"
)

// SignUp - регистрация пользователя.
func (s *UseCase) SignUp(ctx context.Context, login, password string) (string, error) {
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

// SignIn - авторизация пользователя.
func (s *UseCase) SignIn(ctx context.Context, login, password string) (string, error) {
	// Получаем хэш пароль передает User

	hashPassword := authorization.Sha256Hash(password)

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
