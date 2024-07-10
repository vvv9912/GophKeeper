package server

import (
	"GophKeeper/pkg/store"
	"time"
)

// Auth представляет информацию об аутентификации пользователя.
type Auth struct {
	Login    string `json:"login"`    // Логин пользователя.
	Password string `json:"password"` // Пароль пользователя.
}

// User представляет информацию о пользователе.
type User struct {
	Login string `json:"login"` // Логин пользователя.
	JWT   string `json:"jwt"`   // JWT токен.
}

// RespError представляет информацию об ошибке при запросе.
type RespError struct {
	Code    int    `json:"code"`    // Код ошибки.
	Message string `json:"message"` // Описание ошибки.
}

// ReqData представляет информацию о данных (Запрос).
type ReqData struct {
	Name        string `json:"name"`           // Имя данных.
	Description string `json:"description"`    // Описание данных.
	Data        []byte `json:"data,omitempty"` // Данные.
}

// RespData представляет информацию о данных (Ответ).
type RespData struct {
	UserDataId int64      `json:"userDataId"` // ID данных.
	Hash       string     `json:"hash"`       // Хэш данных.
	CreatedAt  *time.Time `json:"createdAt"`  // Дата создания данных.
	UpdateAt   *time.Time `json:"updateAt"`   // Дата обновления данных.
}

// DataFileInfo - информация о файле.
type DataFileInfo struct {
	OriginalFileName string `json:"originalFileName"` // Имя файла.
}

// RespUsersData - информация о пользовательском файле.
type RespUsersData struct {
	InfoUsersData *store.UsersData `json:"infoUsersData"` // Информация о пользовательском файле.
	EncryptData   *store.DataFile  `json:"encryptData"`   // Зашифрованные данные.
}
