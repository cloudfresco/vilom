package common

import (
	log "github.com/sirupsen/logrus"
	"github.com/throttled/throttled"
	"github.com/throttled/throttled/store/goredisstore"
)

// GetHTTPRateLimiter - Get HTTP Rate Limiter
func GetHTTPRateLimiter(store *goredisstore.GoRedisStore, MaxRate int, MaxBurst int) throttled.HTTPRateLimiter {
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
