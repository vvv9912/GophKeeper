package app

import (
	"GophKeeper/internal/Agent/command"
	"GophKeeper/internal/Agent/service"
	"GophKeeper/pkg/logger"
	"GophKeeper/pkg/store"
	"crypto/rsa"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

type App struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

func Run() error {
	return nil
}
func init() {
	fmt.Println(logger.Initialize("info"))

	logger.Log.Info("start app")

	db, err := sqlx.Open("sqlite", "clientdb.db")
	if err != nil {
		panic(err)
		return
	}
	err = store.MigrateSQLITE(db)
	if err != nil {
		panic(err)
		return
	}
	agent := service.NewServiceAgent(db)
	fmt.Println("Cobra start")
	cob := command.NewCobra(agent)
	if err := cob.Start(); err != nil {
		panic(err)
		return
	}

	return

}
