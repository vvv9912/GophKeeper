package app

import (
	"GophKeeper/internal/Server/handler"
	"GophKeeper/internal/Server/service"
	"GophKeeper/pkg/logger"
	"GophKeeper/pkg/store"
	"context"
	"crypto/rsa"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"os"
)

type App struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

//	func NewApp(ctx context.Context, connString string) (*App, error) {
//		conn, err := pgx.Connect(ctx, connString)
//		if err != nil {
//			return nil, err
//		}
//		_ = conn
//		return nil, err
//	}
func Run() error {
	fmt.Println(logger.Initialize("info"))
	logger.Log.Info("start app")
	ctx := context.Background()
	//conn, err := pgx.Connect(ctx, "postgres://postgres:postgres@localhost:5434/postgres?sslmode=disable")
	//if err != nil {
	//	//todo Log
	//	return err
	//}
	//_ = conn
	db, err := sqlx.Open("pgx", "postgres://postgres:postgres@localhost:5434/postgres?sslmode=disable")
	if err != nil {
		return err
	}

	err = store.MigratePostgres(db)
	if err != nil {
		return err
	}

	secretKey := string([]byte("asdahgf53sk41250"))
	services, err := service.NewService(db, nil, nil, secretKey)
	if err != nil {
		return err
	}

	h := handler.NewHandler(services)
	cert := "server/cert.pem"
	key := "server/key.pem"
	//cert := "server/server.crt"
	//key := "server/server.key"
	addr := ":8080"
	server := service.StartServer(ctx, h.InitRoutes(services), addr, cert, key)
	_ = server
	ch := make(chan os.Signal, 1)
	<-ch
	return nil
}
