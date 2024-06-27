package service

import (
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewService(t *testing.T) {
	db := &sqlx.DB{}
	secretLey := "key"

	u, err := NewService(db, secretLey)
	require.NoError(t, err)
	require.NotNil(t, u)
}
