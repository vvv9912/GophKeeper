package handler

import (
	mock_service "GophKeeper/internal/Agent/service/mocks"
	service2 "GophKeeper/internal/Server/service"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"testing"
)

// todo
func TestHandler_HandlerSignUp(t *testing.T) {
	mock_service.MockUseCaser{}
	s, err := service2.NewService(&sqlx.DB{}, "")
	assert.NoError(t, err)
	h := NewHandler(s)

}
