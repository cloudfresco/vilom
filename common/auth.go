package common

import (
	"crypto/rand"
	"crypto/sha512"
	"crypto/subtle"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
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
	UUID4 []byte
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
	if err != nil {
		log.WithFields(log.Fields{
			"msgnum": 268,
		}).Error(err)
	}
	v := ContextData{}
	if resp == "" {
		user := User{}
		row := db.QueryRow(`select id, uuid4, email, role from users where email = ?;`, data.Email)
		err = row.Scan(&user.ID, &user.UUID4, &user.Email, &user.Role)
		if user.ID == 0 {
			log.WithFields(log.Fields{
				"msgnum": 261,
			}).Error("User not found")
			return nil, "", errors.New("User not found")
		}
		if err != nil {
			log.WithFields(log.Fields{
				"msgnum": 262,
			}).Error(err)
			return nil, "", errors.New("User not found")
		}
		IDS, err := UUIDBytesToStr(user.UUID4)
		if err != nil {
			log.WithFields(log.Fields{
				"msgnum": 263,
			}).Error(err)
			return nil, "", errors.New("User not found")
		}
		v.Email = user.Email
		v.UserID = IDS
		roles := []string{}
		if user.Role != "" {
			roles = append(roles, user.Role)
		}
		v.Roles = roles
		usr, err := json.Marshal(v)
		if err != nil {
			log.WithFields(log.Fields{
				"msgnum": 264,
			}).Error(err)
			return nil, "", errors.New("User not found")
		}
		err = redisClient.Set(data.TokenString, usr, 0).Err()
		if err != nil {
			log.WithFields(log.Fields{
				"msgnum": 265,
			}).Error(err)
			return nil, "", errors.New("User not found")
		}
	} else {
		err = json.Unmarshal([]byte(resp), &v)
		if err != nil {
			log.WithFields(log.Fields{
				"msgnum": 266,
			}).Error(err)
		}
	}
	return &v, GetRequestID(), nil
}

// HashPassword - Generate hash password
func HashPassword(password string, requestID string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.WithFields(log.Fields{
			"reqid":  requestID,
			"msgnum": 269,
		}).Error(err)
		return []byte{}, err
	}
	return hash, nil
}

// GenTokenHash - GenTokenHash generates pieces needed for passwd recovery
// hash of the first half of a 64 byte value
// (to be stored in the database and used in SELECT query)
// verifier: hash of the second half of a 64 byte value
// (to be stored in database but never used in SELECT query)
// token: the user-facing base64 encoded selector+verifier
func GenTokenHash(requestID string) (selector, verifier, token string, err error) {
	rawToken := make([]byte, 64)
	if _, err = io.ReadFull(rand.Reader, rawToken); err != nil {
		log.WithFields(log.Fields{
			"reqid":  requestID,
			"msgnum": 270,
		}).Error(err)
		return "", "", "", err
	}
	selectorBytes := sha512.Sum512(rawToken[:32])
	verifierBytes := sha512.Sum512(rawToken[32:])

	return base64.StdEncoding.EncodeToString(selectorBytes[:]),
		base64.StdEncoding.EncodeToString(verifierBytes[:]),
		base64.URLEncoding.EncodeToString(rawToken),
		nil
}

// GetSelectorForPasswdRecoveryToken - Get Selector For Password Recovery Token
func GetSelectorForPasswdRecoveryToken(token string, requestID string) ([64]byte, string, error) {

	rawToken, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		log.WithFields(log.Fields{
			"reqid":  requestID,
			"msgnum": 271,
		}).Error(err)
		return [64]byte{}, "", errors.New("invalid recover token submitted, base64 decode failed")
	}

	if len(rawToken) != 64 {
		log.WithFields(log.Fields{
			"reqid":  requestID,
			"msgnum": 272,
		}).Error("invalid recover token submitted, size was wrong")
		return [64]byte{}, "", errors.New("invalid recover token submitted, size was wrong")
	}

	selectorBytes := sha512.Sum512(rawToken[:32])
	verifierBytes := sha512.Sum512(rawToken[32:])
	selector := base64.StdEncoding.EncodeToString(selectorBytes[:])

	return verifierBytes, selector, nil
}

// ValidatePasswdRecoveryToken - Validate Passwd Recovery Token
func ValidatePasswdRecoveryToken(verifierBytes [64]byte, verifier string, tokenExpiry time.Time, requestID string) error {
	tn, _, _, _, _ := GetTimeDetails()
	if tn.After(tokenExpiry) {
		log.WithFields(log.Fields{
			"reqid":  requestID,
			"msgnum": 273,
		}).Error("Token already expired")
		return errors.New("Token already expired")
	}

	dbVerifierBytes, err := base64.StdEncoding.DecodeString(verifier)
	if err != nil {
		log.WithFields(log.Fields{
			"reqid":  requestID,
			"msgnum": 274,
		}).Error(err)
		return err
	}
	if subtle.ConstantTimeEq(int32(len(verifierBytes)), int32(len(dbVerifierBytes))) != 1 ||
		subtle.ConstantTimeCompare(verifierBytes[:], dbVerifierBytes) != 1 {
		log.WithFields(log.Fields{
			"reqid":  requestID,
			"msgnum": 275,
		}).Error("stored recover verifier does not match provided one")
		return errors.New("stored recover verifier does not match provided one")
	}

	return nil

}

// CheckRoles - used for checking roles
func CheckRoles(AllowedRoles []string, UserRoles []string) error {
	for _, permission := range AllowedRoles {
		if err := checkRoles(UserRoles, permission); err != nil {
			log.WithFields(log.Fields{
				"msgnum": 267,
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
