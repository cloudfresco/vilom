package userservices

import (
	"context"
	"database/sql"
	"errors"

	log "github.com/sirupsen/logrus"

	"github.com/cloudfresco/vilom/common"
)

/* error message range: 2300-2999 */

// For validation of Ugroup fields
const (
	UgroupNameLenMin = 1
	UgroupNameLenMax = 50
	UgroupDescLenMin = 1
	UgroupDescLenMax = 1000
)

// Ugroup - Ugroup view representation
type Ugroup struct {
	ID    uint   `json:"id,omitempty"`
	UUID4 []byte `json:"-"`
	IDS   string `json:"id_s,omitempty"`

	UgroupName string `json:"ugroup_name,omitempty"`
	UgroupDesc string `json:"ugroup_desc,omitempty"`
	Levelc     uint   `json:"levelc,omitempty"`
	ParentID   uint   `json:"parent_id,omitempty"`
	NumChd     uint   `json:"num_chd,omitempty"`

	common.StatusDates

	Users []*User
}

// UgroupChd - UgroupChd view representation
type UgroupChd struct {
	ID          uint   `json:"id,omitempty"`
	UUID4       []byte `json:"-"`
	UgroupID    uint   `json:"ugroup_id,omitempty"`
	UgroupChdID uint   `json:"ugroup_chd_id,omitempty"`

	common.StatusDates
}

// UgroupUser - UgroupUser view representation
type UgroupUser struct {
	ID    uint   `json:"id,omitempty"`
	UUID4 []byte `json:"-"`
	IDS   string `json:"id_s,omitempty"`

	UgroupID uint `json:"ugroup_id,omitempty"`
	UserID   uint `json:"user_id,omitempty"`

	common.StatusDates
}

// UgroupServiceIntf - interface for Ugroup Service
type UgroupServiceIntf interface {
	CreateUgroup(ctx context.Context, form *Ugroup, userEmail string, requestID string) (*Ugroup, error)
	CreateChild(ctx context.Context, form *Ugroup, userEmail string, requestID string) (*Ugroup, error)
	AddUserToGroup(ctx context.Context, form *UgroupUser, ID string, userEmail string, requestID string) error
	GetUgroups(ctx context.Context, limit string, nextCursor string, userEmail string, requestID string) (*UgroupCursor, error)
	TopLevelUgroups(ctx context.Context, userEmail string, requestID string) ([]*Ugroup, error)
	GetParent(ctx context.Context, ID string, userEmail string, requestID string) (*Ugroup, error)
	GetUgroup(ctx context.Context, ID string, userEmail string, requestID string) (*Ugroup, error)
	GetUgroupByID(ctx context.Context, ID string, userEmail string, requestID string) (*Ugroup, error)
	GetUgroupByIDuint(ctx context.Context, ID uint, userEmail string, requestID string) (*Ugroup, error)
	GetChildUgroups(ctx context.Context, ID string, userEmail string, requestID string) ([]*Ugroup, error)
	UpdateUgroup(ctx context.Context, ID string, form *Ugroup, UserID string, userEmail string, requestID string) error
	DeleteUgroup(ctx context.Context, ID string, userEmail string, requestID string) error
	DeleteUserFromGroup(ctx context.Context, form *UgroupUser, ID string, userEmail string, requestID string) error
}

// UgroupService - For accessing Ugroup services
type UgroupService struct {
	DBService    *common.DBService
	RedisService *common.RedisService
}

// NewUgroupService - Create Ugroup Service
func NewUgroupService(dbOpt *common.DBService, redisOpt *common.RedisService) *UgroupService {
	return &UgroupService{
		DBService:    dbOpt,
		RedisService: redisOpt,
	}
}

// UgroupCursor - used to get groups
type UgroupCursor struct {
	Ugroups    []*Ugroup
	NextCursor string `json:"next_cursor,omitempty"`
}

// CreateUgroup - Create ugroup
func (u *UgroupService) CreateUgroup(ctx context.Context, form *Ugroup, userEmail string, requestID string) (*Ugroup, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2306}).Error(err)
		return nil, err
	default:
		db := u.DBService.DB
		insertUgroupStmt, err := u.insertUgroupPrepare(ctx, userEmail, requestID)
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

		tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2307}).Error(err)
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

		err = u.insertUgroup(ctx, insertUgroupStmt, tx, &ug, userEmail, requestID)

		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2309}).Error(err)
			err = tx.Rollback()
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2307}).Error(err)
				return nil, err
			}
			err = insertUgroupStmt.Close()
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2307}).Error(err)
				return nil, err
			}
			return nil, err
		}

		err = tx.Commit()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2310}).Error(err)
			return nil, err
		}

		err = insertUgroupStmt.Close()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2310}).Error(err)
			return nil, err
		}

		return &ug, nil
	}
}

// insertUgroupPrepare - Insert Ugroup Prepare Statements
func (u *UgroupService) insertUgroupPrepare(ctx context.Context, userEmail string, requestID string) (*sql.Stmt, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2328}).Error(err)
		return nil, err
	default:
		db := u.DBService.DB
		stmt, err := db.PrepareContext(ctx, `insert into ugroups
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
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2329}).Error(err)
			return nil, err
		}
		return stmt, nil
	}
}

// insertUgroup - Insert Ugroup details into database
func (u *UgroupService) insertUgroup(ctx context.Context, stmt *sql.Stmt, tx *sql.Tx, ug *Ugroup, userEmail string, requestID string) error {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2328}).Error(err)
		return err
	default:
		res, err := tx.StmtContext(ctx, stmt).Exec(
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
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2330}).Error(err)
			return err
		}
		uID, err := res.LastInsertId()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2331}).Error(err)
			return err
		}
		ug.ID = uint(uID)
		uuid4Str, err := common.UUIDBytesToStr(ug.UUID4)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2332}).Error(err)
			return err
		}
		ug.IDS = uuid4Str
		return nil
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

		db := u.DBService.DB
		insertUgroupStmt, insertChildStmt, updateNumChildrenStmt, err := u.createChildPrepareStmts(ctx, userEmail, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2313}).Error(err)
			return nil, err
		}

		tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2313}).Error(err)
			return nil, err
		}

		ug, err := u.createChild(ctx, insertUgroupStmt, insertChildStmt, updateNumChildrenStmt, tx, form, userEmail, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2313}).Error(err)
			err = tx.Rollback()
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2313}).Error(err)
				return nil, err
			}
			err = u.createChildPrepareStmtsClose(ctx, insertUgroupStmt, insertChildStmt, updateNumChildrenStmt, userEmail, requestID)
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2313}).Error(err)
				return nil, err
			}
			return nil, err
		}

		err = tx.Commit()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2321}).Error(err)
			return nil, err
		}

		err = u.createChildPrepareStmtsClose(ctx, insertUgroupStmt, insertChildStmt, updateNumChildrenStmt, userEmail, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2320}).Error(err)
			return nil, err
		}

		return ug, nil
	}
}

//createChildPrepareStmts - Prepare Statements
func (u *UgroupService) createChildPrepareStmts(ctx context.Context, userEmail string, requestID string) (*sql.Stmt, *sql.Stmt, *sql.Stmt, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2312}).Error(err)
		return nil, nil, nil, err
	default:
		insertUgroupStmt, err := u.insertUgroupPrepare(ctx, userEmail, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2312}).Error(err)
			return nil, nil, nil, err
		}
		insertChildStmt, err := u.insertChildPrepare(ctx, userEmail, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2312}).Error(err)
			return nil, nil, nil, err
		}
		updateNumChildrenStmt, err := u.updateNumChildrenPrepare(ctx, userEmail, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2312}).Error(err)
			return nil, nil, nil, err
		}
		return insertUgroupStmt, insertChildStmt, updateNumChildrenStmt, nil
	}
}

//createChildPrepareStmtsClose - Prepare Statements Close
func (u *UgroupService) createChildPrepareStmtsClose(ctx context.Context, insertUgroupStmt *sql.Stmt, insertChildStmt *sql.Stmt, updateNumChildrenStmt *sql.Stmt, userEmail string, requestID string) error {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2312}).Error(err)
		return err
	default:
		err := insertUgroupStmt.Close()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2312}).Error(err)
			return err
		}
		err = insertChildStmt.Close()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2312}).Error(err)
			return err
		}
		err = updateNumChildrenStmt.Close()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2312}).Error(err)
			return err
		}
		return nil
	}
}

// createChild - Create Child ugroup
func (u *UgroupService) createChild(ctx context.Context, insertUgroupStmt *sql.Stmt, insertChildStmt *sql.Stmt, updateNumChildrenStmt *sql.Stmt, tx *sql.Tx, form *Ugroup, userEmail string, requestID string) (*Ugroup, error) {
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

		err = u.insertUgroup(ctx, insertUgroupStmt, tx, &ug, userEmail, requestID)

		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2315}).Error(err)
			return nil, err
		}

		Ugroupchd := UgroupChd{}
		Ugroupchd.UUID4, err = common.GetUUIDBytes()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2316}).Error(err)
			return nil, err
		}
		Ugroupchd.UgroupID = parent.ID
		Ugroupchd.UgroupChdID = ug.ID
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

		err = u.insertChild(ctx, insertChildStmt, tx, &Ugroupchd, userEmail, requestID)

		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2317}).Error(err)
			return nil, err
		}

		numchd := parent.NumChd + 1
		err = u.updateNumChildren(ctx, updateNumChildrenStmt, tx, numchd, parent.ID, userEmail, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2317}).Error(err)
			return nil, err
		}
		return &ug, err
	}
}

// insertChildPrepare - Insert Child Ugroup Prepare Statement
func (u *UgroupService) insertChildPrepare(ctx context.Context, userEmail string, requestID string) (*sql.Stmt, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2334}).Error(err)
		return nil, err
	default:
		db := u.DBService.DB
		stmt, err := db.PrepareContext(ctx, `insert into ugroup_chds
	  (
    uuid4, 
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
					?,?,?,?);`)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2335}).Error(err)
			return nil, err
		}
		return stmt, nil
	}
}

// insertChild - Insert Child Ugroup details into database
func (u *UgroupService) insertChild(ctx context.Context, stmt *sql.Stmt, tx *sql.Tx, ugroupchd *UgroupChd, userEmail string, requestID string) error {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2334}).Error(err)
		return err
	default:
		res, err := tx.StmtContext(ctx, stmt).Exec(
			ugroupchd.UUID4,
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
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2336}).Error(err)
			return err
		}
		uID, err := res.LastInsertId()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2337}).Error(err)
			return err
		}
		ugroupchd.ID = uint(uID)
		return nil
	}
}

// updateNumChildrenPrepare - updateNumChildren Prepare Statement
func (u *UgroupService) updateNumChildrenPrepare(ctx context.Context, userEmail string, requestID string) (*sql.Stmt, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2317}).Error(err)
		return nil, err
	default:
		db := u.DBService.DB
		stmt, err := db.PrepareContext(ctx, `update ugroups set 
				  num_chd = ?,
				  updated_at = ?, 
					updated_day = ?, 
					updated_week = ?, 
					updated_month = ?, 
					updated_year = ? where id = ? and statusc = ?;`)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2317}).Error(err)
			return nil, err
		}
		return stmt, nil
	}
}

// updateNumChildren - Update number of child in ugroup
func (u *UgroupService) updateNumChildren(ctx context.Context, stmt *sql.Stmt, tx *sql.Tx, numchd uint, parentID uint, userEmail string, requestID string) error {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2317}).Error(err)
		return err
	default:
		tn, tnday, tnweek, tnmonth, tnyear := common.GetTimeDetails()
		_, err := tx.StmtContext(ctx, stmt).Exec(
			numchd,
			tn,
			tnday,
			tnweek,
			tnmonth,
			tnyear,
			parentID,
			common.Active)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2317}).Error(err)
			return err
		}
		return nil
	}
}

// AddUserToGroup - Add user to ugroup
func (u *UgroupService) AddUserToGroup(ctx context.Context, form *UgroupUser, ID string, userEmail string, requestID string) error {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2322}).Error(err)
		return err
	default:
		db := u.DBService.DB
		ug, err := u.GetUgroupByID(ctx, ID, userEmail, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2323}).Error(err)
			return err
		}

		if ug.NumChd > 0 {
			err = errors.New("Cannot add user to group")
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2324}).Error(err)
			return err
		}

		insertUgroupUserStmt, err := u.insertUgroupUserPrepare(ctx, userEmail, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2390}).Error(err)
			return err
		}

		tn, tnday, tnweek, tnmonth, tnyear := common.GetTimeDetails()

		Uguser := UgroupUser{}
		Uguser.UUID4, err = common.GetUUIDBytes()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2325}).Error(err)
			return err
		}

		tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2390}).Error(err)
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

		err = u.insertUgroupUser(ctx, insertUgroupUserStmt, tx, &Uguser, userEmail, requestID)

		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2326}).Error(err)
			err = tx.Rollback()
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2326}).Error(err)
				return err
			}
			err = insertUgroupUserStmt.Close()
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2326}).Error(err)
				return err
			}
			return err
		}

		err = tx.Commit()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2327}).Error(err)
			return err
		}

		err = insertUgroupUserStmt.Close()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2327}).Error(err)
			return err
		}
		return nil
	}
}

// insertUgroupUserPrepare - Insert Ugroup User Prepare Statement
func (u *UgroupService) insertUgroupUserPrepare(ctx context.Context, userEmail string, requestID string) (*sql.Stmt, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2366}).Error(err)
		return nil, err
	default:
		db := u.DBService.DB
		stmt, err := db.PrepareContext(ctx, `insert into ugroups_users
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
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2367}).Error(err)
			return nil, err
		}
		return stmt, nil
	}
}

// insertUgroupUser - Insert Ugroup User details into database
func (u *UgroupService) insertUgroupUser(ctx context.Context, stmt *sql.Stmt, tx *sql.Tx, Uguser *UgroupUser, userEmail string, requestID string) error {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2366}).Error(err)
		return err
	default:
		res, err := tx.StmtContext(ctx, stmt).Exec(
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
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2368}).Error(err)
			return err
		}
		uID, err := res.LastInsertId()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2369}).Error(err)
			return err
		}
		Uguser.ID = uint(uID)
		uuid4Str, err := common.UUIDBytesToStr(Uguser.UUID4)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2370}).Error(err)
			return err
		}
		Uguser.IDS = uuid4Str
		return nil
	}
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
			limit = u.DBService.LimitSQLRows
		}
		query := "(statusc = ?)"
		if nextCursor == "" {
			query = query + " order by id desc " + " limit " + limit + ";"
		} else {
			nextCursor = common.DecodeCursor(nextCursor)
			query = query + " " + "and" + " " + "id <= " + nextCursor + " order by id desc " + " limit " + limit + ";"
		}

		ugroups := []*Ugroup{}
		db := u.DBService.DB
		rows, err := db.QueryContext(ctx, `select 
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
			updated_year from ugroups where `+query, common.Active)
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
			uuid4Str, err := common.UUIDBytesToStr(ug.UUID4)
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2303}).Error(err)
				return nil, err
			}
			ug.IDS = uuid4Str
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

// TopLevelUgroups - Get top level ugroups
func (u *UgroupService) TopLevelUgroups(ctx context.Context, userEmail string, requestID string) ([]*Ugroup, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2379}).Error(err)
		return nil, err
	default:
		Ugroups := []*Ugroup{}
		db := u.DBService.DB
		rows, err := db.QueryContext(ctx, `select 
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
		updated_year from ugroups where levelc = ? and statusc = ?;`, 0, common.Active)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2380}).Error(err)
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
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2381}).Error(err)
				return nil, err
			}
			uuid4Str, err := common.UUIDBytesToStr(ug.UUID4)
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2382}).Error(err)
				return nil, err
			}
			ug.IDS = uuid4Str
			Ugroups = append(Ugroups, &ug)
		}
		err = rows.Close()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2383}).Error(err)
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
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2384}).Error(err)
		return nil, err
	default:
		db := u.DBService.DB
		ugroup, err := u.GetUgroupByID(ctx, ID, userEmail, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2385}).Error(err)
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
		updated_year from ugroups where id = ? and statusc = ?;`, ugroup.ParentID, common.Active)

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
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2386}).Error(err)
			return nil, err
		}
		uuid4Str, err := common.UUIDBytesToStr(ug.UUID4)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2387}).Error(err)
			return nil, err
		}
		ug.IDS = uuid4Str
		return &ug, nil
	}
}

// GetUgroup - Get ugroup details with users by ID
func (u *UgroupService) GetUgroup(ctx context.Context, ID string, userEmail string, requestID string) (*Ugroup, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2345}).Error(err)
		return nil, err
	default:
		db := u.DBService.DB
		ugroup, err := u.GetUgroupByID(ctx, ID, userEmail, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2403}).Error(err)
			return nil, err
		}

		var isPresent bool
		row := db.QueryRowContext(ctx, `select exists (select 1 from ugroups_users where ugroup_id = ?);`, ugroup.ID)
		err = row.Scan(&isPresent)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2404}).Error(err)
			return nil, err
		}
		if !isPresent {
			return ugroup, nil
		}

		uuid4byte, err := common.UUIDStrToBytes(ID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2346}).Error(err)
			return nil, err
		}

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
		v.updated_year from ugroups ug inner join ugroups_users ugu on (ug.id = ugu.ugroup_id) inner join users v on (ugu.user_id = v.id) where ug.uuid4 = ?`, uuid4byte)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2347}).Error(err)
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
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2348}).Error(err)
			}
			uuid4Str1, err := common.UUIDBytesToStr(poh.UUID4)
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2349}).Error(err)
				return nil, err
			}
			poh.IDS = uuid4Str1

			uuid4Str, err := common.UUIDBytesToStr(user.UUID4)
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2350}).Error(err)
				return nil, err
			}
			user.IDS = uuid4Str

			poh.Users = append(poh.Users, &user)
		}

		err = rows.Close()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2351}).Error(err)
			return nil, err
		}

		err = rows.Err()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2352}).Error(err)
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
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2353}).Error(err)
		return nil, err
	default:
		uuid4byte, err := common.UUIDStrToBytes(ID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2354}).Error(err)
			return nil, err
		}
		db := u.DBService.DB
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
		updated_year from ugroups where uuid4 = ? and statusc = ?;`, uuid4byte, common.Active)

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
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2355}).Error(err)
			return nil, err
		}
		uuid4Str, err := common.UUIDBytesToStr(ug.UUID4)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2356}).Error(err)
			return nil, err
		}
		ug.IDS = uuid4Str
		return &ug, nil
	}
}

// GetUgroupByIDuint - Get Ugroup By ID(uint)
func (u *UgroupService) GetUgroupByIDuint(ctx context.Context, ID uint, userEmail string, requestID string) (*Ugroup, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2357}).Error(err)
		return nil, err
	default:
		db := u.DBService.DB
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
		updated_year from ugroups where id = ? and statusc = ?;`, ID, common.Active)

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
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2358}).Error(err)
			return nil, err
		}
		uuid4Str, err := common.UUIDBytesToStr(ug.UUID4)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2359}).Error(err)
			return nil, err
		}
		ug.IDS = uuid4Str
		return &ug, nil
	}
}

// GetChildUgroups - Get child ugroups
func (u *UgroupService) GetChildUgroups(ctx context.Context, ID string, userEmail string, requestID string) ([]*Ugroup, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2372}).Error(err)
		return nil, err
	default:
		ugroup, err := u.GetUgroupByID(ctx, ID, userEmail, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2373}).Error(err)
			return nil, err
		}
		Ugroups := []*Ugroup{}
		db := u.DBService.DB
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
		ug.updated_year from ugroups ug inner join ugroup_chds ugch on (ug.id = ugch.ugroup_chd_id) where ((ugch.ugroup_id = ?) and (ug.statusc = 1))`, ugroup.ID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2374}).Error(err)
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
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2375}).Error(err)
				return nil, err
			}
			uuid4Str, err := common.UUIDBytesToStr(ug.UUID4)
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2376}).Error(err)
				return nil, err
			}
			ug.IDS = uuid4Str
			Ugroups = append(Ugroups, &ug)
		}
		err = rows.Close()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2377}).Error(err)
			return nil, err
		}
		err = rows.Err()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2378}).Error(err)
			return nil, err
		}

		return Ugroups, nil
	}

}

//UpdateUgroup - Update Ugroup
func (u *UgroupService) UpdateUgroup(ctx context.Context, ID string, form *Ugroup, UserID string, userEmail string, requestID string) error {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2394}).Error(err)
		return err
	default:
		ugroup, err := u.GetUgroup(ctx, ID, userEmail, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2395}).Error(err)
			return err
		}

		db := u.DBService.DB
		stmt, err := db.PrepareContext(ctx, `update ugroups set 
		  ugroup_name = ?,
      ugroup_desc = ?,
			updated_at = ?, 
			updated_day = ?, 
			updated_week = ?, 
			updated_month = ?, 
			updated_year = ? where id = ? and statusc = ?;`)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2397}).Error(err)
			return err
		}

		tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2396}).Error(err)
			return err
		}

		tn, tnday, tnweek, tnmonth, tnyear := common.GetTimeDetails()
		_, err = tx.StmtContext(ctx, stmt).Exec(
			form.UgroupName,
			form.UgroupDesc,
			tn,
			tnday,
			tnweek,
			tnmonth,
			tnyear,
			ugroup.ID,
			common.Active)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2399}).Error(err)
			err = tx.Rollback()
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2400}).Error(err)
				return err
			}
			err = stmt.Close()
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2400}).Error(err)
				return err
			}
			err = tx.Rollback()
			return err
		}

		err = tx.Commit()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2402}).Error(err)
			return err
		}

		err = stmt.Close()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2401}).Error(err)
			return err
		}

		return nil
	}
}

// DeleteUgroup - Delete ugroup
func (u *UgroupService) DeleteUgroup(ctx context.Context, ID string, userEmail string, requestID string) error {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2339}).Error(err)
		return err
	default:
		uuid4byte, err := common.UUIDStrToBytes(ID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2340}).Error(err)
			return err
		}
		db := u.DBService.DB
		stmt, err := db.PrepareContext(ctx, `update ugroups set 
		  statusc = ?,
			updated_at = ?, 
			updated_day = ?, 
			updated_week = ?, 
			updated_month = ?, 
			updated_year = ? where uuid4= ?;`)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2341}).Error(err)
			return err
		}

		tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2391}).Error(err)
			return err
		}
		tn, tnday, tnweek, tnmonth, tnyear := common.GetTimeDetails()

		_, err = tx.StmtContext(ctx, stmt).Exec(
			common.Inactive,
			tn,
			tnday,
			tnweek,
			tnmonth,
			tnyear,
			uuid4byte)

		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2342}).Error(err)
			err = tx.Rollback()
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2392}).Error(err)
				return err
			}
			err = stmt.Close()
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2392}).Error(err)
				return err
			}

			return err
		}

		err = tx.Commit()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2344}).Error(err)
			err = tx.Rollback()
			return err
		}

		err = stmt.Close()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2343}).Error(err)
			return err
		}

		return nil
	}
}

// DeleteUserFromGroup - Delete user from group
func (u *UgroupService) DeleteUserFromGroup(ctx context.Context, form *UgroupUser, ID string, userEmail string, requestID string) error {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2360}).Error(err)
		return err
	default:
		db := u.DBService.DB
		stmt, err := db.PrepareContext(ctx, `delete from ugroups_users where user_id= ? and ugroup_id = ?;`)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2362}).Error(err)
			return err
		}
		tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2361}).Error(err)
			return err
		}
		_, err = tx.StmtContext(ctx, stmt).Exec(form.UserID, ID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2363}).Error(err)
			err = tx.Rollback()
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2393}).Error(err)
				return err
			}
			err = stmt.Close()
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2393}).Error(err)
				return err
			}

			return err
		}

		err = tx.Commit()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2365}).Error(err)
			err = tx.Rollback()
			return err
		}

		err = stmt.Close()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 2364}).Error(err)
			return err
		}

		return nil
	}
}
