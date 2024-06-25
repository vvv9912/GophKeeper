package config

import (
	"github.com/cristalhq/aconfig"
	"github.com/cristalhq/aconfig/aconfighcl"
	"sync"
)

// Config - конфигурация сервера.
type Config struct {
	CertFile    string `hcl:"certfile" env:"CertFile" `        // путь к файлу с сертификатом.
	KeyFile     string `hcl:"keyfile" env:"KeyFile" `          // путь к файлу с приватным ключом.
	ServerDNS   string `hcl:"server_dns" env:"SERVER_DNS"`     // DNS сервера.
	DatabaseDNS string `hcl:"database_dns" env:"DATABASE_DNS"` // DNS базы данных.
	SecretKey   string `hcl:"SecretKey" env:"PATH_SECRET_KEY"` // путь к файлу с приватным ключом для jwt.
	LevelLogger string `hcl:"levelLogger" env:"LEVEL_LOGGER"`  // уровень логирования.
}

func (c *Config) setDefaultValues() {
	c.ServerDNS = ":8080"
	c.CertFile = "server/cert.pem"
	c.KeyFile = "server/key.pem"
	c.SecretKey = "asdahgf53sk41250"
	c.LevelLogger = "debug"
	c.DatabaseDNS = "postgres://postgres:postgres@localhost:5434/postgres?sslmode=disable"
}

var (
	cfg  Config
	once sync.Once
)

// InitConfig - инициализация конфигурации.
func InitConfig() (err error) {
	once.Do(func() {

		cfg.setDefaultValues()

		loader := aconfig.LoaderFor(&cfg, aconfig.Config{
			EnvPrefix: "NFB",
			Files:     []string{"./config.hcl", "./config.local.hcl"},
			FileDecoders: map[string]aconfig.FileDecoder{
				".hcl": aconfighcl.New(),
			},
		})

		if err = loader.Load(); err != nil {
			return
		}
	})

	return err
}

// Get - получение конфигурации.
func Get() Config {
	return cfg
}
