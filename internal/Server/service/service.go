package service

import (
	"GophKeeper/internal/Server/authorization"
	"GophKeeper/pkg/store"
	"GophKeeper/pkg/store/postgresql"
	"context"
	"crypto/rsa"
	"github.com/jmoiron/sqlx"
	"time"
)

// Auth - интерфейс аутентификации (расчет токена и получения данных из него).
type Auth interface {
	BuildJWTString(userId int64) (string, error)
	GetUserId(tokenString string) (int64, error)
}

// StoreAuth - интерфейс для работы с БД пользователя.
type StoreAuth interface {
	CreateUser(ctx context.Context, login, password string) (int64, error)
	GetUserId(ctx context.Context, login string, password string) (int64, error)
}

// Data - интерфейс для работы с данными пользователя.
type Data interface {
	CreateCredentials(ctx context.Context, userId int64, data []byte, name, description string) error
	CreateCreditCard(ctx context.Context, userId int64, data []byte, name, description string) error
	CreateFile(ctx context.Context, userId int64, data []byte, name, description string) error
	ChangeData(ctx context.Context, userId int64, lastTimeUpdate time.Time) ([]byte, error)
}

// StoreData - интерфейс для работы с БД данных пользователя.
type StoreData interface {
	CreateCredentials(ctx context.Context, userId int64, data []byte, name, description, hash string) error
	CreateCreditCard(ctx context.Context, userId int64, data []byte, name, description, hash string) error
	CreateFileData(ctx context.Context, userId int64, data []byte, name, description, hash string) error
	ChangeData(ctx context.Context, userId int64, lastTimeUpdate time.Time) ([]store.UsersData, error)
}

// Service - структура сервисного слоя.
type Service struct {
	Auth
	StoreAuth
	Data
	StoreData
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

// Service - Конструктор структуры сервисного слоя.
func NewService(db *sqlx.DB, privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey) *Service {
	nDb := postgresql.NewDatabase(db)
	//todo
	return &Service{Auth: authorization.NewAutorization(5*time.Minute, ""),
		StoreAuth:  nDb,
		Data:       NewServiceData(nDb),
		StoreData:  nDb,
		privateKey: privateKey,
		publicKey:  publicKey}
}
