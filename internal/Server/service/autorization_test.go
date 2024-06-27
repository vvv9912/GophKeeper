package service

import (
	"context"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUseCase_SignUp(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuth := NewMockAuth(ctrl)
	mockStore := NewMockStoreAuth(ctrl)
	storeData := NewMockStoreData(ctrl)
	fileSaver := NewMockFileSaver(ctrl)

	u := &UseCase{
		Auth:      mockAuth,
		StoreAuth: mockStore,
		StoreData: storeData,
		FileSaver: fileSaver,
	}

	mockStore.EXPECT().CreateUser(gomock.Any(), gomock.Any(), gomock.Any()).Return(int64(1), nil)
	mockAuth.EXPECT().BuildJWTString(int64(1)).Return("token", nil)

	_, err := u.SignUp(context.TODO(), "login", "password")
	assert.NoError(t, err)
}
func TestUseCase_SignUpBadCreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuth := NewMockAuth(ctrl)
	mockStore := NewMockStoreAuth(ctrl)
	storeData := NewMockStoreData(ctrl)
	fileSaver := NewMockFileSaver(ctrl)

	u := &UseCase{
		Auth:      mockAuth,
		StoreAuth: mockStore,
		StoreData: storeData,
		FileSaver: fileSaver,
	}

	mockStore.EXPECT().CreateUser(gomock.Any(), gomock.Any(), gomock.Any()).Return(int64(0), fmt.Errorf("error"))

	_, err := u.SignUp(context.TODO(), "login", "password")
	assert.Error(t, err)
}
func TestUseCase_SignUpBadBuildJwt(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuth := NewMockAuth(ctrl)
	mockStore := NewMockStoreAuth(ctrl)
	storeData := NewMockStoreData(ctrl)
	fileSaver := NewMockFileSaver(ctrl)

	u := &UseCase{
		Auth:      mockAuth,
		StoreAuth: mockStore,
		StoreData: storeData,
		FileSaver: fileSaver,
	}

	mockStore.EXPECT().CreateUser(gomock.Any(), gomock.Any(), gomock.Any()).Return(int64(1), nil)
	mockAuth.EXPECT().BuildJWTString(int64(1)).Return("", fmt.Errorf("error"))

	_, err := u.SignUp(context.TODO(), "login", "password")
	assert.Error(t, err)
}

func TestUseCase_SignIn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuth := NewMockAuth(ctrl)
	mockStore := NewMockStoreAuth(ctrl)
	storeData := NewMockStoreData(ctrl)
	fileSaver := NewMockFileSaver(ctrl)

	u := &UseCase{
		Auth:      mockAuth,
		StoreAuth: mockStore,
		StoreData: storeData,
		FileSaver: fileSaver,
	}

	mockStore.EXPECT().GetUserId(gomock.Any(), gomock.Any(), gomock.Any()).Return(int64(1), nil)
	mockAuth.EXPECT().BuildJWTString(int64(1)).Return("token", nil)

	_, err := u.SignIn(context.TODO(), "login", "password")
	assert.NoError(t, err)
}
func TestUseCase_SignInBadCreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuth := NewMockAuth(ctrl)
	mockStore := NewMockStoreAuth(ctrl)
	storeData := NewMockStoreData(ctrl)
	fileSaver := NewMockFileSaver(ctrl)

	u := &UseCase{
		Auth:      mockAuth,
		StoreAuth: mockStore,
		StoreData: storeData,
		FileSaver: fileSaver,
	}

	mockStore.EXPECT().GetUserId(gomock.Any(), gomock.Any(), gomock.Any()).Return(int64(0), fmt.Errorf("error"))

	_, err := u.SignIn(context.TODO(), "login", "password")
	assert.Error(t, err)
}
func TestUseCase_SignInBadBuildJwt(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuth := NewMockAuth(ctrl)
	mockStore := NewMockStoreAuth(ctrl)
	storeData := NewMockStoreData(ctrl)
	fileSaver := NewMockFileSaver(ctrl)

	u := &UseCase{
		Auth:      mockAuth,
		StoreAuth: mockStore,
		StoreData: storeData,
		FileSaver: fileSaver,
	}

	mockStore.EXPECT().GetUserId(gomock.Any(), gomock.Any(), gomock.Any()).Return(int64(1), nil)
	mockAuth.EXPECT().BuildJWTString(int64(1)).Return("", fmt.Errorf("error"))

	_, err := u.SignIn(context.TODO(), "login", "password")
	assert.Error(t, err)
}
