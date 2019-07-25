package common

import (
	"os"
	"path/filepath"
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

/* error message range: 500-749 */

// DBMysql for DbType is mysql
const DBMysql string = "mysql"

// DBPgsql for DbType is pgsql
const DBPgsql string = "pgsql"

// DBOptions - for db config
type DBOptions struct {
	DB           string `mapstructure:"db"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	Host         string `mapstructure:"hostname"`
	Port         string `mapstructure:"port"`
	Schema       string `mapstructure:"db_schema"`
	LimitSQLRows string `mapstructure:"limit_sql_rows"`
}

// RedisOptions - for redis config
type RedisOptions struct {
	Addr string `mapstructure:"addr"`
}

// MailerOptions - for mailer config
type MailerOptions struct {
	User     string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Port     int    `mapstructure:"port"`
	Server   string `mapstructure:"server"`
}

// ServerOptions - for server config
type ServerOptions struct {
	ServerAddr string `mapstructure:"server_addr"`
	ServerTLS  string `mapstructure:"server_tls"`
	CaCertPath string `mapstructure:"ca_cert_path"`
	CertPath   string `mapstructure:"cert_path"`
	KeyPath    string `mapstructure:"key_path"`
}

// RateOptions - for rate limiting requests
type RateOptions struct {
	UserMaxRate    int `mapstructure:"user_max_rate"`
	UserMaxBurst   int `mapstructure:"user_max_burst"`
	UgroupMaxRate  int `mapstructure:"ugroup_max_rate"`
	UgroupMaxBurst int `mapstructure:"ugroup_max_burst"`
	CatMaxRate     int `mapstructure:"cat_max_rate"`
	CatMaxBurst    int `mapstructure:"cat_max_burst"`
	TopicMaxRate   int `mapstructure:"topic_max_rate"`
	TopicMaxBurst  int `mapstructure:"topic_max_burst"`
	MsgMaxRate     int `mapstructure:"msg_max_rate"`
	MsgMaxBurst    int `mapstructure:"msg_max_burst"`
	UbadgeMaxRate  int `mapstructure:"ubadge_max_rate"`
	UbadgeMaxBurst int `mapstructure:"ubadge_max_burst"`
	SearchMaxRate  int `mapstructure:"search_max_rate"`
	SearchMaxBurst int `mapstructure:"search_max_burst"`
	UMaxRate       int `mapstructure:"u_max_rate"`
	UMaxBurst      int `mapstructure:"u_max_burst"`
}

// JWTOptions - for JWT config
type JWTOptions struct {
	JWTKey      []byte
	JWTDuration int
}

// OauthOptions - for oauth config
type OauthOptions struct {
	ClientID     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
}

// UserOptions - for user login
type UserOptions struct {
	ConfirmTokenDuration string `mapstructure:"confirm_token_duration"`
	ResetTokenDuration   string `mapstructure:"reset_token_duration"`
}

func getDbConfig(v *viper.Viper) (*DBOptions, error) {
	var LimitSQLRows string

	dbOpt := DBOptions{}
	dbOpt.DB = v.GetString("VILOM_DB")
	dbOpt.User = v.GetString("VILOM_DBUSER")
	dbOpt.Password = v.GetString("VILOM_DBPASS")
	dbOpt.Host = v.GetString("VILOM_DBHOST")
	dbOpt.Port = v.GetString("VILOM_DBPORT")
	dbOpt.Schema = v.GetString("VILOM_DBNAME")

	if err := v.UnmarshalKey("limit_sql_rows", &LimitSQLRows); err != nil {
		log.WithFields(log.Fields{
			"msgnum": 507,
		}).Error(err)
	}
	dbOpt.LimitSQLRows = LimitSQLRows

	return &dbOpt, nil
}

func getRedisConfig(v *viper.Viper) (*RedisOptions, error) {
	redisOpt := RedisOptions{}
	redisOpt.Addr = v.GetString("VILOM_REDIS_ADDRESS")
	return &redisOpt, nil
}

func getMailerConfig(v *viper.Viper) (*MailerOptions, error) {
	mailerOpt := MailerOptions{}
	mailerOpt.Server = v.GetString("VILOM_MAILER_SERVER")
	MailerPort, err := strconv.Atoi(v.GetString("VILOM_MAILER_PORT"))
	if err != nil {
		log.WithFields(log.Fields{
			"msgnum": 503,
		}).Error(err)
		return nil, err
	}
	mailerOpt.Port = MailerPort
	mailerOpt.User = v.GetString("VILOM_MAILER_USER")
	mailerOpt.Password = v.GetString("VILOM_MAILER_PASS")
	return &mailerOpt, nil
}

func getServerConfig(v *viper.Viper) (*ServerOptions, error) {
	serverOpt := ServerOptions{}
	serverOpt.ServerAddr = v.GetString("VILOM_SERVER_ADDRESS")
	serverOpt.ServerTLS = v.GetString("VILOM_SERVER_TLS")
	serverOpt.CaCertPath = v.GetString("VILOM_CA_CERT_PATH")
	serverOpt.CertPath = v.GetString("VILOM_CERT_PATH")
	serverOpt.KeyPath = v.GetString("VILOM_KEY_PATH")
	return &serverOpt, nil
}

func getRateConfig(v *viper.Viper) (*RateOptions, error) {
	rateOpt := RateOptions{}
	if err := v.UnmarshalKey("rate_limit", &rateOpt); err != nil {
		log.WithFields(log.Fields{
			"msgnum": 506,
		}).Error(err)
		return nil, err
	}
	return &rateOpt, nil
}

func getJWTConfig(v *viper.Viper) (*JWTOptions, error) {
	var err error

	jwtOpt := JWTOptions{}
	jwtOpt.JWTKey = []byte(v.GetString("VILOM_JWT_KEY"))
	jwtOpt.JWTDuration, err = strconv.Atoi(v.GetString("VILOM_JWT_DURATION"))
	if err != nil {
		log.WithFields(log.Fields{
			"msgnum": 504,
		}).Error(err)
		return nil, err
	}
	return &jwtOpt, nil
}

func getOauthConfig(v *viper.Viper) (*OauthOptions, error) {
	oauthOpt := OauthOptions{}
	oauthOpt.ClientID = v.GetString("GOOGLE_OAUTH2_CLIENT_ID")
	oauthOpt.ClientSecret = v.GetString("GOOGLE_OAUTH2_CLIENT_SECRET")
	return &oauthOpt, nil
}

func getUserConfig(v *viper.Viper) (*UserOptions, error) {
	userOpt := UserOptions{}
	if err := v.UnmarshalKey("user_options", &userOpt); err != nil {
		log.WithFields(log.Fields{
			"msgnum": 508,
		}).Error(err)
		return nil, err
	}
	return &userOpt, nil
}

// GetConfig - Bring in Configuration info from ENV and Confi.json
func GetConfig() (*DBOptions, *RedisOptions, *MailerOptions, *ServerOptions, *RateOptions, *JWTOptions, *OauthOptions, *UserOptions, error) {

	v := viper.New()
	v.AutomaticEnv()

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

	dbOpt, err := getDbConfig(v)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, err
	}

	redisOpt, err := getRedisConfig(v)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, err
	}

	mailerOpt, err := getMailerConfig(v)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, err
	}

	serverOpt, err := getServerConfig(v)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, err
	}

	rateOpt, err := getRateConfig(v)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, err
	}

	jwtOpt, err := getJWTConfig(v)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, err
	}

	oauthOpt, err := getOauthConfig(v)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, err
	}

	userOpt, err := getUserConfig(v)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, err
	}

	return dbOpt, redisOpt, mailerOpt, serverOpt, rateOpt, jwtOpt, oauthOpt, userOpt, nil
}
