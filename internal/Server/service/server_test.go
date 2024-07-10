package service

import (
	"context"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestStartServer_Success(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	handler := http.NewServeMux()
	addr := "127.0.0.1:8080"
	cert := "path/to/valid/cert.pem"
	key := "path/to/valid/key.pem"

	server := StartServer(ctx, handler, addr, cert, key)
	assert.NotNil(t, server)
}
