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
	CreateCredentials(ctx context.Context, userId int64, data []byte, name, description, hash string) (*UsersData, error)
	CreateCreditCard(ctx context.Context, userId int64, data []byte, name, description, hash string) (*UsersData, error)
	CreateFileData(ctx context.Context, userId int64, data []byte, name, description, hash string) (*UsersData, error)
	ChangeData(ctx context.Context, userId int64, userDataId int64, lastTimeUpdate time.Time) (bool, error)
	GetData(ctx context.Context, userId int64, usersDataId int64) (*UsersData, *DataFile, error)
	UpdateData(ctx context.Context, userId int64, userDataId int64, data []byte, hash string) (*UsersData, error)
	RemoveData(ctx context.Context, userId, usersDataId int64) error
	CreateFileDataChunks(ctx context.Context, userId int64, data []byte, name string, description string, hash string, metaData *MetaData) (*UsersData, error)
	GetMetaData(ctx context.Context, userId, userDataId int64) (*MetaData, error)
	GetListData(ctx context.Context, userId int64) ([]UsersData, error)
}

type DataClient interface {
	CreateFileData(ctx context.Context, data []byte, userDataId int64, name, description, hash string) error
	CreateCredentials(ctx context.Context, data []byte, userDataId int64, name, description, hash string) error
	CreateCreditCard(ctx context.Context, data []byte, userDataId int64, name, description, hash string) error
	GetJWTToken(ctx context.Context) (string, error)
	SetJWTToken(ctx context.Context, JWTToken string) error

	CreateBinaryFile(ctx context.Context, data []byte, userDataId int64, name, description, hash string, metaData *MetaData) error
}
