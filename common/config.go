package common

import (
	"os"
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
	DB                    string `mapstructure:"db"`
	Host                  string `mapstructure:"hostname"`
	Port                  string `mapstructure:"port"`
	User                  string `mapstructure:"user"`
	Password              string `mapstructure:"password"`
	Schema                string `mapstructure:"db_schema"`
	LimitSQLRows          string `mapstructure:"limit_sql_rows"`
	MySQLTestFilePath     string `mapstructure:"mysql_test_file_path"`
	MySQLSchemaFilePath   string `mapstructure:"mysql_schema_file_path"`
	MySQLTruncateFilePath string `mapstructure:"mysql_truncate_file_path"`
	PgSQLTestFilePath     string `mapstructure:"pgsql_test_file_path"`
	PgSQLSchemaFilePath   string `mapstructure:"pgsql_schema_file_path"`
	PgSQLTruncateFilePath string `mapstructure:"pgsql_truncate_file_path"`
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

// LogOptions - for logging
type LogOptions struct {
	Path  string `mapstructure:"log_file_path"`
	Level string `mapstructure:"log_level"`
}

// RoleOptions - for Role
type RoleOptions struct {
	Roles                 []Role `mapstructure:"roles"`
	RolesPolicyConfigPath string
	RolesTableName        string `mapstructure:"roles_table"`
}

// Role -  for user roles
type Role struct {
	PType string `mapstructure:"ptype"`
	V0    string `mapstructure:"v0"`
	V1    string `mapstructure:"v1"`
	V2    string `mapstructure:"v2"`
	V3    string `mapstructure:"v3"`
	V4    string `mapstructure:"v4"`
	V5    string `mapstructure:"v5"`
}

// GetDbConfig -- read DB config options
func GetDbConfig(v *viper.Viper) (*DBOptions, error) {
	var LimitSQLRows string

	dbOpt := DBOptions{}
	dbOpt.DB = v.GetString("VILOM_DB")
	dbOpt.Host = v.GetString("VILOM_DBHOST")
	dbOpt.Port = v.GetString("VILOM_DBPORT")
	dbOpt.User = v.GetString("VILOM_DBUSER")
	dbOpt.Password = v.GetString("VILOM_DBPASS")
	dbOpt.Schema = v.GetString("VILOM_DBNAME")
	dbOpt.MySQLTestFilePath = ""
	dbOpt.MySQLSchemaFilePath = ""
	dbOpt.MySQLTruncateFilePath = ""
	dbOpt.PgSQLTestFilePath = ""
	dbOpt.PgSQLSchemaFilePath = ""
	dbOpt.PgSQLTruncateFilePath = ""

	if err := v.UnmarshalKey("limit_sql_rows", &LimitSQLRows); err != nil {
		log.WithFields(log.Fields{
			"msgnum": 507,
		}).Error(err)
	}
	dbOpt.LimitSQLRows = LimitSQLRows

	return &dbOpt, nil
}

// GetRedisConfig -- read redis config options
func GetRedisConfig(v *viper.Viper) (*RedisOptions, error) {
	redisOpt := RedisOptions{}
	redisOpt.Addr = v.GetString("VILOM_REDIS_ADDRESS")
	return &redisOpt, nil
}

// GetMailerConfig -- read mailer config options
func GetMailerConfig(v *viper.Viper) (*MailerOptions, error) {
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

// GetServerConfig -- read server config options
func GetServerConfig(v *viper.Viper) (*ServerOptions, error) {
	serverOpt := ServerOptions{}
	serverOpt.ServerAddr = v.GetString("VILOM_SERVER_ADDRESS")
	serverOpt.ServerTLS = v.GetString("VILOM_SERVER_TLS")
	serverOpt.CaCertPath = v.GetString("VILOM_CA_CERT_PATH")
	serverOpt.CertPath = v.GetString("VILOM_CERT_PATH")
	serverOpt.KeyPath = v.GetString("VILOM_KEY_PATH")
	return &serverOpt, nil
}

// GetRateConfig -- read rate config options
func GetRateConfig(v *viper.Viper) (*RateOptions, error) {
	rateOpt := RateOptions{}
	if err := v.UnmarshalKey("rate_limit", &rateOpt); err != nil {
		log.WithFields(log.Fields{
			"msgnum": 506,
		}).Error(err)
		return nil, err
	}
	return &rateOpt, nil
}

// GetJWTConfig -- read JWT config options
func GetJWTConfig(v *viper.Viper) (*JWTOptions, error) {
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

// GetOauthConfig -- read oauth config options
func GetOauthConfig(v *viper.Viper) (*OauthOptions, error) {
	oauthOpt := OauthOptions{}
	oauthOpt.ClientID = v.GetString("GOOGLE_OAUTH2_CLIENT_ID")
	oauthOpt.ClientSecret = v.GetString("GOOGLE_OAUTH2_CLIENT_SECRET")
	return &oauthOpt, nil
}

// GetUserConfig -- read user config options
func GetUserConfig(v *viper.Viper) (*UserOptions, error) {
	userOpt := UserOptions{}
	if err := v.UnmarshalKey("user_options", &userOpt); err != nil {
		log.WithFields(log.Fields{
			"msgnum": 508,
		}).Error(err)
		return nil, err
	}
	return &userOpt, nil
}

// GetLogConfig -- read log config options
func GetLogConfig(v *viper.Viper) (*LogOptions, error) {
	logOpt := LogOptions{}
	logOpt.Path = v.GetString("VILOM_LOG_FILE_PATH")
	logOpt.Level = v.GetString("VILOM_LOG_LEVEL")
	return &logOpt, nil
}

// GetRoleConfig -- read Roles config options
func GetRoleConfig(v *viper.Viper) (*RoleOptions, error) {
	roleOpt := RoleOptions{}
	if err := v.Unmarshal(&roleOpt); err != nil {
		log.WithFields(log.Fields{
			"msgnum": 751,
		}).Error(err)
		return nil, err
	}
	roleOpt.RolesPolicyConfigPath = v.GetString("VILOM_ROLES_POLICY_CONFIG_PATH")
	return &roleOpt, nil
}

// GetViper -- init viper
func GetViper() (*viper.Viper, error) {
	v := viper.New()
	v.AutomaticEnv()

	v.SetConfigName("config")
	configFilePath := v.GetString("VILOM_CONFIG_FILE_PATH")
	v.AddConfigPath(configFilePath)

	if err := v.ReadInConfig(); err != nil {
		log.WithFields(log.Fields{
			"msgnum": 505,
		}).Error(err)
		return nil, err
	}
	return v, nil
}

// SetUpLogging -- create log file, and other log init
func SetUpLogging(logOpt *LogOptions) {
	var logFilePath string
	var logLevel log.Level
	var err error
	var f *os.File

	logFilePath = logOpt.Path

	switch logOpt.Level {
	case "PanicLevel":
		logLevel = log.PanicLevel
	case "FatalLevel":
		logLevel = log.FatalLevel
	case "ErrorLevel":
		logLevel = log.ErrorLevel
	case "WarnLevel":
		logLevel = log.WarnLevel
	case "InfoLevel":
		logLevel = log.InfoLevel
	case "DebugLevel":
		logLevel = log.DebugLevel
	case "TraceLevel":
		logLevel = log.TraceLevel
	default:
		logLevel = log.FatalLevel
	}

	// open the log file
	f, err = os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.WithFields(log.Fields{
			"msgnum": 100,
		}).Error(err)
		os.Exit(1)
	}

	log.SetOutput(f)
	log.SetFormatter(&log.JSONFormatter{})

	log.SetLevel(logLevel)
	return
}

// CreateDBService -- init DB
func CreateDBService(dbOpt *DBOptions) (*DBService, error) {
	dbService, err := NewDBService(dbOpt)
	if err != nil {
		log.WithFields(log.Fields{
			"msgnum": 750,
		}).Error(err)
		return nil, err
	}
	return dbService, nil
}

// CreateRedisService -- init redis
func CreateRedisService(redisOpt *RedisOptions) (*RedisService, error) {
	redisService, err := NewRedisService(redisOpt)
	if err != nil {
		log.WithFields(log.Fields{
			"msgnum": 750,
		}).Error(err)
		return nil, err
	}
	return redisService, nil
}

// CreateMailerService -- init mailer
func CreateMailerService(mailerOpt *MailerOptions) (*MailerService, error) {
	mailerService, err := NewMailerService(mailerOpt)
	if err != nil {
		log.WithFields(log.Fields{
			"msgnum": 750,
		}).Error(err)
		return nil, err
	}
	return mailerService, nil
}
