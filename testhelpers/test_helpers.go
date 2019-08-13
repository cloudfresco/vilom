package testhelpers

import (
	"context"
	"database/sql"
	"github.com/cloudfresco/vilom/common"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io/ioutil"
	"strconv"
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
	dbOpt.MySQLSchemaFilePath = v.GetString("VILOM_DBSQL_MYSQL_SCHEMA")
	dbOpt.MySQLTruncateFilePath = v.GetString("VILOM_DBSQL_MYSQL_TRUNCATE")
	dbOpt.PgSQLTestFilePath = v.GetString("VILOM_DBSQL_PGSQL_TEST")
	dbOpt.PgSQLSchemaFilePath = v.GetString("VILOM_DBSQL_PGSQL_SCHEMA")
	dbOpt.PgSQLTruncateFilePath = v.GetString("VILOM_DBSQL_PGSQL_TRUNCATE")

	if err := v.UnmarshalKey("limit_sql_rows", &LimitSQLRows); err != nil {
		log.Fatal(err)
	}
	dbOpt.LimitSQLRows = LimitSQLRows

	return &dbOpt, nil
}

func getTestJWTConfig(v *viper.Viper) (*common.JWTOptions, error) {
	var err error

	jwtOpt := common.JWTOptions{}
	jwtOpt.JWTKey = []byte(v.GetString("VILOM_JWT_KEY_TEST"))
	jwtOpt.JWTDuration, err = strconv.Atoi(v.GetString("VILOM_JWT_DURATION_TEST"))
	if err != nil {
		log.Fatal(err)
	}
	return &jwtOpt, nil
}

func getTestConfigOpt() (*common.DBOptions, *common.RedisOptions, *common.ServerOptions, *common.UserOptions, *common.LogOptions) {

	v, err := common.GetViper()
	if err != nil {
		log.Fatal(err)
	}

	dbOpt, err := getTestDbConfig(v)
	if err != nil {
		log.Fatal(err)
	}

	redisOpt, err := common.GetRedisConfig(v)
	if err != nil {
		log.Fatal(err)
	}

	serverOpt, err := common.GetServerConfig(v)
	if err != nil {
		log.Fatal(err)
	}

	userOpt, err := common.GetUserConfig(v)
	if err != nil {
		log.Fatal(err)
	}

	logOpt, err := common.GetLogConfig(v)
	if err != nil {
		log.Fatal(err)
	}

	return dbOpt, redisOpt, serverOpt, userOpt, logOpt
}

// InitTest - used for initialization of the test DB etc.
func InitTest() (*common.DBService, *common.RedisService, *common.ServerOptions, *common.UserOptions, error) {

	dbOpt, redisOpt, serverOpt, userOpt, logOpt := getTestConfigOpt()

	common.SetUpLogging(logOpt)

	dbService, err := common.CreateDBService(dbOpt)
	if err != nil {
		log.Fatal(err)
	}

	redisService, err := common.CreateRedisService(redisOpt)
	if err != nil {
		log.Fatal(err)
	}

	return dbService, redisService, serverOpt, userOpt, nil

}

func getTestConfigOptController() (*common.DBOptions, *common.RedisOptions, *common.ServerOptions, *common.RateOptions, *common.JWTOptions, *common.OauthOptions, *common.UserOptions, *common.LogOptions) {

	v, err := common.GetViper()
	if err != nil {
		log.Fatal(err)
	}

	dbOpt, err := getTestDbConfig(v)
	if err != nil {
		log.Fatal(err)
	}

	redisOpt, err := common.GetRedisConfig(v)
	if err != nil {
		log.Fatal(err)
	}

	serverOpt, err := common.GetServerConfig(v)
	if err != nil {
		log.Fatal(err)
	}

	rateOpt, err := common.GetRateConfig(v)
	if err != nil {
		log.Fatal(err)
	}

	jwtOpt, err := getTestJWTConfig(v)
	if err != nil {
		log.Fatal(err)
	}

	oauthOpt, err := common.GetOauthConfig(v)
	if err != nil {
		log.Fatal(err)
	}

	userOpt, err := common.GetUserConfig(v)
	if err != nil {
		log.Fatal(err)
	}

	logOpt, err := common.GetLogConfig(v)
	if err != nil {
		log.Fatal(err)
	}

	return dbOpt, redisOpt, serverOpt, rateOpt, jwtOpt, oauthOpt, userOpt, logOpt
}

// InitTestController - used for initialization of the test controllers
func InitTestController() (*common.DBService, *common.RedisService, *common.ServerOptions, *common.RateOptions, *common.JWTOptions, *common.OauthOptions, *common.UserOptions, error) {

	dbOpt, redisOpt, serverOpt, rateOpt, jwtOpt, oauthOpt, userOpt, logOpt := getTestConfigOptController()

	common.SetUpLogging(logOpt)
	common.SetJWTOpt(jwtOpt)

	dbService, err := common.CreateDBService(dbOpt)
	if err != nil {
		log.Fatal(err)
	}

	redisService, err := common.CreateRedisService(redisOpt)
	if err != nil {
		log.Fatal(err)
	}

	return dbService, redisService, serverOpt, rateOpt, jwtOpt, oauthOpt, userOpt, nil

}

// LoadSQL -- drop db, create db, use db, load data
func LoadSQL(dbService *common.DBService) error {
	var err error
	ctx := context.Background()

	if dbService.DBType == common.DBMysql {
		err = execSQLFile(ctx, dbService.MySQLTruncateFilePath, dbService.DB)
		if err != nil {
			log.Println(err)
			return err
		}

		err = execSQLFile(ctx, dbService.MySQLTestFilePath, dbService.DB)
		if err != nil {
			log.Println(err)
			return err
		}

	} else if dbService.DBType == common.DBPgsql {
		err = execSQLFile(ctx, dbService.PgSQLTruncateFilePath, dbService.DB)
		if err != nil {
			log.Println(err)
			return err
		}

		err = execSQLFile(ctx, dbService.PgSQLTestFilePath, dbService.DB)
		if err != nil {
			log.Println(err)
			return err
		}
	}

	return nil
}

func execSQLFile(ctx context.Context, sqlFilePath string, db *sql.DB) error {

	content, err := ioutil.ReadFile(sqlFilePath)

	if err != nil {
		log.Println(err)
		return err
	}

	tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		log.Fatal(err)
	}

	sqlLines := strings.Split(string(content), ";\n")

	for _, sqlLine := range sqlLines {

		if sqlLine != "" {
			_, err := tx.ExecContext(ctx, sqlLine)
			if err != nil {
				log.Println(err)
				if rollbackErr := tx.Rollback(); rollbackErr != nil {
					log.Printf("Load SQL failed: %v, unable to rollback: %v\n", err, rollbackErr)
					return err
				}
			}
		}
	}
	err = tx.Commit()
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
