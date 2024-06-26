package service

import (
	"GophKeeper/pkg/store"
	"context"
	"github.com/jmoiron/sqlx"
	"net/http"
	"time"
)

//
//go:generate mockgen -source=service.go -destination=mocks/mock.go

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

//// Data - интерфейс для работы с данными пользователя.
//type Data interface {
//	CreateCredentials(ctx context.Context, userId int64, data []byte, name, description string) (*RespData, error)
//	CreateCreditCard(ctx context.Context, userId int64, data []byte, name, description string) (*RespData, error)
//	CreateFile(ctx context.Context, userId int64, data []byte, name, description string) (*RespData, error)
//	ChangeAllData(ctx context.Context, userId int64, lastTimeUpdate time.Time) ([]byte, error)
//	GetData(ctx context.Context, userId int64, userDataId int64) ([]byte, error)
//
//	RemoveData(ctx context.Context, userId, userDataId int64) error
//	UploadFile(additionalPath string, r *http.Request) (bool, *TmpFile, error)
//	GetListData(ctx context.Context, userId int64) ([]byte, error)
//	UpdateData(ctx context.Context, userId int64, userDataId int64, data []byte) ([]byte, error)
//
//	CreateFileChunks(ctx context.Context, userId int64, tmpFile *TmpFile, name, description string) (*RespData, error)
//	GetFileSize(ctx context.Context, userId int64, userDataId int64) ([]byte, error)
//	GetFileChunks(ctx context.Context, userId int64, userDataId int64, r *http.Request) ([]byte, error)
//}

// StoreData - интерфейс для работы с БД данных пользователя.
type StoreData interface {
	CreateCredentials(ctx context.Context, userId int64, data []byte, name, description, hash string) (*store.UsersData, error)
	CreateCreditCard(ctx context.Context, userId int64, data []byte, name, description, hash string) (*store.UsersData, error)
	CreateFileData(ctx context.Context, userId int64, data []byte, name, description, hash string) (*store.UsersData, error)
	CreateFileDataChunks(ctx context.Context, userId int64, data []byte, name, description, hash string, metaData *store.MetaData) (*store.UsersData, error)
	ChangeAllData(ctx context.Context, userId int64, lastTimeUpdate time.Time) ([]store.UsersData, error)
	ChangeData(ctx context.Context, userId int64, userDataId int64, lastTimeUpdate time.Time) (bool, error)
	GetData(ctx context.Context, userId int64, usersDataId int64) (*store.UsersData, *store.DataFile, error)

	UpdateData(ctx context.Context, userId, userDataId int64, data []byte, hash string) (*store.UsersData, error)
	RemoveData(ctx context.Context, userId, usersDataId int64) error
	GetFileSize(ctx context.Context, userId int64, userDataId int64) (int64, error)
	GetMetaData(ctx context.Context, userId, userDataId int64) (*store.MetaData, error)
	GetListData(ctx context.Context, userId int64) ([]store.UsersData, error)
	UpdateBinaryFile(ctx context.Context, userId int64, userDataId int64, data []byte, hash string, metaData []byte) (*store.UsersData, error)
}
type UseCaser interface {
	// SignUp - регистрация пользователя.
	SignUp(ctx context.Context, login, password string) (string, error)
	// SignIn - авторизация пользователя.
	SignIn(ctx context.Context, login, password string) (string, error)
	// CreateCredentials - Создание пары логин/пароль.
	CreateCredentials(ctx context.Context, userId int64, data []byte, name, description string) (*RespData, error)
	// CreateCreditCard - Создание пары данные банковских карт.
	CreateCreditCard(ctx context.Context, userId int64, data []byte, name, description string) (*RespData, error)
	// CreateFileChunks - Создание бинарных данных.
	CreateFileChunks(ctx context.Context, userId int64, tmpFile *TmpFile, name, description string, encryptedData []byte) (*RespData, error)
	// CreateFile - Создание  данных (файл).
	CreateFile(ctx context.Context, userId int64, data []byte, name, description string) (*RespData, error)
	// ChangeData - проверка изменения данных.
	ChangeData(ctx context.Context, userId int64, userDataId int64, lastTimeUpdate time.Time) ([]byte, error)
	// ChangeAllData - список изменненых данных.
	ChangeAllData(ctx context.Context, userId int64, lastTimeUpdate time.Time) ([]byte, error)
	// GetFileSize - получение размера бинарного файла.
	GetFileSize(ctx context.Context, userId int64, userDataId int64) ([]byte, error)
	// GetFileChunks - получение чанков бинарного файла.
	GetFileChunks(ctx context.Context, userId int64, userDataId int64, r *http.Request) ([]byte, error)
	// GetData - получение данных.
	GetData(ctx context.Context, userId int64, userDataId int64) ([]byte, error)
	// UpdateData - обновление данных.
	UpdateData(ctx context.Context, userId int64, userDataId int64, data []byte) ([]byte, error)
	// RemoveData - удаление данных (выставление флага в бд).
	RemoveData(ctx context.Context, userId, userDataId int64) error
	// GetListData - получение списка данных для пользователя.
	GetListData(ctx context.Context, userId int64) ([]byte, error)
	// UploadFile - загрузка файла.
	UploadFile(additionalPath string, r *http.Request) (bool, *TmpFile, error)
	// UpdateBinaryFile - Обновление бинарных данных.
	UpdateBinaryFile(ctx context.Context, userId int64, userDataId int64, tmpFile *TmpFile, encryptedData []byte) (*RespData, error)
}

// Service - структура сервисного слоя.
type Service struct {
	UseCaser // интерфейс UseCase. //todo add interface
	Auth
	StoreAuth
}

// NewService - Конструктор структуры сервисного слоя.
func NewService(db *sqlx.DB, secretKey string) (*Service, error) {
	u, err := NewUseCase(db, secretKey)
	if err != nil {
		return nil, err
	}
	return &Service{
		UseCaser: u,
	}, nil
}
