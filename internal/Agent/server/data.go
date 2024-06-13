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
	"github.com/google/uuid"
	"go.uber.org/zap"
	"math"
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

// Передача бинарного файла
func (a *AgentServer) PostCrateFileStartChunks(ctx context.Context, data []byte, fileName string, uuidChunk string, nStart int, nEnd int, maxSize int, reqData []byte) (string, *RespData, error) {
	req := a.client.R()

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// Добавление файла в форму
	fileWriter, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		logger.Log.Error("Bad req", zap.Error(err))
		return "", nil, err
	}

	// Записываем данные
	_, err = fileWriter.Write(data)
	if err != nil {
		logger.Log.Error("Bad req", zap.Error(err))
		return "", nil, err
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
		fmt.Println("чанк")
		resp, err = req.SetContext(ctx).SetBody(&buf).SetMultipartFields(
			&resty.MultipartField{
				Param:       "file",
				FileName:    fileName,
				ContentType: writer.FormDataContentType(),
				Reader:      bytes.NewReader(data)},
		).Post(a.host + pathFileChunks)
	} else {
		// Последний чанк, передаем информацию о файле
		fmt.Println("Последний чанк")
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
		return "", nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		var respError RespError
		err = json.Unmarshal(resp.Body(), &respError)
		if err != nil {
			logger.Log.Error("Bad resp", zap.Error(err), zap.Int("status_code", resp.StatusCode()))
			return "", nil, err
		}

		return "", nil, errors.New(respError.Message)

	}

	uuidChunk = resp.Header().Get("Uuid-chunk")
	if len(resp.Body()) == 0 {
		return uuidChunk, nil, nil
	}

	var respData RespData
	err = json.Unmarshal(resp.Body(), &respData)
	if err != nil {
		logger.Log.Error("Bad resp", zap.Error(err))
		return "", nil, err
	}
	return uuidChunk, &respData, nil
}

// Передача любых текстовых данных
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
		fileSize, err := a.getFileSize(ctx, Data.InfoUsersData.UserDataId)
		if err != nil {
			return nil, err
		}
		// Скачиваем
		pathSaveFile, err := a.getFileData(ctx, Data.InfoUsersData.UserDataId, fileSize)
		if err != nil {
			return nil, err
		}
		//todo rename
		return []byte("файл скачен путь: " + pathSaveFile), nil
	}

	// Готовим ответ todo
	return Data.EncryptData.EncryptData, nil
}

func (a *AgentServer) getFileSize(ctx context.Context, userDataId int64) (int64, error) {
	req := a.client.R()
	req.SetHeaders(map[string]string{
		"Authorization": "Bearer " + a.JWTToken,
	})
	resp, err := req.SetContext(ctx).Get(a.host + pathGetFileSize + "/" + strconv.Itoa(int(userDataId)))
	if err != nil {
		logger.Log.Error("Bad req", zap.Error(err))
		return 0, err
	}

	if resp.StatusCode() != http.StatusOK {
		var respError RespError
		err = json.Unmarshal(resp.Body(), &respError)
		if err != nil {
			return 0, err
		}
		return 0, errors.New(respError.Message)
	}
	reqFileSize := struct {
		FileSize int64 `json:"fileSize"`
	}{}

	err = json.Unmarshal(resp.Body(), &reqFileSize)
	if err != nil {
		return 0, err
	}
	return reqFileSize.FileSize, nil
}

func (a *AgentServer) getFileData(ctx context.Context, userDataId int64, fileSize int64) (string, error) {

	sizeChunk := 1024 * 1024

	nChunk := math.Ceil(float64(fileSize) / float64(sizeChunk))

	saveFile, err := NewSaveFile(uuid.NewString())
	if err != nil {
		return "", err
	}

	for i := 1; i <= int(nChunk); i++ {
		startChunk := (i - 1) * sizeChunk

		endChunk := i * sizeChunk

		if i == int(nChunk) {
			endChunk = int(fileSize)
		}

		fileChunk, err := a.getFileChunks(ctx, userDataId, fileSize, startChunk, endChunk)
		if err != nil {
			return "", err
		}
		if _, err := saveFile.Write(fileChunk); err != nil {
			logger.Log.Error("Error save file", zap.Error(err))
			return "", err
		}
	}
	return saveFile.GetPathFile(), nil
}

func (a *AgentServer) getFileChunks(ctx context.Context, userDataId int64, fileSize int64, startChunk int, endChunk int) ([]byte, error) {
	req := a.client.R()
	req.SetHeaders(map[string]string{
		"Authorization": "Bearer " + a.JWTToken,
		"Content-Range": fmt.Sprintf("bytes %d-%d/%d", startChunk, endChunk, fileSize),
	})
	resp, err := req.SetContext(ctx).Get(a.host + pathFileChunks + "/" + strconv.Itoa(int(userDataId)))
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

	return resp.Body(), err
}
func (a *AgentServer) GetListData(ctx context.Context) ([]byte, error) {
	req := a.client.R()
	req.SetHeaders(map[string]string{
		"Authorization": "Bearer " + a.JWTToken,
	})
	resp, err := req.SetContext(ctx).Get(a.host + pathGetListData)
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
	return resp.Body(), err
}

func (a *AgentServer) CheckUpdate(ctx context.Context, userDataid int64, updateAt *time.Time) (bool, error) {

	req := a.client.R()
	req.SetHeaders(map[string]string{
		"Authorization":    "Bearer " + a.JWTToken,
		"Last-Time-Update": updateAt.Format("2006-01-02 15:04:05.999999"),
	})
	resp, err := req.SetContext(ctx).Post(a.host + pathCheckUpdate + "/" + strconv.Itoa(int(userDataid)))
	if err != nil {
		logger.Log.Error("Bad req", zap.Error(err))
		return false, err
	}
	if resp.StatusCode() != http.StatusOK {
		var respError RespError
		err = json.Unmarshal(resp.Body(), &respError)
		if err != nil {
			return false, err
		}
		return false, errors.New(respError.Message)
	}

	var response struct {
		Status bool `json:"status"`
	}

	err = json.Unmarshal(resp.Body(), &response)
	if err != nil {
		logger.Log.Error("Bad req", zap.Error(err))
		return false, err
	}

	return response.Status, err
}

// PostUpdateData - обновление данных (кроме бинарных)
func (a *AgentServer) PostUpdateData(ctx context.Context, userDataId int64, data []byte) (*RespData, error) {
	req := a.client.R()

	req.SetHeaders(map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + a.JWTToken,
	})

	reqData := ReqData{
		Data: data,
	}

	resp, err := req.SetContext(ctx).SetBody(reqData).Post(a.host + pathUpdateData + "/" + strconv.Itoa(int(userDataId)))
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

// Передача бинарного файла
func (a *AgentServer) PostUpdateBinaryFile(ctx context.Context, data []byte, fileName string, uuidChunk string, nStart int, nEnd int, maxSize int, reqData []byte, userDataId int64) (string, *RespData, error) {
	req := a.client.R()

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// Добавление файла в форму
	fileWriter, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		logger.Log.Error("Bad req", zap.Error(err))
		return "", nil, err
	}

	// Записываем данные
	_, err = fileWriter.Write(data)
	if err != nil {
		logger.Log.Error("Bad req", zap.Error(err))
		return "", nil, err
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
		).Post(a.host + pathUpdateBinary + "/" + strconv.Itoa(int(userDataId)))
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
		).Post(a.host + pathUpdateBinary + "/" + strconv.Itoa(int(userDataId)))
	}

	if err != nil {
		logger.Log.Error("Bad req", zap.Error(err))
		return "", nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		var respError RespError
		err = json.Unmarshal(resp.Body(), &respError)
		if err != nil {
			logger.Log.Error("Bad resp", zap.Error(err), zap.Int("status_code", resp.StatusCode()))
			return "", nil, err
		}

		return "", nil, errors.New(respError.Message)

	}

	uuidChunk = resp.Header().Get("Uuid-chunk")
	if len(resp.Body()) == 0 {
		return uuidChunk, nil, nil
	}

	var respData RespData
	err = json.Unmarshal(resp.Body(), &respData)
	if err != nil {
		logger.Log.Error("Bad resp", zap.Error(err))
		return "", nil, err
	}
	return uuidChunk, &respData, nil
}
