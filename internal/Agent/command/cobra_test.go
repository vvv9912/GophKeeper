package command

import (
	"GophKeeper/internal/Agent/service"
	mock_service "GophKeeper/internal/Agent/service/mocks"
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCobra_Start(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUse := mock_service.NewMockUseCaser(ctrl)
	serv := service.Service{mockUse}

	err := NewCobra(&serv).Start(context.Background())
	assert.Error(t, err)
}
