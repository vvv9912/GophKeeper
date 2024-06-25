package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetConfig(t *testing.T) {

	cfg := Config{}
	cfg.setDefaultValues()

	retrievedConfig := cfg

	assert.Equal(t, ":8080", retrievedConfig.ServerDNS)
	assert.Equal(t, "server/cert.pem", retrievedConfig.CertFile)
	assert.Equal(t, "server/key.pem", retrievedConfig.KeyFile)
	assert.Equal(t, "debug", retrievedConfig.LevelLogger)
	assert.Equal(t, "asdahgf53sk41250", retrievedConfig.SecretKey)
	assert.Equal(t, "postgres://postgres:postgres@localhost:5434/postgres?sslmode=disable", retrievedConfig.DatabaseDNS)
}
