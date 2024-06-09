package service

import (
	"GophKeeper/internal/Agent/server"
	"GophKeeper/pkg/store/sqllite"
	"context"
	"crypto/rsa"
	"crypto/tls"
	"github.com/jmoiron/sqlx"
)

type AgentService interface{}

func NewServiceAgent(db *sqlx.DB) *Service {

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
	PostCrateFileStartChunks(ctx context.Context, data []byte, fileName string, uuidChunk string, nStart int, nEnd int, maxSize int, reqData []byte) (string, error)
	GetData(ctx context.Context, userDataId int64) ([]byte, error)
	GetListData(ctx context.Context) ([]byte, error)
}

type StorageData interface {
	CreateFileData(ctx context.Context, data []byte, userDataId int64, name, description, hash string) error
	CreateCredentials(ctx context.Context, data []byte, userDataId int64, name, description, hash string) error
	CreateCreditCard(ctx context.Context, data []byte, userDataId int64, name, description, hash string) error
	GetJWTToken(ctx context.Context) (string, error)
	SetJWTToken(ctx context.Context, JWTToken string) error
}

type Service struct {
	//AgentService
	AuthService
	DataInterface
	StorageData
	//UserId   int64
	JWTToken string
}
