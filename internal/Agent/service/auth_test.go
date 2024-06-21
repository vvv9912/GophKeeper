package service

import (
	mock_service "GophKeeper/internal/Agent/service/mocks"
	"context"
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

	ss := &Service{
		AuthService:   mockAuth,
		DataInterface: nil,
		StorageData:   mockStorage,
		Encrypter:     nil,
		JWTToken:      "",
	}
	err := ss.setJwtToken(context.TODO())
	require.NoError(t, err)
}
