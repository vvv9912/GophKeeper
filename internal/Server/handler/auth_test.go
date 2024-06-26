package handler

import (
	service2 "GophKeeper/internal/Server/service"
	mock_service "GophKeeper/internal/Server/service/mocks"
	"encoding/json"
	"fmt"
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
func TestHandler_HandlerSignUpBadAnswer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	u := mock_service.NewMockUseCaser(ctrl)
	s := service2.Service{
		UseCaser: u,
	}
	u.EXPECT().SignUp(gomock.Any(), "testuser", "testpass").Return("", fmt.Errorf("error"))
	reqBody := `{"login": "testuser", "password": "testpass"}`
	req := httptest.NewRequest(http.MethodPost, "/signup", strings.NewReader(reqBody))
	w := httptest.NewRecorder()

	h := Handler{service: &s}
	h.HandlerSignUp(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

}
func TestHandler_HandlerSignUpErrBadBody(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	u := mock_service.NewMockUseCaser(ctrl)
	s := service2.Service{
		UseCaser: u,
	}

	reqBody := `sword": "testpass"}`
	req := httptest.NewRequest(http.MethodPost, "/signup", strings.NewReader(reqBody))
	w := httptest.NewRecorder()

	h := Handler{service: &s}
	h.HandlerSignUp(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

}

func TestHandler_HandlerSignIn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	u := mock_service.NewMockUseCaser(ctrl)
	s := service2.Service{
		UseCaser: u,
	}
	u.EXPECT().SignIn(gomock.Any(), "testuser", "testpass").Return("mockJWT", nil)
	reqBody := `{"login": "testuser", "password": "testpass"}`
	req := httptest.NewRequest(http.MethodPost, "/signup", strings.NewReader(reqBody))
	w := httptest.NewRecorder()

	h := Handler{service: &s}
	h.HandlerSignIn(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var user User
	err := json.NewDecoder(w.Body).Decode(&user)
	assert.NoError(t, err)
	assert.Equal(t, "testuser", user.Login)
	assert.Equal(t, "mockJWT", user.JWT)
}
func TestHandler_HandlerSignInErrBadBody(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	u := mock_service.NewMockUseCaser(ctrl)
	s := service2.Service{
		UseCaser: u,
	}

	reqBody := `sword": "testpass"}`
	req := httptest.NewRequest(http.MethodPost, "/signIn", strings.NewReader(reqBody))
	w := httptest.NewRecorder()

	h := Handler{service: &s}
	h.HandlerSignIn(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

}
func TestHandler_HandlerSignInBadAnswer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	u := mock_service.NewMockUseCaser(ctrl)
	s := service2.Service{
		UseCaser: u,
	}
	u.EXPECT().SignIn(gomock.Any(), "testuser", "testpass").Return("", fmt.Errorf("error"))
	reqBody := `{"login": "testuser", "password": "testpass"}`
	req := httptest.NewRequest(http.MethodPost, "/signup", strings.NewReader(reqBody))
	w := httptest.NewRecorder()

	h := Handler{service: &s}
	h.HandlerSignIn(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

}

func TestHandlerPing(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/ping", nil)
	w := httptest.NewRecorder()
	HandlerPing(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
