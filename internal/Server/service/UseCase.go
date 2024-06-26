package service

import (
	"GophKeeper/pkg/logger"
	"GophKeeper/pkg/store/postgresql"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"time"
)

type UseCase struct {
	Auth      // интерфейс аутентификации.
	StoreAuth // интерфейс для работы с БД пользователя.
	StoreData // интерфейс для работы с БД данных пользователя.
	SaveFiles // временное хранилище (сохр при работе с чанками).
}

func NewUseCase(db *sqlx.DB, secretKey string) (*UseCase, error) {
	nDb := postgresql.NewDatabase(db)
	saveFiles, err := NewSaveFiles(10 * time.Minute)

	//todo

	if err != nil {
		logger.Log.Error("Error creating save files", zap.Error(err))
		return nil, err
	}

	return &UseCase{
		//	Auth:      authorization.NewAutorization(9000*time.Minute, secretKey),
		StoreAuth: nDb,
		SaveFiles: *saveFiles,
		StoreData: nDb,
	}, nil

}
