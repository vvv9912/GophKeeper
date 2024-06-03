package service

import (
	"GophKeeper/internal/Agent/server"
)

type AgentService interface{}

type AuthService interface {
	SetJWTToken(token string)
	SetUserId(id int64)
	SignIn(ctx context.Context, login, password string) (*server.User, error)
	SignUp(ctx context.Context, login, password string) (*server.User, error)
}

type DataInterface interface {
	PostCredentials(ctx context.Context, data *server.ReqData) error
	PostCrateFile(ctx context.Context, data *server.ReqData) error
	PostCreditCard(ctx context.Context, data *server.ReqData) error
}

type Service struct {
	AgentService
	AuthService
	DataInterface
	UserId   int64
	JWTToken string
}
