package msgcontrollers

import (
	"net/http"

	"github.com/cloudfresco/vilom/common"
	"github.com/cloudfresco/vilom/msg/msgservices"
	"github.com/cloudfresco/vilom/user/userservices"
	"github.com/throttled/throttled/v2/store/goredisstore"
)

// Init the msg controllers
func Init(workspaceservice msgservices.WorkspaceServiceIntf, channelService msgservices.ChannelServiceIntf, msgService msgservices.MessageServiceIntf, userService userservices.UserServiceIntf, rateOpt *common.RateOptions, jwtOpt *common.JWTOptions, mux *http.ServeMux, store *goredisstore.GoRedisStore) {

	cc := NewWorkspaceController(workspaceservice, userService)
	tc := NewChannelController(channelService, userService)
	mc := NewMessageController(msgService, userService)

	hrlCat := common.GetHTTPRateLimiter(store, rateOpt.WorkspaceMaxRate, rateOpt.WorkspaceMaxBurst)
	hrlChannel := common.GetHTTPRateLimiter(store, rateOpt.ChannelMaxRate, rateOpt.ChannelMaxBurst)
	hrlMsg := common.GetHTTPRateLimiter(store, rateOpt.MsgMaxRate, rateOpt.MsgMaxBurst)
	mux.Handle("/v0.1/workspaces", common.AddMiddleware(hrlCat.RateLimit(cc),
		common.AuthenticateMiddleware,
		common.CorsMiddleware))
	mux.Handle("/v0.1/workspaces/", common.AddMiddleware(hrlCat.RateLimit(cc),
		common.AuthenticateMiddleware,
		common.CorsMiddleware))
	mux.Handle("/v0.1/channels/", common.AddMiddleware(hrlChannel.RateLimit(tc),
		common.AuthenticateMiddleware,
		common.CorsMiddleware))
	mux.Handle("/v0.1/messages/", common.AddMiddleware(hrlMsg.RateLimit(mc),
		common.AuthenticateMiddleware,
		common.CorsMiddleware))
}
