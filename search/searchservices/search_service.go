package searchservices

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"database/sql"
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/analysis/analyzer/custom"
	"github.com/blevesearch/bleve/analysis/analyzer/keyword"
	"github.com/blevesearch/bleve/analysis/lang/en"
	"github.com/blevesearch/bleve/analysis/token/edgengram"
	"github.com/blevesearch/bleve/analysis/token/lowercase"
	"github.com/blevesearch/bleve/analysis/tokenizer/unicode"
	"github.com/blevesearch/bleve/mapping"

	log "github.com/sirupsen/logrus"

	"github.com/cloudfresco/vilom/common"
	"github.com/cloudfresco/vilom/msg/msgservices"
)

/* error message range: 7300-7999 */

// BleveForm - Search form
type BleveForm struct {
	SearchText string
}

// SearchServiceIntf - interface for Search Service
type SearchServiceIntf interface {
	Search(form *BleveForm, userEmail string, requestID string) (*bleve.SearchResult, error)
}

// SearchService -  For accessing  search service
type SearchService struct {
	DBService    *common.DBService
	RedisService *common.RedisService
	SearchIndex  bleve.Index
}

// NewSearchService - Create search service
func NewSearchService(dbOpt *common.DBService, redisOpt *common.RedisService, searchIndex bleve.Index) *SearchService {
	return &SearchService{
		DBService:    dbOpt,
		RedisService: redisOpt,
		SearchIndex:  searchIndex,
	}
}

var bSearchIndex bleve.Index

// InitSearch -
func InitSearch(p string, db *sql.DB) bleve.Index {
	indexPath := ""
	pwd, _ := os.Getwd()
	indexPath = pwd + filepath.FromSlash("/files/search/channels.bleve")
	productIndex, err := bleve.OpenUsing(indexPath, map[string]interface{}{
		"read_only": true,
	})

	if err == bleve.ErrorIndexPathDoesNotExist {
		productMapping, err := buildIndexMapping()
		if err != nil {
			log.WithFields(log.Fields{
				"msgnum": 7002,
			}).Error(err)
		}
		productIndex, err = bleve.New(indexPath, productMapping)
		if err != nil {
			log.WithFields(log.Fields{
				"msgnum": 7003,
			}).Error(err)
		}

	} else if err != nil {
		log.WithFields(log.Fields{
			"msgnum": 7004,
		}).Error(err)
	}

	err = IndexChannels(db, productIndex)
	if err != nil {
		log.WithFields(log.Fields{
			"msgnum": 7005,
		}).Error(err)
	}

	bSearchIndex = productIndex
	return bSearchIndex
}

// buildIndexMapping - used for
func buildIndexMapping() (mapping.IndexMapping, error) {

	edgeNgram325FieldMapping := bleve.NewTextFieldMapping()
	edgeNgram325FieldMapping.Analyzer = "enWithEdgeNgram325"

	// a generic reusable mapping for english text
	englishTextFieldMapping := bleve.NewTextFieldMapping()
	englishTextFieldMapping.Analyzer = en.AnalyzerName

	// a generic reusable mapping for keyword text
	keywordFieldMapping := bleve.NewTextFieldMapping()
	keywordFieldMapping.Analyzer = keyword.Name

	channelMapping := bleve.NewDocumentMapping()

	//disabledMapping := bleve.NewDocumentDisabledMapping()

	// name
	channelMapping.AddFieldMappingsAt("Name", englishTextFieldMapping)

	// description
	channelMapping.AddFieldMappingsAt("Description",
		englishTextFieldMapping, edgeNgram325FieldMapping)

	// messagetext
	channelMapping.AddFieldMappingsAt("MessageText", englishTextFieldMapping, edgeNgram325FieldMapping)

	indexMapping := bleve.NewIndexMapping()
	indexMapping.AddDocumentMapping("Name", channelMapping)
	err := indexMapping.AddCustomTokenFilter("edgeNgram325",
		map[string]interface{}{
			"type": edgengram.Name,
			"min":  3.0,
			"max":  25.0,
		})
	if err != nil {
		log.WithFields(log.Fields{
			"msgnum": 7006,
		}).Error(err)
		return nil, err
	}

	err = indexMapping.AddCustomAnalyzer("enWithEdgeNgram325",
		map[string]interface{}{
			"type":      custom.Name,
			"tokenizer": unicode.Name,
			"token_filters": []string{
				en.PossessiveName,
				lowercase.Name,
				en.StopName,
				"edgeNgram325",
			},
		})
	if err != nil {
		log.WithFields(log.Fields{
			"msgnum": 7007,
		}).Error(err)
		return nil, err
	}
	indexMapping.TypeField = "Type"
	indexMapping.DefaultAnalyzer = "en"

	return indexMapping, nil
}

// IndexChannels - used for
func IndexChannels(db *sql.DB, index bleve.Index) error {
	batch := index.NewBatch()
	var channelMsgMap map[string]string
	docID := ""
	channels, err := getChannels(db)

	if err != nil {
		log.WithFields(log.Fields{
			"msgnum": 7008,
		}).Error(err)
		return err
	}
	for _, channel := range channels {
		messages, err := getMessagesByChannelID(channel.ID, db)
		if err != nil {
			log.WithFields(log.Fields{
				"msgnum": 7009,
			}).Error(err)
			return err
		}
		for _, message := range messages {
			channelMsgMap = map[string]string{"Type": "Name"}

			channelMsgMap["Name"] = channel.ChannelName
			channelMsgMap["Description"] = channel.ChannelDesc
			channelMsgMap["Pid"] = strconv.FormatUint(uint64(channel.ID), 10)
			channelMsgMap["MessageText"] = message.Mtext

			docID = fmt.Sprintf("%d##%d", channel.ID, message.ID)

			p, err := json.Marshal(channelMsgMap)

			if err == nil {
				var prodInterface interface{}

				err := json.Unmarshal(p, &prodInterface)
				if err != nil {
					log.WithFields(log.Fields{
						"msgnum": 7010,
					}).Error(err)
					return err
				}

				err = batch.Index(docID, channelMsgMap)
				if err != nil {
					log.WithFields(log.Fields{
						"msgnum": 7011,
					}).Error(err)
					return err
				}
				if batch.Size() >= 100 {
					err := index.Batch(batch)
					if err != nil {
						log.WithFields(log.Fields{
							"msgnum": 7012,
						}).Error(err)
						return err
					}
					batch = index.NewBatch()
				}
			} else {
				log.WithFields(log.Fields{
					"msgnum": 7013,
				}).Error(err)
			}
		}
	}

	if batch.Size() > 0 {
		err := index.Batch(batch)
		if err != nil {
			log.WithFields(log.Fields{
				"msgnum": 7014,
			}).Error(err)
			return err
		}
	}
	return nil
}

// Search - used for
func (t *SearchService) Search(form *BleveForm, userEmail string, requestID string) (*bleve.SearchResult, error) {
	query := bleve.NewMatchQuery(form.SearchText)
	search := bleve.NewSearchRequest(query)
	fields := []string{"Name", "Description", "Pid", "MessageText"}

	search.Fields = fields
	searchResults, err := t.SearchIndex.Search(search)
	if err != nil {
		log.WithFields(log.Fields{
			"user":   userEmail,
			"reqid":  requestID,
			"msgnum": 7015,
		}).Error(err)
		return nil, err
	}
	return searchResults, nil
}

// getMessagesByChannelID - Get messages by channel id
func getMessagesByChannelID(ID uint, db *sql.DB) ([]*msgservices.MessageText, error) {
	msgs := []*msgservices.MessageText{}
	rows, err := db.Query(`select 
    id,
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
		updated_year from message_texts where channel_id = ?`, ID)

	if err != nil {
		log.WithFields(log.Fields{
			"msgnum": 7016,
		}).Error(err)
	}
	for rows.Next() {
		msgtxt := msgservices.MessageText{}
		err = rows.Scan(
			&msgtxt.ID,
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
			log.WithFields(log.Fields{
				"msgnum": 7017,
			}).Error(err)
			return nil, err
		}
		msgs = append(msgs, &msgtxt)
	}

	err = rows.Close()
	if err != nil {
		log.WithFields(log.Fields{
			"msgnum": 7018,
		}).Error(err)
		return nil, err
	}

	return msgs, nil
}

// getChannels - Get channels
func getChannels(db *sql.DB) ([]*msgservices.Channel, error) {

	pohs := []*msgservices.Channel{}
	rows, err := db.Query(`select 
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
		updated_year from channels`)

	if err != nil {
		log.WithFields(log.Fields{
			"msgnum": 7019,
		}).Error(err)
	}
	for rows.Next() {
		poh := msgservices.Channel{}
		err = rows.Scan(
			&poh.ID,
			&poh.UUID4,
			&poh.ChannelName,
			&poh.ChannelDesc,
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
			&poh.WorkspaceID,
			&poh.UserID,
			&poh.UgroupID,
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
			log.WithFields(log.Fields{
				"msgnum": 7020,
			}).Error(err)
			return nil, err
		}
		uuid4Str, err := common.UUIDBytesToStr(poh.UUID4)
		if err != nil {
			log.WithFields(log.Fields{
				"msgnum": 7021,
			}).Error(err)
			log.Println(err)
		}
		poh.IDS = uuid4Str
		pohs = append(pohs, &poh)

	}

	err = rows.Close()
	if err != nil {
		log.WithFields(log.Fields{
			"msgnum": 7022,
		}).Error(err)
		return nil, err
	}

	return pohs, nil
}
