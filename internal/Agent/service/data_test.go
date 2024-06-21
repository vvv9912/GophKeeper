package service

import (
	"GophKeeper/internal/Agent/server"
	mock_service "GophKeeper/internal/Agent/service/mocks"
	"context"
	"github.com/golang/mock/gomock"
	"testing"
	"time"
)

func TestCreateFileData(t *testing.T) {

	// Выставить jwt
	auth := func(s *mock_service.MockAuthService, ctx context.Context) {
		s.EXPECT().GetJWTToken().Return("testToken")
	}

	// Выставить encrypt
	encr := func(s *mock_service.MockEncrypter, data []byte) {
		s.EXPECT().Encrypt(gomock.Any()).Return([]byte("EncryptData"), nil)
	}
	c := gomock.NewController(t)
	defer c.Finish()

	//
	// Запрос
	PostD := func(s *mock_service.MockDataInterface, ctx context.Context, req *server.ReqData) {
		s.EXPECT().PostCredentials(gomock.Any(), gomock.Any()).Return(&server.RespData{
			UserDataId: 1,
			Hash:       "testHash",
			CreatedAt:  nil,
			UpdateAt:   nil,
		}, nil)
	}

	//storage
	storage := func(s *mock_service.MockStorageData, ctx context.Context, data []byte, userDataId int64, name, description, hash string, createdAt *time.Time, UpdateAt *time.Time) {
		s.EXPECT().CreateCredentials(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	}
	//Авториз
	mockAuthService := mock_service.NewMockAuthService(c)

	auth(mockAuthService, context.TODO())
	// Шифр
	mockEncryptSer := mock_service.NewMockEncrypter(c)

	encr(mockEncryptSer, []byte("testData"))

	//Запрос

	mockPost := mock_service.NewMockDataInterface(c)
	PostD(mockPost, context.TODO(), &server.ReqData{
		Name:        "test",
		Description: "test description",
		Data:        []byte("test data"),
	})

	//storage
	mockStorage := mock_service.NewMockStorageData(c)
	storage(mockStorage, context.TODO(), []byte("EncryptData"), 1, "testName", "testDescr", "testHash", nil, nil)

	ss := &Service{
		AuthService:   mockAuthService,
		DataInterface: mockPost,
		StorageData:   mockStorage,
		Encrypter:     mockEncryptSer,
		JWTToken:      "",
	}

	reqData := &server.ReqData{
		Name:        "testName",
		Description: "testDescr",
		Data:        []byte("testData"),
	}
	ss.CreateCredentials(context.TODO(), reqData)

}

func TestService_CreateFileData(t *testing.T) {

	// Выставить jwt
	auth := func(s *mock_service.MockAuthService, ctx context.Context) {
		s.EXPECT().GetJWTToken().Return("testToken")
	}

	c := gomock.NewController(t)
	defer c.Finish()

	//
	// Запрос
	PostD := func(s *mock_service.MockDataInterface, ctx context.Context, req *server.ReqData) {
		s.EXPECT().PostCrateFile(gomock.Any(), gomock.Any()).Return(&server.RespData{
			UserDataId: 1,
			Hash:       "testHash",
			CreatedAt:  nil,
			UpdateAt:   nil,
		}, nil)
	}

	//storage
	storage := func(s *mock_service.MockStorageData, ctx context.Context, data []byte, userDataId int64, name, description, hash string, createdAt *time.Time, UpdateAt *time.Time) {
		s.EXPECT().CreateFileData(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	}
	//Авториз
	mockAuthService := mock_service.NewMockAuthService(c)

	auth(mockAuthService, context.TODO())

	//Запрос

	mockPost := mock_service.NewMockDataInterface(c)
	PostD(mockPost, context.TODO(), &server.ReqData{
		Name:        "test",
		Description: "test description",
		Data:        []byte("test data"),
	})

	//storage
	mockStorage := mock_service.NewMockStorageData(c)
	storage(mockStorage, context.TODO(), []byte("EncryptData"), 1, "testName", "testDescr", "testHash", nil, nil)

	ss := &Service{
		AuthService:   mockAuthService,
		DataInterface: mockPost,
		StorageData:   mockStorage,
		Encrypter:     nil,
		JWTToken:      "",
	}

	reqData := &server.ReqData{
		Name:        "testName",
		Description: "testDescr",
		Data:        []byte("testData"),
	}
	ss.CreateFileData(context.TODO(), reqData)

}

func TestService_CreateCreditCard(t *testing.T) {

	// Выставить jwt
	auth := func(s *mock_service.MockAuthService, ctx context.Context) {
		s.EXPECT().GetJWTToken().Return("testToken")
	}

	// Выставить encrypt
	encr := func(s *mock_service.MockEncrypter, data []byte) {
		s.EXPECT().Encrypt(gomock.Any()).Return([]byte("EncryptData"), nil)
	}
	c := gomock.NewController(t)
	defer c.Finish()

	//
	// Запрос
	PostD := func(s *mock_service.MockDataInterface, ctx context.Context, req *server.ReqData) {
		s.EXPECT().PostCreditCard(gomock.Any(), gomock.Any()).Return(&server.RespData{
			UserDataId: 1,
			Hash:       "testHash",
			CreatedAt:  nil,
			UpdateAt:   nil,
		}, nil)
	}

	//storage
	storage := func(s *mock_service.MockStorageData, ctx context.Context, data []byte, userDataId int64, name, description, hash string, createdAt *time.Time, UpdateAt *time.Time) {
		s.EXPECT().CreateCreditCard(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	}
	//Авториз
	mockAuthService := mock_service.NewMockAuthService(c)

	auth(mockAuthService, context.TODO())
	// Шифр
	mockEncryptSer := mock_service.NewMockEncrypter(c)

	encr(mockEncryptSer, []byte("testData"))

	//Запрос

	mockPost := mock_service.NewMockDataInterface(c)
	PostD(mockPost, context.TODO(), &server.ReqData{
		Name:        "test",
		Description: "test description",
		Data:        []byte("test data"),
	})

	//storage
	mockStorage := mock_service.NewMockStorageData(c)
	storage(mockStorage, context.TODO(), []byte("EncryptData"), 1, "testName", "testDescr", "testHash", nil, nil)

	ss := &Service{
		AuthService:   mockAuthService,
		DataInterface: mockPost,
		StorageData:   mockStorage,
		Encrypter:     mockEncryptSer,
		JWTToken:      "",
	}

	reqData := &server.ReqData{
		Name:        "testName",
		Description: "testDescr",
		Data:        []byte("testData"),
	}
	ss.CreateCreditCard(context.TODO(), reqData)
}
