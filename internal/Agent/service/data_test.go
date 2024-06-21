package service

import (
	"GophKeeper/internal/Agent/server"
	mock_service "GophKeeper/internal/Agent/service/mocks"
	"context"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
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
		name string
		//	fields  fields
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

			ss := &Service{
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

	ss := &Service{
		AuthService:   nil,
		DataInterface: mockPost,
		StorageData:   nil,
		Encrypter:     nil,
		JWTToken:      "",
	}

	err := ss.Ping(context.TODO())
	require.NoError(t, err)
}

// server доступенн
func TestService_GetData(t *testing.T) {

}

// server недоступенн

func TestService_GetData1(t *testing.T) {

}
