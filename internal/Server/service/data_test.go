package service

import (
	"GophKeeper/pkg/store"
	"context"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestUseCase_CreateCredentials(t *testing.T) {
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

	storeData.EXPECT().CreateCredentials(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(&store.UsersData{}, nil)

	_, err := u.CreateCredentials(context.TODO(), 1, []byte("data"), "name", "description")
	assert.NoError(t, err)
}
func TestUseCase_CreateCredentialsBadNane(t *testing.T) {
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

	_, err := u.CreateCredentials(context.TODO(), 1, []byte("data"), "", "description")
	assert.Error(t, err)
}
func TestUseCase_CreateCredentialsBadCreate(t *testing.T) {
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

	storeData.EXPECT().CreateCredentials(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(&store.UsersData{}, fmt.Errorf("Error"))

	_, err := u.CreateCredentials(context.TODO(), 1, []byte("data"), "name", "description")
	assert.Error(t, err)
}

// CreateCreaditCard
func TestUseCase_CreateCreditCard(t *testing.T) {
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
	storeData.EXPECT().CreateCreditCard(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(&store.UsersData{}, nil)
	_, err := u.CreateCreditCard(context.TODO(), 1, []byte("data"), "name", "description")
	assert.NoError(t, err)
}
func TestUseCase_CreateCreditCardBadNane(t *testing.T) {
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

	_, err := u.CreateCreditCard(context.TODO(), 1, []byte("data"), "", "description")
	assert.Error(t, err)
}
func TestUseCase_CreateCreditCardBadCreate(t *testing.T) {
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

	storeData.EXPECT().CreateCreditCard(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("Error"))

	_, err := u.CreateCreditCard(context.TODO(), 1, []byte("data"), "name", "description")
	assert.Error(t, err)
}

// createPathIfNotExists
func Test_createPathIfNotExists(t *testing.T) {
	dirPath := "testDir"
	defer os.RemoveAll(dirPath)

	err := createPathIfNotExists(dirPath)

	assert.NoError(t, err)
	_, statErr := os.Stat(dirPath)
	assert.False(t, os.IsNotExist(statErr))
}
func TestCreatePathIfNotExists_DirectoryAlreadyExists(t *testing.T) {
	dirPath := "testDir/t"
	os.MkdirAll(dirPath, os.ModePerm)
	defer os.RemoveAll(dirPath)

	err := createPathIfNotExists(dirPath)

	assert.NoError(t, err)
}
func TestCreatePathIfNotExists_PermissionError(t *testing.T) {
	dirPath := "/root/testDir"

	err := createPathIfNotExists(dirPath)

	assert.Error(t, err)
	assert.True(t, os.IsPermission(err))
}

//moveFile

func Test_moveFile(t *testing.T) {
	src := "test_src.txt"
	dst := "test_dst.txt"

	// Create a dummy source file
	err := ioutil.WriteFile(src, []byte("dummy content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}

	defer os.Remove(src)
	defer os.Remove(dst)

	err = moveFile(src, dst)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if _, err := os.Stat(dst); os.IsNotExist(err) {
		t.Errorf("Expected destination file to exist, but it does not")
	}
}
func TestMoveFileInvalidDestinationPath(t *testing.T) {
	src := "test_src.txt"
	dst := "/invalid_path/test_dst.txt"

	// Create a dummy source file
	err := ioutil.WriteFile(src, []byte("dummy content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}

	defer os.Remove(src)

	err = moveFile(src, dst)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

// CreateFileChunks

func TestUseCase_CreateFileChunks(t *testing.T) {
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

	storeData.EXPECT().CreateFileDataChunks(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(&store.UsersData{}, nil)

	userId := int64(1)
	namef := uuid.NewString()
	tmpFile := &TmpFile{
		PathFileSave: "tmp/testfile.txt",
		Uuid:         namef,
		Size:         1024,
	}

	err := os.MkdirAll(path.Dir(tmpFile.PathFileSave), os.ModePerm)
	require.NoError(t, err)
	f, err := os.Create(tmpFile.PathFileSave)
	require.NoError(t, err)
	defer f.Close()
	defer os.Remove(tmpFile.PathFileSave)
	_, err = u.CreateFileChunks(context.TODO(), userId, tmpFile, "name", "description", []byte("data"))
	require.NoError(t, err)
}
func TestUseCase_CreateFileChunksMoveFile(t *testing.T) {
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

	userId := int64(1)
	namef := uuid.NewString()
	tmpFile := &TmpFile{
		PathFileSave: "",
		Uuid:         namef,
		Size:         1024,
	}

	_, err := u.CreateFileChunks(context.TODO(), userId, tmpFile, "name", "description", []byte("data"))
	require.Error(t, err)
}
func TestUseCase_CreateFileChunksCreateFileData(t *testing.T) {
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

	storeData.EXPECT().CreateFileDataChunks(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(&store.UsersData{}, errors.New("error"))

	userId := int64(1)
	namef := uuid.NewString()
	tmpFile := &TmpFile{
		PathFileSave: "tmp/testfile.txt",
		Uuid:         namef,
		Size:         1024,
	}

	err := os.MkdirAll(path.Dir(tmpFile.PathFileSave), os.ModePerm)
	require.NoError(t, err)
	f, err := os.Create(tmpFile.PathFileSave)
	require.NoError(t, err)
	defer f.Close()
	defer os.Remove(tmpFile.PathFileSave)
	_, err = u.CreateFileChunks(context.TODO(), userId, tmpFile, "name", "description", []byte("data"))
	require.Error(t, err)
}
func TestUseCase_CreateFileChunksErrCreateData(t *testing.T) {
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

	userId := int64(1)
	namef := uuid.NewString()
	tmpFile := &TmpFile{
		PathFileSave: "tmp/testfile.txt",
		Uuid:         namef,
		Size:         1024,
	}

	err := os.MkdirAll(path.Dir(tmpFile.PathFileSave), os.ModePerm)
	require.NoError(t, err)
	f, err := os.Create(tmpFile.PathFileSave)
	require.NoError(t, err)
	defer f.Close()
	defer os.Remove(tmpFile.PathFileSave)
	_, err = u.CreateFileChunks(context.TODO(), userId, tmpFile, "", "description", []byte("data"))
	require.Error(t, err)
}

// CreateFile
func TestUseCase_CreateFile(t *testing.T) {
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

	storeData.EXPECT().CreateFileData(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(&store.UsersData{}, nil)

	_, err := u.CreateFile(context.TODO(), 1, []byte("data"), "name", "description")
	assert.NoError(t, err)
}
func TestUseCase_CreateFileBadCreateData(t *testing.T) {
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

	_, err := u.CreateFile(context.TODO(), 1, []byte("data"), "", "description")
	assert.Error(t, err)
}
func TestUseCase_CreateFileBadCreateFileData(t *testing.T) {
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

	storeData.EXPECT().CreateFileData(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error"))

	_, err := u.CreateFile(context.TODO(), 1, []byte("data"), "name", "description")
	assert.Error(t, err)
}

// createData

func TestUseCase_createData(t *testing.T) {
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
	_, err := u.createData(context.TODO(), 1, []byte("data"), "name", "description")
	assert.NoError(t, err)
}
func TestUseCase_createDataDatanill(t *testing.T) {
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
	_, err := u.createData(context.TODO(), 1, nil, "name", "description")
	assert.Error(t, err)
}
func TestUseCase_createDataNameNull(t *testing.T) {
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
	_, err := u.createData(context.TODO(), 1, []byte("data"), "", "description")
	assert.Error(t, err)
}
func TestUseCase_createDataDescNull(t *testing.T) {
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
	_, err := u.createData(context.TODO(), 1, []byte("data"), "Name", "")
	assert.Error(t, err)
}
func TestUseCase_createDataUserIdNull(t *testing.T) {
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
	_, err := u.createData(context.TODO(), 0, []byte("data"), "Name", "")
	assert.Error(t, err)
}
