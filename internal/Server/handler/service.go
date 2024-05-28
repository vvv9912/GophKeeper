package handler

import (
	"GophKeeper/internal/Server/service"
	"fmt"
	"net/http"
)

type Handler struct {
	service service.Service
}

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
