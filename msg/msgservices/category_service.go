package msgservices

import (
	"context"
	"database/sql"
	"errors"

	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"

	"github.com/cloudfresco/vilom/common"
	"github.com/cloudfresco/vilom/user/userservices"
)

// Category - Category view representation
type Category struct {
	ID           uint
	UUID4        []byte
	IDS          string
	CategoryName string
	CategoryDesc string
	NumViews     uint
	NumTopics    uint
	Levelc       uint
	ParentID     uint
	NumChd       uint

	UgroupID uint
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
func (c *CategoryService) GetCategories(ctx context.Context, limit string, nextCursor string, userEmail string, requestID string) (*CategoryCursor, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4300}).Error(err)
		return nil, err
	default:
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
			updated_year from categories where `+query, 0)
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
			uUID4Str, err := common.UUIDBytesToStr(cat.UUID4)
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4303}).Error(err)
				return nil, err
			}
			cat.IDS = uUID4Str
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

// GetCategory - Get Category
func (c *CategoryService) GetCategory(ctx context.Context, ID string, userEmail string, requestID string) (*Category, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4306}).Error(err)
		return nil, err
	default:
		uUID4byte, err := common.UUIDStrToBytes(ID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4307}).Error(err)
			return nil, err
		}
		cat := Category{}
		row := c.Db.QueryRowContext(ctx, `select
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
			updated_year from categories where uuid4 = ?;`, uUID4byte)

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
		uUID4Str, err := common.UUIDBytesToStr(cat.UUID4)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4309}).Error(err)
			return nil, err
		}
		cat.IDS = uUID4Str

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
		row := c.Db.QueryRowContext(ctx, `select
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
			updated_year from categories where id = ?;`, ID)

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
		uUID4Str, err := common.UUIDBytesToStr(cat.UUID4)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4312}).Error(err)
			return nil, err
		}
		cat.IDS = uUID4Str
		return &cat, nil
	}
}

// UpdateCategory - Update Category
func (c *CategoryService) UpdateCategory(ctx context.Context, form *Category, ID string, userEmail string, requestID string) error {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4313}).Error(err)
		return err
	default:
		uUID4byte, err := common.UUIDStrToBytes(ID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4314}).Error(err)
			return err
		}

		tn, tnday, tnweek, tnmonth, tnyear := common.GetTimeDetails()
		UpdatedDay := tnday
		UpdatedWeek := tnweek
		UpdatedMonth := tnmonth
		UpdatedYear := tnyear

		stmt, err := c.Db.PrepareContext(ctx, `update categories set 
				  category_name = ?,
				  updated_at = ?, 
					updated_day = ?, 
					updated_week = ?, 
					updated_month = ?, 
					updated_year = ? where uuid4 = ?;`)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4315}).Error(err)
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
			uUID4byte)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4316}).Error(err)
			err = stmt.Close()
			return err
		}
		err = stmt.Close()

		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4317}).Error(err)
			return err
		}

		return nil
	}
}

// Create - Create Category
func (c *CategoryService) Create(ctx context.Context, form *Category, UserID string, userEmail string, requestID string) (*Category, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4318}).Error(err)
		return nil, err
	default:
		userserv := &userservices.UserService{Config: c.Config, Db: c.Db, RedisClient: c.RedisClient}
		user, err := userserv.GetUser(ctx, UserID, userEmail, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4319}).Error(err)
			return nil, err
		}
		db := c.Db
		tx, err := db.Begin()
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

		Cat, err := c.InsertCategory(ctx, tx, cat, userEmail, requestID)

		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4321}).Error(err)
			err = tx.Rollback()
			return nil, err
		}

		err = tx.Commit()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4322}).Error(err)
			err = tx.Rollback()
			return nil, err
		}
		return Cat, nil
	}
}

// InsertCategory - Insert category details into database
func (c *CategoryService) InsertCategory(ctx context.Context, tx *sql.Tx, cat Category, userEmail string, requestID string) (*Category, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4323}).Error(err)
		return nil, err
	default:
		stmt, err := tx.PrepareContext(ctx, `insert into categories
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
			err = stmt.Close()
			return nil, err
		}
		res, err := stmt.ExecContext(ctx,
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
			err = stmt.Close()
			return nil, err
		}

		uID, err := res.LastInsertId()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4326}).Error(err)
			err = stmt.Close()
			return nil, err
		}
		cat.ID = uint(uID)
		uUID4Str, err := common.UUIDBytesToStr(cat.UUID4)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 43257}).Error(err)
			return nil, err
		}
		cat.IDS = uUID4Str
		err = stmt.Close()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4328}).Error(err)
			return nil, err
		}
		return &cat, nil
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
		uUID4byte, err := common.UUIDStrToBytes(ID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4330}).Error(err)
			return nil, err
		}
		db := c.Db
		cat := &Category{}
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
			v.updated_year from categories c inner join topics v on (c.id = v.category_id) where c.uuid4 = ?`, uUID4byte)

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
				uUID4Str1, err := common.UUIDBytesToStr(cat.UUID4)
				if err != nil {
					log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4335}).Error(err)
					return nil, err
				}
				cat.IDS = uUID4Str1

				uUID4Str, err := common.UUIDBytesToStr(topc.UUID4)
				if err != nil {
					log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4336}).Error(err)
					return nil, err
				}
				topc.IDS = uUID4Str
				cat.Topics = append(cat.Topics, &topc)
			}

			err = rows.Close()
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4337}).Error(err)
				return nil, err
			}
		} else {
			cat = ctegry
		}
		return cat, nil
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
		userserv := &userservices.UserService{Config: c.Config, Db: c.Db, RedisClient: c.RedisClient}
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

		db := c.Db
		tx, err := db.Begin()

		tn, tnday, tnweek, tnmonth, tnyear := common.GetTimeDetails()
		cat := Category{}
		cat.UUID4, err = common.GetUUIDBytes()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4341}).Error(err)
			err = tx.Rollback()
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

		Cat, err := c.InsertCategory(ctx, tx, cat, userEmail, requestID)

		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4342}).Error(err)
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
		catchd.CreatedDay = tnday
		catchd.CreatedWeek = tnweek
		catchd.CreatedMonth = tnmonth
		catchd.CreatedYear = tnyear
		catchd.UpdatedDay = tnday
		catchd.UpdatedWeek = tnweek
		catchd.UpdatedMonth = tnmonth
		catchd.UpdatedYear = tnyear

		_, err = c.InsertChild(ctx, tx, catchd, userEmail, requestID)

		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4343}).Error(err)
			err = tx.Rollback()
			return nil, err
		}

		UpdatedDay := tnday
		UpdatedWeek := tnweek
		UpdatedMonth := tnmonth
		UpdatedYear := tnyear

		stmt, err := tx.PrepareContext(ctx, `update categories set 
				  num_chd = ?,
				  updated_at = ?, 
					updated_day = ?, 
					updated_week = ?, 
					updated_month = ?, 
					updated_year = ? where id = ?;`)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4344}).Error(err)
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
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4345}).Error(err)
			err = stmt.Close()
			err = tx.Rollback()
			return nil, err
		}

		err = stmt.Close()

		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4346}).Error(err)
			err = tx.Rollback()
			return nil, err
		}

		err = tx.Commit()
		return Cat, nil
	}
}

// InsertChild - Insert child category details into database
func (c *CategoryService) InsertChild(ctx context.Context, tx *sql.Tx, catchd CategoryChd, userEmail string, requestID string) (*CategoryChd, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4347}).Error(err)
		return nil, err
	default:
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
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4348}).Error(err)
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
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4349}).Error(err)
			err = stmt.Close()
			return nil, err
		}

		uID, err := res.LastInsertId()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4350}).Error(err)
			err = stmt.Close()
			return nil, err
		}
		catchd.ID = uint(uID)
		err = stmt.Close()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4351}).Error(err)
			return nil, err
		}
		return &catchd, nil
	}
}

// GetTopLevelCategories - Get top level categories
func (c *CategoryService) GetTopLevelCategories(ctx context.Context, userEmail string, requestID string) ([]*Category, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4352}).Error(err)
		return nil, err
	default:
		cats := []*Category{}
		rows, err := c.Db.QueryContext(ctx, `select 
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
			updated_year from categories where levelc = ? and statusc = ?;`, 0, 1)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4353}).Error(err)
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
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4354}).Error(err)
				return nil, err
			}
			uUID4Str, err := common.UUIDBytesToStr(cat.UUID4)
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4355}).Error(err)
				return nil, err
			}
			cat.IDS = uUID4Str
			cats = append(cats, &cat)
		}
		err = rows.Close()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4356}).Error(err)
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
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4357}).Error(err)
		return nil, err
	default:
		category, err := c.GetCategory(ctx, ID, userEmail, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4358}).Error(err)
			return nil, err
		}
		pohs := []*Category{}
		rows, err := c.Db.QueryContext(ctx, `select 
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
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4359}).Error(err)
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
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4360}).Error(err)
				return nil, err
			}
			uUID4Str, err := common.UUIDBytesToStr(cat.UUID4)
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4361}).Error(err)
				return nil, err
			}
			cat.IDS = uUID4Str

			pohs = append(pohs, &cat)
		}
		err = rows.Close()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4362}).Error(err)
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
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4363}).Error(err)
		return nil, err
	default:
		category, err := c.GetCategory(ctx, ID, userEmail, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4364}).Error(err)
			return nil, err
		}
		cat := Category{}
		row := c.Db.QueryRowContext(ctx, `select
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
			updated_year from categories where id = ?;`, category.ParentID)

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
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4365}).Error(err)
			return nil, err
		}
		uUID4Str, err := common.UUIDBytesToStr(cat.UUID4)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 4366}).Error(err)
			return nil, err
		}
		cat.IDS = uUID4Str

		return &cat, nil
	}
}
