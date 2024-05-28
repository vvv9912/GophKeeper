package app

import (
	"crypto/rsa"
	"github.com/jackc/pgx/v5"
)

type App struct {
	db         *pgx.Conn
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

//
//func NewApp(ctx context.Context, connString string) (*App, error) {
//	conn, err := pgx.Connect(ctx, connString)
//	if err != nil {
//		return nil, err
//	}
//	_ = conn
//	return nil, err
//}

func Run() {
}
