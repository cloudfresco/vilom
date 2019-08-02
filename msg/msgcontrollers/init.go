package msgcontrollers

import (
	//log "github.com/sirupsen/logrus"
	"net/http"

	"github.com/cloudfresco/vilom/common"
	"github.com/cloudfresco/vilom/msg/msgservices"
	"github.com/cloudfresco/vilom/user/userservices"
	"github.com/throttled/throttled/store/goredisstore"
)

// Init the msg controllers
func Init(catService msgservices.CategoryServiceIntf, topicService msgservices.TopicServiceIntf, msgService msgservices.MessageServiceIntf, userService userservices.UserServiceIntf, rateOpt *common.RateOptions, jwtOpt *common.JWTOptions, mux *http.ServeMux, store *goredisstore.GoRedisStore) (*CategoryController, *TopicController, *MessageController) {

	cc := NewCategoryController(catService, userService)
	tc := NewTopicController(topicService, userService)
	mc := NewMessageController(msgService, userService)

	hrlCat := common.GetHTTPRateLimiter(store, rateOpt.CatMaxRate, rateOpt.CatMaxBurst)
	hrlTopic := common.GetHTTPRateLimiter(store, rateOpt.TopicMaxRate, rateOpt.TopicMaxBurst)
	hrlMsg := common.GetHTTPRateLimiter(store, rateOpt.MsgMaxRate, rateOpt.MsgMaxBurst)
	mux.Handle("/v0.1/categories/", common.AddMiddleware(hrlCat.RateLimit(cc),
		common.AuthenticateMiddleware,
		common.CorsMiddleware))
	mux.Handle("/v0.1/topics/", common.AddMiddleware(hrlTopic.RateLimit(tc),
		common.AuthenticateMiddleware,
		common.CorsMiddleware))
	mux.Handle("/v0.1/messages/", common.AddMiddleware(hrlMsg.RateLimit(mc),
		common.AuthenticateMiddleware,
		common.CorsMiddleware))
	return cc, tc, mc
}
