package service

import (
	"GophKeeper/internal/Agent/Encrypt"
	"GophKeeper/internal/Agent/server"
	"GophKeeper/pkg/store"
	"GophKeeper/pkg/store/sqllite"
	"context"
	"crypto/rsa"
	"crypto/tls"
	"github.com/jmoiron/sqlx"
	"time"
)

type AgentService interface{}

func NewServiceAgent(db *sqlx.DB, key []byte) *Service {
	//todo config
	cert1, err := tls.LoadX509KeyPair("certs/cert.pem", "certs/key.pem")
	if err != nil {
		panic(err)
	}

	priv := cert1.PrivateKey.(*rsa.PrivateKey)
	pub := priv.PublicKey

	serv := server.NewAgentServer(&pub, priv, "https://localhost:8080")

	return &Service{
		AuthService:   serv,
		DataInterface: serv,
		StorageData:   sqllite.NewDatabase(db),
		Encrypter:     Encrypt.NewEncrypt(key),
	}
}

type AuthService interface {
	SignIn(ctx context.Context, login, password string) (*server.User, error)
	SignUp(ctx context.Context, login, password string) (*server.User, error)
	SetJWTToken(token string)
	GetJWTToken() string
	GetPublicKey() *rsa.PublicKey
	GetPrivateKey() *rsa.PrivateKey
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

type Service struct {
	//AgentService
	AuthService
	DataInterface
	StorageData
	Encrypter
	//UserId   int64
	JWTToken string
}
