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
	"github.com/spf13/viper"
	gomail "gopkg.in/gomail.v2"

	"github.com/cloudfresco/vilom/common"
)

// InitConfig - Confirguration
func InitConfig() (*common.RedisOptions, *sql.DB, *redis.Client, *common.OauthOptions, *gomail.Dialer, *common.KeyOptions, string, string, *common.JWTOptions, *common.RateLimiterOptions, string, *common.UserOptions, error) {
	v := viper.New()
	v.AutomaticEnv()
	var redisOpt common.RedisOptions
	redisOpt.Addr = v.GetString("VILOM_REDIS_ADDRESS")

	var dbOpt common.DbOptions
	dbOpt.User = v.GetString("VILOM_DBUSER")
	dbOpt.Password = v.GetString("VILOM_DBPASS")
	dbOpt.Host = v.GetString("VILOM_DBHOST")
	dbOpt.Port = v.GetString("VILOM_DBPORT")
	dbOpt.Schema = v.GetString("VILOM_DBNAME")

	db, err := sql.Open("mysql", fmt.Sprint(dbOpt.User, ":", dbOpt.Password, "@(", dbOpt.Host,
		":", dbOpt.Port, ")/", dbOpt.Schema, "?charset=utf8mb4&parseTime=True"))
	if err != nil {
		log.WithFields(log.Fields{
			"msgnum": 500,
		}).Error(err)
	}
	// make sure connection is available
	err = db.Ping()
	if err != nil {
		log.WithFields(log.Fields{
			"msgnum": 501,
		}).Error(err)
	} else {
		log.WithFields(log.Fields{
			"msgnum": 502,
		}).Info("Connected to Sql DB")
	}

	redisClient := redis.NewClient(&redis.Options{
		PoolSize:    10, // default
		IdleTimeout: 30 * time.Second,
		Addr:        redisOpt.Addr,
		Password:    "", // no password set
		DB:          0,  // use default DB
	})

	var oauth common.OauthOptions
	oauth.ClientID = v.GetString("GOOGLE_OAUTH2_CLIENT_ID")
	oauth.ClientSecret = v.GetString("GOOGLE_OAUTH2_CLIENT_SECRET")

	var mailerOpt common.MailerOptions
	mailerOpt.Server = v.GetString("VILOM_MAILER_SERVER")
	MailerPort, err := strconv.Atoi(v.GetString("VILOM_MAILER_PORT"))
	if err != nil {
		log.WithFields(log.Fields{
			"msgnum": 503,
		}).Error(err)
	}
	mailerOpt.Port = MailerPort
	mailerOpt.User = v.GetString("VILOM_MAILER_USER")
	mailerOpt.Password = v.GetString("VILOM_MAILER_PASS")

	mailer := gomail.NewDialer(mailerOpt.Server, mailerOpt.Port, mailerOpt.User, mailerOpt.Password)

	var keyOpt common.KeyOptions
	keyOpt.CaCertPath = v.GetString("VILOM_CA_CERT_PATH")
	keyOpt.CertPath = v.GetString("VILOM_CERT_PATH")
	keyOpt.KeyPath = v.GetString("VILOM_KEY_PATH")

	serverTLS := v.GetString("VILOM_SERVER_TLS")
	serverAddr := v.GetString("VILOM_SERVER_ADDRESS")

	var jwtOpt common.JWTOptions

	JWTKey := v.GetString("VILOM_JWT_KEY")
	jwtOpt.JWTKey = []byte(JWTKey)
	JWTDuration, err := strconv.Atoi(v.GetString("VILOM_JWT_DURATION"))
	if err != nil {
		log.WithFields(log.Fields{
			"msgnum": 504,
		}).Error(err)
	}
	jwtOpt.JWTDuration = JWTDuration

	v1 := viper.New()
	v1.SetConfigName("config")
	pwd, _ := os.Getwd()
	viewpath := pwd + filepath.FromSlash("/config")
	v1.AddConfigPath(viewpath)

	if err := v1.ReadInConfig(); err != nil {
		log.WithFields(log.Fields{
			"msgnum": 505,
		}).Error(err)
		os.Exit(1)
	}

	var rateOpt common.RateLimiterOptions
	if err := v1.UnmarshalKey("ratelimit", &rateOpt); err != nil {
		log.WithFields(log.Fields{
			"msgnum": 506,
		}).Error(err)
	}

	var limit string
	if err := v1.UnmarshalKey("limit", &limit); err != nil {
		log.WithFields(log.Fields{
			"msgnum": 507,
		}).Error(err)
	}

	var userOpt common.UserOptions
	if err := v1.UnmarshalKey("useroptions", &userOpt); err != nil {
		log.WithFields(log.Fields{
			"msgnum": 508,
		}).Error(err)
	}

	return &redisOpt, db, redisClient, &oauth, mailer, &keyOpt, serverTLS, serverAddr, &jwtOpt, &rateOpt, limit, &userOpt, nil
}
