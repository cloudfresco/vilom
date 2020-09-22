package msgservices

import (
	"log"
	"os"
	"testing"

	"github.com/cloudfresco/vilom/common"
	"github.com/cloudfresco/vilom/testhelpers"
)

var dbService *common.DBService
var redisService *common.RedisService
var serverOpt *common.ServerOptions
var userOpt *common.UserOptions
var roleOpt *common.RoleOptions
var Layout string

func TestMain(m *testing.M) {
	var err error

	dbService, redisService, serverOpt, userOpt, roleOpt, err = testhelpers.InitTest()
	if err != nil {
		log.Fatal(err)
	}
	Layout = "2006-01-02T15:04:05Z"

	os.Exit(m.Run())

}
