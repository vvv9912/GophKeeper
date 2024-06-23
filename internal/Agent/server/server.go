package server

import (
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
	pathGetFileSize  = "/api/data/fileSize"
	pathGetListData  = "/api/data"
	pathPing         = "/api/ping"
	pathUpdateData   = "/api/data/update"
	pathCheckUpdate  = "/api/data/CheckUpdate"
	pathUpdateBinary = "/api/data/updateBinary"
)

type AgentServer struct {
	host     string
	JWTToken string
	client   *resty.Client
}

func NewAgentServer(certFile, keyFile string, host string) *AgentServer {
	client := resty.New()

	client.SetTLSClientConfig(&tls.Config{
		InsecureSkipVerify: true,
	})

	cert1, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		log.Panicln("ERROR client certificate: %s", err)
	}

	client.SetCertificates(cert1)

	return &AgentServer{
		client: client,
		host:   host,
	}
}
