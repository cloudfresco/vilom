package usercontrollers

import (
	"net/http"

	"github.com/cloudfresco/vilom/common"
	"github.com/cloudfresco/vilom/user/userservices"
	"github.com/throttled/throttled/store/goredisstore"
)

// Init the user controllers
func Init(userService userservices.UserServiceIntf, ugroupService userservices.UgroupServiceIntf, ubadgeService userservices.UbadgeServiceIntf, rateOpt *common.RateOptions, jwtOpt *common.JWTOptions, mux *http.ServeMux, store *goredisstore.GoRedisStore) {

	usc := NewUserController(userService)
	uc := NewUController(userService)
	ugc := NewUgroupController(ugroupService, userService)
	ubc := NewUbadgeController(ubadgeService, userService)

	hrlUser := common.GetHTTPRateLimiter(store, rateOpt.UserMaxRate, rateOpt.UserMaxBurst)
	hrlU := common.GetHTTPRateLimiter(store, rateOpt.UMaxRate, rateOpt.UMaxBurst)
	hrlUgroup := common.GetHTTPRateLimiter(store, rateOpt.UgroupMaxRate, rateOpt.UgroupMaxBurst)
	hrlUbadge := common.GetHTTPRateLimiter(store, rateOpt.UbadgeMaxRate, rateOpt.UbadgeMaxBurst)

	mux.Handle("/v0.1/users", common.AddMiddleware(hrlUser.RateLimit(usc),
		common.AuthenticateMiddleware,
		common.CorsMiddleware))
	mux.Handle("/v0.1/users/", common.AddMiddleware(hrlUser.RateLimit(usc),
		common.AuthenticateMiddleware,
		common.CorsMiddleware))
	mux.Handle("/v0.1/u/", common.AddMiddleware(hrlU.RateLimit(uc), common.CorsMiddleware))
	mux.Handle("/v0.1/ugroups", common.AddMiddleware(hrlUgroup.RateLimit(ugc),
		common.AuthenticateMiddleware,
		common.CorsMiddleware))
	mux.Handle("/v0.1/ugroups/", common.AddMiddleware(hrlUgroup.RateLimit(ugc),
		common.AuthenticateMiddleware,
		common.CorsMiddleware))
	mux.Handle("/v0.1/ubadges", common.AddMiddleware(hrlUbadge.RateLimit(ubc),
		common.AuthenticateMiddleware,
		common.CorsMiddleware))
	mux.Handle("/v0.1/ubadges/", common.AddMiddleware(hrlUbadge.RateLimit(ubc),
		common.AuthenticateMiddleware,
		common.CorsMiddleware))
}
