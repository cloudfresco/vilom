package common

import (
	"bytes"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/go-redis/redis"
	"github.com/rs/xid"
	gomail "gopkg.in/gomail.v2"
)

// StatusDates - Used for all structs
type StatusDates struct {
	Statusc      uint
	CreatedAt    time.Time
	UpdatedAt    time.Time
	CreatedDay   uint
	CreatedWeek  uint
	CreatedMonth uint
	CreatedYear  uint
	UpdatedDay   uint
	UpdatedWeek  uint
	UpdatedMonth uint
	UpdatedYear  uint
}

// Active - used for status of all struct
const Active = 1

// Key - Its key type
type Key string

// KeyEmailToken - used for context key
const KeyEmailToken Key = "emailtoken"

// Error - used for
type Error struct {
	ErrorCode      string `json:"error_code"`
	ErrorMsg       string `json:"error_msg"`
	HTTPStatusCode int    `json:"status"`
	RequestID      string `json:"request_id"`
}

// RedisOptions - used for minimal RedisOptions view representation
type RedisOptions struct {
	Addr string `mapstructure:"addr"`
}

// KeyOptions - used for
type KeyOptions struct {
	CaCertPath string `mapstructure:"CaCerTPath"`
	CertPath   string `mapstructure:"CertPath"`
	KeyPath    string `mapstructure:"KeyPath"`
	ServerAddr string `mapstructure:"ServerAddr"`
}

// OauthOptions - used for
type OauthOptions struct {
	ClientID     string `mapstructure:"ClientID"`
	ClientSecret string `mapstructure:"ClientSecret"`
}

// DbOptions - used for
type DbOptions struct {
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Host     string `mapstructure:"hostname"`
	Port     string `mapstructure:"port"`
	Schema   string `mapstructure:"database"`
}

// MailerOptions - used for
type MailerOptions struct {
	User     string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Port     int    `mapstructure:"port"`
	Server   string `mapstructure:"server"`
}

// JWTOptions - used for
type JWTOptions struct {
	JWTKey      []byte
	JWTDuration int
}

// RateLimiterOptions - used for
type RateLimiterOptions struct {
	UserMaxRate    int `mapstructure:"usermaxrate"`
	UserMaxBurst   int `mapstructure:"usermaxburst"`
	UgroupMaxRate  int `mapstructure:"ugroupmaxrate"`
	UgroupMaxBurst int `mapstructure:"ugroupmaxburst"`
	CatMaxRate     int `mapstructure:"catmaxrate"`
	CatMaxBurst    int `mapstructure:"catmaxburst"`
	TopicMaxRate   int `mapstructure:"topicmaxrate"`
	TopicMaxBurst  int `mapstructure:"topicmaxburst"`
	MsgMaxRate     int `mapstructure:"msgmaxrate"`
	MsgMaxBurst    int `mapstructure:"msgmaxburst"`
	UbadgeMaxRate  int `mapstructure:"ubadgemaxrate"`
	UbadgeMaxBurst int `mapstructure:"ubadgemaxburst"`
	SearchMaxRate  int `mapstructure:"searchmaxrate"`
	SearchMaxBurst int `mapstructure:"searchmaxburst"`
	UMaxRate       int `mapstructure:"umaxrate"`
	UMaxBurst      int `mapstructure:"umaxburst"`
}

// UserOptions - used for
type UserOptions struct {
	ConfirmTokenDuration string `mapstructure:"confirmtokenduration"`
	ResetTokenDuration   string `mapstructure:"resettokenduration"`
}

// Email - used for email
type Email struct {
	From    string
	To      string
	Subject string
	Body    string
	Cc      string
}

// User - used for
type User struct {
	ID    uint
	IDS   string
	Email string
	Role  string
}

// ParseURL - used for
func ParseURL(urlString string) ([]string, url.Values, error) {
	pathString, queryString, err := GetPathQueryString(urlString)
	if err != nil {
		log.WithFields(log.Fields{
			"msgnum": 250,
		}).Error(err)
		return []string{}, nil, err
	}
	pathParts := GetPathParts(pathString)

	return pathParts, queryString, nil
}

// GetPathQueryString -- given url string, returns the path, and the
// query string
// Eg. "/v1/users?limit=5&cursor=s4R0Z6ecFTzTC4j=" will return
// "/v1/users", ["limit"]="5", ["cursor"]="s4R0Z6ecFTzTC4j="
func GetPathQueryString(s string) (string, url.Values, error) {

	u, err := url.Parse(s)

	if err != nil {
		log.WithFields(log.Fields{
			"msgnum": 251,
		}).Error(err)
		return "", nil, err
	}

	p := u.Path
	q := u.Query()

	return p, q, nil
}

// GetPathParts - used for spliting routes
func GetPathParts(url string) []string {

	var pathParts []string

	sliceOfSubstrings := strings.Split(url, "/")

	for _, p := range sliceOfSubstrings {
		if p != "" {
			pathParts = append(pathParts, p)
		}
	}

	return pathParts
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

// ContextData - used for
type ContextData struct {
	Email  string
	UserID string
	Roles  []string
}

// ContextStruct - used for
type ContextStruct struct {
	Email       string
	TokenString string
}

// GetAuthUserDetails - used for fetching redis details
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
	requestID := GetRequestID()
	return &v, requestID, nil
}

// RenderJSON - send JSON response
func RenderJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.WithFields(log.Fields{
			"msgnum": 258,
		}).Error(err)
		http.Error(w, err.Error(), 400)
		return
	}
	return
}

// RenderErrorJSON - send JSON response
func RenderErrorJSON(w http.ResponseWriter, errorCode string, errorMsg string, httpStatusCode int, requestID string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	e := Error{ErrorCode: errorCode, ErrorMsg: errorMsg, HTTPStatusCode: httpStatusCode, RequestID: requestID}
	err := json.NewEncoder(w).Encode(e)
	if err != nil {
		log.WithFields(log.Fields{
			"msgnum": 259,
		}).Error(err)
		http.Error(w, err.Error(), 400)
		return
	}
	return
}

// GetRequestID - used for GetRequestID generation
func GetRequestID() string {
	return fmt.Sprintf("%x", xid.New().String())
}

// GetUID - used for id generation
func GetUID() string {
	return fmt.Sprintf("%x", xid.New().String())
}

// ParseTemplate - used for parsing template
func ParseTemplate(templateFileName string, data interface{}) (string, error) {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		log.WithFields(log.Fields{
			"msgnum": 260,
		}).Error(err)
		return "", err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		log.WithFields(log.Fields{
			"msgnum": 261,
		}).Error(err)
		return "", err
	}
	body := buf.String()
	return body, nil
}

// SendMail - used for sending email
func SendMail(msg Email, gomailer *gomail.Dialer) error {
	m := gomail.NewMessage()
	m.SetHeader("From", gomailer.Username)
	m.SetHeader("To", msg.To)
	m.SetHeader("Subject", msg.Subject)
	m.SetBody("text/html", msg.Body)

	err := gomailer.DialAndSend(m)
	if err != nil {
		log.WithFields(log.Fields{
			"msgnum": 262,
		}).Error(err)
		return err
	}
	return nil
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

// EncodeCursor - encode cursor
func EncodeCursor(cursor uint) string {
	cursorStr := strconv.FormatUint(uint64(cursor), 10)
	return base64.StdEncoding.EncodeToString([]byte(cursorStr))
}

// DecodeCursor - decode cursor
func DecodeCursor(cursor string) string {
	cursorBytes, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		log.WithFields(log.Fields{
			"msgnum": 264,
		}).Error(err)
		return ""
	}
	return string(cursorBytes)
}
