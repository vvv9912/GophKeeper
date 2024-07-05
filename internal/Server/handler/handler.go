package handler

import (
	"GophKeeper/internal/Server/middleware"
	"GophKeeper/internal/Server/service"
	"github.com/go-chi/chi/v5"
	"net/http"
)

// Handler - обработчик HTTP-запросов.
type Handler struct {
	service *service.Service // интерфейс сервиса.
}

// NewHandler - конструктор обработчика HTTP-запросов.
func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

// InitRoutes - инициализация маршрутов.
func (h *Handler) InitRoutes() http.Handler {

	r := chi.NewRouter()
	mw := middleware.NewMw(h.service.Auth)

	apiR := r.Route("/api", func(r chi.Router) {})

	apiR.Post("/ping", HandlerPing)
	// /api/signIn
	apiR.Post("/signIn", h.HandlerSignIn)
	// /api/signIn
	apiR.Post("/signUp", h.HandlerSignUp)

	// /api/data
	apiR.Route("/data", func(r chi.Router) {
		r.Use(mw.MiddlewareAuth)
		//_ = mw
		r.Get("/", h.HandlerGetListData)

		r.Post("/credentials", h.HandlerPostCredentials)
		r.Post("/file", h.HandlerPostCrateFile)
		r.Post("/fileChunks", h.HandlerPostChunkCrateFile)
		r.Post("/creditCard", h.HandlerPostCreditCard)

		r.Get("/changes", h.HandlerCheckChanges)
		r.Get("/{userDataId:[0-9]+}", h.HandlerGetData)
		r.Get("/fileSize/{userDataId:[0-9]+}", h.HandlerGetFileSize)
		r.Get("/fileChunks/{userDataId:[0-9]+}", h.HandlerGetFile)
		r.Post("/CheckUpdate/{userDataId:[0-9]+}", h.HandlerCheckUpdateData)
		r.Post("/update/{userDataId:[0-9]+}", h.HandlerUpdateData)
		r.Post("/updateBinary/{userDataId:[0-9]+}", h.HandlerUpdateBinaryFile)

		r.Patch("/", h.HandlerUpdateData)
		r.Delete("/{userDataId:[0-9]+}", h.HandlerRemoveData)
	})

	return r
}
