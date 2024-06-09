package service

import (
	"GophKeeper/internal/Server/authorization"
	"GophKeeper/pkg/logger"
	"GophKeeper/pkg/store"
	"GophKeeper/pkg/store/postgresql"
	"context"
	"crypto/rsa"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"net/http"
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
	CreateCreditCard(ctx context.Context, userId int64, data []byte, name, description string) (*RespData, error)
	CreateFile(ctx context.Context, userId int64, data []byte, name, description string) (*RespData, error)
	ChangeAllData(ctx context.Context, userId int64, lastTimeUpdate time.Time) ([]byte, error)
	GetData(ctx context.Context, userId int64, userDataId int64) ([]byte, error)
	//UpdateData(ctx context.Context, userId int64, usersData *store.UpdateUsersData, data []byte) error

	RemoveData(ctx context.Context, userId, userDataId int64) error
	UploadFile(additionalPath string, r *http.Request) (bool, *TmpFile, error)
	GetListData(ctx context.Context, userId int64) ([]byte, error)
	UpdateData(ctx context.Context, userId int64, userDataId int64, data []byte) ([]byte, error)

	CreateFileChunks(ctx context.Context, userId int64, tmpFile *TmpFile, name, description string) (*RespData, error)
	GetFileSize(ctx context.Context, userId int64, userDataId int64) ([]byte, error)
	GetFileChunks(ctx context.Context, userId int64, userDataId int64, r *http.Request) ([]byte, error)
}

// StoreData - интерфейс для работы с БД данных пользователя.
type StoreData interface {
	CreateCredentials(ctx context.Context, userId int64, data []byte, name, description, hash string) (*store.UsersData, error)
	CreateCreditCard(ctx context.Context, userId int64, data []byte, name, description, hash string) (*store.UsersData, error)
	CreateFileData(ctx context.Context, userId int64, data []byte, name, description, hash string) (*store.UsersData, error)
	CreateFileDataChunks(ctx context.Context, userId int64, data []byte, name, description, hash string, metaData *store.MetaData) (*store.UsersData, error)
	ChangeAllData(ctx context.Context, userId int64, lastTimeUpdate time.Time) ([]store.UsersData, error)
	ChangeData(ctx context.Context, userId int64, userDataId int64, lastTimeUpdate time.Time) (bool, error)
	GetData(ctx context.Context, userId int64, usersDataId int64) (*store.UsersData, *store.DataFile, error)
	//UpdateData(ctx context.Context, updateData *store.UpdateUsersData, data []byte) error
	UpdateData(ctx context.Context, userId, userDataId int64, data []byte, hash string) (*store.UsersData, error)
	RemoveData(ctx context.Context, userId, usersDataId int64) error
	GetFileSize(ctx context.Context, userId int64, userDataId int64) (int64, error)
	GetMetaData(ctx context.Context, userId, userDataId int64) (*store.MetaData, error)
	GetListData(ctx context.Context, userId int64) ([]store.UsersData, error)
	UpdateBinaryFile(ctx context.Context, userId int64, userDataId int64, data []byte, hash string, metaData []byte) (*store.UsersData, error)
}

// Service - структура сервисного слоя.
type Service struct {
	Auth
	StoreAuth
	StoreData
	SaveFiles
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

// Service - Конструктор структуры сервисного слоя.
func NewService(db *sqlx.DB, privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey, secretKey string) (*Service, error) {
	nDb := postgresql.NewDatabase(db)
	saveFiles, err := NewSaveFiles(10 * time.Minute)

	//todo

	if err != nil {
		logger.Log.Error("Error creating save files", zap.Error(err))
		return nil, err
	}
	return &Service{Auth: authorization.NewAutorization(9000*time.Minute, secretKey),
		StoreAuth:  nDb,
		SaveFiles:  *saveFiles,
		StoreData:  nDb,
		privateKey: privateKey,
		publicKey:  publicKey}, nil
}
