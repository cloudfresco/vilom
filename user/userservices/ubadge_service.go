package userservices

import (
	"context"
	"database/sql"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/go-redis/redis"
	"github.com/palantir/stacktrace"

	"github.com/cloudfresco/vilom/common"
)

// Ubadge - Ubadge view representation
type Ubadge struct {
	ID  uint
	IDS string

	UbadgeName string
	UbadgeDesc string

	common.StatusDates
	Users []*User
}

// UbadgeUser - Ubadge User view representation
type UbadgeUser struct {
	ID  uint
	IDS string

	UbadgeID uint
	UserID   uint

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
	NextCursor string
}

// GetUbadges - Get Ubadges
func (u *UbadgeService) GetUbadges(ctx context.Context, limit string, nextCursor string) (*UbadgeCursor, error) {
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
			id_s,
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
		log.Error(stacktrace.Propagate(err, ""))
		return nil, err
	}

	for rows.Next() {
		ubadge := Ubadge{}
		err = rows.Scan(&ubadge.ID,
			&ubadge.IDS,
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
			log.Error(stacktrace.Propagate(err, ""))
			return nil, err
		}
		ubadges = append(ubadges, &ubadge)
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

	next := ubadges[len(ubadges)-1].ID
	next = next - 1
	nextc := common.EncodeCursor(next)
	x := UbadgeCursor{ubadges, nextc}
	return &x, nil
}

// Create - Create Ubadge
func (u *UbadgeService) Create(ctx context.Context, form *Ubadge) (*Ubadge, error) {
	db := u.Db
	tx, err := db.Begin()
	if err != nil {
		log.Error(stacktrace.Propagate(err, ""))
		return nil, err
	}

	tn := time.Now().UTC()
	_, week := tn.ISOWeek()
	day := tn.YearDay()

	Ubadge := Ubadge{}
	Ubadge.IDS = common.GetUID()
	Ubadge.UbadgeName = form.UbadgeName
	Ubadge.UbadgeDesc = form.UbadgeDesc
	Ubadge.Statusc = common.Active
	Ubadge.CreatedAt = tn
	Ubadge.UpdatedAt = tn
	Ubadge.CreatedDay = uint(day)
	Ubadge.CreatedWeek = uint(week)
	Ubadge.CreatedMonth = uint(tn.Month())
	Ubadge.CreatedYear = uint(tn.Year())
	Ubadge.UpdatedDay = uint(day)
	Ubadge.UpdatedWeek = uint(week)
	Ubadge.UpdatedMonth = uint(tn.Month())
	Ubadge.UpdatedYear = uint(tn.Year())

	ugrp, err := u.InsertUbadge(ctx, tx, Ubadge)

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

// AddUserToGroup - Add user to ubadge
func (u *UbadgeService) AddUserToGroup(ctx context.Context, form *UbadgeUser, ID string) error {
	db := u.Db
	ubadge, err := u.GetUbadge(ctx, ID)
	if err != nil {
		log.Error(stacktrace.Propagate(err, ""))
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		log.Error(stacktrace.Propagate(err, ""))
		return err
	}
	tn := time.Now().UTC()
	_, week := tn.ISOWeek()
	day := tn.YearDay()

	Uguser := UbadgeUser{}
	Uguser.IDS = common.GetUID()
	Uguser.UbadgeID = ubadge.ID
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

	_, err = u.InsertUbadgeUser(ctx, tx, Uguser)

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

// InsertUbadge - Insert Ubadge details into database
func (u *UbadgeService) InsertUbadge(ctx context.Context, tx *sql.Tx, Ubadge Ubadge) (*Ubadge, error) {
	stmt, err := tx.PrepareContext(ctx, `insert into ubadges
	  (
		id_s,
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
		log.Error(stacktrace.Propagate(err, ""))
		err = stmt.Close()
		return nil, err
	}
	res, err := stmt.ExecContext(ctx,
		Ubadge.IDS,
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
	Ubadge.ID = uint(uID)
	err = stmt.Close()
	if err != nil {
		log.Error(stacktrace.Propagate(err, ""))
		return nil, err
	}

	return &Ubadge, nil
}

// Delete - Delele Ubadge
func (u *UbadgeService) Delete(ctx context.Context, ID string) error {
	db := u.Db
	tx, err := db.Begin()
	if err != nil {
		log.Error(stacktrace.Propagate(err, ""))
		return err
	}
	stmt, err := tx.PrepareContext(ctx, "delete from ubadges where id_s= ?;")
	if err != nil {
		log.Error(stacktrace.Propagate(err, ""))
		err = stmt.Close()
		err = tx.Rollback()
		return err
	}

	_, err = stmt.ExecContext(ctx, ID)
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

// GetUbadge - Get Ubadge Details
func (u *UbadgeService) GetUbadge(ctx context.Context, ID string) (*Ubadge, error) {
	db := u.Db
	poh := Ubadge{}
	rows, err := db.QueryContext(ctx, `select 
    p.id,
		p.id_s,
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
		v.updated_year from ubadges p inner join ubadges_users ubu on (p.id = ubu.ubadge_id) inner join users v on (ubu.user_id = v.id) where p.id_s = ?`, ID)

	if err != nil {
		log.Error(stacktrace.Propagate(err, ""))
		return nil, err
	}
	for rows.Next() {
		user := User{}
		err = rows.Scan(
			&poh.ID,
			&poh.IDS,
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
			return nil, err
		}
		poh.Users = append(poh.Users, &user)
	}

	err = rows.Close()
	if err != nil {
		log.Error(stacktrace.Propagate(err, ""))
		return nil, err
	}

	if err != nil {
		log.Error(stacktrace.Propagate(err, ""))
		return nil, err
	}

	return &poh, nil
}

// GetUbadgeByID - Get Ubadge by ID
func (u *UbadgeService) GetUbadgeByID(ctx context.Context, ID string) (*Ubadge, error) {
	db := u.Db
	Ubadge := Ubadge{}
	row := db.QueryRowContext(ctx, `select
    id,
		id_s,
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
		updated_year from ubadges where id_s = ?;`, ID)

	err := row.Scan(
		&Ubadge.ID,
		&Ubadge.IDS,
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
		log.Error(stacktrace.Propagate(err, ""))
		return nil, err
	}

	return &Ubadge, nil
}

// DeleteUserFromGroup - Delete user from Ubadge
func (u *UbadgeService) DeleteUserFromGroup(ctx context.Context, form *UbadgeUser, ID string) error {
	db := u.Db
	tx, err := db.Begin()
	if err != nil {
		log.Error(stacktrace.Propagate(err, ""))
		return err
	}
	stmt, err := tx.PrepareContext(ctx, `delete from ubadges_users where user_id= ? and ubadge_id = ?;`)
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

// InsertUbadgeUser - Insert Ubadge User details into database
func (u *UbadgeService) InsertUbadgeUser(ctx context.Context, tx *sql.Tx, Uguser UbadgeUser) (*UbadgeUser, error) {
	stmt, err := tx.PrepareContext(ctx, `insert into ubadges_users
	  (
		id_s,
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
		log.Error(stacktrace.Propagate(err, ""))
		return nil, err
	}
	res, err := stmt.ExecContext(ctx,
		Uguser.IDS,
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
