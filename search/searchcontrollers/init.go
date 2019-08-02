package searchcontrollers

import (
	//log "github.com/sirupsen/logrus"
	"net/http"

	"github.com/cloudfresco/vilom/common"
	"github.com/cloudfresco/vilom/search/searchservices"
	"github.com/cloudfresco/vilom/user/userservices"
	"github.com/throttled/throttled/store/goredisstore"
)

// Init the search controller
func Init(searchService searchservices.SearchServiceIntf, userService userservices.UserServiceIntf, rateOpt *common.RateOptions, jwtOpt *common.JWTOptions, mux *http.ServeMux, store *goredisstore.GoRedisStore) *SearchController {

	sc := NewSearchController(searchService, userService)

	hrlSearch := common.GetHTTPRateLimiter(store, rateOpt.SearchMaxRate, rateOpt.SearchMaxBurst)

	mux.Handle("/v0.1/search/", common.AddMiddleware(hrlSearch.RateLimit(sc),
		common.AuthenticateMiddleware,
		common.CorsMiddleware))
	return sc
}
