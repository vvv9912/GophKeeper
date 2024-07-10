package service

import (
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewUseCase(t *testing.T) {
	db := &sqlx.DB{}
	secretLey := "key"

	u, err := NewUseCase(db, secretLey)
	require.NoError(t, err)
	require.NotNil(t, u)
}
