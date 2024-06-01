package handler

import (
	"GophKeeper/pkg/customErrors"
	"GophKeeper/pkg/logger"
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func (h *Handler) HandlerPostCredentials(w http.ResponseWriter, r *http.Request) {
	var err error
	var resp []byte

	defer func() {
		if err != nil {
			deferHandler(err, w)
		} else {
			w.WriteHeader(http.StatusOK)
			_, err = w.Write(resp)
			if err != nil {
				logger.Log.Error("Error send resp", zap.Error(err))
			}
		}
	}()

	userId, err := getUserId(r)
	if err != nil {
		return
	}

	var Cred ReqData

	err = json.NewDecoder(r.Body).Decode(&Cred)
	if err != nil {
		logger.Log.Error("Unmarshal json failed", zap.Error(err))
		err = customErrors.NewCustomError(err, http.StatusBadRequest, "Error reading request body")
		return
	}

	err = h.service.Data.CreateCredentials(r.Context(), userId, Cred.Data, Cred.Name, Cred.Description)
	if err != nil {
		return
	}

	resp = []byte("Credential successfully created")

}

func (h *Handler) HandlerPostCreditCard(w http.ResponseWriter, r *http.Request) {
	var err error
	var resp []byte

	defer func() {
		if err != nil {
			deferHandler(err, w)
		} else {
			w.WriteHeader(http.StatusOK)
			_, err = w.Write(resp)
			if err != nil {
				logger.Log.Error("Error send resp", zap.Error(err))
			}
		}
	}()

	userId, err := getUserId(r)
	if err != nil {
		return
	}

	var Cred ReqData

	err = json.NewDecoder(r.Body).Decode(&Cred)
	if err != nil {
		logger.Log.Error("Unmarshal json failed", zap.Error(err))
		err = customErrors.NewCustomError(err, http.StatusBadRequest, "Error reading request body")
		return
	}

	err = h.service.Data.CreateCreditCard(r.Context(), userId, Cred.Data, Cred.Name, Cred.Description)
	if err != nil {
		return
	}

	resp = []byte("Credential successfully created")

}

func (h *Handler) HandlerPostCrateFile(w http.ResponseWriter, r *http.Request) {
	var err error
	var resp []byte

	defer func() {
		if err != nil {
			deferHandler(err, w)
		} else {
			w.WriteHeader(http.StatusOK)
			_, err = w.Write(resp)
			if err != nil {
				logger.Log.Error("Error send resp", zap.Error(err))
			}
		}
	}()

	userId, err := getUserId(r)
	if err != nil {
		return
	}

	var Cred ReqData

	err = json.NewDecoder(r.Body).Decode(&Cred)
	if err != nil {
		logger.Log.Error("Unmarshal json failed", zap.Error(err))
		err = customErrors.NewCustomError(err, http.StatusBadRequest, "Error reading request body")
		return
	}

	err = h.service.Data.CreateFile(r.Context(), userId, Cred.Data, Cred.Name, Cred.Description)
	if err != nil {
		return
	}

	resp = []byte("Credential successfully created")

}

func (h *Handler) HandlerCheckChanges(w http.ResponseWriter, r *http.Request) {
	var err error
	var resp []byte

	defer func() {
		if err != nil {
			deferHandler(err, w)
		} else {
			w.WriteHeader(http.StatusOK)
			_, err = w.Write(resp)
			if err != nil {
				logger.Log.Error("Error send resp", zap.Error(err))
			}
		}
	}()
	userId, err := getUserId(r)
	if err != nil {
		return
	}
	strLastTime := r.Header.Get("Last-Time-Update")
	if strLastTime == "" {
		err = customErrors.NewCustomError(nil, http.StatusBadRequest, "Last-Time-Update header is required")
		logger.Log.Error("Last-Time-Update header is required", zap.Error(err))
		return
	}

	LastTimeUpdate, err := time.Parse("2006-01-02 15:04:05.999999", strLastTime)
	if err != nil {
		err = customErrors.NewCustomError(err, http.StatusBadRequest, "Last-Time-Update header is invalid")
		logger.Log.Error("Last-Time-Update header is invalid", zap.Error(err))
	}

	resp, err = h.service.Data.ChangeData(r.Context(), userId, LastTimeUpdate)
}
