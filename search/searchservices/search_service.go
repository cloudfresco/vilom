package searchservices

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"database/sql"
	log "github.com/Sirupsen/logrus"
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/analysis/analyzer/custom"
	"github.com/blevesearch/bleve/analysis/analyzer/keyword"
	"github.com/blevesearch/bleve/analysis/lang/en"
	"github.com/blevesearch/bleve/analysis/token/edgengram"
	"github.com/blevesearch/bleve/analysis/token/lowercase"
	"github.com/blevesearch/bleve/analysis/tokenizer/unicode"
	"github.com/blevesearch/bleve/mapping"
	"github.com/go-redis/redis"

	"github.com/cloudfresco/vilom/common"
	"github.com/cloudfresco/vilom/msg/msgservices"
)

// BleveForm - Search form
type BleveForm struct {
	SearchText string
}

// SearchService -  For accessing  search service
type SearchService struct {
	Config      *common.RedisOptions
	Db          *sql.DB
	RedisClient *redis.Client
	SearchIndex bleve.Index
}

// NewSearchService - Create search service
func NewSearchService(config *common.RedisOptions,
	db *sql.DB,
	redisClient *redis.Client,
	searchIndex bleve.Index) *SearchService {
	return &SearchService{config, db, redisClient, searchIndex}
}

var bSearchIndex bleve.Index

// InitSearch -
func InitSearch(p string, db *sql.DB) bleve.Index {
	indexPath := ""
	pwd, _ := os.Getwd()
	indexPath = pwd + filepath.FromSlash("/files/search/topics.bleve")
	productIndex, err := bleve.OpenUsing(indexPath, map[string]interface{}{
		"read_only": true,
	})

	if err == bleve.ErrorIndexPathDoesNotExist {
		productMapping, err := BuildIndexMapping()
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

	err = IndexTopics(db, productIndex)
	if err != nil {
		log.WithFields(log.Fields{
			"msgnum": 7005,
		}).Error(err)
	}

	bSearchIndex = productIndex
	return bSearchIndex
}

// BuildIndexMapping - used for
func BuildIndexMapping() (mapping.IndexMapping, error) {

	edgeNgram325FieldMapping := bleve.NewTextFieldMapping()
	edgeNgram325FieldMapping.Analyzer = "enWithEdgeNgram325"

	// a generic reusable mapping for english text
	englishTextFieldMapping := bleve.NewTextFieldMapping()
	englishTextFieldMapping.Analyzer = en.AnalyzerName

	// a generic reusable mapping for keyword text
	keywordFieldMapping := bleve.NewTextFieldMapping()
	keywordFieldMapping.Analyzer = keyword.Name

	topicMapping := bleve.NewDocumentMapping()

	//disabledMapping := bleve.NewDocumentDisabledMapping()

	// name
	topicMapping.AddFieldMappingsAt("Name", englishTextFieldMapping)

	// description
	topicMapping.AddFieldMappingsAt("Description",
		englishTextFieldMapping, edgeNgram325FieldMapping)

	// messagetext
	topicMapping.AddFieldMappingsAt("MessageText", englishTextFieldMapping, edgeNgram325FieldMapping)

	indexMapping := bleve.NewIndexMapping()
	indexMapping.AddDocumentMapping("Name", topicMapping)
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

// IndexTopics - used for
func IndexTopics(db *sql.DB, index bleve.Index) error {
	count := 0
	batch := index.NewBatch()
	prod := make(map[string]string)
	docID := ""
	topics, err := GetTopics(db)

	if err != nil {
		log.WithFields(log.Fields{
			"msgnum": 7008,
		}).Error(err)
		return err
	}
	for _, topic := range topics {
		messages, err := GetMessagesByTopicID(topic.ID, db)
		if err != nil {
			log.WithFields(log.Fields{
				"msgnum": 7009,
			}).Error(err)
			return err
		}
		for _, message := range messages {
			prod = map[string]string{"Type": "Name"}

			prod["Name"] = topic.TopicName
			prod["Description"] = topic.TopicDesc
			prod["Pid"] = strconv.FormatUint(uint64(topic.ID), 10)
			prod["MessageText"] = message.Mtext

			docID = fmt.Sprintf("%d##%d", topic.ID, message.ID)

			p, err := json.Marshal(prod)

			if err == nil {
				var prodInterface interface{}

				err := json.Unmarshal(p, &prodInterface)
				if err != nil {
					log.WithFields(log.Fields{
						"msgnum": 7010,
					}).Error(err)
					return err
				}

				err = batch.Index(docID, prod)
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
					count += batch.Size()
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
		count += batch.Size()
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

// GetMessagesByTopicID - Get messages by topic id
func GetMessagesByTopicID(ID uint, db *sql.DB) ([]*msgservices.MessageText, error) {
	msgs := []*msgservices.MessageText{}
	rows, err := db.Query(`select 
    id,
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
		updated_year from message_texts where topic_id = ?`, ID)

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

// GetTopics - Get topics
func GetTopics(db *sql.DB) ([]*msgservices.Topic, error) {

	pohs := []*msgservices.Topic{}
	rows, err := db.Query(`select 
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
		updated_year from topics`)

	if err != nil {
		log.WithFields(log.Fields{
			"msgnum": 7019,
		}).Error(err)
	}
	for rows.Next() {
		poh := msgservices.Topic{}
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
		uUID4Str, err := common.UUIDBytesToStr(poh.UUID4)
		if err != nil {
			log.WithFields(log.Fields{
				"msgnum": 7021,
			}).Error(err)
			log.Println(err)
		}
		poh.IDS = uUID4Str
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
