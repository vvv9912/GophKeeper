package server

import (
	"GophKeeper/pkg/store"
	"time"
)

type Auth struct {
	Login    string `json:"login" `
	Password string `json:"password" `
}
type User struct {
	Login string `json:"login" `
	JWT   string `json:"jwt"`
}

type RespError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ReqData struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Data        []byte `json:"data,omitempty"`
}

type RespData struct {
	UserDataId int64      `json:"userDataId"`
	Hash       string     `json:"hash"`
	CreatedAt  *time.Time `json:"createdAt"`
	UpdateAt   *time.Time `json:"updateAt"`
}

type DataFileInfo struct {
	OriginalFileName string `json:"originalFileName"`
}
type RespUsersData struct {
	InfoUsersData *store.UsersData `json:"infoUsersData"`
	EncryptData   *store.DataFile  `json:"encryptData"`
}
