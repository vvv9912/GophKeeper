package handler

import (
	"GophKeeper/internal/Server/middleware"
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

func (h *Handler) InitRoutes(services *service.Service) http.Handler {

	r := chi.NewRouter()
	mw := middleware.Mw{services.Auth}

	apiR := r.Route("/api", func(r chi.Router) {})
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

		r.Patch("/", h.HandlerUpdateData)
		r.Delete("/{userDataId:[0-9]+}", h.HandlerRemoveData)
	})

	return r
}
