package common

import (
	"database/sql"
	"fmt"
	log "github.com/sirupsen/logrus"
)

// DBIntf - Interface to the Database
type DBIntf interface {
	DBClose() error
}

// DBService - Database type and Pointer to access Db
type DBService struct {
	DBType       string
	DB           *sql.DB
	LimitSQLRows string
}

// NewDBService - get connection to DB and create a DBService struct
func NewDBService(dbOpt *DBOptions) (*DBService, error) {

	var db *sql.DB
	var err error

	if dbOpt.DB == DBMysql {
		db, err = sql.Open(dbOpt.DB, fmt.Sprint(dbOpt.User, ":", dbOpt.Password, "@(", dbOpt.Host,
			":", dbOpt.Port, ")/", dbOpt.Schema, "?charset=utf8mb4&parseTime=True"))
		if err != nil {
			log.WithFields(log.Fields{
				"msgnum": 500,
			}).Error(err)
			return nil, err
		}
	} else if dbOpt.DB == DBPgsql {

	}
	// make sure connection is available
	err = db.Ping()
	if err != nil {
		log.WithFields(log.Fields{
			"msgnum": 501,
		}).Error(err)
		return nil, err
	}

	dbService := &DBService{}
	dbService.DBType = dbOpt.DB
	dbService.DB = db
	dbService.LimitSQLRows = dbOpt.LimitSQLRows

	return dbService, nil
}

// DBClose - Close connection to database
func (dbService *DBService) DBClose() error {
	err := dbService.DB.Close()
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}