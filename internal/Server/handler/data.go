package handler

import (
	"GophKeeper/pkg/customErrors"
	"GophKeeper/pkg/logger"
	"GophKeeper/pkg/store"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"time"
)

func (h *Handler) HandlerPostCredentials(w http.ResponseWriter, r *http.Request) {
	var err error
	var resp []byte

	defer func() {
		if err != nil {
			deferHandler(err, w)
			return
		}
		_, err = w.Write(resp)
		if err != nil {
			logger.Log.Error("Error writing response", zap.Error(err))
		}
		w.WriteHeader(http.StatusOK)

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

	response, err := h.service.Data.CreateCredentials(r.Context(), userId, Cred.Data, Cred.Name, Cred.Description)
	if err != nil {
		return
	}

	resp, err = json.Marshal(response)
	if err != nil {
		logger.Log.Error("Marshal response failed", zap.Error(err))
		err = customErrors.NewCustomError(err, http.StatusInternalServerError, "Error reading request body")
		return
	}

}

func (h *Handler) HandlerPostCreditCard(w http.ResponseWriter, r *http.Request) {
	var err error
	var resp []byte

	defer func() {
		if err != nil {
			deferHandler(err, w)
			return
		}
		_, err = w.Write(resp)
		if err != nil {
			logger.Log.Error("Error writing response", zap.Error(err))
		}
		w.WriteHeader(http.StatusOK)

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
			return
		}
		_, err = w.Write(resp)
		if err != nil {
			logger.Log.Error("Error writing response", zap.Error(err))
		}
		w.WriteHeader(http.StatusOK)

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
			return
		}
		_, err = w.Write(resp)
		if err != nil {
			logger.Log.Error("Error writing response", zap.Error(err))
		}
		w.WriteHeader(http.StatusOK)

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

func (h *Handler) HandlerGetData(w http.ResponseWriter, r *http.Request) {
	var err error
	var resp []byte

	defer func() {
		if err != nil {
			deferHandler(err, w)
			return
		}
		_, err = w.Write(resp)
		if err != nil {
			logger.Log.Error("Error writing response", zap.Error(err))
		}
		w.WriteHeader(http.StatusOK)

	}()

	userId, err := getUserId(r)
	if err != nil {
		return
	}

	strUserDataId := chi.URLParam(r, "userDataId")

	userDataId, err := strconv.Atoi(strUserDataId)
	if err != nil {
		err = customErrors.NewCustomError(err, http.StatusBadRequest, "UserDataId is invalid")
		logger.Log.Error("UserDataId is invalid", zap.Error(err))
		return
	}

	resp, err = h.service.Data.GetData(r.Context(), userId, int64(userDataId))
	if err != nil {
		return
	}
}

func (h *Handler) HandlerUpdateData(w http.ResponseWriter, r *http.Request) {
	var err error
	var resp []byte

	defer func() {
		if err != nil {
			deferHandler(err, w)
			return
		}
		_, err = w.Write(resp)
		if err != nil {
			logger.Log.Error("Error writing response", zap.Error(err))
		}
		w.WriteHeader(http.StatusOK)

	}()

	userId, err := getUserId(r)
	if err != nil {
		return
	}

	var updateData *store.UpdateUsersData
	err = json.NewDecoder(r.Body).Decode(&updateData)
	if err != nil {
		logger.Log.Error("Unmarshal json failed", zap.Error(err))
		err = customErrors.NewCustomError(err, http.StatusBadRequest, "Error reading request body")
	}

	err = h.service.Data.UpdateData(r.Context(), int64(userId), updateData, updateData.EncryptData)
	if err != nil {
		return
	}

}
func (h *Handler) HandlerRemoveData(w http.ResponseWriter, r *http.Request) {
	var err error
	var resp []byte

	defer func() {
		if err != nil {
			deferHandler(err, w)
			return
		}
		_, err = w.Write(resp)
		if err != nil {
			logger.Log.Error("Error writing response", zap.Error(err))
		}
		w.WriteHeader(http.StatusOK)

	}()

	userId, err := getUserId(r)
	if err != nil {
		return
	}

	strUserDataId := chi.URLParam(r, "userDataId")
	userDataId, err := strconv.Atoi(strUserDataId)
	if err != nil {
		err = customErrors.NewCustomError(err, http.StatusBadRequest, "UserDataId is invalid")
		logger.Log.Error("UserDataId is invalid", zap.Error(err))
		return
	}

	err = h.service.Data.RemoveData(r.Context(), userId, int64(userDataId))
	if err != nil {
		return
	}

}
