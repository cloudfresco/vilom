package config

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/go-redis/redis"
	"github.com/palantir/stacktrace"
	"github.com/spf13/viper"
	gomail "gopkg.in/gomail.v2"

	"github.com/cloudfresco/vilom/common"
)

// InitConfig - Confirguration
func InitConfig() (*common.RedisOptions, *sql.DB, *redis.Client, *common.OauthOptions, *gomail.Dialer, *common.KeyOptions, string, string, *common.JWTOptions, *common.RateLimiterOptions, string, *common.UserOptions, error) {
	v := viper.New()
	v.AutomaticEnv()
	var redisObj common.RedisOptions
	redisObj.Addr = v.GetString("VILOM_REDIS_ADDRESS")

	var dbObj common.DbOptions
	dbObj.User = v.GetString("VILOM_DBUSER")
	dbObj.Password = v.GetString("VILOM_DBPASS")
	dbObj.Host = v.GetString("VILOM_DBHOST")
	dbObj.Port = v.GetString("VILOM_DBPORT")
	dbObj.Schema = v.GetString("VILOM_DBNAME")

	db, err := sql.Open("mysql", fmt.Sprint(dbObj.User, ":", dbObj.Password, "@(", dbObj.Host,
		":", dbObj.Port, ")/", dbObj.Schema, "?charset=utf8mb4&parseTime=True"))
	if err != nil {
		log.Error(stacktrace.Propagate(err, ""))
	}
	// make sure connection is available
	err = db.Ping()
	if err != nil {
		log.Error(stacktrace.Propagate(err, ""))
	} else {
		log.Info("Connected to Sql DB")
	}

	redisClient := redis.NewClient(&redis.Options{
		PoolSize:    10, // default
		IdleTimeout: 30 * time.Second,
		Addr:        redisObj.Addr,
		Password:    "", // no password set
		DB:          0,  // use default DB
	})

	var oauth common.OauthOptions
	oauth.ClientID = v.GetString("GOOGLE_OAUTH2_CLIENT_ID")
	oauth.ClientSecret = v.GetString("GOOGLE_OAUTH2_CLIENT_SECRET")

	var mailerObj common.MailerOptions
	mailerObj.Server = v.GetString("VILOM_MAILER_SERVER")
	MailerPort, err := strconv.Atoi(v.GetString("VILOM_MAILER_PORT"))
	if err != nil {
		log.Error(stacktrace.Propagate(err, ""))
	}
	mailerObj.Port = MailerPort
	mailerObj.User = v.GetString("VILOM_MAILER_USER")
	mailerObj.Password = v.GetString("VILOM_MAILER_PASS")

	mailer := gomail.NewDialer(mailerObj.Server, mailerObj.Port, mailerObj.User, mailerObj.Password)

	var keyObj common.KeyOptions
	keyObj.CaCertPath = v.GetString("VILOM_CA_CERT_PATH")
	keyObj.CertPath = v.GetString("VILOM_CERT_PATH")
	keyObj.KeyPath = v.GetString("VILOM_KEY_PATH")

	serverTLS := v.GetString("VILOM_SERVER_TLS")
	serverAddr := v.GetString("VILOM_SERVER_ADDRESS")

	var jwtObj common.JWTOptions

	JWTKey := v.GetString("VILOM_JWT_KEY")
	jwtObj.JWTKey = []byte(JWTKey)
	JWTDuration, err := strconv.Atoi(v.GetString("VILOM_JWT_Duration"))
	if err != nil {
		log.Error(stacktrace.Propagate(err, ""))
	}
	jwtObj.JWTDuration = JWTDuration

	v1 := viper.New()
	v1.SetConfigName("config")
	pwd, _ := os.Getwd()
	viewpath := pwd + filepath.FromSlash("/config")
	v1.AddConfigPath(viewpath)

	if err := v1.ReadInConfig(); err != nil {
		log.Error(stacktrace.Propagate(err, ""))
		os.Exit(1)
	}

	var rateObj common.RateLimiterOptions
	if err := v1.UnmarshalKey("ratelimit", &rateObj); err != nil {
		log.Error(stacktrace.Propagate(err, ""))
	}

	var limit string
	if err := v1.UnmarshalKey("limit", &limit); err != nil {
		log.Error(stacktrace.Propagate(err, ""))
	}

	var userObj common.UserOptions
	if err := v1.UnmarshalKey("useroptions", &userObj); err != nil {
		log.Error(stacktrace.Propagate(err, ""))
	}

	return &redisObj, db, redisClient, &oauth, mailer, &keyObj, serverTLS, serverAddr, &jwtObj, &rateObj, limit, &userObj, nil
}
