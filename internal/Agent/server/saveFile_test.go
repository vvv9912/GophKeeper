package server

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"testing"
)

func TestNewSaveFile(t *testing.T) {
	// Create a temporary test directory
	tempDir := "./tmp_test"
	err := os.Mkdir(tempDir, 0755)
	if err != nil {
		t.Errorf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	fileName := "testFile.txt"
	expectedPath := path.Join(tempDir, fileName)

	saveFile, err := NewSaveFile(fileName)
	if err != nil {
		t.Fatalf("NewSaveFile() returned an error: %v", err)
	}

	// Check if the file was created
	if _, err := os.Stat("tmp/agent/testFile.txt"); os.IsNotExist(err) {
		t.Errorf("File was not created at the expected path: %v", expectedPath)
	}

	// Check if SaveFile struct was created correctly
	assert.Equal(t, fileName, saveFile.FileName, "FileName field is incorrect")
	assert.Equal(t, "tmp/agent/testFile.txt", saveFile.pathFile, "pathFile field is incorrect")

	// Close the file
	err = saveFile.f.Close()
	if err != nil {
		t.Errorf("Failed to close the file: %v", err)
	}
}
func TestSaveFile_GetPathFile(t *testing.T) {
	fileName := "test.txt"
	saveFile, err := NewSaveFile(fileName)
	if err != nil {
		t.Fatalf("Error creating SaveFile: %v", err)
	}

	if saveFile.GetPathFile() != "tmp/agent/test.txt" {
		t.Errorf("Expected path:  tmp/agent/test.txt, got: %s", saveFile.GetPathFile())
	}
}

func TestSaveFile_Write(t *testing.T) {
	fileName := "test.txt"
	saveFile, err := NewSaveFile(fileName)
	if err != nil {
		t.Fatalf("Error creating SaveFile: %v", err)
	}

	data := []byte("Hello, World!")
	n, err := saveFile.Write(data)

	if err != nil {
		t.Errorf("Error writing to file: %v", err)
	}
	if n != len(data) {
		t.Errorf("Expected to write %d bytes, wrote %d bytes", len(data), n)
	}
}

func TestSaveFile_CloseFile(t *testing.T) {
	fileName := "test.txt"
	saveFile, err := NewSaveFile(fileName)
	if err != nil {
		t.Fatalf("Error creating SaveFile: %v", err)
	}

	err = saveFile.CloseFile()
	if err != nil {
		t.Errorf("Error closing file: %v", err)
	}
}
