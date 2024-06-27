package service

import (
	"bytes"
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

func TestNewSaveFiles(t *testing.T) {
	s, err := NewSaveFiles(time.Minute)
	require.NoError(t, err)
	require.NotNil(t, s.Chunks)
}

func TestSaveFiles_addNewSaveFile(t *testing.T) {
	s := &SaveFiles{
		Chunks:      make(map[string]TmpFile),
		defaultPath: "/tmp",
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "test.txt")
	if err != nil {
		t.Fatal(err)
	}
	part.Write([]byte("This is a test file"))
	writer.Close()

	req := httptest.NewRequest("POST", "/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Content-Range", "bytes 0-19/20")

	fileUpload, tmpFile, err := s.addNewSaveFile("uploads", req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if fileUpload {
		t.Fatalf("Expected fileUpload to be true, got %v", fileUpload)
	}

	if tmpFile == nil {
		t.Fatalf("Expected tmpFile to be non-nil")
	}
}
func TestSaveFiles_addNewSaveFileBadCreateFile(t *testing.T) {
	s := &SaveFiles{
		Chunks:      make(map[string]TmpFile),
		defaultPath: "/tmp",
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "test.txt")
	if err != nil {
		t.Fatal(err)
	}
	part.Write([]byte("This is a test file"))
	writer.Close()

	req := httptest.NewRequest("POST", "/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Content-Range", "bytes 0-19/20")

	fileUpload, tmpFile, err := s.addNewSaveFile("/root/", req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if fileUpload {
		t.Fatalf("Expected fileUpload to be true, got %v", fileUpload)
	}

	if tmpFile == nil {
		t.Fatalf("Expected tmpFile to be non-nil")
	}
}
func TestSaveFiles_addNewSaveFileBadFormFile(t *testing.T) {
	s := &SaveFiles{
		Chunks:      make(map[string]TmpFile),
		defaultPath: "/tmp",
	}

	req := httptest.NewRequest("POST", "/upload", nil)

	req.Header.Set("Content-Range", "bytes 0-19/20")

	fileUpload, tmpFile, err := s.addNewSaveFile("uploads", req)
	if err == nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if fileUpload {
		t.Fatalf("Expected fileUpload to be true, got %v", fileUpload)
	}

	if tmpFile != nil {
		t.Fatalf("Expected tmpFile to be non-nil")
	}
}
func TestSaveFiles_addNewSaveFileBadParserContentRage(t *testing.T) {
	s := &SaveFiles{
		Chunks:      make(map[string]TmpFile),
		defaultPath: "/tmp",
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "test.txt")
	if err != nil {
		t.Fatal(err)
	}
	part.Write([]byte("This is a test file"))
	writer.Close()

	req := httptest.NewRequest("POST", "/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Content-Range", "bytses 0-19/20")

	_, _, err = s.addNewSaveFile("uploads", req)
	require.Error(t, err)
}
func TestSaveFiles_addNewSaveFileRangeMin(t *testing.T) {
	s := &SaveFiles{
		Chunks:      make(map[string]TmpFile),
		defaultPath: "/tmp",
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "test.txt")
	if err != nil {
		t.Fatal(err)
	}
	part.Write([]byte("This is a test file"))
	writer.Close()

	req := httptest.NewRequest("POST", "/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Content-Range", "bytes 1-19/20")

	_, _, err = s.addNewSaveFile("uploads", req)
	require.Error(t, err)
}
func TestSaveFiles_addNewSaveFileRangeMax(t *testing.T) {
	s := &SaveFiles{
		Chunks:      make(map[string]TmpFile),
		defaultPath: "/tmp",
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "test.txt")
	if err != nil {
		t.Fatal(err)
	}
	part.Write([]byte("This is a test file"))
	writer.Close()

	req := httptest.NewRequest("POST", "/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Content-Range", "bytes 0-21/20")

	_, _, err = s.addNewSaveFile("uploads", req)
	require.Error(t, err)
}
func TestUploadFile_NewFileChunkUUid(t *testing.T) {
	s := &SaveFiles{
		Chunks:      make(map[string]TmpFile),
		defaultPath: "/tmp",
	}

	r := httptest.NewRequest("POST", "/upload", nil)
	r.Header.Set("Content-Range", "bytes 0-1023/1024")

	_, _, err := s.UploadFile("testpath", r)

	assert.Error(t, err)

}

func TestUploadFile_NewFileChunkNo(t *testing.T) {
	s := &SaveFiles{
		Chunks:      make(map[string]TmpFile),
		defaultPath: "/tmp",
	}

	r := httptest.NewRequest("POST", "/upload", nil)
	r.Header.Set("Uuid-Chunk", "testuuid")
	r.Header.Set("Content-Range", "bytes 0-1023/1024")

	_, _, err := s.UploadFile("testpath", r)

	assert.Error(t, err)
	//assert.True(t, fileUpload)
	//assert.NotNil(t, tmpFile)
}
func TestUploadFile_NewFileChunkGetFile(t *testing.T) {
	s := &SaveFiles{
		Chunks:      make(map[string]TmpFile),
		defaultPath: "/tmp",
	}
	s.Chunks["testuuid"] = TmpFile{
		OriginalFileName: "test.txt",
		PathFileSave:     "/tmp/test.txt",
		Uuid:             "testuuid",
		Size:             1024,
	}

	r := httptest.NewRequest("POST", "/upload", nil)
	r.Header.Set("Uuid-Chunk", "testuuid")
	r.Header.Set("Content-Range", "bytes 0-1023/1024")

	_, _, err := s.UploadFile("testpath", r)

	assert.Error(t, err)

}
func TestUploadFile_NewFileChunk(t *testing.T) {
	s := &SaveFiles{
		Chunks:      make(map[string]TmpFile),
		defaultPath: "/tmp",
	}
	s.Chunks["testuuid"] = TmpFile{
		OriginalFileName: "test.txt",
		PathFileSave:     "/tmp/test.txt",
		Uuid:             "testuuid",
		Size:             1024,
	}
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "test.txt")
	if err != nil {
		t.Fatal(err)
	}
	part.Write([]byte("This is a test file"))
	writer.Close()

	r := httptest.NewRequest("POST", "/upload", body)
	r.Header.Set("Uuid-Chunk", "testuuid")
	r.Header.Set("Content-Range", "bytes 0-1023/1024")
	r.Header.Set("Content-Type", writer.FormDataContentType())
	_, _, err = s.UploadFile("testpath", r)

	assert.NoError(t, err)

}

func TestSaveFiles_DeleteFile(t *testing.T) {
	s := &SaveFiles{
		Chunks: map[string]TmpFile{
			"test-uuid": {
				PathFileSave: "test-path",
			},
		},
		fileSave: false,
	}
	f, err := os.Create("test-path")
	require.NoError(t, err)

	defer func() {
		_ = f.Close()
		_ = os.Remove("test-path")
	}()

	err = s.DeleteFile("test-uuid")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if _, exists := s.Chunks["test-uuid"]; exists {
		t.Fatalf("expected file to be deleted from Chunks map")
	}
}
func TestSaveFiles_DeleteFileNook(t *testing.T) {
	s := &SaveFiles{}
	f, err := os.Create("test-path")
	require.NoError(t, err)

	defer func() {
		_ = f.Close()
		_ = os.Remove("test-path")
	}()

	err = s.DeleteFile("test-uuid")
	require.Error(t, err)

}
func TestSaveFiles_DeleteFileNoDelete(t *testing.T) {
	s := &SaveFiles{
		Chunks: map[string]TmpFile{
			"test-uuid": {
				PathFileSave: "test-path",
			},
		},
		fileSave: true,
	}

	err := s.DeleteFile("test-uuid")

	require.NoError(t, err)

}

func TestSaveFiles_RunCronDeleteFiles(t *testing.T) {

	s := &SaveFiles{
		Chunks: map[string]TmpFile{
			"test-uuid": {
				PathFileSave: "test-path",
				LastUpdate:   time.Now().Add(-time.Hour),
			},
		},
		tmpFileLifeTime: time.Second,
		fileSave:        false,
	}

	ctx, _ := context.WithTimeout(context.Background(), time.Second)
	err := s.RunCronDeleteFiles(ctx)
	time.Sleep(time.Second)
	require.NoError(t, err)

	if _, exists := s.Chunks["test-uuid"]; !exists {
		t.Fatalf("expected file to be deleted from Chunks map")
	}

}
func TestSaveFiles_RunCronDeleteFiles2(t *testing.T) {

	s := &SaveFiles{
		Chunks: map[string]TmpFile{
			"test-uuid": {
				PathFileSave: "test-path",
				LastUpdate:   time.Now().Add(-time.Hour),
			},
		},
		tmpFileLifeTime: time.Microsecond,
		fileSave:        false,
	}

	f, err := os.Create("test-path")
	require.NoError(t, err)
	defer func() {
		f.Close()
		os.Remove("test-path")

	}()

	err = s.RunCronDeleteFiles(context.Background())
	time.Sleep(3 * time.Second)
	require.NoError(t, err)
	fmt.Println(s.Chunks)
	if _, exists := s.Chunks["test-uuid"]; exists {
		t.Fatalf("expected file to be deleted from Chunks map")
	}

}

func TestSaveFiles_generateUuid(t *testing.T) {
	s := &SaveFiles{
		Chunks: map[string]TmpFile{
			"test-uuid": {
				PathFileSave: "test-path",
				LastUpdate:   time.Now().Add(-time.Hour),
			},
		},
		tmpFileLifeTime: time.Microsecond,
		fileSave:        false,
	}
	s.generateUuid()
}

func TestSaveFiles_saveFile(t *testing.T) {
	s := &SaveFiles{
		Chunks: map[string]TmpFile{
			"test-uuid": {
				PathFileSave: "test-path",
				LastUpdate:   time.Now().Add(-time.Hour),
			},
		},
		tmpFileLifeTime: time.Microsecond,
		fileSave:        false,
	}

	f, err := os.Create("test-path")
	require.NoError(t, err)
	defer func() {
		f.Close()
		os.Remove("test-path")

	}()

	content := "Hello, World!"
	file := strings.NewReader(content)
	s.saveFile("test-path", file)
}
func TestSaveFiles_saveFileNoFile(t *testing.T) {
	s := &SaveFiles{
		Chunks: map[string]TmpFile{
			"test-uuid": {
				PathFileSave: "test-path",
				LastUpdate:   time.Now().Add(-time.Hour),
			},
		},
		tmpFileLifeTime: time.Microsecond,
		fileSave:        false,
	}

	content := "Hello, World!"
	file := strings.NewReader(content)
	s.saveFile("", file)
}

func TestParserContentRange(t *testing.T) {
	contentRangeHeader := "bytes 0-499/12345"

	rangeMin, rangeMax, totalFileSize, err := ParserContentRange(contentRangeHeader)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if rangeMin != 0 {
		t.Errorf("Expected rangeMin to be 0, got %d", rangeMin)
	}

	if rangeMax != 499 {
		t.Errorf("Expected rangeMax to be 499, got %d", rangeMax)
	}

	if totalFileSize != 12345 {
		t.Errorf("Expected totalFileSize to be 12345, got %d", totalFileSize)
	}
}
func TestParserContentRangeBad1(t *testing.T) {
	contentRangeHeader := "0-499/12345"

	_, _, _, err := ParserContentRange(contentRangeHeader)
	require.Error(t, err)
}

func TestFileUploadCompleted_SuccessNoSize(t *testing.T) {
	s := &SaveFiles{}
	req, err := http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Range", "bytes 0-99/100")

	_, err = s.FileUploadCompleted(100, req)
	require.Error(t, err)
}
func TestFileUploadCompleted_Success(t *testing.T) {
	s := &SaveFiles{}
	req, err := http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Range", "bytes 0-99/100")

	_, err = s.FileUploadCompleted(99, req)
	require.NoError(t, err)
}
func TestFileUploadCompleted_Success2(t *testing.T) {
	s := &SaveFiles{}
	req, err := http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Range", "bytes 0-100/100")

	_, err = s.FileUploadCompleted(100, req)
	require.NoError(t, err)
}
func TestFileUploadCompleted_BadContentRange(t *testing.T) {
	s := &SaveFiles{}
	req, err := http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Range", "bytesss 0-100/100")

	_, err = s.FileUploadCompleted(100, req)
	require.Error(t, err)
}
