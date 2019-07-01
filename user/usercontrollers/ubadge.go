package usercontrollers

import (
	"encoding/json"
	"net/http"
	"net/url"

	log "github.com/sirupsen/logrus"

	"github.com/cloudfresco/vilom/common"
	"github.com/cloudfresco/vilom/user/userservices"
)

/* error message range: 3000-3299 */

// UbadgeController - Create Ubadge Controller
type UbadgeController struct {
	Service *userservices.UbadgeService
}

// NewUbadgeController - Create Ubadge Handler
func NewUbadgeController(s *userservices.UbadgeService) *UbadgeController {
	return &UbadgeController{s}
}

// ServeHTTP - parse url and call controller action
func (uc *UbadgeController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user, requestID, err := common.GetAuthUserDetails(r, uc.Service.RedisClient, uc.Service.Db)
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
 GET  "/v1/ubadges/"
 GET  "/v1/ubadges/{id}"
*/

func (uc *UbadgeController) processGet(w http.ResponseWriter, r *http.Request, user *common.ContextData, requestID string, pathParts []string, queryString url.Values) {

	if (len(pathParts) == 2) && (pathParts[1] == "ubadges") {
		limit := queryString.Get("limit")
		cursor := queryString.Get("cursor")
		uc.Index(w, r, limit, cursor, user, requestID)
	} else if (len(pathParts) == 3) && (pathParts[1] == "ubadges") {
		uc.Show(w, r, pathParts[2], user, requestID)
	} else {
		common.RenderErrorJSON(w, "1000", "Invalid Request", 400, requestID)
		return
	}
}

// processPost - Parse URL for all the POST paths and call the controller action
/*
 POST  "/v1/ubadges/add"
 POST  "/v1/ubadges/{id}/adduser"
 POST  "/v1/ubadges/{id}/deleteuser"
*/

func (uc *UbadgeController) processPost(w http.ResponseWriter, r *http.Request, user *common.ContextData, requestID string, pathParts []string) {

	if (len(pathParts) == 3) && (pathParts[1] == "ubadges") {
		if pathParts[2] == "add" {
			uc.Create(w, r, user, requestID)
		} else {
			common.RenderErrorJSON(w, "1000", "Invalid Request", 400, requestID)
			return
		}
	} else if (len(pathParts) == 4) && (pathParts[1] == "ubadges") {
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
 PUT  "/v1/ubadges/{id}"
*/

func (uc *UbadgeController) processPut(w http.ResponseWriter, r *http.Request, user *common.ContextData, requestID string, pathParts []string) {

	if (len(pathParts) == 3) && (pathParts[1] == "ubadges") {
		uc.Update(w, r, pathParts[2], user, requestID)
	} else {
		common.RenderErrorJSON(w, "1000", "Invalid Request", 400, requestID)
		return
	}

}

// processDelete - Parse URL for all the delete paths and call the controller action
/*
 DELETE  "/v1/ubadges/{id}"
*/

func (uc *UbadgeController) processDelete(w http.ResponseWriter, r *http.Request, user *common.ContextData, requestID string, pathParts []string) {

	if (len(pathParts) == 3) && (pathParts[1] == "ubadges") {
		uc.Delete(w, r, pathParts[2], user, requestID)
	} else {
		common.RenderErrorJSON(w, "1000", "Invalid Request", 400, requestID)
		return
	}

}

// Index - Get Ubadges
func (uc *UbadgeController) Index(w http.ResponseWriter, r *http.Request, limit string, cursor string, user *common.ContextData, requestID string) {
	ctx := r.Context()

	select {
	case <-ctx.Done():
		common.RenderErrorJSON(w, "1002", "Client closed connection", 402, requestID)
		return
	default:
		ubadges, err := uc.Service.GetUbadges(ctx, limit, cursor, user.Email, requestID)
		if err != nil {
			log.WithFields(log.Fields{
				"user":   user.Email,
				"reqid":  requestID,
				"msgnum": 3000,
			}).Error(err)

			common.RenderErrorJSON(w, "3000", err.Error(), 402, requestID)
			return
		}
		common.RenderJSON(w, ubadges)
	}
}

// Show - Get Ubadge Details
func (uc *UbadgeController) Show(w http.ResponseWriter, r *http.Request, id string, user *common.ContextData, requestID string) {
	ctx := r.Context()

	select {
	case <-ctx.Done():
		common.RenderErrorJSON(w, "1002", "Client closed connection", 402, requestID)
		return
	default:
		ubadge, err := uc.Service.GetUbadge(ctx, id, user.Email, requestID)
		if err != nil {
			log.WithFields(log.Fields{
				"user":   user.Email,
				"reqid":  requestID,
				"msgnum": 3001,
			}).Error(err)
			common.RenderErrorJSON(w, "3001", err.Error(), 402, requestID)
			return
		}

		common.RenderJSON(w, ubadge)
	}
}

// Create - Create Ubadge
func (uc *UbadgeController) Create(w http.ResponseWriter, r *http.Request, user *common.ContextData, requestID string) {
	ctx := r.Context()

	select {
	case <-ctx.Done():
		common.RenderErrorJSON(w, "1002", "Client closed connection", 402, requestID)
		return
	default:
		form := userservices.Ubadge{}
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&form)
		if err != nil {
			log.WithFields(log.Fields{
				"user":   user.Email,
				"reqid":  requestID,
				"msgnum": 3002,
			}).Error(err)

			common.RenderErrorJSON(w, "3002", err.Error(), 402, requestID)
			return
		}
		ubadge, err := uc.Service.Create(ctx, &form, user.Email, requestID)
		if err != nil {
			log.WithFields(log.Fields{
				"user":   user.Email,
				"reqid":  requestID,
				"msgnum": 3003,
			}).Error(err)

			common.RenderErrorJSON(w, "3003", err.Error(), 402, requestID)
			return
		}

		common.RenderJSON(w, ubadge)
	}
}

// Delete - delete ubadge
func (uc *UbadgeController) Delete(w http.ResponseWriter, r *http.Request, id string, user *common.ContextData, requestID string) {
	ctx := r.Context()

	select {
	case <-ctx.Done():
		common.RenderErrorJSON(w, "1002", "Client closed connection", 402, requestID)
		return
	default:
		err := uc.Service.Delete(ctx, id, user.Email, requestID)
		if err != nil {
			log.WithFields(log.Fields{
				"user":   user.Email,
				"reqid":  requestID,
				"msgnum": 3004,
			}).Error(err)

			common.RenderErrorJSON(w, "3004", err.Error(), 402, requestID)
			return
		}

		common.RenderJSON(w, "Group Deleted successfully")
	}
}

// AddUserToGroup - Add user to Ubadge
func (uc *UbadgeController) AddUserToGroup(w http.ResponseWriter, r *http.Request, id string, user *common.ContextData, requestID string) {
	ctx := r.Context()

	select {
	case <-ctx.Done():
		common.RenderErrorJSON(w, "1002", "Client closed connection", 402, requestID)
		return
	default:
		form := userservices.UbadgeUser{}
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&form)
		if err != nil {
			log.WithFields(log.Fields{
				"user":   user.Email,
				"reqid":  requestID,
				"msgnum": 3005,
			}).Error(err)

			common.RenderErrorJSON(w, "3005", err.Error(), 402, requestID)
			return
		}
		err = uc.Service.AddUserToGroup(ctx, &form, id, user.Email, requestID)
		if err != nil {
			log.WithFields(log.Fields{
				"user":   user.Email,
				"reqid":  requestID,
				"msgnum": 3006,
			}).Error(err)

			common.RenderErrorJSON(w, "3006", err.Error(), 402, requestID)
			return
		}

		common.RenderJSON(w, "User Added Successfully")
	}
}

// DeleteUserFromGroup - delete user from Ubadge
func (uc *UbadgeController) DeleteUserFromGroup(w http.ResponseWriter, r *http.Request, id string, user *common.ContextData, requestID string) {
	ctx := r.Context()

	select {
	case <-ctx.Done():
		common.RenderErrorJSON(w, "1002", "Client closed connection", 402, requestID)
		return
	default:
		form := userservices.UbadgeUser{}
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&form)
		if err != nil {
			log.WithFields(log.Fields{
				"user":   user.Email,
				"reqid":  requestID,
				"msgnum": 3007,
			}).Error(err)
			common.RenderErrorJSON(w, "3007", err.Error(), 402, requestID)
			return
		}
		err = uc.Service.DeleteUserFromGroup(ctx, &form, id, user.Email, requestID)
		if err != nil {
			log.WithFields(log.Fields{
				"user":   user.Email,
				"reqid":  requestID,
				"msgnum": 3008,
			}).Error(err)

			common.RenderErrorJSON(w, "3008", err.Error(), 402, requestID)
			return
		}

		common.RenderJSON(w, "User removed Successfully")
	}
}

// Update - Update Ubadge
func (uc *UbadgeController) Update(w http.ResponseWriter, r *http.Request, id string, user *common.ContextData, requestID string) {
	ctx := r.Context()

	select {
	case <-ctx.Done():
		common.RenderErrorJSON(w, "1002", "Client closed connection", 402, requestID)
		return
	default:
		form := userservices.Ubadge{}
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&form)
		if err != nil {
			log.WithFields(log.Fields{"user": user.Email, "reqid": requestID, "msgnum": 3009}).Error(err)
			common.RenderErrorJSON(w, "3009", err.Error(), 402, requestID)
			return
		}
		err = uc.Service.Update(ctx, id, &form, user.UserID, user.Email, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": user.Email, "reqid": requestID, "msgnum": 3010}).Error(err)
			common.RenderErrorJSON(w, "3010", err.Error(), 402, requestID)
			return
		}

		common.RenderJSON(w, "Updated Successfully")
	}
}
