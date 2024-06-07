package service

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
)

func EncryptData(data []byte, publicKey *rsa.PublicKey) ([]byte, error) {
	if publicKey == nil {
		//todo add err
		return nil, fmt.Errorf("public key is nil")
	}

	encryptedData, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, data)
	if err != nil {
		return nil, err
	}
	return encryptedData, nil
}
func DecryptData(encryptedData []byte, privateKey *rsa.PrivateKey) ([]byte, error) {
	if privateKey == nil {
		//todo add err
		return nil, fmt.Errorf("private key is nil")
	}
	decryptedData, err := rsa.DecryptPKCS1v15(nil, privateKey, encryptedData)
	if err != nil {
		return nil, err
	}
	return decryptedData, nil
}
