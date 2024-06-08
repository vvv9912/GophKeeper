package server

import (
	"GophKeeper/pkg/logger"
	"GophKeeper/pkg/store"
	"GophKeeper/pkg/store/sqllite"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
	"mime/multipart"
	"net/http"
	"strconv"
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

func (a *AgentServer) PostCrateFileStartChunks(ctx context.Context, data []byte, fileName string, uuidChunk string, nStart int, nEnd int, maxSize int, reqData []byte) (string, error) {
	req := a.client.R()

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// Добавление файла в форму
	fileWriter, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		return "", err
	}

	// Записываем данные
	_, err = fileWriter.Write(data)
	if err != nil {
		logger.Log.Error("Bad req", zap.Error(err))
		return "", err
	}
	req.SetHeaders(map[string]string{
		"Authorization": "Bearer " + a.JWTToken,
		"Content-Range": fmt.Sprintf("bytes %d-%d/%d", nStart, nEnd, maxSize),
	})

	if uuidChunk != "" {
		req.SetHeaders(map[string]string{
			"Uuid-chunk": uuidChunk,
		})
	}

	var resp *resty.Response
	if nEnd != maxSize {
		resp, err = req.SetContext(ctx).SetBody(&buf).SetMultipartFields(
			&resty.MultipartField{
				Param:       "file",
				FileName:    fileName,
				ContentType: writer.FormDataContentType(),
				Reader:      bytes.NewReader(data)},
		).Post(a.host + pathFileChunks)

	} else {
		// Последний чанк, передаем информацию о файле

		resp, err = req.SetContext(ctx).SetBody(&buf).SetMultipartFields(
			&resty.MultipartField{
				Param:       "file",
				FileName:    fileName,
				ContentType: writer.FormDataContentType(),
				Reader:      bytes.NewReader(data)},
			&resty.MultipartField{
				Param:       "info",
				ContentType: "application/json",
				Reader:      bytes.NewReader(reqData),
			},
		).Post(a.host + pathFileChunks)
	}

	if err != nil {
		logger.Log.Error("Bad req", zap.Error(err))
		return "", err
	}

	if resp.StatusCode() != http.StatusOK {
		var respError RespError
		err = json.Unmarshal(resp.Body(), &respError)
		if err != nil {
			return "", err
		}

		return "", errors.New(respError.Message)

	}

	uuidChunk = resp.Header().Get("Uuid-chunk")

	return uuidChunk, nil
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

func (a *AgentServer) GetData(ctx context.Context, userDataId int64) ([]byte, error) {
	req := a.client.R()
	req.SetHeaders(map[string]string{
		"Authorization": "Bearer " + a.JWTToken,
	})
	resp, err := req.SetContext(ctx).Get(a.host + pathGetData + "/" + strconv.Itoa(int(userDataId)))
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

	var Data RespUsersData
	err = json.Unmarshal(resp.Body(), &Data)
	if err != nil {
		return nil, err
	}

	if Data.InfoUsersData.DataType == sqllite.TypeFile {
		// логика скачивания
	}

	// Готовим ответ
	// Тут расшифрока
	return Data.EncryptData.EncryptData, nil
}
