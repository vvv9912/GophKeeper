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
	CreateCredentials(ctx context.Context, userId int64, data []byte, name, description, hash string) (int64, error)
	CreateCreditCard(ctx context.Context, userId int64, data []byte, name, description, hash string) (int64, error)
	CreateFileData(ctx context.Context, userId int64, data []byte, name, description, hash string) (int64, error)
	ChangeData(ctx context.Context, userId int64, lastTimeUpdate time.Time) ([]UsersData, error)
	GetData(ctx context.Context, userId int64, usersDataId int64) (*UsersData, *DataFile, error)
	UpdateData(ctx context.Context, updateData *UpdateUsersData, data []byte) error
	RemoveData(ctx context.Context, userId, usersDataId int64) error
	CreateFileDataChunks(ctx context.Context, userId int64, data []byte, name, description, hash string, metaData []byte) (int64, error)
	GetMetaData(ctx context.Context, userId, userDataId int64) (*MetaData, error)
}

type DataClient interface {
	CreateFileData(ctx context.Context, data []byte, userDataId int64, name, description, hash string) error
	CreateCredentials(ctx context.Context, data []byte, userDataId int64, name, description, hash string) error
	CreateCreditCard(ctx context.Context, data []byte, userDataId int64, name, description, hash string) error
	GetJWTToken(ctx context.Context) (string, error)
	SetJWTToken(ctx context.Context, JWTToken string) error

	CreateBinaryFile(ctx context.Context, data []byte, userDataId int64, name, description, hash string, metaData *MetaData) error
}
