package service

import (
	"GophKeeper/internal/Agent/server"
	mock_service "GophKeeper/internal/Agent/service/mocks"
	"GophKeeper/pkg/store"
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"testing"
)

func TestUseCase_CreateBinaryFile(t *testing.T) {
	ctx := context.TODO()

	path := path.Join("/tmp", "test.txt")
	name := "test"
	description := "test"
	ch := make(chan string)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDataInterface := mock_service.NewMockDataInterface(ctrl)
	mockEncrypter := mock_service.NewMockEncrypter(ctrl)
	mockStorageData := mock_service.NewMockStorageData(ctrl)
	mockAuthService := mock_service.NewMockAuthService(ctrl)

	mockAuthService.EXPECT().GetJWTToken().Return(";ll;")
	mockEncrypter.EXPECT().EncryptFile(gomock.Any(), gomock.Any()).Return(fmt.Errorf("error encrypting file"))
	useCase := &UseCase{
		DataInterface: mockDataInterface,
		Encrypter:     mockEncrypter,
		StorageData:   mockStorageData,
		AuthService:   mockAuthService,
	}
	err := useCase.CreateBinaryFile(ctx, path, name, description, ch)
	require.Error(t, err)
	require.Equal(t, err.Error(), "error encrypting file")
}

func TestUseCase_CreateBinaryFile7(t *testing.T) {
	ctx := context.TODO()
	path := path.Join("/tmp", "test.txt")
	name := "test"
	description := "test"
	ch := make(chan string)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDataInterface := mock_service.NewMockDataInterface(ctrl)
	mockEncrypter := mock_service.NewMockEncrypter(ctrl)
	mockStorageData := mock_service.NewMockStorageData(ctrl)
	mockAuthService := mock_service.NewMockAuthService(ctrl)

	mockAuthService.EXPECT().GetJWTToken().Return("")
	mockStorageData.EXPECT().GetJWTToken(gomock.Any()).Return("", nil)

	//mockEncrypter.EXPECT().EncryptFile(gomock.Any(), gomock.Any()).Return(fmt.Errorf("error encrypting file"))
	useCase := &UseCase{
		DataInterface: mockDataInterface,
		Encrypter:     mockEncrypter,
		StorageData:   mockStorageData,
		AuthService:   mockAuthService,
	}
	err := useCase.CreateBinaryFile(ctx, path, name, description, ch)
	require.Error(t, err)
	require.Equal(t, err.Error(), "jwt is empty")
}
func TestUseCase_CreateBinaryFile8(t *testing.T) {
	ctx := context.TODO()
	path := path.Join("/tmp", "test.txt")
	name := "test"
	description := "test"
	ch := make(chan string)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDataInterface := mock_service.NewMockDataInterface(ctrl)
	mockEncrypter := mock_service.NewMockEncrypter(ctrl)
	mockStorageData := mock_service.NewMockStorageData(ctrl)
	mockAuthService := mock_service.NewMockAuthService(ctrl)

	mockAuthService.EXPECT().GetJWTToken().Return("")
	mockStorageData.EXPECT().GetJWTToken(gomock.Any()).Return("", fmt.Errorf("error get jwt token"))

	//mockEncrypter.EXPECT().EncryptFile(gomock.Any(), gomock.Any()).Return(fmt.Errorf("error encrypting file"))
	useCase := &UseCase{
		DataInterface: mockDataInterface,
		Encrypter:     mockEncrypter,
		StorageData:   mockStorageData,
		AuthService:   mockAuthService,
	}
	err := useCase.CreateBinaryFile(ctx, path, name, description, ch)
	require.Error(t, err)
	require.Equal(t, err.Error(), "error get jwt token")
}
func TestUseCase_CreateBinaryFile2(t *testing.T) {
	ctx := context.TODO()
	path := path.Join("/tmp", "test.txt")
	name := "test"
	description := "test"
	ch := make(chan string)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDataInterface := mock_service.NewMockDataInterface(ctrl)
	mockEncrypter := mock_service.NewMockEncrypter(ctrl)
	mockStorageData := mock_service.NewMockStorageData(ctrl)
	mockAuthService := mock_service.NewMockAuthService(ctrl)

	mockAuthService.EXPECT().GetJWTToken().Return(";ll;")
	mockEncrypter.EXPECT().EncryptFile(gomock.Any(), gomock.Any()).Return(nil)
	mockEncrypter.EXPECT().Encrypt(gomock.Any()).Return(nil, fmt.Errorf("error encrypting"))
	useCase := &UseCase{
		DataInterface: mockDataInterface,
		Encrypter:     mockEncrypter,
		StorageData:   mockStorageData,
		AuthService:   mockAuthService,
	}
	err := useCase.CreateBinaryFile(ctx, path, name, description, ch)
	require.Error(t, err)
	require.Equal(t, err.Error(), "error encrypting")
}
func TestUseCase_CreateBinaryFile3(t *testing.T) {
	ctx := context.TODO()
	path := path.Join("/tmp", "test.txt")
	name := "test"
	description := "test"
	ch := make(chan string)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDataInterface := mock_service.NewMockDataInterface(ctrl)
	mockEncrypter := mock_service.NewMockEncrypter(ctrl)
	mockStorageData := mock_service.NewMockStorageData(ctrl)
	mockAuthService := mock_service.NewMockAuthService(ctrl)

	mockAuthService.EXPECT().GetJWTToken().Return(";ll;")
	mockEncrypter.EXPECT().EncryptFile(gomock.Any(), gomock.Any()).Return(nil)
	mockEncrypter.EXPECT().Encrypt(gomock.Any()).Return([]byte("sss"), nil)

	useCase := &UseCase{
		DataInterface: mockDataInterface,
		Encrypter:     mockEncrypter,
		StorageData:   mockStorageData,
		AuthService:   mockAuthService,
	}
	err := useCase.CreateBinaryFile(ctx, path, name, description, ch)
	require.Error(t, err)
	require.Equal(t, err.Error(), "Error open file")
}

// сохранение в локальный репозиторий
func TestUseCase_CreateBinaryFile5(t *testing.T) {
	// Create a temporary source file for testing
	tempFile, err := ioutil.TempFile("", "source.txt")
	require.NoError(t, err)
	pathhh := tempFile.Name()
	_, err = tempFile.Write([]byte("Hello, World!"))

	//require.NoError(t, err)
	defer os.Remove(tempFile.Name())

	ctx := context.TODO()
	name := "test"
	description := "test"
	ch := make(chan string)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDataInterface := mock_service.NewMockDataInterface(ctrl)
	mockEncrypter := mock_service.NewMockEncrypter(ctrl)
	mockStorageData := mock_service.NewMockStorageData(ctrl)
	mockAuthService := mock_service.NewMockAuthService(ctrl)

	mockAuthService.EXPECT().GetJWTToken().Return(";ll;")
	mockEncrypter.EXPECT().EncryptFile(gomock.Any(), gomock.Any()).Return(nil).Do(func(a, b string) {
		fullPath, err := filepath.Abs(b)

		_ = fullPath
		d := path.Dir(b)
		_ = d
		err = os.MkdirAll(d, os.ModePerm)
		require.NoError(t, err)

		err = ioutil.WriteFile(b, []byte("Hello, World!"), 0644)
		require.NoError(t, err)

	})
	mockEncrypter.EXPECT().Encrypt(gomock.Any()).Return([]byte("sss"), nil)
	mockDataInterface.EXPECT().PostCrateFileStartChunks(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("", nil, fmt.Errorf("error creating file chunks"))

	useCase := &UseCase{
		DataInterface: mockDataInterface,
		Encrypter:     mockEncrypter,
		StorageData:   mockStorageData,
		AuthService:   mockAuthService,
	}
	err = useCase.CreateBinaryFile(ctx, pathhh, name, description, ch)
	require.Error(t, err)
	require.Equal(t, err.Error(), "error creating file chunks")
}
func TestUseCase_CreateBinaryFile4(t *testing.T) {
	// Create a temporary source file for testing
	tempFile, err := ioutil.TempFile("", "source.txt")
	require.NoError(t, err)
	pathhh := tempFile.Name()
	_, err = tempFile.Write([]byte("Hello, World!"))

	//require.NoError(t, err)
	defer os.Remove(tempFile.Name())

	ctx := context.TODO()
	name := "test"
	description := "test"
	ch := make(chan string)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDataInterface := mock_service.NewMockDataInterface(ctrl)
	mockEncrypter := mock_service.NewMockEncrypter(ctrl)
	mockStorageData := mock_service.NewMockStorageData(ctrl)
	mockAuthService := mock_service.NewMockAuthService(ctrl)

	mockAuthService.EXPECT().GetJWTToken().Return(";ll;")
	mockEncrypter.EXPECT().EncryptFile(gomock.Any(), gomock.Any()).Return(nil).Do(func(a, b string) {
		fullPath, err := filepath.Abs(b)

		_ = fullPath
		d := path.Dir(b)
		_ = d
		err = os.MkdirAll(d, os.ModePerm)
		require.NoError(t, err)

		err = ioutil.WriteFile(b, []byte("Hello, World!"), 0644)
		require.NoError(t, err)

	})
	mockEncrypter.EXPECT().Encrypt(gomock.Any()).Return([]byte("sss"), nil)
	mockDataInterface.EXPECT().PostCrateFileStartChunks(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("s", &server.RespData{
		UserDataId: 1,
		Hash:       "12312",
		CreatedAt:  nil,
		UpdateAt:   nil,
	}, nil)
	mockStorageData.EXPECT().CreateBinaryFile(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("error creating file local"))
	useCase := &UseCase{
		DataInterface: mockDataInterface,
		Encrypter:     mockEncrypter,
		StorageData:   mockStorageData,
		AuthService:   mockAuthService,
	}
	go func() {
		for {
			v, ok := <-ch
			if !ok {
				return
			}
			_ = v
		}
	}()
	err = useCase.CreateBinaryFile(ctx, pathhh, name, description, ch)
	require.Error(t, err)
	require.Equal(t, err.Error(), "error creating file local")
}
func TestUseCase_CreateBinaryFile6(t *testing.T) {
	// Create a temporary source file for testing
	tempFile, err := ioutil.TempFile("", "source.txt")
	require.NoError(t, err)
	pathhh := tempFile.Name()
	_, err = tempFile.Write([]byte("Hello, World!"))

	//require.NoError(t, err)
	defer os.Remove(tempFile.Name())

	ctx := context.TODO()
	name := "test"
	description := "test"
	ch := make(chan string)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDataInterface := mock_service.NewMockDataInterface(ctrl)
	mockEncrypter := mock_service.NewMockEncrypter(ctrl)
	mockStorageData := mock_service.NewMockStorageData(ctrl)
	mockAuthService := mock_service.NewMockAuthService(ctrl)

	mockAuthService.EXPECT().GetJWTToken().Return(";ll;")
	mockEncrypter.EXPECT().EncryptFile(gomock.Any(), gomock.Any()).Return(nil).Do(func(a, b string) {
		fullPath, err := filepath.Abs(b)

		_ = fullPath
		d := path.Dir(b)
		_ = d
		err = os.MkdirAll(d, os.ModePerm)
		require.NoError(t, err)

		err = ioutil.WriteFile(b, []byte("Hello, World!"), 0644)
		require.NoError(t, err)

	})
	mockEncrypter.EXPECT().Encrypt(gomock.Any()).Return([]byte("sss"), nil)
	mockDataInterface.EXPECT().PostCrateFileStartChunks(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("s", nil, nil)
	//	mockStorageData.EXPECT().CreateBinaryFile(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("error creating file local"))
	useCase := &UseCase{
		DataInterface: mockDataInterface,
		Encrypter:     mockEncrypter,
		StorageData:   mockStorageData,
		AuthService:   mockAuthService,
	}
	go func() {
		for {
			v, ok := <-ch
			if !ok {
				return
			}
			_ = v
		}
	}()
	err = useCase.CreateBinaryFile(ctx, pathhh, name, description, ch)
	require.Error(t, err)
	require.Equal(t, err.Error(), "resp is nil")
}

// ////////////////////////////
func TestUseCase_UpdateBinaryFile(t *testing.T) {
	ctx := context.TODO()
	path := path.Join("/tmp", "test.txt")
	userDataId := int64(1)

	ch := make(chan string)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDataInterface := mock_service.NewMockDataInterface(ctrl)
	mockEncrypter := mock_service.NewMockEncrypter(ctrl)
	mockStorageData := mock_service.NewMockStorageData(ctrl)
	mockAuthService := mock_service.NewMockAuthService(ctrl)

	mockAuthService.EXPECT().GetJWTToken().Return(";ll;")
	mockEncrypter.EXPECT().EncryptFile(gomock.Any(), gomock.Any()).Return(fmt.Errorf("error encrypting file"))
	useCase := &UseCase{
		DataInterface: mockDataInterface,
		Encrypter:     mockEncrypter,
		StorageData:   mockStorageData,
		AuthService:   mockAuthService,
	}
	err := useCase.UpdateBinaryFile(ctx, path, userDataId, ch)
	require.Error(t, err)
	require.Equal(t, err.Error(), "error encrypting file")
}
func TestUseCase_UpdateBinaryFile8(t *testing.T) {
	ctx := context.TODO()
	path := path.Join("/tmp", "test.txt")
	userDataId := int64(1)

	ch := make(chan string)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
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
	err := useCase.UpdateBinaryFile(ctx, path, userDataId, ch)
	require.Error(t, err)
	require.Equal(t, err.Error(), "error get jwt token")
}
func TestUseCase_UpdateBinaryFile7(t *testing.T) {
	ctx := context.TODO()
	path := path.Join("/tmp", "test.txt")
	userDataId := int64(1)

	ch := make(chan string)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
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
	err := useCase.UpdateBinaryFile(ctx, path, userDataId, ch)
	require.Error(t, err)
	require.Equal(t, err.Error(), "jwt is empty")
}

func TestUseCase_UpdateBinaryFile2(t *testing.T) {
	ctx := context.TODO()
	path := path.Join("/tmp", "test.txt")
	userDataId := int64(1)
	ch := make(chan string)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDataInterface := mock_service.NewMockDataInterface(ctrl)
	mockEncrypter := mock_service.NewMockEncrypter(ctrl)
	mockStorageData := mock_service.NewMockStorageData(ctrl)
	mockAuthService := mock_service.NewMockAuthService(ctrl)

	mockAuthService.EXPECT().GetJWTToken().Return(";ll;")
	mockEncrypter.EXPECT().EncryptFile(gomock.Any(), gomock.Any()).Return(nil)
	mockEncrypter.EXPECT().Encrypt(gomock.Any()).Return(nil, fmt.Errorf("error encrypting"))

	useCase := &UseCase{
		DataInterface: mockDataInterface,
		Encrypter:     mockEncrypter,
		StorageData:   mockStorageData,
		AuthService:   mockAuthService,
	}
	err := useCase.UpdateBinaryFile(ctx, path, userDataId, ch)
	require.Error(t, err)
	require.Equal(t, err.Error(), "error encrypting")
}
func TestUseCase_UpdateBinaryFile3(t *testing.T) {
	ctx := context.TODO()
	path := path.Join("/tmp", "test.txt")
	userDataId := int64(1)
	ch := make(chan string)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDataInterface := mock_service.NewMockDataInterface(ctrl)
	mockEncrypter := mock_service.NewMockEncrypter(ctrl)
	mockStorageData := mock_service.NewMockStorageData(ctrl)
	mockAuthService := mock_service.NewMockAuthService(ctrl)

	mockAuthService.EXPECT().GetJWTToken().Return(";ll;")
	mockEncrypter.EXPECT().EncryptFile(gomock.Any(), gomock.Any()).Return(nil)
	mockEncrypter.EXPECT().Encrypt(gomock.Any()).Return([]byte("sss"), nil)

	useCase := &UseCase{
		DataInterface: mockDataInterface,
		Encrypter:     mockEncrypter,
		StorageData:   mockStorageData,
		AuthService:   mockAuthService,
	}
	err := useCase.UpdateBinaryFile(ctx, path, userDataId, ch)
	require.Error(t, err)
	require.Equal(t, err.Error(), "Error open file")
}

// сохранение в локальный репозиторий
func TestUseCase_UpdateBinaryFile5(t *testing.T) {
	// Create a temporary source file for testing
	tempFile, err := ioutil.TempFile("", "source.txt")
	require.NoError(t, err)
	pathhh := tempFile.Name()
	_, err = tempFile.Write([]byte("Hello, World!"))

	//require.NoError(t, err)
	defer os.Remove(tempFile.Name())

	ctx := context.TODO()
	userDataId := int64(1)

	ch := make(chan string)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDataInterface := mock_service.NewMockDataInterface(ctrl)
	mockEncrypter := mock_service.NewMockEncrypter(ctrl)
	mockStorageData := mock_service.NewMockStorageData(ctrl)
	mockAuthService := mock_service.NewMockAuthService(ctrl)

	mockAuthService.EXPECT().GetJWTToken().Return(";ll;")
	mockEncrypter.EXPECT().EncryptFile(gomock.Any(), gomock.Any()).Return(nil).Do(func(a, b string) {
		fullPath, err := filepath.Abs(b)

		_ = fullPath
		d := path.Dir(b)
		_ = d
		err = os.MkdirAll(d, os.ModePerm)
		require.NoError(t, err)

		err = ioutil.WriteFile(b, []byte("Hello, World!"), 0644)
		require.NoError(t, err)

	})
	mockEncrypter.EXPECT().Encrypt(gomock.Any()).Return([]byte("sss"), nil)
	mockDataInterface.EXPECT().PostUpdateBinaryFile(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("", nil, fmt.Errorf("error creating file chunks"))

	useCase := &UseCase{
		DataInterface: mockDataInterface,
		Encrypter:     mockEncrypter,
		StorageData:   mockStorageData,
		AuthService:   mockAuthService,
	}
	err = useCase.UpdateBinaryFile(ctx, pathhh, userDataId, ch)
	require.Error(t, err)
	require.Equal(t, err.Error(), "error creating file chunks")
}
func TestUseCase_UpdateBinaryFile4(t *testing.T) {
	// Create a temporary source file for testing
	tempFile, err := ioutil.TempFile("", "source.txt")
	require.NoError(t, err)
	pathhh := tempFile.Name()
	_, err = tempFile.Write([]byte("Hello, World!"))

	//require.NoError(t, err)
	defer os.Remove(tempFile.Name())

	ctx := context.TODO()
	userDataId := int64(1)
	ch := make(chan string)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDataInterface := mock_service.NewMockDataInterface(ctrl)
	mockEncrypter := mock_service.NewMockEncrypter(ctrl)
	mockStorageData := mock_service.NewMockStorageData(ctrl)
	mockAuthService := mock_service.NewMockAuthService(ctrl)

	mockAuthService.EXPECT().GetJWTToken().Return(";ll;")
	mockEncrypter.EXPECT().EncryptFile(gomock.Any(), gomock.Any()).Return(nil).Do(func(a, b string) {
		fullPath, err := filepath.Abs(b)

		_ = fullPath
		d := path.Dir(b)
		_ = d
		err = os.MkdirAll(d, os.ModePerm)
		require.NoError(t, err)

		err = ioutil.WriteFile(b, []byte("Hello, World!"), 0644)
		require.NoError(t, err)

	})
	mockEncrypter.EXPECT().Encrypt(gomock.Any()).Return([]byte("sss"), nil)

	mockDataInterface.EXPECT().PostUpdateBinaryFile(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("1", &server.RespData{
		UserDataId: 1,
		Hash:       "12312",
		CreatedAt:  nil,
		UpdateAt:   nil,
	}, nil)

	mockStorageData.EXPECT().UpdateDataBinary(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("error creating file local"))
	useCase := &UseCase{
		DataInterface: mockDataInterface,
		Encrypter:     mockEncrypter,
		StorageData:   mockStorageData,
		AuthService:   mockAuthService,
	}
	go func() {
		for {
			v, ok := <-ch
			if !ok {
				return
			}
			_ = v
		}
	}()
	err = useCase.UpdateBinaryFile(ctx, pathhh, userDataId, ch)
	require.Error(t, err)
	require.Equal(t, err.Error(), "error creating file local")
}
func TestUseCase_UpdateBinaryFile6(t *testing.T) {
	// Create a temporary source file for testing
	tempFile, err := ioutil.TempFile("", "source.txt")
	require.NoError(t, err)
	pathhh := tempFile.Name()
	_, err = tempFile.Write([]byte("Hello, World!"))

	//require.NoError(t, err)
	defer os.Remove(tempFile.Name())

	ctx := context.TODO()
	userDataId := int64(1)
	ch := make(chan string)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDataInterface := mock_service.NewMockDataInterface(ctrl)
	mockEncrypter := mock_service.NewMockEncrypter(ctrl)
	mockStorageData := mock_service.NewMockStorageData(ctrl)
	mockAuthService := mock_service.NewMockAuthService(ctrl)

	mockAuthService.EXPECT().GetJWTToken().Return(";ll;")
	mockEncrypter.EXPECT().EncryptFile(gomock.Any(), gomock.Any()).Return(nil).Do(func(a, b string) {
		fullPath, err := filepath.Abs(b)

		_ = fullPath
		d := path.Dir(b)
		_ = d
		err = os.MkdirAll(d, os.ModePerm)
		require.NoError(t, err)

		err = ioutil.WriteFile(b, []byte("Hello, World!"), 0644)
		require.NoError(t, err)

	})
	mockEncrypter.EXPECT().Encrypt(gomock.Any()).Return([]byte("sss"), nil)

	mockDataInterface.EXPECT().PostUpdateBinaryFile(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("1", nil, nil)

	go func() {
		for {
			v, ok := <-ch
			if !ok {
				return
			}
			_ = v
		}
	}()
	useCase := &UseCase{
		DataInterface: mockDataInterface,
		Encrypter:     mockEncrypter,
		StorageData:   mockStorageData,
		AuthService:   mockAuthService,
	}
	err = useCase.UpdateBinaryFile(ctx, pathhh, userDataId, ch)
	require.Error(t, err)
	require.Equal(t, err.Error(), "resp is nil")
}

// ////////////////////////////
func TestUseCase_createEncryptedFile(t *testing.T) {
	path := path.Join("/tmp", "test.txt")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDataInterface := mock_service.NewMockDataInterface(ctrl)
	mockEncrypter := mock_service.NewMockEncrypter(ctrl)
	mockStorageData := mock_service.NewMockStorageData(ctrl)
	mockAuthService := mock_service.NewMockAuthService(ctrl)

	//mockAuthService.EXPECT().GetJWTToken().Return(";ll;")
	mockEncrypter.EXPECT().EncryptFile(gomock.Any(), gomock.Any()).Return(fmt.Errorf("error encrypting file"))
	useCase := &UseCase{
		DataInterface: mockDataInterface,
		Encrypter:     mockEncrypter,
		StorageData:   mockStorageData,
		AuthService:   mockAuthService,
	}
	_, err := useCase.createEncryptedFile(path)
	require.Error(t, err)
	require.Equal(t, err.Error(), "error encrypting file")
}
func TestUseCase_createEncryptedFile3(t *testing.T) {

	// Create a temporary source file for testing
	tempFile, err := ioutil.TempFile("", "source.txt")
	require.NoError(t, err)
	path := tempFile.Name()
	_, err = tempFile.Write([]byte("Hello, World!"))

	//require.NoError(t, err)
	defer os.Remove(tempFile.Name())
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDataInterface := mock_service.NewMockDataInterface(ctrl)
	mockEncrypter := mock_service.NewMockEncrypter(ctrl)
	mockStorageData := mock_service.NewMockStorageData(ctrl)
	mockAuthService := mock_service.NewMockAuthService(ctrl)

	//mockAuthService.EXPECT().GetJWTToken().Return(";ll;")
	mockEncrypter.EXPECT().EncryptFile(gomock.Any(), gomock.Any()).Return(fmt.Errorf("error encrypting file"))

	useCase := &UseCase{
		DataInterface: mockDataInterface,
		Encrypter:     mockEncrypter,
		StorageData:   mockStorageData,
		AuthService:   mockAuthService,
	}
	_, err = useCase.createEncryptedFile(path)
	require.Error(t, err)
	require.Equal(t, err.Error(), "error encrypting file")
}
func Test_copyFile(t *testing.T) {

	//src := "/tmp/source.txt"

	//newPath := "/tmp/new/"
	newPath := path.Join("/tmp", "new")
	newNameFile := "new_file.txt"

	// Create a temporary source file for testing
	tempFile, err := ioutil.TempFile("", "source.txt")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())

	// Write some data to the temporary source file
	testData := []byte("Hello, World!")
	_, err = tempFile.Write(testData)
	require.NoError(t, err)

	err = tempFile.Close()
	require.NoError(t, err)

	err = copyFile(tempFile.Name(), newPath, newNameFile)
	require.NoError(t, err)

	// Check if the new file exists in the specified path
	newFilePath := path.Join(newPath, newNameFile)
	_, err = os.Stat(newFilePath)
	require.NoError(t, err)

}

func TestUseCase_decryptFile(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDataInterface := mock_service.NewMockDataInterface(ctrl)
	mockEncrypter := mock_service.NewMockEncrypter(ctrl)
	mockStorageData := mock_service.NewMockStorageData(ctrl)
	mockAuthService := mock_service.NewMockAuthService(ctrl)

	//mockAuthService.EXPECT().GetJWTToken().Return(";ll;")
	mockEncrypter.EXPECT().DecryptFile(gomock.Any(), gomock.Any()).Return(fmt.Errorf("error encrypting file"))
	useCase := &UseCase{
		DataInterface: mockDataInterface,
		Encrypter:     mockEncrypter,
		StorageData:   mockStorageData,
		AuthService:   mockAuthService,
	}
	_, err := useCase.decryptFile(context.TODO(), &store.MetaData{
		FileName: "",
		PathSave: "",
		Size:     0,
	}, "sass")
	require.Error(t, err)
	require.Equal(t, err.Error(), "error encrypting file")
}

func TestUseCase_PrepareReqBinaryFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEncrypter := mock_service.NewMockEncrypter(ctrl)

	useCase := &UseCase{
		Encrypter: mockEncrypter,
	}

	originalFileName := "test.txt"
	name := "Test Name"
	description := "Test Description"

	//infoOriginalFile := server.DataFileInfo{OriginalFileName: originalFileName}
	//dataOriginalFile, _ := json.Marshal(infoOriginalFile)

	mockEncrypter.EXPECT().Encrypt(gomock.Any()).Return([]byte("test encrypt"), nil)

	//reqData := &server.ReqData{
	//	Name:        name,
	//	Description: description,
	//	Data:        dataOriginalFile,
	//}

	//reqDataJson, _ := json.Marshal(reqData)

	_, _, err := useCase.prepareReqBinaryFile(originalFileName, name, description)

	//assert.NotNil(t, resultReqData)
	//assert.NotNil(t, resultReqDataJson)
	assert.Nil(t, err)
}

// json error
func TestUseCase_PrepareReqBinaryFile2(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEncrypter := mock_service.NewMockEncrypter(ctrl)

	useCase := &UseCase{
		Encrypter: mockEncrypter,
	}

	originalFileName := "test.txt"
	name := "Test Name"
	description := "Test Description"

	//infoOriginalFile := server.DataFileInfo{OriginalFileName: originalFileName}
	//dataOriginalFile, _ := json.Marshal(infoOriginalFile)

	mockEncrypter.EXPECT().Encrypt(gomock.Any()).Return([]byte("test encrypt"), nil)

	//reqData := &server.ReqData{
	//	Name:        name,
	//	Description: description,
	//	Data:        dataOriginalFile,
	//}

	//reqDataJson, _ := json.Marshal(reqData)

	_, _, err := useCase.prepareReqBinaryFile(originalFileName, name, description)

	//assert.NotNil(t, resultReqData)
	//assert.NotNil(t, resultReqDataJson)
	assert.Nil(t, err)
}

func Test_copyFile1(t *testing.T) {
	src := "/tmp/source.txt"
	newPath := path.Join("/tmp", "new")
	newNameFile := "new_file.txt"

	// Create a temporary source file for testing
	err := copyFile(src, newPath, newNameFile)
	require.Error(t, err)
}
func Test_copyFile2(t *testing.T) {
	src := ""
	newPath := ""
	newNameFile := ""

	// Create a temporary source file for testing
	err := copyFile(src, newPath, newNameFile)
	require.Error(t, err)
}
func Test_copyFile3(t *testing.T) {

	newPath := path.Join("/tmp", "new")
	newNameFile := path.Join("new_file", "fff", "ss", "ss.txt")

	// Create a temporary source file for testing
	tempFile, err := ioutil.TempFile("", "source.txt")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())

	// Write some data to the temporary source file
	testData := []byte("Hello, World!")
	_, err = tempFile.Write(testData)
	require.NoError(t, err)

	err = tempFile.Close()
	require.NoError(t, err)

	err = copyFile(tempFile.Name(), newPath, newNameFile)
	require.Error(t, err)
}

func TestPrepareReqBinaryFile_SuccessfullyMarshalsDataFileInfo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := mock_service.NewMockAuthService(ctrl)
	mockDataInterface := mock_service.NewMockDataInterface(ctrl)
	mockStorageData := mock_service.NewMockStorageData(ctrl)
	mockEncrypter := mock_service.NewMockEncrypter(ctrl)

	uc := &UseCase{
		AuthService:   mockAuthService,
		DataInterface: mockDataInterface,
		StorageData:   mockStorageData,
		Encrypter:     mockEncrypter,
	}

	originalFileName := "testfile.txt"
	name := "Test File"
	description := "This is a test file"
	//var a server.ReqData
	or := server.DataFileInfo{
		OriginalFileName: originalFileName,
	}

	orr, err := json.Marshal(or)

	require.NoError(t, err)

	mockEncrypter.EXPECT().Encrypt(gomock.Any()).Return(orr, nil).Times(1)

	reqData, reqDataJson, err := uc.prepareReqBinaryFile(originalFileName, name, description)

	require.NoError(t, err)
	require.NotNil(t, reqData)
	require.NotNil(t, reqDataJson)

	var infoOriginalFile server.DataFileInfo
	err = json.Unmarshal(reqData.Data, &infoOriginalFile)
	require.NoError(t, err)
	require.Equal(t, originalFileName, infoOriginalFile.OriginalFileName)
}

func TestUseCase_saveLocalFile(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := mock_service.NewMockAuthService(ctrl)
	mockDataInterface := mock_service.NewMockDataInterface(ctrl)
	mockStorageData := mock_service.NewMockStorageData(ctrl)
	mockEncrypter := mock_service.NewMockEncrypter(ctrl)

	uc := &UseCase{
		AuthService:   mockAuthService,
		DataInterface: mockDataInterface,
		StorageData:   mockStorageData,
		Encrypter:     mockEncrypter,
	}
	err := uc.saveLocalFile(context.TODO(), nil, "", "", "", nil, nil)
	require.Error(t, err)
}
