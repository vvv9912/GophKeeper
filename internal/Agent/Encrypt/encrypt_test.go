package Encrypt

import (
	"crypto/rand"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestEncrypt_Encrypt(t *testing.T) {
	key := make([]byte, 32)
	rand.Read(key)
	encr := NewEncrypt(key)
	e, err := encr.Encrypt([]byte("hello world"))
	require.NoError(t, err)
	d, err := encr.Decrypt(e)
	require.NoError(t, err)
	require.Equal(t, []byte("hello world"), d)
}
