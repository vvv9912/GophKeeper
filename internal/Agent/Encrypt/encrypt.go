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
	block cipher.Block
}

func NewEncrypt(key []byte) (*Encrypt, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		logger.Log.Error("Error new cipher", zap.Error(err))
		return nil, err
	}

	return &Encrypt{block: block}, nil
}

func (e *Encrypt) Encrypt(text []byte) ([]byte, error) {

	ciphertext := make([]byte, aes.BlockSize+len(text))

	iv := ciphertext[:aes.BlockSize]

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		logger.Log.Error("Error read full", zap.Error(err))
		return nil, err
	}

	stream := cipher.NewCFBEncrypter(e.block, iv)

	stream.XORKeyStream(ciphertext[aes.BlockSize:], text)

	return ciphertext, nil

}
func (e *Encrypt) Decrypt(ciphertext []byte) ([]byte, error) {

	if len(ciphertext) < aes.BlockSize {
		logger.Log.Error("cipherText too short")
		return nil, fmt.Errorf("cipherText too short")
	}
	iv := ciphertext[:aes.BlockSize]

	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(e.block, iv)

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

	logger.Log.Debug("dir", zap.String("dir", dir))

	if err := os.MkdirAll(dir, 0755); err != nil {
		logger.Log.Error("Error create dir", zap.Error(err))
		return err
	}

	outputFile, err := os.OpenFile(outputFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logger.Log.Error("Error create file", zap.Error(err))
		return err
	}
	defer outputFile.Close()

	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		logger.Log.Error("Error read iv", zap.Error(err))
		return err
	}

	stream := cipher.NewCFBEncrypter(e.block, iv)
	writer := &cipher.StreamWriter{S: stream, W: outputFile}

	if _, err := io.Copy(writer, inputFile); err != nil {
		logger.Log.Error("Error copy", zap.Error(err))
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
		logger.Log.Error("Error create dir", zap.Error(err))
		return err
	}
	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		logger.Log.Error("Error create file", zap.Error(err))
		return err
	}
	defer outputFile.Close()

	iv := make([]byte, aes.BlockSize)
	if _, err := inputFile.Read(iv); err != nil {
		logger.Log.Error("Error read iv", zap.Error(err))
		return err
	}

	stream := cipher.NewCFBDecrypter(e.block, iv)
	reader := &cipher.StreamReader{S: stream, R: inputFile}

	if _, err := io.Copy(outputFile, reader); err != nil {
		logger.Log.Error("Error copy file", zap.Error(err))
		return err
	}

	return nil
}
