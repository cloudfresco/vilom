package msgservices

import (
	"context"
	"database/sql"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/go-redis/redis"
	"github.com/palantir/stacktrace"

	"github.com/cloudfresco/vilom/common"
	"github.com/cloudfresco/vilom/user/userservices"
)

// Category - Category view representation
type Category struct {
	ID           uint
	IDS          string
	CategoryName string
	CategoryDesc string
	NumViews     uint `sql:"default:'0'"`
	NumTopics    uint `sql:"default:'0'"`
	Levelc       uint `gorm:"type:tinyint"`
	ParentID     uint
	NumChd       uint `gorm:"type:smallint"`

	UgroupID uint `sql:"default:'0'"`
	UserID   uint

	common.StatusDates
	Topics []*Topic
}

// CategoryChd - CategoryChd view representation
type CategoryChd struct {
	ID uint

	CategoryID    uint
	CategoryChdID uint

	common.StatusDates
}

// CategoryService - For accessing category services
type CategoryService struct {
	Config       *common.RedisOptions
	Db           *sql.DB
	RedisClient  *redis.Client
	LimitDefault string
}

// NewCategoryService - Create category service
func NewCategoryService(config *common.RedisOptions,
	db *sql.DB,
	redisClient *redis.Client,
	limitDefault string) *CategoryService {
	return &CategoryService{config, db, redisClient, limitDefault}
}

// CategoryCursor - used to get categories
type CategoryCursor struct {
	Categories []*Category
	NextCursor string
}

// GetCategories - Get Categories
func (c *CategoryService) GetCategories(ctx context.Context, limit string, nextCursor string) (*CategoryCursor, error) {
	if limit == "" {
		limit = c.LimitDefault
	}
	query := "(levelc = ?)"
	if nextCursor == "" {
		query = query + " order by id desc " + " limit " + limit + ";"
	} else {
		nextCursor = common.DecodeCursor(nextCursor)
		query = query + " " + "and" + " " + "id <= " + nextCursor + " order by id desc " + " limit " + limit + ";"
	}

	cats := []*Category{}
	rows, err := c.Db.QueryContext(ctx, `select 
      id, 
			id_s,
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
			updated_year from categories where `+query, 0)
	if err != nil {
		log.Error(stacktrace.Propagate(err, ""))
		return nil, err
	}

	for rows.Next() {
		cat := Category{}
		err = rows.Scan(
			&cat.ID,
			&cat.IDS,
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
			log.Error(stacktrace.Propagate(err, ""))
			return nil, err
		}
		cats = append(cats, &cat)
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

	next := cats[len(cats)-1].ID
	next = next - 1
	nextc := common.EncodeCursor(next)
	x := CategoryCursor{cats, nextc}
	return &x, nil
}

// GetCategory - Get Category
func (c *CategoryService) GetCategory(ctx context.Context, ID string) (*Category, error) {
	cat := Category{}
	row := c.Db.QueryRowContext(ctx, `select
      id,
			id_s,
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
			updated_year from categories where id_s = ?;`, ID)

	err := row.Scan(
		&cat.ID,
		&cat.IDS,
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
		log.Error(stacktrace.Propagate(err, ""))
		return nil, err
	}

	return &cat, nil
}

// GetCategoryByID - Get Category By ID
func (c *CategoryService) GetCategoryByID(ctx context.Context, ID uint) (*Category, error) {
	cat := Category{}
	row := c.Db.QueryRowContext(ctx, `select
      id,
			id_s,
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
			updated_year from categories where id = ?;`, ID)

	err := row.Scan(
		&cat.ID,
		&cat.IDS,
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
		log.Error(stacktrace.Propagate(err, ""))
		return nil, err
	}

	return &cat, nil
}

// UpdateCategory - Update Category
func (c *CategoryService) UpdateCategory(ctx context.Context, form *Category, ID string) error {
	tn := time.Now().UTC()
	_, week := tn.ISOWeek()
	day := tn.YearDay()

	UpdatedDay := uint(day)
	UpdatedWeek := uint(week)
	UpdatedMonth := uint(tn.Month())
	UpdatedYear := uint(tn.Year())

	stmt, err := c.Db.PrepareContext(ctx, `update categories set 
				  category_name = ?,
				  updated_at = ?, 
					updated_day = ?, 
					updated_week = ?, 
					updated_month = ?, 
					updated_year = ? where id_s= ?;`)
	if err != nil {
		log.Error(stacktrace.Propagate(err, ""))
		err = stmt.Close()
		return err
	}

	_, err = stmt.ExecContext(ctx,
		form.CategoryName,
		tn,
		UpdatedDay,
		UpdatedWeek,
		UpdatedMonth,
		UpdatedYear,
		ID)
	if err != nil {
		log.Error(stacktrace.Propagate(err, ""))
		err = stmt.Close()
		return err
	}
	err = stmt.Close()

	if err != nil {
		log.Error(stacktrace.Propagate(err, ""))
		return err
	}

	return nil
}

// Create - Create Category
func (c *CategoryService) Create(ctx context.Context, form *Category, UserID string) (*Category, error) {
	userserv := &userservices.UserService{Config: c.Config, Db: c.Db, RedisClient: c.RedisClient}
	user, err := userserv.GetUser(ctx, UserID)
	if err != nil {
		log.Error(stacktrace.Propagate(err, ""))
		return nil, err
	}
	db := c.Db
	tx, err := db.Begin()
	tn := time.Now().UTC()
	_, week := tn.ISOWeek()
	day := tn.YearDay()
	cat := Category{}
	cat.IDS = common.GetUID()
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
	cat.CreatedDay = uint(day)
	cat.CreatedWeek = uint(week)
	cat.CreatedMonth = uint(tn.Month())
	cat.CreatedYear = uint(tn.Year())
	cat.UpdatedDay = uint(day)
	cat.UpdatedWeek = uint(week)
	cat.UpdatedMonth = uint(tn.Month())
	cat.UpdatedYear = uint(tn.Year())

	Cat, err := c.InsertCategory(ctx, tx, cat)

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
	return Cat, nil
}

// InsertCategory - Insert category details into database
func (c *CategoryService) InsertCategory(ctx context.Context, tx *sql.Tx, cat Category) (*Category, error) {
	stmt, err := tx.PrepareContext(ctx, `insert into categories
	  ( 
			id_s,
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
		log.Error(stacktrace.Propagate(err, ""))
		err = stmt.Close()
		return nil, err
	}
	res, err := stmt.ExecContext(ctx,
		cat.IDS,
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
	cat.ID = uint(uID)
	err = stmt.Close()
	if err != nil {
		log.Error(stacktrace.Propagate(err, ""))
		return nil, err
	}
	return &cat, nil
}

// GetCategoryWithTopics - Get category with topics
func (c *CategoryService) GetCategoryWithTopics(ctx context.Context, ID string) (*Category, error) {
	db := c.Db
	cat := &Category{}
	ctegry, err := c.GetCategory(ctx, ID)
	if err != nil {
		log.Error(stacktrace.Propagate(err, ""))
		return nil, err
	}
	var isPresent bool
	row := db.QueryRowContext(ctx, `select exists (select 1 from topics where category_id = ?);`, ctegry.ID)
	err = row.Scan(&isPresent)
	if err != nil {
		log.Error(stacktrace.Propagate(err, ""))
		return nil, err
	}
	if isPresent {

		rows, err := db.QueryContext(ctx, `select 
		  c.id,
			c.id_s,
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
			v.id_s,
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
			v.updated_year from categories c inner join topics v on (c.id = v.category_id) where c.id_s = ?`, ID)

		if err != nil {
			log.Error(stacktrace.Propagate(err, ""))
			return nil, err
		}
		for rows.Next() {
			topc := Topic{}
			err = rows.Scan(
				&cat.ID,
				&cat.IDS,
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
				&topc.IDS,
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
				log.Error(stacktrace.Propagate(err, ""))
				return nil, err
			}

			cat.Topics = append(cat.Topics, &topc)
		}

		err = rows.Close()
		if err != nil {
			log.Error(stacktrace.Propagate(err, ""))
			return nil, err
		}
	} else {
		cat = ctegry
	}
	return cat, nil
}

// CreateChild - Create Child Category
func (c *CategoryService) CreateChild(ctx context.Context, form *Category, UserID string) (*Category, error) {
	userserv := &userservices.UserService{Config: c.Config, Db: c.Db, RedisClient: c.RedisClient}
	user, err := userserv.GetUser(ctx, UserID)
	if err != nil {
		log.Error(stacktrace.Propagate(err, ""))
		return nil, err
	}

	parent, err := c.GetCategoryByID(ctx, form.ParentID)
	if err != nil {
		log.Error(stacktrace.Propagate(err, ""))
		return nil, err
	}

	db := c.Db
	tx, err := db.Begin()

	tn := time.Now().UTC()
	_, week := tn.ISOWeek()
	day := tn.YearDay()
	cat := Category{}
	cat.IDS = common.GetUID()
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
	cat.CreatedDay = uint(day)
	cat.CreatedWeek = uint(week)
	cat.CreatedMonth = uint(tn.Month())
	cat.CreatedYear = uint(tn.Year())
	cat.UpdatedDay = uint(day)
	cat.UpdatedWeek = uint(week)
	cat.UpdatedMonth = uint(tn.Month())
	cat.UpdatedYear = uint(tn.Year())

	Cat, err := c.InsertCategory(ctx, tx, cat)

	if err != nil {
		log.Error(stacktrace.Propagate(err, ""))
		err = tx.Rollback()
		return nil, err
	}

	catchd := CategoryChd{}
	catchd.CategoryID = parent.ID
	catchd.CategoryChdID = Cat.ID
	/*  StatusDates  */
	catchd.Statusc = common.Active
	catchd.CreatedAt = tn
	catchd.UpdatedAt = tn
	catchd.CreatedDay = uint(day)
	catchd.CreatedWeek = uint(week)
	catchd.CreatedMonth = uint(tn.Month())
	catchd.CreatedYear = uint(tn.Year())
	catchd.UpdatedDay = uint(day)
	catchd.UpdatedWeek = uint(week)
	catchd.UpdatedMonth = uint(tn.Month())
	catchd.UpdatedYear = uint(tn.Year())

	_, err = c.InsertChild(ctx, tx, catchd)

	if err != nil {
		log.Error(stacktrace.Propagate(err, ""))
		err = tx.Rollback()
		return nil, err
	}

	UpdatedDay := uint(day)
	UpdatedWeek := uint(week)
	UpdatedMonth := uint(tn.Month())
	UpdatedYear := uint(tn.Year())

	stmt, err := tx.PrepareContext(ctx, `update categories set 
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
	return Cat, nil
}

// InsertChild - Insert child category details into database
func (c *CategoryService) InsertChild(ctx context.Context, tx *sql.Tx, catchd CategoryChd) (*CategoryChd, error) {
	stmt, err := tx.PrepareContext(ctx, `insert into category_chds
	  ( 
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
					?,?,?);`)
	if err != nil {
		log.Error(stacktrace.Propagate(err, ""))
		err = stmt.Close()
		return nil, err
	}
	res, err := stmt.ExecContext(ctx,
		catchd.CategoryID,
		catchd.CategoryChdID,
		/*  StatusDates  */
		catchd.Statusc,
		catchd.CreatedAt,
		catchd.UpdatedAt,
		catchd.CreatedDay,
		catchd.CreatedWeek,
		catchd.CreatedMonth,
		catchd.CreatedYear,
		catchd.UpdatedDay,
		catchd.UpdatedWeek,
		catchd.UpdatedMonth,
		catchd.UpdatedYear)

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
	catchd.ID = uint(uID)
	err = stmt.Close()
	if err != nil {
		log.Error(stacktrace.Propagate(err, ""))
		return nil, err
	}
	return &catchd, nil
}

// GetTopLevelCategories - Get top level categories
func (c *CategoryService) GetTopLevelCategories(ctx context.Context) ([]*Category, error) {
	cats := []*Category{}
	rows, err := c.Db.QueryContext(ctx, `select 
      id, 
			id_s,
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
			updated_year from categories where levelc = ? and statusc = ?;`, 0, 1)
	if err != nil {
		log.Error(stacktrace.Propagate(err, ""))
		return nil, err
	}

	for rows.Next() {
		cat := Category{}
		err = rows.Scan(
			&cat.ID,
			&cat.IDS,
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
			log.Error(stacktrace.Propagate(err, ""))
			return nil, err
		}
		cats = append(cats, &cat)
	}
	err = rows.Close()
	if err != nil {
		log.Error(stacktrace.Propagate(err, ""))
		return nil, err
	}
	return cats, nil
}

// GetChildCategories - Get child categories
func (c *CategoryService) GetChildCategories(ctx context.Context, ID string) ([]*Category, error) {
	category, err := c.GetCategory(ctx, ID)
	if err != nil {
		log.Error(stacktrace.Propagate(err, ""))
		return nil, err
	}
	pohs := []*Category{}
	rows, err := c.Db.QueryContext(ctx, `select 
		    c.id,
				c.id_s,
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
		log.Error(stacktrace.Propagate(err, ""))
		return nil, err
	}

	for rows.Next() {
		cat := Category{}
		err = rows.Scan(
			&cat.ID,
			&cat.IDS,
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
			log.Error(stacktrace.Propagate(err, ""))
			return nil, err
		}
		pohs = append(pohs, &cat)
	}
	err = rows.Close()
	if err != nil {
		log.Error(stacktrace.Propagate(err, ""))
		return nil, err
	}
	return pohs, nil

}

// GetParentCategory - Get Parent Category
func (c *CategoryService) GetParentCategory(ctx context.Context, ID string) (*Category, error) {
	category, err := c.GetCategory(ctx, ID)
	if err != nil {
		log.Error(stacktrace.Propagate(err, ""))
		return nil, err
	}
	cat := Category{}
	row := c.Db.QueryRowContext(ctx, `select
      id,
			id_s,
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
			updated_year from categories where id = ?;`, category.ParentID)

	err = row.Scan(
		&cat.ID,
		&cat.IDS,
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
		log.Error(stacktrace.Propagate(err, ""))
		return nil, err
	}

	return &cat, nil
}
