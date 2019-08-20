package msgservices

import (
	"context"
	"database/sql"
	"errors"

	log "github.com/sirupsen/logrus"

	"github.com/cloudfresco/vilom/common"
	"github.com/cloudfresco/vilom/user/userservices"
)

/* error message range: 5300-5999 */

// Topic - Topic view representation
type Topic struct {
	ID    uint   `json:"id,omitempty"`
	UUID4 []byte `json:"-"`
	IDS   string `json:"id_s,omitempty"`

	TopicName string `json:"topic_name,omitempty"`
	TopicDesc string `json:"topic_desc,omitempty"`
	NumTags   uint   `json:"num_tags,omitempty"`
	Tag1      string `json:"tag1,omitempty"`
	Tag2      string `json:"tag2,omitempty"`
	Tag3      string `json:"tag3,omitempty"`
	Tag4      string `json:"tag4,omitempty"`
	Tag5      string `json:"tag5,omitempty"`
	Tag6      string `json:"tag6,omitempty"`
	Tag7      string `json:"tag7,omitempty"`
	Tag8      string `json:"tag8,omitempty"`
	Tag9      string `json:"tag9,omitempty"`
	Tag10     string `json:"tag10,omitempty"`

	NumViews    uint `json:"num_views,omitempty"`
	NumMessages uint `json:"num_messages,omitempty"`

	CategoryID uint `json:"category_id,omitempty"`
	UserID     uint `json:"user_id,omitempty"`
	UgroupID   uint `json:"ugroup_id,omitempty"`

	common.StatusDates
	Messages []*Message

	//only for logic purpose to create message with topic
	Mtext   string `json:"-"`
	Mattach string `json:"-"`
}

// TopicsUser - TopicsUser view representation
type TopicsUser struct {
	ID          uint   `json:"id,omitempty"`
	UUID4       []byte `json:"-"`
	IDS         string `json:"id_s,omitempty"`
	TopicID     uint   `json:"topic_id,omitempty"`
	NumMessages uint   `json:"num_messages,omitempty"`
	NumViews    uint   `json:"num_views,omitempty"`
	UserID      uint   `json:"user_id,omitempty"`
	UgroupID    uint   `json:"ugroup_id,omitempty"`

	common.StatusDates
}

// UserTopic - UserTopic view representation
type UserTopic struct {
	ID       uint   `json:"id,omitempty"`
	UUID4    []byte `json:"-"`
	TopicID  uint   `json:"topic_id,omitempty"`
	UserID   uint   `json:"user_id,omitempty"`
	UgroupID uint   `json:"ugroup_id,omitempty"`

	common.StatusDates
}

// TopicServiceIntf - interface for Topic Service
type TopicServiceIntf interface {
	CreateTopic(ctx context.Context, form *Topic, UserID string, userEmail string, requestID string) (*Topic, error)
	ShowTopic(ctx context.Context, ID string, UserID string, userEmail string, requestID string) (*Topic, error)
	GetTopicByID(ctx context.Context, ID uint, userEmail string, requestID string) (*Topic, error)
	GetTopic(ctx context.Context, ID string, userEmail string, requestID string) (*Topic, error)
	GetTopicByName(ctx context.Context, topicname string, userEmail string, requestID string) (*Topic, error)
	GetTopicWithMessages(ctx context.Context, ID string, userEmail string, requestID string) (*Topic, error)
	GetTopicMessages(ctx context.Context, uuid4byte []byte, userEmail string, requestID string) (*Topic, error)
	GetTopicsUser(ctx context.Context, ID uint, UserID uint, userEmail string, requestID string) (*TopicsUser, error)
	UpdateTopic(ctx context.Context, ID string, form *Topic, UserID string, userEmail string, requestID string) error
	UpdateNumMessagesPrepare(ctx context.Context, userEmail string, requestID string) (*sql.Stmt, error)
	UpdateNumMessages(ctx context.Context, stmt *sql.Stmt, tx *sql.Tx, numMessages uint, ID uint, userEmail string, requestID string) error
	DeleteTopic(ctx context.Context, ID string, userEmail string, requestID string) error
}

// TopicService - For accessing topic services
type TopicService struct {
	DBService    *common.DBService
	RedisService *common.RedisService
}

// NewTopicService - Create topic service
func NewTopicService(dbOpt *common.DBService, redisOpt *common.RedisService) *TopicService {
	return &TopicService{
		DBService:    dbOpt,
		RedisService: redisOpt,
	}
}

// CreateTopic - Create topic
func (t *TopicService) CreateTopic(ctx context.Context, form *Topic, UserID string, userEmail string, requestID string) (*Topic, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5321}).Error(err)
		return nil, err
	default:
		userserv := &userservices.UserService{DBService: t.DBService, RedisService: t.RedisService}
		user, err := userserv.GetUser(ctx, UserID, userEmail, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5322}).Error(err)
			return nil, err
		}
		db := t.DBService.DB
		catserv := &CategoryService{DBService: t.DBService, RedisService: t.RedisService}
		category, err := catserv.GetCategoryByID(ctx, form.CategoryID, userEmail, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5326}).Error(err)
			return nil, err
		}

		insertTopicStmt, err := t.insertTopicPrepare(ctx, userEmail, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5323}).Error(err)
			return nil, err
		}
		updateNumTopicsStmt, err := catserv.UpdateNumTopicsPrepare(ctx, userEmail, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5323}).Error(err)
			return nil, err
		}
		insertUserTopicStmt, err := t.insertUserTopicPrepare(ctx, userEmail, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5323}).Error(err)
			return nil, err
		}
		tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5323}).Error(err)
			return nil, err
		}

		topic, err := t.createTopic(ctx, insertTopicStmt, tx, form, user.ID, userEmail, requestID)

		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5325}).Error(err)
			err = tx.Rollback()
			return nil, err
		}

		numtopics := category.NumTopics + 1
		err = catserv.UpdateNumTopics(ctx, updateNumTopicsStmt, tx, numtopics, category.ID, userEmail, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5381}).Error(err)
			err = tx.Rollback()
			return nil, err
		}
		if form.Mtext != "" {
			msgserv := &MessageService{DBService: t.DBService, RedisService: t.RedisService}
			msgform := Message{}
			msgform.CategoryID = category.ID
			msgform.TopicID = topic.ID
			msgform.UgroupID = form.UgroupID
			msgform.Mtext = form.Mtext
			msgform.Mattach = form.Mattach
			_, err = msgserv.CreateMessage(ctx, &msgform, UserID, false, userEmail, requestID)
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5330}).Error(err)
				err = tx.Rollback()
				return nil, err
			}
		}

		err = t.createUserTopic(ctx, insertUserTopicStmt, tx, topic.ID, user.ID, userEmail, requestID)

		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5340}).Error(err)
			err = tx.Rollback()
			return nil, err
		}

		err = tx.Commit()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5341}).Error(err)
			err = tx.Rollback()
			return nil, err
		}
		return topic, nil
	}
}

// create topic - create topic
func (t *TopicService) createTopic(ctx context.Context, stmt *sql.Stmt, tx *sql.Tx, form *Topic, userID uint, userEmail string, requestID string) (*Topic, error) {
	var err error
	tn, tnday, tnweek, tnmonth, tnyear := common.GetTimeDetails()

	topic := Topic{}
	topic.UUID4, err = common.GetUUIDBytes()
	if err != nil {
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5324}).Error(err)
		return nil, err
	}
	topic.TopicName = form.TopicName
	topic.TopicDesc = form.TopicDesc
	topic.NumTags = form.NumTags
	topic.Tag1 = form.Tag1
	topic.Tag2 = form.Tag2
	topic.Tag3 = form.Tag3
	topic.Tag4 = form.Tag4
	topic.Tag5 = form.Tag5
	topic.Tag6 = form.Tag6
	topic.Tag7 = form.Tag7
	topic.Tag8 = form.Tag8
	topic.Tag9 = form.Tag9
	topic.Tag10 = form.Tag10
	topic.NumViews = uint(0)
	topic.NumMessages = uint(0)
	topic.CategoryID = form.CategoryID
	topic.UserID = userID
	topic.UgroupID = form.UgroupID
	/*  StatusDates  */
	topic.Statusc = common.Active
	topic.CreatedAt = tn
	topic.UpdatedAt = tn
	topic.CreatedDay = tnday
	topic.CreatedWeek = tnweek
	topic.CreatedMonth = tnmonth
	topic.CreatedYear = tnyear
	topic.UpdatedDay = tnday
	topic.UpdatedWeek = tnweek
	topic.UpdatedMonth = tnmonth
	topic.UpdatedYear = tnyear

	err = t.insertTopic(ctx, stmt, tx, &topic, userEmail, requestID)

	if err != nil {
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5325}).Error(err)
		return nil, err
	}
	return &topic, nil
}

// insertTopicPrepare - Insert topic Prepare Statement
func (t *TopicService) insertTopicPrepare(ctx context.Context, userEmail string, requestID string) (*sql.Stmt, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5342}).Error(err)
		return nil, err
	default:
		db := t.DBService.DB
		stmt, err := db.PrepareContext(ctx, `insert into topics
	  ( uuid4,
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
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5343}).Error(err)
			return nil, err
		}
		return stmt, nil
	}
}

// insertTopic - Insert topic details into database
func (t *TopicService) insertTopic(ctx context.Context, stmt *sql.Stmt, tx *sql.Tx, topc *Topic, userEmail string, requestID string) error {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5342}).Error(err)
		return err
	default:
		res, err := tx.StmtContext(ctx, stmt).Exec(
			topc.UUID4,
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
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5344}).Error(err)
			return err
		}
		uID, err := res.LastInsertId()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5345}).Error(err)
			return err
		}
		topc.ID = uint(uID)
		uuid4Str, err := common.UUIDBytesToStr(topc.UUID4)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5346}).Error(err)
			return err
		}
		topc.IDS = uuid4Str
		return nil
	}
}

// createUserTopic - create user topic
func (t *TopicService) createUserTopic(ctx context.Context, stmt *sql.Stmt, tx *sql.Tx, topicID uint, userID uint, userEmail string, requestID string) error {
	var err error
	tn, tnday, tnweek, tnmonth, tnyear := common.GetTimeDetails()
	ut := UserTopic{}
	ut.UUID4, err = common.GetUUIDBytes()
	if err != nil {
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5339}).Error(err)
		return err
	}
	ut.TopicID = topicID
	ut.UserID = userID
	ut.UgroupID = uint(0)
	/*  StatusDates  */
	ut.Statusc = common.Active
	ut.CreatedAt = tn
	ut.UpdatedAt = tn
	ut.CreatedDay = tnday
	ut.CreatedWeek = tnweek
	ut.CreatedMonth = tnmonth
	ut.CreatedYear = tnyear
	ut.UpdatedDay = tnday
	ut.UpdatedWeek = tnweek
	ut.UpdatedMonth = tnmonth
	ut.UpdatedYear = tnyear

	err = t.insertUserTopic(ctx, stmt, tx, &ut, userEmail, requestID)

	if err != nil {
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5340}).Error(err)
		return err
	}
	return nil
}

// insertUserTopicPrepare - Insert user topics Prepare Statement
func (t *TopicService) insertUserTopicPrepare(ctx context.Context, userEmail string, requestID string) (*sql.Stmt, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5367}).Error(err)
		return nil, err
	default:
		db := t.DBService.DB
		stmt, err := db.PrepareContext(ctx, `insert into user_topics
	  (
    uuid4,
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
					?,?,?,?,?);`)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5368}).Error(err)
			return nil, err
		}
		return stmt, nil
	}
}

// insertUserTopic - Insert user topics details into database
func (t *TopicService) insertUserTopic(ctx context.Context, stmt *sql.Stmt, tx *sql.Tx, poh *UserTopic, userEmail string, requestID string) error {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5367}).Error(err)
		return err
	default:
		res, err := tx.StmtContext(ctx, stmt).Exec(
			poh.UUID4,
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
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5369}).Error(err)
			return err
		}
		uID, err := res.LastInsertId()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5370}).Error(err)
			return err
		}
		poh.ID = uint(uID)
		return nil
	}
}

// ShowTopic - Get topic details
func (t *TopicService) ShowTopic(ctx context.Context, ID string, UserID string, userEmail string, requestID string) (*Topic, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5300}).Error(err)
		return nil, err
	default:
		db := t.DBService.DB
		topic, err := t.GetTopicWithMessages(ctx, ID, userEmail, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5301}).Error(err)
			return nil, err
		}
		//update topic_users table
		userserv := &userservices.UserService{DBService: t.DBService, RedisService: t.RedisService}
		user, err := userserv.GetUser(ctx, UserID, userEmail, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5302}).Error(err)
			return nil, err
		}
		var isPresent bool
		row := db.QueryRowContext(ctx, `select exists (select 1 from topics_users where topic_id = ? and user_id = ?);`, topic.ID, user.ID)
		err = row.Scan(&isPresent)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5303}).Error(err)
			return nil, err
		}

		updateTopicUsersStmt, insertTopicsUserStmt, err := t.showTopicPrepareStmts(ctx, userEmail, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5379}).Error(err)
			return nil, err
		}

		tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5372}).Error(err)
			err = tx.Rollback()
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5379}).Error(err)
				return nil, err
			}
			err = t.showTopicPrepareStmtsClose(ctx, updateTopicUsersStmt, insertTopicsUserStmt, userEmail, requestID)
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5379}).Error(err)
				return nil, err
			}
			return nil, err
		}
		err = t.showTopicUpdateTopicUsers(ctx, updateTopicUsersStmt, insertTopicsUserStmt, topic, tx, user.ID, isPresent, userEmail, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5372}).Error(err)
			err = tx.Rollback()
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5379}).Error(err)
				return nil, err
			}
			err = t.showTopicPrepareStmtsClose(ctx, updateTopicUsersStmt, insertTopicsUserStmt, userEmail, requestID)
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5379}).Error(err)
				return nil, err
			}
			return nil, err
		}

		err = tx.Commit()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5309}).Error(err)
			return nil, err
		}
		err = t.showTopicPrepareStmtsClose(ctx, updateTopicUsersStmt, insertTopicsUserStmt, userEmail, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5379}).Error(err)
			return nil, err
		}
		return topic, nil
	}
}

//showTopicPrepareStmts - Prepare Statements
func (t *TopicService) showTopicPrepareStmts(ctx context.Context, userEmail string, requestID string) (*sql.Stmt, *sql.Stmt, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6311}).Error(err)
		return nil, nil, err
	default:
		updateTopicUsersStmt, err := t.updateTopicUsersPrepare(ctx, userEmail, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5379}).Error(err)
			return nil, nil, err
		}

		insertTopicsUserStmt, err := t.insertTopicsUserPrepare(ctx, userEmail, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5379}).Error(err)
			return nil, nil, err
		}

		return updateTopicUsersStmt, insertTopicsUserStmt, nil
	}
}

//showTopicPrepareStmtsClose - Close Prepare Statements
func (t *TopicService) showTopicPrepareStmtsClose(ctx context.Context, updateTopicUsersStmt *sql.Stmt, insertTopicsUserStmt *sql.Stmt, userEmail string, requestID string) error {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5379}).Error(err)
		return err
	default:
		err := updateTopicUsersStmt.Close()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5379}).Error(err)
			return err
		}

		err = insertTopicsUserStmt.Close()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5379}).Error(err)
			return err
		}

		return nil
	}
}

// showTopicUpdateTopicUsers - update topic users details
func (t *TopicService) showTopicUpdateTopicUsers(ctx context.Context, updateTopicUsersStmt *sql.Stmt, insertTopicsUserStmt *sql.Stmt, topic *Topic, tx *sql.Tx, userID uint, isPresent bool, userEmail string, requestID string) error {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5300}).Error(err)
		return err
	default:
		var err error
		tn, tnday, tnweek, tnmonth, tnyear := common.GetTimeDetails()
		if isPresent {
			//update
			topicsuser, err := t.GetTopicsUser(ctx, topic.ID, userID, userEmail, requestID)
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5373}).Error(err)
				return err
			}

			numViews := topicsuser.NumViews + 1
			err = t.updateTopicUsers(ctx, updateTopicUsersStmt, tx, topic.NumMessages, numViews, topicsuser.ID, userEmail, requestID)
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5379}).Error(err)
				return err
			}
		} else {
			//create
			tu := TopicsUser{}
			tu.UUID4, err = common.GetUUIDBytes()
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5307}).Error(err)
				return err
			}
			tu.TopicID = topic.ID
			tu.NumMessages = topic.NumMessages
			tu.NumViews = 1
			tu.UserID = userID
			tu.UgroupID = uint(0)
			tu.Statusc = common.Active
			tu.CreatedAt = tn
			tu.UpdatedAt = tn
			tu.CreatedDay = tnday
			tu.CreatedWeek = tnweek
			tu.CreatedMonth = tnmonth
			tu.CreatedYear = tnyear
			tu.UpdatedDay = tnday
			tu.UpdatedWeek = tnweek
			tu.UpdatedMonth = tnmonth
			tu.UpdatedYear = tnyear

			_, err := t.insertTopicsUser(ctx, insertTopicsUserStmt, tx, tu, userEmail, requestID)

			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5308}).Error(err)
				return err
			}

		}
		return nil
	}
}

// updateTopicUsersPrepare - update topic users prepare statement
func (t *TopicService) updateTopicUsersPrepare(ctx context.Context, userEmail string, requestID string) (*sql.Stmt, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5332}).Error(err)
		return nil, err
	default:
		db := t.DBService.DB
		stmt, err := db.PrepareContext(ctx, `update topics_users set 
					num_messages = ?,
          num_views = ?,
					updated_at = ?, 
					updated_day = ?, 
					updated_week = ?, 
					updated_month = ?, 
					updated_year = ? where id = ? and statusc = ?;`)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5304}).Error(err)
			return nil, err
		}
		return stmt, nil
	}
}

// updateTopicUsers - update topic users
func (t *TopicService) updateTopicUsers(ctx context.Context, stmt *sql.Stmt, tx *sql.Tx, numMessages uint, numViews uint, topicsuserID uint, userEmail string, requestID string) error {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5332}).Error(err)
		return err
	default:
		tn, tnday, tnweek, tnmonth, tnyear := common.GetTimeDetails()
		_, err := tx.StmtContext(ctx, stmt).Exec(
			numMessages,
			numViews,
			tn,
			tnday,
			tnweek,
			tnmonth,
			tnyear,
			topicsuserID,
			common.Active)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5305}).Error(err)
			return err
		}
		return nil
	}
}

// insertTopicsUserPrepare - Insert topic user Prepare Statement
func (t *TopicService) insertTopicsUserPrepare(ctx context.Context, userEmail string, requestID string) (*sql.Stmt, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5361}).Error(err)
		return nil, err
	default:
		db := t.DBService.DB
		stmt, err := db.PrepareContext(ctx, `insert into topics_users
	  (uuid4,
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
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5362}).Error(err)
			return nil, err
		}
		return stmt, nil
	}
}

// insertTopicsUser - Insert topic user details into database
func (t *TopicService) insertTopicsUser(ctx context.Context, stmt *sql.Stmt, tx *sql.Tx, poh TopicsUser, userEmail string, requestID string) (*TopicsUser, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5361}).Error(err)
		return nil, err
	default:

		res, err := tx.StmtContext(ctx, stmt).Exec(
			poh.UUID4,
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
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5363}).Error(err)
			return nil, err
		}
		uID, err := res.LastInsertId()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5364}).Error(err)
			return nil, err
		}
		poh.ID = uint(uID)
		uuid4Str, err := common.UUIDBytesToStr(poh.UUID4)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5365}).Error(err)
			return nil, err
		}
		poh.IDS = uuid4Str
		return &poh, nil
	}
}

// GetTopicByID - Get topic by ID
func (t *TopicService) GetTopicByID(ctx context.Context, ID uint, userEmail string, requestID string) (*Topic, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5348}).Error(err)
		return nil, err
	default:
		poh := Topic{}
		db := t.DBService.DB
		row := db.QueryRowContext(ctx, `select
    id,
		uuid4,
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
		updated_year from topics where id = ? and statusc = ?`, ID, common.Active)

		err := row.Scan(
			&poh.ID,
			&poh.UUID4,
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
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5349}).Error(err)
			return nil, err
		}
		uuid4Str, err := common.UUIDBytesToStr(poh.UUID4)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5350}).Error(err)
			return nil, err
		}
		poh.IDS = uuid4Str
		return &poh, nil
	}
}

// GetTopic - Get topic
func (t *TopicService) GetTopic(ctx context.Context, ID string, userEmail string, requestID string) (*Topic, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5351}).Error(err)
		return nil, err
	default:
		uuid4byte, err := common.UUIDStrToBytes(ID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5352}).Error(err)
			return nil, err
		}
		poh := Topic{}
		db := t.DBService.DB
		row := db.QueryRowContext(ctx, `select
    id,
		uuid4,
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
		updated_year from topics where uuid4 = ? and statusc = ?`, uuid4byte, common.Active)

		err = row.Scan(
			&poh.ID,
			&poh.UUID4,
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
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5353}).Error(err)
			return nil, err
		}
		uuid4Str, err := common.UUIDBytesToStr(poh.UUID4)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5354}).Error(err)
			return nil, err
		}
		poh.IDS = uuid4Str
		return &poh, nil
	}
}

// GetTopicByName - Get topic by name
func (t *TopicService) GetTopicByName(ctx context.Context, topicname string, userEmail string, requestID string) (*Topic, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5355}).Error(err)
		return nil, err
	default:
		poh := Topic{}
		db := t.DBService.DB
		row := db.QueryRowContext(ctx, `select
    id,
		uuid4,
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
		updated_year from topics where topic_name = ? and statusc = ?`, topicname, common.Active)

		err := row.Scan(
			&poh.ID,
			&poh.UUID4,
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
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5356}).Error(err)
			return nil, err
		}
		uuid4Str, err := common.UUIDBytesToStr(poh.UUID4)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5357}).Error(err)
			return nil, err
		}
		poh.IDS = uuid4Str
		return &poh, nil
	}
}

// GetTopicWithMessages - Get topic with messages
func (t *TopicService) GetTopicWithMessages(ctx context.Context, ID string, userEmail string, requestID string) (*Topic, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5310}).Error(err)
		return nil, err
	default:
		uuid4byte, err := common.UUIDStrToBytes(ID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5311}).Error(err)
			return nil, err
		}
		db := t.DBService.DB
		poh := &Topic{}
		tpc, err := t.GetTopic(ctx, ID, userEmail, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5312}).Error(err)
			return nil, err
		}
		var isPresent bool
		row := db.QueryRowContext(ctx, `select exists (select 1 from messages where topic_id = ?);`, tpc.ID)
		err = row.Scan(&isPresent)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5313}).Error(err)
			return nil, err
		}
		if isPresent {
			poh, err = t.GetTopicMessages(ctx, uuid4byte, userEmail, requestID)
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5380}).Error(err)
				return nil, err
			}
		} else {
			poh = tpc
		}

		if len(poh.Messages) > 0 {
			msgserv := &MessageService{DBService: t.DBService, RedisService: t.RedisService}
			Messages, err := msgserv.GetMessagesWithTextAttach(ctx, poh.Messages, userEmail, requestID)
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5320}).Error(err)
			}
			poh.Messages = Messages
		}
		return poh, nil
	}
}

// GetTopicMessages - get topic with messages
func (t *TopicService) GetTopicMessages(ctx context.Context, uuid4byte []byte, userEmail string, requestID string) (*Topic, error) {
	db := t.DBService.DB
	poh := Topic{}
	rows, err := db.QueryContext(ctx, `select 
      p.id,
			p.uuid4,
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
			m.uuid4,
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
			m.updated_year from topics p inner join messages m on (p.id = m.topic_id) where p.uuid4 = ?`, uuid4byte)

	if err != nil {
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5314}).Error(err)
		return nil, err
	}
	for rows.Next() {
		msg := Message{}
		err = rows.Scan(
			&poh.ID,
			&poh.UUID4,
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
			&msg.UUID4,
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
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5315}).Error(err)
			return nil, err
		}
		uuid4Str1, err := common.UUIDBytesToStr(poh.UUID4)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5316}).Error(err)
			return nil, err
		}
		poh.IDS = uuid4Str1

		uuid4Str, err := common.UUIDBytesToStr(msg.UUID4)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5317}).Error(err)
			return nil, err
		}
		msg.IDS = uuid4Str
		poh.Messages = append(poh.Messages, &msg)
	}

	err = rows.Close()
	if err != nil {
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5318}).Error(err)
		return nil, err
	}

	err = rows.Err()
	if err != nil {
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5319}).Error(err)
		return nil, err
	}
	return &poh, nil
}

// GetTopicsUser - Get user topics
func (t *TopicService) GetTopicsUser(ctx context.Context, ID uint, UserID uint, userEmail string, requestID string) (*TopicsUser, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5358}).Error(err)
		return nil, err
	default:
		poh := TopicsUser{}
		db := t.DBService.DB
		row := db.QueryRowContext(ctx, `select
    id,
		uuid4,
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
		updated_year from topics_users where topic_id = ? and user_id = ? and statusc = ?`, ID, UserID, common.Active)

		err := row.Scan(
			&poh.ID,
			&poh.UUID4,
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
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5359}).Error(err)
			return nil, err
		}
		uuid4Str, err := common.UUIDBytesToStr(poh.UUID4)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5360}).Error(err)
			return nil, err
		}
		poh.IDS = uuid4Str
		return &poh, nil
	}
}

//UpdateTopic - Update topic
func (t *TopicService) UpdateTopic(ctx context.Context, ID string, form *Topic, UserID string, userEmail string, requestID string) error {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5382}).Error(err)
		return err
	default:
		topic, err := t.GetTopic(ctx, ID, userEmail, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5383}).Error(err)
			return err
		}

		db := t.DBService.DB

		tn, tnday, tnweek, tnmonth, tnyear := common.GetTimeDetails()
		stmt, err := db.PrepareContext(ctx, `update topics set 
		  topic_name = ?,
      topic_desc = ?,
			updated_at = ?, 
			updated_day = ?, 
			updated_week = ?, 
			updated_month = ?, 
			updated_year = ? where id = ? and statusc = ?;`)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5385}).Error(err)
			return err
		}
		tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5384}).Error(err)
			return err
		}

		_, err = tx.StmtContext(ctx, stmt).Exec(
			form.TopicName,
			form.TopicDesc,
			tn,
			tnday,
			tnweek,
			tnmonth,
			tnyear,
			topic.ID,
			common.Active)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5387}).Error(err)
			err = tx.Rollback()
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5379}).Error(err)
				return err
			}
			err = stmt.Close()
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5388}).Error(err)
				return err
			}
			return err
		}

		err = tx.Commit()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5390}).Error(err)
			return err
		}
		err = stmt.Close()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5389}).Error(err)
			return err
		}

		return nil
	}
}

// UpdateNumMessagesPrepare - UpdateNumMessages prepare statement
func (t *TopicService) UpdateNumMessagesPrepare(ctx context.Context, userEmail string, requestID string) (*sql.Stmt, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5331}).Error(err)
		return nil, err
	default:
		db := t.DBService.DB
		stmt, err := db.PrepareContext(ctx, `update topics set 
		  num_messages = ?,
			updated_at = ?, 
			updated_day = ?, 
			updated_week = ?, 
			updated_month = ?, 
			updated_year = ? where id = ? and statusc = ?;`)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5334}).Error(err)
			return nil, err
		}
		return stmt, nil
	}
}

// UpdateNumMessages - update number of messages in topics
func (t *TopicService) UpdateNumMessages(ctx context.Context, stmt *sql.Stmt, tx *sql.Tx, numMessages uint, ID uint, userEmail string, requestID string) error {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5331}).Error(err)
		return err
	default:
		tn, tnday, tnweek, tnmonth, tnyear := common.GetTimeDetails()
		_, err := tx.StmtContext(ctx, stmt).Exec(
			numMessages,
			tn,
			tnday,
			tnweek,
			tnmonth,
			tnyear,
			ID,
			common.Active)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5335}).Error(err)
			return err
		}
		return nil
	}
}

// DeleteTopic - Delete topic
func (t *TopicService) DeleteTopic(ctx context.Context, ID string, userEmail string, requestID string) error {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5376}).Error(err)
		return err
	default:
		uuid4byte, err := common.UUIDStrToBytes(ID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5377}).Error(err)
			return err
		}
		db := t.DBService.DB
		tn, tnday, tnweek, tnmonth, tnyear := common.GetTimeDetails()
		stmt, err := db.PrepareContext(ctx, `update topics set 
		  statusc = ?,
			updated_at = ?, 
			updated_day = ?, 
			updated_week = ?, 
			updated_month = ?, 
			updated_year = ? where uuid4= ?;`)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5379}).Error(err)
			return err
		}

		tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5378}).Error(err)
			err = tx.Rollback()
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5379}).Error(err)
				return err
			}
			err = stmt.Close()
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5381}).Error(err)
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
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5380}).Error(err)
			err = tx.Rollback()
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5381}).Error(err)
				return err
			}
			err = stmt.Close()
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5381}).Error(err)
				return err
			}

			return err
		}

		err = tx.Commit()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5383}).Error(err)
			err = tx.Rollback()
			return err
		}

		err = stmt.Close()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5382}).Error(err)
			return err
		}
		return nil
	}
}
