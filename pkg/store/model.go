package store

import "time"

type UsersData struct {
	UserDataId  int64     `db:"user_data_id"`
	UserId      int64     `db:"user_id"`
	DataId      int64     `db:"data_id"`
	DataType    int       `db:"data_type"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	Hash        string    `db:"hash"`
	CreatedAt   time.Time `db:"created_at"`
	UpdateAt    time.Time `db:"update_at"`
	IsDeleted   bool      `db:"is_deleted"`
}
