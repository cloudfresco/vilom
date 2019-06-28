package userservices

import (
	"context"
	"database/sql"
	"errors"

	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"

	"github.com/cloudfresco/vilom/common"
)

/* error message range: 3300-3999 */

// Ubadge - Ubadge view representation
type Ubadge struct {
	ID    uint   `json:"id,omitempty"`
	UUID4 []byte `json:"-"`
	IDS   string `json:"id_s,omitempty"`

	UbadgeName string `json:"ubadge_name,omitempty"`
	UbadgeDesc string `json:"ubadge_desc,omitempty"`

	common.StatusDates
	Users []*User
}

// UbadgeUser - Ubadge User view representation
type UbadgeUser struct {
	ID    uint   `json:"id,omitempty"`
	UUID4 []byte `json:"-"`
	IDS   string `json:"id_s,omitempty"`

	UbadgeID uint `json:"ubadge_id,omitempty"`
	UserID   uint `json:"user_id,omitempty"`

	common.StatusDates
}

// UbadgeService - For accessing Ubadge services
type UbadgeService struct {
	Config       *common.RedisOptions
	Db           *sql.DB
	RedisClient  *redis.Client
	LimitDefault string
}

// NewUbadgeService - Create Ubadge Service
func NewUbadgeService(config *common.RedisOptions,
	db *sql.DB,
	redisClient *redis.Client,
	limitDefault string) *UbadgeService {
	return &UbadgeService{config, db, redisClient, limitDefault}
}

// UbadgeCursor - used to get ubadges
type UbadgeCursor struct {
	Ubadges    []*Ubadge
	NextCursor string `json:"next_cursor,omitempty"`
}

// GetUbadges - Get Ubadges
func (u *UbadgeService) GetUbadges(ctx context.Context, limit string, nextCursor string, userEmail string, requestID string) (*UbadgeCursor, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 3300}).Error(err)
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
			query = query + "where " + "id <= " + nextCursor + " order by id desc " + "limit " + limit + ";"
		}

		ubadges := []*Ubadge{}
		rows, err := u.Db.QueryContext(ctx, `select 
      id,
			uuid4,
			ubadge_name,
			ubadge_desc,
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
			updated_year from ubadges `+query)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 3301}).Error(err)
			return nil, err
		}

		for rows.Next() {
			ubadge := Ubadge{}
			err = rows.Scan(&ubadge.ID,
				&ubadge.UUID4,
				&ubadge.UbadgeName,
				&ubadge.UbadgeDesc,
				&ubadge.Statusc,
				&ubadge.CreatedAt,
				&ubadge.UpdatedAt,
				&ubadge.CreatedDay,
				&ubadge.CreatedWeek,
				&ubadge.CreatedMonth,
				&ubadge.CreatedYear,
				&ubadge.UpdatedDay,
				&ubadge.UpdatedWeek,
				&ubadge.UpdatedMonth,
				&ubadge.UpdatedYear)
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 3302}).Error(err)
				return nil, err
			}

			uuid4Str, err := common.UUIDBytesToStr(ubadge.UUID4)
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 3303}).Error(err)
				return nil, err
			}
			ubadge.IDS = uuid4Str
			ubadges = append(ubadges, &ubadge)
		}
		err = rows.Close()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 3304}).Error(err)
			return nil, err
		}

		err = rows.Err()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 3305}).Error(err)
			return nil, err
		}
		x := UbadgeCursor{}
		if len(ubadges) != 0 {
			next := ubadges[len(ubadges)-1].ID
			next = next - 1
			nextc := common.EncodeCursor(next)
			x = UbadgeCursor{ubadges, nextc}
		} else {
			x = UbadgeCursor{ubadges, "0"}
		}
		return &x, nil
	}
}

// Create - Create Ubadge
func (u *UbadgeService) Create(ctx context.Context, form *Ubadge, userEmail string, requestID string) (*Ubadge, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 3306}).Error(err)
		return nil, err
	default:
		db := u.Db
		tx, err := db.Begin()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 3307}).Error(err)
			return nil, err
		}

		tn, tnday, tnweek, tnmonth, tnyear := common.GetTimeDetails()
		Ubadge := Ubadge{}
		Ubadge.UUID4, err = common.GetUUIDBytes()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 3308}).Error(err)
			err = tx.Rollback()
			return nil, err
		}
		Ubadge.UbadgeName = form.UbadgeName
		Ubadge.UbadgeDesc = form.UbadgeDesc
		Ubadge.Statusc = common.Active
		Ubadge.CreatedAt = tn
		Ubadge.UpdatedAt = tn
		Ubadge.CreatedDay = tnday
		Ubadge.CreatedWeek = tnweek
		Ubadge.CreatedMonth = tnmonth
		Ubadge.CreatedYear = tnyear
		Ubadge.UpdatedDay = tnday
		Ubadge.UpdatedWeek = tnweek
		Ubadge.UpdatedMonth = tnmonth
		Ubadge.UpdatedYear = tnyear

		err = u.InsertUbadge(ctx, tx, &Ubadge, userEmail, requestID)

		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 3309}).Error(err)
			err = tx.Rollback()
			return nil, err
		}

		err = tx.Commit()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 3310}).Error(err)
			err = tx.Rollback()
			return nil, err
		}

		return &Ubadge, nil
	}
}

// AddUserToGroup - Add user to ubadge
func (u *UbadgeService) AddUserToGroup(ctx context.Context, form *UbadgeUser, ID string, userEmail string, requestID string) error {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 3311}).Error(err)
		return err
	default:
		db := u.Db
		ubadge, err := u.GetUbadge(ctx, ID, userEmail, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 3312}).Error(err)
			return err
		}

		tx, err := db.Begin()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 3313}).Error(err)
			return err
		}
		tn, tnday, tnweek, tnmonth, tnyear := common.GetTimeDetails()

		Uguser := UbadgeUser{}
		Uguser.UUID4, err = common.GetUUIDBytes()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 3314}).Error(err)
			err = tx.Rollback()
			return err
		}

		Uguser.UbadgeID = ubadge.ID
		Uguser.UserID = form.UserID
		Uguser.Statusc = common.Active
		Uguser.CreatedAt = tn
		Uguser.UpdatedAt = tn
		Uguser.CreatedDay = tnday
		Uguser.CreatedWeek = tnweek
		Uguser.CreatedMonth = tnmonth
		Uguser.CreatedYear = tnyear
		Uguser.UpdatedDay = tnday
		Uguser.UpdatedWeek = tnweek
		Uguser.UpdatedMonth = tnmonth
		Uguser.UpdatedYear = tnyear

		err = u.InsertUbadgeUser(ctx, tx, &Uguser, userEmail, requestID)

		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 3315}).Error(err)
			err = tx.Rollback()
			return err
		}

		err = tx.Commit()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 3316}).Error(err)
			err = tx.Rollback()
			return err
		}
		return nil
	}
}

//Update - Update Ubadge
func (u *UbadgeService) Update(ctx context.Context, ID string, form *Ubadge, UserID string, userEmail string, requestID string) error {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 3358}).Error(err)
		return err
	default:
		ubadge, err := u.GetUbadge(ctx, ID, userEmail, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 3359}).Error(err)
			return err
		}

		db := u.Db
		tx, err := db.Begin()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 3360}).Error(err)
			return err
		}

		tn, tnday, tnweek, tnmonth, tnyear := common.GetTimeDetails()
		stmt, err := tx.PrepareContext(ctx, `update ubadges set 
		  ubadge_name = ?,
      ubadge_desc = ?,
			updated_at = ?, 
			updated_day = ?, 
			updated_week = ?, 
			updated_month = ?, 
			updated_year = ? where id = ?;`)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 3361}).Error(err)
			err = stmt.Close()
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 3362}).Error(err)
				err = tx.Rollback()
				return err
			}
			err = tx.Rollback()
			return err
		}
		_, err = stmt.ExecContext(ctx,
			form.UbadgeName,
			form.UbadgeDesc,
			tn,
			tnday,
			tnweek,
			tnmonth,
			tnyear,
			ubadge.ID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 3363}).Error(err)
			err = stmt.Close()
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 3364}).Error(err)
				err = tx.Rollback()
				return err
			}
			err = tx.Rollback()
			return err
		}
		err = stmt.Close()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 3365}).Error(err)
			return err
		}

		err = tx.Commit()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 3366}).Error(err)
			return err
		}
		return nil
	}
}

// InsertUbadge - Insert Ubadge details into database
func (u *UbadgeService) InsertUbadge(ctx context.Context, tx *sql.Tx, Ubadge *Ubadge, userEmail string, requestID string) error {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 3317}).Error(err)
		return err
	default:
		stmt, err := tx.PrepareContext(ctx, `insert into ubadges
	  (
		uuid4,
		ubadge_name,
		ubadge_desc,
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
					?,?,?,?,?);`)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 3318}).Error(err)
			err = stmt.Close()
			return err
		}
		res, err := stmt.ExecContext(ctx,
			Ubadge.UUID4,
			Ubadge.UbadgeName,
			Ubadge.UbadgeDesc,
			Ubadge.Statusc,
			Ubadge.CreatedAt,
			Ubadge.UpdatedAt,
			Ubadge.CreatedDay,
			Ubadge.CreatedWeek,
			Ubadge.CreatedMonth,
			Ubadge.CreatedYear,
			Ubadge.UpdatedDay,
			Ubadge.UpdatedWeek,
			Ubadge.UpdatedMonth,
			Ubadge.UpdatedYear)

		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 3319}).Error(err)
			err = stmt.Close()
			return err
		}
		uID, err := res.LastInsertId()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 3320}).Error(err)
			err = stmt.Close()
			return err
		}
		Ubadge.ID = uint(uID)
		uuid4Str, err := common.UUIDBytesToStr(Ubadge.UUID4)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 3321}).Error(err)
			return err
		}
		Ubadge.IDS = uuid4Str
		err = stmt.Close()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 3322}).Error(err)
			return err
		}

		return nil
	}
}

// Delete - Delele Ubadge
func (u *UbadgeService) Delete(ctx context.Context, ID string, userEmail string, requestID string) error {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 3323}).Error(err)
		return err
	default:
		uuid4byte, err := common.UUIDStrToBytes(ID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 3324}).Error(err)
			return err
		}
		db := u.Db
		tx, err := db.Begin()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 3325}).Error(err)
			return err
		}
		stmt, err := tx.PrepareContext(ctx, "delete from ubadges where uuid4= ?;")
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 3326}).Error(err)
			err = stmt.Close()
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 3353}).Error(err)
				err = tx.Rollback()
				return err
			}
			err = tx.Rollback()
			return err
		}

		_, err = stmt.ExecContext(ctx, uuid4byte)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 3354}).Error(err)
			err = stmt.Close()
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 3355}).Error(err)
				err = tx.Rollback()
				return err
			}
			err = tx.Rollback()
			return err
		}
		err = stmt.Close()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 3327}).Error(err)
			err = tx.Rollback()
			return err
		}
		err = tx.Commit()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 3328}).Error(err)
			err = tx.Rollback()
			return err
		}
		return nil
	}
}

// GetUbadge - Get Ubadge Details
func (u *UbadgeService) GetUbadge(ctx context.Context, ID string, userEmail string, requestID string) (*Ubadge, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 3329}).Error(err)
		return nil, err
	default:
		uuid4byte, err := common.UUIDStrToBytes(ID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 3330}).Error(err)
			return nil, err
		}
		db := u.Db
		poh := Ubadge{}
		rows, err := db.QueryContext(ctx, `select 
    p.id,
		p.uuid4,
		p.ubadge_name,
		p.ubadge_desc,
		p.statusc,
		p.created_at,
		p.updated_at,
		p.created_day,
		p.created_week,
		p.created_month,
		p.created_year,
		p.updated_day,
		p.updated_week,
		p.updated_month,
		p.updated_year,
    v.id,
		v.uuid4,
    v.auth_token,
		v.email,
		v.first_name,
		v.last_name,
		v.role,
		v.password,
		v.active,
		v.email_confirmation_token,
		v.email_token_sent_at,
		v.email_token_expiry,
		v.email_confirmed_at,
		v.new_email,
		v.new_email_reset_token,
		v.new_email_token_sent_at,
		v.new_email_token_expiry,
		v.new_email_confirmed_at,
		v.password_reset_token,
		v.password_token_sent_at,
		v.password_token_expiry,
		v.password_confirmed_at,
		v.timezone,
		v.sign_in_count,
		v.current_sign_in_at,
		v.last_sign_in_at,
		v.statusc,
		v.created_at,
		v.updated_at,
		v.created_day,
		v.created_week,
		v.created_month,
		v.created_year,
		v.updated_day,
		v.updated_week,
		v.updated_month,
		v.updated_year from ubadges p inner join ubadges_users ubu on (p.id = ubu.ubadge_id) inner join users v on (ubu.user_id = v.id) where p.uuid4 = ?`, uuid4byte)

		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 3331}).Error(err)
			return nil, err
		}
		for rows.Next() {
			user := User{}
			err = rows.Scan(
				&poh.ID,
				&poh.UUID4,
				&poh.UbadgeName,
				&poh.UbadgeDesc,
				&poh.Statusc,
				&poh.CreatedAt,
				&poh.UpdatedAt,
				&poh.CreatedDay,
				&poh.CreatedWeek,
				&poh.CreatedMonth,
				&poh.CreatedYear,
				&poh.UpdatedDay,
				&poh.UpdatedWeek,
				&poh.UpdatedMonth,
				&poh.UpdatedYear,
				&user.ID,
				&user.UUID4,
				&user.AuthToken,
				&user.Email,
				&user.FirstName,
				&user.LastName,
				&user.Role,
				&user.Password,
				&user.Active,
				&user.EmailConfirmationToken,
				&user.EmailTokenSentAt,
				&user.EmailTokenExpiry,
				&user.EmailConfirmedAt,
				&user.NewEmail,
				&user.NewEmailResetToken,
				&user.NewEmailTokenSentAt,
				&user.NewEmailTokenExpiry,
				&user.NewEmailConfirmedAt,
				&user.PasswordResetToken,
				&user.PasswordTokenSentAt,
				&user.PasswordTokenExpiry,
				&user.PasswordConfirmedAt,
				&user.Timezone,
				&user.SignInCount,
				&user.CurrentSignInAt,
				&user.LastSignInAt,
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
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 3332}).Error(err)
				return nil, err
			}
			uuid4Str1, err := common.UUIDBytesToStr(poh.UUID4)
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 3333}).Error(err)
				return nil, err
			}
			poh.IDS = uuid4Str1

			uuid4Str, err := common.UUIDBytesToStr(user.UUID4)
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 3334}).Error(err)
				return nil, err
			}
			user.IDS = uuid4Str
			poh.Users = append(poh.Users, &user)
		}

		err = rows.Close()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 3335}).Error(err)
			return nil, err
		}

		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 3336}).Error(err)
			return nil, err
		}

		return &poh, nil
	}
}

// GetUbadgeByID - Get Ubadge by ID
func (u *UbadgeService) GetUbadgeByID(ctx context.Context, ID string, userEmail string, requestID string) (*Ubadge, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 3337}).Error(err)
		return nil, err
	default:
		uuid4byte, err := common.UUIDStrToBytes(ID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 3338}).Error(err)
			return nil, err
		}
		db := u.Db
		Ubadge := Ubadge{}
		row := db.QueryRowContext(ctx, `select
    id,
		uuid4,
		ubadge_name,
		ubadge_desc,
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
		updated_year from ubadges where uuid4 = ?;`, uuid4byte)

		err = row.Scan(
			&Ubadge.ID,
			&Ubadge.UUID4,
			&Ubadge.UbadgeName,
			&Ubadge.UbadgeDesc,
			&Ubadge.Statusc,
			&Ubadge.CreatedAt,
			&Ubadge.UpdatedAt,
			&Ubadge.CreatedDay,
			&Ubadge.CreatedWeek,
			&Ubadge.CreatedMonth,
			&Ubadge.CreatedYear,
			&Ubadge.UpdatedDay,
			&Ubadge.UpdatedWeek,
			&Ubadge.UpdatedMonth,
			&Ubadge.UpdatedYear)

		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 3339}).Error(err)
			return nil, err
		}
		uuid4Str, err := common.UUIDBytesToStr(Ubadge.UUID4)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 3340}).Error(err)
			return nil, err
		}
		Ubadge.IDS = uuid4Str
		return &Ubadge, nil
	}
}

// DeleteUserFromGroup - Delete user from Ubadge
func (u *UbadgeService) DeleteUserFromGroup(ctx context.Context, form *UbadgeUser, ID string, userEmail string, requestID string) error {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 3341}).Error(err)
		return err
	default:
		db := u.Db
		tx, err := db.Begin()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 3342}).Error(err)
			return err
		}
		stmt, err := tx.PrepareContext(ctx, `delete from ubadges_users where user_id= ? and ubadge_id = ?;`)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 3343}).Error(err)
			err = stmt.Close()
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 3356}).Error(err)
				err = tx.Rollback()
				return err
			}
			err = tx.Rollback()
			return err
		}

		_, err = stmt.ExecContext(ctx, form.UserID, ID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 3344}).Error(err)
			err = stmt.Close()
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 3357}).Error(err)
				err = tx.Rollback()
				return err
			}
			err = tx.Rollback()
			return err
		}
		err = stmt.Close()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 3345}).Error(err)
			err = tx.Rollback()
			return err
		}
		err = tx.Commit()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 3346}).Error(err)
			err = tx.Rollback()
			return err
		}
		return nil
	}
}

// InsertUbadgeUser - Insert Ubadge User details into database
func (u *UbadgeService) InsertUbadgeUser(ctx context.Context, tx *sql.Tx, Uguser *UbadgeUser, userEmail string, requestID string) error {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 3347}).Error(err)
		return err
	default:
		stmt, err := tx.PrepareContext(ctx, `insert into ubadges_users
	  (
		uuid4,
		ubadge_id,
		user_id,
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
					?,?,?,?);`)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 3348}).Error(err)
			return err
		}
		res, err := stmt.ExecContext(ctx,
			Uguser.UUID4,
			Uguser.UbadgeID,
			Uguser.UserID,
			Uguser.Statusc,
			Uguser.CreatedAt,
			Uguser.UpdatedAt,
			Uguser.CreatedDay,
			Uguser.CreatedWeek,
			Uguser.CreatedMonth,
			Uguser.CreatedYear,
			Uguser.UpdatedDay,
			Uguser.UpdatedWeek,
			Uguser.UpdatedMonth,
			Uguser.UpdatedYear)

		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 3349}).Error(err)
			err = stmt.Close()
			return err
		}
		uID, err := res.LastInsertId()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 3350}).Error(err)
			err = stmt.Close()
			return err
		}
		Uguser.ID = uint(uID)
		uuid4Str, err := common.UUIDBytesToStr(Uguser.UUID4)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 3351}).Error(err)
			return err
		}
		Uguser.IDS = uuid4Str

		err = stmt.Close()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 3352}).Error(err)
			return err
		}

		return nil
	}
}
