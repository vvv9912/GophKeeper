package service

import (
	"context"
	"crypto/rsa"
)

type Auth interface {
	BuildJWTString(userId int64) (string, error)
	GetUserId(tokenString string) (int64, error)
}

type StoreAuth interface {
	CreateUser(ctx context.Context, login, password string) (int64, error)
	GetUserId(ctx context.Context, login string, password string) (int64, error)
}

type Service struct {
	Auth
	StoreAuth
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}
