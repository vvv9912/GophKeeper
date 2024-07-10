package service

import (
	"GophKeeper/pkg/store"
	"GophKeeper/pkg/store/postgresql"
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"runtime"
	"testing"
	"time"
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
	if runtime.GOOS != "linux" {
		t.Skip("Skipping test on non-Linux")
	}

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
	namef := uuid.NewString() + ".txt"
	tmpFile := &TmpFile{
		PathFileSave: path.Join("tmp", "testfile.txt"),
		Uuid:         namef,
		Size:         1024,
	}

	err := os.MkdirAll(path.Dir(tmpFile.PathFileSave), os.ModePerm)
	require.NoError(t, err)
	f, err := os.Create(tmpFile.PathFileSave)
	require.NoError(t, err)
	err = f.Close()
	require.NoError(t, err)

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
		PathFileSave: path.Join("tmp", "testfile.txt"),
		Uuid:         namef,
		Size:         1024,
	}

	err := os.MkdirAll(path.Dir(tmpFile.PathFileSave), os.ModePerm)
	require.NoError(t, err)
	f, err := os.Create(tmpFile.PathFileSave)
	require.NoError(t, err)
	err = f.Close()
	require.NoError(t, err)
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
		PathFileSave: path.Join("tmp", "testfile.txt"),
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

// ChangeData

func TestUseCase_ChangeData(t *testing.T) {
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

	storeData.EXPECT().ChangeData(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(true, nil)

	_, err := u.ChangeData(context.TODO(), 1, 1, time.Now())

	assert.NoError(t, err)
}
func TestUseCase_ChangeDataBad(t *testing.T) {
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

	storeData.EXPECT().ChangeData(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(false, nil)

	_, err := u.ChangeData(context.TODO(), 1, 1, time.Now())

	assert.NoError(t, err)
}

func TestUseCase_GetFileSize(t *testing.T) {
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

	storeData.EXPECT().GetFileSize(gomock.Any(), gomock.Any(), gomock.Any()).Return(int64(1), nil)

	_, err := u.GetFileSize(context.TODO(), 1, 1)
	assert.NoError(t, err)
}
func TestUseCase_GetFileSizeBadUserId(t *testing.T) {
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

	_, err := u.GetFileSize(context.TODO(), 0, 1)
	assert.Error(t, err)
}

func TestUseCase_GetFileSizeBadGetFileSize(t *testing.T) {
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

	storeData.EXPECT().GetFileSize(gomock.Any(), gomock.Any(), gomock.Any()).Return(int64(0), fmt.Errorf("error"))

	_, err := u.GetFileSize(context.TODO(), 1, 1)
	assert.Error(t, err)
}

// GetFileChunks

func TestUseCase_GetFileChunks(t *testing.T) {
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

	pathFile := "testfile.txt"
	expectedData := []byte("0123456789")
	// Create a test file
	err := os.WriteFile(pathFile, expectedData, 0644)
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}
	defer os.Remove(pathFile)

	contentRange := "bytes 0-10/10"

	req, err := http.NewRequest("GET", "http://localhost:8080/", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Range", contentRange)

	storeData.EXPECT().GetMetaData(gomock.Any(), gomock.Any(), gomock.Any()).Return(&store.MetaData{
		Size:     10,
		FileName: "testfile.txt",
		PathSave: "",
	}, nil)

	_, err = u.GetFileChunks(context.TODO(), 1, 1, req)
	assert.NoError(t, err)
}
func TestUseCase_GetFileChunksBadUserId(t *testing.T) {
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

	contentRange := "bytesh 0-10/10"

	req, err := http.NewRequest("GET", "http://localhost:8080/", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Range", contentRange)

	_, err = u.GetFileChunks(context.TODO(), 0, 1, req)
	assert.Error(t, err)
}
func TestUseCase_GetFileChunksBadParserContentRange(t *testing.T) {
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

	contentRange := "bytesh 0-10/10"

	req, err := http.NewRequest("GET", "http://localhost:8080/", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Range", contentRange)

	_, err = u.GetFileChunks(context.TODO(), 1, 1, req)
	assert.Error(t, err)
}
func TestUseCase_GetFileChunksErrorGetMera(t *testing.T) {
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

	pathFile := "testfile.txt"
	expectedData := []byte("0123456789")
	// Create a test file
	err := os.WriteFile(pathFile, expectedData, 0644)
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}
	defer os.Remove(pathFile)

	contentRange := "bytes 0-10/10"

	req, err := http.NewRequest("GET", "http://localhost:8080/", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Range", contentRange)

	storeData.EXPECT().GetMetaData(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error"))

	_, err = u.GetFileChunks(context.TODO(), 1, 1, req)
	assert.Error(t, err)
}

func TestUseCase_GetFileChunksBadSize(t *testing.T) {
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

	pathFile := "testfile.txt"
	expectedData := []byte("0123456789")
	// Create a test file
	err := os.WriteFile(pathFile, expectedData, 0644)
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}
	defer os.Remove(pathFile)

	contentRange := "bytes 0-10/10"

	req, err := http.NewRequest("GET", "http://localhost:8080/", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Range", contentRange)

	storeData.EXPECT().GetMetaData(gomock.Any(), gomock.Any(), gomock.Any()).Return(&store.MetaData{
		Size:     100,
		FileName: "testfile.txt",
		PathSave: "",
	}, nil)

	_, err = u.GetFileChunks(context.TODO(), 1, 1, req)
	assert.Error(t, err)
}
func TestUseCase_GetFileChunksBadRangeMin(t *testing.T) {
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

	pathFile := "testfile.txt"
	expectedData := []byte("0123456789")
	// Create a test file
	err := os.WriteFile(pathFile, expectedData, 0644)
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}
	defer os.Remove(pathFile)

	contentRange := "bytes 100-10/10"

	req, err := http.NewRequest("GET", "http://localhost:8080/", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Range", contentRange)

	storeData.EXPECT().GetMetaData(gomock.Any(), gomock.Any(), gomock.Any()).Return(&store.MetaData{
		Size:     10,
		FileName: "testfile.txt",
		PathSave: "",
	}, nil)

	_, err = u.GetFileChunks(context.TODO(), 1, 1, req)
	assert.Error(t, err)
}

func TestUseCase_GetFileChunksBadSizeNoTotalSize(t *testing.T) {
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

	pathFile := "testfile.txt"
	expectedData := []byte("0123456789")
	// Create a test file
	err := os.WriteFile(pathFile, expectedData, 0644)
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}
	defer os.Remove(pathFile)

	contentRange := "bytes 0-100/10"

	req, err := http.NewRequest("GET", "http://localhost:8080/", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Range", contentRange)

	storeData.EXPECT().GetMetaData(gomock.Any(), gomock.Any(), gomock.Any()).Return(&store.MetaData{
		Size:     10,
		FileName: "testfile.txt",
		PathSave: "",
	}, nil)

	_, err = u.GetFileChunks(context.TODO(), 1, 1, req)
	assert.Error(t, err)
}

func TestUseCase_GetFileChunksErrorGetFile(t *testing.T) {
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

	pathFile := "testfile.txt"
	expectedData := []byte("0123456789")
	// Create a test file
	err := os.WriteFile(pathFile, expectedData, 0644)
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}
	defer os.Remove(pathFile)

	contentRange := "bytes 0-10/10"

	req, err := http.NewRequest("GET", "http://localhost:8080/", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Range", contentRange)

	storeData.EXPECT().GetMetaData(gomock.Any(), gomock.Any(), gomock.Any()).Return(&store.MetaData{
		Size:     10,
		FileName: "testfile.txt",
		PathSave: "testfile.txt",
	}, nil)

	_, err = u.GetFileChunks(context.TODO(), 1, 1, req)
	assert.Error(t, err)
}

//getFile

func TestUseCase_getFile(t *testing.T) {
	ctx := context.Background()
	pathFile := "testfile.txt"
	byteStart := 0
	byteEnd := 10
	expectedData := []byte("0123456789")

	// Create a test file
	err := os.WriteFile(pathFile, expectedData, 0644)
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}
	defer os.Remove(pathFile)

	useCase := &UseCase{}

	data, err := useCase.getFile(ctx, pathFile, byteStart, byteEnd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !bytes.Equal(data, expectedData) {
		t.Errorf("expected %v, got %v", expectedData, data)
	}
}
func TestUseCase_getFileBadFileOpen(t *testing.T) {
	ctx := context.Background()
	pathFile := "testfile.txt"
	byteStart := 0
	byteEnd := 10

	useCase := &UseCase{}

	_, err := useCase.getFile(ctx, pathFile, byteStart, byteEnd)
	require.Error(t, err)

}
func TestUseCase_getFileBadSeek(t *testing.T) {
	ctx := context.Background()
	pathFile := "testfile.txt"
	byteStart := -10000
	byteEnd := 10
	expectedData := []byte("0123456789")

	// Create a test file
	err := os.WriteFile(pathFile, expectedData, 0644)
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}
	defer os.Remove(pathFile)

	useCase := &UseCase{}

	_, err = useCase.getFile(ctx, pathFile, byteStart, byteEnd)
	require.Error(t, err)
}

// GetData

func TestUseCase_GetData(t *testing.T) {
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

	storeData.EXPECT().GetData(gomock.Any(), gomock.Any(), gomock.Any()).Return(&store.UsersData{}, &store.DataFile{}, nil)

	_, err := u.GetData(context.TODO(), 1, 3)
	assert.NoError(t, err)
}
func TestUseCase_GetDataBadUserId(t *testing.T) {
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

	_, err := u.GetData(context.TODO(), 0, 3)
	assert.Error(t, err)
}
func TestUseCase_GetDataBadGetData(t *testing.T) {
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

	storeData.EXPECT().GetData(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil, fmt.Errorf("error"))

	_, err := u.GetData(context.TODO(), 1, 3)
	assert.Error(t, err)
}

// UpdateData
func TestUseCase_UpdateData(t *testing.T) {
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

	storeData.EXPECT().UpdateData(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(&store.UsersData{}, nil)

	_, err := u.UpdateData(context.TODO(), 1, 3, []byte("test"))
	assert.NoError(t, err)
}
func TestUseCase_UpdateDataBadGetUserId(t *testing.T) {
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

	_, err := u.UpdateData(context.TODO(), 0, 3, []byte("test"))
	assert.Error(t, err)
}
func TestUseCase_UpdateDataBadUpdateData(t *testing.T) {
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

	storeData.EXPECT().UpdateData(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error"))

	_, err := u.UpdateData(context.TODO(), 1, 3, []byte("test"))
	assert.Error(t, err)
}

// RemoveData
func TestUseCase_RemoveData(t *testing.T) {
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

	storeData.EXPECT().RemoveData(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	err := u.RemoveData(context.TODO(), int64(1), int64(1))
	assert.NoError(t, err)
}
func TestUseCase_RemoveDataBadUserid(t *testing.T) {
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

	err := u.RemoveData(context.TODO(), int64(0), int64(1))
	assert.Error(t, err)
}
func TestUseCase_RemoveDataBadUserDataId(t *testing.T) {
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

	err := u.RemoveData(context.TODO(), int64(1), int64(0))
	assert.Error(t, err)
}
func TestUseCase_RemoveDataBadRemove(t *testing.T) {
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

	storeData.EXPECT().RemoveData(gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("error"))

	err := u.RemoveData(context.TODO(), int64(1), int64(1))
	assert.Error(t, err)
}

// GetListData

func TestUseCase_GetListData(t *testing.T) {
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
	data := []store.UsersData{
		{
			UserId:   1,
			DataType: postgresql.TypeCreditCardData,
		},
		{
			UserId: 2,
		},
	}

	storeData.EXPECT().GetListData(gomock.Any(), gomock.Any()).Return(data, nil)

	_, err := u.GetListData(context.TODO(), int64(1))
	assert.NoError(t, err)
}
func TestUseCase_GetListDataBadUserId(t *testing.T) {
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

	_, err := u.GetListData(context.TODO(), int64(0))
	assert.Error(t, err)
}
func TestUseCase_GetListDataBadGetListData(t *testing.T) {
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

	storeData.EXPECT().GetListData(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error"))

	_, err := u.GetListData(context.TODO(), int64(1))
	assert.Error(t, err)
}

// UploadFile
func TestUseCase_UploadFile(t *testing.T) {
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
	fileSaver.EXPECT().UploadFile(gomock.Any(), gomock.Any()).Return(true, &TmpFile{}, nil)

	_, _, err := u.UploadFile("path", &http.Request{})
	require.NoError(t, err)
}

// UpdateBinaryFile
func TestUseCase_UpdateBinaryFile(t *testing.T) {
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
	namef := uuid.NewString() + ".txt"
	tmpFile := &TmpFile{
		PathFileSave: path.Join("tmp", "testfile.txt"),
		Uuid:         namef,
		Size:         1024,
	}
	err := os.MkdirAll(path.Dir(tmpFile.PathFileSave), os.ModePerm)
	require.NoError(t, err)
	f, err := os.Create(tmpFile.PathFileSave)
	require.NoError(t, err)
	err = f.Close()
	require.NoError(t, err)
	defer os.Remove(tmpFile.PathFileSave)

	storeData.EXPECT().UpdateBinaryFile(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(&store.UsersData{}, nil)
	_, err = u.UpdateBinaryFile(context.TODO(), 1, 1, tmpFile, []byte("data"))
	assert.NoError(t, err)
}

func TestUseCase_UpdateBinaryFileBadMove(t *testing.T) {
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
	namef := uuid.NewString() + ".txt"
	tmpFile := &TmpFile{
		PathFileSave: path.Join("tmp", "testfile.txt"),
		Uuid:         namef,
		Size:         1024,
	}

	_, err := u.UpdateBinaryFile(context.TODO(), userId, 1, tmpFile, []byte("data"))

	require.Error(t, err)
}
func TestUseCase_UpdateBinaryFileBadUpdate(t *testing.T) {
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
	namef := uuid.NewString() + ".txt"
	tmpFile := &TmpFile{
		PathFileSave: path.Join("tmp", "testfile.txt"),
		Uuid:         namef,
		Size:         1024,
	}
	err := os.MkdirAll(path.Dir(tmpFile.PathFileSave), os.ModePerm)
	require.NoError(t, err)
	f, err := os.Create(tmpFile.PathFileSave)
	require.NoError(t, err)
	err = f.Close()
	require.NoError(t, err)
	defer os.Remove(tmpFile.PathFileSave)

	storeData.EXPECT().UpdateBinaryFile(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(&store.UsersData{}, fmt.Errorf("error"))
	_, err = u.UpdateBinaryFile(context.TODO(), 1, 1, tmpFile, []byte("data"))
	assert.Error(t, err)
}

// upload
func TestUseCase_UploadFileBadUpload(t *testing.T) {
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
	fileSaver.EXPECT().UploadFile(gomock.Any(), gomock.Any()).Return(true, &TmpFile{}, fmt.Errorf("error"))

	_, _, err := u.UploadFile("path", &http.Request{})
	require.Error(t, err)
}

//ChangeAllData

func TestUseCase_ChangeAllData(t *testing.T) {
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

	storeData.EXPECT().ChangeAllData(gomock.Any(), gomock.Any(), gomock.Any()).Return([]store.UsersData{}, nil)

	_, err := u.ChangeAllData(context.TODO(), 1, time.Now())

	assert.NoError(t, err)
}
func TestUseCase_ChangeAllDataBaduserId(t *testing.T) {
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

	_, err := u.ChangeAllData(context.TODO(), 0, time.Now())

	assert.Error(t, err)
}
func TestUseCase_ChangeAllDataBadStore(t *testing.T) {
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

	storeData.EXPECT().ChangeAllData(gomock.Any(), gomock.Any(), gomock.Any()).Return([]store.UsersData{}, fmt.Errorf("error"))

	_, err := u.ChangeAllData(context.TODO(), 1, time.Now())

	assert.Error(t, err)
}
