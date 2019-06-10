package msgcontrollers

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/cloudfresco/vilom/common"
	"github.com/cloudfresco/vilom/msg/msgservices"
)

/* error message range: 5000-5299 */

// TopicController - Create Topic Controller
type TopicController struct {
	Service *msgservices.TopicService
}

// NewTopicController - Create Topic Handler
func NewTopicController(s *msgservices.TopicService) *TopicController {
	return &TopicController{s}
}

func (tc *TopicController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user, requestID, err := common.GetAuthUserDetails(r, tc.Service.RedisClient, tc.Service.Db)
	if err != nil {
		common.RenderErrorJSON(w, "1001", err.Error(), 401, requestID)
		return
	}
	var pathParts []string

	path := r.URL.Path
	pathParts = common.GetPathParts(path)

	switch r.Method {
	case http.MethodGet:

		/*
		   GET  "/v1/topics/{id}"
		*/
		if (len(pathParts) == 3) && (pathParts[1] == "topics") {
			tc.Show(w, r, pathParts[2], user, requestID)
		}

	case http.MethodPost:
		/*
					   POST  "/v1/topics/create/"
			       POST  "/v1/topics/topicbyname/"
		*/
		if (len(pathParts) == 3) && (pathParts[1] == "topics") && (pathParts[2] == "create") {
			tc.Create(w, r, user, requestID)
		} else if (len(pathParts) == 3) && (pathParts[1] == "topics") && (pathParts[2] == "topicbyname") {
			tc.Topicbyname(w, r, user, requestID)
		}
	case http.MethodPut:
	case http.MethodDelete:
	default:
		common.RenderErrorJSON(w, "1000", "Invalid Request", 400, requestID)
		return
	}

}

// Show - used to view Topic
func (tc *TopicController) Show(w http.ResponseWriter, r *http.Request, id string, user *common.ContextData, requestID string) {
	ctx := r.Context()

	select {
	case <-ctx.Done():
		common.RenderErrorJSON(w, "1002", "Client closed connection", 402, requestID)
		return
	default:

		topic, err := tc.Service.Show(ctx, id, user.UserID, user.Email, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": user.Email, "reqid": requestID, "msgnum": 5000}).Error(err)
			common.RenderErrorJSON(w, "5000", err.Error(), 402, requestID)
			return
		}

		common.RenderJSON(w, topic)
	}
}

// Create - used to Create Topic
func (tc *TopicController) Create(w http.ResponseWriter, r *http.Request, user *common.ContextData, requestID string) {
	ctx := r.Context()

	select {
	case <-ctx.Done():
		common.RenderErrorJSON(w, "1002", "Client closed connection", 402, requestID)
		return
	default:
		form := msgservices.Topic{}
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&form)
		if err != nil {
			log.WithFields(log.Fields{"user": user.Email, "reqid": requestID, "msgnum": 5001}).Error(err)
			common.RenderErrorJSON(w, "5001", err.Error(), 402, requestID)
			return
		}
		topic, err := tc.Service.Create(ctx, &form, user.UserID, user.Email, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": user.Email, "reqid": requestID, "msgnum": 5002}).Error(err)
			common.RenderErrorJSON(w, "5002", err.Error(), 402, requestID)
			return
		}

		common.RenderJSON(w, topic)
	}
}

// Topicbyname - used to get Topic by name
func (tc *TopicController) Topicbyname(w http.ResponseWriter, r *http.Request, user *common.ContextData, requestID string) {
	ctx := r.Context()

	select {
	case <-ctx.Done():
		common.RenderErrorJSON(w, "1002", "Client closed connection", 402, requestID)
		return
	default:
		form := msgservices.Topic{}
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&form)
		if err != nil {
			log.WithFields(log.Fields{"user": user.Email, "reqid": requestID, "msgnum": 5003}).Error(err)
			common.RenderErrorJSON(w, "5003", err.Error(), 402, requestID)
			return
		}
		topc, err := tc.Service.GetTopicByName(ctx, form.TopicName, user.Email, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": user.Email, "reqid": requestID, "msgnum": 5004}).Error(err)
			common.RenderErrorJSON(w, "5004", err.Error(), 402, requestID)
			return
		}

		common.RenderJSON(w, topc)
	}
}
