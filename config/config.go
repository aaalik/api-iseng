package config

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/lib/pq"
	"github.com/subosito/gotenv"

	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"
	"github.com/sirupsen/logrus"
)

const (
	configSQLDBDriver    = "DB_DRIVER"
	configSQLDBWriteHost = "DB_WRITER_HOST"
	configSQLDBWritePort = "DB_WRITER_PORT"
	configSQLDBWriteUser = "DB_WRITER_USER"
	configSQLDBWritePass = "DB_WRITER_PASS"
	configSQLDBWriteName = "DB_WRITER_NAME"

	configSQLDBReadHost = "DB_READER_HOST"
	configSQLDBReadPort = "DB_READER_PORT"
	configSQLDBReadUser = "DB_READER_USER"
	configSQLDBReadPass = "DB_READER_PASS"
	configSQLDBReadName = "DB_READER_NAME"

	configHostAddress = "HOST_ADDRESS"
	configStopTimeout = "STOP_TIMEOUT"
)

type Config struct {
	Host        Host
	StopTimeout int
	SQLDB       SQLDB
}

type Host struct {
	Address string
}

type HttpService struct {
	URL string
}

type SQLDB struct {
	Driver string
	Write  SQLDBConfig
	Read   SQLDBConfig
}

type SQLDBConfig struct {
	Host string
	Port string
	User string
	Pass string
	Name string
}

func NewConfig(mustLoad bool) Config {
	gotenv.Load(".env")

	stopTimeout, err := strconv.Atoi(os.Getenv(configStopTimeout))
	if err != nil {
		panic(err)
	}

	return Config{
		Host: Host{
			Address: os.Getenv(configHostAddress),
		},
		StopTimeout: stopTimeout,
		SQLDB: SQLDB{
			Driver: os.Getenv(configSQLDBDriver),
			Write: SQLDBConfig{
				Host: os.Getenv(configSQLDBWriteHost),
				Port: os.Getenv(configSQLDBWritePort),
				User: os.Getenv(configSQLDBWriteUser),
				Pass: os.Getenv(configSQLDBWritePass),
				Name: os.Getenv(configSQLDBWriteName),
			},
			Read: SQLDBConfig{
				Host: os.Getenv(configSQLDBReadHost),
				Port: os.Getenv(configSQLDBReadPort),
				User: os.Getenv(configSQLDBReadUser),
				Pass: os.Getenv(configSQLDBReadPass),
				Name: os.Getenv(configSQLDBReadName),
			},
		},
	}
}

func (cf Config) NewLogrus() *logrus.Logger {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	log := logrus.StandardLogger()
	return log
}

func (cf Config) NewSQL() (reader, writer *sqlx.DB) {
	reader, err := sqlx.Connect(
		cf.SQLDB.Driver,
		fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			cf.SQLDB.Read.Host,
			cf.SQLDB.Read.Port,
			cf.SQLDB.Read.User,
			cf.SQLDB.Read.Pass,
			cf.SQLDB.Read.Name,
		),
	)
	if err != nil {
		panic(err)
	}

	err = reader.Ping()
	if err != nil {
		panic(err)
	}

	reader.Mapper = reflectx.NewMapper("json")
	reader.SetMaxIdleConns(0)

	writer, err = sqlx.Connect(
		cf.SQLDB.Driver,
		fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			cf.SQLDB.Write.Host,
			cf.SQLDB.Write.Port,
			cf.SQLDB.Write.User,
			cf.SQLDB.Write.Pass,
			cf.SQLDB.Write.Name,
		),
	)
	if err != nil {
		panic(err)
	}

	err = writer.Ping()
	if err != nil {
		panic(err)
	}

	writer.Mapper = reflectx.NewMapper("json")
	writer.SetMaxIdleConns(0)
	return
}

func (cf Config) NewHTTPClient() http.Client {
	client := http.Client{}
	client.Timeout = time.Second * 60
	return client
}
