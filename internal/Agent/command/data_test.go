package command

import (
	"GophKeeper/internal/Agent/service"
	mock_service "GophKeeper/internal/Agent/service/mocks"
	"context"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"testing"
)

func TestCobra_CreateBinaryFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUse := mock_service.NewMockUseCaser(ctrl)

	serv := service.Service{mockUse}

	c := NewCobra(&serv)

	args := []string{"path", "name", "description"}

	mockUse.EXPECT().CreateBinaryFile(gomock.Any(), "path", "name", "description", gomock.Any()).Return(nil)
	cobr := cobra.Command{}
	cobr.SetContext(context.Background())

	c.CreateBinaryFile(&cobr, args)
}
func TestCobra_CreateBinaryFile2(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUse := mock_service.NewMockUseCaser(ctrl)

	serv := service.Service{mockUse}

	c := NewCobra(&serv)

	args := []string{"path", "name"}

	//mockUse.EXPECT().CreateBinaryFile(gomock.Any(), "path", "name", "description", gomock.Any()).Return(nil)
	cobr := cobra.Command{}
	cobr.SetContext(context.Background())

	c.CreateBinaryFile(&cobr, args)
}
func TestCobra_CreateBinaryFile3(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUse := mock_service.NewMockUseCaser(ctrl)

	serv := service.Service{mockUse}

	c := NewCobra(&serv)

	args := []string{"path", "name", "description"}

	mockUse.EXPECT().CreateBinaryFile(gomock.Any(), "path", "name", "description", gomock.Any()).Return(fmt.Errorf("err"))
	cobr := cobra.Command{}
	cobr.SetContext(context.Background())

	c.CreateBinaryFile(&cobr, args)
}

func TestCobra_CreateCredentials(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUse := mock_service.NewMockUseCaser(ctrl)

	serv := service.Service{mockUse}

	c := NewCobra(&serv)

	args := []string{"log", "pass", "desc", "name"}

	mockUse.EXPECT().CreateCredentials(gomock.Any(), gomock.Any()).Return(nil)
	cobr := cobra.Command{}
	cobr.SetContext(context.Background())

	c.CreateCredentials(&cobr, args)
}
func TestCobra_CreateCredentials2(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUse := mock_service.NewMockUseCaser(ctrl)

	serv := service.Service{mockUse}

	c := NewCobra(&serv)

	args := []string{"log", "pass", "desc", "name"}

	mockUse.EXPECT().CreateCredentials(gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
	cobr := cobra.Command{}
	cobr.SetContext(context.Background())

	c.CreateCredentials(&cobr, args)
}
func TestCobra_CreateCredentials3(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUse := mock_service.NewMockUseCaser(ctrl)

	serv := service.Service{mockUse}

	c := NewCobra(&serv)

	args := []string{"log", "pass", "desc"}

	//	mockUse.EXPECT().CreateCredentials(gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
	cobr := cobra.Command{}
	cobr.SetContext(context.Background())

	c.CreateCredentials(&cobr, args)
}

func TestCobra_CreateCreditCard(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUse := mock_service.NewMockUseCaser(ctrl)

	serv := service.Service{mockUse}

	c := NewCobra(&serv)

	args := []string{"name", "desc", "nameBank", "0203", "4578", "055"}

	mockUse.EXPECT().CreateCreditCard(gomock.Any(), gomock.Any()).Return(nil)
	cobr := cobra.Command{}
	cobr.SetContext(context.Background())

	c.CreateCreditCard(&cobr, args)
}
func TestCobra_CreateCreditCard2(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUse := mock_service.NewMockUseCaser(ctrl)

	serv := service.Service{mockUse}

	c := NewCobra(&serv)

	args := []string{"name", "desc", "nameBank", "0203", "4578", "055"}

	mockUse.EXPECT().CreateCreditCard(gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
	cobr := cobra.Command{}
	cobr.SetContext(context.Background())

	c.CreateCreditCard(&cobr, args)
}
func TestCobra_CreateCreditCard3(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUse := mock_service.NewMockUseCaser(ctrl)

	serv := service.Service{mockUse}

	c := NewCobra(&serv)

	args := []string{"name"}

	cobr := cobra.Command{}
	cobr.SetContext(context.Background())

	c.CreateCreditCard(&cobr, args)
}
func TestCobra_CreateCreditCard4(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUse := mock_service.NewMockUseCaser(ctrl)

	serv := service.Service{mockUse}

	c := NewCobra(&serv)

	args := []string{"name", "desc", "nameBank", "a-a-a-", "4578", "055"}

	cobr := cobra.Command{}
	cobr.SetContext(context.Background())

	c.CreateCreditCard(&cobr, args)
}
func TestCobra_CreateCreditCard5(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUse := mock_service.NewMockUseCaser(ctrl)

	serv := service.Service{mockUse}

	c := NewCobra(&serv)

	args := []string{"name", "desc", "nameBank", "0203", "d-d-", "055"}

	cobr := cobra.Command{}
	cobr.SetContext(context.Background())

	c.CreateCreditCard(&cobr, args)
}
func TestCobra_CreateCreditCard6(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUse := mock_service.NewMockUseCaser(ctrl)

	serv := service.Service{mockUse}

	c := NewCobra(&serv)

	args := []string{"name", "desc", "nameBank", "0203", "4555", "-01a"}

	cobr := cobra.Command{}
	cobr.SetContext(context.Background())

	c.CreateCreditCard(&cobr, args)
}

func TestCobra_GetListData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUse := mock_service.NewMockUseCaser(ctrl)

	serv := service.Service{mockUse}

	c := NewCobra(&serv)

	args := []string{}

	mockUse.EXPECT().GetListData(gomock.Any()).Return([]byte("data"), nil)
	cobr := cobra.Command{}
	cobr.SetContext(context.Background())

	c.GetListData(&cobr, args)
}
func TestCobra_GetListData2(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUse := mock_service.NewMockUseCaser(ctrl)

	serv := service.Service{mockUse}

	c := NewCobra(&serv)

	args := []string{}

	mockUse.EXPECT().GetListData(gomock.Any()).Return([]byte("data"), fmt.Errorf("err"))
	cobr := cobra.Command{}
	cobr.SetContext(context.Background())

	c.GetListData(&cobr, args)
}

func TestCobra_GetData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUse := mock_service.NewMockUseCaser(ctrl)

	serv := service.Service{mockUse}

	c := NewCobra(&serv)

	args := []string{"1"}

	mockUse.EXPECT().GetData(gomock.Any(), gomock.Any()).Return([]byte("data"), nil)
	cobr := cobra.Command{}
	cobr.SetContext(context.Background())

	c.GetData(&cobr, args)
}
func TestCobra_GetData2(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUse := mock_service.NewMockUseCaser(ctrl)

	serv := service.Service{mockUse}

	c := NewCobra(&serv)

	args := []string{}

	//	mockUse.EXPECT().GetData(gomock.Any(), gomock.Any()).Return([]byte("data"), nil)
	cobr := cobra.Command{}
	cobr.SetContext(context.Background())

	c.GetData(&cobr, args)
}
func TestCobra_GetData3(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUse := mock_service.NewMockUseCaser(ctrl)

	serv := service.Service{mockUse}

	c := NewCobra(&serv)

	args := []string{"--1"}

	//mockUse.EXPECT().GetData(gomock.Any(), gomock.Any()).Return([]byte("data"), nil)
	cobr := cobra.Command{}
	cobr.SetContext(context.Background())

	c.GetData(&cobr, args)
}
func TestCobra_GetData4(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUse := mock_service.NewMockUseCaser(ctrl)

	serv := service.Service{mockUse}

	c := NewCobra(&serv)

	args := []string{"1"}

	mockUse.EXPECT().GetData(gomock.Any(), gomock.Any()).Return([]byte("data"), fmt.Errorf("err"))
	cobr := cobra.Command{}
	cobr.SetContext(context.Background())

	c.GetData(&cobr, args)
}

func TestCobra_UpdateCredentials(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUse := mock_service.NewMockUseCaser(ctrl)

	serv := service.Service{mockUse}

	c := NewCobra(&serv)

	args := []string{"1", "login", "password"}

	mockUse.EXPECT().UpdateData(gomock.Any(), gomock.Any(), gomock.Any()).Return([]byte("data"), nil)
	cobr := cobra.Command{}
	cobr.SetContext(context.Background())

	c.UpdateCredentials(&cobr, args)
}
func TestCobra_UpdateCredentials2(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUse := mock_service.NewMockUseCaser(ctrl)

	serv := service.Service{mockUse}

	c := NewCobra(&serv)

	args := []string{"1", "login"}

	//mockUse.EXPECT().UpdateData(gomock.Any(), gomock.Any(), gomock.Any()).Return([]byte("data"), nil)
	cobr := cobra.Command{}
	cobr.SetContext(context.Background())

	c.UpdateCredentials(&cobr, args)
}
func TestCobra_UpdateCredentials3(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUse := mock_service.NewMockUseCaser(ctrl)

	serv := service.Service{mockUse}

	c := NewCobra(&serv)

	args := []string{"--a21--", "login", "password"}

	//mockUse.EXPECT().UpdateData(gomock.Any(), gomock.Any(), gomock.Any()).Return([]byte("data"), nil)
	cobr := cobra.Command{}
	cobr.SetContext(context.Background())

	c.UpdateCredentials(&cobr, args)
}
func TestCobra_UpdateCredentials4(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUse := mock_service.NewMockUseCaser(ctrl)

	serv := service.Service{mockUse}

	c := NewCobra(&serv)

	args := []string{"1", "login", "password"}

	mockUse.EXPECT().UpdateData(gomock.Any(), gomock.Any(), gomock.Any()).Return([]byte("data"), fmt.Errorf("err"))
	cobr := cobra.Command{}
	cobr.SetContext(context.Background())

	c.UpdateCredentials(&cobr, args)
}

func TestCobra_UpdateCreditCard(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUse := mock_service.NewMockUseCaser(ctrl)

	serv := service.Service{mockUse}

	c := NewCobra(&serv)

	args := []string{"1", "nameBank", "0203", "4578", "055"}

	mockUse.EXPECT().UpdateData(gomock.Any(), gomock.Any(), gomock.Any()).Return([]byte("data"), nil)
	cobr := cobra.Command{}
	cobr.SetContext(context.Background())

	c.UpdateCreditCard(&cobr, args)
}
func TestCobra_UpdateCreditCard2(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUse := mock_service.NewMockUseCaser(ctrl)

	serv := service.Service{mockUse}

	c := NewCobra(&serv)

	args := []string{"1", "nameBank", "0203", "4578", "055"}

	mockUse.EXPECT().UpdateData(gomock.Any(), gomock.Any(), gomock.Any()).Return([]byte("data"), fmt.Errorf("err"))
	cobr := cobra.Command{}
	cobr.SetContext(context.Background())

	c.UpdateCreditCard(&cobr, args)
}
func TestCobra_UpdateCreditCard3(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUse := mock_service.NewMockUseCaser(ctrl)

	serv := service.Service{mockUse}

	c := NewCobra(&serv)

	args := []string{"0203"}

	cobr := cobra.Command{}
	cobr.SetContext(context.Background())

	c.UpdateCreditCard(&cobr, args)
}
func TestCobra_UpdateCreditCard4(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUse := mock_service.NewMockUseCaser(ctrl)

	serv := service.Service{mockUse}

	c := NewCobra(&serv)

	args := []string{"---a", "nameBank", "0203", "4578", "055"}

	cobr := cobra.Command{}
	cobr.SetContext(context.Background())

	c.UpdateCreditCard(&cobr, args)
}
func TestCobra_UpdateCreditCard5(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUse := mock_service.NewMockUseCaser(ctrl)

	serv := service.Service{mockUse}

	c := NewCobra(&serv)

	args := []string{"1", "nameBank", "---a", "4578", "055"}

	cobr := cobra.Command{}
	cobr.SetContext(context.Background())

	c.UpdateCreditCard(&cobr, args)
}
func TestCobra_UpdateCreditCard6(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUse := mock_service.NewMockUseCaser(ctrl)

	serv := service.Service{mockUse}

	c := NewCobra(&serv)

	args := []string{"1", "nameBank", "0203", "---a", "055"}

	cobr := cobra.Command{}
	cobr.SetContext(context.Background())

	c.UpdateCreditCard(&cobr, args)
}
func TestCobra_UpdateCreditCard7(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUse := mock_service.NewMockUseCaser(ctrl)

	serv := service.Service{mockUse}

	c := NewCobra(&serv)

	args := []string{"1", "nameBank", "0203", "555", "---a"}

	cobr := cobra.Command{}
	cobr.SetContext(context.Background())

	c.UpdateCreditCard(&cobr, args)
}

func TestCobra_UpdateBinaryFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUse := mock_service.NewMockUseCaser(ctrl)

	serv := service.Service{mockUse}

	c := NewCobra(&serv)

	args := []string{"1", "path"}

	cobr := cobra.Command{}
	cobr.SetContext(context.Background())

	mockUse.EXPECT().UpdateBinaryFile(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	c.UpdateBinaryFile(&cobr, args)
}
func TestCobra_UpdateBinaryFile2(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUse := mock_service.NewMockUseCaser(ctrl)

	serv := service.Service{mockUse}

	c := NewCobra(&serv)

	args := []string{"1", "path"}

	cobr := cobra.Command{}
	cobr.SetContext(context.Background())

	mockUse.EXPECT().UpdateBinaryFile(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))

	c.UpdateBinaryFile(&cobr, args)
}
func TestCobra_UpdateBinaryFile3(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUse := mock_service.NewMockUseCaser(ctrl)

	serv := service.Service{mockUse}

	c := NewCobra(&serv)

	args := []string{"path"}

	cobr := cobra.Command{}
	cobr.SetContext(context.Background())

	c.UpdateBinaryFile(&cobr, args)
}
func TestCobra_UpdateBinaryFile4(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUse := mock_service.NewMockUseCaser(ctrl)

	serv := service.Service{mockUse}

	c := NewCobra(&serv)

	args := []string{"---a", "path"}

	cobr := cobra.Command{}
	cobr.SetContext(context.Background())

	c.UpdateBinaryFile(&cobr, args)
}
