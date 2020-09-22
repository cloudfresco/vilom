package common

import (
	"context"
	"crypto/rand"
	"crypto/sha512"
	"crypto/subtle"
	"database/sql"
	"encoding/base64"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	casbindb "github.com/Blank-Xu/sql-adapter"
	"github.com/casbin/casbin/v2"
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

var jwtOpt *JWTOptions

// SetJWTOpt set JWT opt used in auth middleware
func SetJWTOpt(jwt *JWTOptions) {
	jwtOpt = jwt
}

// GetJWTOpt get JWT opt used in auth middleware
func GetJWTOpt() *JWTOptions {
	return jwtOpt
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

// LoadEnforcer - used for checking roles
func LoadEnforcer(dbOpt *DBService, roleOpt *RoleOptions) (*casbin.Enforcer, error) {
	e := &casbin.Enforcer{}
	if dbOpt.DBType == DBMysql {
		// Initialize an adapter and use it in a Casbin enforcer:
		casbinMysql, err := casbindb.NewAdapter(dbOpt.DB, DBMysql, roleOpt.RolesTableName)
		if err != nil {
			log.WithFields(log.Fields{
				"msgnum": 268,
			}).Error(err)
			return nil, err
		}

		ctx := context.Background()
		tx, err := dbOpt.DB.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
		if err != nil {
			log.WithFields(log.Fields{
				"msgnum": 269,
			}).Error(err)
			return nil, err
		}

		_, err = tx.ExecContext(ctx, `truncate `+roleOpt.RolesTableName+``)
		if err != nil {
			log.WithFields(log.Fields{
				"msgnum": 270,
			}).Error(err)
			return nil, err
		}
		for _, r := range roleOpt.Roles {
			stmt, err := tx.Prepare(`insert into ` + roleOpt.RolesTableName + `
					(p_type,
					v0,
					v1,
					v2,
					v3,
					v4,
					v5)
				values (?,?,?,?,?,?,?);`)

			if err != nil {
				log.WithFields(log.Fields{
					"msgnum": 271,
				}).Error(err)
				return nil, err
			}
			_, err = stmt.Exec(
				r.PType,
				r.V0,
				r.V1,
				r.V2,
				r.V3,
				r.V4,
				r.V5)
			if err != nil {
				log.WithFields(log.Fields{
					"msgnum": 272,
				}).Error(err)
				return nil, err
			}
		}

		err = tx.Commit()
		if err != nil {
			log.WithFields(log.Fields{
				"msgnum": 273,
			}).Error(err)
			return nil, err
		}
		// Load policy from file

		e, err = casbin.NewEnforcer(roleOpt.RolesPolicyConfigPath, casbinMysql)
		if err != nil {
			log.WithFields(log.Fields{
				"msgnum": 274,
			}).Error(err)
			return nil, err
		}

		// Load the policy from DB.
		if err = e.LoadPolicy(); err != nil {
			log.WithFields(log.Fields{
				"msgnum": 275,
			}).Error(err)
			return nil, err
		}
	} else if dbOpt.DBType == DBPgsql {

	}
	return e, nil
}
