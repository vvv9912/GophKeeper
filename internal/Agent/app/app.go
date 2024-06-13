package app

import (
	"GophKeeper/internal/Agent/command"
	"GophKeeper/internal/Agent/service"
	"GophKeeper/pkg/logger"
	"GophKeeper/pkg/store"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

type App struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

func Run() error {
	return nil
}
func init() {
	fmt.Println(logger.Initialize("info"))

	logger.Log.Info("start app")

	db, err := sqlx.Open("sqlite", "clientdb.db")
	if err != nil {
		panic(err)
		return
	}
	err = store.MigrateSQLITE(db)
	if err != nil {
		panic(err)
		return
	}
	key := []byte("12345678901234567890123456789012")
	//key = generateKey()
	agent := service.NewServiceAgent(db, key)
	fmt.Println("Cobra start")

	cob := command.NewCobra(agent)

	//cob.UpdateBinaryFile(&cobra.Command{}, []string{"15", "/home/vlad/Загрузки/FileZilla_3.66.1_x86_64-linux-gnu.tar.xz"})
	if err := cob.Start(); err != nil {
		panic(err)
		return
	}

	return

}
func generateKey() []byte {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return nil
	}
	s := string(key)
	_ = s
	fmt.Println(string(key))
	return key
}
