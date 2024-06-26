package handler

import (
	service2 "GophKeeper/internal/Server/service"
	mock_service2 "GophKeeper/internal/Server/service/mocks"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"testing"
)

// todo
func TestHandler_HandlerSignUp(t *testing.T) {
	mock_service2.New
	s, err := service2.NewService(&sqlx.DB{}, "")
	assert.NoError(t, err)
	h := NewHandler(s)

}
