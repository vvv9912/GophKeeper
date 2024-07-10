package command

import (
	"GophKeeper/internal/Agent/service"
	mock_service "GophKeeper/internal/Agent/service/mocks"
	"context"
	"github.com/golang/mock/gomock"
	"testing"
)

func TestNewCobra(t *testing.T) {
	NewCobra(nil)
}

func TestCobra_initCommand(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUse := mock_service.NewMockUseCaser(ctrl)

	serv := &service.Service{
		UseCaser: mockUse,
	}

	cobr := NewCobra(serv)

	cobr.initCommand(context.Background())

}
