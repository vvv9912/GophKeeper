package service

import (
	"GophKeeper/internal/Agent/server"
	"GophKeeper/pkg/store"
	"context"
	"github.com/jmoiron/sqlx"
	"time"
)

//go:generate mockgen -source=service.go -destination=mocks/service.go -package=mocks
func NewServiceAgent(db *sqlx.DB, key []byte, certFile, keyFile string, serverDns string) *Service {

	return &Service{
		UseCaser: NewUseCase(db, key, certFile, keyFile, serverDns),
	}
}

type AuthService interface {
	SignIn(ctx context.Context, login, password string) (*server.User, error)
	SignUp(ctx context.Context, login, password string) (*server.User, error)
	SetJWTToken(token string)
	GetJWTToken() string
}

type DataInterface interface {
	PostCredentials(ctx context.Context, data *server.ReqData) (*server.RespData, error)
	PostCrateFile(ctx context.Context, data *server.ReqData) (*server.RespData, error)
	PostCreditCard(ctx context.Context, data *server.ReqData) (*server.RespData, error)
	PostCrateFileStartChunks(ctx context.Context, data []byte, fileName string, uuidChunk string, nStart int, nEnd int, maxSize int, reqData []byte) (string, *server.RespData, error)
	GetData(ctx context.Context, userDataId int64) ([]byte, error)
	GetListData(ctx context.Context) ([]byte, error)
	Ping(ctx context.Context) error
	CheckUpdate(ctx context.Context, userDataid int64, updateAt *time.Time) (bool, error)
	PostUpdateData(ctx context.Context, userDataId int64, data []byte) (*server.RespData, error)
	PostUpdateBinaryFile(ctx context.Context, data []byte, fileName string, uuidChunk string, nStart int, nEnd int, maxSize int, reqData []byte, userDataId int64) (string, *server.RespData, error)
}

type StorageData interface {
	CreateFileData(ctx context.Context, data []byte, userDataId int64, name, description, hash string, createdAt *time.Time, UpdateAt *time.Time) error
	CreateCredentials(ctx context.Context, data []byte, userDataId int64, name, description, hash string, createdAt *time.Time, UpdateAt *time.Time) error
	CreateCreditCard(ctx context.Context, data []byte, userDataId int64, name, description, hash string, createdAt *time.Time, UpdateAt *time.Time) error
	GetJWTToken(ctx context.Context) (string, error)
	SetJWTToken(ctx context.Context, JWTToken string) error
	CreateBinaryFile(ctx context.Context, data []byte, userDataId int64, name, description, hash string, createdAt *time.Time, UpdateAt *time.Time, metaData *store.MetaData) error
	GetMetaData(ctx context.Context, userDataId int64) (*store.MetaData, error)
	GetData(ctx context.Context, usersDataId int64) (*store.UsersData, *store.DataFile, error)
	GetInfoData(ctx context.Context, userDataId int64) (*store.UsersData, error)
	UpdateData(ctx context.Context, dataId int64, data []byte, hash string, updateAt *time.Time) error
	UpdateDataBinary(ctx context.Context, userDataId int64, data []byte, hash string, updateAt *time.Time, metaData []byte) error
}

type Encrypter interface {
	Encrypt(data []byte) ([]byte, error)
	Decrypt(data []byte) ([]byte, error)
	EncryptFile(inputFilePath string, outputFilePath string) error
	DecryptFile(inputFilePath string, outputFilePath string) error
}

type UseCaser interface {
	// CreateCredentials - Создание данных логин/пароль
	CreateCredentials(ctx context.Context, data *server.ReqData) error
	// CreateFileData - создание файла
	CreateFileData(ctx context.Context, data *server.ReqData) error
	// CreateCreditCard - создание данных о кредитной карте
	CreateCreditCard(ctx context.Context, data *server.ReqData) error
	// PingServer - пинг сервера
	PingServer(ctx context.Context) bool
	// GetData - получение данных любого формата
	GetData(ctx context.Context, userDataId int64) ([]byte, error)
	// CheckNewData - проверка на новые данные
	CheckNewData(ctx context.Context, userDataId int64) (bool, error)
	// GetDataFromAgentStorage - получение данных из хранилища агента
	GetDataFromAgentStorage(ctx context.Context, userDataId int64) ([]byte, error)
	// GetListData - получение списка актуальных данных пользователя
	GetListData(ctx context.Context) ([]byte, error)
	// UpdateData - обновление данных пользователя (кроме бинарного файла)
	UpdateData(ctx context.Context, userDataId int64, data []byte) ([]byte, error)
	// CreateBinaryFile - создание файла бинарного
	CreateBinaryFile(ctx context.Context, path string, name, description string, ch chan<- string) error
	// UpdateBinaryFile - обновление данных бинарного формата
	UpdateBinaryFile(ctx context.Context, path string, userDataId int64, ch chan<- string) error
	SignIn(ctx context.Context, username, password string) (string, error)
	SignUp(ctx context.Context, username, password string) (string, error)
}
type Service struct {
	UseCaser
}
