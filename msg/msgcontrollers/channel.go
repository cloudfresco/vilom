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

// ChannelController - Create Channel Controller
type ChannelController struct {
	Service  msgservices.ChannelServiceIntf
	Serviceu userservices.UserServiceIntf
}

// NewChannelController - Create Channel Handler
func NewChannelController(s msgservices.ChannelServiceIntf, su userservices.UserServiceIntf) *ChannelController {
	return &ChannelController{
		Service:  s,
		Serviceu: su,
	}
}

// ServeHTTP - parse url and call controller action
func (tc *ChannelController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
 GET  "/v1/channels/{id}"
*/
func (tc *ChannelController) processGet(w http.ResponseWriter, r *http.Request, user *common.ContextData, requestID string, pathParts []string) {

	if (len(pathParts) == 3) && (pathParts[1] == "channels") {
		tc.ShowChannel(w, r, pathParts[2], user, requestID)
	} else {
		common.RenderErrorJSON(w, "1000", "Invalid Request", 400, requestID)
		return
	}

}

// processPost - Parse URL for all the POST paths and call the controller action
/*
 POST  "/v1/channels/create/"
 POST  "/v1/channels/channelbyname/"
*/
func (tc *ChannelController) processPost(w http.ResponseWriter, r *http.Request, user *common.ContextData, requestID string, pathParts []string) {

	if (len(pathParts) == 3) && (pathParts[1] == "channels") {
		if pathParts[2] == "create" {
			tc.CreateChannel(w, r, user, requestID)
		} else if pathParts[2] == "channelbyname" {
			tc.GetChannelByName(w, r, user, requestID)
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
 PUT  "/v1/channels/{id}"
*/

func (tc *ChannelController) processPut(w http.ResponseWriter, r *http.Request, user *common.ContextData, requestID string, pathParts []string) {

	if (len(pathParts) == 3) && (pathParts[1] == "channels") {
		tc.UpdateChannel(w, r, pathParts[2], user, requestID)
	} else {
		common.RenderErrorJSON(w, "1000", "Invalid Request", 400, requestID)
		return
	}

}

// processDelete - Parse URL for all the delete paths and call the controller action
/*
 DELETE  "/v1/channels/{id}"
*/

func (tc *ChannelController) processDelete(w http.ResponseWriter, r *http.Request, user *common.ContextData, requestID string, pathParts []string) {

	if (len(pathParts) == 3) && (pathParts[1] == "channels") {
		tc.DeleteChannel(w, r, pathParts[2], user, requestID)
	} else {
		common.RenderErrorJSON(w, "1000", "Invalid Request", 400, requestID)
		return
	}

}

// ShowChannel - used to view Channel
func (tc *ChannelController) ShowChannel(w http.ResponseWriter, r *http.Request, id string, user *common.ContextData, requestID string) {
	ctx := r.Context()

	select {
	case <-ctx.Done():
		common.RenderErrorJSON(w, "1002", "Client closed connection", 402, requestID)
		return
	default:

		channel, err := tc.Service.ShowChannel(ctx, id, user.UserID, user.Email, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": user.Email, "reqid": requestID, "msgnum": 5000}).Error(err)
			common.RenderErrorJSON(w, "5000", err.Error(), 402, requestID)
			return
		}

		common.RenderJSON(w, channel)
	}
}

// CreateChannel - used to Create Channel
func (tc *ChannelController) CreateChannel(w http.ResponseWriter, r *http.Request, user *common.ContextData, requestID string) {
	ctx := r.Context()

	select {
	case <-ctx.Done():
		common.RenderErrorJSON(w, "1002", "Client closed connection", 402, requestID)
		return
	default:
		form := msgservices.Channel{}
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&form)
		if err != nil {
			log.WithFields(log.Fields{"user": user.Email, "reqid": requestID, "msgnum": 5001}).Error(err)
			common.RenderErrorJSON(w, "5001", err.Error(), 402, requestID)
			return
		}
		v := common.NewValidator()
		v.IsStrLenBetMinMax("Channel Name", form.ChannelName, msgservices.ChannelNameLenMin, msgservices.ChannelNameLenMax)
		v.IsStrLenBetMinMax("Channel Description", form.ChannelDesc, msgservices.ChannelDescLenMin, msgservices.ChannelDescLenMax)
		if v.IsValid() {
			common.RenderErrorJSON(w, "5008", v.Error(), 402, requestID)
			return
		}
		channel, err := tc.Service.CreateChannel(ctx, &form, user.UserID, user.Email, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": user.Email, "reqid": requestID, "msgnum": 5002}).Error(err)
			common.RenderErrorJSON(w, "5002", err.Error(), 402, requestID)
			return
		}

		common.RenderJSON(w, channel)
	}
}

// GetChannelByName - used to get Channel by name
func (tc *ChannelController) GetChannelByName(w http.ResponseWriter, r *http.Request, user *common.ContextData, requestID string) {
	ctx := r.Context()

	select {
	case <-ctx.Done():
		common.RenderErrorJSON(w, "1002", "Client closed connection", 402, requestID)
		return
	default:
		form := msgservices.Channel{}
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&form)
		if err != nil {
			log.WithFields(log.Fields{"user": user.Email, "reqid": requestID, "msgnum": 5003}).Error(err)
			common.RenderErrorJSON(w, "5003", err.Error(), 402, requestID)
			return
		}
		topc, err := tc.Service.GetChannelByName(ctx, form.ChannelName, user.Email, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": user.Email, "reqid": requestID, "msgnum": 5004}).Error(err)
			common.RenderErrorJSON(w, "5004", err.Error(), 402, requestID)
			return
		}

		common.RenderJSON(w, topc)
	}
}

// UpdateChannel - Update channel
func (tc *ChannelController) UpdateChannel(w http.ResponseWriter, r *http.Request, id string, user *common.ContextData, requestID string) {
	ctx := r.Context()

	select {
	case <-ctx.Done():
		common.RenderErrorJSON(w, "1002", "Client closed connection", 402, requestID)
		return
	default:
		form := msgservices.Channel{}
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&form)
		if err != nil {
			log.WithFields(log.Fields{"user": user.Email, "reqid": requestID, "msgnum": 5005}).Error(err)
			common.RenderErrorJSON(w, "5005", err.Error(), 402, requestID)
			return
		}
		err = tc.Service.UpdateChannel(ctx, id, &form, user.UserID, user.Email, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": user.Email, "reqid": requestID, "msgnum": 5006}).Error(err)
			common.RenderErrorJSON(w, "5006", err.Error(), 402, requestID)
			return
		}

		common.RenderJSON(w, "Updated Successfully")
	}
}

// DeleteChannel - delete channel
func (tc *ChannelController) DeleteChannel(w http.ResponseWriter, r *http.Request, id string, user *common.ContextData, requestID string) {
	ctx := r.Context()

	select {
	case <-ctx.Done():
		common.RenderErrorJSON(w, "1002", "Client closed connection", 402, requestID)
		return
	default:
		err := tc.Service.DeleteChannel(ctx, id, user.Email, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": user.Email, "reqid": requestID, "msgnum": 5007}).Error(err)
			common.RenderErrorJSON(w, "5007", err.Error(), 402, requestID)
			return
		}

		common.RenderJSON(w, "Deleted Successfully")
	}
}
