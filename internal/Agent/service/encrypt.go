package service

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"golang.org/x/crypto/chacha20poly1305"
)

func generateKey() ([]byte, error) {
	// Generate a random 256-bit key
	key := make([]byte, chacha20poly1305.KeySize)
	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}

	return key, nil
}
func encryptData(data []byte, key []byte) ([]byte, error) {
	// Create a new ChaCha20 cipher
	c, err := chacha20poly1305.New(key)
	if err != nil {
		return nil, err
	}

	// Generate a random nonce
	nonce := make([]byte, c.NonceSize())
	_, err = rand.Read(nonce)
	if err != nil {
		return nil, err
	}

	// Encrypt the data
	encrypted := c.Seal(nil, nonce, data, nil)

	return encrypted, nil
}

func EncryptData(data []byte, publicKey *rsa.PublicKey) ([]byte, error) {
	if publicKey == nil {
		//todo add err
		return nil, fmt.Errorf("public key is nil")
	}

	//publicKeyDER := x509.MarshalPKCS1PublicKey(publicKey)

	// Encrypt the text
	encryptedData, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, data)
	if err != nil {
		return nil, err
	}

	//encryptedData, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, data)
	//if err != nil {
	//	return nil, err
	//}
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
