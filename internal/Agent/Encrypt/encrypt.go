package Encrypt

import (
	"GophKeeper/pkg/logger"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"go.uber.org/zap"
	"io"
	"os"
	"path"
)

type Encrypt struct {
	key []byte
}

func NewEncrypt(key []byte) *Encrypt {
	return &Encrypt{key: key}
}

func (e *Encrypt) Encrypt(text []byte) ([]byte, error) {
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return nil, err
	}
	b := []byte(text)
	ciphertext := make([]byte, aes.BlockSize+len(b))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], b)
	return ciphertext, nil

}
func (e *Encrypt) Decrypt(ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return nil, err
	}
	if len(ciphertext) < aes.BlockSize {
		return nil, fmt.Errorf("cipherText too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)
	return ciphertext, nil
}

func (e *Encrypt) EncryptFile(inputFilePath string, outputFilePath string) error {
	inputFile, err := os.Open(inputFilePath)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	dir := path.Dir(outputFilePath)
	fmt.Println(dir)
	if err := os.MkdirAll(dir, 0755); err != nil {
		logger.Log.Error("Error create dir", zap.Error(err))
		return err
	}

	outputFile, err := os.OpenFile(outputFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	block, err := aes.NewCipher(e.key)
	if err != nil {
		return err
	}

	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	writer := &cipher.StreamWriter{S: stream, W: outputFile}

	if _, err := io.Copy(writer, inputFile); err != nil {
		return err
	}

	return nil
}
func (e *Encrypt) DecryptFile(inputFilePath string, outputFilePath string) error {
	inputFile, err := os.Open(inputFilePath)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	dir, _ := path.Split(outputFilePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	block, err := aes.NewCipher(e.key)
	if err != nil {
		return err
	}

	iv := make([]byte, aes.BlockSize)
	if _, err := inputFile.Read(iv); err != nil {
		return err
	}

	stream := cipher.NewCFBDecrypter(block, iv)
	reader := &cipher.StreamReader{S: stream, R: inputFile}

	if _, err := io.Copy(outputFile, reader); err != nil {
		return err
	}

	return nil
}
