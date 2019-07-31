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
var mailerService *common.MailerService
var serverOpt *common.ServerOptions
var rateOpt *common.RateOptions
var jwtOpt *common.JWTOptions
var oauthOpt *common.OauthOptions
var userOpt *common.UserOptions
var Layout string

func TestMain(m *testing.M) {
	var err error

	dbService, redisService, mailerService, serverOpt, rateOpt, jwtOpt, oauthOpt, userOpt, err = testhelpers.InitTest()
	if err != nil {
		log.Fatal(err)
	}
	Layout = "2006-01-02T15:04:05Z"

	os.Exit(m.Run())

}
