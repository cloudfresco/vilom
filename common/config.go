package common

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	gomail "gopkg.in/gomail.v2"
)

/* error message range: 500-749 */

// RedisOptions - for redis config
type RedisOptions struct {
	Addr string `mapstructure:"addr"`
}

// KeyOptions - for server config
type KeyOptions struct {
	CaCertPath string `mapstructure:"CaCerTPath"`
	CertPath   string `mapstructure:"CertPath"`
	KeyPath    string `mapstructure:"KeyPath"`
	ServerAddr string `mapstructure:"ServerAddr"`
}

// OauthOptions - for oauth config
type OauthOptions struct {
	ClientID     string `mapstructure:"ClientID"`
	ClientSecret string `mapstructure:"ClientSecret"`
}

// DbMysql for DbType is mysql
const DbMysql string = "mysql"

// DbPgsql for DbType is pgsql
const DbPgsql string = "pgsql"

// DbOptions - for db config
type DbOptions struct {
	DB       string `mapstructure:"db"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Host     string `mapstructure:"hostname"`
	Port     string `mapstructure:"port"`
	Schema   string `mapstructure:"database"`
}

// MailerOptions - for mailer config
type MailerOptions struct {
	User     string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Port     int    `mapstructure:"port"`
	Server   string `mapstructure:"server"`
}

// JWTOptions - for JWT config
type JWTOptions struct {
	JWTKey      []byte
	JWTDuration int
}

// RateLimiterOptions - for rate limiting requests
type RateLimiterOptions struct {
	UserMaxRate    int `mapstructure:"usermaxrate"`
	UserMaxBurst   int `mapstructure:"usermaxburst"`
	UgroupMaxRate  int `mapstructure:"ugroupmaxrate"`
	UgroupMaxBurst int `mapstructure:"ugroupmaxburst"`
	CatMaxRate     int `mapstructure:"catmaxrate"`
	CatMaxBurst    int `mapstructure:"catmaxburst"`
	TopicMaxRate   int `mapstructure:"topicmaxrate"`
	TopicMaxBurst  int `mapstructure:"topicmaxburst"`
	MsgMaxRate     int `mapstructure:"msgmaxrate"`
	MsgMaxBurst    int `mapstructure:"msgmaxburst"`
	UbadgeMaxRate  int `mapstructure:"ubadgemaxrate"`
	UbadgeMaxBurst int `mapstructure:"ubadgemaxburst"`
	SearchMaxRate  int `mapstructure:"searchmaxrate"`
	SearchMaxBurst int `mapstructure:"searchmaxburst"`
	UMaxRate       int `mapstructure:"umaxrate"`
	UMaxBurst      int `mapstructure:"umaxburst"`
}

// UserOptions - for user login
type UserOptions struct {
	ConfirmTokenDuration string `mapstructure:"confirmtokenduration"`
	ResetTokenDuration   string `mapstructure:"resettokenduration"`
}

// GetConfig - Confirguration
func GetConfig() (*RedisOptions, *sql.DB, *redis.Client, *OauthOptions, *gomail.Dialer, *KeyOptions, string, string, *JWTOptions, *RateLimiterOptions, string, *UserOptions, error) {

	v := viper.New()
	v.AutomaticEnv()
	var redisOpt RedisOptions
	redisOpt.Addr = v.GetString("VILOM_REDIS_ADDRESS")

	var dbOpt DbOptions
	var db *sql.DB
	var err error
	dbOpt.DB = v.GetString("VILOM_DB")
	dbOpt.User = v.GetString("VILOM_DBUSER")
	dbOpt.Password = v.GetString("VILOM_DBPASS")
	dbOpt.Host = v.GetString("VILOM_DBHOST")
	dbOpt.Port = v.GetString("VILOM_DBPORT")
	dbOpt.Schema = v.GetString("VILOM_DBNAME")

	if dbOpt.DB == DbMysql {
		db, err = sql.Open("mysql", fmt.Sprint(dbOpt.User, ":", dbOpt.Password, "@(", dbOpt.Host,
			":", dbOpt.Port, ")/", dbOpt.Schema, "?charset=utf8mb4&parseTime=True"))
		if err != nil {
			log.WithFields(log.Fields{
				"msgnum": 500,
			}).Error(err)
		}
	} else if dbOpt.DB == DbPgsql {

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

	var oauth OauthOptions
	oauth.ClientID = v.GetString("GOOGLE_OAUTH2_CLIENT_ID")
	oauth.ClientSecret = v.GetString("GOOGLE_OAUTH2_CLIENT_SECRET")

	var mailerOpt MailerOptions
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

	var keyOpt KeyOptions
	keyOpt.CaCertPath = v.GetString("VILOM_CA_CERT_PATH")
	keyOpt.CertPath = v.GetString("VILOM_CERT_PATH")
	keyOpt.KeyPath = v.GetString("VILOM_KEY_PATH")

	serverTLS := v.GetString("VILOM_SERVER_TLS")
	serverAddr := v.GetString("VILOM_SERVER_ADDRESS")

	var jwtOpt JWTOptions

	JWTKey := v.GetString("VILOM_JWT_KEY")
	jwtOpt.JWTKey = []byte(JWTKey)
	JWTDuration, err := strconv.Atoi(v.GetString("VILOM_JWT_DURATION"))
	if err != nil {
		log.WithFields(log.Fields{
			"msgnum": 504,
		}).Error(err)
	}
	jwtOpt.JWTDuration = JWTDuration

	v.SetConfigName("config")
	pwd, _ := os.Getwd()
	viewpath := pwd + filepath.FromSlash("/common")
	v.AddConfigPath(viewpath)

	if err := v.ReadInConfig(); err != nil {
		log.WithFields(log.Fields{
			"msgnum": 505,
		}).Error(err)
		os.Exit(1)
	}

	var rateOpt RateLimiterOptions
	if err := v.UnmarshalKey("ratelimit", &rateOpt); err != nil {
		log.WithFields(log.Fields{
			"msgnum": 506,
		}).Error(err)
	}

	var limit string
	if err := v.UnmarshalKey("limit", &limit); err != nil {
		log.WithFields(log.Fields{
			"msgnum": 507,
		}).Error(err)
	}

	var userOpt UserOptions
	if err := v.UnmarshalKey("useroptions", &userOpt); err != nil {
		log.WithFields(log.Fields{
			"msgnum": 508,
		}).Error(err)
	}

	return &redisOpt, db, redisClient, &oauth, mailer, &keyOpt, serverTLS, serverAddr, &jwtOpt, &rateOpt, limit, &userOpt, nil
}
