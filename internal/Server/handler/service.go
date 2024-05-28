package handler

import (
	"GophKeeper/internal/Server/service"
	"GophKeeper/pkg/customErrors"
	"GophKeeper/pkg/logger"
	"encoding/json"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"net/http"
)

type Handler struct {
	service service.Service
}

// getUserId - получение id пользователя из контекста request
func getUserId(r *http.Request) (userId int64, err error) {
	value := r.Context().Value("UserId")

	if value == nil {
		err := fmt.Errorf("UserId is empty")
		return 0, err
	}

	userId, ok := value.(int64)
	if !ok {
		err := fmt.Errorf("UserId is not int64")
		return 0, err
	}

	return userId, nil
}

// deferHandler - defer function for get error
func deferHandler(err error, w http.ResponseWriter) {

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
