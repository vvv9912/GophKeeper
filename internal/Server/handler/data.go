package handler

import (
	"GophKeeper/pkg/customErrors"
	"GophKeeper/pkg/logger"
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
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
	userId := int64(1)
	//userId, err := getUserId(r)
	//if err != nil {
	//	return
	//}

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
