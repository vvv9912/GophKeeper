package main

import (
	"GophKeeper/internal/Agent/app"
	"context"
	"log"
	"os/signal"
	"syscall"
)

var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
)

func main() {
	log.Println("Build version:", buildVersion)
	log.Println("Build date:", buildDate)
	log.Println("Build commit:", buildCommit)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	if err := app.Run(ctx); err != nil {
		log.Fatal("Error during app.Run:", err)
	}

	<-ctx.Done()

}
