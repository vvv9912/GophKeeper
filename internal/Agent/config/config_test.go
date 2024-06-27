package config

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
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
func TestInitConfig_OverrideValuesFromFiles(t *testing.T) {
	f, err := os.Create("config.hcl")
	require.NoError(t, err)
	f.Close()
	defer os.Remove("config.hcl")
	InitConfig()

	g := Get()

	assert.Equal(t, ":8080", g.ServerDNS)
	assert.Equal(t, "agent/cert.pem", g.CertFile)
	assert.Equal(t, "agent/key.pem", g.KeyFile)
	assert.Equal(t, "clientdb.db", g.PathDatabaseFile)
	assert.Equal(t, "FileAgent/tmp", g.PathTmpStorage)
	assert.Equal(t, "FileAgent/storage", g.PathLocalStorage)
	assert.Equal(t, "FileAgent/userData", g.PathUserData)
	assert.Equal(t, "12345678901234567890123456789012", g.PathSecretKey)
	assert.Equal(t, "debug", g.LevelLogger)
}
