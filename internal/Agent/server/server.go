package server

import (
	"crypto/rsa"
	"crypto/tls"
	"github.com/go-resty/resty/v2"
	"log"
)

var (
	pathSignIn      = "/api/signIn"
	pathSignUp      = "/api/signUp"
	pathGetData     = "/api/data"
	pathChanges     = "/api/changes"
	pathFile        = "/api/data/file"
	pathCredentials = "/api/data/credentials"
	pathCreditCard  = "/api/data/creditCard"
	pathFileChunks  = "/api/data/fileChunks"
	//pathGetFileChunks = "/api/data/GetFile"
	pathGetFileSize = "/api/data/fileSize"
	pathGetListData = "/api/data"
	pathPing        = "/api/ping"
	pathUpdateData  = "/api/data/update"
	pathCheckUpdate = "/api/data/CheckUpdate"
)

type AgentServer struct {
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
	host       string
	JWTToken   string
	client     *resty.Client
}

func NewAgentServer(publicKey *rsa.PublicKey, privateKey *rsa.PrivateKey, host string) *AgentServer {
	client := resty.New()

	//todo
	client.SetTLSClientConfig(&tls.Config{
		InsecureSkipVerify: true,
	})
	cert1, err := tls.LoadX509KeyPair("certs/cert.pem", "certs/key.pem")
	if err != nil {
		log.Fatalf("ERROR client certificate: %s", err)
	}

	client.SetCertificates(cert1)

	return &AgentServer{
		client:     client,
		publicKey:  publicKey,
		privateKey: privateKey,
		host:       host,
	}
}
