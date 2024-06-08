package service

type RespData struct {
	UserDataId int64  `json:"userDataId"`
	Hash       string `json:"hash"`
}

type SaveFile struct {
	FileName         string `json:"fileName"`
	OriginalFileName string `json:"originalFileName"`
	PathSave         string `json:"pathSave"`
	Size             int64  `json:"size"`
}
