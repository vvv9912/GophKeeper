package server

import (
	"GophKeeper/pkg/logger"
	"context"
	"encoding/json"
	"errors"
	"go.uber.org/zap"
	"net/http"
)

func (a *AgentServer) SetJWTToken(token string) {
	a.JWTToken = token
}
func (a *AgentServer) GetJWTToken() string {
	return a.JWTToken
}
func (a *AgentServer) SignIn(ctx context.Context, login, password string) (*User, error) {

	req := a.client.R()

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

	req := a.client.R()

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

func (a *AgentServer) Ping(ctx context.Context) error {

	req := a.client.R()

	resp, err := req.SetContext(ctx).Post(a.host + pathPing)
	if err != nil {
		logger.Log.Error("Bad req", zap.Error(err))
		return err
	}
	if resp.StatusCode() != http.StatusOK {

		return errors.New("Server not available")
	}

	return nil

}
