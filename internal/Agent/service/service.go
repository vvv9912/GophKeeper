package service

import (
	"GophKeeper/internal/Agent/server"
	"GophKeeper/pkg/store/sqllite"
	"context"
	"github.com/jmoiron/sqlx"
)

type AgentService interface{}

func NewServiceAgent(db *sqlx.DB) *Service {
	serv := server.NewAgentServer(nil, nil, "https://localhost:8080")

	return &Service{
		AuthService:   serv,
		DataInterface: serv,
		StorageData:   sqllite.NewDatabase(db),
	}
}

type AuthService interface {
	SetJWTToken(token string)
	SignIn(ctx context.Context, login, password string) (*server.User, error)
	SignUp(ctx context.Context, login, password string) (*server.User, error)
	GetJWTToken() string
}

type DataInterface interface {
	PostCredentials(ctx context.Context, data *server.ReqData) (*server.RespData, error)
	PostCrateFile(ctx context.Context, data *server.ReqData) (*server.RespData, error)
	PostCreditCard(ctx context.Context, data *server.ReqData) (*server.RespData, error)
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
