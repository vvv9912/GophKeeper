package config

import (
	"github.com/cristalhq/aconfig"
	"github.com/cristalhq/aconfig/aconfighcl"
	"sync"
)

// Config - структура конфигурации
type Config struct {
	CertFile         string `hcl:"certfile" env:"CertFile" `                            // Сертификат.
	KeyFile          string `hcl:"keyfile" env:"KeyFile" `                              // Ключ.
	ServerDNS        string `hcl:"server_dns" env:"SERVER_DNS"`                         // DNS адрес сервера.
	PathDatabaseFile string `hcl:"pathFileLogger,omitempty" env:"PATH_FILE_LOGGER"`     // Путь к базе данных.
	PathTmpStorage   string `hcl:"pathTmpStorage,omitempty" env:"PATH_TMP_STORAGE"`     // Путь к временному хранилищу.
	PathLocalStorage string `hcl:"pathLocalStorage,omitempty" env:"PATH_LOCAL_STORAGE"` // Путь к локальному хранилищу.
	PathUserData     string `hcl:"pathUserData,omitempty" env:"PATH_USER_DATA"`         // Путь пользовательских сохраненных данных (разархивированные данные).
	PathSecretKey    string `hcl:"pathSecretKey" env:"PATH_SECRET_KEY"`                 // Путь к файлу с секретным ключом.
	LevelLogger      string `hcl:"levelLogger" env:"LEVEL_LOGGER"`                      // Уровень логгера.
}

func (c *Config) setDefaultValues() {
	c.ServerDNS = ":8080"
	c.CertFile = "agent/cert.pem"
	c.KeyFile = "agent/key.pem"
	c.PathDatabaseFile = "clientdb.db"
	c.PathTmpStorage = "FileAgent/tmp"
	c.PathLocalStorage = "FileAgent/storage"
	c.PathUserData = "FileAgent/userData"
	c.PathSecretKey = "12345678901234567890123456789012"
	c.LevelLogger = "debug"

}

var (
	cfg  Config
	once sync.Once
)

// InitConfig - инициализация конфигурации
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

// Get - получение конфигурации
func Get() Config {
	return cfg
}
