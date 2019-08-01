package routes

import (
	"context"
	"fmt"
	"net/http"

	"github.com/blevesearch/bleve"
	"github.com/dgrijalva/jwt-go"

	log "github.com/sirupsen/logrus"
	"github.com/throttled/throttled"
	"github.com/throttled/throttled/store/goredisstore"

	"github.com/cloudfresco/vilom/common"

	"github.com/cloudfresco/vilom/msg/msgcontrollers"
	"github.com/cloudfresco/vilom/msg/msgservices"
	"github.com/cloudfresco/vilom/search/searchcontrollers"
	"github.com/cloudfresco/vilom/search/searchservices"
	"github.com/cloudfresco/vilom/user/usercontrollers"
	"github.com/cloudfresco/vilom/user/userservices"
)

/* error message range: 750-999 */

// AppState - holds pointers to services and options used through the
// journey of a request received from a client to the response
type AppState struct {
	DBService     *common.DBService
	RedisService  *common.RedisService
	MailerService *common.MailerService
	ServerOptions *common.ServerOptions
	RateOptions   *common.RateOptions
	JWTOptions    *common.JWTOptions
	OauthOptions  *common.OauthOptions
	UserOptions   *common.UserOptions
	SearchIndex   bleve.Index
}

// Init - Fill up AppState Struct
func (appState *AppState) Init(dbService *common.DBService, redisService *common.RedisService, mailerService *common.MailerService, serverOpt *common.ServerOptions, rateOpt *common.RateOptions, jwtOpt *common.JWTOptions, oauthOpt *common.OauthOptions, userOpt *common.UserOptions) error {

	appState.DBService = dbService
	appState.RedisService = redisService
	appState.MailerService = mailerService
	appState.ServerOptions = serverOpt
	appState.RateOptions = rateOpt
	appState.JWTOptions = jwtOpt
	appState.OauthOptions = oauthOpt
	appState.UserOptions = userOpt
	return nil
}

// AddMiddleware - adds middleware to a Handler
func AddMiddleware(h http.Handler, middleware ...func(http.Handler) http.Handler) http.Handler {
	for _, mw := range middleware {
		h = mw(h)
	}
	return h
}

// CorsMiddleware - Enable CORS with various options
func (appState *AppState) CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers",
				"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Allow-Origin")
		}
		// Stop here if its Preflighted OPTIONS request
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Expose-Headers", "Authorization")
			w.Header().Set("Access-Control-Max-Age", "86400")
			return
		}
		next.ServeHTTP(w, r)
	})
}

// AuthenticateMiddleware - Authenticate Token from request
func (appState *AppState) AuthenticateMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString, err := common.GetAuthBearerToken(r)
		if err != nil {
			log.WithFields(log.Fields{
				"msgnum": 751,
			}).Error(err)
			http.Error(w, "Error parsing token", http.StatusUnauthorized)
			return
		}
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// validate the alg is
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				http.Error(w, "Unexpected signing method", http.StatusUnauthorized)
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
			return appState.JWTOptions.JWTKey, nil
		})
		v := common.ContextStruct{}
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			v.Email = claims["EmailAddr"].(string)
			v.TokenString = tokenString
		} else {
			log.WithFields(log.Fields{
				"msgnum": 752,
			}).Error(err)
			return
		}
		ctx := context.WithValue(r.Context(), common.KeyEmailToken, v)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetHTTPRateLimiter - Get HTTP Rate Limiter
func (appState *AppState) GetHTTPRateLimiter(store *goredisstore.GoRedisStore, MaxRate int, MaxBurst int) throttled.HTTPRateLimiter {
	quota := throttled.RateQuota{MaxRate: throttled.PerMin(MaxRate), MaxBurst: MaxBurst}

	rateLimiter, err := throttled.NewGCRARateLimiter(store, quota)
	if err != nil {
		log.WithFields(log.Fields{
			"msgnum": 753,
		}).Error(err)
	}

	httpRateLimiter := throttled.HTTPRateLimiter{
		RateLimiter: rateLimiter,
		VaryBy:      &throttled.VaryBy{RemoteAddr: true},
	}
	return httpRateLimiter
}

// CreateServices create all the services
func (appState *AppState) CreateServices() (userservices.UserServiceIntf, userservices.UgroupServiceIntf, userservices.UbadgeServiceIntf, msgservices.CategoryServiceIntf, msgservices.TopicServiceIntf, msgservices.MessageServiceIntf, searchservices.SearchServiceIntf) {

	userService := userservices.NewUserService(appState.DBService, appState.RedisService, appState.MailerService, appState.JWTOptions, appState.UserOptions)
	ugroupService := userservices.NewUgroupService(appState.DBService, appState.RedisService)
	ubadgeService := userservices.NewUbadgeService(appState.DBService, appState.RedisService)
	catService := msgservices.NewCategoryService(appState.DBService, appState.RedisService)
	topicService := msgservices.NewTopicService(appState.DBService, appState.RedisService)
	msgService := msgservices.NewMessageService(appState.DBService, appState.RedisService)
	searchService := searchservices.NewSearchService(appState.DBService, appState.RedisService, appState.SearchIndex)

	return userService, ugroupService, ubadgeService, catService, topicService, msgService, searchService

}

// CreateControllers create all the controllers
func (appState *AppState) CreateControllers(userService userservices.UserServiceIntf, ugroupService userservices.UgroupServiceIntf, ubadgeService userservices.UbadgeServiceIntf, catService msgservices.CategoryServiceIntf, topicService msgservices.TopicServiceIntf, msgService msgservices.MessageServiceIntf, searchService searchservices.SearchServiceIntf) (*usercontrollers.UserController, *usercontrollers.UController, *usercontrollers.UgroupController, *usercontrollers.UbadgeController, *msgcontrollers.CategoryController, *msgcontrollers.TopicController, *msgcontrollers.MessageController, *searchcontrollers.SearchController) {

	usc := usercontrollers.NewUserController(userService)
	uc := usercontrollers.NewUController(userService)
	ugc := usercontrollers.NewUgroupController(ugroupService, userService)
	ubc := usercontrollers.NewUbadgeController(ubadgeService, userService)

	cc := msgcontrollers.NewCategoryController(catService, userService)
	tc := msgcontrollers.NewTopicController(topicService, userService)
	mc := msgcontrollers.NewMessageController(msgService, userService)

	sc := searchcontrollers.NewSearchController(searchService, userService)

	return usc, uc, ugc, ubc, cc, tc, mc, sc
}

// CreateRateLimiters create all the rate limiters
func (appState *AppState) CreateRateLimiters() (throttled.HTTPRateLimiter, throttled.HTTPRateLimiter, throttled.HTTPRateLimiter, throttled.HTTPRateLimiter, throttled.HTTPRateLimiter, throttled.HTTPRateLimiter, throttled.HTTPRateLimiter, throttled.HTTPRateLimiter) {

	store, err := goredisstore.New(appState.RedisService.RedisClient, "throttled:")
	if err != nil {
		log.WithFields(log.Fields{
			"msgnum": 754,
		}).Error(err)
	}

	hrlUser := appState.GetHTTPRateLimiter(store, appState.RateOptions.UserMaxRate, appState.RateOptions.UserMaxBurst)
	hrlU := appState.GetHTTPRateLimiter(store, appState.RateOptions.UMaxRate, appState.RateOptions.UMaxBurst)
	hrlUgroup := appState.GetHTTPRateLimiter(store, appState.RateOptions.UgroupMaxRate, appState.RateOptions.UgroupMaxBurst)
	hrlUbadge := appState.GetHTTPRateLimiter(store, appState.RateOptions.UbadgeMaxRate, appState.RateOptions.UbadgeMaxBurst)

	hrlCat := appState.GetHTTPRateLimiter(store, appState.RateOptions.CatMaxRate, appState.RateOptions.CatMaxBurst)
	hrlTopic := appState.GetHTTPRateLimiter(store, appState.RateOptions.TopicMaxRate, appState.RateOptions.TopicMaxBurst)
	hrlMsg := appState.GetHTTPRateLimiter(store, appState.RateOptions.MsgMaxRate, appState.RateOptions.MsgMaxBurst)

	hrlSearch := appState.GetHTTPRateLimiter(store, appState.RateOptions.SearchMaxRate, appState.RateOptions.SearchMaxBurst)

	return hrlUser, hrlU, hrlUgroup, hrlUbadge, hrlCat, hrlTopic, hrlMsg, hrlSearch
}

// CreateRoutes - Create all the routes
func (appState *AppState) CreateRoutes() *http.ServeMux {

	userService, ugroupService, ubadgeService, catService, topicService, msgService, searchService := appState.CreateServices()

	usc, uc, ugc, ubc, cc, tc, mc, sc := appState.CreateControllers(userService, ugroupService, ubadgeService, catService, topicService, msgService, searchService)

	hrlUser, hrlU, hrlUgroup, hrlUbadge, hrlCat, hrlTopic, hrlMsg, hrlSearch := appState.CreateRateLimiters()

	mux := http.NewServeMux()
	mux.Handle("/v0.1/users", AddMiddleware(hrlUser.RateLimit(usc),
		appState.AuthenticateMiddleware,
		appState.CorsMiddleware))
	mux.Handle("/v0.1/users/", AddMiddleware(hrlUser.RateLimit(usc),
		appState.AuthenticateMiddleware,
		appState.CorsMiddleware))
	mux.Handle("/v0.1/u/", AddMiddleware(hrlU.RateLimit(uc), appState.CorsMiddleware))
	mux.Handle("/v0.1/ugroups/", AddMiddleware(hrlUgroup.RateLimit(ugc),
		appState.AuthenticateMiddleware,
		appState.CorsMiddleware))
	mux.Handle("/v0.1/ubadges/", AddMiddleware(hrlUbadge.RateLimit(ubc),
		appState.AuthenticateMiddleware,
		appState.CorsMiddleware))
	mux.Handle("/v0.1/categories/", AddMiddleware(hrlCat.RateLimit(cc),
		appState.AuthenticateMiddleware,
		appState.CorsMiddleware))
	mux.Handle("/v0.1/topics/", AddMiddleware(hrlTopic.RateLimit(tc),
		appState.AuthenticateMiddleware,
		appState.CorsMiddleware))
	mux.Handle("/v0.1/messages/", AddMiddleware(hrlMsg.RateLimit(mc),
		appState.AuthenticateMiddleware,
		appState.CorsMiddleware))
	mux.Handle("/v0.1/search/", AddMiddleware(hrlSearch.RateLimit(sc),
		appState.AuthenticateMiddleware,
		appState.CorsMiddleware))

	return mux
}
