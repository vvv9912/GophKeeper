package app

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestReadKeyFromFile_Success(t *testing.T) {
	filePath := "test_key_valid.txt"
	keyContent := "12345678901234567890123456789012"
	err := os.WriteFile(filePath, []byte(keyContent), 0644)
	assert.NoError(t, err)
	defer os.Remove(filePath)

	key, err := readKeyFromFile(filePath)
	assert.NoError(t, err)
	assert.Equal(t, []byte(keyContent), key)
}
func TestReadKeyFromFile_FileNotExist(t *testing.T) {
	filePath := "non_existent_file.txt"

	_, err := readKeyFromFile(filePath)
	assert.Error(t, err)
}
func TestReadKeyFromFile_EmptyFile(t *testing.T) {
	filePath := "test_key_empty.txt"
	err := os.WriteFile(filePath, []byte(""), 0644)
	assert.NoError(t, err)
	defer os.Remove(filePath)

	_, err = readKeyFromFile(filePath)
	assert.Error(t, err)
}
