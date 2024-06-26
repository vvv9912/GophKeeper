package handler

import (
	service2 "GophKeeper/internal/Server/service"
	mock_service "GophKeeper/internal/Server/service/mocks"
	"context"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
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
