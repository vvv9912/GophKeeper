package handler

import (
	"GophKeeper/internal/Agent/server"
	service2 "GophKeeper/internal/Server/service"
	mock_service "GophKeeper/internal/Server/service/mocks"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestHandler_HandlerPostCredentials(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	u := mock_service.NewMockUseCaser(ctrl)
	s := service2.Service{
		UseCaser: u,
	}
	u.EXPECT().CreateCredentials(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(&service2.RespData{}, nil)
	reqData := ReqData{
		Name:        "Name",
		Description: "Desc",
		Data:        []byte("data"),
	}

	reqBody, err := json.Marshal(reqData)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/Credentials", strings.NewReader(string(reqBody)))
	ctx := context.Background()
	ctx = context.WithValue(ctx, "UserId", int64(1))

	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	h := Handler{service: &s}
	h.HandlerPostCredentials(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

}
func TestHandler_HandlerPostCredentialsBadUserId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	u := mock_service.NewMockUseCaser(ctrl)
	s := service2.Service{
		UseCaser: u,
	}
	reqData := ReqData{
		Name:        "Name",
		Description: "Desc",
		Data:        []byte("data"),
	}

	reqBody, err := json.Marshal(reqData)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/Credentials", strings.NewReader(string(reqBody)))

	w := httptest.NewRecorder()

	h := Handler{service: &s}
	h.HandlerPostCredentials(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

}
func TestHandler_HandlerPostCredentialsBadDecodeBody(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	u := mock_service.NewMockUseCaser(ctrl)
	s := service2.Service{
		UseCaser: u,
	}

	reqBody, err := json.Marshal("data")
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/Credentials", strings.NewReader(string(reqBody)))
	ctx := context.Background()
	ctx = context.WithValue(ctx, "UserId", int64(1))

	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	h := Handler{service: &s}
	h.HandlerPostCredentials(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

}

func TestHandler_HandlerPostCredentialsBadCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	u := mock_service.NewMockUseCaser(ctrl)
	s := service2.Service{
		UseCaser: u,
	}
	u.EXPECT().CreateCredentials(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error"))
	reqData := ReqData{
		Name:        "Name",
		Description: "Desc",
		Data:        []byte("data"),
	}

	reqBody, err := json.Marshal(reqData)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/Credentials", strings.NewReader(string(reqBody)))
	ctx := context.Background()
	ctx = context.WithValue(ctx, "UserId", int64(1))

	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	h := Handler{service: &s}
	h.HandlerPostCredentials(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

}

// test creditCard
func TestHandler_HandlerPostCreditCard(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	u := mock_service.NewMockUseCaser(ctrl)
	s := service2.Service{
		UseCaser: u,
	}
	u.EXPECT().CreateCreditCard(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(&service2.RespData{}, nil)
	reqData := ReqData{
		Name:        "Name",
		Description: "Desc",
		Data:        []byte("data"),
	}

	reqBody, err := json.Marshal(reqData)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/CreditCard", strings.NewReader(string(reqBody)))
	ctx := context.Background()
	ctx = context.WithValue(ctx, "UserId", int64(1))

	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	h := Handler{service: &s}
	h.HandlerPostCreditCard(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

}
func TestHandler_HandlerPostCreditCardBadUserId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	u := mock_service.NewMockUseCaser(ctrl)
	s := service2.Service{
		UseCaser: u,
	}
	reqData := ReqData{
		Name:        "Name",
		Description: "Desc",
		Data:        []byte("data"),
	}

	reqBody, err := json.Marshal(reqData)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/CreditCard", strings.NewReader(string(reqBody)))

	w := httptest.NewRecorder()

	h := Handler{service: &s}
	h.HandlerPostCreditCard(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

}
func TestHandler_HandlerPostCreditCardBadJson(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	u := mock_service.NewMockUseCaser(ctrl)
	s := service2.Service{
		UseCaser: u,
	}

	reqBody, err := json.Marshal("data")
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/CreditCard", strings.NewReader(string(reqBody)))
	ctx := context.Background()
	ctx = context.WithValue(ctx, "UserId", int64(1))

	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	h := Handler{service: &s}
	h.HandlerPostCreditCard(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

}
func TestHandler_HandlerPostCreditCardBadCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	u := mock_service.NewMockUseCaser(ctrl)
	s := service2.Service{
		UseCaser: u,
	}
	u.EXPECT().CreateCreditCard(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error"))
	reqData := ReqData{
		Name:        "Name",
		Description: "Desc",
		Data:        []byte("data"),
	}

	reqBody, err := json.Marshal(reqData)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/CreditCard", strings.NewReader(string(reqBody)))
	ctx := context.Background()
	ctx = context.WithValue(ctx, "UserId", int64(1))

	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	h := Handler{service: &s}
	h.HandlerPostCreditCard(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

}

func Test_validateFile(t *testing.T) {
	header := &multipart.FileHeader{Size: 1}
	err := validateFile(header)
	assert.NoError(t, err)
}
func TestValidateFileLogsSuccesss(t *testing.T) {
	header := &multipart.FileHeader{Size: 0}
	err := validateFile(header)
	assert.Error(t, err)
}

// test PostChunkCreateFile
func TestHandler_HandlerPostChunkCrateFilePart1Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	u := mock_service.NewMockUseCaser(ctrl)
	s := service2.Service{
		UseCaser: u,
	}
	u.EXPECT().UploadFile(gomock.Any(), gomock.Any()).Return(false, &service2.TmpFile{}, nil)
	//u.EXPECT().CreateFileChunks(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(&service2.RespData{}, nil)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "test.txt")
	part.Write([]byte("test content"))
	writer.Close()
	req := httptest.NewRequest(http.MethodPost, "/upload", body)
	ctx := context.Background()
	ctx = context.WithValue(ctx, "UserId", int64(1))

	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	w := httptest.NewRecorder()

	h := Handler{service: &s}
	h.HandlerPostChunkCrateFile(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

}

func TestHandler_HandlerPostChunkCrateFilePart1BadUserId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	u := mock_service.NewMockUseCaser(ctrl)
	s := service2.Service{
		UseCaser: u,
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "test.txt")
	part.Write([]byte("test content"))
	writer.Close()
	req := httptest.NewRequest(http.MethodPost, "/upload", body)

	req.Header.Set("Content-Type", writer.FormDataContentType())

	w := httptest.NewRecorder()

	h := Handler{service: &s}
	h.HandlerPostChunkCrateFile(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

}

func TestHandler_HandlerPostChunkCrateFilePart1BadGetFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	u := mock_service.NewMockUseCaser(ctrl)
	s := service2.Service{
		UseCaser: u,
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "test.txt")
	part.Write([]byte("test content"))
	writer.Close()
	req := httptest.NewRequest(http.MethodPost, "/upload", nil)
	ctx := context.Background()
	ctx = context.WithValue(ctx, "UserId", int64(1))

	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	w := httptest.NewRecorder()

	h := Handler{service: &s}
	h.HandlerPostChunkCrateFile(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

}
func TestHandler_HandlerPostChunkCrateFilePart1BadGetFile2(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	u := mock_service.NewMockUseCaser(ctrl)
	s := service2.Service{
		UseCaser: u,
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "test.txt")
	part.Write([]byte(""))
	writer.Close()
	req := httptest.NewRequest(http.MethodPost, "/upload", body)
	ctx := context.Background()
	ctx = context.WithValue(ctx, "UserId", int64(1))

	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	w := httptest.NewRecorder()

	h := Handler{service: &s}
	h.HandlerPostChunkCrateFile(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

}
func TestHandler_HandlerPostChunkCrateFilePart2Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	u := mock_service.NewMockUseCaser(ctrl)
	s := service2.Service{
		UseCaser: u,
	}
	u.EXPECT().UploadFile(gomock.Any(), gomock.Any()).Return(true, &service2.TmpFile{}, nil)
	u.EXPECT().CreateFileChunks(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(&service2.RespData{}, nil)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "test.txt")
	assert.NoError(t, err)

	part.Write([]byte("test content"))

	reqData := &server.ReqData{
		Name:        "test.txt",
		Description: "desc",
		Data:        []byte("test content"),
	}

	reqDataBytes, err := json.Marshal(reqData)
	assert.NoError(t, err)
	part2, err := writer.CreateFormField("info")
	assert.NoError(t, err)
	part2.Write(reqDataBytes)

	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/upload", body)
	ctx := context.Background()
	ctx = context.WithValue(ctx, "UserId", int64(1))

	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	w := httptest.NewRecorder()

	h := Handler{service: &s}
	h.HandlerPostChunkCrateFile(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

}
func TestHandler_HandlerPostChunkCrateFilePart2BadInfo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	u := mock_service.NewMockUseCaser(ctrl)
	s := service2.Service{
		UseCaser: u,
	}
	u.EXPECT().UploadFile(gomock.Any(), gomock.Any()).Return(true, &service2.TmpFile{}, nil)
	//u.EXPECT().CreateFileChunks(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(&service2.RespData{}, nil)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "test.txt")
	assert.NoError(t, err)

	part.Write([]byte("test content"))

	part2, err := writer.CreateFormField("info")
	assert.NoError(t, err)
	part2.Write([]byte("not json"))

	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/upload", body)
	ctx := context.Background()
	ctx = context.WithValue(ctx, "UserId", int64(1))

	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	w := httptest.NewRecorder()

	h := Handler{service: &s}
	h.HandlerPostChunkCrateFile(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

}

func TestHandler_HandlerPostChunkCrateFilePart2BadCreateFileChunks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	u := mock_service.NewMockUseCaser(ctrl)
	s := service2.Service{
		UseCaser: u,
	}
	u.EXPECT().UploadFile(gomock.Any(), gomock.Any()).Return(true, &service2.TmpFile{}, nil)
	u.EXPECT().CreateFileChunks(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error"))

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "test.txt")
	assert.NoError(t, err)

	part.Write([]byte("test content"))

	reqData := &server.ReqData{
		Name:        "test.txt",
		Description: "desc",
		Data:        []byte("test content"),
	}

	reqDataBytes, err := json.Marshal(reqData)
	assert.NoError(t, err)
	part2, err := writer.CreateFormField("info")
	assert.NoError(t, err)
	part2.Write(reqDataBytes)

	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/upload", body)
	ctx := context.Background()
	ctx = context.WithValue(ctx, "UserId", int64(1))

	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	w := httptest.NewRecorder()

	h := Handler{service: &s}
	h.HandlerPostChunkCrateFile(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

}

// HandlerGetListData
func TestHandler_HandlerGetListData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	u := mock_service.NewMockUseCaser(ctrl)
	s := service2.Service{
		UseCaser: u,
	}

	u.EXPECT().GetListData(gomock.Any(), gomock.Any()).Return([]byte("data"), nil)

	req := httptest.NewRequest(http.MethodPost, "/GetListData", nil)
	ctx := context.Background()
	ctx = context.WithValue(ctx, "UserId", int64(1))

	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	h := Handler{service: &s}
	h.HandlerGetListData(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
func TestHandler_HandlerGetListDataBadUserId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	u := mock_service.NewMockUseCaser(ctrl)
	s := service2.Service{
		UseCaser: u,
	}

	req := httptest.NewRequest(http.MethodPost, "/GetListData", nil)

	w := httptest.NewRecorder()

	h := Handler{service: &s}
	h.HandlerGetListData(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
func TestHandler_HandlerGetListDataBadGetList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	u := mock_service.NewMockUseCaser(ctrl)
	s := service2.Service{
		UseCaser: u,
	}

	u.EXPECT().GetListData(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error"))

	req := httptest.NewRequest(http.MethodPost, "/GetListData", nil)
	ctx := context.Background()
	ctx = context.WithValue(ctx, "UserId", int64(1))

	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	h := Handler{service: &s}
	h.HandlerGetListData(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

//HandlerCheckUpdateData

func TestHandler_HandlerCheckUpdateData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	u := mock_service.NewMockUseCaser(ctrl)
	s := service2.Service{
		UseCaser: u,
	}

	u.EXPECT().ChangeData(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]byte("data"), nil)

	req := httptest.NewRequest(http.MethodPost, "/CheckUpdate", nil)
	ctx := context.Background()
	ctx = context.WithValue(ctx, "UserId", int64(1))

	req = req.WithContext(ctx)
	timee := time.Now()
	req.Header.Set("Last-Time-Update", timee.Format("2006-01-02 15:04:05.999999"))
	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("userDataId", "123")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
	w := httptest.NewRecorder()

	h := Handler{service: &s}
	h.HandlerCheckUpdateData(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
func TestHandler_HandlerCheckUpdateDataBadUserId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	u := mock_service.NewMockUseCaser(ctrl)
	s := service2.Service{
		UseCaser: u,
	}

	//u.EXPECT().ChangeData(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]byte("data"), nil)

	req := httptest.NewRequest(http.MethodPost, "/CheckUpdate", nil)

	timee := time.Now()
	req.Header.Set("Last-Time-Update", timee.Format("2006-01-02 15:04:05.999999"))
	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("userDataId", "123")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
	w := httptest.NewRecorder()

	h := Handler{service: &s}
	h.HandlerCheckUpdateData(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
func TestHandler_HandlerCheckUpdateDataBadTimeGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	u := mock_service.NewMockUseCaser(ctrl)
	s := service2.Service{
		UseCaser: u,
	}

	req := httptest.NewRequest(http.MethodPost, "/CheckUpdate", nil)
	ctx := context.Background()
	ctx = context.WithValue(ctx, "UserId", int64(1))

	req = req.WithContext(ctx)

	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("userDataId", "123")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
	w := httptest.NewRecorder()

	h := Handler{service: &s}
	h.HandlerCheckUpdateData(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
func TestHandler_HandlerCheckUpdateDataBadTimeParse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	u := mock_service.NewMockUseCaser(ctrl)
	s := service2.Service{
		UseCaser: u,
	}

	req := httptest.NewRequest(http.MethodPost, "/CheckUpdate", nil)
	ctx := context.Background()
	ctx = context.WithValue(ctx, "UserId", int64(1))

	req = req.WithContext(ctx)
	timee := time.Now()

	req.Header.Set("Last-Time-Update", timee.String())
	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("userDataId", "123")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
	w := httptest.NewRecorder()

	h := Handler{service: &s}
	h.HandlerCheckUpdateData(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
func TestHandler_HandlerCheckUpdateDataBadUserDataId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	u := mock_service.NewMockUseCaser(ctrl)
	s := service2.Service{
		UseCaser: u,
	}

	req := httptest.NewRequest(http.MethodPost, "/CheckUpdate", nil)
	ctx := context.Background()
	ctx = context.WithValue(ctx, "UserId", int64(1))

	req = req.WithContext(ctx)
	timee := time.Now()
	req.Header.Set("Last-Time-Update", timee.Format("2006-01-02 15:04:05.999999"))

	w := httptest.NewRecorder()

	h := Handler{service: &s}
	h.HandlerCheckUpdateData(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
func TestHandler_HandlerCheckUpdateDataBadChangeData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	u := mock_service.NewMockUseCaser(ctrl)
	s := service2.Service{
		UseCaser: u,
	}

	u.EXPECT().ChangeData(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error"))

	req := httptest.NewRequest(http.MethodPost, "/CheckUpdate", nil)
	ctx := context.Background()
	ctx = context.WithValue(ctx, "UserId", int64(1))

	req = req.WithContext(ctx)
	timee := time.Now()
	req.Header.Set("Last-Time-Update", timee.Format("2006-01-02 15:04:05.999999"))
	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("userDataId", "123")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
	w := httptest.NewRecorder()

	h := Handler{service: &s}
	h.HandlerCheckUpdateData(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

//HandlerCheckChanges

func TestHandler_HandlerCheckChanges(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	u := mock_service.NewMockUseCaser(ctrl)
	s := service2.Service{
		UseCaser: u,
	}

	u.EXPECT().ChangeAllData(gomock.Any(), gomock.Any(), gomock.Any()).Return([]byte("data"), nil)

	req := httptest.NewRequest(http.MethodPost, "/CheckUpdate", nil)
	ctx := context.Background()
	ctx = context.WithValue(ctx, "UserId", int64(1))

	req = req.WithContext(ctx)
	timee := time.Now()
	req.Header.Set("Last-Time-Update", timee.Format("2006-01-02 15:04:05.999999"))

	w := httptest.NewRecorder()

	h := Handler{service: &s}
	h.HandlerCheckChanges(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
func TestHandler_HandlerCheckChangesBadUserId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	u := mock_service.NewMockUseCaser(ctrl)
	s := service2.Service{
		UseCaser: u,
	}

	req := httptest.NewRequest(http.MethodPost, "/CheckUpdate", nil)

	timee := time.Now()
	req.Header.Set("Last-Time-Update", timee.Format("2006-01-02 15:04:05.999999"))

	w := httptest.NewRecorder()

	h := Handler{service: &s}
	h.HandlerCheckChanges(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
func TestHandler_HandlerCheckChangesBadGetTime(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	u := mock_service.NewMockUseCaser(ctrl)
	s := service2.Service{
		UseCaser: u,
	}

	req := httptest.NewRequest(http.MethodPost, "/CheckUpdate", nil)
	ctx := context.Background()
	ctx = context.WithValue(ctx, "UserId", int64(1))

	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	h := Handler{service: &s}
	h.HandlerCheckChanges(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
func TestHandler_HandlerCheckChangesBadParseTime(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	u := mock_service.NewMockUseCaser(ctrl)
	s := service2.Service{
		UseCaser: u,
	}

	req := httptest.NewRequest(http.MethodPost, "/CheckUpdate", nil)
	ctx := context.Background()
	ctx = context.WithValue(ctx, "UserId", int64(1))

	req = req.WithContext(ctx)
	timee := time.Now()
	req.Header.Set("Last-Time-Update", timee.String())

	w := httptest.NewRecorder()

	h := Handler{service: &s}
	h.HandlerCheckChanges(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
func TestHandler_HandlerCheckChangesErrorChangeAllData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	u := mock_service.NewMockUseCaser(ctrl)
	s := service2.Service{
		UseCaser: u,
	}

	u.EXPECT().ChangeAllData(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("Error"))

	req := httptest.NewRequest(http.MethodPost, "/CheckUpdate", nil)
	ctx := context.Background()
	ctx = context.WithValue(ctx, "UserId", int64(1))

	req = req.WithContext(ctx)
	timee := time.Now()
	req.Header.Set("Last-Time-Update", timee.Format("2006-01-02 15:04:05.999999"))

	w := httptest.NewRecorder()

	h := Handler{service: &s}
	h.HandlerCheckChanges(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

//HandlerGetFile

func TestHandler_HandlerGetFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	u := mock_service.NewMockUseCaser(ctrl)
	s := service2.Service{
		UseCaser: u,
	}

	u.EXPECT().GetFileChunks(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]byte("data"), nil)

	req := httptest.NewRequest(http.MethodPost, "/CheckUpdate", nil)
	ctx := context.Background()
	ctx = context.WithValue(ctx, "UserId", int64(1))

	req = req.WithContext(ctx)
	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("userDataId", "123")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
	w := httptest.NewRecorder()

	h := Handler{service: &s}
	h.HandlerGetFile(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
func TestHandler_HandlerGetFileErrGetUserId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	u := mock_service.NewMockUseCaser(ctrl)
	s := service2.Service{
		UseCaser: u,
	}

	req := httptest.NewRequest(http.MethodPost, "/CheckUpdate", nil)

	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("userDataId", "123")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
	w := httptest.NewRecorder()

	h := Handler{service: &s}
	h.HandlerGetFile(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
func TestHandler_HandlerGetFileGetUserDataId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	u := mock_service.NewMockUseCaser(ctrl)
	s := service2.Service{
		UseCaser: u,
	}

	req := httptest.NewRequest(http.MethodPost, "/CheckUpdate", nil)
	ctx := context.Background()
	ctx = context.WithValue(ctx, "UserId", int64(1))

	req = req.WithContext(ctx)
	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("userDataId", "--")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
	w := httptest.NewRecorder()

	h := Handler{service: &s}
	h.HandlerGetFile(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
func TestHandler_HandlerGetFileBadGetFileChunk(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	u := mock_service.NewMockUseCaser(ctrl)
	s := service2.Service{
		UseCaser: u,
	}

	u.EXPECT().GetFileChunks(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("Error"))

	req := httptest.NewRequest(http.MethodPost, "/CheckUpdate", nil)
	ctx := context.Background()
	ctx = context.WithValue(ctx, "UserId", int64(1))

	req = req.WithContext(ctx)
	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("userDataId", "123")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
	w := httptest.NewRecorder()

	h := Handler{service: &s}
	h.HandlerGetFile(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// HandlerGetData

func TestHandler_HandlerGetData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	u := mock_service.NewMockUseCaser(ctrl)
	s := service2.Service{
		UseCaser: u,
	}

	u.EXPECT().GetData(gomock.Any(), gomock.Any(), gomock.Any()).Return([]byte("data"), nil)

	req := httptest.NewRequest(http.MethodPost, "/CheckUpdate", nil)

	ctx := context.Background()
	ctx = context.WithValue(ctx, "UserId", int64(1))
	req = req.WithContext(ctx)

	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("userDataId", "123")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

	w := httptest.NewRecorder()

	h := Handler{service: &s}
	h.HandlerGetData(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
func TestHandler_HandlerGetDataBadUserId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	u := mock_service.NewMockUseCaser(ctrl)
	s := service2.Service{
		UseCaser: u,
	}

	req := httptest.NewRequest(http.MethodPost, "/CheckUpdate", nil)

	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("userDataId", "123")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
	w := httptest.NewRecorder()

	h := Handler{service: &s}
	h.HandlerGetData(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
func TestHandler_HandlerGetDataErrUserDataId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	u := mock_service.NewMockUseCaser(ctrl)
	s := service2.Service{
		UseCaser: u,
	}

	req := httptest.NewRequest(http.MethodPost, "/CheckUpdate", nil)
	ctx := context.Background()
	ctx = context.WithValue(ctx, "UserId", int64(1))

	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	h := Handler{service: &s}
	h.HandlerGetData(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
func TestHandler_HandlerGetDataErrGerData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	u := mock_service.NewMockUseCaser(ctrl)
	s := service2.Service{
		UseCaser: u,
	}

	u.EXPECT().GetData(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("Error"))

	req := httptest.NewRequest(http.MethodPost, "/CheckUpdate", nil)

	ctx := context.Background()
	ctx = context.WithValue(ctx, "UserId", int64(1))
	req = req.WithContext(ctx)

	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("userDataId", "123")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

	w := httptest.NewRecorder()

	h := Handler{service: &s}
	h.HandlerGetData(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

//HandlerGetFileSize

func TestHandler_HandlerGetFileSize(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	u := mock_service.NewMockUseCaser(ctrl)
	s := service2.Service{
		UseCaser: u,
	}

	u.EXPECT().GetFileSize(gomock.Any(), gomock.Any(), gomock.Any()).Return([]byte("data"), nil)

	req := httptest.NewRequest(http.MethodPost, "/CheckUpdate", nil)

	ctx := context.Background()
	ctx = context.WithValue(ctx, "UserId", int64(1))
	req = req.WithContext(ctx)

	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("userDataId", "123")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

	w := httptest.NewRecorder()

	h := Handler{service: &s}
	h.HandlerGetFileSize(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

}
func TestHandler_HandlerGetFileSizeBadUserId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	u := mock_service.NewMockUseCaser(ctrl)
	s := service2.Service{
		UseCaser: u,
	}

	req := httptest.NewRequest(http.MethodPost, "/CheckUpdate", nil)

	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("userDataId", "123")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

	w := httptest.NewRecorder()

	h := Handler{service: &s}
	h.HandlerGetFileSize(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

}
func TestHandler_HandlerGetFileSizeBadUserDataId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	u := mock_service.NewMockUseCaser(ctrl)
	s := service2.Service{
		UseCaser: u,
	}

	req := httptest.NewRequest(http.MethodPost, "/CheckUpdate", nil)

	ctx := context.Background()
	ctx = context.WithValue(ctx, "UserId", int64(1))
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	h := Handler{service: &s}
	h.HandlerGetFileSize(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

}
func TestHandler_HandlerGetFileSizeBadUserDataId0(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	u := mock_service.NewMockUseCaser(ctrl)
	s := service2.Service{
		UseCaser: u,
	}

	req := httptest.NewRequest(http.MethodPost, "/CheckUpdate", nil)

	ctx := context.Background()
	ctx = context.WithValue(ctx, "UserId", int64(1))
	req = req.WithContext(ctx)

	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("userDataId", "0")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

	w := httptest.NewRecorder()

	h := Handler{service: &s}
	h.HandlerGetFileSize(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

}
func TestHandler_HandlerGetFileSizeBadGetFileSize(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	u := mock_service.NewMockUseCaser(ctrl)
	s := service2.Service{
		UseCaser: u,
	}

	u.EXPECT().GetFileSize(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error"))

	req := httptest.NewRequest(http.MethodPost, "/CheckUpdate", nil)

	ctx := context.Background()
	ctx = context.WithValue(ctx, "UserId", int64(1))
	req = req.WithContext(ctx)

	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("userDataId", "123")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

	w := httptest.NewRecorder()

	h := Handler{service: &s}
	h.HandlerGetFileSize(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

}

func TestHandler_HandlerUpdateData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	u := mock_service.NewMockUseCaser(ctrl)
	s := service2.Service{
		UseCaser: u,
	}

	u.EXPECT().UpdateData(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]byte("data"), nil)

	body := []byte(`{"Name": "data"}`)

	req := httptest.NewRequest(http.MethodPost, "/CheckUpdate", bytes.NewBuffer(body))

	ctx := context.Background()
	ctx = context.WithValue(ctx, "UserId", int64(1))
	req = req.WithContext(ctx)

	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("userDataId", "123")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

	w := httptest.NewRecorder()

	h := Handler{service: &s}
	h.HandlerUpdateData(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

}
func TestHandler_HandlerUpdateDataBadUserId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	u := mock_service.NewMockUseCaser(ctrl)
	s := service2.Service{
		UseCaser: u,
	}

	body := []byte(`{"Name": "data"}`)

	req := httptest.NewRequest(http.MethodPost, "/CheckUpdate", bytes.NewBuffer(body))

	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("userDataId", "123")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

	w := httptest.NewRecorder()

	h := Handler{service: &s}
	h.HandlerUpdateData(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

}
func TestHandler_HandlerUpdateDataBadUserDataId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	u := mock_service.NewMockUseCaser(ctrl)
	s := service2.Service{
		UseCaser: u,
	}

	body := []byte(`{"Name": "data"}`)

	req := httptest.NewRequest(http.MethodPost, "/CheckUpdate", bytes.NewBuffer(body))

	ctx := context.Background()
	ctx = context.WithValue(ctx, "UserId", int64(1))
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	h := Handler{service: &s}
	h.HandlerUpdateData(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

}
func TestHandler_HandlerUpdateDataBadBody(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	u := mock_service.NewMockUseCaser(ctrl)
	s := service2.Service{
		UseCaser: u,
	}

	body := []byte(`-`)

	req := httptest.NewRequest(http.MethodPost, "/CheckUpdate", bytes.NewBuffer(body))

	ctx := context.Background()
	ctx = context.WithValue(ctx, "UserId", int64(1))
	req = req.WithContext(ctx)

	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("userDataId", "123")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

	w := httptest.NewRecorder()

	h := Handler{service: &s}
	h.HandlerUpdateData(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

}
func TestHandler_HandlerUpdateDataBadUpdateData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	u := mock_service.NewMockUseCaser(ctrl)
	s := service2.Service{
		UseCaser: u,
	}

	u.EXPECT().UpdateData(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error"))

	body := []byte(`{"Name": "data"}`)

	req := httptest.NewRequest(http.MethodPost, "/CheckUpdate", bytes.NewBuffer(body))

	ctx := context.Background()
	ctx = context.WithValue(ctx, "UserId", int64(1))
	req = req.WithContext(ctx)

	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("userDataId", "123")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

	w := httptest.NewRecorder()

	h := Handler{service: &s}
	h.HandlerUpdateData(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

}

func TestHandler_HandlerRemoveData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	u := mock_service.NewMockUseCaser(ctrl)
	s := service2.Service{
		UseCaser: u,
	}

	u.EXPECT().RemoveData(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	req := httptest.NewRequest(http.MethodPost, "/CheckUpdate", nil)

	ctx := context.Background()
	ctx = context.WithValue(ctx, "UserId", int64(1))
	req = req.WithContext(ctx)

	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("userDataId", "123")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

	w := httptest.NewRecorder()

	h := Handler{service: &s}
	h.HandlerRemoveData(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

}
func TestHandler_HandlerRemoveDataBadUserId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	u := mock_service.NewMockUseCaser(ctrl)
	s := service2.Service{
		UseCaser: u,
	}

	req := httptest.NewRequest(http.MethodPost, "/CheckUpdate", nil)

	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("userDataId", "123")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

	w := httptest.NewRecorder()

	h := Handler{service: &s}
	h.HandlerRemoveData(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

}

func TestHandler_HandlerRemoveDataBadUserDataId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	u := mock_service.NewMockUseCaser(ctrl)
	s := service2.Service{
		UseCaser: u,
	}

	req := httptest.NewRequest(http.MethodPost, "/CheckUpdate", nil)

	ctx := context.Background()
	ctx = context.WithValue(ctx, "UserId", int64(1))
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	h := Handler{service: &s}
	h.HandlerRemoveData(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

}
func TestHandler_HandlerRemoveDataBadRemoveData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	u := mock_service.NewMockUseCaser(ctrl)
	s := service2.Service{
		UseCaser: u,
	}

	u.EXPECT().RemoveData(gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("error"))

	req := httptest.NewRequest(http.MethodPost, "/CheckUpdate", nil)

	ctx := context.Background()
	ctx = context.WithValue(ctx, "UserId", int64(1))
	req = req.WithContext(ctx)

	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("userDataId", "123")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

	w := httptest.NewRecorder()

	h := Handler{service: &s}
	h.HandlerRemoveData(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

}

func TestHandler_HandlerUpdateBinaryFilePart1Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	u := mock_service.NewMockUseCaser(ctrl)
	s := service2.Service{
		UseCaser: u,
	}
	u.EXPECT().UploadFile(gomock.Any(), gomock.Any()).Return(false, &service2.TmpFile{}, nil)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "test.txt")
	part.Write([]byte("test content"))
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/upload", body)

	ctx := context.Background()
	ctx = context.WithValue(ctx, "UserId", int64(1))
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", writer.FormDataContentType())

	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("userDataId", "123")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

	w := httptest.NewRecorder()

	h := Handler{service: &s}
	h.HandlerUpdateBinaryFile(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

}
func TestHandler_HandlerUpdateBinaryFilePart1BadUserId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	u := mock_service.NewMockUseCaser(ctrl)
	s := service2.Service{
		UseCaser: u,
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "test.txt")
	part.Write([]byte("test content"))
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/upload", body)

	req.Header.Set("Content-Type", writer.FormDataContentType())

	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("userDataId", "123")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

	w := httptest.NewRecorder()

	h := Handler{service: &s}
	h.HandlerUpdateBinaryFile(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

}
func TestHandler_HandlerUpdateBinaryFilePart1SuccessBadUserDataId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	u := mock_service.NewMockUseCaser(ctrl)
	s := service2.Service{
		UseCaser: u,
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "test.txt")
	part.Write([]byte("test content"))
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/upload", body)

	ctx := context.Background()
	ctx = context.WithValue(ctx, "UserId", int64(1))
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", writer.FormDataContentType())

	w := httptest.NewRecorder()

	h := Handler{service: &s}
	h.HandlerUpdateBinaryFile(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

}
func TestHandler_HandlerUpdateBinaryFilePart1GetFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	u := mock_service.NewMockUseCaser(ctrl)
	s := service2.Service{
		UseCaser: u,
	}

	req := httptest.NewRequest(http.MethodPost, "/upload", nil)

	ctx := context.Background()
	ctx = context.WithValue(ctx, "UserId", int64(1))
	req = req.WithContext(ctx)

	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("userDataId", "123")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

	w := httptest.NewRecorder()

	h := Handler{service: &s}
	h.HandlerUpdateBinaryFile(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

}
func TestHandler_HandlerUpdateBinaryFilePart1FileIsNull(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	u := mock_service.NewMockUseCaser(ctrl)
	s := service2.Service{
		UseCaser: u,
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "test.txt")
	part.Write([]byte(""))
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/upload", body)

	ctx := context.Background()
	ctx = context.WithValue(ctx, "UserId", int64(1))
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", writer.FormDataContentType())

	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("userDataId", "123")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

	w := httptest.NewRecorder()

	h := Handler{service: &s}
	h.HandlerUpdateBinaryFile(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

}
func TestHandler_HandlerUpdateBinaryFilePart1BadUploadFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	u := mock_service.NewMockUseCaser(ctrl)
	s := service2.Service{
		UseCaser: u,
	}
	u.EXPECT().UploadFile(gomock.Any(), gomock.Any()).Return(false, nil, fmt.Errorf("error"))

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "test.txt")
	part.Write([]byte("test content"))
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/upload", body)

	ctx := context.Background()
	ctx = context.WithValue(ctx, "UserId", int64(1))
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", writer.FormDataContentType())

	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("userDataId", "123")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

	w := httptest.NewRecorder()

	h := Handler{service: &s}
	h.HandlerUpdateBinaryFile(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

}
func TestHandler_HandlerUpdateBinaryFilePart1BadUploadFileTmpIsNill(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	u := mock_service.NewMockUseCaser(ctrl)
	s := service2.Service{
		UseCaser: u,
	}
	u.EXPECT().UploadFile(gomock.Any(), gomock.Any()).Return(false, nil, nil)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "test.txt")
	part.Write([]byte("test content"))
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/upload", body)

	ctx := context.Background()
	ctx = context.WithValue(ctx, "UserId", int64(1))
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", writer.FormDataContentType())

	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("userDataId", "123")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

	w := httptest.NewRecorder()

	h := Handler{service: &s}
	h.HandlerUpdateBinaryFile(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

}
func TestHandler_HandlerUpdateBinaryFilePart2Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	u := mock_service.NewMockUseCaser(ctrl)
	s := service2.Service{
		UseCaser: u,
	}
	u.EXPECT().UploadFile(gomock.Any(), gomock.Any()).Return(true, &service2.TmpFile{}, nil)
	u.EXPECT().UpdateBinaryFile(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(&service2.RespData{}, nil)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "test.txt")
	assert.NoError(t, err)

	part.Write([]byte("test content"))

	reqData := &server.ReqData{
		Name:        "test.txt",
		Description: "desc",
		Data:        []byte("test content"),
	}

	reqDataBytes, err := json.Marshal(reqData)
	assert.NoError(t, err)
	part2, err := writer.CreateFormField("info")
	assert.NoError(t, err)
	part2.Write(reqDataBytes)

	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/upload", body)

	ctx := context.Background()
	ctx = context.WithValue(ctx, "UserId", int64(1))
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", writer.FormDataContentType())

	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("userDataId", "123")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

	w := httptest.NewRecorder()

	h := Handler{service: &s}
	h.HandlerUpdateBinaryFile(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

}
func TestHandler_HandlerUpdateBinaryFilePart2BadGetInfo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	u := mock_service.NewMockUseCaser(ctrl)
	s := service2.Service{
		UseCaser: u,
	}
	u.EXPECT().UploadFile(gomock.Any(), gomock.Any()).Return(true, &service2.TmpFile{}, nil)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "test.txt")
	assert.NoError(t, err)

	part.Write([]byte("test content"))

	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/upload", body)

	ctx := context.Background()
	ctx = context.WithValue(ctx, "UserId", int64(1))
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", writer.FormDataContentType())

	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("userDataId", "123")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

	w := httptest.NewRecorder()

	h := Handler{service: &s}
	h.HandlerUpdateBinaryFile(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

}
func TestHandler_HandlerUpdateBinaryFilePart2BadUpdateBinaryFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	u := mock_service.NewMockUseCaser(ctrl)
	s := service2.Service{
		UseCaser: u,
	}
	u.EXPECT().UploadFile(gomock.Any(), gomock.Any()).Return(true, &service2.TmpFile{}, nil)
	u.EXPECT().UpdateBinaryFile(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error"))

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "test.txt")
	assert.NoError(t, err)

	part.Write([]byte("test content"))

	reqData := &server.ReqData{
		Name:        "test.txt",
		Description: "desc",
		Data:        []byte("test content"),
	}

	reqDataBytes, err := json.Marshal(reqData)
	assert.NoError(t, err)
	part2, err := writer.CreateFormField("info")
	assert.NoError(t, err)
	part2.Write(reqDataBytes)

	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/upload", body)

	ctx := context.Background()
	ctx = context.WithValue(ctx, "UserId", int64(1))
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", writer.FormDataContentType())

	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("userDataId", "123")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

	w := httptest.NewRecorder()

	h := Handler{service: &s}
	h.HandlerUpdateBinaryFile(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

}

func TestHandler_HandlerPostCrateFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	u := mock_service.NewMockUseCaser(ctrl)
	s := service2.Service{
		UseCaser: u,
	}
	u.EXPECT().CreateFile(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(&service2.RespData{}, nil)
	reqData := ReqData{
		Name:        "Name",
		Description: "Desc",
		Data:        []byte("data"),
	}

	reqBody, err := json.Marshal(reqData)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/postCreateFile", strings.NewReader(string(reqBody)))

	ctx := context.Background()
	ctx = context.WithValue(ctx, "UserId", int64(1))
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	h := Handler{service: &s}
	h.HandlerPostCrateFile(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
func TestHandler_HandlerPostCrateFileBadUserId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	u := mock_service.NewMockUseCaser(ctrl)
	s := service2.Service{
		UseCaser: u,
	}
	reqData := ReqData{
		Name:        "Name",
		Description: "Desc",
		Data:        []byte("data"),
	}

	reqBody, err := json.Marshal(reqData)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/postCreateFile", strings.NewReader(string(reqBody)))

	w := httptest.NewRecorder()

	h := Handler{service: &s}
	h.HandlerPostCrateFile(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
func TestHandler_HandlerPostCrateFileBadDecodeBody(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	u := mock_service.NewMockUseCaser(ctrl)
	s := service2.Service{
		UseCaser: u,
	}

	req := httptest.NewRequest(http.MethodPost, "/postCreateFile", strings.NewReader(string("---")))

	ctx := context.Background()
	ctx = context.WithValue(ctx, "UserId", int64(1))
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	h := Handler{service: &s}
	h.HandlerPostCrateFile(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandler_HandlerPostCrateFileBadCreateFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	u := mock_service.NewMockUseCaser(ctrl)
	s := service2.Service{
		UseCaser: u,
	}
	u.EXPECT().CreateFile(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error"))
	reqData := ReqData{
		Name:        "Name",
		Description: "Desc",
		Data:        []byte("data"),
	}

	reqBody, err := json.Marshal(reqData)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/postCreateFile", strings.NewReader(string(reqBody)))

	ctx := context.Background()
	ctx = context.WithValue(ctx, "UserId", int64(1))
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	h := Handler{service: &s}
	h.HandlerPostCrateFile(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
