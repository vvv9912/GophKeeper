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
	assert.Equal(t, "agent/cert.pem", retrievedConfig.CertFile)
	assert.Equal(t, "agent/key.pem", retrievedConfig.KeyFile)
	assert.Equal(t, "clientdb.db", retrievedConfig.PathDatabaseFile)
	assert.Equal(t, "FileAgent/tmp", retrievedConfig.PathTmpStorage)
	assert.Equal(t, "FileAgent/storage", retrievedConfig.PathLocalStorage)
	assert.Equal(t, "FileAgent/userData", retrievedConfig.PathUserData)
	assert.Equal(t, "12345678901234567890123456789012", retrievedConfig.PathSecretKey)
	assert.Equal(t, "debug", retrievedConfig.LevelLogger)
}
