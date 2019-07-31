package testhelpers

import (
	"context"
	"github.com/cloudfresco/vilom/common"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"strings"
)

func getTestDbConfig(v *viper.Viper) (*common.DBOptions, error) {
	var LimitSQLRows string

	dbOpt := common.DBOptions{}
	dbOpt.DB = v.GetString("VILOM_DB")
	dbOpt.Host = v.GetString("VILOM_DBHOST")
	dbOpt.Port = v.GetString("VILOM_DBPORT")
	dbOpt.User = v.GetString("VILOM_DBUSER_TEST")
	dbOpt.Password = v.GetString("VILOM_DBPASS_TEST")
	dbOpt.Schema = v.GetString("VILOM_DBNAME_TEST")
	dbOpt.MySQLTestFilePath = v.GetString("VILOM_DBSQL_MYSQL_TEST")
	dbOpt.PgSQLTestFilePath = v.GetString("VILOM_DBSQL_PGSQL_TEST")

	if err := v.UnmarshalKey("limit_sql_rows", &LimitSQLRows); err != nil {
		log.WithFields(log.Fields{
			"msgnum": 507,
		}).Error(err)
	}
	dbOpt.LimitSQLRows = LimitSQLRows

	return &dbOpt, nil
}

func getTestConfigOpt() (*common.DBOptions, *common.RedisOptions, *common.MailerOptions, *common.ServerOptions, *common.RateOptions, *common.JWTOptions, *common.OauthOptions, *common.UserOptions, *common.LogOptions) {

	v, err := common.GetViper()
	if err != nil {
		log.WithFields(log.Fields{
			"msgnum": 103,
		}).Error(err)
		os.Exit(1)
	}

	dbOpt, err := getTestDbConfig(v)
	if err != nil {
		log.WithFields(log.Fields{
			"msgnum": 103,
		}).Error(err)
		os.Exit(1)
	}

	redisOpt, err := common.GetRedisConfig(v)
	if err != nil {
		log.WithFields(log.Fields{
			"msgnum": 103,
		}).Error(err)
		os.Exit(1)
	}

	mailerOpt, err := common.GetMailerConfig(v)
	if err != nil {
		log.WithFields(log.Fields{
			"msgnum": 103,
		}).Error(err)
		os.Exit(1)
	}

	serverOpt, err := common.GetServerConfig(v)
	if err != nil {
		log.WithFields(log.Fields{
			"msgnum": 103,
		}).Error(err)
		os.Exit(1)
	}

	rateOpt, err := common.GetRateConfig(v)
	if err != nil {
		log.WithFields(log.Fields{
			"msgnum": 103,
		}).Error(err)
		os.Exit(1)
	}

	jwtOpt, err := common.GetJWTConfig(v)
	if err != nil {
		log.WithFields(log.Fields{
			"msgnum": 103,
		}).Error(err)
		os.Exit(1)
	}

	oauthOpt, err := common.GetOauthConfig(v)
	if err != nil {
		log.WithFields(log.Fields{
			"msgnum": 103,
		}).Error(err)
		os.Exit(1)
	}

	userOpt, err := common.GetUserConfig(v)
	if err != nil {
		log.WithFields(log.Fields{
			"msgnum": 103,
		}).Error(err)
		os.Exit(1)
	}

	logOpt, err := common.GetLogConfig(v)
	if err != nil {
		log.WithFields(log.Fields{
			"msgnum": 103,
		}).Error(err)
		os.Exit(1)
	}

	return dbOpt, redisOpt, mailerOpt, serverOpt, rateOpt, jwtOpt, oauthOpt, userOpt, logOpt
}

// InitTest - used for initialization of the test DB etc.
func InitTest() (*common.DBService, *common.RedisService, *common.MailerService, *common.ServerOptions, *common.RateOptions, *common.JWTOptions, *common.OauthOptions, *common.UserOptions, error) {

	dbOpt, redisOpt, mailerOpt, serverOpt, rateOpt, jwtOpt, oauthOpt, userOpt, logOpt := getTestConfigOpt()

	common.SetUpLogging(logOpt)

	dbService, err := common.CreateDBService(dbOpt)
	if err != nil {
		log.WithFields(log.Fields{
			"msgnum": 750,
		}).Error(err)
		return nil, nil, nil, nil, nil, nil, nil, nil, err
	}

	redisService, err := common.CreateRedisService(redisOpt)
	if err != nil {
		log.WithFields(log.Fields{
			"msgnum": 750,
		}).Error(err)
		return nil, nil, nil, nil, nil, nil, nil, nil, err
	}

	mailerService, err := common.CreateMailerService(mailerOpt)
	if err != nil {
		log.WithFields(log.Fields{
			"msgnum": 750,
		}).Error(err)
		return nil, nil, nil, nil, nil, nil, nil, nil, err
	}

	return dbService, redisService, mailerService, serverOpt, rateOpt, jwtOpt, oauthOpt, userOpt, nil

}

// LoadSQL -- load data into all tables
func LoadSQL(dbService *common.DBService) error {
	ctx := context.Background()
	content, err := ioutil.ReadFile(dbService.MySQLTestFilePath)

	if err != nil {
		log.Println(err)
		return err
	}

	sqlLines := strings.Split(string(content), ";\n")

	for _, sqlLine := range sqlLines {

		if sqlLine != "" {
			_, err := dbService.DB.ExecContext(ctx, sqlLine)
			if err != nil {
				log.Println(err)
				return err
			}
		}
	}
	return nil
}

// DeleteSQL -- delete data from all tables
func DeleteSQL(dbService *common.DBService) error {
	ctx := context.Background()
	tables := []string{}
	tableSchema := dbService.Schema
	sql := "select table_name from information_schema.tables where table_schema = " + " '" + tableSchema + "' " + ";"
	rows, err := dbService.DB.QueryContext(ctx, sql)
	if err != nil {
		log.Println(err)
		return err
	}
	var tableName string
	for rows.Next() {
		err = rows.Scan(&tableName)
		if err != nil {
			log.Println(err)
			err = rows.Close()
			if err != nil {
				log.Println(err)
				return err
			}
			return err
		}
		tables = append(tables, tableName)
	}
	err = rows.Close()
	if err != nil {
		log.Println(err)
		return err
	}

	err = rows.Err()
	if err != nil {
		log.Println(err)
		return err
	}

	for _, tableName := range tables {
		sql = "truncate " + tableName
		_, err := dbService.DB.ExecContext(ctx, sql)
		if err != nil {
			log.Println(err)
			return err
		}
	}

	return nil
}
