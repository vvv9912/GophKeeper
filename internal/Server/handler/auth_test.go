package handler

import (
	service2 "GophKeeper/internal/Server/service"
	mock_service "GophKeeper/internal/Server/service/mocks"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// todo
func TestHandler_HandlerSignUp(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	u := mock_service.NewMockUseCaser(ctrl)
	s := service2.Service{
		UseCaser: u,
	}
	u.EXPECT().SignUp(gomock.Any(), "testuser", "testpass").Return("mockJWT", nil)
	reqBody := `{"login": "testuser", "password": "testpass"}`
	req := httptest.NewRequest(http.MethodPost, "/signup", strings.NewReader(reqBody))
	w := httptest.NewRecorder()

	h := Handler{service: &s}
	h.HandlerSignUp(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var user User
	err := json.NewDecoder(w.Body).Decode(&user)
	assert.NoError(t, err)
	assert.Equal(t, "testuser", user.Login)
	assert.Equal(t, "mockJWT", user.JWT)
}
func TestHandler_HandlerSignUpErrBadBody(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	u := mock_service.NewMockUseCaser(ctrl)
	s := service2.Service{
		UseCaser: u,
	}
	u.EXPECT().SignUp(gomock.Any(), "testuser", "testpass").Return("mockJWT", nil)
	reqBody := `sword": "testpass"}`
	req := httptest.NewRequest(http.MethodPost, "/signup", strings.NewReader(reqBody))
	w := httptest.NewRecorder()

	h := Handler{service: &s}
	h.HandlerSignUp(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var user User
	err := json.NewDecoder(w.Body).Decode(&user)
	assert.Error(t, err)

}
