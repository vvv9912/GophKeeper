package service

import "GophKeeper/pkg/store"

type Data struct {
	InfoUsersData *store.UsersData `json:"infoUsersData"`
	EncryptData   *store.DataFile  `json:"encryptData"`
}

var PathStorage = "FileAgent/storage"
var PathTmp = "FileAgent/tmp"
var PathUserData = "FileAgent/userData"

// NewPath - создание новых путей сохранения данных
func NewPath(path ...string) {
	if len(path) >= 1 {
		PathStorage = path[0]
	}

	if len(path) >= 2 {
		PathTmp = path[1]
	}

	if len(path) >= 3 {
		PathUserData = path[2]
	}
}
