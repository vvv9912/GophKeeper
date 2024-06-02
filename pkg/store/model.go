package store

import "time"

type UsersData struct {
	UserDataId  int64      `json:"userDataId,omitempty"`
	UserId      int64      `json:"userId,omitempty"`
	DataId      int64      `json:"dataId,omitempty"`
	DataType    int        `json:"dataType,omitempty"`
	Name        string     `json:"name,omitempty"`
	Description string     `json:"description,omitempty"`
	Hash        string     `json:"hash,omitempty"`
	CreatedAt   *time.Time `json:"createdAt,omitempty"`
	UpdateAt    *time.Time `json:"updateAt,omitempty"`
	IsDeleted   bool       `json:"isDeleted,omitempty"`
}

type DataFile struct {
	DataId      int    `json:"dataId"`
	EncryptData []byte `json:"EncryptData"`
}

type UpdateUsersData struct {
	UserDataId  int64      `json:"userDataId"`
	UserId      int64      `json:"userId,omitempty"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Hash        string     `json:"hash"`
	UpdateAt    *time.Time `json:"updateAt"`
	EncryptData []byte     `json:"EncryptData"`
}
