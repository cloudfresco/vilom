package common

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/go-redis/redis"
)

// ContextData - details of a user stored in the Redis cache
type ContextData struct {
	Email  string
	UserID string
	Roles  []string
}

// Key - type of the key used in the request context
type Key string

// KeyEmailToken - used for the request context key
const KeyEmailToken Key = "emailtoken"

// ContextStruct - stored in the request context
// set in AuthMiddleware
type ContextStruct struct {
	Email       string
	TokenString string
}

// User - details of the user from the database
type User struct {
	ID    uint
	IDS   string
	Email string
	Role  string
}

// GetAuthBearerToken - extract the BEARER token from the auth header
func GetAuthBearerToken(r *http.Request) (string, error) {

	var APIkey string
	bearer := r.Header.Get("Authorization")
	if len(bearer) > 7 && strings.ToUpper(bearer[0:6]) == "BEARER" {
		APIkey = bearer[7:]
	} else {
		log.WithFields(log.Fields{
			"msgnum": 252,
		}).Error("APIkey Not Found")
		return "", errors.New("APIkey Not Found ")
	}
	return APIkey, nil
}

// GetAuthUserDetails - used for fetching redis details
// In AuthMiddleware, we are setting the Email and Auth Token in the
// request context
// These are extracted from the  request context here, into ContextStruct
// Then we check if this auth token has been stored in Redis
// (the Redis key is the auth token)
// If the auth token has not been stored in Redis, we run a query
// to get the details of the user from the db, and store then in Redis
// for future requests to use
func GetAuthUserDetails(r *http.Request, redisClient *redis.Client, db *sql.DB) (*ContextData, string, error) {
	data := r.Context().Value(KeyEmailToken).(ContextStruct)
	resp, err := redisClient.Get(data.TokenString).Result()
	v := ContextData{}
	if resp == "" {
		user := User{}
		row := db.QueryRow(`select id, id_s, email, role from users where email = ?;`, data.Email)
		err = row.Scan(&user.ID, &user.IDS, &user.Email, &user.Role)
		if user.ID == 0 {
			log.WithFields(log.Fields{
				"msgnum": 253,
			}).Error("User not found")
			return nil, "", errors.New("User not found")
		}
		if err != nil {
			log.WithFields(log.Fields{
				"msgnum": 254,
			}).Error(err)
			return nil, "", errors.New("User not found")
		}
		v.Email = user.Email
		v.UserID = user.IDS
		roles := []string{}
		if user.Role != "" {
			roles = append(roles, user.Role)
		}
		v.Roles = roles
		usr, err := json.Marshal(v)
		if err != nil {
			log.WithFields(log.Fields{
				"msgnum": 255,
			}).Error(err)
			return nil, "", errors.New("User not found")
		}
		err = redisClient.Set(data.TokenString, usr, 0).Err()
		if err != nil {
			log.WithFields(log.Fields{
				"msgnum": 256,
			}).Error(err)
			return nil, "", errors.New("User not found")
		}
	} else {
		err = json.Unmarshal([]byte(resp), &v)
		if err != nil {
			log.WithFields(log.Fields{
				"msgnum": 257,
			}).Error(err)
		}
	}
	return &v, GetRequestID(), nil
}

// CheckRoles - used for checking roles
func CheckRoles(AllowedRoles []string, UserRoles []string) error {
	for _, permission := range AllowedRoles {
		if err := checkRoles(UserRoles, permission); err != nil {
			log.WithFields(log.Fields{
				"msgnum": 263,
			}).Error(err)
			return err
		}
		break
	}
	return nil
}

func checkRoles(roles []string, role string) error {
	if roles == nil {
		return errors.New("No user supplied")
	}

	if role == "" {
		return errors.New("You must supply a valid permission to check against")
	}

	for _, roleName := range roles {
		if role == roleName {
			return nil
		}
	}

	return errors.New("User not authorized")
}
