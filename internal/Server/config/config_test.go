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
	assert.Equal(t, "server/cert.pem", retrievedConfig.CertFile)
	assert.Equal(t, "server/key.pem", retrievedConfig.KeyFile)
	assert.Equal(t, "debug", retrievedConfig.LevelLogger)
	assert.Equal(t, "asdahgf53sk41250", retrievedConfig.SecretKey)
	assert.Equal(t, "postgres://postgres:postgres@localhost:5434/postgres?sslmode=disable", retrievedConfig.DatabaseDNS)
}

func TestInitConfig_OverrideValuesFromFiles(t *testing.T) {
	f, err := os.Create("config.hcl")
	require.NoError(t, err)
	defer os.Remove("config.hcl")
	err = f.Close()
	require.NoError(t, err)
	InitConfig()

	g := Get()

	assert.Equal(t, ":8080", g.ServerDNS)
	assert.Equal(t, "server/cert.pem", g.CertFile)
	assert.Equal(t, "server/key.pem", g.KeyFile)
	assert.Equal(t, "debug", g.LevelLogger)
	assert.Equal(t, "asdahgf53sk41250", g.SecretKey)
	assert.Equal(t, "postgres://postgres:postgres@localhost:5434/postgres?sslmode=disable", g.DatabaseDNS)
}
