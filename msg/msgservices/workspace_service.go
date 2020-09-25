package msgservices

import (
	"context"
	"database/sql"
	"errors"

	log "github.com/sirupsen/logrus"

	"github.com/cloudfresco/vilom/common"
	"github.com/cloudfresco/vilom/user/userservices"
)

/* error message range: 4300-4999 */

// For validation of workspace fields
const (
	WorkspaceNameLenMin = 1
	WorkspaceNameLenMax = 50
	WorkspaceDescLenMin = 1
	WorkspaceDescLenMax = 1000
)

// Workspace - Workspace view representation
type Workspace struct {
	ID            uint   `json:"id,omitempty"`
	UUID4         []byte `json:"-"`
	IDS           string `json:"id_s,omitempty"`
	WorkspaceName string `json:"workspace_name,omitempty"`
	WorkspaceDesc string `json:"workspace_desc,omitempty"`
	NumViews      uint   `json:"num_views,omitempty"`
	NumChannels   uint   `json:"num_channels,omitempty"`
	Levelc        uint   `json:"levelc,omitempty"`
	ParentID      uint   `json:"parent_id,omitempty"`
	NumChd        uint   `json:"num_chd,omitempty"`

	UgroupID uint `json:"ugroup_id,omitempty"`
	UserID   uint `json:"user_id,omitempty"`

	common.StatusDates
	Channels []*Channel
}

// WorkspaceChd - WorkspaceChd view representation
type WorkspaceChd struct {
	ID             uint   `json:"id,omitempty"`
	UUID4          []byte `json:"-"`
	WorkspaceID    uint   `json:"workspace_id,omitempty"`
	WorkspaceChdID uint   `json:"workspace_chd_id,omitempty"`

	common.StatusDates
}

// WorkspaceServiceIntf - interface for Workspace Service
type WorkspaceServiceIntf interface {
	CreateWorkspace(ctx context.Context, form *Workspace, UserID string, userEmail string, requestID string) (*Workspace, error)
	CreateChild(ctx context.Context, form *Workspace, UserID string, userEmail string, requestID string) (*Workspace, error)
	GetWorkspaces(ctx context.Context, limit string, nextCursor string, userEmail string, requestID string) (*WorkspaceCursor, error)
	GetWorkspaceWithChannels(ctx context.Context, ID string, userEmail string, requestID string) (*Workspace, error)
	GetWorkspace(ctx context.Context, ID string, userEmail string, requestID string) (*Workspace, error)
	GetWorkspaceByID(ctx context.Context, ID uint, userEmail string, requestID string) (*Workspace, error)
	GetTopLevelWorkspaces(ctx context.Context, userEmail string, requestID string) ([]*Workspace, error)
	GetChildWorkspaces(ctx context.Context, ID string, userEmail string, requestID string) ([]*Workspace, error)
	GetParentWorkspace(ctx context.Context, ID string, userEmail string, requestID string) (*Workspace, error)
	UpdateWorkspace(ctx context.Context, ID string, form *Workspace, UserID string, userEmail string, requestID string) error
	DeleteWorkspace(ctx context.Context, ID string, userEmail string, requestID string) error
}

// WorkspaceService - For accessing workspace services
type WorkspaceService struct {
	DBService    *common.DBService
	RedisService *common.RedisService
}

// NewWorkspaceService - Create workspace service
func NewWorkspaceService(dbOpt *common.DBService, redisOpt *common.RedisService) *WorkspaceService {
	return &WorkspaceService{
		DBService:    dbOpt,
		RedisService: redisOpt,
	}
}

// WorkspaceCursor - used to get workspaces
type WorkspaceCursor struct {
	Workspaces []*Workspace
	NextCursor string `json:"next_cursor,omitempty"`
}

// CreateWorkspace - Create Workspace
func (c *WorkspaceService) CreateWorkspace(ctx context.Context, form *Workspace, UserID string, userEmail string, requestID string) (*Workspace, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4318}).Error(err)
		return nil, err
	default:
		userserv := &userservices.UserService{DBService: c.DBService, RedisService: c.RedisService}
		user, err := userserv.GetUser(ctx, UserID, userEmail, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4319}).Error(err)
			return nil, err
		}
		db := c.DBService.DB
		insertWorkspaceStmt, err := c.insertWorkspacePrepare(ctx, userEmail, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4369}).Error(err)
			return nil, err
		}
		tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4368}).Error(err)
			return nil, err
		}
		tn, tnday, tnweek, tnmonth, tnyear := common.GetTimeDetails()
		workspace := Workspace{}
		workspace.UUID4, err = common.GetUUIDBytes()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4320}).Error(err)
			return nil, err
		}
		workspace.WorkspaceName = form.WorkspaceName
		workspace.WorkspaceDesc = form.WorkspaceDesc
		workspace.NumViews = 0
		workspace.NumChannels = 0
		workspace.Levelc = 0
		workspace.ParentID = uint(0)
		workspace.NumChd = 0
		workspace.UgroupID = uint(0)
		workspace.UserID = user.ID
		/*  StatusDates  */
		workspace.Statusc = common.Active
		workspace.CreatedAt = tn
		workspace.UpdatedAt = tn
		workspace.CreatedDay = tnday
		workspace.CreatedWeek = tnweek
		workspace.CreatedMonth = tnmonth
		workspace.CreatedYear = tnyear
		workspace.UpdatedDay = tnday
		workspace.UpdatedWeek = tnweek
		workspace.UpdatedMonth = tnmonth
		workspace.UpdatedYear = tnyear

		err = c.insertWorkspace(ctx, insertWorkspaceStmt, tx, &workspace, userEmail, requestID)

		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4321}).Error(err)
			err = tx.Rollback()
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4328}).Error(err)
				return nil, err
			}
			err = insertWorkspaceStmt.Close()
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4328}).Error(err)
				return nil, err
			}
			return nil, err
		}

		err = tx.Commit()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4322}).Error(err)
			err = tx.Rollback()
			return nil, err
		}

		err = insertWorkspaceStmt.Close()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4328}).Error(err)
			return nil, err
		}

		return &workspace, nil
	}
}

// insertWorkspacePrepare - Insert Workspace prepare statement
func (c *WorkspaceService) insertWorkspacePrepare(ctx context.Context, userEmail string, requestID string) (*sql.Stmt, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4323}).Error(err)
		return nil, err
	default:
		db := c.DBService.DB
		stmt, err := db.PrepareContext(ctx, `insert into workspaces
	  ( 
			uuid4,
			workspace_name,
			workspace_desc,
			num_views,
			num_channels,
			levelc,
			parent_id,
			num_chd,
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
					?,?,?,?,?,?,?,?,?,?,
          ?);`)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4324}).Error(err)
			return nil, err
		}
		return stmt, nil
	}
}

// insertWorkspace - Insert workspace details into database
func (c *WorkspaceService) insertWorkspace(ctx context.Context, stmt *sql.Stmt, tx *sql.Tx, workspace *Workspace, userEmail string, requestID string) error {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4323}).Error(err)
		return err
	default:
		res, err := tx.StmtContext(ctx, stmt).Exec(
			workspace.UUID4,
			workspace.WorkspaceName,
			workspace.WorkspaceDesc,
			workspace.NumViews,
			workspace.NumChannels,
			workspace.Levelc,
			workspace.ParentID,
			workspace.NumChd,
			workspace.UgroupID,
			workspace.UserID,
			/*  StatusDates  */
			workspace.Statusc,
			workspace.CreatedAt,
			workspace.UpdatedAt,
			workspace.CreatedDay,
			workspace.CreatedWeek,
			workspace.CreatedMonth,
			workspace.CreatedYear,
			workspace.UpdatedDay,
			workspace.UpdatedWeek,
			workspace.UpdatedMonth,
			workspace.UpdatedYear)

		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4325}).Error(err)
			return err
		}

		uID, err := res.LastInsertId()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4326}).Error(err)
			return err
		}
		workspace.ID = uint(uID)
		uuid4Str, err := common.UUIDBytesToStr(workspace.UUID4)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4327}).Error(err)
			return err
		}
		workspace.IDS = uuid4Str
		return nil
	}
}

// CreateChild - Create Child Workspace
func (c *WorkspaceService) CreateChild(ctx context.Context, form *Workspace, UserID string, userEmail string, requestID string) (*Workspace, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4338}).Error(err)
		return nil, err
	default:
		db := c.DBService.DB
		insertWorkspaceStmt, insertChildStmt, updateNumChildrenStmt, err := c.createChildPrepareStmts(ctx, userEmail, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4340}).Error(err)
			return nil, err
		}

		tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4369}).Error(err)
			err = c.createChildPrepareStmtsClose(ctx, insertWorkspaceStmt, insertChildStmt, updateNumChildrenStmt, userEmail, requestID)
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4340}).Error(err)
				return nil, err
			}
			return nil, err
		}

		workspace, err := c.createChild(ctx, insertWorkspaceStmt, insertChildStmt, updateNumChildrenStmt, tx, form, UserID, userEmail, requestID)

		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4342}).Error(err)
			err = tx.Rollback()
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4340}).Error(err)
				return nil, err
			}
			err = c.createChildPrepareStmtsClose(ctx, insertWorkspaceStmt, insertChildStmt, updateNumChildrenStmt, userEmail, requestID)
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4340}).Error(err)
				return nil, err
			}
			return nil, err
		}

		err = tx.Commit()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4372}).Error(err)
			return nil, err
		}

		err = c.createChildPrepareStmtsClose(ctx, insertWorkspaceStmt, insertChildStmt, updateNumChildrenStmt, userEmail, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4340}).Error(err)
			return nil, err
		}
		return workspace, nil
	}
}

//createChildPrepareStmts - Prepare Statements
func (c *WorkspaceService) createChildPrepareStmts(ctx context.Context, userEmail string, requestID string) (*sql.Stmt, *sql.Stmt, *sql.Stmt, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6311}).Error(err)
		return nil, nil, nil, err
	default:
		insertWorkspaceStmt, err := c.insertWorkspacePrepare(ctx, userEmail, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4369}).Error(err)
			return nil, nil, nil, err
		}
		insertChildStmt, err := c.insertChildPrepare(ctx, userEmail, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4369}).Error(err)
			return nil, nil, nil, err
		}
		updateNumChildrenStmt, err := c.updateNumChildrenPrepare(ctx, userEmail, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4369}).Error(err)
			return nil, nil, nil, err
		}
		return insertWorkspaceStmt, insertChildStmt, updateNumChildrenStmt, nil
	}
}

//createChildPrepareStmtsClose - Close Prepare Statements
func (c *WorkspaceService) createChildPrepareStmtsClose(ctx context.Context, insertWorkspaceStmt *sql.Stmt, insertChildStmt *sql.Stmt, updateNumChildrenStmt *sql.Stmt, userEmail string, requestID string) error {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6311}).Error(err)
		return err
	default:
		err := insertWorkspaceStmt.Close()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4328}).Error(err)
			return err
		}
		err = insertChildStmt.Close()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4328}).Error(err)
			return err
		}
		err = updateNumChildrenStmt.Close()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4328}).Error(err)
			return err
		}

		return nil
	}
}

// createChild - Create Child Workspace
func (c *WorkspaceService) createChild(ctx context.Context, insertWorkspaceStmt *sql.Stmt, insertChildStmt *sql.Stmt, updateNumChildrenStmt *sql.Stmt, tx *sql.Tx, form *Workspace, UserID string, userEmail string, requestID string) (*Workspace, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4338}).Error(err)
		return nil, err
	default:
		userserv := &userservices.UserService{DBService: c.DBService, RedisService: c.RedisService}
		user, err := userserv.GetUser(ctx, UserID, userEmail, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4339}).Error(err)
			return nil, err
		}

		parent, err := c.GetWorkspaceByID(ctx, form.ParentID, userEmail, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4340}).Error(err)
			return nil, err
		}

		tn, tnday, tnweek, tnmonth, tnyear := common.GetTimeDetails()
		workspace := Workspace{}
		workspace.UUID4, err = common.GetUUIDBytes()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4341}).Error(err)
			return nil, err
		}
		workspace.WorkspaceName = form.WorkspaceName
		workspace.WorkspaceDesc = form.WorkspaceDesc
		workspace.NumViews = 0
		workspace.NumChannels = 0
		workspace.Levelc = parent.Levelc + 1
		workspace.ParentID = parent.ID
		workspace.NumChd = 0
		workspace.UgroupID = uint(0)
		workspace.UserID = user.ID
		/*  StatusDates  */
		workspace.Statusc = common.Active
		workspace.CreatedAt = tn
		workspace.UpdatedAt = tn
		workspace.CreatedDay = tnday
		workspace.CreatedWeek = tnweek
		workspace.CreatedMonth = tnmonth
		workspace.CreatedYear = tnyear
		workspace.UpdatedDay = tnday
		workspace.UpdatedWeek = tnweek
		workspace.UpdatedMonth = tnmonth
		workspace.UpdatedYear = tnyear

		err = c.insertWorkspace(ctx, insertWorkspaceStmt, tx, &workspace, userEmail, requestID)

		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4342}).Error(err)
			return nil, err
		}

		catChd := WorkspaceChd{}
		catChd.UUID4, err = common.GetUUIDBytes()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4343}).Error(err)
			return nil, err
		}
		catChd.WorkspaceID = parent.ID
		catChd.WorkspaceChdID = workspace.ID
		/*  StatusDates  */
		catChd.Statusc = common.Active
		catChd.CreatedAt = tn
		catChd.UpdatedAt = tn
		catChd.CreatedDay = tnday
		catChd.CreatedWeek = tnweek
		catChd.CreatedMonth = tnmonth
		catChd.CreatedYear = tnyear
		catChd.UpdatedDay = tnday
		catChd.UpdatedWeek = tnweek
		catChd.UpdatedMonth = tnmonth
		catChd.UpdatedYear = tnyear

		err = c.insertChild(ctx, insertChildStmt, tx, &catChd, userEmail, requestID)

		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4344}).Error(err)
			return nil, err
		}

		numchd := parent.NumChd + 1
		err = c.updateNumChildren(ctx, updateNumChildrenStmt, tx, numchd, parent.ID, userEmail, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4373}).Error(err)
			return nil, err
		}

		return &workspace, nil
	}
}

// insertChildPrepare - Insert child prepare statement
func (c *WorkspaceService) insertChildPrepare(ctx context.Context, userEmail string, requestID string) (*sql.Stmt, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4348}).Error(err)
		return nil, err
	default:
		db := c.DBService.DB
		stmt, err := db.PrepareContext(ctx, `insert into workspace_chds
	  ( 
    uuid4,
		workspace_id,
		workspace_chd_id,
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
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4349}).Error(err)
			return nil, err
		}
		return stmt, nil
	}
}

// insertChild - Insert child workspace details into database
func (c *WorkspaceService) insertChild(ctx context.Context, stmt *sql.Stmt, tx *sql.Tx, catChd *WorkspaceChd, userEmail string, requestID string) error {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4348}).Error(err)
		return err
	default:
		res, err := tx.StmtContext(ctx, stmt).Exec(
			catChd.UUID4,
			catChd.WorkspaceID,
			catChd.WorkspaceChdID,
			/*  StatusDates  */
			catChd.Statusc,
			catChd.CreatedAt,
			catChd.UpdatedAt,
			catChd.CreatedDay,
			catChd.CreatedWeek,
			catChd.CreatedMonth,
			catChd.CreatedYear,
			catChd.UpdatedDay,
			catChd.UpdatedWeek,
			catChd.UpdatedMonth,
			catChd.UpdatedYear)

		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4350}).Error(err)
			return err
		}

		uID, err := res.LastInsertId()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4351}).Error(err)
			return err
		}
		catChd.ID = uint(uID)
		return nil
	}
}

// updateNumChildrenPrepare - updateNumChildren Prepare Statement
func (c *WorkspaceService) updateNumChildrenPrepare(ctx context.Context, userEmail string, requestID string) (*sql.Stmt, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4375}).Error(err)
		return nil, err
	default:
		db := c.DBService.DB
		stmt, err := db.PrepareContext(ctx, `update workspaces set 
				  num_chd = ?,
				  updated_at = ?, 
					updated_day = ?, 
					updated_week = ?, 
					updated_month = ?, 
					updated_year = ? where id = ? and statusc = ?;`)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4345}).Error(err)
			return nil, err
		}
		return stmt, nil
	}
}

// updateNumChildren - Update number of child in workspace
func (c *WorkspaceService) updateNumChildren(ctx context.Context, stmt *sql.Stmt, tx *sql.Tx, numchd uint, parentID uint, userEmail string, requestID string) error {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4375}).Error(err)
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
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4346}).Error(err)
			return err
		}
		return nil
	}
}

// GetWorkspaces - Get Workspaces
func (c *WorkspaceService) GetWorkspaces(ctx context.Context, limit string, nextCursor string, userEmail string, requestID string) (*WorkspaceCursor, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4300}).Error(err)
		return nil, err
	default:
		if limit == "" {
			limit = c.DBService.LimitSQLRows
		}
		query := "(levelc = ? and statusc = ?)"
		if nextCursor == "" {
			query = query + " order by id desc " + " limit " + limit + ";"
		} else {
			nextCursor = common.DecodeCursor(nextCursor)
			query = query + " " + "and" + " " + "id <= " + nextCursor + " order by id desc " + " limit " + limit + ";"
		}
		db := c.DBService.DB
		workspaces := []*Workspace{}
		rows, err := db.QueryContext(ctx, `select 
      id, 
			uuid4,
			workspace_name,
			workspace_desc,
			num_views,
			num_channels,
			levelc,
			parent_id,
			num_chd,
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
			updated_year from workspaces where `+query, 0, common.Active)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4301}).Error(err)
			return nil, err
		}

		for rows.Next() {
			workspace := Workspace{}
			err = rows.Scan(
				&workspace.ID,
				&workspace.UUID4,
				&workspace.WorkspaceName,
				&workspace.WorkspaceDesc,
				&workspace.NumViews,
				&workspace.NumChannels,
				&workspace.Levelc,
				&workspace.ParentID,
				&workspace.NumChd,
				&workspace.UgroupID,
				&workspace.UserID,
				/*  StatusDates  */
				&workspace.Statusc,
				&workspace.CreatedAt,
				&workspace.UpdatedAt,
				&workspace.CreatedDay,
				&workspace.CreatedWeek,
				&workspace.CreatedMonth,
				&workspace.CreatedYear,
				&workspace.UpdatedDay,
				&workspace.UpdatedWeek,
				&workspace.UpdatedMonth,
				&workspace.UpdatedYear)
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4302}).Error(err)
				return nil, err
			}
			uuid4Str, err := common.UUIDBytesToStr(workspace.UUID4)
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4303}).Error(err)
				return nil, err
			}
			workspace.IDS = uuid4Str
			workspaces = append(workspaces, &workspace)
		}
		err = rows.Close()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4304}).Error(err)
			return nil, err
		}

		err = rows.Err()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4305}).Error(err)
			return nil, err
		}
		x := WorkspaceCursor{}
		if len(workspaces) != 0 {
			next := workspaces[len(workspaces)-1].ID
			next = next - 1
			nextc := common.EncodeCursor(next)
			x = WorkspaceCursor{workspaces, nextc}
		} else {
			x = WorkspaceCursor{workspaces, "0"}
		}
		return &x, nil
	}
}

// GetWorkspaceWithChannels - Get workspace with channels
func (c *WorkspaceService) GetWorkspaceWithChannels(ctx context.Context, ID string, userEmail string, requestID string) (*Workspace, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4329}).Error(err)
		return nil, err
	default:
		uuid4byte, err := common.UUIDStrToBytes(ID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4330}).Error(err)
			return nil, err
		}
		db := c.DBService.DB
		workspace := Workspace{}
		ctegry, err := c.GetWorkspace(ctx, ID, userEmail, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4331}).Error(err)
			return nil, err
		}
		var isPresent bool
		row := db.QueryRowContext(ctx, `select exists (select 1 from channels where workspace_id = ?);`, ctegry.ID)
		err = row.Scan(&isPresent)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4332}).Error(err)
			return nil, err
		}
		if isPresent {

			rows, err := db.QueryContext(ctx, `select 
		  c.id,
			c.uuid4,
			c.workspace_name,
			c.workspace_desc,
			c.num_views,
			c.num_channels,
			c.levelc,
			c.parent_id,
			c.num_chd,
			c.ugroup_id,
			c.user_id,
			c.statusc,
			c.created_at,
			c.updated_at,
			c.created_day,
			c.created_week,
			c.created_month,
			c.created_year,
			c.updated_day,
			c.updated_week,
			c.updated_month,
			c.updated_year,
		  v.id,
			v.uuid4,
			v.channel_name,
			v.channel_desc,
			v.num_tags,
			v.tag1,
			v.tag2,
			v.tag3,
			v.tag4,
			v.tag5,
			v.tag6,
			v.tag7,
			v.tag8,
			v.tag9,
			v.tag10,
			v.num_views,
			v.num_messages,
			v.workspace_id,
			v.user_id,
			v.ugroup_id,
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
			v.updated_year from workspaces c inner join channels v on (c.id = v.workspace_id) where c.uuid4 = ?`, uuid4byte)

			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4333}).Error(err)
				return nil, err
			}
			for rows.Next() {
				topc := Channel{}
				err = rows.Scan(
					&workspace.ID,
					&workspace.UUID4,
					&workspace.WorkspaceName,
					&workspace.WorkspaceDesc,
					&workspace.NumViews,
					&workspace.NumChannels,
					&workspace.Levelc,
					&workspace.ParentID,
					&workspace.NumChd,
					&workspace.UgroupID,
					&workspace.UserID,
					/*  StatusDates  */
					&workspace.Statusc,
					&workspace.CreatedAt,
					&workspace.UpdatedAt,
					&workspace.CreatedDay,
					&workspace.CreatedWeek,
					&workspace.CreatedMonth,
					&workspace.CreatedYear,
					&workspace.UpdatedDay,
					&workspace.UpdatedWeek,
					&workspace.UpdatedMonth,
					&workspace.UpdatedYear,
					&topc.ID,
					&topc.UUID4,
					&topc.ChannelName,
					&topc.ChannelDesc,
					&topc.NumTags,
					&topc.Tag1,
					&topc.Tag2,
					&topc.Tag3,
					&topc.Tag4,
					&topc.Tag5,
					&topc.Tag6,
					&topc.Tag7,
					&topc.Tag8,
					&topc.Tag9,
					&topc.Tag10,
					&topc.NumViews,
					&topc.NumMessages,
					&topc.WorkspaceID,
					&topc.UserID,
					&topc.UgroupID,
					/*  StatusDates  */
					&topc.Statusc,
					&topc.CreatedAt,
					&topc.UpdatedAt,
					&topc.CreatedDay,
					&topc.CreatedWeek,
					&topc.CreatedMonth,
					&topc.CreatedYear,
					&topc.UpdatedDay,
					&topc.UpdatedWeek,
					&topc.UpdatedMonth,
					&topc.UpdatedYear)

				if err != nil {
					log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4334}).Error(err)
					return nil, err
				}
				uuid4Str1, err := common.UUIDBytesToStr(workspace.UUID4)
				if err != nil {
					log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4335}).Error(err)
					return nil, err
				}
				workspace.IDS = uuid4Str1

				uuid4Str, err := common.UUIDBytesToStr(topc.UUID4)
				if err != nil {
					log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4336}).Error(err)
					return nil, err
				}
				topc.IDS = uuid4Str
				workspace.Channels = append(workspace.Channels, &topc)
			}

			err = rows.Close()
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4337}).Error(err)
				return nil, err
			}
		} else {
			return ctegry, nil
		}
		return &workspace, nil
	}
}

// GetWorkspace - Get Workspace
func (c *WorkspaceService) GetWorkspace(ctx context.Context, ID string, userEmail string, requestID string) (*Workspace, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4306}).Error(err)
		return nil, err
	default:
		uuid4byte, err := common.UUIDStrToBytes(ID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4307}).Error(err)
			return nil, err
		}
		workspace := Workspace{}
		db := c.DBService.DB
		row := db.QueryRowContext(ctx, `select
      id,
			uuid4,
			workspace_name,
			workspace_desc,
			num_views,
			num_channels,
			levelc,
			parent_id,
			num_chd,
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
			updated_year from workspaces where uuid4 = ? and statusc = ?;`, uuid4byte, common.Active)

		err = row.Scan(
			&workspace.ID,
			&workspace.UUID4,
			&workspace.WorkspaceName,
			&workspace.WorkspaceDesc,
			&workspace.NumViews,
			&workspace.NumChannels,
			&workspace.Levelc,
			&workspace.ParentID,
			&workspace.NumChd,
			&workspace.UgroupID,
			&workspace.UserID,
			&workspace.Statusc,
			&workspace.CreatedAt,
			&workspace.UpdatedAt,
			&workspace.CreatedDay,
			&workspace.CreatedWeek,
			&workspace.CreatedMonth,
			&workspace.CreatedYear,
			&workspace.UpdatedDay,
			&workspace.UpdatedWeek,
			&workspace.UpdatedMonth,
			&workspace.UpdatedYear)

		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4308}).Error(err)
			return nil, err
		}
		uuid4Str, err := common.UUIDBytesToStr(workspace.UUID4)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4309}).Error(err)
			return nil, err
		}
		workspace.IDS = uuid4Str

		return &workspace, nil
	}
}

// GetWorkspaceByID - Get Workspace By ID
func (c *WorkspaceService) GetWorkspaceByID(ctx context.Context, ID uint, userEmail string, requestID string) (*Workspace, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4310}).Error(err)
		return nil, err
	default:
		workspace := Workspace{}
		db := c.DBService.DB
		row := db.QueryRowContext(ctx, `select
      id,
			uuid4,
			workspace_name,
			workspace_desc,
			num_views,
			num_channels,
			levelc,
			parent_id,
			num_chd,
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
			updated_year from workspaces where id = ? and statusc = ?;`, ID, common.Active)

		err := row.Scan(
			&workspace.ID,
			&workspace.UUID4,
			&workspace.WorkspaceName,
			&workspace.WorkspaceDesc,
			&workspace.NumViews,
			&workspace.NumChannels,
			&workspace.Levelc,
			&workspace.ParentID,
			&workspace.NumChd,
			&workspace.UgroupID,
			&workspace.UserID,
			/*  StatusDates  */
			&workspace.Statusc,
			&workspace.CreatedAt,
			&workspace.UpdatedAt,
			&workspace.CreatedDay,
			&workspace.CreatedWeek,
			&workspace.CreatedMonth,
			&workspace.CreatedYear,
			&workspace.UpdatedDay,
			&workspace.UpdatedWeek,
			&workspace.UpdatedMonth,
			&workspace.UpdatedYear)

		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4311}).Error(err)
			return nil, err
		}
		uuid4Str, err := common.UUIDBytesToStr(workspace.UUID4)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4312}).Error(err)
			return nil, err
		}
		workspace.IDS = uuid4Str
		return &workspace, nil
	}
}

// GetTopLevelWorkspaces - Get top level workspaces
func (c *WorkspaceService) GetTopLevelWorkspaces(ctx context.Context, userEmail string, requestID string) ([]*Workspace, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4353}).Error(err)
		return nil, err
	default:
		workspaces := []*Workspace{}
		db := c.DBService.DB
		rows, err := db.QueryContext(ctx, `select 
      id, 
			uuid4,
			workspace_name,
			workspace_desc,
			num_views,
			num_channels,
			levelc,
			parent_id,
			num_chd,
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
			updated_year from workspaces where levelc = ? and statusc = ?;`, 0, common.Active)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4354}).Error(err)
			return nil, err
		}

		for rows.Next() {
			workspace := Workspace{}
			err = rows.Scan(
				&workspace.ID,
				&workspace.UUID4,
				&workspace.WorkspaceName,
				&workspace.WorkspaceDesc,
				&workspace.NumViews,
				&workspace.NumChannels,
				&workspace.Levelc,
				&workspace.ParentID,
				&workspace.NumChd,
				&workspace.UgroupID,
				&workspace.UserID,
				/*  StatusDates  */
				&workspace.Statusc,
				&workspace.CreatedAt,
				&workspace.UpdatedAt,
				&workspace.CreatedDay,
				&workspace.CreatedWeek,
				&workspace.CreatedMonth,
				&workspace.CreatedYear,
				&workspace.UpdatedDay,
				&workspace.UpdatedWeek,
				&workspace.UpdatedMonth,
				&workspace.UpdatedYear)
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4355}).Error(err)
				return nil, err
			}
			uuid4Str, err := common.UUIDBytesToStr(workspace.UUID4)
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4356}).Error(err)
				return nil, err
			}
			workspace.IDS = uuid4Str
			workspaces = append(workspaces, &workspace)
		}
		err = rows.Close()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4357}).Error(err)
			return nil, err
		}
		return workspaces, nil
	}
}

// GetChildWorkspaces - Get child workspaces
func (c *WorkspaceService) GetChildWorkspaces(ctx context.Context, ID string, userEmail string, requestID string) ([]*Workspace, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4358}).Error(err)
		return nil, err
	default:
		workspace, err := c.GetWorkspace(ctx, ID, userEmail, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4359}).Error(err)
			return nil, err
		}
		pohs := []*Workspace{}
		db := c.DBService.DB
		rows, err := db.QueryContext(ctx, `select 
		    c.id,
				c.uuid4,
				c.workspace_name,
				c.workspace_desc,
				c.num_views,
				c.num_channels,
				c.levelc,
				c.parent_id,
				c.num_chd,
				c.ugroup_id,
				c.user_id,
				c.statusc,
				c.created_at,
				c.updated_at,
				c.created_day,
				c.created_week,
				c.created_month,
				c.created_year,
				c.updated_day,
				c.updated_week,
				c.updated_month,
				c.updated_year from workspaces c inner join workspace_chds ch on (c.id = ch.workspace_chd_id) where ch.workspace_id = ?`, workspace.ID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4360}).Error(err)
			return nil, err
		}

		for rows.Next() {
			workspace := Workspace{}
			err = rows.Scan(
				&workspace.ID,
				&workspace.UUID4,
				&workspace.WorkspaceName,
				&workspace.WorkspaceDesc,
				&workspace.NumViews,
				&workspace.NumChannels,
				&workspace.Levelc,
				&workspace.ParentID,
				&workspace.NumChd,
				&workspace.UgroupID,
				&workspace.UserID,
				&workspace.Statusc,
				&workspace.CreatedAt,
				&workspace.UpdatedAt,
				&workspace.CreatedDay,
				&workspace.CreatedWeek,
				&workspace.CreatedMonth,
				&workspace.CreatedYear,
				&workspace.UpdatedDay,
				&workspace.UpdatedWeek,
				&workspace.UpdatedMonth,
				&workspace.UpdatedYear)
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4361}).Error(err)
				return nil, err
			}
			uuid4Str, err := common.UUIDBytesToStr(workspace.UUID4)
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4362}).Error(err)
				return nil, err
			}
			workspace.IDS = uuid4Str

			pohs = append(pohs, &workspace)
		}
		err = rows.Close()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4363}).Error(err)
			return nil, err
		}
		return pohs, nil
	}
}

// GetParentWorkspace - Get Parent Workspace
func (c *WorkspaceService) GetParentWorkspace(ctx context.Context, ID string, userEmail string, requestID string) (*Workspace, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4364}).Error(err)
		return nil, err
	default:
		pworkspace, err := c.GetWorkspace(ctx, ID, userEmail, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4365}).Error(err)
			return nil, err
		}
		workspace := Workspace{}
		db := c.DBService.DB
		row := db.QueryRowContext(ctx, `select
      id,
			uuid4,
			workspace_name,
			workspace_desc,
			num_views,
			num_channels,
			levelc,
			parent_id,
			num_chd,
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
			updated_year from workspaces where id = ? and statusc = ?;`, pworkspace.ParentID, common.Active)

		err = row.Scan(
			&workspace.ID,
			&workspace.UUID4,
			&workspace.WorkspaceName,
			&workspace.WorkspaceDesc,
			&workspace.NumViews,
			&workspace.NumChannels,
			&workspace.Levelc,
			&workspace.ParentID,
			&workspace.NumChd,
			&workspace.UgroupID,
			&workspace.UserID,
			/*  StatusDates  */
			&workspace.Statusc,
			&workspace.CreatedAt,
			&workspace.UpdatedAt,
			&workspace.CreatedDay,
			&workspace.CreatedWeek,
			&workspace.CreatedMonth,
			&workspace.CreatedYear,
			&workspace.UpdatedDay,
			&workspace.UpdatedWeek,
			&workspace.UpdatedMonth,
			&workspace.UpdatedYear)

		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4366}).Error(err)
			return nil, err
		}
		uuid4Str, err := common.UUIDBytesToStr(workspace.UUID4)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4367}).Error(err)
			return nil, err
		}
		workspace.IDS = uuid4Str

		return &workspace, nil
	}
}

//UpdateWorkspace - Update workspace
func (c *WorkspaceService) UpdateWorkspace(ctx context.Context, ID string, form *Workspace, UserID string, userEmail string, requestID string) error {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4376}).Error(err)
		return err
	default:
		workspace, err := c.GetWorkspace(ctx, ID, userEmail, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4377}).Error(err)
			return err
		}

		db := c.DBService.DB
		tn, tnday, tnweek, tnmonth, tnyear := common.GetTimeDetails()
		stmt, err := db.PrepareContext(ctx, `update workspaces set 
		  workspace_name = ?,
      workspace_desc = ?,
			updated_at = ?, 
			updated_day = ?, 
			updated_week = ?, 
			updated_month = ?, 
			updated_year = ? where id = ? and statusc = ?;`)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4379}).Error(err)
			return err
		}
		tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4378}).Error(err)
			return err
		}

		_, err = tx.StmtContext(ctx, stmt).Exec(
			form.WorkspaceName,
			form.WorkspaceDesc,
			tn,
			tnday,
			tnweek,
			tnmonth,
			tnyear,
			workspace.ID,
			common.Active)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4381}).Error(err)
			err = tx.Rollback()
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4390}).Error(err)
				return err
			}
			err = stmt.Close()
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4382}).Error(err)
				return err
			}
			return err
		}
		err = tx.Commit()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4384}).Error(err)
			return err
		}

		err = stmt.Close()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4383}).Error(err)
			return err
		}

		return nil
	}
}

// DeleteWorkspace - Delete workspace
func (c *WorkspaceService) DeleteWorkspace(ctx context.Context, ID string, userEmail string, requestID string) error {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4385}).Error(err)
		return err
	default:
		uuid4byte, err := common.UUIDStrToBytes(ID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4386}).Error(err)
			return err
		}
		db := c.DBService.DB
		tn, tnday, tnweek, tnmonth, tnyear := common.GetTimeDetails()
		stmt, err := db.PrepareContext(ctx, `update workspaces set 
		  statusc = ?,
			updated_at = ?, 
			updated_day = ?, 
			updated_week = ?, 
			updated_month = ?, 
			updated_year = ? where uuid4= ?;`)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4388}).Error(err)
			return err
		}
		tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4387}).Error(err)
			err = stmt.Close()
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4390}).Error(err)
				return err
			}
			return err
		}
		_, err = tx.StmtContext(ctx, stmt).Exec(
			common.Inactive,
			tn,
			tnday,
			tnweek,
			tnmonth,
			tnyear,
			uuid4byte)

		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4389}).Error(err)
			err = tx.Rollback()
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4390}).Error(err)
				return err
			}
			err = stmt.Close()
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4390}).Error(err)
				return err
			}
			return err
		}

		err = tx.Commit()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4392}).Error(err)
			err = tx.Rollback()
			return err
		}

		err = stmt.Close()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4391}).Error(err)
			return err
		}
		return nil
	}
}
