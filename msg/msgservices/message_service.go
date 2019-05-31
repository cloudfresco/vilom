package msgservices

import (
	"context"
	"database/sql"
	"errors"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/go-redis/redis"

	"github.com/cloudfresco/vilom/common"
	"github.com/cloudfresco/vilom/user/userservices"
)

// Message - Message view representation
type Message struct {
	ID  uint
	IDS string

	NumLikes     uint
	NumUpvotes   uint
	NumDownvotes uint

	CategoryID uint
	TopicID    uint
	UserID     uint
	UgroupID   uint

	common.StatusDates

	MessageTexts       []*MessageText
	MessageAttachments []*MessageAttachment

	//only for logic purpose to create message
	Mtext   string
	Mattach string
}

// MessageText - MessageText view representation
type MessageText struct {
	ID    uint
	Mtext string

	CategoryID uint
	TopicID    uint
	MessageID  uint
	UserID     uint
	UgroupID   uint

	common.StatusDates
}

// MessageAttachment - MessageAttachment view representation
type MessageAttachment struct {
	ID         uint
	Mattach    string
	CategoryID uint
	TopicID    uint
	MessageID  uint
	UserID     uint
	UgroupID   uint

	common.StatusDates
}

// UserReply - UserReply view representation
type UserReply struct {
	ID        uint
	TopicID   uint
	MessageID uint `sql:"default:'0'"`
	UserID    uint
	UgroupID  uint `sql:"default:'0'"`

	common.StatusDates
}

// UserLike - UserLike view representation
type UserLike struct {
	ID uint

	TopicID   uint
	MessageID uint `sql:"default:'0'"`

	UgroupID uint `sql:"default:'0'"`
	UserID   uint

	common.StatusDates
}

// UserVote - UserVote view representation
type UserVote struct {
	ID uint

	TopicID   uint
	MessageID uint `sql:"default:'0'"`
	Vote      uint `sql:"default:'0'"`

	UgroupID uint `sql:"default:'0'"`
	UserID   uint

	common.StatusDates
}

// MessageService - For accessing message services
type MessageService struct {
	Config       *common.RedisOptions
	Db           *sql.DB
	RedisClient  *redis.Client
	LimitDefault string
}

// NewMessageService - Create message service
func NewMessageService(config *common.RedisOptions,
	db *sql.DB,
	redisClient *redis.Client,
	limitDefault string) *MessageService {
	return &MessageService{config, db, redisClient, limitDefault}
}

//Create - Create message
func (t *MessageService) Create(ctx context.Context, form *Message, UserID string, userEmail string, requestID string) (*Message, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6300}).Error(err)
		return nil, err
	default:
		userserv := &userservices.UserService{Config: t.Config, Db: t.Db, RedisClient: t.RedisClient}
		user, err := userserv.GetUser(ctx, UserID, userEmail, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6301}).Error(err)
			return nil, err
		}
		db := t.Db
		tx, err := db.Begin()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6302}).Error(err)
			return nil, err
		}
		tn := time.Now().UTC()
		_, week := tn.ISOWeek()
		day := tn.YearDay()

		msg := Message{}

		msg.IDS = common.GetUID()
		msg.NumLikes = uint(0)
		msg.NumUpvotes = uint(0)
		msg.NumDownvotes = uint(0)
		msg.CategoryID = form.CategoryID
		msg.TopicID = form.TopicID
		msg.UserID = user.ID
		msg.UgroupID = form.UgroupID
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

		Message, err := t.InsertMessage(ctx, tx, msg, userEmail, requestID)

		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6303}).Error(err)
			err = tx.Rollback()
			return nil, err
		}

		msgtxt := MessageText{}
		msgtxt.Mtext = form.Mtext
		msgtxt.CategoryID = form.CategoryID
		msgtxt.TopicID = form.TopicID
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

		_, err = t.InsertMessageText(ctx, tx, msgtxt, userEmail, requestID)

		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6304}).Error(err)
			err = tx.Rollback()
			return nil, err
		}

		msgath := MessageAttachment{}
		msgath.Mattach = form.Mattach
		msgath.CategoryID = form.CategoryID
		msgath.TopicID = form.TopicID
		msgath.MessageID = Message.ID
		msgath.UserID = user.ID
		msgath.UgroupID = form.UgroupID
		/*  StatusDates  */
		msgath.Statusc = common.Active
		msgath.CreatedAt = tn.UTC()
		msgath.UpdatedAt = tn.UTC()
		msgath.CreatedDay = uint(day)
		msgath.CreatedWeek = uint(week)
		msgath.CreatedMonth = uint(tn.Month())
		msgath.CreatedYear = uint(tn.Year())
		msgath.UpdatedDay = uint(day)
		msgath.UpdatedWeek = uint(week)
		msgath.UpdatedMonth = uint(tn.Month())
		msgath.UpdatedYear = uint(tn.Year())

		_, err = t.InsertMessageAttachment(ctx, tx, msgath, userEmail, requestID)

		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6305}).Error(err)
			err = tx.Rollback()
			return nil, err
		}

		UpdatedDay := uint(day)
		UpdatedWeek := uint(week)
		UpdatedMonth := uint(tn.Month())
		UpdatedYear := uint(tn.Year())

		topicserv := &TopicService{t.Config, t.Db, t.RedisClient, t.LimitDefault}
		topic, err := topicserv.GetTopicByID(ctx, form.TopicID, userEmail, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6306}).Error(err)
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
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6307}).Error(err)
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

		err = stmt.Close()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6308}).Error(err)
			err = stmt.Close()
			err = tx.Rollback()
			return nil, err
		}

		ur := UserReply{}

		ur.TopicID = form.TopicID
		ur.MessageID = Message.ID
		ur.UserID = user.ID
		ur.UgroupID = form.UgroupID
		/*  StatusDates  */
		ur.Statusc = common.Active
		ur.CreatedAt = tn.UTC()
		ur.UpdatedAt = tn.UTC()
		ur.CreatedDay = uint(day)
		ur.CreatedWeek = uint(week)
		ur.CreatedMonth = uint(tn.Month())
		ur.CreatedYear = uint(tn.Year())
		ur.UpdatedDay = uint(day)
		ur.UpdatedWeek = uint(week)
		ur.UpdatedMonth = uint(tn.Month())
		ur.UpdatedYear = uint(tn.Year())

		_, err = t.InsertUserReply(ctx, tx, ur, userEmail, requestID)

		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6309}).Error(err)
			err = tx.Rollback()
			return nil, err
		}

		err = tx.Commit()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6310}).Error(err)
			return nil, err
		}
		return Message, nil
	}
}

// UserLikeCreate - Create user likes messages
func (t *MessageService) UserLikeCreate(ctx context.Context, form *UserLike, UserID string, userEmail string, requestID string) (*UserLike, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6311}).Error(err)
		return nil, err
	default:
		userserv := &userservices.UserService{Config: t.Config, Db: t.Db, RedisClient: t.RedisClient}
		user, err := userserv.GetUser(ctx, UserID, userEmail, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6312}).Error(err)
			return nil, err
		}
		db := t.Db
		tx, err := db.Begin()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6313}).Error(err)
			return nil, err
		}

		tn := time.Now().UTC()
		_, week := tn.ISOWeek()
		day := tn.YearDay()

		ul := UserLike{}
		ul.TopicID = form.TopicID
		ul.MessageID = form.MessageID
		ul.UgroupID = form.UgroupID
		ul.UserID = user.ID
		/*  StatusDates  */
		ul.Statusc = common.Active
		ul.CreatedAt = tn.UTC()
		ul.UpdatedAt = tn.UTC()
		ul.CreatedDay = uint(day)
		ul.CreatedWeek = uint(week)
		ul.CreatedMonth = uint(tn.Month())
		ul.CreatedYear = uint(tn.Year())
		ul.UpdatedDay = uint(day)
		ul.UpdatedWeek = uint(week)
		ul.UpdatedMonth = uint(tn.Month())
		ul.UpdatedYear = uint(tn.Year())

		UserLk, err := t.InsertUserLike(ctx, tx, ul, userEmail, requestID)

		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6314}).Error(err)
			err = tx.Rollback()
			return nil, err
		}

		err = tx.Commit()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6315}).Error(err)
			return nil, err
		}
		return UserLk, nil
	}
}

// UserVoteCreate - Create User Vote
func (t *MessageService) UserVoteCreate(ctx context.Context, form *UserVote, UserID string, userEmail string, requestID string) (*UserVote, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6316}).Error(err)
		return nil, err
	default:
		userserv := &userservices.UserService{Config: t.Config, Db: t.Db, RedisClient: t.RedisClient}
		user, err := userserv.GetUser(ctx, UserID, userEmail, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6317}).Error(err)
			return nil, err
		}
		db := t.Db
		tx, err := db.Begin()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6318}).Error(err)
			return nil, err
		}

		tn := time.Now().UTC()
		_, week := tn.ISOWeek()
		day := tn.YearDay()

		ul := UserVote{}
		ul.TopicID = form.TopicID
		ul.MessageID = form.MessageID
		ul.UgroupID = form.UgroupID
		ul.UserID = user.ID
		ul.Vote = form.Vote
		/*  StatusDates  */
		ul.Statusc = common.Active
		ul.CreatedAt = tn.UTC()
		ul.UpdatedAt = tn.UTC()
		ul.CreatedDay = uint(day)
		ul.CreatedWeek = uint(week)
		ul.CreatedMonth = uint(tn.Month())
		ul.CreatedYear = uint(tn.Year())
		ul.UpdatedDay = uint(day)
		ul.UpdatedWeek = uint(week)
		ul.UpdatedMonth = uint(tn.Month())
		ul.UpdatedYear = uint(tn.Year())

		UserVt, err := t.InsertUserVote(ctx, tx, ul, userEmail, requestID)

		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6319}).Error(err)
			err = tx.Rollback()
			return nil, err
		}

		err = tx.Commit()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6320}).Error(err)
			err = tx.Rollback()
			return nil, err
		}
		return UserVt, nil
	}
}

// InsertMessage - Insert message details into database
func (t *MessageService) InsertMessage(ctx context.Context, tx *sql.Tx, msg Message, userEmail string, requestID string) (*Message, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6321}).Error(err)
		return nil, err
	default:
		stmt, err := tx.PrepareContext(ctx, `insert into messages
	  ( 
			id_s,
			num_likes,
			num_upvotes,
			num_downvotes,
			category_id,
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
					?,?,?,?,?,?,?,?,?);`)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6322}).Error(err)
			err = stmt.Close()
			return nil, err
		}

		res, err := stmt.ExecContext(ctx,
			msg.IDS,
			msg.NumLikes,
			msg.NumUpvotes,
			msg.NumDownvotes,
			msg.CategoryID,
			msg.TopicID,
			msg.UserID,
			msg.UgroupID,
			msg.Statusc,
			msg.CreatedAt,
			msg.UpdatedAt,
			msg.CreatedDay,
			msg.CreatedWeek,
			msg.CreatedMonth,
			msg.CreatedYear,
			msg.UpdatedDay,
			msg.UpdatedWeek,
			msg.UpdatedMonth,
			msg.UpdatedYear)

		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6323}).Error(err)
			err = stmt.Close()
			return nil, err
		}
		uID, err := res.LastInsertId()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6324}).Error(err)
			err = stmt.Close()
			return nil, err
		}
		msg.ID = uint(uID)
		err = stmt.Close()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6325}).Error(err)
			return nil, err
		}
		return &msg, nil
	}
}

// InsertMessageText - Insert message text details in database
func (t *MessageService) InsertMessageText(ctx context.Context, tx *sql.Tx, msgtxt MessageText, userEmail string, requestID string) (*MessageText, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6326}).Error(err)
		return nil, err
	default:
		stmt, err := tx.PrepareContext(ctx, `insert into message_texts
	  ( 
			mtext,
			category_id,
			topic_id,
			message_id,
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
					?,?,?,?,?,?,?);`)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6327}).Error(err)
			err = stmt.Close()
			return nil, err
		}
		res, err := stmt.ExecContext(ctx,
			msgtxt.Mtext,
			msgtxt.CategoryID,
			msgtxt.TopicID,
			msgtxt.MessageID,
			msgtxt.UgroupID,
			msgtxt.UserID,
			/*  StatusDates  */
			msgtxt.Statusc,
			msgtxt.CreatedAt,
			msgtxt.UpdatedAt,
			msgtxt.CreatedDay,
			msgtxt.CreatedWeek,
			msgtxt.CreatedMonth,
			msgtxt.CreatedYear,
			msgtxt.UpdatedDay,
			msgtxt.UpdatedWeek,
			msgtxt.UpdatedMonth,
			msgtxt.UpdatedYear)

		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6328}).Error(err)
			err = stmt.Close()
			return nil, err
		}
		uID, err := res.LastInsertId()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6329}).Error(err)
			err = stmt.Close()
			return nil, err
		}
		msgtxt.ID = uint(uID)
		err = stmt.Close()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6330}).Error(err)
			return nil, err
		}
		return &msgtxt, nil
	}
}

// InsertMessageAttachment - Insert message attachment details in database
func (t *MessageService) InsertMessageAttachment(ctx context.Context, tx *sql.Tx, msgath MessageAttachment, userEmail string, requestID string) (*MessageAttachment, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6331}).Error(err)
		return nil, err
	default:
		stmt, err := tx.PrepareContext(ctx, `insert into message_attachments
	  ( 
			mattach,
			category_id,
			topic_id,
			message_id,
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
					?,?,?,?,?,?,?);`)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6332}).Error(err)
			err = stmt.Close()
			return nil, err
		}
		res, err := stmt.ExecContext(ctx,
			msgath.Mattach,
			msgath.CategoryID,
			msgath.TopicID,
			msgath.MessageID,
			msgath.UgroupID,
			msgath.UserID,
			msgath.Statusc,
			msgath.CreatedAt,
			msgath.UpdatedAt,
			msgath.CreatedDay,
			msgath.CreatedWeek,
			msgath.CreatedMonth,
			msgath.CreatedYear,
			msgath.UpdatedDay,
			msgath.UpdatedWeek,
			msgath.UpdatedMonth,
			msgath.UpdatedYear)

		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6333}).Error(err)
			err = stmt.Close()
			return nil, err
		}
		uID, err := res.LastInsertId()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6334}).Error(err)
			err = stmt.Close()
			return nil, err
		}
		msgath.ID = uint(uID)
		err = stmt.Close()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6335}).Error(err)
			return nil, err
		}
		return &msgath, nil
	}
}

// InsertUserReply - Insert user reply details into database
func (t *MessageService) InsertUserReply(ctx context.Context, tx *sql.Tx, ur UserReply, userEmail string, requestID string) (*UserReply, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6336}).Error(err)
		return nil, err
	default:
		stmt, err := tx.PrepareContext(ctx, `insert into user_replies
	  ( 
			topic_id,
			message_id,
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
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6337}).Error(err)
			err = stmt.Close()
			return nil, err
		}
		res, err := stmt.ExecContext(ctx,
			ur.TopicID,
			ur.MessageID,
			ur.UserID,
			ur.UgroupID,
			/*  StatusDates  */
			ur.Statusc,
			ur.CreatedAt,
			ur.UpdatedAt,
			ur.CreatedDay,
			ur.CreatedWeek,
			ur.CreatedMonth,
			ur.CreatedYear,
			ur.UpdatedDay,
			ur.UpdatedWeek,
			ur.UpdatedMonth,
			ur.UpdatedYear)

		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6338}).Error(err)
			err = stmt.Close()
			return nil, err
		}
		uID, err := res.LastInsertId()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6339}).Error(err)
			err = stmt.Close()
			return nil, err
		}
		ur.ID = uint(uID)
		err = stmt.Close()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6340}).Error(err)
			return nil, err
		}
		return &ur, nil
	}
}

// InsertUserLike - Insert User like details in database
func (t *MessageService) InsertUserLike(ctx context.Context, tx *sql.Tx, ur UserLike, userEmail string, requestID string) (*UserLike, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6341}).Error(err)
		return nil, err
	default:
		stmt, err := tx.PrepareContext(ctx, `insert into user_likes
	  ( 
			topic_id,
			message_id,
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
					?,?,?);`)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6342}).Error(err)
			err = stmt.Close()
			return nil, err
		}
		res, err := stmt.ExecContext(ctx,
			ur.TopicID,
			ur.MessageID,
			ur.UgroupID,
			ur.UserID,
			/*  StatusDates  */
			ur.Statusc,
			ur.CreatedAt,
			ur.UpdatedAt,
			ur.CreatedDay,
			ur.CreatedWeek,
			ur.CreatedMonth,
			ur.CreatedYear,
			ur.UpdatedDay,
			ur.UpdatedWeek,
			ur.UpdatedMonth,
			ur.UpdatedYear)

		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6343}).Error(err)
			err = stmt.Close()
			return nil, err
		}
		uID, err := res.LastInsertId()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6344}).Error(err)
			err = stmt.Close()
			return nil, err
		}
		ur.ID = uint(uID)
		err = stmt.Close()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6345}).Error(err)
			return nil, err
		}
		return &ur, nil
	}
}

// InsertUserVote - Insert User vote details into database
func (t *MessageService) InsertUserVote(ctx context.Context, tx *sql.Tx, ur UserVote, userEmail string, requestID string) (*UserVote, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6346}).Error(err)
		return nil, err
	default:
		stmt, err := tx.PrepareContext(ctx, `insert into user_votes
	  ( 
			topic_id,
			message_id,
			vote,
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
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6347}).Error(err)
			err = stmt.Close()
			return nil, err
		}
		res, err := stmt.ExecContext(ctx,
			ur.TopicID,
			ur.MessageID,
			ur.Vote,
			ur.UgroupID,
			ur.UserID,
			/*  StatusDates  */
			ur.Statusc,
			ur.CreatedAt,
			ur.UpdatedAt,
			ur.CreatedDay,
			ur.CreatedWeek,
			ur.CreatedMonth,
			ur.CreatedYear,
			ur.UpdatedDay,
			ur.UpdatedWeek,
			ur.UpdatedMonth,
			ur.UpdatedYear)

		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6348}).Error(err)
			err = stmt.Close()
			return nil, err
		}
		uID, err := res.LastInsertId()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6349}).Error(err)
			err = stmt.Close()
			return nil, err
		}
		ur.ID = uint(uID)
		err = stmt.Close()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6350}).Error(err)
			return nil, err
		}
		return &ur, nil
	}
}

// GetMessage - Get message
func (t *MessageService) GetMessage(ctx context.Context, ID string, userEmail string, requestID string) (*Message, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6351}).Error(err)
		return nil, err
	default:
		msg := Message{}
		row := t.Db.QueryRowContext(ctx, `select
      id,
 			id_s,
			num_likes,
			num_upvotes,
			num_downvotes,
			category_id,
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
			updated_year from messages where id_s = ?;`, ID)

		err := row.Scan(
			&msg.ID,
			&msg.IDS,
			&msg.NumLikes,
			&msg.NumUpvotes,
			&msg.NumDownvotes,
			&msg.CategoryID,
			&msg.TopicID,
			&msg.UserID,
			&msg.UgroupID,
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
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6352}).Error(err)
			return nil, err
		}

		return &msg, nil
	}
}

// GetMessagesWithTextAttach - Get messages with attachemnts
func (t *MessageService) GetMessagesWithTextAttach(ctx context.Context, messages []*Message, userEmail string, requestID string) ([]*Message, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6353}).Error(err)
		return nil, err
	default:
		db := t.Db
		pohs := []*Message{}

		for _, message := range messages {
			var isPresent bool
			row := db.QueryRowContext(ctx, `select exists (select 1 from message_texts where message_id = ?);`, message.ID)
			err := row.Scan(&isPresent)
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6354}).Error(err)
			}
			if isPresent {

				rows, err := db.QueryContext(ctx, `select 
				mtext,
				category_id,
				topic_id,
				message_id,
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
				updated_year from message_texts where message_id = ?`, message.ID)

				if err != nil {
					log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6355}).Error(err)
					return nil, err
				}
				for rows.Next() {
					msgtxt := MessageText{}
					err = rows.Scan(
						&msgtxt.Mtext,
						&msgtxt.CategoryID,
						&msgtxt.TopicID,
						&msgtxt.MessageID,
						&msgtxt.UgroupID,
						&msgtxt.UserID,
						/*  StatusDates  */
						&msgtxt.Statusc,
						&msgtxt.CreatedAt,
						&msgtxt.UpdatedAt,
						&msgtxt.CreatedDay,
						&msgtxt.CreatedWeek,
						&msgtxt.CreatedMonth,
						&msgtxt.CreatedYear,
						&msgtxt.UpdatedDay,
						&msgtxt.UpdatedWeek,
						&msgtxt.UpdatedMonth,
						&msgtxt.UpdatedYear)

					if err != nil {
						log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6356}).Error(err)
					}

					message.MessageTexts = append(message.MessageTexts, &msgtxt)
				}

				err = rows.Close()
				if err != nil {
					log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6357}).Error(err)
					return nil, err
				}
			}

			var isPresent1 bool
			row1 := db.QueryRowContext(ctx, `select exists (select 1 from message_attachments where message_id = ?);`, message.ID)
			err = row1.Scan(&isPresent1)
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6358}).Error(err)
			}
			if isPresent1 {
				msgath := MessageAttachment{}
				rows, err := db.QueryContext(ctx, `select 
				mattach,
				category_id,
				topic_id,
				message_id,
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
				updated_year from message_attachments where message_id = ?`, message.ID)

				if err != nil {
					log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6359}).Error(err)
				}
				for rows.Next() {
					err = rows.Scan(
						&msgath.Mattach,
						&msgath.CategoryID,
						&msgath.TopicID,
						&msgath.MessageID,
						&msgath.UgroupID,
						&msgath.UserID,
						/*  StatusDates  */
						&msgath.Statusc,
						&msgath.CreatedAt,
						&msgath.UpdatedAt,
						&msgath.CreatedDay,
						&msgath.CreatedWeek,
						&msgath.CreatedMonth,
						&msgath.CreatedYear,
						&msgath.UpdatedDay,
						&msgath.UpdatedWeek,
						&msgath.UpdatedMonth,
						&msgath.UpdatedYear)

					if err != nil {
						log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6360}).Error(err)
					}

					message.MessageAttachments = append(message.MessageAttachments, &msgath)
				}

				err = rows.Close()
				if err != nil {
					log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6361}).Error(err)
					return nil, err
				}
			}

			pohs = append(pohs, message)
		}

		return pohs, nil
	}
}
