package service

import (
	"GophKeeper/internal/Agent/server"
	"GophKeeper/pkg/store"
	"context"
	"github.com/jmoiron/sqlx"
	"time"
)

//
//go:generate mockgen -source=service.go -destination=mocks/service.go -package=mocks

// NewServiceAgent - Конструктор агента.
func NewServiceAgent(db *sqlx.DB, key []byte, certFile, keyFile string, serverDns string) *Service {

	return &Service{
		UseCaser: NewUseCase(db, key, certFile, keyFile, serverDns),
	}
}

// AuthService - интерфейс авторизации.
type AuthService interface {
	// SignIn - Авторизация пользователя.
	SignIn(ctx context.Context, login, password string) (*server.User, error)
	// SignUp - Регистрация пользователя.
	SignUp(ctx context.Context, login, password string) (*server.User, error)
	// SetJWTToken - Установка JWT токена.
	SetJWTToken(token string)
	// GetJWTToken - Получение JWT токена.
	GetJWTToken() string
}

// DataInterface - интерфейс для работы с данными.
type DataInterface interface {
	// PostCredentials - Запись данных логин/пароль.
	PostCredentials(ctx context.Context, data *server.ReqData) (*server.RespData, error)
	// PostCrateFile - Запись данных о файле.
	PostCrateFile(ctx context.Context, data *server.ReqData) (*server.RespData, error)
	// PostCreditCard - Запись данных о кредитной карте.
	PostCreditCard(ctx context.Context, data *server.ReqData) (*server.RespData, error)
	// PostCrateFileStartChunks - Запись файла по частям.
	PostCrateFileStartChunks(ctx context.Context, data []byte, fileName string, uuidChunk string, nStart int, nEnd int, maxSize int, reqData []byte) (string, *server.RespData, error)
	// GetData - Запрос данных.
	GetData(ctx context.Context, userDataId int64) ([]byte, error)
	// GetListData - Запрос списка данных.
	GetListData(ctx context.Context) ([]byte, error)
	// Ping - Пинг сервера.
	Ping(ctx context.Context) error
	// CheckUpdate - Проверка обновлений.
	CheckUpdate(ctx context.Context, userDataid int64, updateAt *time.Time) (bool, error)
	// PostUpdateData - Обновление данных.
	PostUpdateData(ctx context.Context, userDataId int64, data []byte) (*server.RespData, error)
	// PostUpdateBinaryFile - Обновление файлов бинарных файлов.
	PostUpdateBinaryFile(ctx context.Context, data []byte, fileName string, uuidChunk string, nStart int, nEnd int, maxSize int, reqData []byte, userDataId int64) (string, *server.RespData, error)
}

// StorageData - интерфейс для работы с хранилищем данных.
type StorageData interface {
	// CreateFileData - Создание данных о файле.
	CreateFileData(ctx context.Context, data []byte, userDataId int64, name, description, hash string, createdAt *time.Time, UpdateAt *time.Time) error
	// CreateCredentials - Создание данных логин/пароль.
	CreateCredentials(ctx context.Context, data []byte, userDataId int64, name, description, hash string, createdAt *time.Time, UpdateAt *time.Time) error
	// CreateCreditCard - Создание данных о кредитной карте.
	CreateCreditCard(ctx context.Context, data []byte, userDataId int64, name, description, hash string, createdAt *time.Time, UpdateAt *time.Time) error
	// GetJWTToken - Получение JWT токена.
	GetJWTToken(ctx context.Context) (string, error)
	// SetJWTToken - Установка JWT токена.
	SetJWTToken(ctx context.Context, JWTToken string) error
	// CreateBinaryFile - Создание файлов бинарных файлов.
	CreateBinaryFile(ctx context.Context, data []byte, userDataId int64, name, description, hash string, createdAt *time.Time, UpdateAt *time.Time, metaData *store.MetaData) error
	// GetMetaData - Запрос метаданных.
	GetMetaData(ctx context.Context, userDataId int64) (*store.MetaData, error)
	// GetData - Запрос данных.
	GetData(ctx context.Context, usersDataId int64) (*store.UsersData, *store.DataFile, error)
	// GetInfoData - Запрос информации о данных.
	GetInfoData(ctx context.Context, userDataId int64) (*store.UsersData, error)
	// UpdateData - Обновление данных.
	UpdateData(ctx context.Context, dataId int64, data []byte, hash string, updateAt *time.Time) error
	// UpdateDataBinary - Обновление файлов бинарных файлов.
	UpdateDataBinary(ctx context.Context, userDataId int64, data []byte, hash string, updateAt *time.Time, metaData []byte) error
}

// Encrypter - интерфейс шифрования данных.
type Encrypter interface {
	// Encrypt - шифрование данных.
	Encrypt(data []byte) ([]byte, error)
	// Decrypt - дешифрование данных.
	Decrypt(data []byte) ([]byte, error)
	// EncryptFile - шифрование файлов.
	EncryptFile(inputFilePath string, outputFilePath string) error
	// DecryptFile - дешифрование файлов.
	DecryptFile(inputFilePath string, outputFilePath string) error
}

// UseCaser - интерфейс UseCase.
type UseCaser interface {
	// CreateCredentials - Создание данных логин/пароль.
	CreateCredentials(ctx context.Context, data *server.ReqData) error
	// CreateFileData - создание файла.
	CreateFileData(ctx context.Context, data *server.ReqData) error
	// CreateCreditCard - создание данных о кредитной карте.
	CreateCreditCard(ctx context.Context, data *server.ReqData) error
	// PingServer - пинг сервера.
	PingServer(ctx context.Context) bool
	// GetData - получение данных любого формата.
	GetData(ctx context.Context, userDataId int64) ([]byte, error)
	// CheckNewData - проверка на новые данные.
	CheckNewData(ctx context.Context, userDataId int64) (bool, error)
	// GetDataFromAgentStorage - получение данных из хранилища агента.
	GetDataFromAgentStorage(ctx context.Context, userDataId int64) ([]byte, error)
	// GetListData - получение списка актуальных данных пользователя.
	GetListData(ctx context.Context) ([]byte, error)
	// UpdateData - обновление данных пользователя (кроме бинарного файла).
	UpdateData(ctx context.Context, userDataId int64, data []byte) ([]byte, error)
	// CreateBinaryFile - создание файла бинарного.
	CreateBinaryFile(ctx context.Context, path string, name, description string, ch chan<- string) error
	// UpdateBinaryFile - обновление данных бинарного формата.
	UpdateBinaryFile(ctx context.Context, path string, userDataId int64, ch chan<- string) error
	// SignIn - авторизация пользователя.
	SignIn(ctx context.Context, username, password string) (string, error)
	// SignUp - регистрация пользователя.
	SignUp(ctx context.Context, username, password string) (string, error)
}

// Service - структура сервиса.
type Service struct {
	UseCaser // интерфейс UseCase.
}
