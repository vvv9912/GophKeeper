package service

import "time"

type RespData struct {
	UserDataId int64      `json:"userDataId"`
	Hash       string     `json:"hash"`
	CreatedAt  *time.Time `json:"createdAt"`
	UpdateAt   *time.Time `json:"updateAt"`
}

type SaveFile struct {
	FileName         string `json:"fileName"`
	OriginalFileName string `json:"originalFileName"`
	PathSave         string `json:"pathSave"`
	Size             int64  `json:"size"`
}
