package app

import (
	"GophKeeper/internal/Server/config"
	"GophKeeper/internal/Server/handler"
	"GophKeeper/internal/Server/service"
	"GophKeeper/pkg/logger"
	"GophKeeper/pkg/store"
	"context"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func Run(ctx context.Context) error {

	if err := config.InitConfig(); err != nil {
		panic(err)
	}
	if err := logger.Initialize(config.Get().LevelLogger); err != nil {
		panic(err)
	}
	logger.Log.Info("start app")

	db, err := sqlx.Open("pgx", config.Get().DatabaseDNS)
	if err != nil {
		return err
	}

	err = store.MigratePostgres(db)
	if err != nil {
		return err
	}

	services, err := service.NewService(db, config.Get().SecretKey)
	if err != nil {
		return err
	}

	h := handler.NewHandler(services)

	service.StartServer(ctx, h.InitRoutes(services), config.Get().ServerDNS, config.Get().CertFile, config.Get().KeyFile)

	return nil
}
