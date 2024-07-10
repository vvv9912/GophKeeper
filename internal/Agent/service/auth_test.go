package service

import (
	"GophKeeper/internal/Agent/server"
	mock_service "GophKeeper/internal/Agent/service/mocks"
	"context"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSetJwtToken(t *testing.T) {
	c := gomock.NewController(t)
	defer c.Finish()
	//auth
	auth := func(s *mock_service.MockAuthService, ctx context.Context) {
		s.EXPECT().GetJWTToken().Return("")
		s.EXPECT().SetJWTToken("test JWT TOKEN")
	}
	mockAuth := mock_service.NewMockAuthService(c)
	auth(mockAuth, context.TODO())

	//storage
	storage := func(s *mock_service.MockStorageData) {
		s.EXPECT().GetJWTToken(gomock.Any()).Return("test JWT TOKEN", nil)
	}

	mockStorage := mock_service.NewMockStorageData(c)
	storage(mockStorage)
	//

	ss := &UseCase{
		AuthService:   mockAuth,
		DataInterface: nil,
		StorageData:   mockStorage,
		Encrypter:     nil,
		JWTToken:      "",
	}
	err := ss.setJwtToken(context.TODO())
	require.NoError(t, err)
}
func TestSetJwtToken2(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	//auth
	mockDataInterface := mock_service.NewMockDataInterface(ctrl)
	mockEncrypter := mock_service.NewMockEncrypter(ctrl)
	mockStorageData := mock_service.NewMockStorageData(ctrl)
	mockAuthService := mock_service.NewMockAuthService(ctrl)

	mockAuthService.EXPECT().GetJWTToken().Return("")
	mockStorageData.EXPECT().GetJWTToken(gomock.Any()).Return("", fmt.Errorf("error get jwt token"))
	useCase := &UseCase{
		DataInterface: mockDataInterface,
		Encrypter:     mockEncrypter,
		StorageData:   mockStorageData,
		AuthService:   mockAuthService,
	}
	err := useCase.setJwtToken(context.TODO())
	require.Error(t, err)
	require.Equal(t, err.Error(), "error get jwt token")
}
func TestSetJwtToken3(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	//auth
	mockDataInterface := mock_service.NewMockDataInterface(ctrl)
	mockEncrypter := mock_service.NewMockEncrypter(ctrl)
	mockStorageData := mock_service.NewMockStorageData(ctrl)
	mockAuthService := mock_service.NewMockAuthService(ctrl)

	mockAuthService.EXPECT().GetJWTToken().Return("")
	mockStorageData.EXPECT().GetJWTToken(gomock.Any()).Return("", nil)
	useCase := &UseCase{
		DataInterface: mockDataInterface,
		Encrypter:     mockEncrypter,
		StorageData:   mockStorageData,
		AuthService:   mockAuthService,
	}
	err := useCase.setJwtToken(context.TODO())
	require.Error(t, err)
	require.Equal(t, err.Error(), "jwt is empty")
}

func TestUseCase_SignIn(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDataInterface := mock_service.NewMockDataInterface(ctrl)
	mockEncrypter := mock_service.NewMockEncrypter(ctrl)
	mockStorageData := mock_service.NewMockStorageData(ctrl)
	mockAuthService := mock_service.NewMockAuthService(ctrl)

	//mockAuthService.EXPECT().GetJWTToken().Return(";ll;")
	mockAuthService.EXPECT().SignIn(context.TODO(), "test", "test").Return(nil, fmt.Errorf("error sign in"))
	useCase := &UseCase{
		DataInterface: mockDataInterface,
		Encrypter:     mockEncrypter,
		StorageData:   mockStorageData,
		AuthService:   mockAuthService,
	}
	_, err := useCase.SignIn(context.TODO(), "test", "test")
	require.Error(t, err)
	require.Equal(t, err.Error(), "error sign in")
}

func TestUseCase_SignIn2(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDataInterface := mock_service.NewMockDataInterface(ctrl)
	mockEncrypter := mock_service.NewMockEncrypter(ctrl)
	mockStorageData := mock_service.NewMockStorageData(ctrl)
	mockAuthService := mock_service.NewMockAuthService(ctrl)

	//mockAuthService.EXPECT().GetJWTToken().Return(";ll;")
	mockAuthService.EXPECT().SignIn(gomock.Any(), gomock.Any(), gomock.Any()).Return(&server.User{
		Login: "aaaaa",
		JWT:   "bb",
	}, nil)
	mockStorageData.EXPECT().SetJWTToken(gomock.Any(), gomock.Any()).Return(fmt.Errorf("error sign in"))
	mockAuthService.EXPECT().SetJWTToken(gomock.Any())
	useCase := &UseCase{
		DataInterface: mockDataInterface,
		Encrypter:     mockEncrypter,
		StorageData:   mockStorageData,
		AuthService:   mockAuthService,
	}
	_, err := useCase.SignIn(context.TODO(), "test", "test")
	require.Error(t, err)
	require.Equal(t, err.Error(), "error sign in")
}
func TestUseCase_SignIn3(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDataInterface := mock_service.NewMockDataInterface(ctrl)
	mockEncrypter := mock_service.NewMockEncrypter(ctrl)
	mockStorageData := mock_service.NewMockStorageData(ctrl)
	mockAuthService := mock_service.NewMockAuthService(ctrl)

	//mockAuthService.EXPECT().GetJWTToken().Return(";ll;")
	mockAuthService.EXPECT().SignIn(gomock.Any(), gomock.Any(), gomock.Any()).Return(&server.User{
		Login: "aaaaa",
		JWT:   "bb",
	}, nil)
	mockStorageData.EXPECT().SetJWTToken(gomock.Any(), gomock.Any()).Return(nil)
	mockAuthService.EXPECT().SetJWTToken(gomock.Any())
	useCase := &UseCase{
		DataInterface: mockDataInterface,
		Encrypter:     mockEncrypter,
		StorageData:   mockStorageData,
		AuthService:   mockAuthService,
	}
	_, err := useCase.SignIn(context.TODO(), "test", "test")
	require.NoError(t, err)

}
func TestUseCase_SignUp(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDataInterface := mock_service.NewMockDataInterface(ctrl)
	mockEncrypter := mock_service.NewMockEncrypter(ctrl)
	mockStorageData := mock_service.NewMockStorageData(ctrl)
	mockAuthService := mock_service.NewMockAuthService(ctrl)

	//mockAuthService.EXPECT().GetJWTToken().Return(";ll;")
	mockAuthService.EXPECT().SignUp(context.TODO(), "test", "test").Return(nil, fmt.Errorf("error sign in"))
	useCase := &UseCase{
		DataInterface: mockDataInterface,
		Encrypter:     mockEncrypter,
		StorageData:   mockStorageData,
		AuthService:   mockAuthService,
	}
	_, err := useCase.SignUp(context.TODO(), "test", "test")
	require.Error(t, err)
	require.Equal(t, err.Error(), "error sign in")
}

func TestUseCase_SignUp2(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDataInterface := mock_service.NewMockDataInterface(ctrl)
	mockEncrypter := mock_service.NewMockEncrypter(ctrl)
	mockStorageData := mock_service.NewMockStorageData(ctrl)
	mockAuthService := mock_service.NewMockAuthService(ctrl)

	//mockAuthService.EXPECT().GetJWTToken().Return(";ll;")
	mockAuthService.EXPECT().SignUp(gomock.Any(), gomock.Any(), gomock.Any()).Return(&server.User{
		Login: "aaaaa",
		JWT:   "bb",
	}, nil)
	mockStorageData.EXPECT().SetJWTToken(gomock.Any(), gomock.Any()).Return(fmt.Errorf("error sign in"))
	mockAuthService.EXPECT().SetJWTToken(gomock.Any())
	useCase := &UseCase{
		DataInterface: mockDataInterface,
		Encrypter:     mockEncrypter,
		StorageData:   mockStorageData,
		AuthService:   mockAuthService,
	}
	_, err := useCase.SignUp(context.TODO(), "test", "test")
	require.Error(t, err)
	require.Equal(t, err.Error(), "error sign in")
}
func TestUseCase_SignUp3(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDataInterface := mock_service.NewMockDataInterface(ctrl)
	mockEncrypter := mock_service.NewMockEncrypter(ctrl)
	mockStorageData := mock_service.NewMockStorageData(ctrl)
	mockAuthService := mock_service.NewMockAuthService(ctrl)

	//mockAuthService.EXPECT().GetJWTToken().Return(";ll;")
	mockAuthService.EXPECT().SignUp(gomock.Any(), gomock.Any(), gomock.Any()).Return(&server.User{
		Login: "aaaaa",
		JWT:   "bb",
	}, nil)
	mockStorageData.EXPECT().SetJWTToken(gomock.Any(), gomock.Any()).Return(nil)
	mockAuthService.EXPECT().SetJWTToken(gomock.Any())
	useCase := &UseCase{
		DataInterface: mockDataInterface,
		Encrypter:     mockEncrypter,
		StorageData:   mockStorageData,
		AuthService:   mockAuthService,
	}
	_, err := useCase.SignUp(context.TODO(), "test", "test")
	require.NoError(t, err)

}
