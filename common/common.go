package common

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"html/template"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rs/xid"
	log "github.com/sirupsen/logrus"
	gomail "gopkg.in/gomail.v2"
)

/* error message range: 250-499 */

// StatusDates - Used in all the database tables
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

// Active - value of status
const Active = 1

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

// Email - for sending email notifications
type Email struct {
	From    string
	To      string
	Subject string
	Body    string
	Cc      string
}

// ParseURL - parses a url into a slice (GetPathParts) and
// the query string (GetPathQueryString)
func ParseURL(urlString string) ([]string, url.Values, error) {
	pathString, queryString, err := GetPathQueryString(urlString)
	if err != nil {
		log.WithFields(log.Fields{"msgnum": 250}).Error(err)
		return []string{}, nil, err
	}

	return GetPathParts(pathString), queryString, nil
}

// GetPathQueryString -- given url string, returns the path, and the
// query string
// Eg. "/v1/users?limit=5&cursor=s4R0Z6ecFTzTC4j=" will return
// "/v1/users", ["limit"]="5", ["cursor"]="s4R0Z6ecFTzTC4j="
func GetPathQueryString(s string) (string, url.Values, error) {

	u, err := url.Parse(s)

	if err != nil {
		log.WithFields(log.Fields{"msgnum": 251}).Error(err)
		return "", nil, err
	}

	return u.Path, u.Query(), nil
}

// GetPathParts - given a url, returns a slice
// of the parts of the url
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

// RenderJSON - send JSON response
func RenderJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.WithFields(log.Fields{"msgnum": 252}).Error(err)
		http.Error(w, err.Error(), 400)
		return
	}
	return
}

// RenderErrorJSON - send error JSON response
func RenderErrorJSON(w http.ResponseWriter, errorCode string, errorMsg string, httpStatusCode int, requestID string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	e := Error{ErrorCode: errorCode, ErrorMsg: errorMsg, HTTPStatusCode: httpStatusCode, RequestID: requestID}
	err := json.NewEncoder(w).Encode(e)
	if err != nil {
		log.WithFields(log.Fields{"msgnum": 253}).Error(err)
		http.Error(w, err.Error(), 400)
		return
	}
	return
}

// GetRequestID - used for RequestID generation
func GetRequestID() string {
	return xid.New().String()
}

// GetUUID - used for UUID generation
func GetUUID() uuid.UUID {
	return uuid.New()
}

// GetUUIDBytes - used for UUID generation, to save in the db
func GetUUIDBytes() ([]byte, error) {
	return uuid.New().MarshalBinary()
}

// UUIDBytesToStr - convert a UUID retrieved from the DB as str,
// to string for sending to the client
func UUIDBytesToStr(b []byte) (string, error) {
	u, err := uuid.FromBytes(b)
	if err != nil {
		log.WithFields(log.Fields{"msgnum": 254}).Error(err)
		return "", err
	}
	return u.String(), nil
}

// UUIDStrToUUID - convert a UUID str into UUID
func UUIDStrToUUID(s string) (uuid.UUID, error) {
	u, err := uuid.Parse(s)
	if err != nil {
		log.WithFields(log.Fields{"msgnum": 255}).Error(err)
		return u, err
	}
	return u, nil
}

// UUIDStrToBytes - convert a UUID str into bytes
func UUIDStrToBytes(s string) ([]byte, error) {
	u, err := uuid.Parse(s)
	if err != nil {
		log.WithFields(log.Fields{"msgnum": 256}).Error(err)
		return nil, err
	}
	return u.MarshalBinary()
}

// ParseTemplate - used for parsing template (for emails)
func ParseTemplate(templateFileName string, data interface{}) (string, error) {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		log.WithFields(log.Fields{"msgnum": 257}).Error(err)
		return "", err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		log.WithFields(log.Fields{"msgnum": 258}).Error(err)
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
		log.WithFields(log.Fields{"msgnum": 259}).Error(err)
		return err
	}
	return nil
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
		log.WithFields(log.Fields{"msgnum": 260}).Error(err)
		return ""
	}
	return string(cursorBytes)
}

// GetTimeDetails - used to populate created_by and updated_by fields
// when inserting/updating records in the database
func GetTimeDetails() (time.Time, uint, uint, uint, uint) {
	tn := time.Now().UTC().Truncate(time.Second)
	tnday := uint(tn.YearDay())
	_, tnweek := tn.ISOWeek()
	tnmonth := uint(tn.Month())
	tnyear := uint(tn.Year())
	return tn, tnday, uint(tnweek), tnmonth, tnyear
}
