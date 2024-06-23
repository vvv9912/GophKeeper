package service

import (
	"GophKeeper/internal/Agent/model"
	"GophKeeper/internal/Agent/server"
	mock_service "GophKeeper/internal/Agent/service/mocks"
	"GophKeeper/pkg/store"
	"GophKeeper/pkg/store/sqllite"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math/big"
	"os"
	"testing"
	"time"
)

// mockData := &mock_service.MockStorageData{}
// mockData.On("CreateCredentials", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
func TestService_CreateCredentials(t *testing.T) {
	type mockAuth func(s *mock_service.MockAuthService, ctx context.Context)
	type mockDataI func(s *mock_service.MockDataInterface, ctx context.Context, req *server.ReqData)
	type mockStorage func(s *mock_service.MockStorageData)
	type mockEncrypt func(s *mock_service.MockEncrypter, data []byte)

	type args struct {
		mockAuth    mockAuth
		mockDataI   mockDataI
		mockStorage mockStorage
		mockEncrypt mockEncrypt
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				mockAuth: func(s *mock_service.MockAuthService, ctx context.Context) {
					s.EXPECT().GetJWTToken().Return("testToken")
				},
				mockDataI: func(s *mock_service.MockDataInterface, ctx context.Context, req *server.ReqData) {
					s.EXPECT().PostCredentials(gomock.Any(), gomock.Any()).Return(&server.RespData{
						UserDataId: 1,
						Hash:       "testHash",
						CreatedAt:  nil,
						UpdateAt:   nil,
					}, nil)
				},
				mockStorage: func(s *mock_service.MockStorageData) {
					s.EXPECT().CreateCredentials(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				},
				mockEncrypt: func(s *mock_service.MockEncrypter, data []byte) {
					s.EXPECT().Encrypt(gomock.Any()).Return([]byte("EncryptData"), nil)
				},
			},
			wantErr: false,
		},
		{
			name: "false err mockDataInterface",
			args: args{
				mockAuth: func(s *mock_service.MockAuthService, ctx context.Context) {
					s.EXPECT().GetJWTToken().Return("testToken")
				},
				mockDataI: func(s *mock_service.MockDataInterface, ctx context.Context, req *server.ReqData) {
					s.EXPECT().PostCredentials(gomock.Any(), gomock.Any()).Return(nil, errors.New("Error custom"))
				},
				mockStorage: func(s *mock_service.MockStorageData) {
				},
				mockEncrypt: func(s *mock_service.MockEncrypter, data []byte) {
					s.EXPECT().Encrypt(gomock.Any()).Return([]byte("EncryptData"), nil)
				},
			},
			wantErr: true,
		},
		{
			name: "false err mockStorage",
			args: args{
				mockAuth: func(s *mock_service.MockAuthService, ctx context.Context) {
					s.EXPECT().GetJWTToken().Return("testToken")
				},
				mockDataI: func(s *mock_service.MockDataInterface, ctx context.Context, req *server.ReqData) {
					s.EXPECT().PostCredentials(gomock.Any(), gomock.Any()).Return(&server.RespData{
						UserDataId: 1,
						Hash:       "testHash",
						CreatedAt:  nil,
						UpdateAt:   nil,
					}, nil)
				},
				mockStorage: func(s *mock_service.MockStorageData) {
					s.EXPECT().CreateCredentials(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("Error custom"))
				},
				mockEncrypt: func(s *mock_service.MockEncrypter, data []byte) {
					s.EXPECT().Encrypt(gomock.Any()).Return([]byte("EncryptData"), nil)
				},
			},
			wantErr: true,
		},
		{
			name: "false err mockEncrypt",
			args: args{
				mockAuth: func(s *mock_service.MockAuthService, ctx context.Context) {
					s.EXPECT().GetJWTToken().Return("testToken")
				},
				mockDataI: func(s *mock_service.MockDataInterface, ctx context.Context, req *server.ReqData) {
					//s.EXPECT().PostCredentials(gomock.Any(), gomock.Any()).Return(nil, errors.New("Error custom"))
				},
				mockStorage: func(s *mock_service.MockStorageData) {
					//s.EXPECT().CreateCredentials(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				},
				mockEncrypt: func(s *mock_service.MockEncrypter, data []byte) {
					s.EXPECT().Encrypt(gomock.Any()).Return([]byte("EncryptData"), errors.New("Error custom"))
				},
			},
			wantErr: true,
		},
		// TODO: Add test cases.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			//Авториз
			mockAuthService := mock_service.NewMockAuthService(ctrl)

			tt.args.mockAuth(mockAuthService, context.TODO())
			// Шифр
			mockEncryptSer := mock_service.NewMockEncrypter(ctrl)

			tt.args.mockEncrypt(mockEncryptSer, []byte("testData"))

			//Запрос

			mockPost := mock_service.NewMockDataInterface(ctrl)
			tt.args.mockDataI(mockPost, context.TODO(), &server.ReqData{
				Name:        "test",
				Description: "test description",
				Data:        []byte("test data"),
			})

			//storage
			mockStoragee := mock_service.NewMockStorageData(ctrl)
			tt.args.mockStorage(mockStoragee)

			ss := &UseCase{
				AuthService:   mockAuthService,
				DataInterface: mockPost,
				StorageData:   mockStoragee,
				Encrypter:     mockEncryptSer,
				JWTToken:      "",
			}

			reqData := &server.ReqData{
				Name:        "testName",
				Description: "testDescr",
				Data:        []byte("testData"),
			}

			if err := ss.CreateCredentials(context.TODO(), reqData); (err != nil) != tt.wantErr {
				t.Errorf("CreateCredentials() error = %v, wantErr %v", err, tt.wantErr)
			}

		})
	}
}
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

	ss := &UseCase{
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
	err := ss.CreateCredentials(context.TODO(), reqData)
	require.NoError(t, err)
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

	ss := &UseCase{
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
	err := ss.CreateFileData(context.TODO(), reqData)
	require.NoError(t, err)
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

	ss := &UseCase{
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
	err := ss.CreateCreditCard(context.TODO(), reqData)
	require.NoError(t, err)
}

func TestService_PingServer(t *testing.T) {

	//
	// Запрос
	PostD := func(s *mock_service.MockDataInterface) {
		s.EXPECT().Ping(gomock.Any()).Return(nil)
	}
	c := gomock.NewController(t)
	defer c.Finish()
	//Запрос

	mockPost := mock_service.NewMockDataInterface(c)
	PostD(mockPost)

	ss := &UseCase{
		AuthService:   nil,
		DataInterface: mockPost,
		StorageData:   nil,
		Encrypter:     nil,
		JWTToken:      "",
	}

	err := ss.Ping(context.TODO())
	require.NoError(t, err)
}

func TestGetData_ServerUnavailable_LocalStorageFails(t *testing.T) {
	ctx := context.TODO()
	userDataId := int64(1)
	c := gomock.NewController(t)
	defer c.Finish()

	mockAuthService := mock_service.NewMockAuthService(c)
	mockDataInterface := mock_service.NewMockDataInterface(c)
	mockStorageData := mock_service.NewMockStorageData(c)
	mockEncrypter := mock_service.NewMockEncrypter(c)

	mockDataInterface.EXPECT().Ping(ctx).Return(errors.New("server unavailable"))
	mockStorageData.EXPECT().GetData(ctx, userDataId).Return(nil, nil, errors.New("local storage error"))

	useCase := &UseCase{
		AuthService:   mockAuthService,
		DataInterface: mockDataInterface,
		StorageData:   mockStorageData,
		Encrypter:     mockEncrypter,
	}
	s := Service{
		UseCaser: useCase,
	}

	result, err := s.GetData(ctx, userDataId)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "local storage error", err.Error())

}
func TestGetData_ServerUnavailable_jwtFailed(t *testing.T) {
	ctx := context.TODO()
	userDataId := int64(1)
	c := gomock.NewController(t)
	defer c.Finish()

	mockAuthService := mock_service.NewMockAuthService(c)
	mockDataInterface := mock_service.NewMockDataInterface(c)
	mockStorageData := mock_service.NewMockStorageData(c)
	mockEncrypter := mock_service.NewMockEncrypter(c)

	mockDataInterface.EXPECT().Ping(ctx).Return(nil)
	mockAuthService.EXPECT().GetJWTToken().Return("")
	mockStorageData.EXPECT().GetJWTToken(ctx).Return("", errors.New("jwt failed"))

	useCase := &UseCase{
		AuthService:   mockAuthService,
		DataInterface: mockDataInterface,
		StorageData:   mockStorageData,
		Encrypter:     mockEncrypter,
	}
	s := Service{
		UseCaser: useCase,
	}

	result, err := s.GetData(ctx, userDataId)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "jwt failed", err.Error())

}

func TestGetData_ServerFailedCheckNewData(t *testing.T) {
	ctx := context.TODO()
	userDataId := int64(1)
	c := gomock.NewController(t)
	defer c.Finish()

	mockAuthService := mock_service.NewMockAuthService(c)
	mockDataInterface := mock_service.NewMockDataInterface(c)
	mockStorageData := mock_service.NewMockStorageData(c)
	mockEncrypter := mock_service.NewMockEncrypter(c)

	mockDataInterface.EXPECT().Ping(ctx).Return(nil)
	mockAuthService.EXPECT().GetJWTToken().Return(";ll;")
	mockStorageData.EXPECT().GetInfoData(ctx, userDataId).Return(nil, errors.New("Failed update"))
	useCase := &UseCase{
		AuthService:   mockAuthService,
		DataInterface: mockDataInterface,
		StorageData:   mockStorageData,
		Encrypter:     mockEncrypter,
	}
	s := Service{
		UseCaser: useCase,
	}

	result, err := s.GetData(ctx, userDataId)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "Failed update", err.Error())

}
func TestGetData_ServerFailedGetDataInt(t *testing.T) {
	ctx := context.TODO()
	userDataId := int64(1)
	ttt := time.Now()
	c := gomock.NewController(t)
	defer c.Finish()

	mockAuthService := mock_service.NewMockAuthService(c)
	mockDataInterface := mock_service.NewMockDataInterface(c)
	mockStorageData := mock_service.NewMockStorageData(c)
	mockEncrypter := mock_service.NewMockEncrypter(c)

	mockDataInterface.EXPECT().Ping(ctx).Return(nil)
	mockAuthService.EXPECT().GetJWTToken().Return(";ll;").AnyTimes()
	mockStorageData.EXPECT().GetInfoData(ctx, userDataId).Return(&store.UsersData{UpdateAt: &ttt}, nil)
	mockDataInterface.EXPECT().CheckUpdate(ctx, userDataId, gomock.Any()).Return(true, nil)
	mockDataInterface.EXPECT().GetData(ctx, userDataId).Return(nil, errors.New("Failed update"))

	useCase := &UseCase{
		AuthService:   mockAuthService,
		DataInterface: mockDataInterface,
		StorageData:   mockStorageData,
		Encrypter:     mockEncrypter,
	}
	s := Service{
		UseCaser: useCase,
	}

	result, err := s.GetData(ctx, userDataId)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "Failed update", err.Error())

}
func TestGetData_ServerFailedDecr(t *testing.T) {
	ctx := context.TODO()
	userDataId := int64(1)
	ttt := time.Now()
	c := gomock.NewController(t)
	defer c.Finish()

	mockAuthService := mock_service.NewMockAuthService(c)
	mockDataInterface := mock_service.NewMockDataInterface(c)
	mockStorageData := mock_service.NewMockStorageData(c)
	mockEncrypter := mock_service.NewMockEncrypter(c)

	mockDataInterface.EXPECT().Ping(ctx).Return(nil)
	mockAuthService.EXPECT().GetJWTToken().Return(";ll;").AnyTimes()
	mockStorageData.EXPECT().GetInfoData(ctx, userDataId).Return(&store.UsersData{UpdateAt: &ttt}, nil)
	mockDataInterface.EXPECT().CheckUpdate(ctx, userDataId, gomock.Any()).Return(true, nil)
	mockDataInterface.EXPECT().GetData(ctx, userDataId).Return([]byte("data"), nil)
	mockEncrypter.EXPECT().Decrypt(gomock.Any()).Return(nil, errors.New("Failed decr"))

	useCase := &UseCase{
		AuthService:   mockAuthService,
		DataInterface: mockDataInterface,
		StorageData:   mockStorageData,
		Encrypter:     mockEncrypter,
	}
	s := Service{
		UseCaser: useCase,
	}

	result, err := s.GetData(ctx, userDataId)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "Failed decr", err.Error())

}

func TestCheckNewData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorageData := mock_service.NewMockStorageData(ctrl)
	mockAuthService := mock_service.NewMockAuthService(ctrl)
	mockDataInterface := mock_service.NewMockDataInterface(ctrl)
	mockEncrypter := mock_service.NewMockEncrypter(ctrl)

	userDataId := int64(123)
	ttt := time.Now()
	ctx := context.TODO()

	mockStorageData.EXPECT().GetInfoData(ctx, userDataId).Return(&store.UsersData{UpdateAt: &ttt}, nil)
	mockAuthService.EXPECT().GetJWTToken().Return("test JWT token").AnyTimes()
	mockDataInterface.EXPECT().CheckUpdate(ctx, userDataId, gomock.Any()).Return(true, nil)

	useCase := &UseCase{
		StorageData:   mockStorageData,
		AuthService:   mockAuthService,
		DataInterface: mockDataInterface,
		Encrypter:     mockEncrypter,
	}

	ok, err := useCase.CheckNewData(ctx, userDataId)

	assert.NoError(t, err)
	assert.True(t, ok)

}
func TestCheckErrNewData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorageData := mock_service.NewMockStorageData(ctrl)
	mockAuthService := mock_service.NewMockAuthService(ctrl)
	mockDataInterface := mock_service.NewMockDataInterface(ctrl)
	mockEncrypter := mock_service.NewMockEncrypter(ctrl)

	userDataId := int64(123)

	ctx := context.TODO()

	mockStorageData.EXPECT().GetInfoData(ctx, userDataId).Return(nil, fmt.Errorf("err get info"))
	//mockAuthService.EXPECT().GetJWTToken().Return("test JWT token").AnyTimes()
	//mockDataInterface.EXPECT().CheckUpdate(ctx, userDataId, gomock.Any()).Return(true, nil)

	useCase := &UseCase{
		StorageData:   mockStorageData,
		AuthService:   mockAuthService,
		DataInterface: mockDataInterface,
		Encrypter:     mockEncrypter,
	}

	ok, err := useCase.CheckNewData(ctx, userDataId)

	assert.Error(t, err)
	assert.Equal(t, "err get info", err.Error())
	assert.Equal(t, ok, false)

}
func TestGetDataFromAgentStorage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorageData := mock_service.NewMockStorageData(ctrl)
	mockEncrypter := mock_service.NewMockEncrypter(ctrl)

	userDataId := int64(123)
	ctx := context.TODO()

	// Expected calls
	mockStorageData.EXPECT().GetData(ctx, userDataId).Return(nil, nil, fmt.Errorf("error getting data"))
	//mockEncrypter.EXPECT().Decrypt(gomock.Any()).Return(nil, fmt.Errorf("error decrypting"))
	//mockStorageData.EXPECT().GetMetaData(ctx, userDataId).Return(nil, fmt.Errorf("error getting metadata"))

	useCase := &UseCase{
		StorageData: mockStorageData,
		Encrypter:   mockEncrypter,
	}

	_, err := useCase.GetDataFromAgentStorage(ctx, userDataId)

	assert.Error(t, err)
}
func TestGetDataFromAgentStorage2(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorageData := mock_service.NewMockStorageData(ctrl)
	mockEncrypter := mock_service.NewMockEncrypter(ctrl)

	userDataId := int64(123)
	ctx := context.TODO()

	mockStorageData.EXPECT().GetData(ctx, userDataId).Return(&store.UsersData{}, &store.DataFile{}, nil)
	mockEncrypter.EXPECT().Decrypt(gomock.Any()).Return(nil, fmt.Errorf("error decrypting"))
	//mockStorageData.EXPECT().GetMetaData(ctx, userDataId).Return(nil, fmt.Errorf("error getting metadata"))

	useCase := &UseCase{
		StorageData: mockStorageData,
		Encrypter:   mockEncrypter,
	}

	_, err := useCase.GetDataFromAgentStorage(ctx, userDataId)

	assert.Error(t, err)
}
func TestGetDataFromAgentStorage3(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorageData := mock_service.NewMockStorageData(ctrl)
	mockEncrypter := mock_service.NewMockEncrypter(ctrl)

	userDataId := int64(123)
	ctx := context.TODO()

	mockStorageData.EXPECT().GetData(ctx, userDataId).Return(&store.UsersData{}, &store.DataFile{}, nil)
	mockEncrypter.EXPECT().Decrypt(gomock.Any()).Return([]byte("test data"), nil)

	useCase := &UseCase{
		StorageData: mockStorageData,
		Encrypter:   mockEncrypter,
	}

	_, err := useCase.GetDataFromAgentStorage(ctx, userDataId)

	assert.NoError(t, err)
}

func TestUseCase_GetListData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorageData := mock_service.NewMockStorageData(ctrl)
	mockAuthService := mock_service.NewMockAuthService(ctrl)
	mockDataInterface := mock_service.NewMockDataInterface(ctrl)
	mockEncrypter := mock_service.NewMockEncrypter(ctrl)

	ctx := context.TODO()

	mockAuthService.EXPECT().GetJWTToken().Return(";ll;")
	mockDataInterface.EXPECT().GetListData(ctx).Return([]byte("test data"), nil)

	useCase := &UseCase{
		StorageData:   mockStorageData,
		AuthService:   mockAuthService,
		DataInterface: mockDataInterface,
		Encrypter:     mockEncrypter,
	}

	d, err := useCase.GetListData(ctx)

	assert.Equal(t, "test data", string(d))
	assert.NoError(t, err)

}

// err test
func TestUpdateData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDataInterface := mock_service.NewMockDataInterface(ctrl)
	mockEncrypter := mock_service.NewMockEncrypter(ctrl)
	mockStorageData := mock_service.NewMockStorageData(ctrl)
	mockAuthService := mock_service.NewMockAuthService(ctrl)

	userDataId := int64(123)
	data := []byte("example data")
	ctx := context.TODO()

	//expectedResponse := []byte("Data updated")
	expectedError := errors.New("error message")

	// Expected calls
	mockAuthService.EXPECT().GetJWTToken().Return(";ll;")
	mockEncrypter.EXPECT().Encrypt(data).Return(data, nil)
	mockDataInterface.EXPECT().PostUpdateData(ctx, userDataId, data).Return(nil, expectedError)
	//mockStorageData.EXPECT().UpdateData(ctx, userDataId, data, gomock.Any(), gomock.Any()).Return(expectedError)

	useCase := &UseCase{
		DataInterface: mockDataInterface,
		Encrypter:     mockEncrypter,
		StorageData:   mockStorageData,
		AuthService:   mockAuthService,
	}

	// Call the function
	resp, err := useCase.UpdateData(ctx, userDataId, data)

	// Assertions for successful response
	require.Error(t, err)
	require.Equal(t, err.Error(), expectedError.Error())
	require.Nil(t, resp)
	// Assertions for error response
}

// success test
func TestUpdateDataSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDataInterface := mock_service.NewMockDataInterface(ctrl)
	mockEncrypter := mock_service.NewMockEncrypter(ctrl)
	mockStorageData := mock_service.NewMockStorageData(ctrl)
	mockAuthService := mock_service.NewMockAuthService(ctrl)

	userDataId := int64(123)
	data := []byte("example data")
	ctx := context.TODO()

	expectedResponse := []byte("Data updated")

	mockAuthService.EXPECT().GetJWTToken().Return(";ll;")
	mockEncrypter.EXPECT().Encrypt(data).Return(data, nil)
	mockDataInterface.EXPECT().PostUpdateData(ctx, userDataId, data).Return(&server.RespData{}, nil)
	mockStorageData.EXPECT().UpdateData(ctx, userDataId, data, gomock.Any(), gomock.Any()).Return(nil)

	useCase := &UseCase{
		DataInterface: mockDataInterface,
		Encrypter:     mockEncrypter,
		StorageData:   mockStorageData,
		AuthService:   mockAuthService,
	}

	resp, err := useCase.UpdateData(ctx, userDataId, data)

	require.NoError(t, err)
	require.Equal(t, resp, expectedResponse)

}

func TestEncryptData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEncrypter := mock_service.NewMockEncrypter(ctrl)
	mockEncrypter.EXPECT().Encrypt(gomock.Any()).Return(nil, fmt.Errorf("error"))

	useCase := &UseCase{
		Encrypter: mockEncrypter,
	}

	err := useCase.encryptData(&server.ReqData{})

	require.Error(t, err)

}
func TestEncryptData2(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	e := "encrypt"
	d := "decrypt"

	r := server.ReqData{Data: []byte(d)}

	mockEncrypter := mock_service.NewMockEncrypter(ctrl)
	mockEncrypter.EXPECT().Encrypt(gomock.Any()).Return([]byte(e), nil)

	useCase := &UseCase{
		Encrypter: mockEncrypter,
	}

	err := useCase.encryptData(&r)

	require.Equal(t, e, string(r.Data))
	require.NoError(t, err)

}
func TestDecryptData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEncrypter := mock_service.NewMockEncrypter(ctrl)
	mockEncrypter.EXPECT().Decrypt(gomock.Any()).Return(nil, fmt.Errorf("error"))

	useCase := &UseCase{
		Encrypter: mockEncrypter,
	}

	err := useCase.decryptData(&server.ReqData{})

	require.Error(t, err)

}
func TestDecryptData2(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	e := "encrypt"
	d := "decrypt"

	r := server.ReqData{Data: []byte(d)}

	mockEncrypter := mock_service.NewMockEncrypter(ctrl)
	mockEncrypter.EXPECT().Decrypt(gomock.Any()).Return([]byte(e), nil)

	useCase := &UseCase{
		Encrypter: mockEncrypter,
	}

	err := useCase.decryptData(&r)

	require.Equal(t, e, string(r.Data))
	require.NoError(t, err)

}

func TestUseCase_encryptData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDataInterface := mock_service.NewMockDataInterface(ctrl)
	mockEncrypter := mock_service.NewMockEncrypter(ctrl)
	mockStorageData := mock_service.NewMockStorageData(ctrl)
	mockAuthService := mock_service.NewMockAuthService(ctrl)

	//mockAuthService.EXPECT().GetJWTToken().Return(";ll;")
	mockEncrypter.EXPECT().Encrypt(gomock.Any()).Return(nil, fmt.Errorf("error encrypting file"))
	useCase := &UseCase{
		DataInterface: mockDataInterface,
		Encrypter:     mockEncrypter,
		StorageData:   mockStorageData,
		AuthService:   mockAuthService,
	}
	err := useCase.encryptData(&server.ReqData{})
	require.Error(t, err)
	require.Equal(t, err.Error(), "error encrypting file")
}
func TestUseCase_decryptData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDataInterface := mock_service.NewMockDataInterface(ctrl)
	mockEncrypter := mock_service.NewMockEncrypter(ctrl)
	mockStorageData := mock_service.NewMockStorageData(ctrl)
	mockAuthService := mock_service.NewMockAuthService(ctrl)

	//mockAuthService.EXPECT().GetJWTToken().Return(";ll;")
	mockEncrypter.EXPECT().Decrypt(gomock.Any()).Return(nil, fmt.Errorf("error encrypting file"))
	useCase := &UseCase{
		DataInterface: mockDataInterface,
		Encrypter:     mockEncrypter,
		StorageData:   mockStorageData,
		AuthService:   mockAuthService,
	}
	err := useCase.decryptData(&server.ReqData{})
	require.Error(t, err)
	require.Equal(t, err.Error(), "error encrypting file")
}

func TestUseCase_GetDataFromAgentStorage(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDataInterface := mock_service.NewMockDataInterface(ctrl)
	mockEncrypter := mock_service.NewMockEncrypter(ctrl)
	mockStorageData := mock_service.NewMockStorageData(ctrl)
	mockAuthService := mock_service.NewMockAuthService(ctrl)

	mockStorageData.EXPECT().GetData(gomock.Any(), gomock.Any()).Return(nil, nil, fmt.Errorf("error get Data"))

	useCase := &UseCase{
		DataInterface: mockDataInterface,
		Encrypter:     mockEncrypter,
		StorageData:   mockStorageData,
		AuthService:   mockAuthService,
	}
	_, err := useCase.GetDataFromAgentStorage(context.TODO(), int64(1))
	require.Error(t, err)
	require.Equal(t, err.Error(), "error get Data")
}
func TestUseCase_GetDataFromAgentStorage2(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDataInterface := mock_service.NewMockDataInterface(ctrl)
	mockEncrypter := mock_service.NewMockEncrypter(ctrl)
	mockStorageData := mock_service.NewMockStorageData(ctrl)
	mockAuthService := mock_service.NewMockAuthService(ctrl)

	mockStorageData.EXPECT().GetData(gomock.Any(), gomock.Any()).Return(&store.UsersData{}, &store.DataFile{}, nil)
	mockEncrypter.EXPECT().Decrypt(gomock.Any()).Return(nil, fmt.Errorf("error get crypted data"))
	useCase := &UseCase{
		DataInterface: mockDataInterface,
		Encrypter:     mockEncrypter,
		StorageData:   mockStorageData,
		AuthService:   mockAuthService,
	}
	_, err := useCase.GetDataFromAgentStorage(context.TODO(), int64(1))
	require.Error(t, err)
	require.Equal(t, err.Error(), "error get crypted data")
}
func TestUseCase_GetDataFromAgentStorage3(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDataInterface := mock_service.NewMockDataInterface(ctrl)
	mockEncrypter := mock_service.NewMockEncrypter(ctrl)
	mockStorageData := mock_service.NewMockStorageData(ctrl)
	mockAuthService := mock_service.NewMockAuthService(ctrl)

	mockStorageData.EXPECT().GetData(gomock.Any(), gomock.Any()).Return(&store.UsersData{
		DataType: sqllite.TypeBinaryFile,
	}, &store.DataFile{}, nil)

	mockEncrypter.EXPECT().Decrypt(gomock.Any()).Return(nil, nil)
	useCase := &UseCase{
		DataInterface: mockDataInterface,
		Encrypter:     mockEncrypter,
		StorageData:   mockStorageData,
		AuthService:   mockAuthService,
	}
	mockStorageData.EXPECT().GetMetaData(gomock.Any(), gomock.Any()).Return(&store.MetaData{}, fmt.Errorf("error get meta data"))
	_, err := useCase.GetDataFromAgentStorage(context.TODO(), int64(1))
	require.Error(t, err)
	require.Equal(t, err.Error(), "error get meta data")
}
func TestUseCase_GetDataFromAgentStorage4(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDataInterface := mock_service.NewMockDataInterface(ctrl)
	mockEncrypter := mock_service.NewMockEncrypter(ctrl)
	mockStorageData := mock_service.NewMockStorageData(ctrl)
	mockAuthService := mock_service.NewMockAuthService(ctrl)

	mockStorageData.EXPECT().GetData(gomock.Any(), gomock.Any()).Return(&store.UsersData{
		DataType: sqllite.TypeBinaryFile,
	}, &store.DataFile{}, nil)

	mockEncrypter.EXPECT().Decrypt(gomock.Any()).Return(nil, nil)
	useCase := &UseCase{
		DataInterface: mockDataInterface,
		Encrypter:     mockEncrypter,
		StorageData:   mockStorageData,
		AuthService:   mockAuthService,
	}
	mockStorageData.EXPECT().GetMetaData(gomock.Any(), gomock.Any()).Return(&store.MetaData{}, nil)

	_, err := useCase.GetDataFromAgentStorage(context.TODO(), int64(1))
	require.Error(t, err)
	require.Equal(t, err.Error(), "unexpected end of JSON input")
}
func TestUseCase_GetDataFromAgentStorage5(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDataInterface := mock_service.NewMockDataInterface(ctrl)
	mockEncrypter := mock_service.NewMockEncrypter(ctrl)
	mockStorageData := mock_service.NewMockStorageData(ctrl)
	mockAuthService := mock_service.NewMockAuthService(ctrl)

	mockStorageData.EXPECT().GetData(gomock.Any(), gomock.Any()).Return(&store.UsersData{
		DataType: sqllite.TypeBinaryFile,
	}, &store.DataFile{
		EncryptData: []byte("test"),
	}, nil)
	j := model.Data{
		DecryptData: []byte("test"),
	}
	msg, err := json.Marshal(j)
	require.NoError(t, err)

	mockEncrypter.EXPECT().Decrypt(gomock.Any()).Return(msg, nil)
	useCase := &UseCase{
		DataInterface: mockDataInterface,
		Encrypter:     mockEncrypter,
		StorageData:   mockStorageData,
		AuthService:   mockAuthService,
	}
	mockStorageData.EXPECT().GetMetaData(gomock.Any(), gomock.Any()).Return(&store.MetaData{}, nil)
	mockEncrypter.EXPECT().DecryptFile(gomock.Any(), gomock.Any()).Return(fmt.Errorf("error decrypt file"))

	_, err = useCase.GetDataFromAgentStorage(context.TODO(), int64(1))
	require.Error(t, err)
	require.Equal(t, err.Error(), "error decrypt file")
}
func TestUseCase_GetDataFromAgentStorage6(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDataInterface := mock_service.NewMockDataInterface(ctrl)
	mockEncrypter := mock_service.NewMockEncrypter(ctrl)
	mockStorageData := mock_service.NewMockStorageData(ctrl)
	mockAuthService := mock_service.NewMockAuthService(ctrl)

	mockStorageData.EXPECT().GetData(gomock.Any(), gomock.Any()).Return(&store.UsersData{
		DataType: sqllite.TypeCreditCardData,
	}, &store.DataFile{
		EncryptData: []byte("test"),
	}, nil)
	j := model.Data{
		DecryptData: []byte("test"),
	}
	msg, err := json.Marshal(j)
	require.NoError(t, err)

	mockEncrypter.EXPECT().Decrypt(gomock.Any()).Return(msg, nil)
	useCase := &UseCase{
		DataInterface: mockDataInterface,
		Encrypter:     mockEncrypter,
		StorageData:   mockStorageData,
		AuthService:   mockAuthService,
	}

	resp, err := useCase.GetDataFromAgentStorage(context.TODO(), int64(1))
	require.NoError(t, err)
	require.Equal(t, resp, []byte("Данные {\"decryptData\":\"dGVzdA==\"}"))
}

func TestNewUseCase(t *testing.T) {
	db := &sqlx.DB{}
	key := []byte("shortkey")
	certFile := "path/to/certFile"
	keyFile := "path/to/keyFile"
	serverDns := "example.com"

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	NewUseCase(db, key, certFile, keyFile, serverDns)
}
func TestNewUseCase_InvalidCertFiles(t *testing.T) {
	db := &sqlx.DB{}
	key := []byte("examplekey123456")
	certFile := "invalid/path/to/certFile"
	keyFile := "invalid/path/to/keyFile"
	serverDns := "example.com"

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	NewUseCase(db, key, certFile, keyFile, serverDns)
}
func TestNewUseCase2(t *testing.T) {
	db := &sqlx.DB{}
	key := []byte("examplekey123456")
	certFile := "path/to/certFile"
	keyFile := "path/to/keyFile"
	serverDns := "example.com"

	a := func() (string, string, error) {
		// Генерация закрытого ключа
		privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			return "", "", err
		}

		// Генерация сертификата
		template := x509.Certificate{
			SerialNumber: big.NewInt(1),
			Subject: pkix.Name{
				Country:      []string{"US"},
				Organization: []string{"Example"},
			},
			NotBefore:   time.Now(),
			NotAfter:    time.Now().AddDate(1, 0, 0),
			KeyUsage:    x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
			ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		}

		certBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
		if err != nil {
			return "", "", err
		}

		// Запись сертификата в файл cert.pem
		certFile := "cert.pem"
		certOut, err := os.Create(certFile)
		if err != nil {
			return "", "", err
		}
		pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: certBytes})
		certOut.Close()

		// Запись закрытого ключа в файл key.pem
		keyFile := "key.pem"
		keyOut, err := os.Create(keyFile)
		if err != nil {
			return "", "", err
		}
		pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)})
		keyOut.Close()

		return certFile, keyFile, nil

	}
	certFile, keyFile, err := a()
	defer func() {
		os.Remove(certFile)
		os.Remove(keyFile)
	}()
	require.NoError(t, err)
	useCase := NewUseCase(db, key, certFile, keyFile, serverDns)
	require.NotNil(t, useCase)
}
func TestNewUseCase3(t *testing.T) {

	key := []byte("examplekey123456")
	certFile := "path/to/certFile"
	keyFile := "path/to/keyFile"
	serverDns := "example.com"

	a := func() (string, string, error) {
		// Генерация закрытого ключа
		privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			return "", "", err
		}

		// Генерация сертификата
		template := x509.Certificate{
			SerialNumber: big.NewInt(1),
			Subject: pkix.Name{
				Country:      []string{"US"},
				Organization: []string{"Example"},
			},
			NotBefore:   time.Now(),
			NotAfter:    time.Now().AddDate(1, 0, 0),
			KeyUsage:    x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
			ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		}

		certBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
		if err != nil {
			return "", "", err
		}

		// Запись сертификата в файл cert.pem
		certFile := "cert.pem"
		certOut, err := os.Create(certFile)
		if err != nil {
			return "", "", err
		}
		pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: certBytes})
		certOut.Close()

		// Запись закрытого ключа в файл key.pem
		keyFile := "key.pem"
		keyOut, err := os.Create(keyFile)
		if err != nil {
			return "", "", err
		}
		pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)})
		keyOut.Close()

		return certFile, keyFile, nil

	}
	certFile, keyFile, err := a()
	defer func() {
		os.Remove(certFile)
		os.Remove(keyFile)
	}()
	require.NoError(t, err)
	useCase := NewUseCase(nil, key, certFile, keyFile, serverDns)
	require.NotNil(t, useCase)
}

func TestCreateFileData_SuccessfullyEncryptsData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDataInterface := mock_service.NewMockDataInterface(ctrl)
	mockEncrypter := mock_service.NewMockEncrypter(ctrl)
	mockStorageData := mock_service.NewMockStorageData(ctrl)
	mockAuthService := mock_service.NewMockAuthService(ctrl)

	useCase := &UseCase{
		AuthService:   mockAuthService,
		DataInterface: mockDataInterface,
		StorageData:   mockStorageData,
		Encrypter:     mockEncrypter,
	}

	mockAuthService.EXPECT().GetJWTToken().Return("")

	mockStorageData.EXPECT().GetJWTToken(gomock.Any()).Return("", fmt.Errorf("error"))
	err := useCase.CreateFileData(context.TODO(), &server.ReqData{})
	assert.Error(t, err)

}

func TestCreateFileData_SuccessfullyEncryptsData3(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDataInterface := mock_service.NewMockDataInterface(ctrl)
	mockEncrypter := mock_service.NewMockEncrypter(ctrl)
	mockStorageData := mock_service.NewMockStorageData(ctrl)
	mockAuthService := mock_service.NewMockAuthService(ctrl)

	useCase := &UseCase{
		AuthService:   mockAuthService,
		DataInterface: mockDataInterface,
		StorageData:   mockStorageData,
		Encrypter:     mockEncrypter,
	}

	mockAuthService.EXPECT().GetJWTToken().Return("")

	mockStorageData.EXPECT().GetJWTToken(gomock.Any()).Return("", nil)
	err := useCase.CreateFileData(context.TODO(), &server.ReqData{})
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "jwt is empty")

}

func TestUseCase_CreateCreditCard(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDataInterface := mock_service.NewMockDataInterface(ctrl)
	mockEncrypter := mock_service.NewMockEncrypter(ctrl)
	mockStorageData := mock_service.NewMockStorageData(ctrl)
	mockAuthService := mock_service.NewMockAuthService(ctrl)

	useCase := &UseCase{
		AuthService:   mockAuthService,
		DataInterface: mockDataInterface,
		StorageData:   mockStorageData,
		Encrypter:     mockEncrypter,
	}

	mockAuthService.EXPECT().GetJWTToken().Return("")

	mockStorageData.EXPECT().GetJWTToken(gomock.Any()).Return("", nil)
	err := useCase.CreateCreditCard(context.TODO(), &server.ReqData{})
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "jwt is empty")
}
func TestUseCase_CreateCreditCard2(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDataInterface := mock_service.NewMockDataInterface(ctrl)
	mockEncrypter := mock_service.NewMockEncrypter(ctrl)
	mockStorageData := mock_service.NewMockStorageData(ctrl)
	mockAuthService := mock_service.NewMockAuthService(ctrl)

	useCase := &UseCase{
		AuthService:   mockAuthService,
		DataInterface: mockDataInterface,
		StorageData:   mockStorageData,
		Encrypter:     mockEncrypter,
	}

	mockAuthService.EXPECT().GetJWTToken().Return("")

	mockStorageData.EXPECT().GetJWTToken(gomock.Any()).Return("", fmt.Errorf("error"))
	err := useCase.CreateCreditCard(context.TODO(), &server.ReqData{})
	assert.Error(t, err)
}
func TestUseCase_CreateCreditCard3(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDataInterface := mock_service.NewMockDataInterface(ctrl)
	mockEncrypter := mock_service.NewMockEncrypter(ctrl)
	mockStorageData := mock_service.NewMockStorageData(ctrl)
	mockAuthService := mock_service.NewMockAuthService(ctrl)

	useCase := &UseCase{
		AuthService:   mockAuthService,
		DataInterface: mockDataInterface,
		StorageData:   mockStorageData,
		Encrypter:     mockEncrypter,
	}

	mockAuthService.EXPECT().GetJWTToken().Return("tplkemn")
	mockEncrypter.EXPECT().Encrypt(gomock.Any()).Return([]byte("sss"), fmt.Errorf("error"))

	err := useCase.CreateCreditCard(context.TODO(), &server.ReqData{})
	assert.Error(t, err)
}

func TestUseCase_CreateFileData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDataInterface := mock_service.NewMockDataInterface(ctrl)
	mockEncrypter := mock_service.NewMockEncrypter(ctrl)
	mockStorageData := mock_service.NewMockStorageData(ctrl)
	mockAuthService := mock_service.NewMockAuthService(ctrl)

	useCase := &UseCase{
		AuthService:   mockAuthService,
		DataInterface: mockDataInterface,
		StorageData:   mockStorageData,
		Encrypter:     mockEncrypter,
	}

	mockAuthService.EXPECT().GetJWTToken().Return("")

	mockStorageData.EXPECT().GetJWTToken(gomock.Any()).Return("", fmt.Errorf("error"))
	err := useCase.CreateFileData(context.TODO(), &server.ReqData{})
	assert.Error(t, err)
}
func TestUseCase_CreateFileData2(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDataInterface := mock_service.NewMockDataInterface(ctrl)
	mockEncrypter := mock_service.NewMockEncrypter(ctrl)
	mockStorageData := mock_service.NewMockStorageData(ctrl)
	mockAuthService := mock_service.NewMockAuthService(ctrl)

	useCase := &UseCase{
		AuthService:   mockAuthService,
		DataInterface: mockDataInterface,
		StorageData:   mockStorageData,
		Encrypter:     mockEncrypter,
	}

	mockAuthService.EXPECT().GetJWTToken().Return("")

	mockStorageData.EXPECT().GetJWTToken(gomock.Any()).Return("", nil)
	err := useCase.CreateFileData(context.TODO(), &server.ReqData{})
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "jwt is empty")
}

func TestUseCase_GetData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDataInterface := mock_service.NewMockDataInterface(ctrl)
	mockEncrypter := mock_service.NewMockEncrypter(ctrl)
	mockStorageData := mock_service.NewMockStorageData(ctrl)
	mockAuthService := mock_service.NewMockAuthService(ctrl)

	useCase := &UseCase{
		AuthService:   mockAuthService,
		DataInterface: mockDataInterface,
		StorageData:   mockStorageData,
		Encrypter:     mockEncrypter,
	}

	mockDataInterface.EXPECT().Ping(gomock.Any()).Return(nil)

	mockAuthService.EXPECT().GetJWTToken().Return("")

	mockStorageData.EXPECT().GetJWTToken(gomock.Any()).Return("", nil)

	_, err := useCase.GetData(context.TODO(), int64(1))
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "jwt is empty")
}
func TestUseCase_GetData2(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDataInterface := mock_service.NewMockDataInterface(ctrl)
	mockEncrypter := mock_service.NewMockEncrypter(ctrl)
	mockStorageData := mock_service.NewMockStorageData(ctrl)
	mockAuthService := mock_service.NewMockAuthService(ctrl)

	useCase := &UseCase{
		AuthService:   mockAuthService,
		DataInterface: mockDataInterface,
		StorageData:   mockStorageData,
		Encrypter:     mockEncrypter,
	}

	mockDataInterface.EXPECT().Ping(gomock.Any()).Return(nil)

	mockAuthService.EXPECT().GetJWTToken().Return("")

	mockStorageData.EXPECT().GetJWTToken(gomock.Any()).Return("", fmt.Errorf("error"))

	_, err := useCase.GetData(context.TODO(), int64(1))
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "error")
}
func TestUseCase_GetData3(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDataInterface := mock_service.NewMockDataInterface(ctrl)
	mockEncrypter := mock_service.NewMockEncrypter(ctrl)
	mockStorageData := mock_service.NewMockStorageData(ctrl)
	mockAuthService := mock_service.NewMockAuthService(ctrl)

	useCase := &UseCase{
		AuthService:   mockAuthService,
		DataInterface: mockDataInterface,
		StorageData:   mockStorageData,
		Encrypter:     mockEncrypter,
	}

	mockDataInterface.EXPECT().Ping(gomock.Any()).Return(nil)

	mockAuthService.EXPECT().GetJWTToken().Return("kk")

	mockStorageData.EXPECT().GetInfoData(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error"))

	_, err := useCase.GetData(context.TODO(), int64(1))
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "error")
}
func TestUseCase_GetData4(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDataInterface := mock_service.NewMockDataInterface(ctrl)
	mockEncrypter := mock_service.NewMockEncrypter(ctrl)
	mockStorageData := mock_service.NewMockStorageData(ctrl)
	mockAuthService := mock_service.NewMockAuthService(ctrl)

	useCase := &UseCase{
		AuthService:   mockAuthService,
		DataInterface: mockDataInterface,
		StorageData:   mockStorageData,
		Encrypter:     mockEncrypter,
	}

	mockDataInterface.EXPECT().Ping(gomock.Any()).Return(nil)

	mockAuthService.EXPECT().GetJWTToken().Return("kk").AnyTimes()

	mockStorageData.EXPECT().GetInfoData(gomock.Any(), gomock.Any()).Return(&store.UsersData{
		UserDataId:  0,
		UserId:      0,
		DataId:      0,
		DataType:    0,
		Name:        "",
		Description: "",
		Hash:        "",
		CreatedAt:   nil,
		UpdateAt:    nil,
		IsDeleted:   false,
	}, nil)

	mockDataInterface.EXPECT().CheckUpdate(gomock.Any(), gomock.Any(), gomock.Any()).Return(false, fmt.Errorf("error"))

	_, err := useCase.GetData(context.TODO(), int64(1))
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "error")
}
func TestUseCase_GetData5(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDataInterface := mock_service.NewMockDataInterface(ctrl)
	mockEncrypter := mock_service.NewMockEncrypter(ctrl)
	mockStorageData := mock_service.NewMockStorageData(ctrl)
	mockAuthService := mock_service.NewMockAuthService(ctrl)

	useCase := &UseCase{
		AuthService:   mockAuthService,
		DataInterface: mockDataInterface,
		StorageData:   mockStorageData,
		Encrypter:     mockEncrypter,
	}

	mockDataInterface.EXPECT().Ping(gomock.Any()).Return(nil)

	mockAuthService.EXPECT().GetJWTToken().Return("kk").AnyTimes()

	mockStorageData.EXPECT().GetInfoData(gomock.Any(), gomock.Any()).Return(&store.UsersData{
		UserDataId:  0,
		UserId:      0,
		DataId:      0,
		DataType:    0,
		Name:        "",
		Description: "",
		Hash:        "",
		CreatedAt:   nil,
		UpdateAt:    nil,
		IsDeleted:   false,
	}, nil)

	//false/ok
	mockDataInterface.EXPECT().CheckUpdate(gomock.Any(), gomock.Any(), gomock.Any()).Return(false, nil)

	mockStorageData.EXPECT().GetData(gomock.Any(), gomock.Any()).Return(nil, nil, fmt.Errorf("error"))
	_, err := useCase.GetData(context.TODO(), int64(1))
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "error")
}

func TestUseCase_GetListData1(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDataInterface := mock_service.NewMockDataInterface(ctrl)
	mockEncrypter := mock_service.NewMockEncrypter(ctrl)
	mockStorageData := mock_service.NewMockStorageData(ctrl)
	mockAuthService := mock_service.NewMockAuthService(ctrl)

	useCase := &UseCase{
		AuthService:   mockAuthService,
		DataInterface: mockDataInterface,
		StorageData:   mockStorageData,
		Encrypter:     mockEncrypter,
	}

	mockAuthService.EXPECT().GetJWTToken().Return("")

	mockStorageData.EXPECT().GetJWTToken(gomock.Any()).Return("", fmt.Errorf("error"))

	_, err := useCase.GetListData(context.TODO())
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "error")
}
func TestUseCase_GetListData2(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDataInterface := mock_service.NewMockDataInterface(ctrl)
	mockEncrypter := mock_service.NewMockEncrypter(ctrl)
	mockStorageData := mock_service.NewMockStorageData(ctrl)
	mockAuthService := mock_service.NewMockAuthService(ctrl)

	useCase := &UseCase{
		AuthService:   mockAuthService,
		DataInterface: mockDataInterface,
		StorageData:   mockStorageData,
		Encrypter:     mockEncrypter,
	}

	mockAuthService.EXPECT().GetJWTToken().Return("")

	mockStorageData.EXPECT().GetJWTToken(gomock.Any()).Return("", nil)

	_, err := useCase.GetListData(context.TODO())
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "jwt is empty")
}
func TestUseCase_GetListData3(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDataInterface := mock_service.NewMockDataInterface(ctrl)
	mockEncrypter := mock_service.NewMockEncrypter(ctrl)
	mockStorageData := mock_service.NewMockStorageData(ctrl)
	mockAuthService := mock_service.NewMockAuthService(ctrl)

	useCase := &UseCase{
		AuthService:   mockAuthService,
		DataInterface: mockDataInterface,
		StorageData:   mockStorageData,
		Encrypter:     mockEncrypter,
	}

	mockAuthService.EXPECT().GetJWTToken().Return("sss")

	mockDataInterface.EXPECT().GetListData(gomock.Any()).Return([]byte("data"), nil)

	resp, err := useCase.GetListData(context.TODO())
	assert.NoError(t, err)
	assert.Equal(t, string(resp), "data")
}
func TestUseCase_GetListData4(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDataInterface := mock_service.NewMockDataInterface(ctrl)
	mockEncrypter := mock_service.NewMockEncrypter(ctrl)
	mockStorageData := mock_service.NewMockStorageData(ctrl)
	mockAuthService := mock_service.NewMockAuthService(ctrl)

	useCase := &UseCase{
		AuthService:   mockAuthService,
		DataInterface: mockDataInterface,
		StorageData:   mockStorageData,
		Encrypter:     mockEncrypter,
	}

	mockAuthService.EXPECT().GetJWTToken().Return("sss")

	mockDataInterface.EXPECT().GetListData(gomock.Any()).Return([]byte("data"), fmt.Errorf("error"))

	resp, err := useCase.GetListData(context.TODO())
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "error")
	assert.Nil(t, resp)
}
