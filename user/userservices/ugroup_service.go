package userservices

import (
	"context"
	"database/sql"
	"errors"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/go-redis/redis"
	"github.com/palantir/stacktrace"

	"github.com/cloudfresco/vilom/common"
)

// Ugroup - Ugroup view representation
type Ugroup struct {
	ID uint

	IDS string

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
	ID  uint
	IDS string

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
func (u *UgroupService) GetUgroups(ctx context.Context, limit string, nextCursor string) (*UgroupCursor, error) {
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
			id_s,
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
		log.Println(err)
	}

	for rows.Next() {
		ug := Ugroup{}
		err = rows.Scan(&ug.ID,
			&ug.IDS,
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
			log.Error(stacktrace.Propagate(err, ""))
			return nil, err
		}
		ugroups = append(ugroups, &ug)
	}
	err = rows.Close()
	if err != nil {
		log.Error(stacktrace.Propagate(err, ""))
		return nil, err
	}

	err = rows.Err()
	if err != nil {
		log.Error(stacktrace.Propagate(err, ""))
		return nil, err
	}

	next := ugroups[len(ugroups)-1].ID
	next = next - 1
	nextc := common.EncodeCursor(next)
	x := UgroupCursor{ugroups, nextc}
	return &x, nil
}

// Create - Create ugroup
func (u *UgroupService) Create(ctx context.Context, form *Ugroup) (*Ugroup, error) {
	db := u.Db
	tx, err := db.Begin()
	if err != nil {
		log.Error(stacktrace.Propagate(err, ""))
		return nil, err
	}

	tn := time.Now().UTC()
	_, week := tn.ISOWeek()
	day := tn.YearDay()

	ug := Ugroup{}
	ug.IDS = common.GetUID()
	ug.UgroupName = form.UgroupName
	ug.UgroupDesc = form.UgroupDesc
	ug.Levelc = 0
	ug.ParentID = uint(0)
	ug.NumChd = 0
	ug.Statusc = common.Active
	ug.CreatedAt = tn
	ug.UpdatedAt = tn
	ug.CreatedDay = uint(day)
	ug.CreatedWeek = uint(week)
	ug.CreatedMonth = uint(tn.Month())
	ug.CreatedYear = uint(tn.Year())
	ug.UpdatedDay = uint(day)
	ug.UpdatedWeek = uint(week)
	ug.UpdatedMonth = uint(tn.Month())
	ug.UpdatedYear = uint(tn.Year())

	ugrp, err := u.InsertUgroup(ctx, tx, ug)

	if err != nil {
		log.Error(stacktrace.Propagate(err, ""))
		err = tx.Rollback()
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		log.Error(stacktrace.Propagate(err, ""))
		err = tx.Rollback()
		return nil, err
	}

	return ugrp, nil
}

// CreateChild - Create child of ugroup
func (u *UgroupService) CreateChild(ctx context.Context, form *Ugroup) (*Ugroup, error) {
	parent, err := u.GetUgroupByIDuint(ctx, form.ParentID)
	if err != nil {
		log.Error(stacktrace.Propagate(err, ""))
		return nil, err
	}

	db := u.Db
	tx, err := db.Begin()
	if err != nil {
		log.Error(stacktrace.Propagate(err, ""))
		return nil, err
	}

	tn := time.Now().UTC()
	_, week := tn.ISOWeek()
	day := tn.YearDay()
	ug := Ugroup{}
	ug.IDS = common.GetUID()
	ug.UgroupName = form.UgroupName
	ug.UgroupDesc = form.UgroupDesc
	ug.Levelc = parent.Levelc + 1
	ug.ParentID = parent.ID
	ug.NumChd = 0
	ug.Statusc = common.Active
	ug.CreatedAt = tn
	ug.UpdatedAt = tn
	ug.CreatedDay = uint(day)
	ug.CreatedWeek = uint(week)
	ug.CreatedMonth = uint(tn.Month())
	ug.CreatedYear = uint(tn.Year())
	ug.UpdatedDay = uint(day)
	ug.UpdatedWeek = uint(week)
	ug.UpdatedMonth = uint(tn.Month())
	ug.UpdatedYear = uint(tn.Year())

	ugrp, err := u.InsertUgroup(ctx, tx, ug)

	if err != nil {
		log.Error(stacktrace.Propagate(err, ""))
		err = tx.Rollback()
		return nil, err
	}

	Ugroupchd := UgroupChd{}
	Ugroupchd.UgroupID = parent.ID
	Ugroupchd.UgroupChdID = ugrp.ID
	Ugroupchd.Statusc = common.Active
	Ugroupchd.CreatedAt = tn
	Ugroupchd.UpdatedAt = tn
	Ugroupchd.CreatedDay = uint(day)
	Ugroupchd.CreatedWeek = uint(week)
	Ugroupchd.CreatedMonth = uint(tn.Month())
	Ugroupchd.CreatedYear = uint(tn.Year())
	Ugroupchd.UpdatedDay = uint(day)
	Ugroupchd.UpdatedWeek = uint(week)
	Ugroupchd.UpdatedMonth = uint(tn.Month())
	Ugroupchd.UpdatedYear = uint(tn.Year())

	_, err = u.InsertChild(ctx, tx, Ugroupchd)

	if err != nil {
		log.Error(stacktrace.Propagate(err, ""))
		err = tx.Rollback()
		return nil, err
	}

	UpdatedDay := uint(day)
	UpdatedWeek := uint(week)
	UpdatedMonth := uint(tn.Month())
	UpdatedYear := uint(tn.Year())

	stmt, err := tx.PrepareContext(ctx, `update ugroups set 
				  num_chd = ?,
				  updated_at = ?, 
					updated_day = ?, 
					updated_week = ?, 
					updated_month = ?, 
					updated_year = ? where id = ?;`)
	if err != nil {
		log.Error(stacktrace.Propagate(err, ""))
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
		log.Error(stacktrace.Propagate(err, ""))
		err = stmt.Close()
		err = tx.Rollback()
		return nil, err
	}

	err = stmt.Close()

	if err != nil {
		log.Error(stacktrace.Propagate(err, ""))
		err = tx.Rollback()
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		log.Error(stacktrace.Propagate(err, ""))
		err = tx.Rollback()
		return nil, err
	}
	return ugrp, nil
}

// AddUserToGroup - Add user to ugroup
func (u *UgroupService) AddUserToGroup(ctx context.Context, form *UgroupUser, ID string) error {
	db := u.Db
	ug, err := u.GetUgroupByID(ctx, ID)
	if err != nil {
		log.Error(stacktrace.Propagate(err, ""))
		return err
	}

	if ug.NumChd > 0 {
		err = errors.New("Cannot add user to group")
		log.Error(stacktrace.Propagate(err, ""))
		return err
	}

	tx, err := db.Begin()
	tn := time.Now().UTC()
	_, week := tn.ISOWeek()
	day := tn.YearDay()

	Uguser := UgroupUser{}
	Uguser.IDS = common.GetUID()
	Uguser.UgroupID = ug.ID
	Uguser.UserID = form.UserID
	Uguser.Statusc = common.Active
	Uguser.CreatedAt = tn
	Uguser.UpdatedAt = tn
	Uguser.CreatedDay = uint(day)
	Uguser.CreatedWeek = uint(week)
	Uguser.CreatedMonth = uint(tn.Month())
	Uguser.CreatedYear = uint(tn.Year())
	Uguser.UpdatedDay = uint(day)
	Uguser.UpdatedWeek = uint(week)
	Uguser.UpdatedMonth = uint(tn.Month())
	Uguser.UpdatedYear = uint(tn.Year())

	_, err = u.InsertUgroupUser(ctx, tx, Uguser)

	if err != nil {
		log.Error(stacktrace.Propagate(err, ""))
		err = tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		log.Error(stacktrace.Propagate(err, ""))
		err = tx.Rollback()
		return err
	}
	return nil
}

// InsertUgroup - Insert Ugroup details into database
func (u *UgroupService) InsertUgroup(ctx context.Context, tx *sql.Tx, ug Ugroup) (*Ugroup, error) {
	stmt, err := tx.PrepareContext(ctx, `insert into ugroups
	  (
		id_s,
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
		log.Error(stacktrace.Propagate(err, ""))
		return nil, err
	}
	res, err := stmt.ExecContext(ctx,
		ug.IDS,
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
		log.Error(stacktrace.Propagate(err, ""))
		err = stmt.Close()
		return nil, err
	}
	uID, err := res.LastInsertId()
	if err != nil {
		log.Error(stacktrace.Propagate(err, ""))
		return nil, err
	}
	ug.ID = uint(uID)
	err = stmt.Close()
	if err != nil {
		log.Error(stacktrace.Propagate(err, ""))
		return nil, err
	}
	return &ug, nil
}

// InsertChild - Insert Child Ugroup details into database
func (u *UgroupService) InsertChild(ctx context.Context, tx *sql.Tx, ugroupchd UgroupChd) (*UgroupChd, error) {
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
		log.Error(stacktrace.Propagate(err, ""))
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
		log.Error(stacktrace.Propagate(err, ""))
		err = stmt.Close()
		return nil, err
	}
	uID, err := res.LastInsertId()
	if err != nil {
		log.Error(stacktrace.Propagate(err, ""))
		err = stmt.Close()
		return nil, err
	}
	ugroupchd.ID = uint(uID)
	err = stmt.Close()
	if err != nil {
		log.Error(stacktrace.Propagate(err, ""))
		return nil, err
	}
	return &ugroupchd, nil
}

// Delete - Delete ugroup
func (u *UgroupService) Delete(ctx context.Context, ID string) error {
	db := u.Db
	tx, err := db.Begin()
	stmt, err := tx.PrepareContext(ctx, "delete from ugroups where id_s= ?;")
	if err != nil {
		log.Error(stacktrace.Propagate(err, ""))
		err = tx.Rollback()
		return err
	}

	_, err = stmt.ExecContext(ctx, ID)
	if err != nil {
		log.Error(stacktrace.Propagate(err, ""))
		err = stmt.Close()
		err = tx.Rollback()
		return err
	}

	err = stmt.Close()
	if err != nil {
		log.Error(stacktrace.Propagate(err, ""))
		err = tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		log.Error(stacktrace.Propagate(err, ""))
		err = tx.Rollback()
		return err
	}
	return nil
}

// GetUgroup - Get ugroup details with users by ID
func (u *UgroupService) GetUgroup(ctx context.Context, ID string) (*Ugroup, error) {
	db := u.Db
	poh := Ugroup{}

	rows, err := db.QueryContext(ctx, `select 
    ug.id,
		ug.id_s,
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
		v.id_s,
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
		v.updated_year from ugroups ug inner join ugroups_users ugu on (ug.id = ugu.ugroup_id) inner join users v on (ugu.user_id = v.id) where ug.id_s = ?`, ID)
	if err != nil {
		log.Error(stacktrace.Propagate(err, ""))
		return nil, err
	}
	for rows.Next() {
		user := User{}
		err = rows.Scan(
			&poh.ID,
			&poh.IDS,
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
			&user.IDS,
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
			log.Error(stacktrace.Propagate(err, ""))
		}
		poh.Users = append(poh.Users, &user)
	}

	err = rows.Close()
	if err != nil {
		log.Error(stacktrace.Propagate(err, ""))
		return nil, err
	}

	err = rows.Err()
	if err != nil {
		log.Error(stacktrace.Propagate(err, ""))
		return nil, err
	}

	return &poh, nil
}

// GetUgroupByID - Get Ugroup By ID
func (u *UgroupService) GetUgroupByID(ctx context.Context, ID string) (*Ugroup, error) {
	db := u.Db
	ug := Ugroup{}
	row := db.QueryRowContext(ctx, `select
    id,
		id_s,
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
		updated_year from ugroups where id_s = ?;`, ID)

	err := row.Scan(
		&ug.ID,
		&ug.IDS,
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
		log.Error(stacktrace.Propagate(err, ""))
		return nil, err
	}

	return &ug, nil
}

// GetUgroupByIDuint - Get Ugroup By ID(uint)
func (u *UgroupService) GetUgroupByIDuint(ctx context.Context, ID uint) (*Ugroup, error) {
	db := u.Db
	ug := Ugroup{}
	row := db.QueryRowContext(ctx, `select
    id,
		id_s,
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
		&ug.IDS,
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
		log.Error(stacktrace.Propagate(err, ""))
		return nil, err
	}

	return &ug, nil
}

// DeleteUserFromGroup - Delete user from group
func (u *UgroupService) DeleteUserFromGroup(ctx context.Context, form *UgroupUser, ID string) error {
	db := u.Db
	tx, err := db.Begin()
	if err != nil {
		log.Error(stacktrace.Propagate(err, ""))
		return err
	}
	stmt, err := tx.PrepareContext(ctx, `delete from ugroups_users where user_id= ? and ugroup_id = ?;`)
	if err != nil {
		log.Error(stacktrace.Propagate(err, ""))
		err = stmt.Close()
		err = tx.Rollback()
		return err
	}

	_, err = stmt.ExecContext(ctx, form.UserID, ID)
	if err != nil {
		log.Error(stacktrace.Propagate(err, ""))
		err = stmt.Close()
		err = tx.Rollback()
		return err
	}
	err = stmt.Close()
	if err != nil {
		log.Error(stacktrace.Propagate(err, ""))
		err = tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		log.Error(stacktrace.Propagate(err, ""))
		err = tx.Rollback()
		return err
	}
	return nil
}

// InsertUgroupUser - Insert Ugroup User details into database
func (u *UgroupService) InsertUgroupUser(ctx context.Context, tx *sql.Tx, Uguser UgroupUser) (*UgroupUser, error) {
	stmt, err := tx.PrepareContext(ctx, `insert into ugroups_users
	  (
		id_s,
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
		log.Error(stacktrace.Propagate(err, ""))
		err = stmt.Close()
		return nil, err
	}
	res, err := stmt.ExecContext(ctx,
		Uguser.IDS,
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
		log.Error(stacktrace.Propagate(err, ""))
		err = stmt.Close()
		return nil, err
	}
	uID, err := res.LastInsertId()
	if err != nil {
		log.Error(stacktrace.Propagate(err, ""))
		err = stmt.Close()
		return nil, err
	}
	Uguser.ID = uint(uID)
	err = stmt.Close()
	if err != nil {
		log.Error(stacktrace.Propagate(err, ""))
		return nil, err
	}
	return &Uguser, nil
}

// GetChildUgroups - Get child ugroups
func (u *UgroupService) GetChildUgroups(ctx context.Context, ID string) ([]*Ugroup, error) {
	ugroup, err := u.GetUgroupByID(ctx, ID)
	if err != nil {
		log.Error(stacktrace.Propagate(err, ""))
		return nil, err
	}
	Ugroups := []*Ugroup{}
	rows, err := u.Db.QueryContext(ctx, `select 
    ug.id,
		ug.id_s,
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
		log.Error(stacktrace.Propagate(err, ""))
		return nil, err
	}

	for rows.Next() {
		ug := Ugroup{}
		err = rows.Scan(
			&ug.ID,
			&ug.IDS,
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
			log.Error(stacktrace.Propagate(err, ""))
			return nil, err
		}
		Ugroups = append(Ugroups, &ug)
	}
	err = rows.Close()
	if err != nil {
		log.Error(stacktrace.Propagate(err, ""))
		return nil, err
	}
	err = rows.Err()
	if err != nil {
		log.Error(stacktrace.Propagate(err, ""))
		return nil, err
	}

	return Ugroups, nil

}

// TopLevelUgroups - Get top level ugroups
func (u *UgroupService) TopLevelUgroups(ctx context.Context) ([]*Ugroup, error) {
	Ugroups := []*Ugroup{}
	rows, err := u.Db.QueryContext(ctx, `select 
    id,
		id_s,
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
		log.Error(stacktrace.Propagate(err, ""))
		return nil, err
	}

	for rows.Next() {
		ug := Ugroup{}
		err = rows.Scan(
			&ug.ID,
			&ug.IDS,
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
			log.Error(stacktrace.Propagate(err, ""))
			return nil, err
		}
		Ugroups = append(Ugroups, &ug)
	}
	err = rows.Close()
	if err != nil {
		log.Error(stacktrace.Propagate(err, ""))
		return nil, err
	}
	return Ugroups, nil

}

// GetParent - Get parent ugroup
func (u *UgroupService) GetParent(ctx context.Context, ID string) (*Ugroup, error) {
	db := u.Db
	ugroup, err := u.GetUgroupByID(ctx, ID)
	if err != nil {
		log.Error(stacktrace.Propagate(err, ""))
		return nil, err
	}
	ug := Ugroup{}
	row := db.QueryRowContext(ctx, `select
    id,
		id_s,
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
		&ug.IDS,
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
		log.Error(stacktrace.Propagate(err, ""))
		return nil, err
	}

	return &ug, nil
}
