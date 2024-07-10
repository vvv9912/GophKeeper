package service

import (
	"crypto/rand"
	"errors"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestReader_ReadFile(t *testing.T) {
	f, err := os.CreateTemp("", "test.txt")
	if err != nil {
		t.Errorf("Error creating temporary file: %v", err)
	}
	defer os.Remove(f.Name())
	ttt := make([]byte, 256)
	_, err = rand.Read(ttt)
	if err != nil {
		t.Errorf("Error Read: %v", err)
	}
	_, err = f.Write(ttt)
	if err != nil {
		t.Errorf("Error write: %v", err)
	}

	r := NewReader(f.Name())
	n, err := r.NumChunk()
	if err != nil {
		t.Errorf("Error get n chunk: %v", err)
	}
	// Test reading a full chunk
	data, err := r.ReadFile(n)
	if err != nil {
		t.Errorf("Error reading file: %v", err)
	}
	expectedData := ttt
	if !reflect.DeepEqual(data, expectedData) {
		t.Errorf("Data is not as expected")
	}

}

func TestReader_ReadFile1(t *testing.T) {
	file, _ := os.CreateTemp("", "testfile")
	defer os.Remove(file.Name())
	file.WriteString("1234567890")
	reader := &Reader{
		SizeChunk: 5,
		Path:      file.Name(),
		f:         file,
		size:      10,
		maxChunk:  2,
		NameFile:  "testfile",
	}

	// Test valid chunk
	data, err := reader.ReadFile(1)
	assert.NoError(t, err)
	assert.Equal(t, []byte("12345"), data)
}
func TestReadFileHandlesNumChunkZeroOrNegative(t *testing.T) {

	// Initialize Reader
	file, _ := os.CreateTemp("", "testfile")
	defer os.Remove(file.Name())
	file.WriteString("1234567890")
	reader := &Reader{
		SizeChunk: 5,
		Path:      file.Name(),
		f:         file,
		size:      10,
		maxChunk:  2,
		NameFile:  "testfile",
	}

	// Test numChunk as 0
	data, err := reader.ReadFile(0)
	assert.Error(t, err)
	assert.Nil(t, data)

	// Test numChunk as negative value
	data, err = reader.ReadFile(-1)
	assert.Error(t, err)
	assert.Nil(t, data)
}
func TestNumChunk_SuccessfullyCalculatesNumberOfChunks(t *testing.T) {

	// Create a temporary file
	tmpFile, err := ioutil.TempFile("", "testfile")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	// Write data to the file
	data := make([]byte, 1024*1024) // 1MB
	if _, err := tmpFile.Write(data); err != nil {
		t.Fatal(err)
	}

	reader := &Reader{
		SizeChunk: 1024 * 1024, // 1MB
		Path:      tmpFile.Name(),
	}

	numChunks, err := reader.NumChunk()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expectedChunks := 1
	if numChunks != expectedChunks {
		t.Errorf("expected %d chunks, got %d", expectedChunks, numChunks)
	}
}
func TestNumChunk_ErrorWhenFileCannotBeOpened(t *testing.T) {

	reader := &Reader{
		SizeChunk: 1024 * 1024, // 1MB
		Path:      "nonexistentfile",
	}

	_, err := reader.NumChunk()
	if err == nil {
		t.Fatal("expected error, got none")
	}
	if !errors.Is(err, ErrOpenFile) {
		t.Errorf("expected error %v, got %v", ErrOpenFile, err)
	}
}
func TestNumChunk_ErrorWhenFileInfoCannotBeRetrieved(t *testing.T) {

	// Create a temporary file and remove it to simulate error in Stat()
	tmpFile, err := ioutil.TempFile("", "testfile")
	if err != nil {
		t.Fatal(err)
	}
	os.Remove(tmpFile.Name())

	reader := &Reader{
		SizeChunk: 1024 * 1024, // 1MB
		Path:      tmpFile.Name(),
	}

	_, err = reader.NumChunk()
	if err == nil {
		t.Fatal("expected error, got none")
	}
}
