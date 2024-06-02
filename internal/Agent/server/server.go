package server

import (
	"crypto/rsa"
)

var (
	pathSignIn      = "api/signIn"
	pathSignUp      = "api/signUp"
	pathGetData     = "api/data"
	pathChanges     = "api/changes"
	pathFile        = "/api/data/file"
	pathCredentials = "/api/data/credentials"
	pathCreditCard  = "/api/data/creditCard"
)

type AgentServer struct {
	publicKey *rsa.PublicKey
	host      string
	JWTToken  []string
}

func NewAgentServer(publicKey *rsa.PublicKey, host string) *AgentServer {
	return &AgentServer{
		publicKey: publicKey,
		host:      host,
	}
}
