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

// For validation of Channel fields
const (
	ChannelNameLenMin = 1
	ChannelNameLenMax = 50
	ChannelDescLenMin = 1
	ChannelDescLenMax = 1000
)

// Channel - Channel view representation
type Channel struct {
	ID    uint   `json:"id,omitempty"`
	UUID4 []byte `json:"-"`
	IDS   string `json:"id_s,omitempty"`

	ChannelName string `json:"channel_name,omitempty"`
	ChannelDesc string `json:"channel_desc,omitempty"`
	NumTags     uint   `json:"num_tags,omitempty"`
	Tag1        string `json:"tag1,omitempty"`
	Tag2        string `json:"tag2,omitempty"`
	Tag3        string `json:"tag3,omitempty"`
	Tag4        string `json:"tag4,omitempty"`
	Tag5        string `json:"tag5,omitempty"`
	Tag6        string `json:"tag6,omitempty"`
	Tag7        string `json:"tag7,omitempty"`
	Tag8        string `json:"tag8,omitempty"`
	Tag9        string `json:"tag9,omitempty"`
	Tag10       string `json:"tag10,omitempty"`

	NumViews    uint `json:"num_views,omitempty"`
	NumMessages uint `json:"num_messages,omitempty"`

	WorkspaceID uint `json:"workspace_id,omitempty"`
	UserID      uint `json:"user_id,omitempty"`
	UgroupID    uint `json:"ugroup_id,omitempty"`

	common.StatusDates
	Messages []*Message

	//only for logic purpose to create message with channel
	Mtext   string `json:"-"`
	Mattach string `json:"-"`
}

// ChannelsUser - ChannelsUser view representation
type ChannelsUser struct {
	ID          uint   `json:"id,omitempty"`
	UUID4       []byte `json:"-"`
	IDS         string `json:"id_s,omitempty"`
	ChannelID   uint   `json:"channel_id,omitempty"`
	NumMessages uint   `json:"num_messages,omitempty"`
	NumViews    uint   `json:"num_views,omitempty"`
	UserID      uint   `json:"user_id,omitempty"`
	UgroupID    uint   `json:"ugroup_id,omitempty"`

	common.StatusDates
}

// UserChannel - UserChannel view representation
type UserChannel struct {
	ID        uint   `json:"id,omitempty"`
	UUID4     []byte `json:"-"`
	ChannelID uint   `json:"channel_id,omitempty"`
	UserID    uint   `json:"user_id,omitempty"`
	UgroupID  uint   `json:"ugroup_id,omitempty"`

	common.StatusDates
}

// ChannelServiceIntf - interface for Channel Service
type ChannelServiceIntf interface {
	CreateChannel(ctx context.Context, form *Channel, UserID string, userEmail string, requestID string) (*Channel, error)
	ShowChannel(ctx context.Context, ID string, UserID string, userEmail string, requestID string) (*Channel, error)
	GetChannelByID(ctx context.Context, ID uint, userEmail string, requestID string) (*Channel, error)
	GetChannel(ctx context.Context, ID string, userEmail string, requestID string) (*Channel, error)
	GetChannelByName(ctx context.Context, channelname string, userEmail string, requestID string) (*Channel, error)
	GetChannelWithMessages(ctx context.Context, ID string, userEmail string, requestID string) (*Channel, error)
	GetChannelMessages(ctx context.Context, uuid4byte []byte, userEmail string, requestID string) (*Channel, error)
	GetChannelsUser(ctx context.Context, ID uint, UserID uint, userEmail string, requestID string) (*ChannelsUser, error)
	UpdateChannel(ctx context.Context, ID string, form *Channel, UserID string, userEmail string, requestID string) error
	DeleteChannel(ctx context.Context, ID string, userEmail string, requestID string) error
}

// ChannelService - For accessing channel services
type ChannelService struct {
	DBService    *common.DBService
	RedisService *common.RedisService
}

// NewChannelService - Create channel service
func NewChannelService(dbOpt *common.DBService, redisOpt *common.RedisService) *ChannelService {
	return &ChannelService{
		DBService:    dbOpt,
		RedisService: redisOpt,
	}
}

// CreateChannel - Create channel
func (t *ChannelService) CreateChannel(ctx context.Context, form *Channel, UserID string, userEmail string, requestID string) (*Channel, error) {
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
		workspaceserv := &WorkspaceService{DBService: t.DBService, RedisService: t.RedisService}
		workspace, err := workspaceserv.GetWorkspaceByID(ctx, form.WorkspaceID, userEmail, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5326}).Error(err)
			return nil, err
		}

		insertChannelStmt, err := t.insertChannelPrepare(ctx, userEmail, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5323}).Error(err)
			return nil, err
		}
		updateNumChannelsStmt, err := t.updateNumChannelsPrepare(ctx, userEmail, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5323}).Error(err)
			return nil, err
		}
		insertUserChannelStmt, err := t.insertUserChannelPrepare(ctx, userEmail, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5323}).Error(err)
			return nil, err
		}
		tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5323}).Error(err)
			return nil, err
		}

		channel, err := t.createChannel(ctx, insertChannelStmt, tx, form, user.ID, userEmail, requestID)

		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5325}).Error(err)
			err = tx.Rollback()
			return nil, err
		}

		numchannels := workspace.NumChannels + 1
		err = t.updateNumChannels(ctx, updateNumChannelsStmt, tx, numchannels, workspace.ID, userEmail, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5381}).Error(err)
			err = tx.Rollback()
			return nil, err
		}
		if form.Mtext != "" {
			msgserv := &MessageService{DBService: t.DBService, RedisService: t.RedisService}
			msgform := Message{}
			msgform.WorkspaceID = workspace.ID
			msgform.ChannelID = channel.ID
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

		err = t.createUserChannel(ctx, insertUserChannelStmt, tx, channel.ID, user.ID, userEmail, requestID)

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
		return channel, nil
	}
}

// create channel - create channel
func (t *ChannelService) createChannel(ctx context.Context, stmt *sql.Stmt, tx *sql.Tx, form *Channel, userID uint, userEmail string, requestID string) (*Channel, error) {
	var err error
	tn, tnday, tnweek, tnmonth, tnyear := common.GetTimeDetails()

	channel := Channel{}
	channel.UUID4, err = common.GetUUIDBytes()
	if err != nil {
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5324}).Error(err)
		return nil, err
	}
	channel.ChannelName = form.ChannelName
	channel.ChannelDesc = form.ChannelDesc
	channel.NumTags = form.NumTags
	channel.Tag1 = form.Tag1
	channel.Tag2 = form.Tag2
	channel.Tag3 = form.Tag3
	channel.Tag4 = form.Tag4
	channel.Tag5 = form.Tag5
	channel.Tag6 = form.Tag6
	channel.Tag7 = form.Tag7
	channel.Tag8 = form.Tag8
	channel.Tag9 = form.Tag9
	channel.Tag10 = form.Tag10
	channel.NumViews = uint(0)
	channel.NumMessages = uint(0)
	channel.WorkspaceID = form.WorkspaceID
	channel.UserID = userID
	channel.UgroupID = form.UgroupID
	/*  StatusDates  */
	channel.Statusc = common.Active
	channel.CreatedAt = tn
	channel.UpdatedAt = tn
	channel.CreatedDay = tnday
	channel.CreatedWeek = tnweek
	channel.CreatedMonth = tnmonth
	channel.CreatedYear = tnyear
	channel.UpdatedDay = tnday
	channel.UpdatedWeek = tnweek
	channel.UpdatedMonth = tnmonth
	channel.UpdatedYear = tnyear

	err = t.insertChannel(ctx, stmt, tx, &channel, userEmail, requestID)

	if err != nil {
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5325}).Error(err)
		return nil, err
	}
	return &channel, nil
}

// insertChannelPrepare - Insert channel Prepare Statement
func (t *ChannelService) insertChannelPrepare(ctx context.Context, userEmail string, requestID string) (*sql.Stmt, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5342}).Error(err)
		return nil, err
	default:
		db := t.DBService.DB
		stmt, err := db.PrepareContext(ctx, `insert into channels
	  ( uuid4,
			channel_name,
			channel_desc,
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
			workspace_id,
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

// insertChannel - Insert channel details into database
func (t *ChannelService) insertChannel(ctx context.Context, stmt *sql.Stmt, tx *sql.Tx, channel *Channel, userEmail string, requestID string) error {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5342}).Error(err)
		return err
	default:
		res, err := tx.StmtContext(ctx, stmt).Exec(
			channel.UUID4,
			channel.ChannelName,
			channel.ChannelDesc,
			channel.NumTags,
			channel.Tag1,
			channel.Tag2,
			channel.Tag3,
			channel.Tag4,
			channel.Tag5,
			channel.Tag6,
			channel.Tag7,
			channel.Tag8,
			channel.Tag9,
			channel.Tag10,
			channel.NumViews,
			channel.NumMessages,
			channel.WorkspaceID,
			channel.UserID,
			channel.UgroupID,
			/*  StatusDates  */
			channel.Statusc,
			channel.CreatedAt,
			channel.UpdatedAt,
			channel.CreatedDay,
			channel.CreatedWeek,
			channel.CreatedMonth,
			channel.CreatedYear,
			channel.UpdatedDay,
			channel.UpdatedWeek,
			channel.UpdatedMonth,
			channel.UpdatedYear)

		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5344}).Error(err)
			return err
		}
		uID, err := res.LastInsertId()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5345}).Error(err)
			return err
		}
		channel.ID = uint(uID)
		uuid4Str, err := common.UUIDBytesToStr(channel.UUID4)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5346}).Error(err)
			return err
		}
		channel.IDS = uuid4Str
		return nil
	}
}

// updateNumChannelsPrepare - UpdateNumChannels Prepare Statement
func (t *ChannelService) updateNumChannelsPrepare(ctx context.Context, userEmail string, requestID string) (*sql.Stmt, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5391}).Error(err)
		return nil, err
	default:
		db := t.DBService.DB
		stmt, err := db.PrepareContext(ctx, `update workspaces set 
    num_channels = ?,
	  updated_at = ?, 
		updated_day = ?, 
		updated_week = ?, 
		updated_month = ?, 
		updated_year = ? where id = ? and statusc = ?;`)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5392}).Error(err)
			return nil, err
		}
		return stmt, nil
	}
}

// updateNumChannels - update number of channels in workspace
func (t *ChannelService) updateNumChannels(ctx context.Context, stmt *sql.Stmt, tx *sql.Tx, numChannels uint, ID uint, userEmail string, requestID string) error {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5393}).Error(err)
		return err
	default:
		tn, tnday, tnweek, tnmonth, tnyear := common.GetTimeDetails()

		_, err := tx.StmtContext(ctx, stmt).Exec(
			numChannels,
			tn,
			tnday,
			tnweek,
			tnmonth,
			tnyear,
			ID,
			common.Active)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5394}).Error(err)
			return err
		}
		return nil
	}
}

// createUserChannel - create user channel
func (t *ChannelService) createUserChannel(ctx context.Context, stmt *sql.Stmt, tx *sql.Tx, channelID uint, userID uint, userEmail string, requestID string) error {
	var err error
	tn, tnday, tnweek, tnmonth, tnyear := common.GetTimeDetails()
	uc := UserChannel{}
	uc.UUID4, err = common.GetUUIDBytes()
	if err != nil {
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5339}).Error(err)
		return err
	}
	uc.ChannelID = channelID
	uc.UserID = userID
	uc.UgroupID = uint(0)
	/*  StatusDates  */
	uc.Statusc = common.Active
	uc.CreatedAt = tn
	uc.UpdatedAt = tn
	uc.CreatedDay = tnday
	uc.CreatedWeek = tnweek
	uc.CreatedMonth = tnmonth
	uc.CreatedYear = tnyear
	uc.UpdatedDay = tnday
	uc.UpdatedWeek = tnweek
	uc.UpdatedMonth = tnmonth
	uc.UpdatedYear = tnyear

	err = t.insertUserChannel(ctx, stmt, tx, &uc, userEmail, requestID)

	if err != nil {
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5340}).Error(err)
		return err
	}
	return nil
}

// insertUserChannelPrepare - Insert user channels Prepare Statement
func (t *ChannelService) insertUserChannelPrepare(ctx context.Context, userEmail string, requestID string) (*sql.Stmt, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5367}).Error(err)
		return nil, err
	default:
		db := t.DBService.DB
		stmt, err := db.PrepareContext(ctx, `insert into user_channels
	  (
    uuid4,
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
					?,?,?,?,?);`)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5368}).Error(err)
			return nil, err
		}
		return stmt, nil
	}
}

// insertUserChannel - Insert user channels details into database
func (t *ChannelService) insertUserChannel(ctx context.Context, stmt *sql.Stmt, tx *sql.Tx, channel *UserChannel, userEmail string, requestID string) error {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5367}).Error(err)
		return err
	default:
		res, err := tx.StmtContext(ctx, stmt).Exec(
			channel.UUID4,
			channel.ChannelID,
			channel.UserID,
			channel.UgroupID,
			channel.Statusc,
			channel.CreatedAt,
			channel.UpdatedAt,
			channel.CreatedDay,
			channel.CreatedWeek,
			channel.CreatedMonth,
			channel.CreatedYear,
			channel.UpdatedDay,
			channel.UpdatedWeek,
			channel.UpdatedMonth,
			channel.UpdatedYear)

		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5369}).Error(err)
			return err
		}
		uID, err := res.LastInsertId()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5370}).Error(err)
			return err
		}
		channel.ID = uint(uID)
		return nil
	}
}

// ShowChannel - Get channel details
func (t *ChannelService) ShowChannel(ctx context.Context, ID string, UserID string, userEmail string, requestID string) (*Channel, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5300}).Error(err)
		return nil, err
	default:
		db := t.DBService.DB
		channel, err := t.GetChannelWithMessages(ctx, ID, userEmail, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5301}).Error(err)
			return nil, err
		}
		//update channel_users table
		userserv := &userservices.UserService{DBService: t.DBService, RedisService: t.RedisService}
		user, err := userserv.GetUser(ctx, UserID, userEmail, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5302}).Error(err)
			return nil, err
		}
		var isPresent bool
		row := db.QueryRowContext(ctx, `select exists (select 1 from channels_users where channel_id = ? and user_id = ?);`, channel.ID, user.ID)
		err = row.Scan(&isPresent)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5303}).Error(err)
			return nil, err
		}

		updateChannelUsersStmt, insertChannelsUserStmt, err := t.showChannelPrepareStmts(ctx, userEmail, requestID)
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
			err = t.showChannelPrepareStmtsClose(ctx, updateChannelUsersStmt, insertChannelsUserStmt, userEmail, requestID)
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5379}).Error(err)
				return nil, err
			}
			return nil, err
		}
		err = t.showChannelUpdateChannelUsers(ctx, updateChannelUsersStmt, insertChannelsUserStmt, channel, tx, user.ID, isPresent, userEmail, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5372}).Error(err)
			err = tx.Rollback()
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5379}).Error(err)
				return nil, err
			}
			err = t.showChannelPrepareStmtsClose(ctx, updateChannelUsersStmt, insertChannelsUserStmt, userEmail, requestID)
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
		err = t.showChannelPrepareStmtsClose(ctx, updateChannelUsersStmt, insertChannelsUserStmt, userEmail, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5379}).Error(err)
			return nil, err
		}
		return channel, nil
	}
}

//showChannelPrepareStmts - Prepare Statements
func (t *ChannelService) showChannelPrepareStmts(ctx context.Context, userEmail string, requestID string) (*sql.Stmt, *sql.Stmt, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 6311}).Error(err)
		return nil, nil, err
	default:
		updateChannelUsersStmt, err := t.updateChannelUsersPrepare(ctx, userEmail, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5379}).Error(err)
			return nil, nil, err
		}

		insertChannelsUserStmt, err := t.insertChannelsUserPrepare(ctx, userEmail, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5379}).Error(err)
			return nil, nil, err
		}

		return updateChannelUsersStmt, insertChannelsUserStmt, nil
	}
}

//showChannelPrepareStmtsClose - Close Prepare Statements
func (t *ChannelService) showChannelPrepareStmtsClose(ctx context.Context, updateChannelUsersStmt *sql.Stmt, insertChannelsUserStmt *sql.Stmt, userEmail string, requestID string) error {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5379}).Error(err)
		return err
	default:
		err := updateChannelUsersStmt.Close()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5379}).Error(err)
			return err
		}

		err = insertChannelsUserStmt.Close()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5379}).Error(err)
			return err
		}

		return nil
	}
}

// showChannelUpdateChannelUsers - update channel users details
func (t *ChannelService) showChannelUpdateChannelUsers(ctx context.Context, updateChannelUsersStmt *sql.Stmt, insertChannelsUserStmt *sql.Stmt, channel *Channel, tx *sql.Tx, userID uint, isPresent bool, userEmail string, requestID string) error {
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
			channelsuser, err := t.GetChannelsUser(ctx, channel.ID, userID, userEmail, requestID)
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5373}).Error(err)
				return err
			}

			numViews := channelsuser.NumViews + 1
			err = t.updateChannelUsers(ctx, updateChannelUsersStmt, tx, channel.NumMessages, numViews, channelsuser.ID, userEmail, requestID)
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5379}).Error(err)
				return err
			}
		} else {
			//create
			cu := ChannelsUser{}
			cu.UUID4, err = common.GetUUIDBytes()
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5307}).Error(err)
				return err
			}
			cu.ChannelID = channel.ID
			cu.NumMessages = channel.NumMessages
			cu.NumViews = 1
			cu.UserID = userID
			cu.UgroupID = uint(0)
			cu.Statusc = common.Active
			cu.CreatedAt = tn
			cu.UpdatedAt = tn
			cu.CreatedDay = tnday
			cu.CreatedWeek = tnweek
			cu.CreatedMonth = tnmonth
			cu.CreatedYear = tnyear
			cu.UpdatedDay = tnday
			cu.UpdatedWeek = tnweek
			cu.UpdatedMonth = tnmonth
			cu.UpdatedYear = tnyear

			_, err := t.insertChannelsUser(ctx, insertChannelsUserStmt, tx, cu, userEmail, requestID)

			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5308}).Error(err)
				return err
			}

		}
		return nil
	}
}

// updateChannelUsersPrepare - update channel users prepare statement
func (t *ChannelService) updateChannelUsersPrepare(ctx context.Context, userEmail string, requestID string) (*sql.Stmt, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5332}).Error(err)
		return nil, err
	default:
		db := t.DBService.DB
		stmt, err := db.PrepareContext(ctx, `update channels_users set 
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

// updateChannelUsers - update channel users
func (t *ChannelService) updateChannelUsers(ctx context.Context, stmt *sql.Stmt, tx *sql.Tx, numMessages uint, numViews uint, channelsuserID uint, userEmail string, requestID string) error {
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
			channelsuserID,
			common.Active)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5305}).Error(err)
			return err
		}
		return nil
	}
}

// insertChannelsUserPrepare - Insert channel user Prepare Statement
func (t *ChannelService) insertChannelsUserPrepare(ctx context.Context, userEmail string, requestID string) (*sql.Stmt, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5361}).Error(err)
		return nil, err
	default:
		db := t.DBService.DB
		stmt, err := db.PrepareContext(ctx, `insert into channels_users
	  (uuid4,
		channel_id,
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

// insertChannelsUser - Insert channel user details into database
func (t *ChannelService) insertChannelsUser(ctx context.Context, stmt *sql.Stmt, tx *sql.Tx, channelUser ChannelsUser, userEmail string, requestID string) (*ChannelsUser, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5361}).Error(err)
		return nil, err
	default:

		res, err := tx.StmtContext(ctx, stmt).Exec(
			channelUser.UUID4,
			channelUser.ChannelID,
			channelUser.NumMessages,
			channelUser.NumViews,
			channelUser.UserID,
			channelUser.UgroupID,
			channelUser.Statusc,
			channelUser.CreatedAt,
			channelUser.UpdatedAt,
			channelUser.CreatedDay,
			channelUser.CreatedWeek,
			channelUser.CreatedMonth,
			channelUser.CreatedYear,
			channelUser.UpdatedDay,
			channelUser.UpdatedWeek,
			channelUser.UpdatedMonth,
			channelUser.UpdatedYear)

		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5363}).Error(err)
			return nil, err
		}
		uID, err := res.LastInsertId()
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5364}).Error(err)
			return nil, err
		}
		channelUser.ID = uint(uID)
		uuid4Str, err := common.UUIDBytesToStr(channelUser.UUID4)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5365}).Error(err)
			return nil, err
		}
		channelUser.IDS = uuid4Str
		return &channelUser, nil
	}
}

// GetChannelByID - Get channel by ID
func (t *ChannelService) GetChannelByID(ctx context.Context, ID uint, userEmail string, requestID string) (*Channel, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5348}).Error(err)
		return nil, err
	default:
		channel := Channel{}
		db := t.DBService.DB
		row := db.QueryRowContext(ctx, `select
    id,
		uuid4,
		channel_name,
		channel_desc,
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
		workspace_id,
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
		updated_year from channels where id = ? and statusc = ?`, ID, common.Active)

		err := row.Scan(
			&channel.ID,
			&channel.UUID4,
			&channel.ChannelName,
			&channel.ChannelDesc,
			&channel.NumTags,
			&channel.Tag1,
			&channel.Tag2,
			&channel.Tag3,
			&channel.Tag4,
			&channel.Tag5,
			&channel.Tag6,
			&channel.Tag7,
			&channel.Tag8,
			&channel.Tag9,
			&channel.Tag10,
			&channel.NumViews,
			&channel.NumMessages,
			&channel.WorkspaceID,
			&channel.UserID,
			&channel.UgroupID,
			/*  StatusDates  */
			&channel.Statusc,
			&channel.CreatedAt,
			&channel.UpdatedAt,
			&channel.CreatedDay,
			&channel.CreatedWeek,
			&channel.CreatedMonth,
			&channel.CreatedYear,
			&channel.UpdatedDay,
			&channel.UpdatedWeek,
			&channel.UpdatedMonth,
			&channel.UpdatedYear)

		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5349}).Error(err)
			return nil, err
		}
		uuid4Str, err := common.UUIDBytesToStr(channel.UUID4)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5350}).Error(err)
			return nil, err
		}
		channel.IDS = uuid4Str
		return &channel, nil
	}
}

// GetChannel - Get channel
func (t *ChannelService) GetChannel(ctx context.Context, ID string, userEmail string, requestID string) (*Channel, error) {
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
		channel := Channel{}
		db := t.DBService.DB
		row := db.QueryRowContext(ctx, `select
    id,
		uuid4,
		channel_name,
		channel_desc,
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
		workspace_id,
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
		updated_year from channels where uuid4 = ? and statusc = ?`, uuid4byte, common.Active)

		err = row.Scan(
			&channel.ID,
			&channel.UUID4,
			&channel.ChannelName,
			&channel.ChannelDesc,
			&channel.NumTags,
			&channel.Tag1,
			&channel.Tag2,
			&channel.Tag3,
			&channel.Tag4,
			&channel.Tag5,
			&channel.Tag6,
			&channel.Tag7,
			&channel.Tag8,
			&channel.Tag9,
			&channel.Tag10,
			&channel.NumViews,
			&channel.NumMessages,
			&channel.WorkspaceID,
			&channel.UserID,
			&channel.UgroupID,
			/*  StatusDates  */
			&channel.Statusc,
			&channel.CreatedAt,
			&channel.UpdatedAt,
			&channel.CreatedDay,
			&channel.CreatedWeek,
			&channel.CreatedMonth,
			&channel.CreatedYear,
			&channel.UpdatedDay,
			&channel.UpdatedWeek,
			&channel.UpdatedMonth,
			&channel.UpdatedYear)

		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5353}).Error(err)
			return nil, err
		}
		uuid4Str, err := common.UUIDBytesToStr(channel.UUID4)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5354}).Error(err)
			return nil, err
		}
		channel.IDS = uuid4Str
		return &channel, nil
	}
}

// GetChannelByName - Get channel by name
func (t *ChannelService) GetChannelByName(ctx context.Context, channelname string, userEmail string, requestID string) (*Channel, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5355}).Error(err)
		return nil, err
	default:
		channel := Channel{}
		db := t.DBService.DB
		row := db.QueryRowContext(ctx, `select
    id,
		uuid4,
		channel_name,
		channel_desc,
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
		workspace_id,
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
		updated_year from channels where channel_name = ? and statusc = ?`, channelname, common.Active)

		err := row.Scan(
			&channel.ID,
			&channel.UUID4,
			&channel.ChannelName,
			&channel.ChannelDesc,
			&channel.NumTags,
			&channel.Tag1,
			&channel.Tag2,
			&channel.Tag3,
			&channel.Tag4,
			&channel.Tag5,
			&channel.Tag6,
			&channel.Tag7,
			&channel.Tag8,
			&channel.Tag9,
			&channel.Tag10,
			&channel.NumViews,
			&channel.NumMessages,
			&channel.WorkspaceID,
			&channel.UserID,
			&channel.UgroupID,
			/*  StatusDates  */
			&channel.Statusc,
			&channel.CreatedAt,
			&channel.UpdatedAt,
			&channel.CreatedDay,
			&channel.CreatedWeek,
			&channel.CreatedMonth,
			&channel.CreatedYear,
			&channel.UpdatedDay,
			&channel.UpdatedWeek,
			&channel.UpdatedMonth,
			&channel.UpdatedYear)

		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5356}).Error(err)
			return nil, err
		}
		uuid4Str, err := common.UUIDBytesToStr(channel.UUID4)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5357}).Error(err)
			return nil, err
		}
		channel.IDS = uuid4Str
		return &channel, nil
	}
}

// GetChannelWithMessages - Get channel with messages
func (t *ChannelService) GetChannelWithMessages(ctx context.Context, ID string, userEmail string, requestID string) (*Channel, error) {
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
		channel := &Channel{}
		chnl, err := t.GetChannel(ctx, ID, userEmail, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5312}).Error(err)
			return nil, err
		}
		var isPresent bool
		row := db.QueryRowContext(ctx, `select exists (select 1 from messages where channel_id = ?);`, chnl.ID)
		err = row.Scan(&isPresent)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5313}).Error(err)
			return nil, err
		}
		if isPresent {
			channel, err = t.GetChannelMessages(ctx, uuid4byte, userEmail, requestID)
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5380}).Error(err)
				return nil, err
			}
		} else {
			channel = chnl
		}

		if len(channel.Messages) > 0 {
			msgserv := &MessageService{DBService: t.DBService, RedisService: t.RedisService}
			Messages, err := msgserv.GetMessagesWithTextAttach(ctx, channel.Messages, userEmail, requestID)
			if err != nil {
				log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5320}).Error(err)
			}
			channel.Messages = Messages
		}
		return channel, nil
	}
}

// GetChannelMessages - get channel with messages
func (t *ChannelService) GetChannelMessages(ctx context.Context, uuid4byte []byte, userEmail string, requestID string) (*Channel, error) {
	db := t.DBService.DB
	channel := Channel{}
	rows, err := db.QueryContext(ctx, `select 
      p.id,
			p.uuid4,
			p.channel_name,
			p.channel_desc,
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
			p.workspace_id,
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
			m.workspace_id,
			m.channel_id,
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
			m.updated_year from channels p inner join messages m on (p.id = m.channel_id) where p.uuid4 = ?`, uuid4byte)

	if err != nil {
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5314}).Error(err)
		return nil, err
	}
	for rows.Next() {
		msg := Message{}
		err = rows.Scan(
			&channel.ID,
			&channel.UUID4,
			&channel.ChannelName,
			&channel.ChannelDesc,
			&channel.NumTags,
			&channel.Tag1,
			&channel.Tag2,
			&channel.Tag3,
			&channel.Tag4,
			&channel.Tag5,
			&channel.Tag6,
			&channel.Tag7,
			&channel.Tag8,
			&channel.Tag9,
			&channel.Tag10,
			&channel.NumViews,
			&channel.NumMessages,
			&channel.WorkspaceID,
			&channel.UgroupID,
			&channel.UserID,
			/*  StatusDates  */
			&channel.Statusc,
			&channel.CreatedAt,
			&channel.UpdatedAt,
			&channel.CreatedDay,
			&channel.CreatedWeek,
			&channel.CreatedMonth,
			&channel.CreatedYear,
			&channel.UpdatedDay,
			&channel.UpdatedWeek,
			&channel.UpdatedMonth,
			&channel.UpdatedYear,
			&msg.ID,
			&msg.UUID4,
			&msg.NumLikes,
			&msg.NumUpvotes,
			&msg.NumDownvotes,
			&msg.WorkspaceID,
			&msg.ChannelID,
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
		uuid4Str1, err := common.UUIDBytesToStr(channel.UUID4)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5316}).Error(err)
			return nil, err
		}
		channel.IDS = uuid4Str1

		uuid4Str, err := common.UUIDBytesToStr(msg.UUID4)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5317}).Error(err)
			return nil, err
		}
		msg.IDS = uuid4Str
		channel.Messages = append(channel.Messages, &msg)
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
	return &channel, nil
}

// GetChannelsUser - Get user channels
func (t *ChannelService) GetChannelsUser(ctx context.Context, ID uint, UserID uint, userEmail string, requestID string) (*ChannelsUser, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5358}).Error(err)
		return nil, err
	default:
		channel := ChannelsUser{}
		db := t.DBService.DB
		row := db.QueryRowContext(ctx, `select
    id,
		uuid4,
		channel_id,
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
		updated_year from channels_users where channel_id = ? and user_id = ? and statusc = ?`, ID, UserID, common.Active)

		err := row.Scan(
			&channel.ID,
			&channel.UUID4,
			&channel.ChannelID,
			&channel.NumMessages,
			&channel.NumViews,
			&channel.UserID,
			&channel.UgroupID,
			/*  StatusDates  */
			&channel.Statusc,
			&channel.CreatedAt,
			&channel.UpdatedAt,
			&channel.CreatedDay,
			&channel.CreatedWeek,
			&channel.CreatedMonth,
			&channel.CreatedYear,
			&channel.UpdatedDay,
			&channel.UpdatedWeek,
			&channel.UpdatedMonth,
			&channel.UpdatedYear)

		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5359}).Error(err)
			return nil, err
		}
		uuid4Str, err := common.UUIDBytesToStr(channel.UUID4)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5360}).Error(err)
			return nil, err
		}
		channel.IDS = uuid4Str
		return &channel, nil
	}
}

//UpdateChannel - Update channel
func (t *ChannelService) UpdateChannel(ctx context.Context, ID string, form *Channel, UserID string, userEmail string, requestID string) error {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5382}).Error(err)
		return err
	default:
		channel, err := t.GetChannel(ctx, ID, userEmail, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": userEmail, "reqid": requestID, "msgnum": 5383}).Error(err)
			return err
		}

		db := t.DBService.DB

		tn, tnday, tnweek, tnmonth, tnyear := common.GetTimeDetails()
		stmt, err := db.PrepareContext(ctx, `update channels set 
		  channel_name = ?,
      channel_desc = ?,
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
			form.ChannelName,
			form.ChannelDesc,
			tn,
			tnday,
			tnweek,
			tnmonth,
			tnyear,
			channel.ID,
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

// DeleteChannel - Delete channel
func (t *ChannelService) DeleteChannel(ctx context.Context, ID string, userEmail string, requestID string) error {
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
		stmt, err := db.PrepareContext(ctx, `update channels set 
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
