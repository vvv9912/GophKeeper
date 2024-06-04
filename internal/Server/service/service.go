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
	CreateCredentials(ctx context.Context, userId int64, data []byte, name, description string) (*RespData, error)
	CreateCreditCard(ctx context.Context, userId int64, data []byte, name, description string) error
	CreateFile(ctx context.Context, userId int64, data []byte, name, description string) error
	ChangeData(ctx context.Context, userId int64, lastTimeUpdate time.Time) ([]byte, error)
	GetData(ctx context.Context, userId int64, userDataId int64) ([]byte, error)
	UpdateData(ctx context.Context, userId int64, usersData *store.UpdateUsersData, data []byte) error
	RemoveData(ctx context.Context, userId, userDataId int64) error
}

// StoreData - интерфейс для работы с БД данных пользователя.
type StoreData interface {
	CreateCredentials(ctx context.Context, userId int64, data []byte, name, description, hash string) (int64, error)
	CreateCreditCard(ctx context.Context, userId int64, data []byte, name, description, hash string) error
	CreateFileData(ctx context.Context, userId int64, data []byte, name, description, hash string) error
	ChangeData(ctx context.Context, userId int64, lastTimeUpdate time.Time) ([]store.UsersData, error)
	GetData(ctx context.Context, userId int64, usersDataId int64) (*store.UsersData, *store.DataFile, error)
	UpdateData(ctx context.Context, updateData *store.UpdateUsersData, data []byte) error
	RemoveData(ctx context.Context, userId, usersDataId int64) error
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
func NewService(db *sqlx.DB, privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey, secretKey string) *Service {
	nDb := postgresql.NewDatabase(db)
	//todo
	return &Service{Auth: authorization.NewAutorization(9000*time.Minute, secretKey),
		StoreAuth:  nDb,
		Data:       NewServiceData(nDb),
		StoreData:  nDb,
		privateKey: privateKey,
		publicKey:  publicKey}
}
