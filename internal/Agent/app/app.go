package app

import (
	"GophKeeper/internal/Agent/server"
	"GophKeeper/internal/Agent/service"
	"GophKeeper/pkg/logger"
	"GophKeeper/pkg/store"
	"context"
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
	ctx := context.Background()
	db, err := sqlx.Open("sqlite", "clientdb.db")
	if err != nil {
		return
	}
	err = store.MigrateSQLITE(db)
	if err != nil {
		return
	}
	agent := service.NewServiceAgent(db)

	//cob := command.NewCobra(agent)
	//if err := cob.Start(); err != nil {
	//	panic(err)
	//	return
	//}

	fmt.Println("next")
	s, err := agent.SignIn(ctx, "sadds", "asddsa")
	if err != nil {
		return
	}

	agent.GetData(ctx, 1)

	agent.CreateFile(ctx, "/home/vlad/Загрузки/customify.0.4.4.zip", " file zip", "test description")

	agent.CreateCredentials(ctx, &server.ReqData{
		Name:        "testName",
		Description: "testDescription",
		Data:        []byte("testData"),
	})
	_ = ctx
	_ = s
	return
}
