package main

import (
	"GophKeeper/internal/Server/app"
	"context"
	"os/signal"
	"syscall"
)

func main() {

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	app.Run(ctx)

	<-ctx.Done()
}
