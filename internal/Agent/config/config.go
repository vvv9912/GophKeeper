package config

import (
	"github.com/cristalhq/aconfig"
	"github.com/cristalhq/aconfig/aconfighcl"
	"sync"
)

type Config struct {
	CertFile         string `hcl:"certfile" env:"CertFile" `
	KeyFile          string `hcl:"keyfile" env:"KeyFile" `
	ServerDNS        string `hcl:"server_dns" env:"SERVER_DNS"`
	PathDatabaseFile string `hcl:"pathFileLogger,omitempty" env:"PATH_FILE_LOGGER"`
	PathTmpStorage   string `hcl:"pathTmpStorage,omitempty" env:"PATH_TMP_STORAGE"`
	PathLocalStorage string `hcl:"pathLocalStorage,omitempty" env:"PATH_LOCAL_STORAGE"`
	PathSecretKey    string `hcl:"pathSecretKey" env:"PATH_SECRET_KEY"`
	LevelLogger      string `hcl:"levelLogger" env:"LEVEL_LOGGER"`
}

func (c *Config) setDefaultValues() {
	c.ServerDNS = "default_server_dns"
	c.PathDatabaseFile = "default_path_database_file"
	c.PathTmpStorage = "default_path_tmp_storage"
	c.PathLocalStorage = "default_path_local_storage"
	c.PathSecretKey = "default_path_secret_key"
	c.LevelLogger = "debug"
}

var (
	cfg  Config
	once sync.Once
)

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
func Get() Config {
	return cfg
}
