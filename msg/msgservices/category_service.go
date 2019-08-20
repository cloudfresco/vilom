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

// Category - Category view representation
type Category struct {
	ID           uint   `json:"id,omitempty"`
	UUID4        []byte `json:"-"`
	IDS          string `json:"id_s,omitempty"`
	CategoryName string `json:"category_name,omitempty"`
	CategoryDesc string `json:"category_desc,omitempty"`
	NumViews     uint   `json:"num_views,omitempty"`
	NumTopics    uint   `json:"num_topics,omitempty"`
	Levelc       uint   `json:"levelc,omitempty"`
	ParentID     uint   `json:"parent_id,omitempty"`
	NumChd       uint   `json:"num_chd,omitempty"`

	UgroupID uint `json:"ugroup_id,omitempty"`
	UserID   uint `json:"user_id,omitempty"`

	common.StatusDates
	Topics []*Topic
}

// CategoryChd - CategoryChd view representation
type CategoryChd struct {
	ID            uint   `json:"id,omitempty"`
	UUID4         []byte `json:"-"`
	CategoryID    uint   `json:"category_id,omitempty"`
	CategoryChdID uint   `json:"category_chd_id,omitempty"`

	common.StatusDates
}

// CategoryServiceIntf - interface for Category Service
type CategoryServiceIntf interface {
	CreateCategory(ctx context.Context, form *Category, UserID string, userEmail string, requestID string) (*Category, error)
	CreateChild(ctx context.Context, form *Category, UserID string, userEmail string, requestID string) (*Category, error)
	GetCategories(ctx context.Context, limit string, nextCursor string, userEmail string, requestID string) (*CategoryCursor, error)
	GetCategoryWithTopics(ctx context.Context, ID string, userEmail string, requestID string) (*Category, error)
	GetCategory(ctx context.Context, ID string, userEmail string, requestID string) (*Category, error)
	GetCategoryByID(ctx context.Context, ID uint, userEmail string, requestID string) (*Category, error)
	GetTopLevelCategories(ctx context.Context, userEmail string, requestID string) ([]*Category, error)
	GetChildCategories(ctx context.Context, ID string, userEmail string, requestID string) ([]*Category, error)
	GetParentCategory(ctx context.Context, ID string, userEmail string, requestID string) (*Category, error)
	UpdateCategory(ctx context.Context, ID string, form *Category, UserID string, userEmail string, requestID string) error
	UpdateNumTopicsPrepare(ctx context.Context, userEmail string, requestID string) (*sql.Stmt, error)
	UpdateNumTopics(ctx context.Context, stmt *sql.Stmt, tx *sql.Tx, numTopics uint, ID uint, userEmail string, requestID string) error
	DeleteCategory(ctx context.Context, ID string, userEmail string, requestID string) error
}

// CategoryService - For accessing category services
type CategoryService struct {
	DBService    *common.DBService
	RedisService *common.RedisService
}

// NewCategoryService - Create category service
func NewCategoryService(dbOpt *common.DBService, redisOpt *common.RedisService) *CategoryService {
	return &CategoryService{
		DBService:    dbOpt,
		RedisService: redisOpt,
	}
}

// CategoryCursor - used to get categories
type CategoryCursor struct {
	Categories []*Category
	NextCursor string `json:"next_cursor,omitempty"`
}

// CreateCategory - Create Category
func (c *CategoryService) CreateCategory(ctx context.Context, form *Category, UserID string, userEmail string, requestID string) (*Category, error) {
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
		insertCategoryStmt, err := c.insertCategoryPrepare(ctx, userEmail, requestID)
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
		cat := Category{}
		cat.UUID4, err = common.GetUUIDBytes()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4320}).Error(err)
			return nil, err
		}
		cat.CategoryName = form.CategoryName
		cat.CategoryDesc = form.CategoryDesc
		cat.NumViews = 0
		cat.NumTopics = 0
		cat.Levelc = 0
		cat.ParentID = uint(0)
		cat.NumChd = 0
		cat.UgroupID = uint(0)
		cat.UserID = user.ID
		/*  StatusDates  */
		cat.Statusc = common.Active
		cat.CreatedAt = tn
		cat.UpdatedAt = tn
		cat.CreatedDay = tnday
		cat.CreatedWeek = tnweek
		cat.CreatedMonth = tnmonth
		cat.CreatedYear = tnyear
		cat.UpdatedDay = tnday
		cat.UpdatedWeek = tnweek
		cat.UpdatedMonth = tnmonth
		cat.UpdatedYear = tnyear

		err = c.insertCategory(ctx, insertCategoryStmt, tx, &cat, userEmail, requestID)

		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4321}).Error(err)
			err = tx.Rollback()
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4328}).Error(err)
				return nil, err
			}
			err = insertCategoryStmt.Close()
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

		err = insertCategoryStmt.Close()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4328}).Error(err)
			return nil, err
		}

		return &cat, nil
	}
}

// insertCategoryPrepare - Insert Category prepare statement
func (c *CategoryService) insertCategoryPrepare(ctx context.Context, userEmail string, requestID string) (*sql.Stmt, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4323}).Error(err)
		return nil, err
	default:
		db := c.DBService.DB
		stmt, err := db.PrepareContext(ctx, `insert into categories
	  ( 
			uuid4,
			category_name,
			category_desc,
			num_views,
			num_topics,
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

// insertCategory - Insert category details into database
func (c *CategoryService) insertCategory(ctx context.Context, stmt *sql.Stmt, tx *sql.Tx, cat *Category, userEmail string, requestID string) error {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4323}).Error(err)
		return err
	default:
		res, err := tx.StmtContext(ctx, stmt).Exec(
			cat.UUID4,
			cat.CategoryName,
			cat.CategoryDesc,
			cat.NumViews,
			cat.NumTopics,
			cat.Levelc,
			cat.ParentID,
			cat.NumChd,
			cat.UgroupID,
			cat.UserID,
			/*  StatusDates  */
			cat.Statusc,
			cat.CreatedAt,
			cat.UpdatedAt,
			cat.CreatedDay,
			cat.CreatedWeek,
			cat.CreatedMonth,
			cat.CreatedYear,
			cat.UpdatedDay,
			cat.UpdatedWeek,
			cat.UpdatedMonth,
			cat.UpdatedYear)

		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4325}).Error(err)
			return err
		}

		uID, err := res.LastInsertId()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4326}).Error(err)
			return err
		}
		cat.ID = uint(uID)
		uuid4Str, err := common.UUIDBytesToStr(cat.UUID4)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4327}).Error(err)
			return err
		}
		cat.IDS = uuid4Str
		return nil
	}
}

// CreateChild - Create Child Category
func (c *CategoryService) CreateChild(ctx context.Context, form *Category, UserID string, userEmail string, requestID string) (*Category, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4338}).Error(err)
		return nil, err
	default:
		db := c.DBService.DB
		insertCategoryStmt, insertChildStmt, updateNumChildrenStmt, err := c.createChildPrepareStmts(ctx, userEmail, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4340}).Error(err)
			return nil, err
		}

		tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4369}).Error(err)
			err = c.createChildPrepareStmtsClose(ctx, insertCategoryStmt, insertChildStmt, updateNumChildrenStmt, userEmail, requestID)
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4340}).Error(err)
				return nil, err
			}
			return nil, err
		}

		cat, err := c.createChild(ctx, insertCategoryStmt, insertChildStmt, updateNumChildrenStmt, tx, form, UserID, userEmail, requestID)

		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4342}).Error(err)
			err = tx.Rollback()
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4340}).Error(err)
				return nil, err
			}
			err = c.createChildPrepareStmtsClose(ctx, insertCategoryStmt, insertChildStmt, updateNumChildrenStmt, userEmail, requestID)
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

		err = c.createChildPrepareStmtsClose(ctx, insertCategoryStmt, insertChildStmt, updateNumChildrenStmt, userEmail, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4340}).Error(err)
			return nil, err
		}
		return cat, nil
	}
}

//createChildPrepareStmts - Prepare Statements
func (c *CategoryService) createChildPrepareStmts(ctx context.Context, userEmail string, requestID string) (*sql.Stmt, *sql.Stmt, *sql.Stmt, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6311}).Error(err)
		return nil, nil, nil, err
	default:
		insertCategoryStmt, err := c.insertCategoryPrepare(ctx, userEmail, requestID)
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
		return insertCategoryStmt, insertChildStmt, updateNumChildrenStmt, nil
	}
}

//createChildPrepareStmtsClose - Close Prepare Statements
func (c *CategoryService) createChildPrepareStmtsClose(ctx context.Context, insertCategoryStmt *sql.Stmt, insertChildStmt *sql.Stmt, updateNumChildrenStmt *sql.Stmt, userEmail string, requestID string) error {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6311}).Error(err)
		return err
	default:
		err := insertCategoryStmt.Close()
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

// createChild - Create Child Category
func (c *CategoryService) createChild(ctx context.Context, insertCategoryStmt *sql.Stmt, insertChildStmt *sql.Stmt, updateNumChildrenStmt *sql.Stmt, tx *sql.Tx, form *Category, UserID string, userEmail string, requestID string) (*Category, error) {
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

		parent, err := c.GetCategoryByID(ctx, form.ParentID, userEmail, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4340}).Error(err)
			return nil, err
		}

		tn, tnday, tnweek, tnmonth, tnyear := common.GetTimeDetails()
		cat := Category{}
		cat.UUID4, err = common.GetUUIDBytes()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4341}).Error(err)
			return nil, err
		}
		cat.CategoryName = form.CategoryName
		cat.CategoryDesc = form.CategoryDesc
		cat.NumViews = 0
		cat.NumTopics = 0
		cat.Levelc = parent.Levelc + 1
		cat.ParentID = parent.ID
		cat.NumChd = 0
		cat.UgroupID = uint(0)
		cat.UserID = user.ID
		/*  StatusDates  */
		cat.Statusc = common.Active
		cat.CreatedAt = tn
		cat.UpdatedAt = tn
		cat.CreatedDay = tnday
		cat.CreatedWeek = tnweek
		cat.CreatedMonth = tnmonth
		cat.CreatedYear = tnyear
		cat.UpdatedDay = tnday
		cat.UpdatedWeek = tnweek
		cat.UpdatedMonth = tnmonth
		cat.UpdatedYear = tnyear

		err = c.insertCategory(ctx, insertCategoryStmt, tx, &cat, userEmail, requestID)

		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4342}).Error(err)
			return nil, err
		}

		catChd := CategoryChd{}
		catChd.UUID4, err = common.GetUUIDBytes()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4343}).Error(err)
			return nil, err
		}
		catChd.CategoryID = parent.ID
		catChd.CategoryChdID = cat.ID
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

		return &cat, nil
	}
}

// insertChildPrepare - Insert child prepare statement
func (c *CategoryService) insertChildPrepare(ctx context.Context, userEmail string, requestID string) (*sql.Stmt, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4348}).Error(err)
		return nil, err
	default:
		db := c.DBService.DB
		stmt, err := db.PrepareContext(ctx, `insert into category_chds
	  ( 
    uuid4,
		category_id,
		category_chd_id,
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

// insertChild - Insert child category details into database
func (c *CategoryService) insertChild(ctx context.Context, stmt *sql.Stmt, tx *sql.Tx, catChd *CategoryChd, userEmail string, requestID string) error {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4348}).Error(err)
		return err
	default:
		res, err := tx.StmtContext(ctx, stmt).Exec(
			catChd.UUID4,
			catChd.CategoryID,
			catChd.CategoryChdID,
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
func (c *CategoryService) updateNumChildrenPrepare(ctx context.Context, userEmail string, requestID string) (*sql.Stmt, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4375}).Error(err)
		return nil, err
	default:
		db := c.DBService.DB
		stmt, err := db.PrepareContext(ctx, `update categories set 
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

// updateNumChildren - Update number of child in category
func (c *CategoryService) updateNumChildren(ctx context.Context, stmt *sql.Stmt, tx *sql.Tx, numchd uint, parentID uint, userEmail string, requestID string) error {
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

// GetCategories - Get Categories
func (c *CategoryService) GetCategories(ctx context.Context, limit string, nextCursor string, userEmail string, requestID string) (*CategoryCursor, error) {
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
		cats := []*Category{}
		rows, err := db.QueryContext(ctx, `select 
      id, 
			uuid4,
			category_name,
			category_desc,
			num_views,
			num_topics,
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
			updated_year from categories where `+query, 0, common.Active)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4301}).Error(err)
			return nil, err
		}

		for rows.Next() {
			cat := Category{}
			err = rows.Scan(
				&cat.ID,
				&cat.UUID4,
				&cat.CategoryName,
				&cat.CategoryDesc,
				&cat.NumViews,
				&cat.NumTopics,
				&cat.Levelc,
				&cat.ParentID,
				&cat.NumChd,
				&cat.UgroupID,
				&cat.UserID,
				/*  StatusDates  */
				&cat.Statusc,
				&cat.CreatedAt,
				&cat.UpdatedAt,
				&cat.CreatedDay,
				&cat.CreatedWeek,
				&cat.CreatedMonth,
				&cat.CreatedYear,
				&cat.UpdatedDay,
				&cat.UpdatedWeek,
				&cat.UpdatedMonth,
				&cat.UpdatedYear)
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4302}).Error(err)
				return nil, err
			}
			uuid4Str, err := common.UUIDBytesToStr(cat.UUID4)
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4303}).Error(err)
				return nil, err
			}
			cat.IDS = uuid4Str
			cats = append(cats, &cat)
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
		x := CategoryCursor{}
		if len(cats) != 0 {
			next := cats[len(cats)-1].ID
			next = next - 1
			nextc := common.EncodeCursor(next)
			x = CategoryCursor{cats, nextc}
		} else {
			x = CategoryCursor{cats, "0"}
		}
		return &x, nil
	}
}

// GetCategoryWithTopics - Get category with topics
func (c *CategoryService) GetCategoryWithTopics(ctx context.Context, ID string, userEmail string, requestID string) (*Category, error) {
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
		cat := Category{}
		ctegry, err := c.GetCategory(ctx, ID, userEmail, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4331}).Error(err)
			return nil, err
		}
		var isPresent bool
		row := db.QueryRowContext(ctx, `select exists (select 1 from topics where category_id = ?);`, ctegry.ID)
		err = row.Scan(&isPresent)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4332}).Error(err)
			return nil, err
		}
		if isPresent {

			rows, err := db.QueryContext(ctx, `select 
		  c.id,
			c.uuid4,
			c.category_name,
			c.category_desc,
			c.num_views,
			c.num_topics,
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
			v.topic_name,
			v.topic_desc,
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
			v.category_id,
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
			v.updated_year from categories c inner join topics v on (c.id = v.category_id) where c.uuid4 = ?`, uuid4byte)

			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4333}).Error(err)
				return nil, err
			}
			for rows.Next() {
				topc := Topic{}
				err = rows.Scan(
					&cat.ID,
					&cat.UUID4,
					&cat.CategoryName,
					&cat.CategoryDesc,
					&cat.NumViews,
					&cat.NumTopics,
					&cat.Levelc,
					&cat.ParentID,
					&cat.NumChd,
					&cat.UgroupID,
					&cat.UserID,
					/*  StatusDates  */
					&cat.Statusc,
					&cat.CreatedAt,
					&cat.UpdatedAt,
					&cat.CreatedDay,
					&cat.CreatedWeek,
					&cat.CreatedMonth,
					&cat.CreatedYear,
					&cat.UpdatedDay,
					&cat.UpdatedWeek,
					&cat.UpdatedMonth,
					&cat.UpdatedYear,
					&topc.ID,
					&topc.UUID4,
					&topc.TopicName,
					&topc.TopicDesc,
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
					&topc.CategoryID,
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
				uuid4Str1, err := common.UUIDBytesToStr(cat.UUID4)
				if err != nil {
					log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4335}).Error(err)
					return nil, err
				}
				cat.IDS = uuid4Str1

				uuid4Str, err := common.UUIDBytesToStr(topc.UUID4)
				if err != nil {
					log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4336}).Error(err)
					return nil, err
				}
				topc.IDS = uuid4Str
				cat.Topics = append(cat.Topics, &topc)
			}

			err = rows.Close()
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4337}).Error(err)
				return nil, err
			}
		} else {
			return ctegry, nil
		}
		return &cat, nil
	}
}

// GetCategory - Get Category
func (c *CategoryService) GetCategory(ctx context.Context, ID string, userEmail string, requestID string) (*Category, error) {
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
		cat := Category{}
		db := c.DBService.DB
		row := db.QueryRowContext(ctx, `select
      id,
			uuid4,
			category_name,
			category_desc,
			num_views,
			num_topics,
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
			updated_year from categories where uuid4 = ? and statusc = ?;`, uuid4byte, common.Active)

		err = row.Scan(
			&cat.ID,
			&cat.UUID4,
			&cat.CategoryName,
			&cat.CategoryDesc,
			&cat.NumViews,
			&cat.NumTopics,
			&cat.Levelc,
			&cat.ParentID,
			&cat.NumChd,
			&cat.UgroupID,
			&cat.UserID,
			&cat.Statusc,
			&cat.CreatedAt,
			&cat.UpdatedAt,
			&cat.CreatedDay,
			&cat.CreatedWeek,
			&cat.CreatedMonth,
			&cat.CreatedYear,
			&cat.UpdatedDay,
			&cat.UpdatedWeek,
			&cat.UpdatedMonth,
			&cat.UpdatedYear)

		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4308}).Error(err)
			return nil, err
		}
		uuid4Str, err := common.UUIDBytesToStr(cat.UUID4)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4309}).Error(err)
			return nil, err
		}
		cat.IDS = uuid4Str

		return &cat, nil
	}
}

// GetCategoryByID - Get Category By ID
func (c *CategoryService) GetCategoryByID(ctx context.Context, ID uint, userEmail string, requestID string) (*Category, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4310}).Error(err)
		return nil, err
	default:
		cat := Category{}
		db := c.DBService.DB
		row := db.QueryRowContext(ctx, `select
      id,
			uuid4,
			category_name,
			category_desc,
			num_views,
			num_topics,
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
			updated_year from categories where id = ? and statusc = ?;`, ID, common.Active)

		err := row.Scan(
			&cat.ID,
			&cat.UUID4,
			&cat.CategoryName,
			&cat.CategoryDesc,
			&cat.NumViews,
			&cat.NumTopics,
			&cat.Levelc,
			&cat.ParentID,
			&cat.NumChd,
			&cat.UgroupID,
			&cat.UserID,
			/*  StatusDates  */
			&cat.Statusc,
			&cat.CreatedAt,
			&cat.UpdatedAt,
			&cat.CreatedDay,
			&cat.CreatedWeek,
			&cat.CreatedMonth,
			&cat.CreatedYear,
			&cat.UpdatedDay,
			&cat.UpdatedWeek,
			&cat.UpdatedMonth,
			&cat.UpdatedYear)

		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4311}).Error(err)
			return nil, err
		}
		uuid4Str, err := common.UUIDBytesToStr(cat.UUID4)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4312}).Error(err)
			return nil, err
		}
		cat.IDS = uuid4Str
		return &cat, nil
	}
}

// GetTopLevelCategories - Get top level categories
func (c *CategoryService) GetTopLevelCategories(ctx context.Context, userEmail string, requestID string) ([]*Category, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4353}).Error(err)
		return nil, err
	default:
		cats := []*Category{}
		db := c.DBService.DB
		rows, err := db.QueryContext(ctx, `select 
      id, 
			uuid4,
			category_name,
			category_desc,
			num_views,
			num_topics,
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
			updated_year from categories where levelc = ? and statusc = ?;`, 0, common.Active)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4354}).Error(err)
			return nil, err
		}

		for rows.Next() {
			cat := Category{}
			err = rows.Scan(
				&cat.ID,
				&cat.UUID4,
				&cat.CategoryName,
				&cat.CategoryDesc,
				&cat.NumViews,
				&cat.NumTopics,
				&cat.Levelc,
				&cat.ParentID,
				&cat.NumChd,
				&cat.UgroupID,
				&cat.UserID,
				/*  StatusDates  */
				&cat.Statusc,
				&cat.CreatedAt,
				&cat.UpdatedAt,
				&cat.CreatedDay,
				&cat.CreatedWeek,
				&cat.CreatedMonth,
				&cat.CreatedYear,
				&cat.UpdatedDay,
				&cat.UpdatedWeek,
				&cat.UpdatedMonth,
				&cat.UpdatedYear)
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4355}).Error(err)
				return nil, err
			}
			uuid4Str, err := common.UUIDBytesToStr(cat.UUID4)
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4356}).Error(err)
				return nil, err
			}
			cat.IDS = uuid4Str
			cats = append(cats, &cat)
		}
		err = rows.Close()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4357}).Error(err)
			return nil, err
		}
		return cats, nil
	}
}

// GetChildCategories - Get child categories
func (c *CategoryService) GetChildCategories(ctx context.Context, ID string, userEmail string, requestID string) ([]*Category, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4358}).Error(err)
		return nil, err
	default:
		category, err := c.GetCategory(ctx, ID, userEmail, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4359}).Error(err)
			return nil, err
		}
		pohs := []*Category{}
		db := c.DBService.DB
		rows, err := db.QueryContext(ctx, `select 
		    c.id,
				c.uuid4,
				c.category_name,
				c.category_desc,
				c.num_views,
				c.num_topics,
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
				c.updated_year from categories c inner join category_chds ch on (c.id = ch.category_chd_id) where ch.category_id = ?`, category.ID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4360}).Error(err)
			return nil, err
		}

		for rows.Next() {
			cat := Category{}
			err = rows.Scan(
				&cat.ID,
				&cat.UUID4,
				&cat.CategoryName,
				&cat.CategoryDesc,
				&cat.NumViews,
				&cat.NumTopics,
				&cat.Levelc,
				&cat.ParentID,
				&cat.NumChd,
				&cat.UgroupID,
				&cat.UserID,
				&cat.Statusc,
				&cat.CreatedAt,
				&cat.UpdatedAt,
				&cat.CreatedDay,
				&cat.CreatedWeek,
				&cat.CreatedMonth,
				&cat.CreatedYear,
				&cat.UpdatedDay,
				&cat.UpdatedWeek,
				&cat.UpdatedMonth,
				&cat.UpdatedYear)
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4361}).Error(err)
				return nil, err
			}
			uuid4Str, err := common.UUIDBytesToStr(cat.UUID4)
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4362}).Error(err)
				return nil, err
			}
			cat.IDS = uuid4Str

			pohs = append(pohs, &cat)
		}
		err = rows.Close()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4363}).Error(err)
			return nil, err
		}
		return pohs, nil
	}
}

// GetParentCategory - Get Parent Category
func (c *CategoryService) GetParentCategory(ctx context.Context, ID string, userEmail string, requestID string) (*Category, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4364}).Error(err)
		return nil, err
	default:
		category, err := c.GetCategory(ctx, ID, userEmail, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4365}).Error(err)
			return nil, err
		}
		cat := Category{}
		db := c.DBService.DB
		row := db.QueryRowContext(ctx, `select
      id,
			uuid4,
			category_name,
			category_desc,
			num_views,
			num_topics,
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
			updated_year from categories where id = ? and statusc = ?;`, category.ParentID, common.Active)

		err = row.Scan(
			&cat.ID,
			&cat.UUID4,
			&cat.CategoryName,
			&cat.CategoryDesc,
			&cat.NumViews,
			&cat.NumTopics,
			&cat.Levelc,
			&cat.ParentID,
			&cat.NumChd,
			&cat.UgroupID,
			&cat.UserID,
			/*  StatusDates  */
			&cat.Statusc,
			&cat.CreatedAt,
			&cat.UpdatedAt,
			&cat.CreatedDay,
			&cat.CreatedWeek,
			&cat.CreatedMonth,
			&cat.CreatedYear,
			&cat.UpdatedDay,
			&cat.UpdatedWeek,
			&cat.UpdatedMonth,
			&cat.UpdatedYear)

		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4366}).Error(err)
			return nil, err
		}
		uuid4Str, err := common.UUIDBytesToStr(cat.UUID4)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4367}).Error(err)
			return nil, err
		}
		cat.IDS = uuid4Str

		return &cat, nil
	}
}

//UpdateCategory - Update category
func (c *CategoryService) UpdateCategory(ctx context.Context, ID string, form *Category, UserID string, userEmail string, requestID string) error {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4376}).Error(err)
		return err
	default:
		category, err := c.GetCategory(ctx, ID, userEmail, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4377}).Error(err)
			return err
		}

		db := c.DBService.DB
		tn, tnday, tnweek, tnmonth, tnyear := common.GetTimeDetails()
		stmt, err := db.PrepareContext(ctx, `update categories set 
		  category_name = ?,
      category_desc = ?,
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
			form.CategoryName,
			form.CategoryDesc,
			tn,
			tnday,
			tnweek,
			tnmonth,
			tnyear,
			category.ID,
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

// UpdateNumTopicsPrepare - UpdateNumTopics Prepare Statement
func (c *CategoryService) UpdateNumTopicsPrepare(ctx context.Context, userEmail string, requestID string) (*sql.Stmt, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4374}).Error(err)
		return nil, err
	default:
		db := c.DBService.DB
		stmt, err := db.PrepareContext(ctx, `update categories set 
    num_topics = ?,
	  updated_at = ?, 
		updated_day = ?, 
		updated_week = ?, 
		updated_month = ?, 
		updated_year = ? where id = ? and statusc = ?;`)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5327}).Error(err)
			return nil, err
		}
		return stmt, nil
	}
}

// UpdateNumTopics - update number of topics in category
func (c *CategoryService) UpdateNumTopics(ctx context.Context, stmt *sql.Stmt, tx *sql.Tx, numTopics uint, ID uint, userEmail string, requestID string) error {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4374}).Error(err)
		return err
	default:
		tn, tnday, tnweek, tnmonth, tnyear := common.GetTimeDetails()

		_, err := tx.StmtContext(ctx, stmt).Exec(
			numTopics,
			tn,
			tnday,
			tnweek,
			tnmonth,
			tnyear,
			ID,
			common.Active)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5328}).Error(err)
			return err
		}
		return nil
	}
}

// DeleteCategory - Delete category
func (c *CategoryService) DeleteCategory(ctx context.Context, ID string, userEmail string, requestID string) error {
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
		stmt, err := db.PrepareContext(ctx, `update categories set 
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
