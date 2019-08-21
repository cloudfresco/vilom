package msgcontrollers

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/cloudfresco/vilom/common"
	"github.com/cloudfresco/vilom/msg/msgservices"
	"github.com/cloudfresco/vilom/user/userservices"
)

/* error message range: 5000-5299 */

// TopicController - Create Topic Controller
type TopicController struct {
	Service  msgservices.TopicServiceIntf
	Serviceu userservices.UserServiceIntf
}

// NewTopicController - Create Topic Handler
func NewTopicController(s msgservices.TopicServiceIntf, su userservices.UserServiceIntf) *TopicController {
	return &TopicController{
		Service:  s,
		Serviceu: su,
	}
}

// ServeHTTP - parse url and call controller action
func (tc *TopicController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user, requestID, err := tc.Serviceu.GetAuthUserDetails(r)
	if err != nil {
		common.RenderErrorJSON(w, "1001", err.Error(), 401, requestID)
		return
	}
	var pathParts []string

	path := r.URL.Path
	pathParts = common.GetPathParts(path)

	switch r.Method {
	case http.MethodGet:
		tc.processGet(w, r, user, requestID, pathParts)
	case http.MethodPost:
		tc.processPost(w, r, user, requestID, pathParts)
	case http.MethodPut:
		tc.processPut(w, r, user, requestID, pathParts)
	case http.MethodDelete:
		tc.processDelete(w, r, user, requestID, pathParts)
	default:
		common.RenderErrorJSON(w, "1000", "Invalid Request", 400, requestID)
		return
	}

}

// processGet - Parse URL for all the GET paths and call the controller action
/*
 GET  "/v1/topics/{id}"
*/
func (tc *TopicController) processGet(w http.ResponseWriter, r *http.Request, user *common.ContextData, requestID string, pathParts []string) {

	if (len(pathParts) == 3) && (pathParts[1] == "topics") {
		tc.ShowTopic(w, r, pathParts[2], user, requestID)
	} else {
		common.RenderErrorJSON(w, "1000", "Invalid Request", 400, requestID)
		return
	}

}

// processPost - Parse URL for all the POST paths and call the controller action
/*
 POST  "/v1/topics/create/"
 POST  "/v1/topics/topicbyname/"
*/
func (tc *TopicController) processPost(w http.ResponseWriter, r *http.Request, user *common.ContextData, requestID string, pathParts []string) {

	if (len(pathParts) == 3) && (pathParts[1] == "topics") {
		if pathParts[2] == "create" {
			tc.CreateTopic(w, r, user, requestID)
		} else if pathParts[2] == "topicbyname" {
			tc.GetTopicByName(w, r, user, requestID)
		} else {
			common.RenderErrorJSON(w, "1000", "Invalid Request", 400, requestID)
			return
		}
	} else {
		common.RenderErrorJSON(w, "1000", "Invalid Request", 400, requestID)
		return
	}
}

// processPut - Parse URL for all the put paths and call the controller action
/*
 PUT  "/v1/topics/{id}"
*/

func (tc *TopicController) processPut(w http.ResponseWriter, r *http.Request, user *common.ContextData, requestID string, pathParts []string) {

	if (len(pathParts) == 3) && (pathParts[1] == "topics") {
		tc.UpdateTopic(w, r, pathParts[2], user, requestID)
	} else {
		common.RenderErrorJSON(w, "1000", "Invalid Request", 400, requestID)
		return
	}

}

// processDelete - Parse URL for all the delete paths and call the controller action
/*
 DELETE  "/v1/topics/{id}"
*/

func (tc *TopicController) processDelete(w http.ResponseWriter, r *http.Request, user *common.ContextData, requestID string, pathParts []string) {

	if (len(pathParts) == 3) && (pathParts[1] == "topics") {
		tc.DeleteTopic(w, r, pathParts[2], user, requestID)
	} else {
		common.RenderErrorJSON(w, "1000", "Invalid Request", 400, requestID)
		return
	}

}

// ShowTopic - used to view Topic
func (tc *TopicController) ShowTopic(w http.ResponseWriter, r *http.Request, id string, user *common.ContextData, requestID string) {
	ctx := r.Context()

	select {
	case <-ctx.Done():
		common.RenderErrorJSON(w, "1002", "Client closed connection", 402, requestID)
		return
	default:

		topic, err := tc.Service.ShowTopic(ctx, id, user.UserID, user.Email, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": user.Email, "reqid": requestID, "msgnum": 5000}).Error(err)
			common.RenderErrorJSON(w, "5000", err.Error(), 402, requestID)
			return
		}

		common.RenderJSON(w, topic)
	}
}

// CreateTopic - used to Create Topic
func (tc *TopicController) CreateTopic(w http.ResponseWriter, r *http.Request, user *common.ContextData, requestID string) {
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
		v := common.NewValidator()
		v.IsStrLenBetMinMax("Topic Name", form.TopicName, msgservices.TopicNameLenMin, msgservices.TopicNameLenMax)
		v.IsStrLenBetMinMax("Topic Description", form.TopicDesc, msgservices.TopicDescLenMin, msgservices.TopicDescLenMax)
		if v.IsValid() {
			common.RenderErrorJSON(w, "5008", v.Error(), 402, requestID)
			return
		}
		topic, err := tc.Service.CreateTopic(ctx, &form, user.UserID, user.Email, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": user.Email, "reqid": requestID, "msgnum": 5002}).Error(err)
			common.RenderErrorJSON(w, "5002", err.Error(), 402, requestID)
			return
		}

		common.RenderJSON(w, topic)
	}
}

// GetTopicByName - used to get Topic by name
func (tc *TopicController) GetTopicByName(w http.ResponseWriter, r *http.Request, user *common.ContextData, requestID string) {
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

// UpdateTopic - Update topic
func (tc *TopicController) UpdateTopic(w http.ResponseWriter, r *http.Request, id string, user *common.ContextData, requestID string) {
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
			log.WithFields(log.Fields{"user": user.Email, "reqid": requestID, "msgnum": 5005}).Error(err)
			common.RenderErrorJSON(w, "5005", err.Error(), 402, requestID)
			return
		}
		err = tc.Service.UpdateTopic(ctx, id, &form, user.UserID, user.Email, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": user.Email, "reqid": requestID, "msgnum": 5006}).Error(err)
			common.RenderErrorJSON(w, "5006", err.Error(), 402, requestID)
			return
		}

		common.RenderJSON(w, "Updated Successfully")
	}
}

// DeleteTopic - delete topic
func (tc *TopicController) DeleteTopic(w http.ResponseWriter, r *http.Request, id string, user *common.ContextData, requestID string) {
	ctx := r.Context()

	select {
	case <-ctx.Done():
		common.RenderErrorJSON(w, "1002", "Client closed connection", 402, requestID)
		return
	default:
		err := tc.Service.DeleteTopic(ctx, id, user.Email, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": user.Email, "reqid": requestID, "msgnum": 5007}).Error(err)
			common.RenderErrorJSON(w, "5007", err.Error(), 402, requestID)
			return
		}

		common.RenderJSON(w, "Deleted Successfully")
	}
}
