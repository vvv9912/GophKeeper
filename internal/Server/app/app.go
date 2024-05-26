package app

import (
	"context"
	"github.com/jackc/pgx/v5"
)

type App struct {
	db *pgx.Conn
}

func NewApp(ctx context.Context, connString string) (*App, error) {
	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		return nil, err
	}
	_ = conn
	return nil, err
}

func Run() {
}
