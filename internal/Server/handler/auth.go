package handler

import (
	"GophKeeper/pkg/customErrors"
	"GophKeeper/pkg/logger"
	"encoding/json"
	"errors"
	"go.uber.org/zap"
	"net/http"
)

// deferHandler - defer function for get error
func deferHandler(err error, w http.ResponseWriter) {
	if err != nil {
		return
	}

	var MyErr *customErrors.CustomError
	if errors.As(err, &MyErr) {

		if MyErr.StatusCode != 0 {

			w.WriteHeader(MyErr.StatusCode)

			err := json.NewEncoder(w).Encode(RespError{
				Code:    MyErr.StatusCode,
				Message: MyErr.Message,
			})
			if err != nil {
				logger.Log.Error("Error writing response", zap.Error(err))
			}

			return
		}

	}

	w.WriteHeader(http.StatusInternalServerError)

	err = json.NewEncoder(w).Encode(RespError{
		Code:    http.StatusInternalServerError,
		Message: err.Error(),
	})
	if err != nil {
		logger.Log.Error("Error writing response", zap.Error(err))
	}

}

// HandlerSignUp - регистрация пользователя
func (h Handler) HandlerSignUp(w http.ResponseWriter, r *http.Request) {
	var err error
	var resp []byte

	defer func() {
		if err != nil {
			deferHandler(err, w)
		} else {
			w.Write(resp)
			w.WriteHeader(http.StatusOK)
		}
	}()

	var auth Auth

	err = json.NewDecoder(r.Body).Decode(&auth)
	if err != nil {
		logger.Log.Error("Unmarshal json failed", zap.Error(err))
		err = customErrors.NewCustomError(err, http.StatusBadRequest, "Error reading request body")
		return
	}

	jwt, err := h.service.SignUp(r.Context(), auth.Login, auth.Password)
	if err != nil {
		return
	}

	user := User{
		Login: auth.Login,
		JWT:   jwt,
	}

	resp, err = json.Marshal(user)
	if err != nil {
		logger.Log.Error("Marshal json failed", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
	}

}

// HandlerSignIn - авторизация пользователя
func (h Handler) HandlerSignIn(w http.ResponseWriter, r *http.Request) {
	var err error
	var resp []byte

	defer func() {
		if err != nil {
			deferHandler(err, w)
		} else {
			w.Write(resp)
			w.WriteHeader(http.StatusOK)
		}
	}()

	var auth Auth

	err = json.NewDecoder(r.Body).Decode(&auth)
	if err != nil {
		logger.Log.Error("Unmarshal json failed", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
	}

	jwt, err := h.service.SignIn(r.Context(), auth.Login, auth.Password)
	if err != nil {
		return
	}

	user := User{
		Login: auth.Login,
		JWT:   jwt,
	}

	resp, err = json.Marshal(user)
	if err != nil {
		logger.Log.Error("Marshal json failed", zap.Error(err))
		err = customErrors.NewCustomError(err, http.StatusInternalServerError, "Error create response")
		return
	}

}
