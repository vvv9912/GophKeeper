package store

import "time"

type UsersData struct {
	UserDataId  int64      `json:"user_data_id,omitempty"`
	UserId      int64      `json:"user_id,omitempty"`
	DataId      int64      `json:"data_id,omitempty"`
	DataType    int        `json:"data_type,omitempty"`
	Name        string     `json:"name,omitempty"`
	Description string     `json:"description,omitempty"`
	Hash        string     `json:"hash,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	UpdateAt    *time.Time `json:"update_at,omitempty"`
	IsDeleted   bool       `json:"is_deleted,omitempty"`
}
