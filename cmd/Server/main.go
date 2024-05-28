package main

import (
	"GophKeeper/internal/Server/service"
	"context"
	"os"
)

func main() {

	Run()
	service.StartServer(context.Background())
	ch := make(chan os.Signal, 1)
	<-ch
}
