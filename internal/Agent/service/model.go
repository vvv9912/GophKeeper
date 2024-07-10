package service

import "GophKeeper/pkg/store"

// Data - информация о пользовательском файле
type Data struct {
	InfoUsersData *store.UsersData `json:"infoUsersData"` // Информация о пользовательском файле.
	EncryptData   *store.DataFile  `json:"encryptData"`   // Зашифрованные данные.
}

var PathStorage = "FileAgent/storage"   // PathStorage - путь к хранилищу.
var PathTmp = "FileAgent/tmp"           // PathTmp - путь к временной папке.
var PathUserData = "FileAgent/userData" // PathUserData - путь к папке с пользовательскими файлами.

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
