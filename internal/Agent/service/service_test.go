package service

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math/big"
	"os"
	"testing"
	"time"
)

func TestSuccessfullyCreatesServiceInstanceWithValidInputs(t *testing.T) {
	key := []byte("examplekey123456")
	certFile := "path/to/certFile"
	keyFile := "path/to/keyFile"
	serverDns := "example.com"
	db := &sqlx.DB{}
	a := func() (string, string, error) {
		// Генерация закрытого ключа
		privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			return "", "", err
		}

		// Генерация сертификата
		template := x509.Certificate{
			SerialNumber: big.NewInt(1),
			Subject: pkix.Name{
				Country:      []string{"US"},
				Organization: []string{"Example"},
			},
			NotBefore:   time.Now(),
			NotAfter:    time.Now().AddDate(1, 0, 0),
			KeyUsage:    x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
			ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		}

		certBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
		if err != nil {
			return "", "", err
		}

		// Запись сертификата в файл cert.pem
		certFile := "cert.pem"
		certOut, err := os.Create(certFile)
		if err != nil {
			return "", "", err
		}
		pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: certBytes})
		certOut.Close()

		// Запись закрытого ключа в файл key.pem
		keyFile := "key.pem"
		keyOut, err := os.Create(keyFile)
		if err != nil {
			return "", "", err
		}
		pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)})
		keyOut.Close()

		return certFile, keyFile, nil

	}
	certFile, keyFile, err := a()
	defer func() {
		os.Remove(certFile)
		os.Remove(keyFile)
	}()
	require.NoError(t, err)

	service := NewServiceAgent(db, key, certFile, keyFile, serverDns)

	assert.NotNil(t, service)
	assert.NotNil(t, service.UseCaser)
}
