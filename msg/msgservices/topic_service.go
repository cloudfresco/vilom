package msgservices

import (
	"context"
	"database/sql"
	"errors"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/go-redis/redis"
	"github.com/palantir/stacktrace"

	"github.com/cloudfresco/vilom/common"
	"github.com/cloudfresco/vilom/user/userservices"
)

// Topic - Topic view representation
type Topic struct {
	ID uint

	IDS       string
	TopicName string
	TopicDesc string
	NumTags   uint
	Tag1      string
	Tag2      string
	Tag3      string
	Tag4      string
	Tag5      string
	Tag6      string
	Tag7      string
	Tag8      string
	Tag9      string
	Tag10     string

	NumViews    uint `sql:"default:'0'"`
	NumMessages uint `sql:"default:'0'"`

	CategoryID uint
	UserID     uint
	UgroupID   uint

	common.StatusDates
	Messages []*Message

	//only for logic purpose to create message with topic
	Mtext   string
	Mattach string
}

// TopicsUser - TopicsUser view representation
type TopicsUser struct {
	ID          uint
	IDS         string
	TopicID     uint
	NumMessages uint `sql:"default:'0'"`
	NumViews    uint `sql:"default:'0'"`
	UserID      uint
	UgroupID    uint

	common.StatusDates
}

// UserTopic - UserTopic view representation
type UserTopic struct {
	ID       uint
	TopicID  uint
	UserID   uint
	UgroupID uint `sql:"default:'0'"`

	common.StatusDates
}

// TopicService - For accessing topic services
type TopicService struct {
	Config       *common.RedisOptions
	Db           *sql.DB
	RedisClient  *redis.Client
	LimitDefault string
}

// NewTopicService - Create topic service
func NewTopicService(config *common.RedisOptions,
	db *sql.DB,
	redisClient *redis.Client,
	limitDefault string) *TopicService {
	return &TopicService{config, db, redisClient, limitDefault}
}

// Show - Get topic details
func (t *TopicService) Show(ctx context.Context, ID string, UserID string) (*Topic, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		return nil, err
	default:
		db := t.Db
		topic, err := t.GetTopicWithMessages(ctx, ID)
		if err != nil {
			log.Error(stacktrace.Propagate(err, ""))
			return nil, err
		}
		//update topic_users table
		userserv := &userservices.UserService{Config: t.Config, Db: t.Db, RedisClient: t.RedisClient}
		user, err := userserv.GetUser(ctx, UserID)
		if err != nil {
			log.Error(stacktrace.Propagate(err, ""))
			return nil, err
		}
		var isPresent bool
		row := db.QueryRowContext(ctx, `select exists (select 1 from topics_users where topic_id = ? and user_id = ?);`, topic.ID)
		err = row.Scan(&isPresent)
		if err != nil {
			log.Error(stacktrace.Propagate(err, ""))
			return nil, err
		}

		tn := time.Now().UTC()
		_, week := tn.ISOWeek()
		day := tn.YearDay()

		tx, err := db.Begin()

		if isPresent {
			//update
			topicsuser, err := t.GetTopicsUser(ctx, topic.ID, user.ID)

			UpdatedDay := uint(day)
			UpdatedWeek := uint(week)
			UpdatedMonth := uint(tn.Month())
			UpdatedYear := uint(tn.Year())

			stmt, err := tx.PrepareContext(ctx, `update topics_users set 
					num_messages = ?,
          num_views = ?,
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
				topic.NumMessages,
				topicsuser.NumViews+1,
				tn,
				UpdatedDay,
				UpdatedWeek,
				UpdatedMonth,
				UpdatedYear,
				topicsuser.ID)
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

		} else {
			//create
			tu := TopicsUser{}
			tu.IDS = common.GetUID()
			tu.TopicID = topic.ID
			tu.NumMessages = topic.NumMessages
			tu.NumViews = 1
			tu.UserID = user.ID
			tu.UgroupID = uint(0)
			tu.Statusc = common.Active
			tu.CreatedAt = tn.UTC()
			tu.UpdatedAt = tn.UTC()
			tu.CreatedDay = uint(day)
			tu.CreatedWeek = uint(week)
			tu.CreatedMonth = uint(tn.Month())
			tu.CreatedYear = uint(tn.Year())
			tu.UpdatedDay = uint(day)
			tu.UpdatedWeek = uint(week)
			tu.UpdatedMonth = uint(tn.Month())
			tu.UpdatedYear = uint(tn.Year())

			_, err := t.InsertTopicsUser(ctx, tx, tu)

			if err != nil {
				log.Error(stacktrace.Propagate(err, ""))
				err = tx.Rollback()
				return nil, err
			}

		}
		err = tx.Commit()
		if err != nil {
			log.Error(stacktrace.Propagate(err, ""))
			return nil, err
		}
		return topic, nil
	}
}

// GetTopicWithMessages - Get topic with messages
func (t *TopicService) GetTopicWithMessages(ctx context.Context, ID string) (*Topic, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		return nil, err
	default:
		db := t.Db
		poh := &Topic{}

		tpc, err := t.GetTopic(ctx, ID)
		if err != nil {
			log.Error(stacktrace.Propagate(err, ""))
			return nil, err
		}
		var isPresent bool
		row := db.QueryRowContext(ctx, `select exists (select 1 from messages where topic_id = ?);`, tpc.ID)
		err = row.Scan(&isPresent)
		if err != nil {
			log.Error(stacktrace.Propagate(err, ""))
			return nil, err
		}
		if isPresent {

			rows, err := db.QueryContext(ctx, `select 
      p.id,
			p.id_s,
			p.topic_name,
			p.topic_desc,
			p.num_tags,
			p.tag1,
			p.tag2,
			p.tag3,
			p.tag4,
			p.tag5,
			p.tag6,
			p.tag7,
			p.tag8,
			p.tag9,
			p.tag10,
			p.num_views,
			p.num_messages,
			p.category_id,
			p.ugroup_id,
			p.user_id,
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
		  m.id,
			m.id_s,
			m.num_likes,
			m.num_upvotes,
			m.num_downvotes,
			m.category_id,
			m.topic_id,
			m.ugroup_id,
			m.user_id,
			m.statusc,
			m.created_at,
			m.updated_at,
			m.created_day,
			m.created_week,
			m.created_month,
			m.created_year,
			m.updated_day,
			m.updated_week,
			m.updated_month,
			m.updated_year from topics p inner join messages m on (p.id = m.topic_id) where p.id_s = ?`, ID)

			if err != nil {
				log.Error(stacktrace.Propagate(err, ""))
				return nil, err
			}
			for rows.Next() {
				msg := Message{}
				err = rows.Scan(
					&poh.ID,
					&poh.IDS,
					&poh.TopicName,
					&poh.TopicDesc,
					&poh.NumTags,
					&poh.Tag1,
					&poh.Tag2,
					&poh.Tag3,
					&poh.Tag4,
					&poh.Tag5,
					&poh.Tag6,
					&poh.Tag7,
					&poh.Tag8,
					&poh.Tag9,
					&poh.Tag10,
					&poh.NumViews,
					&poh.NumMessages,
					&poh.CategoryID,
					&poh.UgroupID,
					&poh.UserID,
					/*  StatusDates  */
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
					&msg.ID,
					&msg.IDS,
					&msg.NumLikes,
					&msg.NumUpvotes,
					&msg.NumDownvotes,
					&msg.CategoryID,
					&msg.TopicID,
					&msg.UgroupID,
					&msg.UserID,
					/*  StatusDates  */
					&msg.Statusc,
					&msg.CreatedAt,
					&msg.UpdatedAt,
					&msg.CreatedDay,
					&msg.CreatedWeek,
					&msg.CreatedMonth,
					&msg.CreatedYear,
					&msg.UpdatedDay,
					&msg.UpdatedWeek,
					&msg.UpdatedMonth,
					&msg.UpdatedYear)

				if err != nil {
					log.Error(stacktrace.Propagate(err, ""))
					return nil, err
				}

				poh.Messages = append(poh.Messages, &msg)
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

		} else {
			poh = tpc
		}

		if len(poh.Messages) > 0 {
			msgserv := &MessageService{t.Config, t.Db, t.RedisClient, t.LimitDefault}
			Messages, err := msgserv.GetMessagesWithTextAttach(ctx, poh.Messages)
			if err != nil {
				log.Error(stacktrace.Propagate(err, ""))
			}
			poh.Messages = Messages
		}
		return poh, nil
	}
}

// Create - Create topic
func (t *TopicService) Create(ctx context.Context, form *Topic, UserID string) (*Topic, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		return nil, err
	default:
		userserv := &userservices.UserService{Config: t.Config, Db: t.Db, RedisClient: t.RedisClient}
		user, err := userserv.GetUser(ctx, UserID)
		if err != nil {
			log.Error(stacktrace.Propagate(err, ""))
			return nil, err
		}
		db := t.Db
		tx, err := db.Begin()
		if err != nil {
			log.Error(stacktrace.Propagate(err, ""))
			return nil, err
		}

		tn := time.Now().UTC()
		_, week := tn.ISOWeek()
		day := tn.YearDay()
		topc := Topic{}
		topc.IDS = common.GetUID()
		topc.TopicName = form.TopicName
		topc.TopicDesc = form.TopicDesc
		topc.NumTags = form.NumTags
		topc.Tag1 = form.Tag1
		topc.Tag2 = form.Tag2
		topc.Tag3 = form.Tag3
		topc.Tag4 = form.Tag4
		topc.Tag5 = form.Tag5
		topc.Tag6 = form.Tag6
		topc.Tag7 = form.Tag7
		topc.Tag8 = form.Tag8
		topc.Tag9 = form.Tag9
		topc.Tag10 = form.Tag10
		topc.NumViews = uint(0)
		topc.NumMessages = uint(0)
		topc.CategoryID = form.CategoryID
		topc.UserID = user.ID
		topc.UgroupID = form.UgroupID
		/*  StatusDates  */
		topc.Statusc = common.Active
		topc.CreatedAt = tn.UTC()
		topc.UpdatedAt = tn.UTC()
		topc.CreatedDay = uint(day)
		topc.CreatedWeek = uint(week)
		topc.CreatedMonth = uint(tn.Month())
		topc.CreatedYear = uint(tn.Year())
		topc.UpdatedDay = uint(day)
		topc.UpdatedWeek = uint(week)
		topc.UpdatedMonth = uint(tn.Month())
		topc.UpdatedYear = uint(tn.Year())

		topic, err := t.InsertTopic(ctx, tx, topc)

		if err != nil {
			log.Error(stacktrace.Propagate(err, ""))
			err = tx.Rollback()
			return nil, err
		}

		UpdatedDay := uint(day)
		UpdatedWeek := uint(week)
		UpdatedMonth := uint(tn.Month())
		UpdatedYear := uint(tn.Year())

		catserv := &CategoryService{t.Config, t.Db, t.RedisClient, t.LimitDefault}
		category, err := catserv.GetCategoryByID(ctx, form.CategoryID)
		if err != nil {
			log.Error(stacktrace.Propagate(err, ""))
			return nil, err
		}
		//update category count
		stmt, err := tx.PrepareContext(ctx, `update categories set 
    num_topics = ?,
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
			category.NumTopics+1,
			tn,
			UpdatedDay,
			UpdatedWeek,
			UpdatedMonth,
			UpdatedYear, category.ID)
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

		msgserv := &MessageService{t.Config, t.Db, t.RedisClient, t.LimitDefault}

		if form.Mtext != "" {
			//insert message
			msg := Message{}
			msg.IDS = common.GetUID()
			msg.NumLikes = 0
			msg.NumUpvotes = 0
			msg.NumDownvotes = 0
			msg.CategoryID = category.ID
			msg.TopicID = topic.ID
			msg.UserID = user.ID
			msg.UgroupID = form.UgroupID
			/*  StatusDates  */
			msg.Statusc = common.Active
			msg.CreatedAt = tn.UTC()
			msg.UpdatedAt = tn.UTC()
			msg.CreatedDay = uint(day)
			msg.CreatedWeek = uint(week)
			msg.CreatedMonth = uint(tn.Month())
			msg.CreatedYear = uint(tn.Year())
			msg.UpdatedDay = uint(day)
			msg.UpdatedWeek = uint(week)
			msg.UpdatedMonth = uint(tn.Month())
			msg.UpdatedYear = uint(tn.Year())
			Message, err := msgserv.InsertMessage(ctx, tx, msg)

			if err != nil {
				log.Error(stacktrace.Propagate(err, ""))
				err = tx.Rollback()
				return nil, err
			}
			//insert message_text
			msgtxt := MessageText{}
			msgtxt.Mtext = form.Mtext
			msgtxt.CategoryID = category.ID
			msgtxt.TopicID = topic.ID
			msgtxt.MessageID = Message.ID
			msgtxt.UserID = user.ID
			msgtxt.UgroupID = form.UgroupID
			/*  StatusDates  */
			msgtxt.Statusc = common.Active
			msgtxt.CreatedAt = tn.UTC()
			msgtxt.UpdatedAt = tn.UTC()
			msgtxt.CreatedDay = uint(day)
			msgtxt.CreatedWeek = uint(week)
			msgtxt.CreatedMonth = uint(tn.Month())
			msgtxt.CreatedYear = uint(tn.Year())
			msgtxt.UpdatedDay = uint(day)
			msgtxt.UpdatedWeek = uint(week)
			msgtxt.UpdatedMonth = uint(tn.Month())
			msgtxt.UpdatedYear = uint(tn.Year())

			_, err = msgserv.InsertMessageText(ctx, tx, msgtxt)

			if err != nil {
				log.Error(stacktrace.Propagate(err, ""))
				err = tx.Rollback()
				return nil, err
			}

			//update messages count in topic table
			stmt, err := tx.PrepareContext(ctx, `update topics set 
		  num_messages = ?,
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
				topic.NumMessages+1,
				tn,
				UpdatedDay,
				UpdatedWeek,
				UpdatedMonth,
				UpdatedYear, topic.ID)
			if err != nil {
				log.Error(stacktrace.Propagate(err, ""))
				err = stmt.Close()
				err = tx.Rollback()
				return nil, err
			}
			err = stmt.Close()
			if err != nil {
				log.Error(stacktrace.Propagate(err, ""))
				return nil, err
			}

			//insert atatchment
			if form.Mattach != "" {
				msgatch := MessageAttachment{}
				msgatch.Mattach = form.Mattach
				msgatch.CategoryID = category.ID
				msgatch.TopicID = topic.ID
				msgatch.MessageID = Message.ID
				msgatch.UserID = user.ID
				msgatch.UgroupID = form.UgroupID
				/*  StatusDates  */
				msgatch.Statusc = common.Active
				msgatch.CreatedAt = tn.UTC()
				msgatch.UpdatedAt = tn.UTC()
				msgatch.CreatedDay = uint(day)
				msgatch.CreatedWeek = uint(week)
				msgatch.CreatedMonth = uint(tn.Month())
				msgatch.CreatedYear = uint(tn.Year())
				msgatch.UpdatedDay = uint(day)
				msgatch.UpdatedWeek = uint(week)
				msgatch.UpdatedMonth = uint(tn.Month())
				msgatch.UpdatedYear = uint(tn.Year())

				_, err := msgserv.InsertMessageAttachment(ctx, tx, msgatch)

				if err != nil {
					log.Error(stacktrace.Propagate(err, ""))
					err = tx.Rollback()
					return nil, err
				}
			}
		}

		//insert user_topics
		ut := UserTopic{}
		ut.TopicID = topic.ID
		ut.UserID = user.ID
		ut.UgroupID = uint(0)
		/*  StatusDates  */
		ut.Statusc = common.Active
		ut.CreatedAt = tn.UTC()
		ut.UpdatedAt = tn.UTC()
		ut.CreatedDay = uint(day)
		ut.CreatedWeek = uint(week)
		ut.CreatedMonth = uint(tn.Month())
		ut.CreatedYear = uint(tn.Year())
		ut.UpdatedDay = uint(day)
		ut.UpdatedWeek = uint(week)
		ut.UpdatedMonth = uint(tn.Month())
		ut.UpdatedYear = uint(tn.Year())

		_, err = t.InsertUserTopic(ctx, tx, ut)

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
		return topic, nil
	}
}

// InsertTopic - Insert topic details into database
func (t *TopicService) InsertTopic(ctx context.Context, tx *sql.Tx, topc Topic) (*Topic, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		return nil, err
	default:
		stmt, err := tx.PrepareContext(ctx, `insert into topics
	  ( id_s,
			topic_name,
			topic_desc,
			num_tags,
			tag1,
			tag2,
			tag3,
			tag4,
			tag5,
			tag6,
			tag7,
			tag8,
			tag9,
			tag10,
			num_views,
			num_messages,
			category_id,
			user_id,
			ugroup_id,
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
					?,?,?,?,?,?,?,?,?,?);`)
		if err != nil {
			log.Error(stacktrace.Propagate(err, ""))
			err = stmt.Close()
			return nil, err
		}
		res, err := stmt.ExecContext(ctx,
			topc.IDS,
			topc.TopicName,
			topc.TopicDesc,
			topc.NumTags,
			topc.Tag1,
			topc.Tag2,
			topc.Tag3,
			topc.Tag4,
			topc.Tag5,
			topc.Tag6,
			topc.Tag7,
			topc.Tag8,
			topc.Tag9,
			topc.Tag10,
			topc.NumViews,
			topc.NumMessages,
			topc.CategoryID,
			topc.UserID,
			topc.UgroupID,
			/*  StatusDates  */
			topc.Statusc,
			topc.CreatedAt,
			topc.UpdatedAt,
			topc.CreatedDay,
			topc.CreatedWeek,
			topc.CreatedMonth,
			topc.CreatedYear,
			topc.UpdatedDay,
			topc.UpdatedWeek,
			topc.UpdatedMonth,
			topc.UpdatedYear)

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
		topc.ID = uint(uID)
		err = stmt.Close()
		if err != nil {
			log.Error(stacktrace.Propagate(err, ""))
			return nil, err
		}
		return &topc, nil
	}
}

// GetTopicByID - Get topic by ID
func (t *TopicService) GetTopicByID(ctx context.Context, ID uint) (*Topic, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		return nil, err
	default:
		poh := Topic{}
		row := t.Db.QueryRowContext(ctx, `select
    id,
		id_s,
		topic_name,
		topic_desc,
		num_tags,
		tag1,
		tag2,
		tag3,
		tag4,
		tag5,
		tag6,
		tag7,
		tag8,
		tag9,
		tag10,
		num_views,
		num_messages,
		category_id,
		user_id,
		ugroup_id,
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
		updated_year from topics where id = ?`, ID)

		err := row.Scan(
			&poh.ID,
			&poh.IDS,
			&poh.TopicName,
			&poh.TopicDesc,
			&poh.NumTags,
			&poh.Tag1,
			&poh.Tag2,
			&poh.Tag3,
			&poh.Tag4,
			&poh.Tag5,
			&poh.Tag6,
			&poh.Tag7,
			&poh.Tag8,
			&poh.Tag9,
			&poh.Tag10,
			&poh.NumViews,
			&poh.NumMessages,
			&poh.CategoryID,
			&poh.UserID,
			&poh.UgroupID,
			/*  StatusDates  */
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
			&poh.UpdatedYear)

		if err != nil {
			log.Error(stacktrace.Propagate(err, ""))
			return nil, err
		}

		return &poh, nil
	}
}

// GetTopic - Get topic
func (t *TopicService) GetTopic(ctx context.Context, ID string) (*Topic, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		return nil, err
	default:
		poh := Topic{}
		row := t.Db.QueryRowContext(ctx, `select
    id,
		id_s,
		topic_name,
		topic_desc,
		num_tags,
		tag1,
		tag2,
		tag3,
		tag4,
		tag5,
		tag6,
		tag7,
		tag8,
		tag9,
		tag10,
		num_views,
		num_messages,
		category_id,
		user_id,
		ugroup_id,
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
		updated_year from topics where id_s = ?`, ID)

		err := row.Scan(
			&poh.ID,
			&poh.IDS,
			&poh.TopicName,
			&poh.TopicDesc,
			&poh.NumTags,
			&poh.Tag1,
			&poh.Tag2,
			&poh.Tag3,
			&poh.Tag4,
			&poh.Tag5,
			&poh.Tag6,
			&poh.Tag7,
			&poh.Tag8,
			&poh.Tag9,
			&poh.Tag10,
			&poh.NumViews,
			&poh.NumMessages,
			&poh.CategoryID,
			&poh.UserID,
			&poh.UgroupID,
			/*  StatusDates  */
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
			&poh.UpdatedYear)

		if err != nil {
			log.Error(stacktrace.Propagate(err, ""))
			return nil, err
		}

		return &poh, nil
	}
}

// GetTopicByName - Get topic by name
func (t *TopicService) GetTopicByName(ctx context.Context, topicname string) (*Topic, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		return nil, err
	default:
		poh := Topic{}
		row := t.Db.QueryRowContext(ctx, `select
    id,
		id_s,
		topic_name,
		topic_desc,
		num_tags,
		tag1,
		tag2,
		tag3,
		tag4,
		tag5,
		tag6,
		tag7,
		tag8,
		tag9,
		tag10,
		num_views,
		num_messages,
		category_id,
		user_id,
		ugroup_id,
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
		updated_year from topics where topic_name = ?`, topicname)

		err := row.Scan(
			&poh.ID,
			&poh.IDS,
			&poh.TopicName,
			&poh.TopicDesc,
			&poh.NumTags,
			&poh.Tag1,
			&poh.Tag2,
			&poh.Tag3,
			&poh.Tag4,
			&poh.Tag5,
			&poh.Tag6,
			&poh.Tag7,
			&poh.Tag8,
			&poh.Tag9,
			&poh.Tag10,
			&poh.NumViews,
			&poh.NumMessages,
			&poh.CategoryID,
			&poh.UserID,
			&poh.UgroupID,
			/*  StatusDates  */
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
			&poh.UpdatedYear)

		if err != nil {
			log.Error(stacktrace.Propagate(err, ""))
			return nil, err
		}

		return &poh, nil
	}
}

// GetTopicsUser - Get user topics
func (t *TopicService) GetTopicsUser(ctx context.Context, ID uint, UserID uint) (*TopicsUser, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		return nil, err
	default:
		poh := TopicsUser{}
		row := t.Db.QueryRowContext(ctx, `select
    id,
		id_s,
		topic_id,
		num_messages,
    num_views,
		user_id,
		ugroup_id,
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
		updated_year from topics_users where topic_id = ? and user_id = ?`, ID, UserID)

		err := row.Scan(
			&poh.ID,
			&poh.IDS,
			&poh.TopicID,
			&poh.NumMessages,
			&poh.NumViews,
			&poh.UserID,
			&poh.UgroupID,
			/*  StatusDates  */
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
			&poh.UpdatedYear)

		if err != nil {
			log.Error(stacktrace.Propagate(err, ""))
			return nil, err
		}

		return &poh, nil
	}
}

// InsertTopicsUser - Insert topic user details into database
func (t *TopicService) InsertTopicsUser(ctx context.Context, tx *sql.Tx, poh TopicsUser) (*TopicsUser, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		return nil, err
	default:
		stmt, err := tx.PrepareContext(ctx, `insert into topics_users
	  (id_s,
		topic_id,
		num_messages,
    num_views,
		user_id,
		ugroup_id,
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
			err = stmt.Close()
			return nil, err
		}
		_, err = stmt.ExecContext(ctx,
			poh.IDS,
			poh.TopicID,
			poh.NumMessages,
			poh.NumViews,
			poh.UserID,
			poh.UgroupID,
			poh.Statusc,
			poh.CreatedAt,
			poh.UpdatedAt,
			poh.CreatedDay,
			poh.CreatedWeek,
			poh.CreatedMonth,
			poh.CreatedYear,
			poh.UpdatedDay,
			poh.UpdatedWeek,
			poh.UpdatedMonth,
			poh.UpdatedYear)

		if err != nil {
			log.Error(stacktrace.Propagate(err, ""))
			err = stmt.Close()
			return nil, err
		}

		err = stmt.Close()

		if err != nil {
			log.Error(stacktrace.Propagate(err, ""))
			return nil, err
		}

		return &poh, nil
	}
}

// InsertUserTopic - Insert user topics details into database
func (t *TopicService) InsertUserTopic(ctx context.Context, tx *sql.Tx, poh UserTopic) (*UserTopic, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		return nil, err
	default:
		stmt, err := tx.PrepareContext(ctx, `insert into user_topics
	  (
		topic_id,
		user_id,
		ugroup_id,
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
		_, err = stmt.ExecContext(ctx,
			poh.TopicID,
			poh.UserID,
			poh.UgroupID,
			poh.Statusc,
			poh.CreatedAt,
			poh.UpdatedAt,
			poh.CreatedDay,
			poh.CreatedWeek,
			poh.CreatedMonth,
			poh.CreatedYear,
			poh.UpdatedDay,
			poh.UpdatedWeek,
			poh.UpdatedMonth,
			poh.UpdatedYear)

		if err != nil {
			log.Error(stacktrace.Propagate(err, ""))
			err = stmt.Close()
			return nil, err
		}

		err = stmt.Close()
		if err != nil {
			log.Error(stacktrace.Propagate(err, ""))
			return nil, err
		}

		return &poh, nil
	}
}
