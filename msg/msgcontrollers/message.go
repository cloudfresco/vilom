package msgcontrollers

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/cloudfresco/vilom/common"
	"github.com/cloudfresco/vilom/msg/msgservices"
	"github.com/cloudfresco/vilom/user/userservices"
)

/* error message range: 6000-6299 */

// MessageController - used for Messages
type MessageController struct {
	Service  msgservices.MessageServiceIntf
	Serviceu userservices.UserServiceIntf
}

// NewMessageController - used for Messages
func NewMessageController(s msgservices.MessageServiceIntf, su userservices.UserServiceIntf) *MessageController {
	return &MessageController{
		Service:  s,
		Serviceu: su,
	}
}

// ServeHTTP - parse url and call controller action
func (mc *MessageController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user, requestID, err := mc.Serviceu.GetAuthUserDetails(r)
	if err != nil {
		common.RenderErrorJSON(w, "1001", err.Error(), 401, requestID)
		return
	}
	var pathParts []string

	path := r.URL.Path
	pathParts = common.GetPathParts(path)

	switch r.Method {
	case http.MethodGet:
		mc.processGet(w, r, user, requestID, pathParts)
	case http.MethodPost:
		mc.processPost(w, r, user, requestID, pathParts)
	case http.MethodPut:
		mc.processPut(w, r, user, requestID, pathParts)
	case http.MethodDelete:
		mc.processDelete(w, r, user, requestID, pathParts)
	default:
		common.RenderErrorJSON(w, "1000", "Invalid Request", 400, requestID)
		return
	}
}

// processGet - Parse URL for all the GET paths and call the controller action
/*
 GET  "/v1/messages/{id}"
*/

func (mc *MessageController) processGet(w http.ResponseWriter, r *http.Request, user *common.ContextData, requestID string, pathParts []string) {

	if (len(pathParts) == 3) && (pathParts[1] == "messages") {
		mc.Show(w, r, pathParts[2], user, requestID)
	} else {
		common.RenderErrorJSON(w, "1000", "Invalid Request", 400, requestID)
		return
	}

}

// processPost - Parse URL for all the POST paths and call the controller action
/*
 POST  "/v1/messages/create/"
 POST  "/v1/messages/like/"
*/
func (mc *MessageController) processPost(w http.ResponseWriter, r *http.Request, user *common.ContextData, requestID string, pathParts []string) {

	if (len(pathParts) == 3) && (pathParts[1] == "messages") {
		if pathParts[2] == "create" {
			mc.Create(w, r, user, requestID)
		} else if pathParts[2] == "like" {
			mc.UserLikeCreate(w, r, user, requestID)
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
 PUT  "/v1/messages/{id}"
*/

func (mc *MessageController) processPut(w http.ResponseWriter, r *http.Request, user *common.ContextData, requestID string, pathParts []string) {

	if (len(pathParts) == 3) && (pathParts[1] == "messages") {
		mc.Update(w, r, pathParts[2], user, requestID)
	} else {
		common.RenderErrorJSON(w, "1000", "Invalid Request", 400, requestID)
		return
	}

}

// processDelete - Parse URL for all the delete paths and call the controller action
/*
 DELETE  "/v1/messages/{id}"
*/

func (mc *MessageController) processDelete(w http.ResponseWriter, r *http.Request, user *common.ContextData, requestID string, pathParts []string) {

	if (len(pathParts) == 3) && (pathParts[1] == "messages") {
		mc.Delete(w, r, pathParts[2], user, requestID)
	} else {
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
		msg, err := mc.Service.CreateMessage(ctx, &form, user.UserID, true, user.Email, requestID)
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
		msg, err := mc.Service.CreateUserLike(ctx, &form, user.UserID, user.Email, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": user.Email, "reqid": requestID, "msgnum": 6004}).Error(err)
			common.RenderErrorJSON(w, "6004", err.Error(), 402, requestID)
			return
		}

		common.RenderJSON(w, msg)
	}
}

// Update - Update Message
func (mc *MessageController) Update(w http.ResponseWriter, r *http.Request, id string, user *common.ContextData, requestID string) {
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
			log.WithFields(log.Fields{"user": user.Email, "reqid": requestID, "msgnum": 6005}).Error(err)
			common.RenderErrorJSON(w, "6005", err.Error(), 402, requestID)
			return
		}
		err = mc.Service.UpdateMessage(ctx, id, &form, user.UserID, user.Email, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": user.Email, "reqid": requestID, "msgnum": 6006}).Error(err)
			common.RenderErrorJSON(w, "6006", err.Error(), 402, requestID)
			return
		}

		common.RenderJSON(w, "Updated Successfully")
	}
}

// Delete - delete message
func (mc *MessageController) Delete(w http.ResponseWriter, r *http.Request, id string, user *common.ContextData, requestID string) {
	ctx := r.Context()

	select {
	case <-ctx.Done():
		common.RenderErrorJSON(w, "1002", "Client closed connection", 402, requestID)
		return
	default:
		err := mc.Service.DeleteMessage(ctx, id, user.Email, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": user.Email, "reqid": requestID, "msgnum": 6007}).Error(err)
			common.RenderErrorJSON(w, "6007", err.Error(), 402, requestID)
			return
		}

		common.RenderJSON(w, "Deleted Successfully")
	}
}
