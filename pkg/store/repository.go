package store

type Auth interface {
	CreateUser(login, password string) (int64, error)
	GetUserId(login, password string) (int64, error)
}
