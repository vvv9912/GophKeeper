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
