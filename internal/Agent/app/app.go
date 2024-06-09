package app

import (
	"GophKeeper/internal/Agent/command"
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
	agent.CheckNewData(ctx, 15)

	fmt.Println("Cobra off")
	fmt.Println("next")
	return
	s, err := agent.SignIn(ctx, "sadds", "asddsa")
	if err != nil {
		return
	}

	agent.GetData(ctx, 1)
	return
	//agent.CreateFile(ctx, "/home/vlad/Загрузки/customify.0.4.4.zip", " file zip", "test description")
	return
	agent.CreateCredentials(ctx, &server.ReqData{
		Name:        "testName",
		Description: "testDescription",
		Data:        []byte("testData"),
	})
	_ = ctx
	_ = s
	return
}
