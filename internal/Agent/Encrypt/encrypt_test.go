package Encrypt

import (
	"crypto/aes"
	"crypto/rand"
	"fmt"
	"github.com/stretchr/testify/require"
	"os"
	"runtime"
	"testing"
)

func TestEncrypt_Encrypt(t *testing.T) {
	key := make([]byte, 32)
	rand.Read(key)
	encr, _ := NewEncrypt(key)
	e, err := encr.Encrypt([]byte("hello world"))
	require.NoError(t, err)
	d, err := encr.Decrypt(e)
	require.NoError(t, err)
	require.Equal(t, []byte("hello world"), d)
}

func TestNewEncrypt(t *testing.T) {
	_, err := NewEncrypt([]byte("hello world"))
	require.Error(t, err)
}

func TestNewEncrypt2(t *testing.T) {
	key := make([]byte, 32)
	rand.Read(key)
	_, err := NewEncrypt(key)
	require.NoError(t, err)
}

func TestEncrypt_Encrypt1(t *testing.T) {
	key := make([]byte, 32)
	rand.Read(key)
	b, err := NewEncrypt(key)
	require.NoError(t, err)

	_, err = b.Encrypt([]byte("hello world"))
	require.NoError(t, err)

}

func TestEncrypt_Decrypt(t *testing.T) {
	key := make([]byte, 32)
	rand.Read(key)
	b, err := NewEncrypt(key)
	require.NoError(t, err)

	_, err = b.Decrypt([]byte(""))
	require.Error(t, err)
}
func TestEncryptFile_CreateOutputDir(t *testing.T) {
	key := []byte("example key 1234")
	block, err := aes.NewCipher(key)
	if err != nil {
		t.Fatalf("Failed to create cipher block: %v", err)
	}

	e := &Encrypt{block: block}

	inputFilePath := "input.txt"
	outputFilePath := "encrypted.txt"
	defer func() {
		_ = os.Remove(inputFilePath)
		_ = os.Remove(outputFilePath)
	}()

	err = os.WriteFile(inputFilePath, []byte("test data"), 0644)
	if err != nil {
		t.Fatalf("Failed to write input file: %v", err)
	}

	err = e.EncryptFile(inputFilePath, outputFilePath)
	if err != nil {
		t.Errorf("EncryptFile() error = %v, wantErr %v", err, false)
	}

	if _, err := os.Stat(outputFilePath); os.IsNotExist(err) {
		t.Errorf("Output file was not created")
	}
}
func TestEncryptFile_Success(t *testing.T) {
	key := []byte("example key 1234")
	block, err := aes.NewCipher(key)
	if err != nil {
		t.Fatalf("Failed to create cipher block: %v", err)
	}

	e := &Encrypt{block: block}

	inputFilePath := "input.txt"
	outputFilePath := "encrypted.txt"
	defer func() {
		_ = os.Remove(inputFilePath)
		_ = os.Remove(outputFilePath)
	}()

	err = os.WriteFile(inputFilePath, []byte("test data"), 0644)
	if err != nil {
		t.Fatalf("Failed to write input file: %v", err)
	}

	err = e.EncryptFile(inputFilePath, outputFilePath)
	if err != nil {
		t.Errorf("EncryptFile() error = %v, wantErr %v", err, false)
	}

	if _, err := os.Stat(outputFilePath); os.IsNotExist(err) {
		t.Errorf("Output file was not created")
	}
}
func TestEncryptFile_InputFileNotExist(t *testing.T) {
	key := []byte("example key 1234")
	block, err := aes.NewCipher(key)
	if err != nil {
		t.Fatalf("Failed to create cipher block: %v", err)
	}

	e := &Encrypt{block: block}

	inputFilePath := "testdata/nonexistent.txt"
	outputFilePath := "testdata/output/encrypted.txt"

	err = e.EncryptFile(inputFilePath, outputFilePath)
	if err == nil {
		t.Errorf("EncryptFile() error = %v, wantErr %v", err, true)
	}
}
func TestEncryptFile_OutputDirPermissionDenied(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("Skipping test on non-Linux")
	}

	key := []byte("example key 1234")
	block, err := aes.NewCipher(key)
	if err != nil {
		t.Fatalf("Failed to create cipher block: %v", err)
	}

	e := &Encrypt{block: block}

	inputFilePath := "input.txt"
	outputFilePath := "/root/encrypted.txt"

	err = os.WriteFile(inputFilePath, []byte("test data"), 0644)
	if err != nil {
		t.Fatalf("Failed to write input file: %v", err)
	}

	err = e.EncryptFile(inputFilePath, outputFilePath)
	if !os.IsPermission(err) {
		t.Errorf("EncryptFile() error = %v, wantErr %v", err, true)
	}
}
func TestEncrypt_DecryptFile23(t *testing.T) {
	key := []byte("example key 1234")
	block, err := aes.NewCipher(key)
	if err != nil {
		t.Fatalf("Failed to create cipher block: %v", err)
	}

	e := &Encrypt{block: block}

	inputFilePath := "input.txt"
	outputFilePath := "encrypted.txt"
	defer func() {
		_ = os.Remove(inputFilePath)
		_ = os.Remove(outputFilePath)
	}()

	err = os.WriteFile(inputFilePath, []byte("test data"), 0644)
	if err != nil {
		t.Fatalf("Failed to write input file: %v", err)
	}

	err = e.EncryptFile(inputFilePath, outputFilePath)
	if err != nil {
		t.Errorf("EncryptFile() error = %v, wantErr %v", err, false)
	}

	if _, err := os.Stat(outputFilePath); os.IsNotExist(err) {
		t.Errorf("Output file was not created")
	}

	err = e.DecryptFile(outputFilePath, "sss/decrypted.txt")
	if err != nil {
		t.Errorf("DecryptFile() error = %v, wantErr %v", err, false)
	}
}
func TestEncrypt_DecryptFile24(t *testing.T) {
	key := []byte("example key 1234")
	block, err := aes.NewCipher(key)
	if err != nil {
		t.Fatalf("Failed to create cipher block: %v", err)
	}

	e := &Encrypt{block: block}

	inputFilePath := "input.txt"
	outputFilePath := "encrypted.txt"
	defer func() {
		_ = os.Remove(inputFilePath)
		_ = os.Remove(outputFilePath)
	}()

	err = os.WriteFile(inputFilePath, []byte("test data"), 0644)
	if err != nil {
		t.Fatalf("Failed to write input file: %v", err)
	}

	err = e.EncryptFile(inputFilePath, outputFilePath)
	if err != nil {
		t.Errorf("EncryptFile() error = %v, wantErr %v", err, false)
	}

	if _, err := os.Stat(outputFilePath); os.IsNotExist(err) {
		t.Errorf("Output file was not created")
	}

	err = e.DecryptFile(outputFilePath, "/root/decrypted.txt")
	require.Error(t, err)
}
func TestEncrypt_DecryptFile22(t *testing.T) {
	key := []byte("example key 1234")
	block, err := aes.NewCipher(key)
	if err != nil {
		t.Fatalf("Failed to create cipher block: %v", err)
	}

	e := &Encrypt{block: block}

	inputFilePath := "input.txt"
	outputFilePath := "ccc/encrypted.txt"
	defer func() {
		_ = os.Remove(inputFilePath)
		_ = os.Remove(outputFilePath)
	}()

	err = os.WriteFile(inputFilePath, []byte("test data"), 0644)
	if err != nil {
		t.Fatalf("Failed to write input file: %v", err)
	}

	err = e.EncryptFile(inputFilePath, outputFilePath)
	if err != nil {
		t.Errorf("EncryptFile() error = %v, wantErr %v", err, false)
	}

	if _, err := os.Stat(outputFilePath); os.IsNotExist(err) {
		t.Errorf("Output file was not created")
	}

	err = e.DecryptFile(outputFilePath, "decrypted.txt")
	if err != nil {
		t.Errorf("DecryptFile() error = %v, wantErr %v", err, false)
	}
}
func TestEncrypt_DecryptFile(t *testing.T) {
	key := []byte("example key 1234")
	block, err := aes.NewCipher(key)
	if err != nil {
		t.Fatalf("Failed to create cipher block: %v", err)
	}

	e := &Encrypt{block: block}

	inputFilePath := "input.txt"
	outputFilePath := "encrypted.txt"
	defer func() {
		_ = os.Remove(inputFilePath)
		_ = os.Remove(outputFilePath)
	}()

	err = os.WriteFile(inputFilePath, []byte("test data"), 0644)
	if err != nil {
		t.Fatalf("Failed to write input file: %v", err)
	}

	err = e.EncryptFile(inputFilePath, outputFilePath)
	if err != nil {
		t.Errorf("EncryptFile() error = %v, wantErr %v", err, false)
	}

	if _, err := os.Stat(outputFilePath); os.IsNotExist(err) {
		t.Errorf("Output file was not created")
	}

	err = e.DecryptFile(outputFilePath, "decrypted.txt")
	if err != nil {
		t.Errorf("DecryptFile() error = %v, wantErr %v", err, false)
	}
}

func TestEncrypt_DecryptFile2(t *testing.T) {
	key := []byte("example key 1234")
	block, err := aes.NewCipher(key)
	if err != nil {
		t.Fatalf("Failed to create cipher block: %v", err)
	}

	e := &Encrypt{block: block}

	inputFilePath := "input.txt"
	outputFilePath := "encrypted.txt"
	defer func() {
		_ = os.Remove(inputFilePath)
		_ = os.Remove(outputFilePath)
	}()

	err = os.WriteFile(inputFilePath, []byte("test data"), 0644)
	if err != nil {
		t.Fatalf("Failed to write input file: %v", err)
	}

	err = e.EncryptFile(inputFilePath, outputFilePath)
	if err != nil {
		t.Errorf("EncryptFile() error = %v, wantErr %v", err, false)
	}

	if _, err := os.Stat(outputFilePath); os.IsNotExist(err) {
		t.Errorf("Output file was not created")
	}
	err = os.Remove(outputFilePath)
	if err != nil {
		// Обработка ошибки удаления файла
		fmt.Printf("Failed to remove file: %v\n", err)
	}

	// Создание нового файла
	newFile, err := os.Create(outputFilePath)
	if err != nil {
		// Обработка ошибки создания нового файла
		fmt.Printf("Failed to create new file: %v\n", err)
	}
	newFile.Close()
	err = e.DecryptFile(outputFilePath, "decr.txt")
	require.Error(t, err)
}
