package store

import (
	"context"
	"time"
)

type Auth interface {
	CreateUser(ctx context.Context, login, password string) (int64, error)
	GetUserId(ctx context.Context, login, password string) (int64, error)
}

type Data interface {
	CreateCredentials(ctx context.Context, userId int64, data []byte, name, description, hash string) error
	CreateCreditCard(ctx context.Context, userId int64, data []byte, name, description, hash string) error
	CreateFileData(ctx context.Context, userId int64, data []byte, name, description, hash string) error
	ChangeData(ctx context.Context, userId int64, lastTimeUpdate time.Time) ([]UsersData, error)
	GetData(ctx context.Context, userId int64, usersDataId int64) (*UsersData, *DataFile, error)
	UpdateData(ctx context.Context, updateData *UpdateUsersData, data []byte) error
	RemoveData(ctx context.Context, usersDataId int64) error
}
