package server

import (
	"GophKeeper/pkg/logger"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
	"log"
	"net/http"
)

func (a *AgentServer) SetJWTToken(token string) {
	a.JWTToken = token
}

func (a *AgentServer) SignIn(ctx context.Context, login, password string) (*User, error) {

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

	req := client.R()

	req.SetHeaders(map[string]string{
		"Content-Type": "application/json",
	})

	auth := Auth{
		Login:    login,
		Password: password,
	}

	resp, err := req.SetBody(auth).Post(a.host + pathSignIn)
	if err != nil {
		logger.Log.Error("Bad req", zap.Error(err))
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		var respError RespError
		err = json.Unmarshal(resp.Body(), &respError)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(respError.Message)

	}

	var user User
	err = json.Unmarshal(resp.Body(), &user)
	if err != nil {
		return nil, err
	}
	return &user, nil

}
func (a *AgentServer) SignUp(ctx context.Context, login, password string) (*User, error) {
	client := resty.New()
	req := client.R()

	req.SetHeaders(map[string]string{
		"Content-Type": "application/json",
	})

	auth := Auth{
		Login:    login,
		Password: password,
	}

	resp, err := req.SetBody(auth).Post(a.host + pathSignUp)
	if err != nil {
		logger.Log.Error("Bad req", zap.Error(err))
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		var respError RespError
		err = json.Unmarshal(resp.Body(), &respError)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(respError.Message)

	}

	var user User
	err = json.Unmarshal(resp.Body(), &user)
	if err != nil {
		return nil, err
	}
	return &user, nil

}
