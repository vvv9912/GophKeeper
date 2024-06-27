package middleware

import (
	"GophKeeper/internal/Server/service"
	mock_service "GophKeeper/internal/Server/service/mocks"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMw_MiddlewareAuthBearerPrefix(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAuth := mock_service.NewMockAuth(ctrl)

	s := &service.Service{
		UseCaser: nil,
		Auth:     mockAuth,
	}

	mw := Mw{s.Auth}

	validToken := "validToken"
	userId := int64(123)

	mockAuth.EXPECT().GetUserId(validToken).Return(userId, nil)

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "BeArEr "+validToken)

	rr := httptest.NewRecorder()
	handler := mw.MiddlewareAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctxUserId := r.Context().Value("UserId").(int64)
		if ctxUserId != userId {
			t.Errorf("expected user ID %d, got %d", userId, ctxUserId)
		}
		w.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}
func TestMw_MiddlewareAuthBearerPrefixBadToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAuth := mock_service.NewMockAuth(ctrl)

	s := &service.Service{
		UseCaser: nil,
		Auth:     mockAuth,
	}

	mw := Mw{s.Auth}

	userId := int64(123)

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := mw.MiddlewareAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctxUserId := r.Context().Value("UserId").(int64)
		if ctxUserId != userId {
			t.Errorf("expected user ID %d, got %d", userId, ctxUserId)
		}
		w.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}
func TestMw_MiddlewareAuthBearerPrefixGetUserId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAuth := mock_service.NewMockAuth(ctrl)

	s := &service.Service{
		UseCaser: nil,
		Auth:     mockAuth,
	}

	mw := Mw{s.Auth}

	validToken := "validToken"
	userId := int64(123)

	mockAuth.EXPECT().GetUserId(validToken).Return(int64(0), fmt.Errorf("error"))

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "BeArEr "+validToken)

	rr := httptest.NewRecorder()
	handler := mw.MiddlewareAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctxUserId := r.Context().Value("UserId").(int64)
		if ctxUserId != userId {
			t.Errorf("expected user ID %d, got %d", userId, ctxUserId)
		}
		w.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}
func TestMw_MiddlewareAuthBearerPrefixGetUserId2(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAuth := mock_service.NewMockAuth(ctrl)

	s := &service.Service{
		UseCaser: nil,
		Auth:     mockAuth,
	}

	mw := Mw{s.Auth}

	validToken := "validToken"
	userId := int64(123)

	mockAuth.EXPECT().GetUserId(validToken).Return(int64(-1), nil)

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "BeArEr "+validToken)

	rr := httptest.NewRecorder()
	handler := mw.MiddlewareAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctxUserId := r.Context().Value("UserId").(int64)
		if ctxUserId != userId {
			t.Errorf("expected user ID %d, got %d", userId, ctxUserId)
		}
		w.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}
