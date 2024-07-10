package handler

import (
	"GophKeeper/internal/Server/service"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewHandler(t *testing.T) {
	s := NewHandler(nil)

	assert.NotNil(t, s)
}

func TestHandler_InitRoutes(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mock_auth := service.NewMockAuth(ctrl)

	s := &service.Service{
		UseCaser: nil,
		Auth:     mock_auth,
	}

	h := NewHandler(s)

	r := h.InitRoutes()

	assert.NotNil(t, r)

}
