package userservices

import (
	"context"
	"database/sql"
	"errors"
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
	ID        uint   `json:"id,omitempty"`
	UUID4     []byte `json:"-"`
	IDS       string `json:"id_s,omitempty"`
	AuthToken string `json:"auth_token,omitempty"`

	Email     string `json:"email,omitempty"`
	Username  string `json:"username,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Role      string `json:"role,omitempty"`
	Password  []byte `json:"password,omitempty"`
	Active    bool   `json:"active,omitempty"`

	EmailConfirmationToken string    `json:"email_confirmation_token,omitempty"`
	EmailSelector          string    `json:"email_selector,omitempty"`
	EmailVerifier          string    `json:"email_verifier,omitempty"`
	EmailTokenSentAt       time.Time `json:"email_token_sent_at,omitempty"`
	EmailTokenExpiry       time.Time `json:"email_token_expiry,omitempty"`
	EmailConfirmedAt       time.Time `json:"email_confirmed_at,omitempty"`

	NewEmail            string    `json:"new_email,omitempty"`
	NewEmailResetToken  string    `json:"new_email_reset_token,omitempty"`
	NewEmailSelector    string    `json:"new_email_selector,omitempty"`
	NewEmailVerifier    string    `json:"new_email_verifier,omitempty"`
	NewEmailTokenSentAt time.Time `json:"new_email_token_sent_at,omitempty"`
	NewEmailTokenExpiry time.Time `json:"new_email_token_expiry,omitempty"`
	NewEmailConfirmedAt time.Time `json:"new_email_confirmed_at,omitempty"`

	PasswordResetToken  string    `json:"password_reset_token,omitempty"`
	PasswordSelector    string    `json:"password_selector,omitempty"`
	PasswordVerifier    string    `json:"password_verifier,omitempty"`
	PasswordTokenSentAt time.Time `json:"password_token_sent_at,omitempty"`
	PasswordTokenExpiry time.Time `json:"password_token_expiry,omitempty"`
	PasswordConfirmedAt time.Time `json:"password_confirmed_at,omitempty"`

	Timezone        string    `json:"timezone,omitempty"`
	SignInCount     uint      `json:"sign_in_count,omitempty"`
	CurrentSignInAt time.Time `json:"current_sign_in_at,omitempty"`
	LastSignInAt    time.Time `json:"last_sign_in_at,omitempty"`

	common.StatusDates

	/* used only for logic purpose */
	Roles       []string `json:"roles,omitempty"`
	PasswordS   string   `json:"password_s,omitempty"`
	Tokenstring string   `json:"tokenstring,omitempty"`
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
	Password        string
	ConfirmPassword string
	CurrentPassword string
	ID              string
}

// ChangeEmailForm - used for Change Email
type ChangeEmailForm struct {
	Email    string
	NewEmail string
}

// ForgotPasswordForm - used for forgot password
type ForgotPasswordForm struct {
	Email string
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

// UserCursor - used for getting users list
type UserCursor struct {
	Users      []*User
	NextCursor string `json:"next_cursor,omitempty"`
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
			uuid4Str, err := common.UUIDBytesToStr(user.UUID4)
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 1510}).Error(err)
				return nil, err
			}
			user.IDS = uuid4Str
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

		password1, err := common.HashPassword(form.PasswordS, requestID)
		if err != nil {
			log.WithFields(log.Fields{
				"reqid":  requestID,
				"msgnum": 1522,
			}).Error(err)

			return nil, err
		}

		selector, verifier, token, err := common.GenTokenHash(requestID)
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
		err = u.InsertUser(ctx, tx, &user, hostURL, requestID)
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
		return &user, nil
	}
}

//Update - Update User
func (u *UserService) Update(ctx context.Context, ID string, form *User, UserID string, userEmail string, requestID string) error {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 1592}).Error(err)
		return err
	default:
		user, err := u.GetUser(ctx, ID, userEmail, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 1593}).Error(err)
			return err
		}

		db := u.Db
		tx, err := db.Begin()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 1594}).Error(err)
			return err
		}

		tn, tnday, tnweek, tnmonth, tnyear := common.GetTimeDetails()
		stmt, err := tx.PrepareContext(ctx, `update users set 
		  first_name = ?,
      last_name = ?,
			updated_at = ?, 
			updated_day = ?, 
			updated_week = ?, 
			updated_month = ?, 
			updated_year = ? where id = ?;`)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 1595}).Error(err)
			err = stmt.Close()
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 1596}).Error(err)
				err = tx.Rollback()
				return err
			}
			err = tx.Rollback()
			return err
		}
		_, err = stmt.ExecContext(ctx,
			form.FirstName,
			form.LastName,
			tn,
			tnday,
			tnweek,
			tnmonth,
			tnyear,
			user.ID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 1597}).Error(err)
			err = stmt.Close()
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 1598}).Error(err)
				err = tx.Rollback()
				return err
			}
			err = tx.Rollback()
			return err
		}
		err = stmt.Close()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 1599}).Error(err)
			return err
		}

		err = tx.Commit()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 1600}).Error(err)
			return err
		}
		return nil
	}
}

// InsertUser - Insert User details to database
func (u *UserService) InsertUser(ctx context.Context, tx *sql.Tx, user *User, hostURL string, requestID string) error {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{
			"reqid":  requestID,
			"msgnum": 1528,
		}).Error(err)

		return err
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

			return err
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
			return err
		}
		uID, err := res.LastInsertId()
		if err != nil {
			log.WithFields(log.Fields{
				"reqid":  requestID,
				"msgnum": 1531,
			}).Error(err)

			err = stmt.Close()
			return err
		}
		user.ID = uint(uID)
		uuid4Str, err := common.UUIDBytesToStr(user.UUID4)
		if err != nil {
			log.WithFields(log.Fields{"reqid": requestID, "msgnum": 1532}).Error(err)
			return err
		}
		user.IDS = uuid4Str
		err = stmt.Close()
		if err != nil {
			log.WithFields(log.Fields{
				"reqid":  requestID,
				"msgnum": 1533,
			}).Error(err)

			return err
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

			return err
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

			return err
		}
		return nil
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
		verifierBytes, selector, err := common.GetSelectorForPasswdRecoveryToken(token, requestID)
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

		err = common.ValidatePasswdRecoveryToken(verifierBytes, user.EmailVerifier, user.EmailTokenExpiry, requestID)
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

		selector, verifier, token, err := common.GenTokenHash(requestID)
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
		password1, err := common.HashPassword(form.Password, requestID)
		if err != nil {
			log.WithFields(log.Fields{
				"reqid":  requestID,
				"msgnum": 1551,
			}).Error(err)

		}

		verifierBytes, selector, err := common.GetSelectorForPasswdRecoveryToken(token, requestID)
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

		err = common.ValidatePasswdRecoveryToken(verifierBytes, user.PasswordVerifier, user.PasswordTokenExpiry, requestID)
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
			if err != nil {
				log.WithFields(log.Fields{
					"reqid":  requestID,
					"msgnum": 1590,
				}).Error(err)

				err = tx.Rollback()
				return err
			}
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
			if err != nil {
				log.WithFields(log.Fields{
					"reqid":  requestID,
					"msgnum": 1591,
				}).Error(err)

				err = tx.Rollback()
				return err
			}
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
		uuid4byte, err := common.UUIDStrToBytes(form.ID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 1561}).Error(err)
			return err
		}
		db := u.Db
		user := User{}
		row := db.QueryRowContext(ctx, `select id, password from users where uuid4 = ?;`, uuid4byte)
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

		password1, err := common.HashPassword(form.Password, requestID)
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

		selector, verifier, token, err := common.GenTokenHash(requestID)
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

		verifierBytes, selector, err := common.GetSelectorForPasswdRecoveryToken(token, requestID)
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

		err = common.ValidatePasswdRecoveryToken(verifierBytes, user.EmailVerifier, user.EmailTokenExpiry, requestID)
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
		uuid4Str, err := common.UUIDBytesToStr(user.UUID4)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 1585}).Error(err)
			return nil, err
		}
		user.IDS = uuid4Str
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
		uuid4byte, err := common.UUIDStrToBytes(ID)
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
		updated_year from users where uuid4 = ?;`, uuid4byte)

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
		uuid4Str, err := common.UUIDBytesToStr(user.UUID4)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 1589}).Error(err)
			return nil, err
		}
		user.IDS = uuid4Str
		return &user, nil
	}
}
