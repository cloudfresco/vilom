package userservices

import (
	"context"
	"crypto/rand"
	"crypto/sha512"
	"crypto/subtle"
	"database/sql"
	"encoding/base64"
	"errors"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	gomail "gopkg.in/gomail.v2"

	"github.com/cloudfresco/vilom/common"
)

/* error message range: 1500-1999 */

// User - User view representation
type User struct {
	ID        uint
	UUID4     []byte
	IDS       string
	AuthToken string

	Email     string
	Username  string
	FirstName string `sql:"not null"`
	LastName  string
	Role      string
	Password  []byte
	Active    bool `sql:"default:false"`

	EmailConfirmationToken string
	EmailSelector          string
	EmailVerifier          string
	EmailTokenSentAt       time.Time
	EmailTokenExpiry       time.Time
	EmailConfirmedAt       time.Time

	NewEmail            string
	NewEmailResetToken  string
	NewEmailSelector    string
	NewEmailVerifier    string
	NewEmailTokenSentAt time.Time
	NewEmailTokenExpiry time.Time
	NewEmailConfirmedAt time.Time

	PasswordResetToken  string
	PasswordSelector    string
	PasswordVerifier    string
	PasswordTokenSentAt time.Time
	PasswordTokenExpiry time.Time
	PasswordConfirmedAt time.Time

	Timezone        string `sql:"default:'Asia/Kolkata'"`
	SignInCount     uint
	CurrentSignInAt time.Time
	LastSignInAt    time.Time

	common.StatusDates

	/* used only for logic purpose */
	Roles       []string
	PasswordS   string
	HostURL     string
	Tokenstring string
}

// LoginForm - user login form
type LoginForm struct {
	Email    string `json:"email" valid:"email,required"`
	Password string `json:"password"`
}

// UserEmailForm - user email form
type UserEmailForm struct {
	Email string `json:"email" valid:"email,required"`
}

// PasswordForm - change password form
type PasswordForm struct {
	Password        string `form:"passwd" binding:"required"`
	ConfirmPassword string `form:"confrmpasswd" binding:"required"`
	CurrentPassword string `form:"cpasswd"`
	ID              string
}

// ChangeEmailForm - used for Change Email
type ChangeEmailForm struct {
	Email    string
	NewEmail string
}

// ForgotPasswordForm - used for forgot password
type ForgotPasswordForm struct {
	Email string `form:"email" binding:"required"`
}

// UserService - For accessing user services
type UserService struct {
	Config       *common.RedisOptions
	Db           *sql.DB
	RedisClient  *redis.Client
	Mailer       *gomail.Dialer
	JWTOptions   *common.JWTOptions
	LimitDefault string
	UserOptions  *common.UserOptions
}

// NewUserService - Create User Service
func NewUserService(config *common.RedisOptions,
	db *sql.DB,
	redisClient *redis.Client,
	mailer *gomail.Dialer,
	jwtOptions *common.JWTOptions,
	limitDefault string,
	userOptions *common.UserOptions) *UserService {
	return &UserService{config, db, redisClient, mailer, jwtOptions, limitDefault, userOptions}
}

// Roles - Used for roles
type Roles []string

// CustomClaims - used to type holds the token claims
type CustomClaims struct {
	EmailAddr string
	jwt.StandardClaims
}

// HashPassword - Generate hash password
func HashPassword(password string, requestID string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.WithFields(log.Fields{
			"reqid":  requestID,
			"msgnum": 1500,
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
			"msgnum": 1501,
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
			"msgnum": 1502,
		}).Error(err)
		return [64]byte{}, "", errors.New("invalid recover token submitted, base64 decode failed")
	}

	if len(rawToken) != 64 {
		log.WithFields(log.Fields{
			"reqid":  requestID,
			"msgnum": 1503,
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
	tn, _, _, _, _ := common.GetTimeDetails()
	if tn.After(tokenExpiry) {
		log.WithFields(log.Fields{
			"reqid":  requestID,
			"msgnum": 1504,
		}).Error("Token already expired")
		return errors.New("Token already expired")
	}

	dbVerifierBytes, err := base64.StdEncoding.DecodeString(verifier)
	if err != nil {
		log.WithFields(log.Fields{
			"reqid":  requestID,
			"msgnum": 1505,
		}).Error(err)
		return err
	}
	if subtle.ConstantTimeEq(int32(len(verifierBytes)), int32(len(dbVerifierBytes))) != 1 ||
		subtle.ConstantTimeCompare(verifierBytes[:], dbVerifierBytes) != 1 {
		log.WithFields(log.Fields{
			"reqid":  requestID,
			"msgnum": 1506,
		}).Error("stored recover verifier does not match provided one")
		return errors.New("stored recover verifier does not match provided one")
	}

	log.Info("validated")
	return nil

}

// UserCursor - used for getting users list
type UserCursor struct {
	Users      []*User
	NextCursor string
}

// GetUsers - Get all users
func (u *UserService) GetUsers(ctx context.Context, limit string, nextCursor string, userEmail string, requestID string) (*UserCursor, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{
			"user":   userEmail,
			"reqid":  requestID,
			"msgnum": 1507,
		}).Error(err)
		return nil, err
	default:
		if limit == "" {
			limit = u.LimitDefault
		}
		query := ""
		if nextCursor == "" {
			query = query + " order by id desc " + " limit " + limit + ";"
		} else {
			nextCursor = common.DecodeCursor(nextCursor)
			query = query + "where " + "id <= " + nextCursor + " order by id desc " + " limit " + limit + ";"
		}
		users := []*User{}
		rows, err := u.Db.QueryContext(ctx, `select id, uuid4, auth_token, first_name, last_name, email, role from users `+query)
		if err != nil {
			log.WithFields(log.Fields{
				"user":   userEmail,
				"reqid":  requestID,
				"msgnum": 1508,
			}).Error(err)
			return nil, err
		}

		for rows.Next() {
			user := User{}
			err = rows.Scan(&user.ID, &user.UUID4, &user.AuthToken, &user.FirstName, &user.LastName, &user.Email, &user.Role)
			if err != nil {
				log.WithFields(log.Fields{
					"user":   userEmail,
					"reqid":  requestID,
					"msgnum": 1509,
				}).Error(err)
				err = rows.Close()
				return nil, err
			}
			uUID4Str, err := common.UUIDBytesToStr(user.UUID4)
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 1510}).Error(err)
				return nil, err
			}
			user.IDS = uUID4Str
			users = append(users, &user)
		}
		err = rows.Close()
		if err != nil {
			log.WithFields(log.Fields{
				"user":   userEmail,
				"reqid":  requestID,
				"msgnum": 1511,
			}).Error(err)
			return nil, err
		}

		err = rows.Err()
		if err != nil {
			log.WithFields(log.Fields{
				"user":   userEmail,
				"reqid":  requestID,
				"msgnum": 1512,
			}).Error(err)
			return nil, err
		}
		x := UserCursor{}
		if len(users) != 0 {
			next := users[len(users)-1].ID
			next = next - 1
			nextc := common.EncodeCursor(next)
			x = UserCursor{users, nextc}
		} else {
			x = UserCursor{users, "0"}
		}
		return &x, nil
	}

}

// Login - used for Login user
func (u *UserService) Login(ctx context.Context, form *LoginForm, requestID string) (*User, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{
			"reqid":  requestID,
			"msgnum": 1513,
		}).Error(err)
		return nil, err
	default:
		db := u.Db
		user := User{}
		row := db.QueryRowContext(ctx, `select id, email, password from users where email = ?;`, form.Email)
		err := row.Scan(
			&user.ID,
			&user.Email,
			&user.Password)

		if err != nil {
			log.WithFields(log.Fields{
				"reqid":  requestID,
				"msgnum": 1514,
			}).Error(err)
			return nil, err
		}

		err = bcrypt.CompareHashAndPassword(user.Password, []byte(form.Password))
		if err != nil {
			log.WithFields(log.Fields{
				"reqid":  requestID,
				"msgnum": 1515,
			}).Error(err)
			return nil, err
		}
		tokenDuration := time.Duration(u.JWTOptions.JWTDuration)
		tokenStr, err := u.CreateJWT(form.Email, tokenDuration, requestID)
		if err != nil {
			log.WithFields(log.Fields{
				"reqid":  requestID,
				"msgnum": 1516,
			}).Error(err)
			return nil, err
		}
		user.Tokenstring = tokenStr
		if err != nil {
			log.WithFields(log.Fields{
				"reqid":  requestID,
				"msgnum": 1517,
			}).Error(err)
			return nil, err
		}
		return &user, err
	}
}

// CreateJWT - Create jwt token
func (u *UserService) CreateJWT(emailAddr string, tokenDuration time.Duration, requestID string) (string, error) {
	tn, _, _, _, _ := common.GetTimeDetails()
	claims := CustomClaims{
		EmailAddr: emailAddr,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: tn.Add(time.Hour * tokenDuration).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token with key
	tokenString, err := token.SignedString(u.JWTOptions.JWTKey)
	if err != nil {
		log.WithFields(log.Fields{
			"reqid":  requestID,
			"msgnum": 1518,
		}).Error("Failed to sign token")
		return "", errors.New("Failed to sign token")
	}

	return tokenString, nil

}

// Create - Create User
func (u *UserService) Create(ctx context.Context, form *User, hostURL string, requestID string) (*User, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{
			"reqid":  requestID,
			"msgnum": 1519,
		}).Error(err)
		return nil, err
	default:
		db := u.Db
		//check if email already exists
		var isPresent bool
		row := db.QueryRowContext(ctx, `select exists (select 1 from users where email = ?)`, form.Email)
		err := row.Scan(&isPresent)
		if err != nil {
			log.WithFields(log.Fields{
				"reqid":  requestID,
				"msgnum": 1520,
			}).Error(err)

			return nil, err
		}
		if isPresent {
			err = errors.New("Email Already Exists")
			log.WithFields(log.Fields{
				"reqid":  requestID,
				"msgnum": 1521,
			}).Error(err)

			return nil, err
		}

		password1, err := HashPassword(form.PasswordS, requestID)
		if err != nil {
			log.WithFields(log.Fields{
				"reqid":  requestID,
				"msgnum": 1522,
			}).Error(err)

			return nil, err
		}

		selector, verifier, token, err := GenTokenHash(requestID)
		if err != nil {
			log.WithFields(log.Fields{
				"reqid":  requestID,
				"msgnum": 1523,
			}).Error(err)

			return nil, err
		}

		tn, tnday, tnweek, tnmonth, tnyear := common.GetTimeDetails()
		tokenExpiry, _ := time.ParseDuration(u.UserOptions.ConfirmTokenDuration)

		user := User{}
		user.UUID4, err = common.GetUUIDBytes()
		if err != nil {
			log.WithFields(log.Fields{
				"reqid":  requestID,
				"msgnum": 1524,
			}).Error(err)

			return nil, err
		}
		user.AuthToken = ""
		user.Email = form.Email
		user.Username = form.Email
		user.FirstName = form.FirstName
		user.LastName = form.LastName
		user.Password = password1
		user.Role = form.Role
		user.Active = false
		user.EmailConfirmationToken = token
		user.EmailSelector = selector
		user.EmailVerifier = verifier
		user.EmailTokenSentAt = tn
		user.EmailTokenExpiry = tn.Add(tokenExpiry)
		user.EmailConfirmedAt = tn
		user.NewEmail = ""
		user.NewEmailResetToken = ""
		user.NewEmailSelector = ""
		user.NewEmailVerifier = ""
		user.NewEmailTokenSentAt = tn
		user.NewEmailTokenExpiry = tn
		user.NewEmailConfirmedAt = tn
		user.PasswordResetToken = ""
		user.PasswordSelector = ""
		user.PasswordVerifier = ""
		user.PasswordTokenSentAt = tn
		user.PasswordTokenExpiry = tn
		user.PasswordConfirmedAt = tn
		user.Timezone = "Asia/Kolkata"
		user.SignInCount = 0
		user.CurrentSignInAt = tn
		user.LastSignInAt = tn
		user.Statusc = common.Active
		user.CreatedAt = tn
		user.UpdatedAt = tn
		user.CreatedDay = tnday
		user.CreatedWeek = tnweek
		user.CreatedMonth = tnmonth
		user.CreatedYear = tnyear
		user.UpdatedDay = tnday
		user.UpdatedWeek = tnweek
		user.UpdatedMonth = tnmonth
		user.UpdatedYear = tnyear

		tx, err := db.Begin()
		if err != nil {
			log.WithFields(log.Fields{
				"reqid":  requestID,
				"msgnum": 1525,
			}).Error(err)

			return nil, err
		}
		usr, err := u.InsertUser(ctx, tx, user, hostURL, requestID)
		if err != nil {
			log.WithFields(log.Fields{
				"reqid":  requestID,
				"msgnum": 1526,
			}).Error(err)

			err = tx.Rollback()
			return nil, err
		}

		err = tx.Commit()
		if err != nil {
			log.WithFields(log.Fields{
				"reqid":  requestID,
				"msgnum": 1527,
			}).Error(err)

			return nil, err
		}
		return usr, nil
	}
}

// InsertUser - Insert User details to database
func (u *UserService) InsertUser(ctx context.Context, tx *sql.Tx, user User, hostURL string, requestID string) (*User, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{
			"reqid":  requestID,
			"msgnum": 1528,
		}).Error(err)

		return nil, err
	default:
		stmt, err := tx.PrepareContext(ctx, `insert into users
	  (
		uuid4,
    auth_token,
		email,
    username,
		first_name,
		last_name,
		role,
		password,
		active,
		email_confirmation_token,
    email_selector,
    email_verifier,
		email_token_sent_at,
		email_token_expiry,
		email_confirmed_at,
		new_email,
		new_email_reset_token,
    new_email_selector,
    new_email_verifier,
		new_email_token_sent_at,
		new_email_token_expiry,
		new_email_confirmed_at,
		password_reset_token,
    password_selector,
    password_verifier,
		password_token_sent_at,
		password_token_expiry,
		password_confirmed_at,
		timezone,
		sign_in_count,
		current_sign_in_at,
		last_sign_in_at,
		statusc,
		created_at,
		updated_at,
		created_day,
		created_week,
		created_month,
		created_year,
		updated_day,
		updated_week,
		updated_month,
		updated_year)
  values (?,?,?,?,?,?,?,?,?,?,
					?,?,?,?,?,?,?,?,?,?,
					?,?,?,?,?,?,?,?,?,?,
					?,?,?,?,?,?,?,?,?,?,
          ?,?,?);`)
		if err != nil {
			log.WithFields(log.Fields{
				"reqid":  requestID,
				"msgnum": 1529,
			}).Error(err)

			return nil, err
		}
		res, err := stmt.ExecContext(ctx,
			user.UUID4,
			user.AuthToken,
			user.Email,
			user.Username,
			user.FirstName,
			user.LastName,
			user.Role,
			user.Password,
			user.Active,
			user.EmailConfirmationToken,
			user.EmailSelector,
			user.EmailVerifier,
			user.EmailTokenSentAt,
			user.EmailTokenExpiry,
			user.EmailConfirmedAt,
			user.NewEmail,
			user.NewEmailResetToken,
			user.NewEmailSelector,
			user.NewEmailVerifier,
			user.NewEmailTokenSentAt,
			user.NewEmailTokenExpiry,
			user.NewEmailConfirmedAt,
			user.PasswordResetToken,
			user.PasswordSelector,
			user.PasswordVerifier,
			user.PasswordTokenSentAt,
			user.PasswordTokenExpiry,
			user.PasswordConfirmedAt,
			user.Timezone,
			user.SignInCount,
			user.CurrentSignInAt,
			user.LastSignInAt,
			user.Statusc,
			user.CreatedAt,
			user.UpdatedAt,
			user.CreatedDay,
			user.CreatedWeek,
			user.CreatedMonth,
			user.CreatedYear,
			user.UpdatedDay,
			user.UpdatedWeek,
			user.UpdatedMonth,
			user.UpdatedYear)
		if err != nil {
			log.WithFields(log.Fields{
				"reqid":  requestID,
				"msgnum": 1530,
			}).Error(err)

			err = stmt.Close()
			return nil, err
		}
		uID, err := res.LastInsertId()
		if err != nil {
			log.WithFields(log.Fields{
				"reqid":  requestID,
				"msgnum": 1531,
			}).Error(err)

			err = stmt.Close()
			return nil, err
		}
		user.ID = uint(uID)
		uUID4Str, err := common.UUIDBytesToStr(user.UUID4)
		if err != nil {
			log.WithFields(log.Fields{"reqid": requestID, "msgnum": 1532}).Error(err)
			return nil, err
		}
		user.IDS = uUID4Str
		err = stmt.Close()
		if err != nil {
			log.WithFields(log.Fields{
				"reqid":  requestID,
				"msgnum": 1533,
			}).Error(err)

			return nil, err
		}

		pwd, _ := os.Getwd()
		viewpath := pwd + filepath.FromSlash("/common/views/confirmation.html")
		templateData := struct {
			Title string
			URL   string
		}{
			Title: "Confirmation",
			URL:   "http://" + hostURL + "/u/confirmation/" + user.EmailConfirmationToken,
		}
		ConfirmationEmail, err := common.ParseTemplate(viewpath, templateData)
		if err != nil {
			log.WithFields(log.Fields{
				"reqid":  requestID,
				"msgnum": 1534,
			}).Error(err)

			return nil, err
		}

		email := common.Email{
			To:      user.Email,
			Subject: "Confirmation",
			Body:    ConfirmationEmail,
		}

		err = common.SendMail(email, u.Mailer)
		if err != nil {
			log.WithFields(log.Fields{
				"reqid":  requestID,
				"msgnum": 1535,
			}).Error(err)

			return nil, err
		}

		return &user, nil
	}
}

// ConfirmEmail - used to confirm email
func (u *UserService) ConfirmEmail(ctx context.Context, token string, requestID string) error {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{
			"reqid":  requestID,
			"msgnum": 1536,
		}).Error(err)

		return err
	default:
		db := u.Db
		verifierBytes, selector, err := GetSelectorForPasswdRecoveryToken(token, requestID)
		if err != nil {
			log.WithFields(log.Fields{
				"reqid":  requestID,
				"msgnum": 1537,
			}).Error(err)

			return err
		}
		user := User{}
		row := db.QueryRowContext(ctx, `select id, email_selector, email_verifier, email_token_expiry from users where email_selector = ?;`, selector)

		err = row.Scan(
			&user.ID,
			&user.EmailSelector,
			&user.EmailVerifier,
			&user.EmailTokenExpiry)

		if err != nil {
			log.WithFields(log.Fields{
				"reqid":  requestID,
				"msgnum": 1538,
			}).Error(err)

			return err
		}

		err = ValidatePasswdRecoveryToken(verifierBytes, user.EmailVerifier, user.EmailTokenExpiry, requestID)
		if err != nil {
			log.WithFields(log.Fields{
				"reqid":  requestID,
				"msgnum": 1539,
			}).Error(err)

			return err
		}

		tn, tnday, tnweek, tnmonth, tnyear := common.GetTimeDetails()

		UpdatedDay := tnday
		UpdatedWeek := tnweek
		UpdatedMonth := tnmonth
		UpdatedYear := tnyear

		stmt, err := db.PrepareContext(ctx, `update users set 
				email_confirmation_token = ?,
				email_selector = ?,
				email_verifier = ?,
		    email_confirmed_at = ?,
		    statusc = ?,
        active = ?,
		    updated_at = ?, 
				updated_day = ?, 
				updated_week = ?, 
				updated_month = ?, 
				updated_year = ? where id= ?;`)
		if err != nil {
			log.WithFields(log.Fields{
				"reqid":  requestID,
				"msgnum": 1540,
			}).Error(err)

			err = stmt.Close()
			return err
		}

		_, err = stmt.ExecContext(ctx,
			"",
			"",
			"",
			tn,
			common.Active,
			true,
			tn,
			UpdatedDay,
			UpdatedWeek,
			UpdatedMonth,
			UpdatedYear,
			user.ID)
		if err != nil {
			log.WithFields(log.Fields{
				"reqid":  requestID,
				"msgnum": 1541,
			}).Error(err)

			err = stmt.Close()
			return err
		}
		err = stmt.Close()

		if err != nil {
			log.WithFields(log.Fields{
				"reqid":  requestID,
				"msgnum": 1542,
			}).Error(err)

			return err
		}

		return nil
	}
}

// ForgotPassword - used to reset forgotten Password
func (u *UserService) ForgotPassword(ctx context.Context, form *ForgotPasswordForm, hostURL string, requestID string) error {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{
			"reqid":  requestID,
			"msgnum": 1543,
		}).Error(err)

		return err
	default:
		db := u.Db
		user, err := u.GetUserByEmail(ctx, form.Email, "", requestID)
		if err != nil {
			log.WithFields(log.Fields{
				"reqid":  requestID,
				"msgnum": 1544,
			}).Error(err)

			return err
		}

		selector, verifier, token, err := GenTokenHash(requestID)
		if err != nil {
			log.WithFields(log.Fields{
				"reqid":  requestID,
				"msgnum": 1545,
			}).Error(err)

			return err
		}
		tn, tnday, tnweek, tnmonth, tnyear := common.GetTimeDetails()

		tokenExpiry, _ := time.ParseDuration(u.UserOptions.ResetTokenDuration)
		resetExpiry := tn.Add(tokenExpiry)

		UpdatedDay := tnday
		UpdatedWeek := tnweek
		UpdatedMonth := tnmonth
		UpdatedYear := tnyear

		stmt, err := db.PrepareContext(ctx, `update users set 
		    password_reset_token = ?,
				password_selector = ?,
				password_verifier = ?,
        password_token_sent_at = ?,
		    password_token_expiry = ?,
		    updated_at = ?, 
				updated_day = ?, 
				updated_week = ?, 
				updated_month = ?, 
				updated_year = ? where id= ?;`)
		if err != nil {
			log.WithFields(log.Fields{
				"reqid":  requestID,
				"msgnum": 1546,
			}).Error(err)

			err = stmt.Close()
			return err
		}

		_, err = stmt.ExecContext(ctx,
			token,
			selector,
			verifier,
			tn,
			resetExpiry,
			tn,
			UpdatedDay,
			UpdatedWeek,
			UpdatedMonth,
			UpdatedYear,
			user.ID)
		if err != nil {
			log.WithFields(log.Fields{
				"reqid":  requestID,
				"msgnum": 1547,
			}).Error(err)

			err = stmt.Close()
			return err
		}
		err = stmt.Close()

		if err != nil {
			log.WithFields(log.Fields{
				"reqid":  requestID,
				"msgnum": 1548,
			}).Error(err)

			return err
		}

		pwd, _ := os.Getwd()
		viewpath := pwd + filepath.FromSlash("/common/views/reset_password.html")

		templateData := struct {
			Title string
			URL   string
		}{
			Title: "Reset Password",
			URL:   "http://" + hostURL + "/u/reset_password/" + token,
		}

		ResetPasswordEmail, err := common.ParseTemplate(viewpath, templateData)
		if err != nil {
			log.WithFields(log.Fields{
				"reqid":  requestID,
				"msgnum": 1549,
			}).Error(err)

			return err
		}

		recipient := user.Email
		email := common.Email{
			To:      recipient,
			Subject: "Reset Passowrd",
			Body:    ResetPasswordEmail,
		}

		err = common.SendMail(email, u.Mailer)
		if err != nil {
			log.WithFields(log.Fields{
				"reqid":  requestID,
				"msgnum": 1550,
			}).Error(err)
			return err
		}

		return nil
	}
}

// ConfirmForgotPassword - used to confirm forgotten password
func (u *UserService) ConfirmForgotPassword(ctx context.Context, form *PasswordForm, token string, requestID string) error {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		return err
	default:
		db := u.Db
		password1, err := HashPassword(form.Password, requestID)
		if err != nil {
			log.WithFields(log.Fields{
				"reqid":  requestID,
				"msgnum": 1551,
			}).Error(err)

		}

		verifierBytes, selector, err := GetSelectorForPasswdRecoveryToken(token, requestID)
		if err != nil {
			log.WithFields(log.Fields{
				"reqid":  requestID,
				"msgnum": 1552,
			}).Error(err)

			return err
		}
		user := User{}
		row := db.QueryRowContext(ctx, `select id, password_selector, password_verifier, password_token_expiry from users where password_selector = ?;`, selector)

		err = row.Scan(
			&user.ID,
			&user.PasswordSelector,
			&user.PasswordVerifier,
			&user.PasswordTokenExpiry)

		if err != nil {
			log.WithFields(log.Fields{
				"reqid":  requestID,
				"msgnum": 1553,
			}).Error(err)

			return err
		}

		err = ValidatePasswdRecoveryToken(verifierBytes, user.PasswordVerifier, user.PasswordTokenExpiry, requestID)
		if err != nil {
			log.WithFields(log.Fields{
				"reqid":  requestID,
				"msgnum": 1554,
			}).Error(err)

			return err
		}

		tn, tnday, tnweek, tnmonth, tnyear := common.GetTimeDetails()

		UpdatedDay := tnday
		UpdatedWeek := tnweek
		UpdatedMonth := tnmonth
		UpdatedYear := tnyear

		tx, err := db.Begin()
		if err != nil {
			log.WithFields(log.Fields{
				"reqid":  requestID,
				"msgnum": 1555,
			}).Error(err)

			return err
		}

		stmt, err := tx.PrepareContext(ctx, `update users set 
		    password_reset_token = ?,
				password_selector = ?,
				password_verifier = ?,
        password_confirmed_at = ?
		    password = ?,
		    statusc = ?,
        active = ?,
		    updated_at = ?, 
				updated_day = ?, 
				updated_week = ?, 
				updated_month = ?, 
				updated_year = ? where id= ?;`)
		if err != nil {
			log.WithFields(log.Fields{
				"reqid":  requestID,
				"msgnum": 1556,
			}).Error(err)

			err = stmt.Close()
			err = tx.Rollback()
			return err
		}

		_, err = stmt.ExecContext(ctx,
			"",
			"",
			"",
			tn,
			password1,
			common.Active,
			true,
			tn,
			UpdatedDay,
			UpdatedWeek,
			UpdatedMonth,
			UpdatedYear,
			user.ID)
		if err != nil {
			log.WithFields(log.Fields{
				"reqid":  requestID,
				"msgnum": 1557,
			}).Error(err)

			err = stmt.Close()
			err = tx.Rollback()
			return err
		}
		err = stmt.Close()
		if err != nil {
			log.WithFields(log.Fields{
				"reqid":  requestID,
				"msgnum": 1558,
			}).Error(err)

			err = tx.Rollback()
			return err
		}

		err = tx.Commit()
		if err != nil {
			log.WithFields(log.Fields{
				"reqid":  requestID,
				"msgnum": 1559,
			}).Error(err)

			return err
		}
		return nil
	}
}

// ChangePassword - used to update password
func (u *UserService) ChangePassword(ctx context.Context, form *PasswordForm, userEmail string, requestID string) error {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{
			"user":   userEmail,
			"reqid":  requestID,
			"msgnum": 1560,
		}).Error(err)
		return err
	default:
		uUID4byte, err := common.UUIDStrToBytes(form.ID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 1561}).Error(err)
			return err
		}
		db := u.Db
		user := User{}
		row := db.QueryRowContext(ctx, `select id, password from users where uuid4 = ?;`, uUID4byte)
		err = row.Scan(
			&user.ID,
			&user.Password)

		if err != nil {
			log.WithFields(log.Fields{
				"user":   userEmail,
				"reqid":  requestID,
				"msgnum": 1562,
			}).Error(err)
			return err
		}

		err = bcrypt.CompareHashAndPassword(user.Password,
			[]byte(form.CurrentPassword))
		if err != nil {
			log.WithFields(log.Fields{
				"user":   userEmail,
				"reqid":  requestID,
				"msgnum": 1563,
			}).Error(err)
			return err
		}

		tn, tnday, tnweek, tnmonth, tnyear := common.GetTimeDetails()

		UpdatedDay := tnday
		UpdatedWeek := tnweek
		UpdatedMonth := tnmonth
		UpdatedYear := tnyear

		password1, err := HashPassword(form.Password, requestID)
		if err != nil {
			log.WithFields(log.Fields{
				"user":   userEmail,
				"reqid":  requestID,
				"msgnum": 1564,
			}).Error(err)
		}
		stmt, err := db.PrepareContext(ctx, `update users set 
		    password = ?,
		    updated_at = ?, 
				updated_day = ?, 
				updated_week = ?, 
				updated_month = ?, 
				updated_year = ? where id= ?;`)
		if err != nil {
			log.WithFields(log.Fields{
				"user":   userEmail,
				"reqid":  requestID,
				"msgnum": 1565,
			}).Error(err)
			err = stmt.Close()
			return err
		}

		_, err = stmt.ExecContext(ctx,
			password1,
			tn,
			UpdatedDay,
			UpdatedWeek,
			UpdatedMonth,
			UpdatedYear,
			user.ID)
		if err != nil {
			log.WithFields(log.Fields{
				"user":   userEmail,
				"reqid":  requestID,
				"msgnum": 1566,
			}).Error(err)
			err = stmt.Close()
			return err
		}
		err = stmt.Close()
		if err != nil {
			log.WithFields(log.Fields{
				"user":   userEmail,
				"reqid":  requestID,
				"msgnum": 1567,
			}).Error(err)
			return err
		}

		return nil
	}
}

// ChangeEmail - Change Email
func (u *UserService) ChangeEmail(ctx context.Context, form *ChangeEmailForm, hostURL string, userEmail string, requestID string) error {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{
			"user":   userEmail,
			"reqid":  requestID,
			"msgnum": 1568,
		}).Error(err)
		return err
	default:
		db := u.Db
		user, err := u.GetUserByEmail(ctx, form.Email, userEmail, requestID)
		if err != nil {
			log.WithFields(log.Fields{
				"user":   userEmail,
				"reqid":  requestID,
				"msgnum": 1569,
			}).Error(err)
			return err
		}

		selector, verifier, token, err := GenTokenHash(requestID)
		if err != nil {
			log.WithFields(log.Fields{
				"user":   userEmail,
				"reqid":  requestID,
				"msgnum": 1570,
			}).Error(err)
			return err
		}

		tn, tnday, tnweek, tnmonth, tnyear := common.GetTimeDetails()
		tokenExpiry, _ := time.ParseDuration(u.UserOptions.ResetTokenDuration)
		resetExpiry := tn.Add(tokenExpiry)

		UpdatedDay := tnday
		UpdatedWeek := tnweek
		UpdatedMonth := tnmonth
		UpdatedYear := tnyear

		stmt, err := db.PrepareContext(ctx, `update users set 
        new_email = ?,
		    new_email_reset_token = ?,
				new_email_selector = ?,
				new_email_verifier = ?,
        new_email_token_sent_at = ?,
		    new_email_token_expiry = ?,
		    updated_at = ?, 
				updated_day = ?, 
				updated_week = ?, 
				updated_month = ?, 
				updated_year = ? where id= ?;`)
		if err != nil {
			log.WithFields(log.Fields{
				"user":   userEmail,
				"reqid":  requestID,
				"msgnum": 1571,
			}).Error(err)
			err = stmt.Close()
			return err
		}

		_, err = stmt.ExecContext(ctx,
			form.NewEmail,
			token,
			selector,
			verifier,
			tn,
			resetExpiry,
			tn,
			UpdatedDay,
			UpdatedWeek,
			UpdatedMonth,
			UpdatedYear,
			user.ID)
		if err != nil {
			log.WithFields(log.Fields{
				"user":   userEmail,
				"reqid":  requestID,
				"msgnum": 1572,
			}).Error(err)
			err = stmt.Close()
			return err
		}
		err = stmt.Close()

		if err != nil {
			log.WithFields(log.Fields{
				"user":   userEmail,
				"reqid":  requestID,
				"msgnum": 1573,
			}).Error(err)
			return err
		}

		pwd, _ := os.Getwd()
		viewpath := pwd + filepath.FromSlash("/common/views/change_email.html")

		templateData := struct {
			Title string
			URL   string
		}{
			Title: "Change Email",
			URL:   "http://" + hostURL + "/users/change_email/" + token,
		}

		ChangeEmail, err := common.ParseTemplate(viewpath, templateData)
		if err != nil {
			log.WithFields(log.Fields{
				"user":   userEmail,
				"reqid":  requestID,
				"msgnum": 1574,
			}).Error(err)
			return err
		}

		recipient := form.NewEmail
		email := common.Email{
			To:      recipient,
			Subject: "Change Email",
			Body:    ChangeEmail,
		}

		err = common.SendMail(email, u.Mailer)
		if err != nil {
			log.WithFields(log.Fields{
				"user":   userEmail,
				"reqid":  requestID,
				"msgnum": 1575,
			}).Error(err)
			return err
		}

		return nil
	}
}

// ConfirmChangeEmail - Confirm change email
func (u *UserService) ConfirmChangeEmail(ctx context.Context, token string, requestID string) error {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{
			"reqid":  requestID,
			"msgnum": 1576,
		}).Error(err)
		return err
	default:
		db := u.Db
		tn, tnday, tnweek, tnmonth, tnyear := common.GetTimeDetails()

		verifierBytes, selector, err := GetSelectorForPasswdRecoveryToken(token, requestID)
		if err != nil {
			log.WithFields(log.Fields{
				"reqid":  requestID,
				"msgnum": 1577,
			}).Error(err)
			return err
		}
		user := User{}
		row := db.QueryRowContext(ctx, `select id, email, new_email, new_email_selector, new_email_verifier, new_email_token_expiry from users where new_email_selector = ?;`, selector)

		err = row.Scan(
			&user.ID,
			&user.Email,
			&user.NewEmail,
			&user.NewEmailSelector,
			&user.NewEmailVerifier,
			&user.NewEmailTokenExpiry)

		if err != nil {
			log.WithFields(log.Fields{
				"reqid":  requestID,
				"msgnum": 1578,
			}).Error(err)
			return err
		}

		err = ValidatePasswdRecoveryToken(verifierBytes, user.EmailVerifier, user.EmailTokenExpiry, requestID)
		if err != nil {
			log.WithFields(log.Fields{
				"reqid":  requestID,
				"msgnum": 1579,
			}).Error(err)
			return err
		}

		UpdatedDay := tnday
		UpdatedWeek := tnweek
		UpdatedMonth := tnmonth
		UpdatedYear := tnyear

		stmt, err := db.PrepareContext(ctx, `update users set 
        new_email_confirmed_at = ?
		    email = ?,
        new_email = ?,
		    new_email_reset_token = ?,
				new_email_selector = ?,
				new_email_verifier = ?,
		    statusc = ?,
        active = ?,
		    updated_at = ?, 
				updated_day = ?, 
				updated_week = ?, 
				updated_month = ?, 
				updated_year = ? where id= ?;`)
		if err != nil {
			log.WithFields(log.Fields{
				"reqid":  requestID,
				"msgnum": 1580,
			}).Error(err)
			err = stmt.Close()
			return err
		}

		_, err = stmt.ExecContext(ctx,
			tn,
			user.NewEmail,
			"",
			"",
			"",
			"",
			common.Active,
			true,
			tn,
			UpdatedDay,
			UpdatedWeek,
			UpdatedMonth,
			UpdatedYear,
			user.ID)
		if err != nil {
			log.WithFields(log.Fields{
				"reqid":  requestID,
				"msgnum": 1581,
			}).Error(err)
			err = stmt.Close()
			return err
		}
		err = stmt.Close()

		if err != nil {
			log.WithFields(log.Fields{
				"reqid":  requestID,
				"msgnum": 1582,
			}).Error(err)
			return err
		}

		return nil
	}
}

// GetUserByEmail - Get user details by email
func (u *UserService) GetUserByEmail(ctx context.Context, Email string, userEmail string, requestID string) (*User, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{
			"user":   userEmail,
			"reqid":  requestID,
			"msgnum": 1583,
		}).Error(err)
		return nil, err
	default:
		db := u.Db
		user := User{}
		row := db.QueryRowContext(ctx, `select
    id,
		uuid4,
		email,
    username,
		first_name,
		last_name,
		role,
		active,
		statusc,
		created_at,
		updated_at,
		created_day,
		created_week,
		created_month,
		created_year,
		updated_day,
		updated_week,
		updated_month,
		updated_year from users where email = ?;`, Email)

		err := row.Scan(
			&user.ID,
			&user.UUID4,
			&user.Email,
			&user.Username,
			&user.FirstName,
			&user.LastName,
			&user.Role,
			&user.Active,
			/*  StatusDates  */
			&user.Statusc,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.CreatedDay,
			&user.CreatedWeek,
			&user.CreatedMonth,
			&user.CreatedYear,
			&user.UpdatedDay,
			&user.UpdatedWeek,
			&user.UpdatedMonth,
			&user.UpdatedYear)

		if err != nil {
			log.WithFields(log.Fields{
				"user":   userEmail,
				"reqid":  requestID,
				"msgnum": 1584,
			}).Error(err)
			return nil, err
		}
		uUID4Str, err := common.UUIDBytesToStr(user.UUID4)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 1585}).Error(err)
			return nil, err
		}
		user.IDS = uUID4Str
		return &user, nil
	}
}

// GetUser - Get user details by ID
func (u *UserService) GetUser(ctx context.Context, ID string, userEmail string, requestID string) (*User, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{
			"reqid":  requestID,
			"msgnum": 1586,
		}).Error(err)
		return nil, err
	default:
		uUID4byte, err := common.UUIDStrToBytes(ID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 1587}).Error(err)
			return nil, err
		}
		db := u.Db
		user := User{}
		row := db.QueryRowContext(ctx, `select
    id,
		uuid4,
		email,
    username,
		first_name,
		last_name,
		role,
		active,
		statusc,
		created_at,
		updated_at,
		created_day,
		created_week,
		created_month,
		created_year,
		updated_day,
		updated_week,
		updated_month,
		updated_year from users where uuid4 = ?;`, uUID4byte)

		err = row.Scan(
			&user.ID,
			&user.UUID4,
			&user.Email,
			&user.Username,
			&user.FirstName,
			&user.LastName,
			&user.Role,
			&user.Active,
			/*  StatusDates  */
			&user.Statusc,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.CreatedDay,
			&user.CreatedWeek,
			&user.CreatedMonth,
			&user.CreatedYear,
			&user.UpdatedDay,
			&user.UpdatedWeek,
			&user.UpdatedMonth,
			&user.UpdatedYear)

		if err != nil {
			log.WithFields(log.Fields{
				"reqid":  requestID,
				"msgnum": 1588,
			}).Error(err)
			return nil, err
		}
		uUID4Str, err := common.UUIDBytesToStr(user.UUID4)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 1589}).Error(err)
			return nil, err
		}
		user.IDS = uUID4Str
		return &user, nil
	}
}
