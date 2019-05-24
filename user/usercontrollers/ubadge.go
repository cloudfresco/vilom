package usercontrollers

import (
	"encoding/json"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/palantir/stacktrace"

	"github.com/cloudfresco/vilom/common"
	"github.com/cloudfresco/vilom/user/userservices"
)

// UbadgeController - Create Ubadge Controller
type UbadgeController struct {
	Service *userservices.UbadgeService
}

// NewUbadgeController - Create Ubadge Handler
func NewUbadgeController(s *userservices.UbadgeService) *UbadgeController {
	return &UbadgeController{s}
}

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

		/*
			     GET  "/v1/ubadges/"
				   GET  "/v1/ubadges/{id}"
		*/

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

	case http.MethodPost:
		/*
			     POST  "/v1/ubadges/add"
					 POST  "/v1/ubadges/{id}/delete"
					 POST  "/v1/ubadges/{id}/adduser"
				   POST  "/v1/ubadges/{id}/deleteuser"
		*/

		if (len(pathParts) == 3) && (pathParts[1] == "ubadges") && (pathParts[2] == "add") {
			uc.Create(w, r, user, requestID)
		} else if (len(pathParts) == 4) && (pathParts[1] == "ubadges") && (pathParts[3] == "delete") {
			uc.Delete(w, r, pathParts[2], user, requestID)
		} else if (len(pathParts) == 4) && (pathParts[1] == "ubadges") && (pathParts[3] == "adduser") {
			uc.AddUserToGroup(w, r, pathParts[2], user, requestID)
		} else if (len(pathParts) == 4) && (pathParts[1] == "ubadges") && (pathParts[3] == "deleteuser") {
			uc.DeleteUserFromGroup(w, r, pathParts[2], user, requestID)
		} else {
			common.RenderErrorJSON(w, "1000", "Invalid Request", 400, requestID)
			return
		}

	case http.MethodPut:
	case http.MethodDelete:
	default:
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
		ubadges, err := uc.Service.GetUbadges(ctx, limit, cursor)
		if err != nil {
			log.Error(stacktrace.Propagate(err, ""))
			common.RenderErrorJSON(w, "1230", err.Error(), 402, requestID)
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
		ubadge, err := uc.Service.GetUbadge(ctx, id)
		if err != nil {
			log.Error(stacktrace.Propagate(err, ""))
			common.RenderErrorJSON(w, "1231", err.Error(), 402, requestID)
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
			log.Error(stacktrace.Propagate(err, ""))
			common.RenderErrorJSON(w, "1232", err.Error(), 402, requestID)
			return
		}
		ubadge, err := uc.Service.Create(ctx, &form)
		if err != nil {
			log.Error(stacktrace.Propagate(err, ""))
			common.RenderErrorJSON(w, "1233", err.Error(), 402, requestID)
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
		err := uc.Service.Delete(ctx, id)
		if err != nil {
			log.Error(stacktrace.Propagate(err, ""))
			common.RenderErrorJSON(w, "1234", err.Error(), 402, requestID)
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
			log.Error(stacktrace.Propagate(err, ""))
			common.RenderErrorJSON(w, "1235", err.Error(), 402, requestID)
			return
		}
		err = uc.Service.AddUserToGroup(ctx, &form, id)
		if err != nil {
			log.Error(stacktrace.Propagate(err, ""))
			common.RenderErrorJSON(w, "1236", err.Error(), 402, requestID)
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
			log.Error(stacktrace.Propagate(err, ""))
			common.RenderErrorJSON(w, "1237", err.Error(), 402, requestID)
			return
		}
		err = uc.Service.DeleteUserFromGroup(ctx, &form, id)
		if err != nil {
			log.Error(stacktrace.Propagate(err, ""))
			common.RenderErrorJSON(w, "1238", err.Error(), 402, requestID)
			return
		}

		common.RenderJSON(w, "User removed Successfully")
	}
}
