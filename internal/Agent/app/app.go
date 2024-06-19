package app

import (
	"GophKeeper/internal/Agent/command"
	"GophKeeper/internal/Agent/config"
	"GophKeeper/internal/Agent/service"
	"GophKeeper/pkg/logger"
	"GophKeeper/pkg/store"
	"errors"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	_ "modernc.org/sqlite"
	"os"
)

func Run() error {
	return nil
}
func init() {
	if err := config.InitConfig(); err != nil {
		panic(err)
	}
	if err := logger.Initialize(config.Get().LevelLogger); err != nil {
		panic(err)
	}

	logger.Log.Info("start app, config: ", zap.Any("config", config.Get()))

	key, err := readKeyFromFile(config.Get().PathSecretKey)
	if err != nil {
		panic(err)
	}

	db, err := sqlx.Open("sqlite", config.Get().PathDatabaseFile)
	if err != nil {
		panic(err)
	}

	err = store.MigrateSQLITE(db)
	if err != nil {
		panic(err)
	}

	agent := service.NewServiceAgent(db, key, config.Get().CertFile, config.Get().KeyFile, config.Get().PathDatabaseFile)
	cob := command.NewCobra(agent)

	//cob.UpdateBinaryFile(&cobra.Command{}, []string{"15", "/home/vlad/Загрузки/FileZilla_3.66.1_x86_64-linux-gnu.tar.xz"})
	if err := cob.Start(); err != nil {
		panic(err)
	}

	return

}

//	func generateKey() []byte {
//		key := make([]byte, 32)
//		_, err := rand.Read(key)
//		if err != nil {
//			return nil
//		}
//		s := string(key)
//		_ = s
//		fmt.Println(string(key))
//		return key
//	}
func readKeyFromFile(filePath string) ([]byte, error) {
	key, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Проверка длины ключа
	if len(key) != 32 {
		return nil, errors.New("Неверная длина ключа")
	}

	return key, nil
}
