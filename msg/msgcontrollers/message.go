package msgcontrollers

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/cloudfresco/vilom/common"
	"github.com/cloudfresco/vilom/msg/msgservices"
)

// MessageController - used for Messages
type MessageController struct {
	Service *msgservices.MessageService
}

// NewMessageController - used for Messages
func NewMessageController(s *msgservices.MessageService) *MessageController {
	return &MessageController{s}
}

func (mc *MessageController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user, requestID, err := common.GetAuthUserDetails(r, mc.Service.RedisClient, mc.Service.Db)
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
		   GET  "/v1/messages/{id}"
		*/

		if (len(pathParts) == 3) && (pathParts[1] == "messages") {
			mc.Show(w, r, pathParts[2], user, requestID)
		}

	case http.MethodPost:
		/*
		   POST  "/v1/messages/create/"
		   POST  "/v1/messages/like/"
		*/
		if (len(pathParts) == 3) && (pathParts[1] == "messages") && (pathParts[2] == "create") {
			mc.Create(w, r, user, requestID)
		} else if (len(pathParts) == 3) && (pathParts[1] == "messages") && (pathParts[2] == "like") {
			mc.UserLikeCreate(w, r, user, requestID)
		}
	case http.MethodPut:
	case http.MethodDelete:
	default:
		common.RenderErrorJSON(w, "1000", "Invalid Request", 400, requestID)
		return
	}
}

// Show - used to view message
func (mc *MessageController) Show(w http.ResponseWriter, r *http.Request, id string, user *common.ContextData, requestID string) {
	ctx := r.Context()

	select {
	case <-ctx.Done():
		common.RenderErrorJSON(w, "1002", "Client closed connection", 402, requestID)
		return
	default:
		msg, err := mc.Service.GetMessage(ctx, id, user.Email, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": user.Email, "reqid": requestID, "msgnum": 6000}).Error(err)
			common.RenderErrorJSON(w, "6000", err.Error(), 402, requestID)
			return
		}

		common.RenderJSON(w, msg)
	}
}

// Create - Create Message
func (mc *MessageController) Create(w http.ResponseWriter, r *http.Request, user *common.ContextData, requestID string) {
	ctx := r.Context()

	select {
	case <-ctx.Done():
		common.RenderErrorJSON(w, "1002", "Client closed connection", 402, requestID)
		return
	default:
		form := msgservices.Message{}
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&form)
		if err != nil {
			log.WithFields(log.Fields{"user": user.Email, "reqid": requestID, "msgnum": 6001}).Error(err)
			common.RenderErrorJSON(w, "6001", err.Error(), 402, requestID)
			return
		}
		msg, err := mc.Service.Create(ctx, &form, user.UserID, user.Email, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": user.Email, "reqid": requestID, "msgnum": 6002}).Error(err)
			common.RenderErrorJSON(w, "6002", err.Error(), 402, requestID)
			return
		}

		common.RenderJSON(w, msg)
	}
}

// UserLikeCreate - Create User Like
func (mc *MessageController) UserLikeCreate(w http.ResponseWriter, r *http.Request, user *common.ContextData, requestID string) {
	ctx := r.Context()

	select {
	case <-ctx.Done():
		common.RenderErrorJSON(w, "1002", "Client closed connection", 402, requestID)
		return
	default:
		form := msgservices.UserLike{}
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&form)
		if err != nil {
			log.WithFields(log.Fields{"user": user.Email, "reqid": requestID, "msgnum": 6003}).Error(err)
			common.RenderErrorJSON(w, "6003", err.Error(), 402, requestID)
			return
		}
		msg, err := mc.Service.UserLikeCreate(ctx, &form, user.UserID, user.Email, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": user.Email, "reqid": requestID, "msgnum": 6004}).Error(err)
			common.RenderErrorJSON(w, "6004", err.Error(), 402, requestID)
			return
		}

		common.RenderJSON(w, msg)
	}
}
