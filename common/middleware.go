package common

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// AddMiddleware - adds middleware to a Handler
func AddMiddleware(h http.Handler, middleware ...func(http.Handler) http.Handler) http.Handler {
	for _, mw := range middleware {
		h = mw(h)
	}
	return h
}

// CorsMiddleware - Enable CORS with various options
func CorsMiddleware(next http.Handler) http.Handler {
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
func AuthenticateMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString, err := GetAuthBearerToken(r)
		if err != nil {
			log.WithFields(log.Fields{
				"msgnum": 751,
			}).Error(err)
			http.Error(w, "Error parsing token", http.StatusUnauthorized)
			return
		}
		jwtOpt := GetJWTOpt()
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// validate the alg is
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				http.Error(w, "Unexpected signing method", http.StatusUnauthorized)
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
			return jwtOpt.JWTKey, nil
		})
		v := ContextStruct{}
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			v.Email = claims["EmailAddr"].(string)
			v.TokenString = tokenString
		} else {
			log.WithFields(log.Fields{
				"msgnum": 752,
			}).Error(err)
			return
		}
		ctx := context.WithValue(r.Context(), KeyEmailToken, v)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
