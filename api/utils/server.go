package utils

import (
	"database/sql"
	"log"
	"os"
	"strconv"

	"cloud.google.com/go/storage"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/twilio/twilio-go"
)

func Getenv(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

type Config struct {
	AppName                                string
	ServerHost                             string
	SaltKey                                []byte
	SecretKey                              []byte
	State                                  string
	MaxProcessNumber                       int
	MySQLHost                              string
	MySQLPort                              string
	MySQLUsername                          string
	MySQLPassword                          string
	MySQLDatabaseName                      string
	RedisHost                              string
	RedisPort                              string
	RedisPassword                          string
	RedisDB                                string
	TwilioAccountSID                       string
	TwilioAuthToken                        string
	TwilioPhoneNumber                      string
	TwilioPhoneNumberVerificationServiceID string
}

// Server Struct
type Server struct {
	Config        Config
	Router        *gin.Engine
	DB            *sql.DB
	RedisDB       *redis.Client
	ObjectStorage *storage.Client
	SMSService    *twilio.RestClient
}

// NewServer is a constructor for Server Struct
func NewServer() (*Server, error) {

	var err error

	server := Server{}

	var maxProcessNumber int
	maxProcessNumber, err = strconv.Atoi(Getenv("MAX_PROCESS_NUMBER", "2"))

	if err != nil {
		return nil, err
	}

	server.Config.MaxProcessNumber = maxProcessNumber
	server.Config.State = Getenv("GIN_MODE", "debug")
	server.Config.ServerHost = Getenv("SERVER_HOST", ":80")

	server.Config.RedisHost = Getenv("REDIS_HOST", "localhost")
	server.Config.RedisPort = Getenv("REDIS_PORT", "6379")
	server.Config.RedisPassword = Getenv("REDIS_PASSWORD", "")
	server.Config.RedisDB = Getenv("REDIS_DB", "0")

	server.Config.MySQLHost = Getenv("MYSQL_HOST", "localhost")
	server.Config.MySQLPort = Getenv("MYSQL_PORT", "3306")
	server.Config.MySQLUsername = Getenv("MYSQL_USERNAME", "root")
	server.Config.MySQLPassword = Getenv("MYSQL_PASSWORD", "")
	server.Config.MySQLDatabaseName = Getenv("MYSQL_DATABASE", "yum")

	server.Config.TwilioAccountSID = Getenv("TWILIO_ACCOUNT_SID", "")
	server.Config.TwilioAuthToken = Getenv("TWILIO_AUTH_TOKEN", "")
	server.Config.TwilioPhoneNumber = Getenv("TWILIO_PHONE_NUMBER", "")
	server.Config.TwilioPhoneNumberVerificationServiceID = Getenv("TWILIO_PHONE_NUMBER_VERIFICATION_SERVICE_ID", "")

	server.Config.SaltKey = []byte("d\x8f\xef\x83`\xb1*\xd5[\xedu\xdb0\x8bJ\x94\xe0\xf0\xa5\xf1\x91\xc7t\xa0")
	server.Config.SecretKey = []byte("\xec\xbb\x81\x1fy\xff\tDi\xca\xc9\xd5\x92f{L\xadNh}fz\xe5\x04HS\x92x\x1f\xf0\xd2c,\xb0\xf2Z\xcfz\ru\x86\xfb)%\x89\xc5\x89Im\x84\xde\xeb\x15\xe6\xe5\x04A\xa5p\xeal\x97\xcb\xb7<\xb8y\xfb\xa0;V h\x0f\xc0YK\r\xa3\x8cq\x9f\x19?\xdf\n\xd8B\r \xe7s-\xd1\x1dG\x1bw\xa1\xef\x8f\xc6\xbe\x98\x90\xa7\xf4g\xc1\xcfn@\xe2\x83\x8b\xfb\xbb+\x94d\xb3\x98fD\x87\xe9\xe6m\x99\xee&_\xf9\xd1p\x99\xe7\x99}\xd9\x1b\x1fIj\x836r\xad\xff\xfd\x8dt\xcdFe\x9c\x8c\xd5S\x8a\xe2U\xad\xbd\xccw\xe6\xaf\xec\x0c\xd54?X\xf1\x15\xf1i\x01\x9er\x120\xb8\x05}~\x92BY\x14\xf1\xf5R\n|\xa5\xf7'\xbb\xe5,\x84\xbf\xe8\x0eH\xc3\x9b`\xc0u\xedj\x10Y\xb7\xcbu\xcf:\x8d\x93\xd6\xd0\xe3z)W*z\xd6\xc6\xb6\xd2'\xbfD\x16`]\x12\xcb\x7f[\xfc\xd0\xed\x869o\xa0\xef\xe0\xa3\xa0")

	var dbClient *sql.DB

	dbClient, err = NewDbConnection(
		server.Config.MySQLHost,
		server.Config.MySQLPort,
		server.Config.MySQLUsername,
		server.Config.MySQLPassword,
		server.Config.MySQLDatabaseName,
	)

	if err != nil {
		return nil, err
	}

	server.DB = dbClient

	var rdb *redis.Client

	rdb, err = NewRedisConnection(
		server.Config.RedisHost,
		server.Config.RedisPort,
		server.Config.RedisPassword,
		server.Config.RedisDB,
	)

	if err != nil {
		return nil, err
	}

	server.RedisDB = rdb

	var objectStorageClient *storage.Client

	objectStorageClient, err = NewObjectStorage()

	if err != nil {
		return nil, err
	}

	server.ObjectStorage = objectStorageClient

	server.SMSService = NewTwillioConnection(
		server.Config.TwilioAccountSID,
		server.Config.TwilioAuthToken,
	)

	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)

	gin.SetMode(server.Config.State)
	server.Router = gin.Default()

	return &server, nil
}

func (server *Server) Close() error {

	var err error

	err = server.DB.Close()

	if err != nil {
		return err
	}

	err = server.RedisDB.Close()

	if err != nil {
		return err
	}

	err = server.ObjectStorage.Close()

	if err != nil {
		return err
	}

	return nil

}
