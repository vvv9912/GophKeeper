package main

import (
	"context"
	"os/signal"
	"syscall"
	"testing"
	"time"
)

func TestSignalNotifications(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		time.Sleep(1 * time.Second)
		cancel()
	}()

	mainCtx, mainCancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer mainCancel()

	select {
	case <-mainCtx.Done():
		t.Log("Context was cancelled")
	case <-time.After(2 * time.Second):
		t.Error("Context was not cancelled as expected")
	}
}
