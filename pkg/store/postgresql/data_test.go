//go:build integrationTest

package postgresql

import (
	"GophKeeper/pkg/store"
	"context"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"log"
	"os"
	"testing"
	"time"
)

type TestDatabaseData interface {
	store.Data

	clean(ctx context.Context) error
}

type TestDatabaseUser interface {
	store.Auth

	clean(ctx context.Context) error
}

type сonfig struct {
	ConnectTimeout time.Duration
	QueryTimeout   time.Duration
	Username       string
	Password       string
	DBName         string
	MigrationVer   int

	Host string
	Port int
}

type PostrgresTestSuite struct {
	suite.Suite
	TestDatabaseData
	TestDatabaseUser

	tc  *tcpostgres.PostgresContainer
	cfg *сonfig
}

func (ts *PostrgresTestSuite) SetupSuite() {

	cfg := &сonfig{
		ConnectTimeout: 5 * time.Second,
		QueryTimeout:   5 * time.Second,
		Username:       "postgres",
		Password:       "test",
		DBName:         "postgres",
		MigrationVer:   1,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	// Set DOCKER_REGISTRY_MIRROR (Зеркало для скачивания)
	//"registry-mirrors": [
	//"https://dockerhub1.beget.com",
	//"https://mirror.gcr.io"
	//]
	err := os.Setenv("DOCKER_REGISTRY_MIRROR", "https://mirror.gcr.io")
	if err != nil {
		panic(err)
	}
	pgc, err := tcpostgres.Run(ctx,
		"postgres:latest",
		tcpostgres.WithDatabase(cfg.DBName),
		tcpostgres.WithUsername(cfg.Username),
		tcpostgres.WithPassword(cfg.Password),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connection").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)

	require.NoError(ts.T(), err)

	cfg.Host, err = pgc.Host(ctx)
	require.NoError(ts.T(), err)

	port, err := pgc.MappedPort(ctx, "5432")
	require.NoError(ts.T(), err)

	cfg.Port = port.Int()

	ts.tc = pgc
	ts.cfg = cfg

	database_dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
	log.Println(database_dsn)
	db, err := sqlx.Connect("pgx", database_dsn)
	require.NoError(ts.T(), err)
	storage := NewDatabase(db)

	ts.TestDatabaseData = storage
	ts.TestDatabaseUser = storage
	err = store.MigratePostgres(db)
	require.NoError(ts.T(), err)

}

func (db *Database) clean(ctx context.Context) error {
	Newctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err := db.db.ExecContext(Newctx, "DELETE FROM users")
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = db.db.ExecContext(Newctx, "DELETE FROM users_data")
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = db.db.ExecContext(Newctx, "DELETE FROM data")
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = db.db.ExecContext(Newctx, "DELETE FROM data_types")
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// Выполняется после всех тестов, уничтожаем контейнер
func (ts *PostrgresTestSuite) TearDownSuite() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	require.NoError(ts.T(), ts.tc.Terminate(ctx))
}

// Вызов тестов
func TestPostgres(t *testing.T) {
	suite.Run(t, new(PostrgresTestSuite))
}

// Выполняется перед запуском каждого теста
func (ts *PostrgresTestSuite) SetupTest() {
	ts.Require().NoError(ts.TestDatabaseData.clean(context.Background()))
}

// Выполняется перед запуском каждого теста
func (ts *PostrgresTestSuite) TearDownTest() {
	ts.Require().NoError(ts.TestDatabaseData.clean(context.Background()))
}

func (ts *PostrgresTestSuite) TestIntegrationDB() {

	ctx := context.Background()
	userId := int64(1)
	data := []byte("test data")
	name := "Test User"
	description := "Test Description"
	hash := "test-hash"

	usersData, err := ts.TestDatabaseData.CreateCredentials(ctx, userId, data, name, description, hash)
	ts.Require().NoError(err)

	ts.Require().Equal(usersData.Name, name)
	ts.Require().Equal(usersData.Description, description)
	ts.Require().Equal(usersData.Hash, hash)
	ts.Require().Equal(usersData.UserId, userId)
	ts.Require().NotEqual(usersData.CreatedAt, nil)
	ts.Require().NotEqual(usersData.UpdateAt, nil)
	ts.Require().NotEqual(usersData.UserDataId, 0)

	usersData, DataFile, err := ts.TestDatabaseData.GetData(ctx, userId, int64(1))
	ts.Require().NoError(err)

	ts.Require().Equal(usersData.Name, name)
	ts.Require().Equal(usersData.Description, description)
	ts.Require().Equal(usersData.Hash, hash)
	ts.Require().Equal(usersData.UserId, userId)
	ts.Require().NotEqual(usersData.CreatedAt, nil)
	ts.Require().NotEqual(usersData.UpdateAt, nil)
	ts.Require().NotEqual(usersData.UserDataId, 0)
	ts.Require().Equal(DataFile.DataId, int(1))
	ts.Require().Equal(DataFile.EncryptData, data)

}

// GetListData
func (ts *PostrgresTestSuite) TestIntegrationDBGetListData() {

	ctx := context.Background()
	userId := int64(1)
	data := []byte("test data")
	name := "Test User"
	description := "Test Description"
	hash := "test-hash"

	usersData, err := ts.TestDatabaseData.CreateCredentials(ctx, userId, data, name, description, hash)
	ts.Require().NoError(err)

	ts.Require().Equal(usersData.Name, name)
	ts.Require().Equal(usersData.Description, description)
	ts.Require().Equal(usersData.Hash, hash)
	ts.Require().Equal(usersData.UserId, userId)
	ts.Require().NotEqual(usersData.CreatedAt, nil)
	ts.Require().NotEqual(usersData.UpdateAt, nil)
	ts.Require().NotEqual(usersData.UserDataId, 0)

	data2 := []byte("test data")
	name2 := "Test User2"
	description2 := "Test Description2"
	hash2 := "test-hash2"
	usersData, err = ts.TestDatabaseData.CreateCreditCard(ctx, userId, data2, name2, description2, hash2)
	ts.Require().NoError(err)

	ts.Require().Equal(usersData.Name, name2)
	ts.Require().Equal(usersData.Description, description2)
	ts.Require().Equal(usersData.Hash, hash2)
	ts.Require().Equal(usersData.UserId, userId)
	ts.Require().NotEqual(usersData.CreatedAt, nil)
	ts.Require().NotEqual(usersData.UpdateAt, nil)
	ts.Require().NotEqual(usersData.UserDataId, 0)

	dataList, err := ts.TestDatabaseData.GetListData(ctx, userId)
	ts.Require().NoError(err)
	ts.Require().Equal(len(dataList), 2)
	ts.Require().Equal(dataList[0].Name, name)
	ts.Require().Equal(dataList[0].Description, description)
	ts.Require().Equal(dataList[0].Hash, hash)
	ts.Require().Equal(dataList[0].UserId, userId)
	ts.Require().NotEqual(dataList[0].CreatedAt, nil)
	ts.Require().NotEqual(dataList[0].UpdateAt, nil)
	ts.Require().NotEqual(dataList[0].UserDataId, 0)
	ts.Require().Equal(dataList[1].Name, name2)
	ts.Require().Equal(dataList[1].Description, description2)
	ts.Require().Equal(dataList[1].Hash, hash2)
	ts.Require().Equal(dataList[1].UserId, userId)
	ts.Require().NotEqual(dataList[1].CreatedAt, nil)
	ts.Require().NotEqual(dataList[1].UpdateAt, nil)
	ts.Require().NotEqual(dataList[1].UserDataId, 0)

}

func (ts *PostrgresTestSuite) TestCreateUsers() {
	userId, err := ts.TestDatabaseUser.CreateUser(context.Background(), "test", "testPassword")
	require.NoError(ts.T(), err)
	require.NotEqual(ts.T(), userId, 0)

	uId, err := ts.TestDatabaseUser.GetUserId(context.Background(), "test", "testPassword")
	require.NoError(ts.T(), err)

	require.Equal(ts.T(), uId, userId)
}
