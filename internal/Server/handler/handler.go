package handler

import (
	"GophKeeper/internal/Server/service"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) InitRoutes() http.Handler {

	r := chi.NewRouter()

	r.Post("/", h.HandlerPostCredentials)
	r.Get("/changes", h.HandlerCheckChanges)
	return r
}

// getUserId - получение id пользователя из контекста request
func getUserId(r *http.Request) (userId int64, err error) {
	//value := r.Context().Value("UserId")
	//
	//if value == nil {
	//	err := fmt.Errorf("UserId is empty")
	//	return 0, err
	//}
	//
	//userId, ok := value.(int64)
	//if !ok {
	//	err := fmt.Errorf("UserId is not int64")
	//	return 0, err
	//}
	userId = 1
	return userId, nil
}
