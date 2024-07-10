package command

import (
	"GophKeeper/internal/Agent/service"
	mock_service "GophKeeper/internal/Agent/service/mocks"
	"context"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSuccessfulSignInWithValidCredentials(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUse := mock_service.NewMockUseCaser(ctrl)

	serv := service.Service{mockUse}

	c := NewCobra(&serv)

	args := []string{"validUser", "validPass"}

	mockUse.EXPECT().SignIn(gomock.Any(), "validuser", "validPass").Return("validJWT", nil)
	cobr := cobra.Command{}
	cobr.SetContext(context.Background())

	c.SignIn(&cobr, args)
}
func TestSignInWithValidCredentials(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUse := mock_service.NewMockUseCaser(ctrl)

	serv := service.Service{mockUse}

	c := NewCobra(&serv)

	args := []string{"validUser", "validPass"}

	mockUse.EXPECT().SignIn(gomock.Any(), "validuser", "validPass").Return("validJWT", fmt.Errorf("error"))
	cobr := cobra.Command{}
	cobr.SetContext(context.Background())

	c.SignIn(&cobr, args)
}

func TestCobra_SignUp(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUse := mock_service.NewMockUseCaser(ctrl)

	serv := service.Service{mockUse}

	c := NewCobra(&serv)

	args := []string{"validUser", "validPass"}

	mockUse.EXPECT().SignUp(gomock.Any(), "validuser", "validPass").Return("validJWT", fmt.Errorf("error"))
	cobr := cobra.Command{}
	cobr.SetContext(context.Background())

	c.SignUp(&cobr, args)
}
func TestCobra_SignUp2(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUse := mock_service.NewMockUseCaser(ctrl)

	serv := service.Service{mockUse}

	c := NewCobra(&serv)

	args := []string{"validUser", "validPass"}

	mockUse.EXPECT().SignUp(gomock.Any(), "validuser", "validPass").Return("validJWT", nil)
	cobr := cobra.Command{}
	cobr.SetContext(context.Background())

	c.SignUp(&cobr, args)
}

func TestValidLoginAndPasswordProvidedAsArguments(t *testing.T) {
	// Mock dependencies

	args := []string{"validLogin", "validPassword"}
	login, password, _ := auth(args)

	assert.Equal(t, "validLogin", login)
	assert.Equal(t, "validPassword", password)
}
