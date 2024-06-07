package server

import (
	"GophKeeper/pkg/logger"
	"GophKeeper/pkg/store"
	"context"
	"encoding/json"
	"errors"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func (a *AgentServer) PostCredentials(ctx context.Context, data *ReqData) (*RespData, error) {
	req := a.client.R()

	req.SetHeaders(map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + a.JWTToken,
	})

	resp, err := req.SetContext(ctx).SetBody(data).Post(a.host + pathCredentials)
	if err != nil {
		logger.Log.Error("Bad req", zap.Error(err))
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		var respError RespError
		err = json.Unmarshal(resp.Body(), &respError)
		if err != nil {
			logger.Log.Error("Bad resp", zap.Error(err), zap.Int("status_code", resp.StatusCode()))
			return nil, err
		}

		return nil, errors.New(respError.Message)
	}
	var respData RespData
	err = json.Unmarshal(resp.Body(), &respData)
	if err != nil {
		logger.Log.Error("Bad resp", zap.Error(err))
		return nil, err
	}

	return &respData, nil
}

func (a *AgentServer) PostCrateFile(ctx context.Context, data *ReqData) (*RespData, error) {

	req := a.client.R()

	req.SetHeaders(map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + a.JWTToken,
	})

	resp, err := req.SetContext(ctx).SetBody(data).Post(a.host + pathFile)
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
	var respData RespData
	err = json.Unmarshal(resp.Body(), &respData)
	if err != nil {
		logger.Log.Error("Bad resp", zap.Error(err))
		return nil, err
	}
	return &respData, nil
}
func (a *AgentServer) PostCreditCard(ctx context.Context, data *ReqData) (*RespData, error) {

	req := a.client.R()

	req.SetHeaders(map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + a.JWTToken,
	})

	resp, err := req.SetContext(ctx).SetBody(data).Post(a.host + pathCreditCard)
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

	var respData RespData
	err = json.Unmarshal(resp.Body(), &respData)
	if err != nil {
		logger.Log.Error("Bad resp", zap.Error(err))
		return nil, err
	}
	return &respData, nil
}
func (a *AgentServer) GetCheckChanges(ctx context.Context, data *ReqData, lastTime time.Time) ([]store.UsersData, error) {

	req := a.client.R()

	req.SetHeaders(map[string]string{
		"Content-Type":     "application/json",
		"Authorization":    "Bearer " + a.JWTToken,
		"Last-Time-Update": lastTime.Format("2006-01-02 15:04:05.999999"),
	})

	resp, err := req.SetContext(ctx).SetBody(data).Post(a.host + pathChanges)
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

	var usersData []store.UsersData
	err = json.Unmarshal(resp.Body(), &usersData)
	if err != nil {
		return nil, err
	}
	return usersData, nil
}
