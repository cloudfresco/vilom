package msgcontrollers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/casbin/casbin/v2"
	"github.com/cloudfresco/vilom/common"
	"github.com/cloudfresco/vilom/msg/msgservices"
	"github.com/cloudfresco/vilom/testhelpers"
	"github.com/cloudfresco/vilom/user/usercontrollers"
	"github.com/cloudfresco/vilom/user/userservices"

	"github.com/throttled/throttled/v2/store/goredisstore"
)

var dbService *common.DBService
var redisService *common.RedisService
var mailerService *common.MailerService
var serverOpt *common.ServerOptions
var rateOpt *common.RateOptions
var jwtOpt *common.JWTOptions
var userOpt *common.UserOptions
var Layout string
var mux *http.ServeMux
var authEnforcer *casbin.Enforcer

func TestMain(m *testing.M) {
	var err error

	dbService, redisService, serverOpt, rateOpt, jwtOpt, _, userOpt, authEnforcer, err = testhelpers.InitTestController()
	if err != nil {
		log.Println(err)
		return
	}
	Layout = "2006-01-02T15:04:05Z"

	catService := msgservices.NewCategoryService(dbService, redisService)
	topicService := msgservices.NewTopicService(dbService, redisService)
	msgService := msgservices.NewMessageService(dbService, redisService)
	userService := userservices.NewUserService(dbService, redisService, mailerService, jwtOpt, userOpt, authEnforcer)
	ugroupService := userservices.NewUgroupService(dbService, redisService)
	ubadgeService := userservices.NewUbadgeService(dbService, redisService)
	store, err := goredisstore.New(redisService.RedisClient, "throttled:")
	if err != nil {
		log.Println(err)
		return
	}

	mux = http.NewServeMux()
	Init(catService, topicService, msgService, userService, rateOpt, jwtOpt, mux, store)
	usercontrollers.Init(userService, ugroupService, ubadgeService, rateOpt, jwtOpt, mux, store)
	os.Exit(m.Run())
}

func LoginUser() string {
	w := httptest.NewRecorder()

	req, err := http.NewRequest("POST", "http://localhost:8000/v0.1/u/login", bytes.NewBuffer([]byte(`{"Email": "abcd145@gmail.com", "Password": "abc1238"}`)))
	if err != nil {
		log.Fatal(err)
		return ""
	}
	mux.ServeHTTP(w, req)

	user := userservices.User{}
	dec := json.NewDecoder(strings.NewReader(w.Body.String()))
	err = dec.Decode(&user)
	if err != nil {
		log.Println(err)
		return ""
	}
	tokenstring := user.Tokenstring
	return tokenstring
}
