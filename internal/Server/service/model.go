package service

import "time"

// RespData - информация о данных.
type RespData struct {
	UserDataId int64      `json:"userDataId"` // идентификатор данных.
	Hash       string     `json:"hash"`       // хэш данных.
	CreatedAt  *time.Time `json:"createdAt"`  // время создания.
	UpdateAt   *time.Time `json:"updateAt"`   // время обновления.
}

// SaveFile - информация о сохраненном файле.
type SaveFile struct {
	FileName         string `json:"fileName"`         // имя файла.
	OriginalFileName string `json:"originalFileName"` // оригинальное имя файла.
	PathSave         string `json:"pathSave"`         // путь к файлу на диске.
	Size             int64  `json:"size"`             // размер.
}
