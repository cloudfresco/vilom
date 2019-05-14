package usercontrollers

import (
	"encoding/json"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/palantir/stacktrace"

	"github.com/cloudfresco/vilom/common"
	"github.com/cloudfresco/vilom/user/userservices"
)

// UgroupController - Create Ugroup Controller
type UgroupController struct {
	Service *userservices.UgroupService
}

// NewUgroupController - Create Ugroup Handler
func NewUgroupController(s *userservices.UgroupService) *UgroupController {
	return &UgroupController{s}
}

func (uc *UgroupController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
						     GET  "/v1/ugroups/"
			           GET  "/v1/ugroups/topgroups"
							   GET  "/v1/ugroups/{id}"
			           GET  "/v1/ugroups/{id}/chdn"
			           GET  "/v1/ugroups/{id}/getparent"
		*/

		if (len(pathParts) == 1) && (pathParts[1] == "ugroups") {
			limit := queryString.Get("limit")
			cursor := queryString.Get("cursor")
			uc.Index(w, r, limit, cursor, user, requestID)
		} else if (len(pathParts) == 3) && (pathParts[1] == "ugroups") && (pathParts[2] == "topgroups") {
			uc.TopLevelGroups(w, r, user, requestID)
		} else if (len(pathParts) == 3) && (pathParts[1] == "ugroups") {
			uc.Show(w, r, pathParts[2], user, requestID)
		} else if (len(pathParts) == 4) && (pathParts[1] == "ugroups") && (pathParts[3] == "chdn") {
			uc.GetChdn(w, r, pathParts[2], user, requestID)
		} else if (len(pathParts) == 4) && (pathParts[1] == "ugroups") && (pathParts[3] == "getparent") {
			uc.GetParent(w, r, pathParts[2], user, requestID)
		} else {
			common.RenderErrorJSON(w, "1000", "Invalid Request", 400, requestID)
			return
		}

	case http.MethodPost:
		/*
			     POST  "/v1/ugroups/create"
			     POST  "/v1/ugroups/chdcreate"
					 POST  "/v1/ugroups/{id}/delete"
					 POST  "/v1/ugroups/{id}/adduser"
				   POST  "/v1/ugroups/{id}/deleteuser"
		*/

		if (len(pathParts) == 3) && (pathParts[1] == "ugroups") && (pathParts[2] == "create") {
			uc.Create(w, r, user, requestID)
		} else if (len(pathParts) == 3) && (pathParts[1] == "ugroups") && (pathParts[2] == "chdcreate") {
			uc.CreateChild(w, r, user, requestID)
		} else if (len(pathParts) == 4) && (pathParts[1] == "ugroups") && (pathParts[3] == "delete") {
			uc.Delete(w, r, pathParts[2], user, requestID)
		} else if (len(pathParts) == 4) && (pathParts[1] == "ugroups") && (pathParts[3] == "adduser") {
			uc.AddUserToGroup(w, r, pathParts[2], user, requestID)
		} else if (len(pathParts) == 4) && (pathParts[1] == "ugroups") && (pathParts[3] == "deleteuser") {
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

// Index - Get Ugroups
func (uc *UgroupController) Index(w http.ResponseWriter, r *http.Request, limit string, cursor string, user *common.ContextData, requestID string) {
	ctx := r.Context()

	select {
	case <-ctx.Done():
		common.RenderErrorJSON(w, "1002", "Client closed connection", 402, requestID)
		return
	default:
		ugroups, err := uc.Service.GetUgroups(limit, cursor)
		if err != nil {
			log.Error(stacktrace.Propagate(err, ""))
			common.RenderErrorJSON(w, "1200", err.Error(), 402, requestID)
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
		ugroups, err := uc.Service.TopLevelUgroups()
		if err != nil {
			log.Error(stacktrace.Propagate(err, ""))
			common.RenderErrorJSON(w, "1201", err.Error(), 402, requestID)
			return
		}

		common.RenderJSON(w, ugroups)
	}
}

// Show - Get ugroup details
func (uc *UgroupController) Show(w http.ResponseWriter, r *http.Request, id string, user *common.ContextData, requestID string) {
	ctx := r.Context()

	select {
	case <-ctx.Done():
		common.RenderErrorJSON(w, "1002", "Client closed connection", 402, requestID)
		return
	default:
		ugroup, err := uc.Service.GetUgroup(id)
		if err != nil {
			log.Error(stacktrace.Propagate(err, ""))
			common.RenderErrorJSON(w, "1202", err.Error(), 402, requestID)
			return
		}

		common.RenderJSON(w, ugroup)
	}
}

// Create - Create Ugroup
func (uc *UgroupController) Create(w http.ResponseWriter, r *http.Request, user *common.ContextData, requestID string) {
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
			log.Error(stacktrace.Propagate(err, ""))
			common.RenderErrorJSON(w, "1203", err.Error(), 402, requestID)
			return
		}
		ugroup, err := uc.Service.Create(&form)
		if err != nil {
			log.Error(stacktrace.Propagate(err, ""))
			common.RenderErrorJSON(w, "1204", err.Error(), 402, requestID)
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
			log.Error(stacktrace.Propagate(err, ""))
			common.RenderErrorJSON(w, "1205", err.Error(), 402, requestID)
			return
		}
		ugroup, err := uc.Service.CreateChild(&form)
		if err != nil {
			log.Error(stacktrace.Propagate(err, ""))
			common.RenderErrorJSON(w, "1206", err.Error(), 402, requestID)
			return
		}

		common.RenderJSON(w, ugroup)
	}
}

// Delete - delete ugroup
func (uc *UgroupController) Delete(w http.ResponseWriter, r *http.Request, id string, user *common.ContextData, requestID string) {
	ctx := r.Context()

	select {
	case <-ctx.Done():
		common.RenderErrorJSON(w, "1002", "Client closed connection", 402, requestID)
		return
	default:
		err := uc.Service.Delete(id)
		if err != nil {
			log.Error(stacktrace.Propagate(err, ""))
			common.RenderErrorJSON(w, "1207", err.Error(), 402, requestID)
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
			log.Error(stacktrace.Propagate(err, ""))
			common.RenderErrorJSON(w, "1208", err.Error(), 402, requestID)
			return
		}
		err = uc.Service.AddUserToGroup(&form, id)
		if err != nil {
			log.Error(stacktrace.Propagate(err, ""))
			common.RenderErrorJSON(w, "1209", err.Error(), 402, requestID)
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
			log.Error(stacktrace.Propagate(err, ""))
			common.RenderErrorJSON(w, "1210", err.Error(), 402, requestID)
			return
		}
		err = uc.Service.DeleteUserFromGroup(&form, id)
		if err != nil {
			log.Error(stacktrace.Propagate(err, ""))
			common.RenderErrorJSON(w, "1211", err.Error(), 402, requestID)
			return
		}

		common.RenderJSON(w, "User removed Successfully")
	}
}

// GetChdn - Get children of ugroup
func (uc *UgroupController) GetChdn(w http.ResponseWriter, r *http.Request, id string, user *common.ContextData, requestID string) {
	ctx := r.Context()

	select {
	case <-ctx.Done():
		common.RenderErrorJSON(w, "1002", "Client closed connection", 402, requestID)
		return
	default:
		ugroups, err := uc.Service.GetChildUgroups(id)
		if err != nil {
			log.Error(stacktrace.Propagate(err, ""))
			common.RenderErrorJSON(w, "1212", err.Error(), 402, requestID)
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
		ugroups, err := uc.Service.GetParent(id)
		if err != nil {
			log.Error(stacktrace.Propagate(err, ""))
			common.RenderErrorJSON(w, "1213", err.Error(), 402, requestID)
			return
		}

		common.RenderJSON(w, ugroups)
	}
}
