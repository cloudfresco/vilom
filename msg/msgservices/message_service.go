package msgservices

import (
	"context"
	"database/sql"
	"errors"

	log "github.com/sirupsen/logrus"

	"github.com/cloudfresco/vilom/common"
	"github.com/cloudfresco/vilom/user/userservices"
)

/* error message range: 6300-6999 */

// For validation of Message fields
const (
	MtextLenMin = 1
	MtextLenMax = 50
)

// Message - Message view representation
type Message struct {
	ID    uint   `json:"id,omitempty"`
	UUID4 []byte `json:"-"`
	IDS   string `json:"id_s,omitempty"`

	NumLikes     uint `json:"num_likes,omitempty"`
	NumUpvotes   uint `json:"num_upvotes,omitempty"`
	NumDownvotes uint `json:"num_downvotes,omitempty"`

	WorkspaceID uint `json:"workspace_id,omitempty"`
	ChannelID   uint `json:"channel_id,omitempty"`
	UserID      uint `json:"user_id,omitempty"`
	UgroupID    uint `json:"ugroup_id,omitempty"`

	common.StatusDates

	MessageTexts       []*MessageText
	MessageAttachments []*MessageAttachment

	//only for logic purpose to create message
	Mtext   string
	Mattach string
}

// MessageText - MessageText view representation
type MessageText struct {
	ID          uint   `json:"id,omitempty"`
	UUID4       []byte `json:"-"`
	Mtext       string `json:"mtext,omitempty"`
	WorkspaceID uint   `json:"workspace_id,omitempty"`
	ChannelID   uint   `json:"channel_id,omitempty"`
	MessageID   uint   `json:"message_id,omitempty"`
	UserID      uint   `json:"user_id,omitempty"`
	UgroupID    uint   `json:"ugroup_id,omitempty"`

	common.StatusDates
}

// MessageAttachment - MessageAttachment view representation
type MessageAttachment struct {
	ID          uint   `json:"id,omitempty"`
	UUID4       []byte `json:"-"`
	Mattach     string `json:"mattach,omitempty"`
	WorkspaceID uint   `json:"workspace_id,omitempty"`
	ChannelID   uint   `json:"channel_id,omitempty"`
	MessageID   uint   `json:"message_id,omitempty"`
	UserID      uint   `json:"user_id,omitempty"`
	UgroupID    uint   `json:"ugroup_id,omitempty"`

	common.StatusDates
}

// UserReply - UserReply view representation
type UserReply struct {
	ID        uint   `json:"id,omitempty"`
	UUID4     []byte `json:"-"`
	ChannelID uint   `json:"channel_id,omitempty"`
	MessageID uint   `json:"message_id,omitempty"`
	UserID    uint   `json:"user_id,omitempty"`
	UgroupID  uint   `json:"ugroup_id,omitempty"`

	common.StatusDates
}

// UserLike - UserLike view representation
type UserLike struct {
	ID        uint   `json:"id,omitempty"`
	UUID4     []byte `json:"-"`
	ChannelID uint   `json:"channel_id,omitempty"`
	MessageID uint   `json:"message_id,omitempty"`
	UserID    uint   `json:"user_id,omitempty"`
	UgroupID  uint   `json:"ugroup_id,omitempty"`

	common.StatusDates
}

// UserVote - UserVote view representation
type UserVote struct {
	ID        uint   `json:"id,omitempty"`
	UUID4     []byte `json:"-"`
	Vote      uint   `json:"vote,omitempty"`
	ChannelID uint   `json:"channel_id,omitempty"`
	MessageID uint   `json:"message_id,omitempty"`
	UserID    uint   `json:"user_id,omitempty"`
	UgroupID  uint   `json:"ugroup_id,omitempty"`

	common.StatusDates
}

// MessageServiceIntf - interface for Message Service
type MessageServiceIntf interface {
	CreateMessage(ctx context.Context, form *Message, UserID string, rplymsg bool, userEmail string, requestID string) (*Message, error)
	CreateUserLike(ctx context.Context, form *UserLike, UserID string, userEmail string, requestID string) (*UserLike, error)
	CreateUserVote(ctx context.Context, form *UserVote, UserID string, userEmail string, requestID string) (*UserVote, error)
	GetMessage(ctx context.Context, ID string, userEmail string, requestID string) (*Message, error)
	GetMessagesWithTextAttach(ctx context.Context, messages []*Message, userEmail string, requestID string) ([]*Message, error)
	GetMessagesTexts(ctx context.Context, messageID uint, userEmail string, requestID string) ([]*MessageText, error)
	GetMessageAttachments(ctx context.Context, messageID uint, userEmail string, requestID string) ([]*MessageAttachment, error)
	UpdateMessage(ctx context.Context, ID string, form *Message, UserID string, userEmail string, requestID string) error
	DeleteMessage(ctx context.Context, ID string, userEmail string, requestID string) error
}

// MessageService - For accessing message services
type MessageService struct {
	DBService    *common.DBService
	RedisService *common.RedisService
}

// NewMessageService - Create message service
func NewMessageService(dbOpt *common.DBService, redisOpt *common.RedisService) *MessageService {
	return &MessageService{
		DBService:    dbOpt,
		RedisService: redisOpt,
	}
}

//CreateMessage - Create message
func (m *MessageService) CreateMessage(ctx context.Context, form *Message, UserID string, rplymsg bool, userEmail string, requestID string) (*Message, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6300}).Error(err)
		return nil, err
	default:
		db := m.DBService.DB

		insertMessageStmt, insertMessageTextStmt, insertMessageAttachmentStmt, updateNumMessagesStmt, insertUserReplyStmt, err := m.createMessagePrepareStmts(ctx, userEmail, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6301}).Error(err)
			return nil, err
		}

		tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6302}).Error(err)
			err = m.createMessagePrepareStmtsClose(ctx, insertMessageStmt, insertMessageTextStmt, insertMessageAttachmentStmt, updateNumMessagesStmt, insertUserReplyStmt, userEmail, requestID)
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6304}).Error(err)
				return nil, err
			}
			return nil, err
		}

		msg, err := m.createMessage(ctx, insertMessageStmt, insertMessageTextStmt, insertMessageAttachmentStmt, updateNumMessagesStmt, insertUserReplyStmt, tx, form, UserID, rplymsg, userEmail, requestID)

		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6305}).Error(err)
			err = tx.Rollback()
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6306}).Error(err)
				return nil, err
			}
			err = m.createMessagePrepareStmtsClose(ctx, insertMessageStmt, insertMessageTextStmt, insertMessageAttachmentStmt, updateNumMessagesStmt, insertUserReplyStmt, userEmail, requestID)
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6307}).Error(err)
				return nil, err
			}
			return nil, err
		}

		err = tx.Commit()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6308}).Error(err)
			return nil, err
		}

		err = m.createMessagePrepareStmtsClose(ctx, insertMessageStmt, insertMessageTextStmt, insertMessageAttachmentStmt, updateNumMessagesStmt, insertUserReplyStmt, userEmail, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6309}).Error(err)
			return nil, err
		}

		return msg, nil
	}

}

//createMessagePrepareStmts - Create message Prepare Statements
func (m *MessageService) createMessagePrepareStmts(ctx context.Context, userEmail string, requestID string) (*sql.Stmt, *sql.Stmt, *sql.Stmt, *sql.Stmt, *sql.Stmt, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6310}).Error(err)
		return nil, nil, nil, nil, nil, err
	default:
		insertMessageStmt, err := m.insertMessagePrepare(ctx, userEmail, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6311}).Error(err)
			return nil, nil, nil, nil, nil, err
		}
		insertMessageTextStmt, err := m.insertMessageTextPrepare(ctx, userEmail, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6312}).Error(err)
			return nil, nil, nil, nil, nil, err
		}
		insertMessageAttachmentStmt, err := m.insertMessageAttachmentPrepare(ctx, userEmail, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6313}).Error(err)
			return nil, nil, nil, nil, nil, err
		}
		updateNumMessagesStmt, err := m.updateNumMessagesPrepare(ctx, userEmail, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6314}).Error(err)
			return nil, nil, nil, nil, nil, err
		}
		insertUserReplyStmt, err := m.insertUserReplyPrepare(ctx, userEmail, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6315}).Error(err)
			return nil, nil, nil, nil, nil, err
		}
		return insertMessageStmt, insertMessageTextStmt, insertMessageAttachmentStmt, updateNumMessagesStmt, insertUserReplyStmt, nil

	}
}

//createMessagePrepareStmtsClose - Close Prepare Statements
func (m *MessageService) createMessagePrepareStmtsClose(ctx context.Context, insertMessageStmt *sql.Stmt, insertMessageTextStmt *sql.Stmt, insertMessageAttachmentStmt *sql.Stmt, updateNumMessagesStmt *sql.Stmt, insertUserReplyStmt *sql.Stmt, userEmail string, requestID string) error {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6316}).Error(err)
		return err
	default:
		err := insertMessageStmt.Close()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6317}).Error(err)
			return err
		}
		err = insertMessageTextStmt.Close()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6318}).Error(err)
			return err
		}
		err = insertMessageAttachmentStmt.Close()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6319}).Error(err)
			return err
		}
		err = updateNumMessagesStmt.Close()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6320}).Error(err)
			return err
		}
		err = insertUserReplyStmt.Close()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6321}).Error(err)
			return err
		}

		return nil
	}
}

//createMessage - Create message
func (m *MessageService) createMessage(ctx context.Context, stmt *sql.Stmt, insertMessageTextStmt *sql.Stmt, insertMessageAttachmentStmt *sql.Stmt, updateNumMessagesStmt *sql.Stmt, insertUserReplyStmt *sql.Stmt, tx *sql.Tx, form *Message, UserID string, rplymsg bool, userEmail string, requestID string) (*Message, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6322}).Error(err)
		return nil, err
	default:
		userserv := &userservices.UserService{DBService: m.DBService, RedisService: m.RedisService}
		user, err := userserv.GetUser(ctx, UserID, userEmail, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6323}).Error(err)
			return nil, err
		}
		tn, tnday, tnweek, tnmonth, tnyear := common.GetTimeDetails()
		msg := Message{}
		msg.UUID4, err = common.GetUUIDBytes()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6324}).Error(err)
			return nil, err
		}
		msg.NumLikes = uint(0)
		msg.NumUpvotes = uint(0)
		msg.NumDownvotes = uint(0)
		msg.WorkspaceID = form.WorkspaceID
		msg.ChannelID = form.ChannelID
		msg.UserID = user.ID
		msg.UgroupID = form.UgroupID
		msg.Statusc = common.Active
		msg.CreatedAt = tn
		msg.UpdatedAt = tn
		msg.CreatedDay = tnday
		msg.CreatedWeek = tnweek
		msg.CreatedMonth = tnmonth
		msg.CreatedYear = tnyear
		msg.UpdatedDay = tnday
		msg.UpdatedWeek = tnweek
		msg.UpdatedMonth = tnmonth
		msg.UpdatedYear = tnyear

		err = m.insertMessage(ctx, stmt, tx, &msg, userEmail, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6325}).Error(err)
			return nil, err
		}

		msgtext, err := m.createMessageText(ctx, insertMessageTextStmt, tx, form, msg.ID, user.ID, userEmail, requestID)

		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6326}).Error(err)
			return nil, err
		}
		msg.MessageTexts = append(msg.MessageTexts, msgtext)
		if form.Mattach != "" {
			msgattach, err := m.createMessageAttachment(ctx, insertMessageAttachmentStmt, tx, form, msg.ID, user.ID, userEmail, requestID)

			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6327}).Error(err)
				return nil, err
			}
			msg.MessageAttachments = append(msg.MessageAttachments, msgattach)
		}
		channelserv := &ChannelService{DBService: m.DBService, RedisService: m.RedisService}
		channel, err := channelserv.GetChannelByID(ctx, form.ChannelID, userEmail, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6328}).Error(err)
			return nil, err
		}

		numMessages := channel.NumMessages + 1
		err = m.updateNumMessages(ctx, updateNumMessagesStmt, tx, numMessages, channel.ID, userEmail, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6329}).Error(err)
			return nil, err
		}

		if rplymsg {
			err = m.createUserReply(ctx, insertUserReplyStmt, tx, form.ChannelID, msg.ID, user.ID, form.UgroupID, userEmail, requestID)
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6330}).Error(err)
				err = tx.Rollback()
				return nil, err
			}
		}

		return &msg, nil
	}
}

// insertMessagePrepare - Insert message details Prepare Statement
func (m *MessageService) insertMessagePrepare(ctx context.Context, userEmail string, requestID string) (*sql.Stmt, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6331}).Error(err)
		return nil, err
	default:
		db := m.DBService.DB
		stmt, err := db.PrepareContext(ctx, `insert into messages
	  ( 
			uuid4,
			num_likes,
			num_upvotes,
			num_downvotes,
			workspace_id,
			channel_id,
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
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6332}).Error(err)
			return nil, err
		}
		return stmt, nil
	}

}

// insertMessage - Insert message details into database
func (m *MessageService) insertMessage(ctx context.Context, stmt *sql.Stmt, tx *sql.Tx, msg *Message, userEmail string, requestID string) error {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6333}).Error(err)
		return err
	default:
		res, err := tx.StmtContext(ctx, stmt).Exec(
			msg.UUID4,
			msg.NumLikes,
			msg.NumUpvotes,
			msg.NumDownvotes,
			msg.WorkspaceID,
			msg.ChannelID,
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
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6334}).Error(err)
			return err
		}
		uID, err := res.LastInsertId()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6335}).Error(err)
			return err
		}
		msg.ID = uint(uID)
		uuid4Str, err := common.UUIDBytesToStr(msg.UUID4)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6336}).Error(err)
			return err
		}
		msg.IDS = uuid4Str
		return nil
	}
}

//CreateMessageText - Create message Text
func (m *MessageService) createMessageText(ctx context.Context, stmt *sql.Stmt, tx *sql.Tx, form *Message, messageID uint, userID uint, userEmail string, requestID string) (*MessageText, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6337}).Error(err)
		return nil, err
	default:
		tn, tnday, tnweek, tnmonth, tnyear := common.GetTimeDetails()
		uuid4, err := common.GetUUIDBytes()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6338}).Error(err)
			return nil, err
		}
		msgtxt := MessageText{}
		msgtxt.UUID4 = uuid4
		msgtxt.Mtext = form.Mtext
		msgtxt.WorkspaceID = form.WorkspaceID
		msgtxt.ChannelID = form.ChannelID
		msgtxt.MessageID = messageID
		msgtxt.UserID = userID
		msgtxt.UgroupID = form.UgroupID
		msgtxt.Statusc = common.Active
		msgtxt.CreatedAt = tn
		msgtxt.UpdatedAt = tn
		msgtxt.CreatedDay = tnday
		msgtxt.CreatedWeek = tnweek
		msgtxt.CreatedMonth = tnmonth
		msgtxt.CreatedYear = tnyear
		msgtxt.UpdatedDay = tnday
		msgtxt.UpdatedWeek = tnweek
		msgtxt.UpdatedMonth = tnmonth
		msgtxt.UpdatedYear = tnyear
		err = m.insertMessageText(ctx, stmt, tx, &msgtxt, userEmail, requestID)

		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6339}).Error(err)
			return nil, err
		}

		return &msgtxt, nil
	}
}

// insertMessageTextPrepare - Insert message text Prepare Statement
func (m *MessageService) insertMessageTextPrepare(ctx context.Context, userEmail string, requestID string) (*sql.Stmt, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6340}).Error(err)
		return nil, err
	default:
		db := m.DBService.DB
		stmt, err := db.PrepareContext(ctx, `insert into message_texts
	  ( 
      uuid4,
			mtext,
			workspace_id,
			channel_id,
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
					?,?,?,?,?,?,?,?);`)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6341}).Error(err)
			return nil, err
		}
		return stmt, nil
	}
}

// insertMessageText - Insert message text details in database
func (m *MessageService) insertMessageText(ctx context.Context, stmt *sql.Stmt, tx *sql.Tx, msgtxt *MessageText, userEmail string, requestID string) error {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6342}).Error(err)
		return err
	default:
		res, err := tx.StmtContext(ctx, stmt).Exec(
			msgtxt.UUID4,
			msgtxt.Mtext,
			msgtxt.WorkspaceID,
			msgtxt.ChannelID,
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
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6343}).Error(err)
			return err
		}
		uID, err := res.LastInsertId()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6344}).Error(err)
			return err
		}
		msgtxt.ID = uint(uID)
		return nil
	}
}

//CreateMessageAttachment - Create message Attachment
func (m *MessageService) createMessageAttachment(ctx context.Context, stmt *sql.Stmt, tx *sql.Tx, form *Message, messageID uint, userID uint, userEmail string, requestID string) (*MessageAttachment, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6345}).Error(err)
		return nil, err
	default:
		tn, tnday, tnweek, tnmonth, tnyear := common.GetTimeDetails()
		msgath := MessageAttachment{}
		uuid4, err := common.GetUUIDBytes()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6346}).Error(err)
			return nil, err
		}
		msgath.UUID4 = uuid4
		msgath.Mattach = form.Mattach
		msgath.WorkspaceID = form.WorkspaceID
		msgath.ChannelID = form.ChannelID
		msgath.MessageID = messageID
		msgath.UserID = userID
		msgath.UgroupID = form.UgroupID
		/*  StatusDates  */
		msgath.Statusc = common.Active
		msgath.CreatedAt = tn
		msgath.UpdatedAt = tn
		msgath.CreatedDay = tnday
		msgath.CreatedWeek = tnweek
		msgath.CreatedMonth = tnmonth
		msgath.CreatedYear = tnyear
		msgath.UpdatedDay = tnday
		msgath.UpdatedWeek = tnweek
		msgath.UpdatedMonth = tnmonth
		msgath.UpdatedYear = tnyear

		err = m.insertMessageAttachment(ctx, stmt, tx, &msgath, userEmail, requestID)

		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6347}).Error(err)
			return nil, err
		}
		return &msgath, nil
	}
}

// insertMessageAttachmentPrepare - Insert message attachment Prepare statement
func (m *MessageService) insertMessageAttachmentPrepare(ctx context.Context, userEmail string, requestID string) (*sql.Stmt, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6348}).Error(err)
		return nil, err
	default:
		db := m.DBService.DB
		stmt, err := db.PrepareContext(ctx, `insert into message_attachments
	  ( 
      uuid4,
			mattach,
			workspace_id,
			channel_id,
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
					?,?,?,?,?,?,?,?);`)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6349}).Error(err)
			return nil, err
		}
		return stmt, nil
	}
}

// insertMessageAttachment - Insert message attachment details in database
func (m *MessageService) insertMessageAttachment(ctx context.Context, stmt *sql.Stmt, tx *sql.Tx, msgath *MessageAttachment, userEmail string, requestID string) error {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6350}).Error(err)
		return err
	default:
		res, err := tx.StmtContext(ctx, stmt).Exec(
			msgath.UUID4,
			msgath.Mattach,
			msgath.WorkspaceID,
			msgath.ChannelID,
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
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6351}).Error(err)
			return err
		}
		uID, err := res.LastInsertId()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6352}).Error(err)
			return err
		}
		msgath.ID = uint(uID)
		return nil
	}
}

// createUserReply - create user reply
func (m *MessageService) createUserReply(ctx context.Context, stmt *sql.Stmt, tx *sql.Tx, channelID uint, messageID uint, userID uint, ugroupID uint, userEmail string, requestID string) error {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6353}).Error(err)
		return err
	default:
		tn, tnday, tnweek, tnmonth, tnyear := common.GetTimeDetails()
		ur := UserReply{}
		uuid4, err := common.GetUUIDBytes()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6354}).Error(err)
			err = tx.Rollback()
			return err
		}
		ur.UUID4 = uuid4
		ur.ChannelID = channelID
		ur.MessageID = messageID
		ur.UserID = userID
		ur.UgroupID = ugroupID
		/*  StatusDates  */
		ur.Statusc = common.Active
		ur.CreatedAt = tn
		ur.UpdatedAt = tn
		ur.CreatedDay = tnday
		ur.CreatedWeek = tnweek
		ur.CreatedMonth = tnmonth
		ur.CreatedYear = tnyear
		ur.UpdatedDay = tnday
		ur.UpdatedWeek = tnweek
		ur.UpdatedMonth = tnmonth
		ur.UpdatedYear = tnyear

		err = m.insertUserReply(ctx, stmt, tx, &ur, userEmail, requestID)

		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6355}).Error(err)
			err = tx.Rollback()
			return err
		}
		return nil
	}
}

// updateNumMessagesPrepare - UpdateNumMessages prepare statement
func (m *MessageService) updateNumMessagesPrepare(ctx context.Context, userEmail string, requestID string) (*sql.Stmt, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6431}).Error(err)
		return nil, err
	default:
		db := m.DBService.DB
		stmt, err := db.PrepareContext(ctx, `update channels set 
		  num_messages = ?,
			updated_at = ?, 
			updated_day = ?, 
			updated_week = ?, 
			updated_month = ?, 
			updated_year = ? where id = ? and statusc = ?;`)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6432}).Error(err)
			return nil, err
		}
		return stmt, nil
	}
}

// updateNumMessages - update number of messages in channels
func (m *MessageService) updateNumMessages(ctx context.Context, stmt *sql.Stmt, tx *sql.Tx, numMessages uint, ID uint, userEmail string, requestID string) error {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6433}).Error(err)
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
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6434}).Error(err)
			return err
		}
		return nil
	}
}

// insertUserReplyPrepare - Insert user reply Prepare statement
func (m *MessageService) insertUserReplyPrepare(ctx context.Context, userEmail string, requestID string) (*sql.Stmt, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6356}).Error(err)
		return nil, err
	default:
		db := m.DBService.DB
		stmt, err := db.PrepareContext(ctx, `insert into user_replies
	  ( 
      uuid4,
			channel_id,
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
					?,?,?,?,?,?);`)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6357}).Error(err)
			return nil, err
		}
		return stmt, nil
	}
}

// insertUserReply - Insert user reply details into database
func (m *MessageService) insertUserReply(ctx context.Context, stmt *sql.Stmt, tx *sql.Tx, ur *UserReply, userEmail string, requestID string) error {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6358}).Error(err)
		return err
	default:
		res, err := tx.StmtContext(ctx, stmt).Exec(
			ur.UUID4,
			ur.ChannelID,
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
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6359}).Error(err)
			return err
		}
		uID, err := res.LastInsertId()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6360}).Error(err)
			return err
		}
		ur.ID = uint(uID)
		return nil
	}
}

// CreateUserLike - Create user likes messages
func (m *MessageService) CreateUserLike(ctx context.Context, form *UserLike, UserID string, userEmail string, requestID string) (*UserLike, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6361}).Error(err)
		return nil, err
	default:
		userserv := &userservices.UserService{DBService: m.DBService, RedisService: m.RedisService}
		user, err := userserv.GetUser(ctx, UserID, userEmail, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6362}).Error(err)
			return nil, err
		}
		db := m.DBService.DB
		insertUserLikeStmt, err := m.insertUserLikePrepare(ctx, userEmail, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6363}).Error(err)
			return nil, err
		}

		tn, tnday, tnweek, tnmonth, tnyear := common.GetTimeDetails()

		ul := UserLike{}
		ul.UUID4, err = common.GetUUIDBytes()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6364}).Error(err)
			return nil, err
		}
		ul.ChannelID = form.ChannelID
		ul.MessageID = form.MessageID
		ul.UgroupID = form.UgroupID
		ul.UserID = user.ID
		/*  StatusDates  */
		ul.Statusc = common.Active
		ul.CreatedAt = tn
		ul.UpdatedAt = tn
		ul.CreatedDay = tnday
		ul.CreatedWeek = tnweek
		ul.CreatedMonth = tnmonth
		ul.CreatedYear = tnyear
		ul.UpdatedDay = tnday
		ul.UpdatedWeek = tnweek
		ul.UpdatedMonth = tnmonth
		ul.UpdatedYear = tnyear

		tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6365}).Error(err)
			return nil, err
		}

		err = m.insertUserLike(ctx, insertUserLikeStmt, tx, &ul, userEmail, requestID)

		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6366}).Error(err)
			err = tx.Rollback()
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6367}).Error(err)
				return nil, err
			}
			err = insertUserLikeStmt.Close()
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6368}).Error(err)
				return nil, err
			}
			return nil, err
		}

		err = tx.Commit()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6369}).Error(err)
			return nil, err
		}
		err = insertUserLikeStmt.Close()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6370}).Error(err)
			return nil, err
		}
		return &ul, nil
	}
}

// insertUserLikePrepare - Insert User like Prepare statement
func (m *MessageService) insertUserLikePrepare(ctx context.Context, userEmail string, requestID string) (*sql.Stmt, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6371}).Error(err)
		return nil, err
	default:
		db := m.DBService.DB
		stmt, err := db.PrepareContext(ctx, `insert into user_likes
	  ( 
      uuid4,
			channel_id,
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
					?,?,?,?,?,?);`)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6372}).Error(err)
			return nil, err
		}
		return stmt, nil
	}
}

// insertUserLike - Insert User like details in database
func (m *MessageService) insertUserLike(ctx context.Context, stmt *sql.Stmt, tx *sql.Tx, ur *UserLike, userEmail string, requestID string) error {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6373}).Error(err)
		return err
	default:
		res, err := tx.StmtContext(ctx, stmt).Exec(
			ur.UUID4,
			ur.ChannelID,
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
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6374}).Error(err)
			return err
		}
		uID, err := res.LastInsertId()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6375}).Error(err)
			return err
		}
		ur.ID = uint(uID)
		return nil
	}
}

// CreateUserVote - Create User Vote
func (m *MessageService) CreateUserVote(ctx context.Context, form *UserVote, UserID string, userEmail string, requestID string) (*UserVote, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6376}).Error(err)
		return nil, err
	default:
		userserv := &userservices.UserService{DBService: m.DBService, RedisService: m.RedisService}
		user, err := userserv.GetUser(ctx, UserID, userEmail, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6377}).Error(err)
			return nil, err
		}
		db := m.DBService.DB
		insertUserVoteStmt, err := m.insertUserVotePrepare(ctx, userEmail, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6378}).Error(err)
			return nil, err
		}

		tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6379}).Error(err)
			return nil, err
		}

		tn, tnday, tnweek, tnmonth, tnyear := common.GetTimeDetails()

		ul := UserVote{}
		ul.UUID4, err = common.GetUUIDBytes()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6380}).Error(err)
			err = tx.Rollback()
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6381}).Error(err)
				return nil, err
			}
			err = insertUserVoteStmt.Close()
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6382}).Error(err)
				return nil, err
			}
			return nil, err
		}
		ul.ChannelID = form.ChannelID
		ul.MessageID = form.MessageID
		ul.UgroupID = form.UgroupID
		ul.UserID = user.ID
		ul.Vote = form.Vote
		/*  StatusDates  */
		ul.Statusc = common.Active
		ul.CreatedAt = tn
		ul.UpdatedAt = tn
		ul.CreatedDay = tnday
		ul.CreatedWeek = tnweek
		ul.CreatedMonth = tnmonth
		ul.CreatedYear = tnyear
		ul.UpdatedDay = tnday
		ul.UpdatedWeek = tnweek
		ul.UpdatedMonth = tnmonth
		ul.UpdatedYear = tnyear

		err = m.insertUserVote(ctx, insertUserVoteStmt, tx, &ul, userEmail, requestID)

		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6383}).Error(err)
			err = tx.Rollback()
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6384}).Error(err)
				return nil, err
			}
			err = insertUserVoteStmt.Close()
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6385}).Error(err)
				return nil, err
			}
			return nil, err
		}

		err = tx.Commit()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6386}).Error(err)
			err = tx.Rollback()
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6387}).Error(err)
				return nil, err
			}
			return nil, err
		}
		err = insertUserVoteStmt.Close()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6388}).Error(err)
			return nil, err
		}
		return &ul, nil
	}
}

// insertUserVotePrepare - Insert User vote Prepare Statement
func (m *MessageService) insertUserVotePrepare(ctx context.Context, userEmail string, requestID string) (*sql.Stmt, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6389}).Error(err)
		return nil, err
	default:
		db := m.DBService.DB
		stmt, err := db.PrepareContext(ctx, `insert into user_votes
	  ( 
      uuid4,
			channel_id,
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
					?,?,?,?,?,?,?);`)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6390}).Error(err)
			return nil, err
		}
		return stmt, nil
	}

}

// insertUserVote - Insert User vote details into database
func (m *MessageService) insertUserVote(ctx context.Context, stmt *sql.Stmt, tx *sql.Tx, ur *UserVote, userEmail string, requestID string) error {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6391}).Error(err)
		return err
	default:
		res, err := tx.StmtContext(ctx, stmt).Exec(
			ur.UUID4,
			ur.ChannelID,
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
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6392}).Error(err)
			return err
		}
		uID, err := res.LastInsertId()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6393}).Error(err)
			return err
		}
		ur.ID = uint(uID)
		return nil
	}
}

// GetMessage - Get message
func (m *MessageService) GetMessage(ctx context.Context, ID string, userEmail string, requestID string) (*Message, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6394}).Error(err)
		return nil, err
	default:
		uuid4byte, err := common.UUIDStrToBytes(ID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6395}).Error(err)
			return nil, err
		}
		msg := Message{}
		db := m.DBService.DB
		row := db.QueryRowContext(ctx, `select
      id,
 			uuid4,
			num_likes,
			num_upvotes,
			num_downvotes,
			workspace_id,
			channel_id,
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
			updated_year from messages where uuid4 = ? and statusc = ? ;`, uuid4byte, common.Active)

		err = row.Scan(
			&msg.ID,
			&msg.UUID4,
			&msg.NumLikes,
			&msg.NumUpvotes,
			&msg.NumDownvotes,
			&msg.WorkspaceID,
			&msg.ChannelID,
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
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6396}).Error(err)
			return nil, err
		}
		uuid4Str, err := common.UUIDBytesToStr(msg.UUID4)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6397}).Error(err)
			return nil, err
		}
		msg.IDS = uuid4Str

		var isPresent bool
		row = db.QueryRowContext(ctx, `select exists (select 1 from message_texts where message_id = ?);`, msg.ID)
		err = row.Scan(&isPresent)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6398}).Error(err)
		}
		if isPresent {
			messageTexts, err := m.GetMessagesTexts(ctx, msg.ID, userEmail, requestID)
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6399}).Error(err)
				return nil, err
			}
			msg.MessageTexts = messageTexts
		}

		var isPresent1 bool
		row1 := db.QueryRowContext(ctx, `select exists (select 1 from message_attachments where message_id = ?);`, msg.ID)
		err = row1.Scan(&isPresent1)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6400}).Error(err)
		}
		if isPresent1 {
			messageAttachments, err := m.GetMessageAttachments(ctx, msg.ID, userEmail, requestID)
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6401}).Error(err)
				return nil, err
			}
			msg.MessageAttachments = messageAttachments
		}

		return &msg, nil
	}
}

// GetMessagesWithTextAttach - Get messages with attachements
func (m *MessageService) GetMessagesWithTextAttach(ctx context.Context, messages []*Message, userEmail string, requestID string) ([]*Message, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6402}).Error(err)
		return nil, err
	default:
		db := m.DBService.DB
		pohs := []*Message{}

		for _, message := range messages {
			var isPresent bool
			row := db.QueryRowContext(ctx, `select exists (select 1 from message_texts where message_id = ?);`, message.ID)
			err := row.Scan(&isPresent)
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6403}).Error(err)
			}
			if isPresent {
				messageTexts, err := m.GetMessagesTexts(ctx, message.ID, userEmail, requestID)
				if err != nil {
					log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6404}).Error(err)
					return nil, err
				}
				message.MessageTexts = messageTexts
			}

			var isPresent1 bool
			row1 := db.QueryRowContext(ctx, `select exists (select 1 from message_attachments where message_id = ?);`, message.ID)
			err = row1.Scan(&isPresent1)
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6405}).Error(err)
			}
			if isPresent1 {
				messageAttachments, err := m.GetMessageAttachments(ctx, message.ID, userEmail, requestID)
				if err != nil {
					log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6406}).Error(err)
					return nil, err
				}
				message.MessageAttachments = messageAttachments
			}

			pohs = append(pohs, message)
		}

		return pohs, nil
	}
}

// GetMessagesTexts - get message texts
func (m *MessageService) GetMessagesTexts(ctx context.Context, messageID uint, userEmail string, requestID string) ([]*MessageText, error) {
	db := m.DBService.DB
	mtexts := []*MessageText{}
	rows, err := db.QueryContext(ctx, `select 
        id,
        uuid4,
				mtext,
				workspace_id,
				channel_id,
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
				updated_year from message_texts where message_id = ? and statusc = ?`, messageID, common.Active)

	if err != nil {
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6407}).Error(err)
		return nil, err
	}
	for rows.Next() {
		msgtxt := MessageText{}
		err = rows.Scan(
			&msgtxt.ID,
			&msgtxt.UUID4,
			&msgtxt.Mtext,
			&msgtxt.WorkspaceID,
			&msgtxt.ChannelID,
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
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6408}).Error(err)
		}

		mtexts = append(mtexts, &msgtxt)
	}

	err = rows.Close()
	if err != nil {
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6409}).Error(err)
		return nil, err
	}
	return mtexts, nil
}

// GetMessageAttachments - get message attachements
func (m *MessageService) GetMessageAttachments(ctx context.Context, messageID uint, userEmail string, requestID string) ([]*MessageAttachment, error) {
	db := m.DBService.DB
	messageAttachments := []*MessageAttachment{}
	rows, err := db.QueryContext(ctx, `select 
        id,
        uuid4,
				mattach,
				workspace_id,
				channel_id,
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
				updated_year from message_attachments where message_id = ? and statusc = ?`, messageID, common.Active)

	if err != nil {
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6410}).Error(err)
	}
	for rows.Next() {
		msgath := MessageAttachment{}
		err = rows.Scan(
			&msgath.ID,
			&msgath.UUID4,
			&msgath.Mattach,
			&msgath.WorkspaceID,
			&msgath.ChannelID,
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
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6411}).Error(err)
		}

		messageAttachments = append(messageAttachments, &msgath)
	}

	err = rows.Close()
	if err != nil {
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6412}).Error(err)
		return nil, err
	}
	return messageAttachments, nil
}

//UpdateMessage - Update message
func (m *MessageService) UpdateMessage(ctx context.Context, ID string, form *Message, UserID string, userEmail string, requestID string) error {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6413}).Error(err)
		return err
	default:
		msg, err := m.GetMessage(ctx, ID, userEmail, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6414}).Error(err)
			return err
		}

		db := m.DBService.DB

		tn, tnday, tnweek, tnmonth, tnyear := common.GetTimeDetails()
		stmt, err := db.PrepareContext(ctx, `update message_texts set 
		  mtext = ?,
			updated_at = ?, 
			updated_day = ?, 
			updated_week = ?, 
			updated_month = ?, 
			updated_year = ? where message_id = ? and statusc = ?;`)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6415}).Error(err)
			return err
		}

		tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6416}).Error(err)
			return err
		}

		_, err = tx.StmtContext(ctx, stmt).Exec(
			form.Mtext,
			tn,
			tnday,
			tnweek,
			tnmonth,
			tnyear,
			msg.ID,
			common.Active)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6417}).Error(err)
			err = tx.Rollback()
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6418}).Error(err)
				return err
			}
			err = stmt.Close()
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6419}).Error(err)
				return err
			}
			return err
		}

		err = tx.Commit()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6420}).Error(err)
			return err
		}

		err = stmt.Close()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6421}).Error(err)
			return err
		}
		return nil
	}
}

// DeleteMessage - Delete message
func (m *MessageService) DeleteMessage(ctx context.Context, ID string, userEmail string, requestID string) error {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6422}).Error(err)
		return err
	default:
		uuid4byte, err := common.UUIDStrToBytes(ID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6423}).Error(err)
			return err
		}
		db := m.DBService.DB
		tn, tnday, tnweek, tnmonth, tnyear := common.GetTimeDetails()
		stmt, err := db.PrepareContext(ctx, `update messages set 
		  statusc = ?,
			updated_at = ?, 
			updated_day = ?, 
			updated_week = ?, 
			updated_month = ?, 
			updated_year = ? where uuid4= ?;`)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6424}).Error(err)
			return err
		}

		tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6425}).Error(err)
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
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6426}).Error(err)
			err = tx.Rollback()
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6427}).Error(err)
				return err
			}
			err = stmt.Close()
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6428}).Error(err)
				return err
			}
			return err
		}

		err = tx.Commit()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6429}).Error(err)
			err = tx.Rollback()
			return err
		}
		err = stmt.Close()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6430}).Error(err)
			return err
		}
		return nil
	}
}
