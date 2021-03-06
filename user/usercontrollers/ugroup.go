package usercontrollers

import (
	"encoding/json"
	"net/http"
	"net/url"

	log "github.com/sirupsen/logrus"

	"github.com/cloudfresco/vilom/common"
	"github.com/cloudfresco/vilom/user/userservices"
)

/* error message range: 2000-2299 */

// UgroupController - Create Ugroup Controller
type UgroupController struct {
	Service  userservices.UgroupServiceIntf
	Serviceu userservices.UserServiceIntf
}

// NewUgroupController - Create Ugroup Handler
func NewUgroupController(s userservices.UgroupServiceIntf, su userservices.UserServiceIntf) *UgroupController {
	return &UgroupController{
		Service:  s,
		Serviceu: su,
	}
}

// ServeHTTP - parse url and call controller action
func (uc *UgroupController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user, requestID, err := uc.Serviceu.GetAuthUserDetails(r)
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
		uc.processGet(w, r, user, requestID, pathParts, queryString)
	case http.MethodPost:
		uc.processPost(w, r, user, requestID, pathParts)
	case http.MethodPut:
		uc.processPut(w, r, user, requestID, pathParts)
	case http.MethodDelete:
		uc.processDelete(w, r, user, requestID, pathParts)
	default:
		common.RenderErrorJSON(w, "1000", "Invalid Request", 400, requestID)
		return
	}
}

// processGet - Parse URL for all the GET paths and call the controller action
/*
 GET  "/v1/ugroups/"
 GET  "/v1/ugroups/topgroups"
 GET  "/v1/ugroups/{id}"
 GET  "/v1/ugroups/{id}/chdn"
 GET  "/v1/ugroups/{id}/getparent"
*/

func (uc *UgroupController) processGet(w http.ResponseWriter, r *http.Request, user *common.ContextData, requestID string, pathParts []string, queryString url.Values) {
	if (len(pathParts) == 2) && (pathParts[1] == "ugroups") {
		limit := queryString.Get("limit")
		cursor := queryString.Get("cursor")
		uc.GetUgroups(w, r, limit, cursor, user, requestID)
	} else if (len(pathParts) == 3) && (pathParts[1] == "ugroups") && (pathParts[2] == "topgroups") {
		uc.TopLevelGroups(w, r, user, requestID)
	} else if (len(pathParts) == 3) && (pathParts[1] == "ugroups") {
		uc.GetUgroup(w, r, pathParts[2], user, requestID)
	} else if (len(pathParts) == 4) && (pathParts[1] == "ugroups") {
		if pathParts[3] == "chdn" {
			uc.GetChildUgroups(w, r, pathParts[2], user, requestID)
		} else if pathParts[3] == "getparent" {
			uc.GetParent(w, r, pathParts[2], user, requestID)
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
 POST  "/v1/ugroups/create"
 POST  "/v1/ugroups/chdcreate"
 POST  "/v1/ugroups/{id}/adduser"
 POST  "/v1/ugroups/{id}/deleteuser"
*/
func (uc *UgroupController) processPost(w http.ResponseWriter, r *http.Request, user *common.ContextData, requestID string, pathParts []string) {
	if (len(pathParts) == 3) && (pathParts[1] == "ugroups") {
		if pathParts[2] == "create" {
			uc.CreateUgroup(w, r, user, requestID)
		} else if pathParts[2] == "chdcreate" {
			uc.CreateChild(w, r, user, requestID)
		} else {
			common.RenderErrorJSON(w, "1000", "Invalid Request", 400, requestID)
			return
		}
	} else if (len(pathParts) == 4) && (pathParts[1] == "ugroups") {
		if pathParts[3] == "adduser" {
			uc.AddUserToGroup(w, r, pathParts[2], user, requestID)
		} else if pathParts[3] == "deleteuser" {
			uc.DeleteUserFromGroup(w, r, pathParts[2], user, requestID)
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
 PUT  "/v1/ugroups/{id}"
*/

func (uc *UgroupController) processPut(w http.ResponseWriter, r *http.Request, user *common.ContextData, requestID string, pathParts []string) {

	if (len(pathParts) == 3) && (pathParts[1] == "ugroups") {
		uc.UpdateUgroup(w, r, pathParts[2], user, requestID)
	} else {
		common.RenderErrorJSON(w, "1000", "Invalid Request", 400, requestID)
		return
	}

}

// processDelete - Parse URL for all the delete paths and call the controller action
/*
 DELETE  "/v1/ugroups/{id}"
*/

func (uc *UgroupController) processDelete(w http.ResponseWriter, r *http.Request, user *common.ContextData, requestID string, pathParts []string) {

	if (len(pathParts) == 3) && (pathParts[1] == "ugroups") {
		uc.DeleteUgroup(w, r, pathParts[2], user, requestID)
	} else {
		common.RenderErrorJSON(w, "1000", "Invalid Request", 400, requestID)
		return
	}

}

// GetUgroups - Get Ugroups
func (uc *UgroupController) GetUgroups(w http.ResponseWriter, r *http.Request, limit string, cursor string, user *common.ContextData, requestID string) {
	ctx := r.Context()

	select {
	case <-ctx.Done():
		common.RenderErrorJSON(w, "1002", "Client closed connection", 402, requestID)
		return
	default:
		ugroups, err := uc.Service.GetUgroups(ctx, limit, cursor, user.Email, requestID)
		if err != nil {
			log.WithFields(log.Fields{
				"user":   user.Email,
				"reqid":  requestID,
				"msgnum": 2000,
			}).Error(err)
			common.RenderErrorJSON(w, "2000", err.Error(), 402, requestID)
			return
		}
		common.RenderJSON(w, ugroups)
	}
}

// TopLevelGroups - Get top level Groups
func (uc *UgroupController) TopLevelGroups(w http.ResponseWriter, r *http.Request, user *common.ContextData, requestID string) {
	ctx := r.Context()

	select {
	case <-ctx.Done():
		common.RenderErrorJSON(w, "1002", "Client closed connection", 402, requestID)
		return
	default:
		ugroups, err := uc.Service.TopLevelUgroups(ctx, user.Email, requestID)
		if err != nil {
			log.WithFields(log.Fields{
				"user":   user.Email,
				"reqid":  requestID,
				"msgnum": 2001,
			}).Error(err)
			common.RenderErrorJSON(w, "2001", err.Error(), 402, requestID)
			return
		}

		common.RenderJSON(w, ugroups)
	}
}

// GetUgroup - Get ugroup details
func (uc *UgroupController) GetUgroup(w http.ResponseWriter, r *http.Request, id string, user *common.ContextData, requestID string) {
	ctx := r.Context()

	select {
	case <-ctx.Done():
		common.RenderErrorJSON(w, "1002", "Client closed connection", 402, requestID)
		return
	default:
		ugroup, err := uc.Service.GetUgroup(ctx, id, user.Email, requestID)
		if err != nil {
			log.WithFields(log.Fields{
				"user":   user.Email,
				"reqid":  requestID,
				"msgnum": 2002,
			}).Error(err)
			common.RenderErrorJSON(w, "2002", err.Error(), 402, requestID)
			return
		}

		common.RenderJSON(w, ugroup)
	}
}

// CreateUgroup - Create Ugroup
func (uc *UgroupController) CreateUgroup(w http.ResponseWriter, r *http.Request, user *common.ContextData, requestID string) {
	ctx := r.Context()

	select {
	case <-ctx.Done():
		common.RenderErrorJSON(w, "1002", "Client closed connection", 402, requestID)
		return
	default:
		form := userservices.Ugroup{}
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&form)
		if err != nil {
			log.WithFields(log.Fields{
				"user":   user.Email,
				"reqid":  requestID,
				"msgnum": 2003,
			}).Error(err)
			common.RenderErrorJSON(w, "2003", err.Error(), 402, requestID)
			return
		}
		v := common.NewValidator()
		v.IsStrLenBetMinMax("Ugroup Name", form.UgroupName, userservices.UgroupNameLenMin, userservices.UgroupNameLenMax)
		v.IsStrLenBetMinMax("Ugroup Description", form.UgroupDesc, userservices.UgroupDescLenMin, userservices.UgroupDescLenMax)
		if v.IsValid() {
			common.RenderErrorJSON(w, "4012", v.Error(), 402, requestID)
			return
		}
		ugroup, err := uc.Service.CreateUgroup(ctx, &form, user.Email, requestID)
		if err != nil {
			log.WithFields(log.Fields{
				"user":   user.Email,
				"reqid":  requestID,
				"msgnum": 2004,
			}).Error(err)
			common.RenderErrorJSON(w, "2004", err.Error(), 402, requestID)
			return
		}

		common.RenderJSON(w, ugroup)
	}
}

// CreateChild - Create child of ugroup
func (uc *UgroupController) CreateChild(w http.ResponseWriter, r *http.Request, user *common.ContextData, requestID string) {
	ctx := r.Context()

	select {
	case <-ctx.Done():
		common.RenderErrorJSON(w, "1002", "Client closed connection", 402, requestID)
		return
	default:
		form := userservices.Ugroup{}
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&form)
		if err != nil {
			log.WithFields(log.Fields{
				"user":   user.Email,
				"reqid":  requestID,
				"msgnum": 2005,
			}).Error(err)
			common.RenderErrorJSON(w, "2005", err.Error(), 402, requestID)
			return
		}
		ugroup, err := uc.Service.CreateChild(ctx, &form, user.Email, requestID)
		if err != nil {
			log.WithFields(log.Fields{
				"user":   user.Email,
				"reqid":  requestID,
				"msgnum": 2006,
			}).Error(err)
			common.RenderErrorJSON(w, "2006", err.Error(), 402, requestID)
			return
		}

		common.RenderJSON(w, ugroup)
	}
}

// DeleteUgroup - delete ugroup
func (uc *UgroupController) DeleteUgroup(w http.ResponseWriter, r *http.Request, id string, user *common.ContextData, requestID string) {
	ctx := r.Context()

	select {
	case <-ctx.Done():
		common.RenderErrorJSON(w, "1002", "Client closed connection", 402, requestID)
		return
	default:
		err := uc.Service.DeleteUgroup(ctx, id, user.Email, requestID)
		if err != nil {
			log.WithFields(log.Fields{
				"user":   user.Email,
				"reqid":  requestID,
				"msgnum": 2007,
			}).Error(err)
			common.RenderErrorJSON(w, "2007", err.Error(), 402, requestID)
			return
		}

		common.RenderJSON(w, "Group Deleted Successfully")
	}
}

// AddUserToGroup - Add user to group
func (uc *UgroupController) AddUserToGroup(w http.ResponseWriter, r *http.Request, id string, user *common.ContextData, requestID string) {
	ctx := r.Context()

	select {
	case <-ctx.Done():
		common.RenderErrorJSON(w, "1002", "Client closed connection", 402, requestID)
		return
	default:
		form := userservices.UgroupUser{}
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&form)
		if err != nil {
			log.WithFields(log.Fields{
				"user":   user.Email,
				"reqid":  requestID,
				"msgnum": 2008,
			}).Error(err)
			common.RenderErrorJSON(w, "2008", err.Error(), 402, requestID)
			return
		}
		err = uc.Service.AddUserToGroup(ctx, &form, id, user.Email, requestID)
		if err != nil {
			log.WithFields(log.Fields{
				"user":   user.Email,
				"reqid":  requestID,
				"msgnum": 2009,
			}).Error(err)
			common.RenderErrorJSON(w, "2009", err.Error(), 402, requestID)
			return
		}

		common.RenderJSON(w, "User Added Successfully")
	}
}

// DeleteUserFromGroup - delete user from group
func (uc *UgroupController) DeleteUserFromGroup(w http.ResponseWriter, r *http.Request, id string, user *common.ContextData, requestID string) {
	ctx := r.Context()

	select {
	case <-ctx.Done():
		common.RenderErrorJSON(w, "1002", "Client closed connection", 402, requestID)
		return
	default:
		form := userservices.UgroupUser{}
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&form)
		if err != nil {
			log.WithFields(log.Fields{
				"user":   user.Email,
				"reqid":  requestID,
				"msgnum": 2010,
			}).Error(err)
			common.RenderErrorJSON(w, "2010", err.Error(), 402, requestID)
			return
		}
		err = uc.Service.DeleteUserFromGroup(ctx, &form, id, user.Email, requestID)
		if err != nil {
			log.WithFields(log.Fields{
				"user":   user.Email,
				"reqid":  requestID,
				"msgnum": 2011,
			}).Error(err)
			common.RenderErrorJSON(w, "2011", err.Error(), 402, requestID)
			return
		}

		common.RenderJSON(w, "User removed Successfully")
	}
}

// GetChildUgroups - Get children of ugroup
func (uc *UgroupController) GetChildUgroups(w http.ResponseWriter, r *http.Request, id string, user *common.ContextData, requestID string) {
	ctx := r.Context()

	select {
	case <-ctx.Done():
		common.RenderErrorJSON(w, "1002", "Client closed connection", 402, requestID)
		return
	default:
		ugroups, err := uc.Service.GetChildUgroups(ctx, id, user.Email, requestID)
		if err != nil {
			log.WithFields(log.Fields{
				"user":   user.Email,
				"reqid":  requestID,
				"msgnum": 2012,
			}).Error(err)
			common.RenderErrorJSON(w, "2012", err.Error(), 402, requestID)
			return
		}

		common.RenderJSON(w, ugroups)
	}
}

// GetParent - Get Parent ugroup of child ugroup
func (uc *UgroupController) GetParent(w http.ResponseWriter, r *http.Request, id string, user *common.ContextData, requestID string) {
	ctx := r.Context()

	select {
	case <-ctx.Done():
		common.RenderErrorJSON(w, "1002", "Client closed connection", 402, requestID)
		return
	default:
		ugroups, err := uc.Service.GetParent(ctx, id, user.Email, requestID)
		if err != nil {
			log.WithFields(log.Fields{
				"user":   user.Email,
				"reqid":  requestID,
				"msgnum": 2013,
			}).Error(err)
			common.RenderErrorJSON(w, "2013", err.Error(), 402, requestID)
			return
		}

		common.RenderJSON(w, ugroups)
	}
}

// UpdateUgroup - Update Ugroup
func (uc *UgroupController) UpdateUgroup(w http.ResponseWriter, r *http.Request, id string, user *common.ContextData, requestID string) {
	ctx := r.Context()

	select {
	case <-ctx.Done():
		common.RenderErrorJSON(w, "1002", "Client closed connection", 402, requestID)
		return
	default:
		form := userservices.Ugroup{}
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&form)
		if err != nil {
			log.WithFields(log.Fields{"user": user.Email, "reqid": requestID, "msgnum": 2014}).Error(err)
			common.RenderErrorJSON(w, "2014", err.Error(), 402, requestID)
			return
		}
		err = uc.Service.UpdateUgroup(ctx, id, &form, user.UserID, user.Email, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": user.Email, "reqid": requestID, "msgnum": 2015}).Error(err)
			common.RenderErrorJSON(w, "2015", err.Error(), 402, requestID)
			return
		}

		common.RenderJSON(w, "Updated Successfully")
	}
}
