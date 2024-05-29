package main

import (
	"GophKeeper/internal/Server/app"
	"GophKeeper/internal/Server/service"
	"context"
	"os"
)

func main() {

	//goose.Up()
	Run()
	app.Run()
	service.StartServer(context.Background(), nil)
	ch := make(chan os.Signal, 1)
	<-ch
}
