package service

import "GophKeeper/pkg/store"

type Data struct {
	InfoUsersData *store.UsersData `json:"infoUsersData"`
	EncryptData   *store.DataFile  `json:"encryptData"`
}

// todo config
const PathStorage = "FileAgent/storage"
const PathTmp = "FileAgent/tmp"
const PathUserData = "FileAgent/userData"
