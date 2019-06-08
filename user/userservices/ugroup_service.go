package userservices

import (
	"context"
	"database/sql"
	"errors"

	log "github.com/Sirupsen/logrus"
	"github.com/go-redis/redis"

	"github.com/cloudfresco/vilom/common"
)

// Ugroup - Ugroup view representation
type Ugroup struct {
	ID    uint
	UUID4 []byte
	IDS   string

	UgroupName string
	UgroupDesc string
	Levelc     uint
	ParentID   uint
	NumChd     uint

	common.StatusDates

	Users []*User
}

// UgroupChd - UgroupChd view representation
type UgroupChd struct {
	ID          uint
	UgroupID    uint
	UgroupChdID uint

	common.StatusDates
}

// UgroupUser - UgroupUser view representation
type UgroupUser struct {
	ID    uint
	UUID4 []byte
	IDS   string

	UgroupID uint
	UserID   uint

	common.StatusDates
}

// UgroupService - For accessing Ugroup services
type UgroupService struct {
	Config       *common.RedisOptions
	Db           *sql.DB
	RedisClient  *redis.Client
	LimitDefault string
}

// NewUgroupService - Create Ugroup Service
func NewUgroupService(config *common.RedisOptions,
	db *sql.DB,
	redisClient *redis.Client,
	limitDefault string) *UgroupService {
	return &UgroupService{config, db, redisClient, limitDefault}
}

// UgroupCursor - used to get groups
type UgroupCursor struct {
	Ugroups    []*Ugroup
	NextCursor string
}

// GetUgroups - Get Groups
func (u *UgroupService) GetUgroups(ctx context.Context, limit string, nextCursor string, userEmail string, requestID string) (*UgroupCursor, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2300}).Error(err)
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

		ugroups := []*Ugroup{}
		rows, err := u.Db.QueryContext(ctx, `select 
      id,
			uuid4,
			ugroup_name,
			ugroup_desc,
			levelc,
			parent_id,
			num_chd,
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
			updated_year from ugroups `+query)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2301}).Error(err)
		}

		for rows.Next() {
			ug := Ugroup{}
			err = rows.Scan(&ug.ID,
				&ug.UUID4,
				&ug.UgroupName,
				&ug.UgroupDesc,
				&ug.Levelc,
				&ug.ParentID,
				&ug.NumChd,
				&ug.Statusc,
				&ug.CreatedAt,
				&ug.UpdatedAt,
				&ug.CreatedDay,
				&ug.CreatedWeek,
				&ug.CreatedMonth,
				&ug.CreatedYear,
				&ug.UpdatedDay,
				&ug.UpdatedWeek,
				&ug.UpdatedMonth,
				&ug.UpdatedYear)
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2302}).Error(err)
				return nil, err
			}
			uUID4Str, err := common.UUIDBytesToStr(ug.UUID4)
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2303}).Error(err)
				return nil, err
			}
			ug.IDS = uUID4Str
			ugroups = append(ugroups, &ug)
		}
		err = rows.Close()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2304}).Error(err)
			return nil, err
		}

		err = rows.Err()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2305}).Error(err)
			return nil, err
		}
		x := UgroupCursor{}
		if len(ugroups) != 0 {
			next := ugroups[len(ugroups)-1].ID
			next = next - 1
			nextc := common.EncodeCursor(next)
			x = UgroupCursor{ugroups, nextc}
		} else {
			x = UgroupCursor{ugroups, "0"}
		}
		return &x, nil
	}
}

// Create - Create ugroup
func (u *UgroupService) Create(ctx context.Context, form *Ugroup, userEmail string, requestID string) (*Ugroup, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2306}).Error(err)
		return nil, err
	default:
		db := u.Db
		tx, err := db.Begin()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2307}).Error(err)
			return nil, err
		}

		tn, tnday, tnweek, tnmonth, tnyear := common.GetTimeDetails()

		ug := Ugroup{}
		ug.UUID4, err = common.GetUUIDBytes()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2308}).Error(err)
			return nil, err
		}
		ug.UgroupName = form.UgroupName
		ug.UgroupDesc = form.UgroupDesc
		ug.Levelc = 0
		ug.ParentID = uint(0)
		ug.NumChd = 0
		ug.Statusc = common.Active
		ug.CreatedAt = tn
		ug.UpdatedAt = tn
		ug.CreatedDay = tnday
		ug.CreatedWeek = tnweek
		ug.CreatedMonth = tnmonth
		ug.CreatedYear = tnyear
		ug.UpdatedDay = tnday
		ug.UpdatedWeek = tnweek
		ug.UpdatedMonth = tnmonth
		ug.UpdatedYear = tnyear

		ugrp, err := u.InsertUgroup(ctx, tx, ug, userEmail, requestID)

		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2309}).Error(err)
			err = tx.Rollback()
			return nil, err
		}

		err = tx.Commit()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2310}).Error(err)
			err = tx.Rollback()
			return nil, err
		}

		return ugrp, nil
	}
}

// CreateChild - Create child of ugroup
func (u *UgroupService) CreateChild(ctx context.Context, form *Ugroup, userEmail string, requestID string) (*Ugroup, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2311}).Error(err)
		return nil, err
	default:
		parent, err := u.GetUgroupByIDuint(ctx, form.ParentID, userEmail, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2312}).Error(err)
			return nil, err
		}

		db := u.Db
		tx, err := db.Begin()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2313}).Error(err)
			return nil, err
		}

		tn, tnday, tnweek, tnmonth, tnyear := common.GetTimeDetails()
		ug := Ugroup{}
		ug.UUID4, err = common.GetUUIDBytes()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2314}).Error(err)
			return nil, err
		}
		ug.UgroupName = form.UgroupName
		ug.UgroupDesc = form.UgroupDesc
		ug.Levelc = parent.Levelc + 1
		ug.ParentID = parent.ID
		ug.NumChd = 0
		ug.Statusc = common.Active
		ug.CreatedAt = tn
		ug.UpdatedAt = tn
		ug.CreatedDay = tnday
		ug.CreatedWeek = tnweek
		ug.CreatedMonth = tnmonth
		ug.CreatedYear = tnyear
		ug.UpdatedDay = tnday
		ug.UpdatedWeek = tnweek
		ug.UpdatedMonth = tnmonth
		ug.UpdatedYear = tnyear

		ugrp, err := u.InsertUgroup(ctx, tx, ug, userEmail, requestID)

		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2315}).Error(err)
			err = tx.Rollback()
			return nil, err
		}

		Ugroupchd := UgroupChd{}
		Ugroupchd.UgroupID = parent.ID
		Ugroupchd.UgroupChdID = ugrp.ID
		Ugroupchd.Statusc = common.Active
		Ugroupchd.CreatedAt = tn
		Ugroupchd.UpdatedAt = tn
		Ugroupchd.CreatedDay = tnday
		Ugroupchd.CreatedWeek = tnweek
		Ugroupchd.CreatedMonth = tnmonth
		Ugroupchd.CreatedYear = tnyear
		Ugroupchd.UpdatedDay = tnday
		Ugroupchd.UpdatedWeek = tnweek
		Ugroupchd.UpdatedMonth = tnmonth
		Ugroupchd.UpdatedYear = tnyear

		_, err = u.InsertChild(ctx, tx, Ugroupchd, userEmail, requestID)

		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2316}).Error(err)
			err = tx.Rollback()
			return nil, err
		}

		UpdatedDay := tnday
		UpdatedWeek := tnweek
		UpdatedMonth := tnmonth
		UpdatedYear := tnyear

		stmt, err := tx.PrepareContext(ctx, `update ugroups set 
				  num_chd = ?,
				  updated_at = ?, 
					updated_day = ?, 
					updated_week = ?, 
					updated_month = ?, 
					updated_year = ? where id = ?;`)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2317}).Error(err)
			err = stmt.Close()
			err = tx.Rollback()
			return nil, err
		}

		_, err = stmt.ExecContext(ctx,
			parent.NumChd+1,
			tn,
			UpdatedDay,
			UpdatedWeek,
			UpdatedMonth,
			UpdatedYear,
			parent.ID)

		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2318}).Error(err)
			err = stmt.Close()
			err = tx.Rollback()
			return nil, err
		}

		err = stmt.Close()

		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2319}).Error(err)
			err = tx.Rollback()
			return nil, err
		}

		err = tx.Commit()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2320}).Error(err)
			err = tx.Rollback()
			return nil, err
		}
		return ugrp, nil
	}
}

// AddUserToGroup - Add user to ugroup
func (u *UgroupService) AddUserToGroup(ctx context.Context, form *UgroupUser, ID string, userEmail string, requestID string) error {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2321}).Error(err)
		return err
	default:
		db := u.Db
		ug, err := u.GetUgroupByID(ctx, ID, userEmail, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2322}).Error(err)
			return err
		}

		if ug.NumChd > 0 {
			err = errors.New("Cannot add user to group")
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2323}).Error(err)
			return err
		}

		tx, err := db.Begin()
		tn, tnday, tnweek, tnmonth, tnyear := common.GetTimeDetails()

		Uguser := UgroupUser{}
		Uguser.UUID4, err = common.GetUUIDBytes()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2324}).Error(err)
			err = tx.Rollback()
			return err
		}
		Uguser.UgroupID = ug.ID
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

		_, err = u.InsertUgroupUser(ctx, tx, Uguser, userEmail, requestID)

		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2325}).Error(err)
			err = tx.Rollback()
			return err
		}

		err = tx.Commit()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2326}).Error(err)
			err = tx.Rollback()
			return err
		}
		return nil
	}
}

// InsertUgroup - Insert Ugroup details into database
func (u *UgroupService) InsertUgroup(ctx context.Context, tx *sql.Tx, ug Ugroup, userEmail string, requestID string) (*Ugroup, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2327}).Error(err)
		return nil, err
	default:
		stmt, err := tx.PrepareContext(ctx, `insert into ugroups
	  (
		uuid4,
		ugroup_name,
		ugroup_desc,
		levelc,
		parent_id,
		num_chd,
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
					?,?,?,?,?,?,?);`)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2328}).Error(err)
			return nil, err
		}
		res, err := stmt.ExecContext(ctx,
			ug.UUID4,
			ug.UgroupName,
			ug.UgroupDesc,
			ug.Levelc,
			ug.ParentID,
			ug.NumChd,
			ug.Statusc,
			ug.CreatedAt,
			ug.UpdatedAt,
			ug.CreatedDay,
			ug.CreatedWeek,
			ug.CreatedMonth,
			ug.CreatedYear,
			ug.UpdatedDay,
			ug.UpdatedWeek,
			ug.UpdatedMonth,
			ug.UpdatedYear)

		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2329}).Error(err)
			err = stmt.Close()
			return nil, err
		}
		uID, err := res.LastInsertId()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2330}).Error(err)
			return nil, err
		}
		ug.ID = uint(uID)
		uUID4Str, err := common.UUIDBytesToStr(ug.UUID4)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2331}).Error(err)
			return nil, err
		}
		ug.IDS = uUID4Str
		err = stmt.Close()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2332}).Error(err)
			return nil, err
		}
		return &ug, nil
	}
}

// InsertChild - Insert Child Ugroup details into database
func (u *UgroupService) InsertChild(ctx context.Context, tx *sql.Tx, ugroupchd UgroupChd, userEmail string, requestID string) (*UgroupChd, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2333}).Error(err)
		return nil, err
	default:
		stmt, err := tx.PrepareContext(ctx, `insert into ugroup_chds
	  ( 
		ugroup_id,
		ugroup_chd_id,
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
					?,?,?);`)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2334}).Error(err)
			return nil, err
		}
		res, err := stmt.ExecContext(ctx,
			ugroupchd.UgroupID,
			ugroupchd.UgroupChdID,
			ugroupchd.Statusc,
			ugroupchd.CreatedAt,
			ugroupchd.UpdatedAt,
			ugroupchd.CreatedDay,
			ugroupchd.CreatedWeek,
			ugroupchd.CreatedMonth,
			ugroupchd.CreatedYear,
			ugroupchd.UpdatedDay,
			ugroupchd.UpdatedWeek,
			ugroupchd.UpdatedMonth,
			ugroupchd.UpdatedYear)

		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2335}).Error(err)
			err = stmt.Close()
			return nil, err
		}
		uID, err := res.LastInsertId()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2336}).Error(err)
			err = stmt.Close()
			return nil, err
		}
		ugroupchd.ID = uint(uID)
		err = stmt.Close()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2337}).Error(err)
			return nil, err
		}
		return &ugroupchd, nil
	}
}

// Delete - Delete ugroup
func (u *UgroupService) Delete(ctx context.Context, ID string, userEmail string, requestID string) error {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2338}).Error(err)
		return err
	default:
		uUID4byte, err := common.UUIDStrToBytes(ID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2339}).Error(err)
			return err
		}
		db := u.Db
		tx, err := db.Begin()
		stmt, err := tx.PrepareContext(ctx, "delete from ugroups where uuid4= ?;")
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2340}).Error(err)
			err = tx.Rollback()
			return err
		}

		_, err = stmt.ExecContext(ctx, uUID4byte)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2341}).Error(err)
			err = stmt.Close()
			err = tx.Rollback()
			return err
		}

		err = stmt.Close()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2342}).Error(err)
			err = tx.Rollback()
			return err
		}

		err = tx.Commit()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2343}).Error(err)
			err = tx.Rollback()
			return err
		}
		return nil
	}
}

// GetUgroup - Get ugroup details with users by ID
func (u *UgroupService) GetUgroup(ctx context.Context, ID string, userEmail string, requestID string) (*Ugroup, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2344}).Error(err)
		return nil, err
	default:
		uUID4byte, err := common.UUIDStrToBytes(ID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2345}).Error(err)
			return nil, err
		}
		db := u.Db
		poh := Ugroup{}

		rows, err := db.QueryContext(ctx, `select 
    ug.id,
		ug.uuid4,
		ug.ugroup_name,
		ug.ugroup_desc,
		ug.levelc,
		ug.parent_id,
		ug.num_chd,
		ug.statusc,
		ug.created_at,
		ug.updated_at,
		ug.created_day,
		ug.created_week,
		ug.created_month,
		ug.created_year,
		ug.updated_day,
		ug.updated_week,
		ug.updated_month,
		ug.updated_year,
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
		v.updated_year from ugroups ug inner join ugroups_users ugu on (ug.id = ugu.ugroup_id) inner join users v on (ugu.user_id = v.id) where ug.uuid4 = ?`, uUID4byte)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2346}).Error(err)
			return nil, err
		}
		for rows.Next() {
			user := User{}
			err = rows.Scan(
				&poh.ID,
				&poh.UUID4,
				&poh.UgroupName,
				&poh.UgroupDesc,
				&poh.Levelc,
				&poh.ParentID,
				&poh.NumChd,
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
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2347}).Error(err)
			}
			uUID4Str1, err := common.UUIDBytesToStr(poh.UUID4)
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2348}).Error(err)
				return nil, err
			}
			poh.IDS = uUID4Str1

			uUID4Str, err := common.UUIDBytesToStr(user.UUID4)
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2349}).Error(err)
				return nil, err
			}
			user.IDS = uUID4Str

			poh.Users = append(poh.Users, &user)
		}

		err = rows.Close()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2350}).Error(err)
			return nil, err
		}

		err = rows.Err()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2351}).Error(err)
			return nil, err
		}

		return &poh, nil
	}
}

// GetUgroupByID - Get Ugroup By ID
func (u *UgroupService) GetUgroupByID(ctx context.Context, ID string, userEmail string, requestID string) (*Ugroup, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2352}).Error(err)
		return nil, err
	default:
		uUID4byte, err := common.UUIDStrToBytes(ID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2353}).Error(err)
			return nil, err
		}
		db := u.Db
		ug := Ugroup{}
		row := db.QueryRowContext(ctx, `select
    id,
		uuid4,
		ugroup_name,
		ugroup_desc,
		levelc,
		parent_id,
		num_chd,
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
		updated_year from ugroups where uuid4 = ?;`, uUID4byte)

		err = row.Scan(
			&ug.ID,
			&ug.UUID4,
			&ug.UgroupName,
			&ug.UgroupDesc,
			&ug.Levelc,
			&ug.ParentID,
			&ug.NumChd,
			&ug.Statusc,
			&ug.CreatedAt,
			&ug.UpdatedAt,
			&ug.CreatedDay,
			&ug.CreatedWeek,
			&ug.CreatedMonth,
			&ug.CreatedYear,
			&ug.UpdatedDay,
			&ug.UpdatedWeek,
			&ug.UpdatedMonth,
			&ug.UpdatedYear)

		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2354}).Error(err)
			return nil, err
		}
		uUID4Str, err := common.UUIDBytesToStr(ug.UUID4)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2355}).Error(err)
			return nil, err
		}
		ug.IDS = uUID4Str
		return &ug, nil
	}
}

// GetUgroupByIDuint - Get Ugroup By ID(uint)
func (u *UgroupService) GetUgroupByIDuint(ctx context.Context, ID uint, userEmail string, requestID string) (*Ugroup, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2356}).Error(err)
		return nil, err
	default:
		db := u.Db
		ug := Ugroup{}
		row := db.QueryRowContext(ctx, `select
    id,
		uuid4,
		ugroup_name,
		ugroup_desc,
		levelc,
		parent_id,
		num_chd,
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
		updated_year from ugroups where id = ?;`, ID)

		err := row.Scan(
			&ug.ID,
			&ug.UUID4,
			&ug.UgroupName,
			&ug.UgroupDesc,
			&ug.Levelc,
			&ug.ParentID,
			&ug.NumChd,
			&ug.Statusc,
			&ug.CreatedAt,
			&ug.UpdatedAt,
			&ug.CreatedDay,
			&ug.CreatedWeek,
			&ug.CreatedMonth,
			&ug.CreatedYear,
			&ug.UpdatedDay,
			&ug.UpdatedWeek,
			&ug.UpdatedMonth,
			&ug.UpdatedYear)

		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2357}).Error(err)
			return nil, err
		}
		uUID4Str, err := common.UUIDBytesToStr(ug.UUID4)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2358}).Error(err)
			return nil, err
		}
		ug.IDS = uUID4Str
		return &ug, nil
	}
}

// DeleteUserFromGroup - Delete user from group
func (u *UgroupService) DeleteUserFromGroup(ctx context.Context, form *UgroupUser, ID string, userEmail string, requestID string) error {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2359}).Error(err)
		return err
	default:
		db := u.Db
		tx, err := db.Begin()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2360}).Error(err)
			return err
		}
		stmt, err := tx.PrepareContext(ctx, `delete from ugroups_users where user_id= ? and ugroup_id = ?;`)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2361}).Error(err)
			err = stmt.Close()
			err = tx.Rollback()
			return err
		}

		_, err = stmt.ExecContext(ctx, form.UserID, ID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2362}).Error(err)
			err = stmt.Close()
			err = tx.Rollback()
			return err
		}
		err = stmt.Close()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2363}).Error(err)
			err = tx.Rollback()
			return err
		}

		err = tx.Commit()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2364}).Error(err)
			err = tx.Rollback()
			return err
		}
		return nil
	}
}

// InsertUgroupUser - Insert Ugroup User details into database
func (u *UgroupService) InsertUgroupUser(ctx context.Context, tx *sql.Tx, Uguser UgroupUser, userEmail string, requestID string) (*UgroupUser, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2365}).Error(err)
		return nil, err
	default:
		stmt, err := tx.PrepareContext(ctx, `insert into ugroups_users
	  (
		uuid4,
		ugroup_id,
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
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2366}).Error(err)
			err = stmt.Close()
			return nil, err
		}
		res, err := stmt.ExecContext(ctx,
			Uguser.UUID4,
			Uguser.UgroupID,
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
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2367}).Error(err)
			err = stmt.Close()
			return nil, err
		}
		uID, err := res.LastInsertId()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2368}).Error(err)
			err = stmt.Close()
			return nil, err
		}
		Uguser.ID = uint(uID)
		uUID4Str, err := common.UUIDBytesToStr(Uguser.UUID4)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2369}).Error(err)
			return nil, err
		}
		Uguser.IDS = uUID4Str
		err = stmt.Close()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2370}).Error(err)
			return nil, err
		}
		return &Uguser, nil
	}
}

// GetChildUgroups - Get child ugroups
func (u *UgroupService) GetChildUgroups(ctx context.Context, ID string, userEmail string, requestID string) ([]*Ugroup, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2371}).Error(err)
		return nil, err
	default:
		ugroup, err := u.GetUgroupByID(ctx, ID, userEmail, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2372}).Error(err)
			return nil, err
		}
		Ugroups := []*Ugroup{}
		rows, err := u.Db.QueryContext(ctx, `select 
    ug.id,
		ug.uuid4,
		ug.ugroup_name,
		ug.ugroup_desc,
		ug.levelc,
		ug.parent_id,
		ug.num_chd,
		ug.statusc,
    ug.created_at,
    ug.updated_at,
		ug.created_day,
		ug.created_week,
		ug.created_month,
		ug.created_year,
		ug.updated_day,
		ug.updated_week,
		ug.updated_month,
		ug.updated_year from ugroups ug inner join ugroup_chds ugch on (ug.id = ugch.ugroup_chd_id) where ((ugch.ugroup_id = ?) and (ug.statusc = 1))`, ugroup.ID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2373}).Error(err)
			return nil, err
		}

		for rows.Next() {
			ug := Ugroup{}
			err = rows.Scan(
				&ug.ID,
				&ug.UUID4,
				&ug.UgroupName,
				&ug.UgroupDesc,
				&ug.Levelc,
				&ug.ParentID,
				&ug.NumChd,
				&ug.Statusc,
				&ug.CreatedAt,
				&ug.UpdatedAt,
				&ug.CreatedDay,
				&ug.CreatedWeek,
				&ug.CreatedMonth,
				&ug.CreatedYear,
				&ug.UpdatedDay,
				&ug.UpdatedWeek,
				&ug.UpdatedMonth,
				&ug.UpdatedYear)
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2374}).Error(err)
				return nil, err
			}
			uUID4Str, err := common.UUIDBytesToStr(ug.UUID4)
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2375}).Error(err)
				return nil, err
			}
			ug.IDS = uUID4Str
			Ugroups = append(Ugroups, &ug)
		}
		err = rows.Close()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2376}).Error(err)
			return nil, err
		}
		err = rows.Err()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2377}).Error(err)
			return nil, err
		}

		return Ugroups, nil
	}

}

// TopLevelUgroups - Get top level ugroups
func (u *UgroupService) TopLevelUgroups(ctx context.Context, userEmail string, requestID string) ([]*Ugroup, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2378}).Error(err)
		return nil, err
	default:
		Ugroups := []*Ugroup{}
		rows, err := u.Db.QueryContext(ctx, `select 
    id,
		uuid4,
		ugroup_name,
		ugroup_desc,
		levelc,
		parent_id,
		num_chd,
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
		updated_year from ugroups where ((levelc = 0) and (statusc = 1))`)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2379}).Error(err)
			return nil, err
		}

		for rows.Next() {
			ug := Ugroup{}
			err = rows.Scan(
				&ug.ID,
				&ug.UUID4,
				&ug.UgroupName,
				&ug.UgroupDesc,
				&ug.Levelc,
				&ug.ParentID,
				&ug.NumChd,
				&ug.Statusc,
				&ug.CreatedAt,
				&ug.UpdatedAt,
				&ug.CreatedDay,
				&ug.CreatedWeek,
				&ug.CreatedMonth,
				&ug.CreatedYear,
				&ug.UpdatedDay,
				&ug.UpdatedWeek,
				&ug.UpdatedMonth,
				&ug.UpdatedYear)
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2380}).Error(err)
				return nil, err
			}
			uUID4Str, err := common.UUIDBytesToStr(ug.UUID4)
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2381}).Error(err)
				return nil, err
			}
			ug.IDS = uUID4Str
			Ugroups = append(Ugroups, &ug)
		}
		err = rows.Close()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2382}).Error(err)
			return nil, err
		}
		return Ugroups, nil
	}
}

// GetParent - Get parent ugroup
func (u *UgroupService) GetParent(ctx context.Context, ID string, userEmail string, requestID string) (*Ugroup, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2383}).Error(err)
		return nil, err
	default:
		db := u.Db
		ugroup, err := u.GetUgroupByID(ctx, ID, userEmail, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2384}).Error(err)
			return nil, err
		}
		ug := Ugroup{}
		row := db.QueryRowContext(ctx, `select
    id,
		uuid4,
		ugroup_name,
		ugroup_desc,
		levelc,
		parent_id,
		num_chd,
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
		updated_year from ugroups where id = ?;`, ugroup.ParentID)

		err = row.Scan(
			&ug.ID,
			&ug.UUID4,
			&ug.UgroupName,
			&ug.UgroupDesc,
			&ug.Levelc,
			&ug.ParentID,
			&ug.NumChd,
			&ug.Statusc,
			&ug.CreatedAt,
			&ug.UpdatedAt,
			&ug.CreatedDay,
			&ug.CreatedWeek,
			&ug.CreatedMonth,
			&ug.CreatedYear,
			&ug.UpdatedDay,
			&ug.UpdatedWeek,
			&ug.UpdatedMonth,
			&ug.UpdatedYear)

		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2385}).Error(err)
			return nil, err
		}
		uUID4Str, err := common.UUIDBytesToStr(ug.UUID4)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2386}).Error(err)
			return nil, err
		}
		ug.IDS = uUID4Str
		return &ug, nil
	}
}
