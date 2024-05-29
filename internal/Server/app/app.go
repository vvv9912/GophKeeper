package app

import (
	"GophKeeper/pkg/store"
	"GophKeeper/pkg/store/postgresql"
	"context"
	"crypto/rsa"
	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
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

func Run() error {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, "postgres://postgres:postgres@localhost:5434/postgres?sslmode=disable")
	if err != nil {
		//todo Log
		return err
	}
	_ = conn
	db, err := sqlx.Open("pgx", "postgres://postgres:postgres@localhost:5434/postgres?sslmode=disable")

	ss := postgresql.NewDatabase(db)
	ss.CreateUser("login2", "pass")
	err = store.Migrate(db)
	if err != nil {
		return err
	}
	return nil
}
