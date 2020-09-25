package msgcontrollers

import (
	"encoding/json"
	"net/http"
	"net/url"

	log "github.com/sirupsen/logrus"

	"github.com/cloudfresco/vilom/common"
	"github.com/cloudfresco/vilom/msg/msgservices"
	"github.com/cloudfresco/vilom/user/userservices"
)

/* error message range: 4000-4299 */

// WorkspaceController - Create Workspace Controller
type WorkspaceController struct {
	Service  msgservices.WorkspaceServiceIntf
	Serviceu userservices.UserServiceIntf
}

// NewWorkspaceController - Create Workspace Handler
func NewWorkspaceController(s msgservices.WorkspaceServiceIntf, su userservices.UserServiceIntf) *WorkspaceController {
	return &WorkspaceController{
		Service:  s,
		Serviceu: su,
	}
}

// ServeHTTP - parse url and call controller action
func (cc *WorkspaceController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user, requestID, err := cc.Serviceu.GetAuthUserDetails(r)
	if err != nil {
		common.RenderErrorJSON(w, "1001", err.Error(), 401, requestID)
		return
	}
	pathParts, queryString, err := common.ParseURL(r.URL.String())
	if err != nil {
		common.RenderErrorJSON(w, "1000", "Invalid Request", 400, requestID)
		return
	}

	switch r.Method {
	case http.MethodGet:
		cc.processGet(w, r, user, requestID, pathParts, queryString)
	case http.MethodPost:
		cc.processPost(w, r, user, requestID, pathParts)
	case http.MethodPut:
		cc.processPut(w, r, user, requestID, pathParts)
	case http.MethodDelete:
		cc.processDelete(w, r, user, requestID, pathParts)
	default:
		common.RenderErrorJSON(w, "1000", "Invalid Request", 400, requestID)
		return
	}

}

// processGet - Parse URL for all the GET paths and call the controller action
/*
 GET  "/v1/workspaces/"
 GET  "/v1/workspaces/{id}"
 GET  "/v1/workspaces/topworkspaces"
 GET  "/v1/workspaces/{id}/chdn"
 GET  "/v1/workspaces/{id}/getparent"
*/

func (cc *WorkspaceController) processGet(w http.ResponseWriter, r *http.Request, user *common.ContextData, requestID string, pathParts []string, queryString url.Values) {

	if (len(pathParts) == 2) && (pathParts[1] == "workspaces") {
		limit := queryString.Get("limit")
		cursor := queryString.Get("cursor")
		cc.GetWorkspaces(w, r, limit, cursor, user, requestID)
	} else if len(pathParts) == 3 {
		if pathParts[2] == "topworkspaces" {
			cc.GetTopLevelWorkspaces(w, r, user, requestID)
		} else if pathParts[1] == "workspaces" {
			cc.GetWorkspaceWithChannels(w, r, pathParts[2], user, requestID)
		} else {
			common.RenderErrorJSON(w, "1000", "Invalid Request", 400, requestID)
			return
		}
	} else if (len(pathParts) == 4) && (pathParts[1] == "workspaces") {
		if pathParts[3] == "chdn" {
			cc.GetChildWorkspaces(w, r, pathParts[2], user, requestID)
		} else if pathParts[3] == "getparent" {
			cc.GetParentWorkspace(w, r, pathParts[2], user, requestID)
		} else {
			common.RenderErrorJSON(w, "1000", "Invalid Request", 400, requestID)
			return
		}
	} else {
		common.RenderErrorJSON(w, "1000", "Invalid Request", 400, requestID)
		return
	}
}

// processPost - Parse URL for all the POST paths and call the controller action
/*
 POST  "/v1/workspaces/create/"
 POST  "/v1/workspaces/chdcreate/"
*/
func (cc *WorkspaceController) processPost(w http.ResponseWriter, r *http.Request, user *common.ContextData, requestID string, pathParts []string) {
	if (len(pathParts) == 3) && (pathParts[1] == "workspaces") {
		if pathParts[2] == "create" {
			cc.CreateWorkspace(w, r, user, requestID)
		} else if pathParts[2] == "chdcreate" {
			cc.CreateChild(w, r, user, requestID)
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
 PUT  "/v1/workspaces/{id}"
*/

func (cc *WorkspaceController) processPut(w http.ResponseWriter, r *http.Request, user *common.ContextData, requestID string, pathParts []string) {

	if (len(pathParts) == 3) && (pathParts[1] == "workspaces") {
		cc.UpdateWorkspace(w, r, pathParts[2], user, requestID)
	} else {
		common.RenderErrorJSON(w, "1000", "Invalid Request", 400, requestID)
		return
	}

}

// processDelete - Parse URL for all the delete paths and call the controller action
/*
 DELETE  "/v1/workspaces/{id}"
*/

func (cc *WorkspaceController) processDelete(w http.ResponseWriter, r *http.Request, user *common.ContextData, requestID string, pathParts []string) {

	if (len(pathParts) == 3) && (pathParts[1] == "workspaces") {
		cc.DeleteWorkspace(w, r, pathParts[2], user, requestID)
	} else {
		common.RenderErrorJSON(w, "1000", "Invalid Request", 400, requestID)
		return
	}

}

// GetWorkspaces - used to view all workspaces
func (cc *WorkspaceController) GetWorkspaces(w http.ResponseWriter, r *http.Request, limit string, cursor string, user *common.ContextData, requestID string) {
	ctx := r.Context()

	select {
	case <-ctx.Done():
		common.RenderErrorJSON(w, "1002", "Client closed connection", 402, requestID)
		return
	default:
		workspaces, err := cc.Service.GetWorkspaces(ctx, limit, cursor, user.Email, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": user.Email, "reqid": requestID, "msgnum": 4000}).Error(err)
			common.RenderErrorJSON(w, "4000", err.Error(), 402, requestID)
			return
		}

		common.RenderJSON(w, workspaces)
	}
}

// GetWorkspaceWithChannels - used to view workspace
func (cc *WorkspaceController) GetWorkspaceWithChannels(w http.ResponseWriter, r *http.Request, id string, user *common.ContextData, requestID string) {
	ctx := r.Context()

	select {
	case <-ctx.Done():
		common.RenderErrorJSON(w, "1002", "Client closed connection", 402, requestID)
		return
	default:
		workspace, err := cc.Service.GetWorkspaceWithChannels(ctx, id, user.Email, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": user.Email, "reqid": requestID, "msgnum": 4001}).Error(err)
			common.RenderErrorJSON(w, "4001", err.Error(), 402, requestID)
			return
		}

		common.RenderJSON(w, workspace)
	}
}

// CreateWorkspace - used to Create Workspace
func (cc *WorkspaceController) CreateWorkspace(w http.ResponseWriter, r *http.Request, user *common.ContextData, requestID string) {
	ctx := r.Context()

	select {
	case <-ctx.Done():
		common.RenderErrorJSON(w, "1002", "Client closed connection", 402, requestID)
		return
	default:
		form := msgservices.Workspace{}
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&form)
		if err != nil {
			log.WithFields(log.Fields{"user": user.Email, "reqid": requestID, "msgnum": 4002}).Error(err)
			common.RenderErrorJSON(w, "4002", err.Error(), 402, requestID)
			return
		}
		v := common.NewValidator()
		v.IsStrLenBetMinMax("Workspace Name", form.WorkspaceName, msgservices.WorkspaceNameLenMin, msgservices.WorkspaceNameLenMax)
		v.IsStrLenBetMinMax("Workspace Description", form.WorkspaceDesc, msgservices.WorkspaceDescLenMin, msgservices.WorkspaceDescLenMax)
		if v.IsValid() {
			common.RenderErrorJSON(w, "4012", v.Error(), 402, requestID)
			return
		}
		workspace, err := cc.Service.CreateWorkspace(ctx, &form, user.UserID, user.Email, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": user.Email, "reqid": requestID, "msgnum": 4003}).Error(err)
			common.RenderErrorJSON(w, "4003", err.Error(), 402, requestID)
			return
		}

		common.RenderJSON(w, workspace)
	}
}

// CreateChild - used to Create SubWorkspace
func (cc *WorkspaceController) CreateChild(w http.ResponseWriter, r *http.Request, user *common.ContextData, requestID string) {
	ctx := r.Context()

	select {
	case <-ctx.Done():
		common.RenderErrorJSON(w, "1002", "Client closed connection", 402, requestID)
		return
	default:
		form := msgservices.Workspace{}
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&form)
		if err != nil {
			log.WithFields(log.Fields{"user": user.Email, "reqid": requestID, "msgnum": 4004}).Error(err)
			common.RenderErrorJSON(w, "4004", err.Error(), 402, requestID)
			return
		}
		workspace, err := cc.Service.CreateChild(ctx, &form, user.UserID, user.Email, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": user.Email, "reqid": requestID, "msgnum": 4005}).Error(err)
			common.RenderErrorJSON(w, "4005", err.Error(), 402, requestID)
			return
		}

		common.RenderJSON(w, workspace)
	}
}

// GetTopLevelWorkspaces - Get all top level workspaces
func (cc *WorkspaceController) GetTopLevelWorkspaces(w http.ResponseWriter, r *http.Request, user *common.ContextData, requestID string) {
	ctx := r.Context()

	select {
	case <-ctx.Done():
		common.RenderErrorJSON(w, "1002", "Client closed connection", 402, requestID)
		return
	default:
		workspaces, err := cc.Service.GetTopLevelWorkspaces(ctx, user.Email, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": user.Email, "reqid": requestID, "msgnum": 4006}).Error(err)
			common.RenderErrorJSON(w, "4006", err.Error(), 402, requestID)
			return
		}

		common.RenderJSON(w, workspaces)
	}
}

// GetChildWorkspaces - Get children of workspace
func (cc *WorkspaceController) GetChildWorkspaces(w http.ResponseWriter, r *http.Request, id string, user *common.ContextData, requestID string) {
	ctx := r.Context()

	select {
	case <-ctx.Done():
		common.RenderErrorJSON(w, "1002", "Client closed connection", 402, requestID)
		return
	default:
		workspaces, err := cc.Service.GetChildWorkspaces(ctx, id, user.Email, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": user.Email, "reqid": requestID, "msgnum": 4007}).Error(err)
			common.RenderErrorJSON(w, "4007", err.Error(), 402, requestID)
			return
		}

		common.RenderJSON(w, workspaces)
	}
}

// GetParentWorkspace - Get parent workspace
func (cc *WorkspaceController) GetParentWorkspace(w http.ResponseWriter, r *http.Request, id string, user *common.ContextData, requestID string) {
	ctx := r.Context()

	select {
	case <-ctx.Done():
		common.RenderErrorJSON(w, "1002", "Client closed connection", 402, requestID)
		return
	default:
		workspace, err := cc.Service.GetParentWorkspace(ctx, id, user.Email, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": user.Email, "reqid": requestID, "msgnum": 4008}).Error(err)
			common.RenderErrorJSON(w, "4008", err.Error(), 402, requestID)
			return
		}

		common.RenderJSON(w, workspace)
	}
}

// UpdateWorkspace - Update Workspace
func (cc *WorkspaceController) UpdateWorkspace(w http.ResponseWriter, r *http.Request, id string, user *common.ContextData, requestID string) {
	ctx := r.Context()

	select {
	case <-ctx.Done():
		common.RenderErrorJSON(w, "1002", "Client closed connection", 402, requestID)
		return
	default:
		form := msgservices.Workspace{}
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&form)
		if err != nil {
			log.WithFields(log.Fields{"user": user.Email, "reqid": requestID, "msgnum": 4009}).Error(err)
			common.RenderErrorJSON(w, "4009", err.Error(), 402, requestID)
			return
		}
		err = cc.Service.UpdateWorkspace(ctx, id, &form, user.UserID, user.Email, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": user.Email, "reqid": requestID, "msgnum": 4010}).Error(err)
			common.RenderErrorJSON(w, "4010", err.Error(), 402, requestID)
			return
		}

		common.RenderJSON(w, "Updated Successfully")
	}
}

// DeleteWorkspace - delete workspace
func (cc *WorkspaceController) DeleteWorkspace(w http.ResponseWriter, r *http.Request, id string, user *common.ContextData, requestID string) {
	ctx := r.Context()

	select {
	case <-ctx.Done():
		common.RenderErrorJSON(w, "1002", "Client closed connection", 402, requestID)
		return
	default:
		err := cc.Service.DeleteWorkspace(ctx, id, user.Email, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": user.Email, "reqid": requestID, "msgnum": 4011}).Error(err)
			common.RenderErrorJSON(w, "4011", err.Error(), 402, requestID)
			return
		}

		common.RenderJSON(w, "Deleted Successfully")
	}
}
