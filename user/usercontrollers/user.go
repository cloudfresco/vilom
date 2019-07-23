package usercontrollers

import (
	"encoding/json"
	"net/http"
	"net/url"

	log "github.com/sirupsen/logrus"

	"github.com/cloudfresco/vilom/common"
	"github.com/cloudfresco/vilom/user/userservices"
)

/* error message range: 1300-1499 */

// UserController - used for
type UserController struct {
	Service userservices.UserServiceIntf
}

// NewUserController - Used to create a users handler
func NewUserController(s userservices.UserServiceIntf) *UserController {
	return &UserController{
		Service: s,
	}
}

// ServeHTTP - parse url and call controller action
func (uc *UserController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user, requestID, err := uc.Service.GetAuthUserDetails(r)
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
	GET  "/v1/users/"
	GET  "/v1/users/{id}"
*/

func (uc *UserController) processGet(w http.ResponseWriter, r *http.Request, user *common.ContextData, requestID string, pathParts []string, queryString url.Values) {

	if (len(pathParts) == 2) && (pathParts[1] == "users") {
		limit := queryString.Get("limit")
		cursor := queryString.Get("cursor")
		uc.GetUsers(w, r, limit, cursor, user, requestID)
	} else if (len(pathParts) == 3) && (pathParts[1] == "users") {
		uc.GetUser(w, r, pathParts[2], user, requestID)
	} else {
		common.RenderErrorJSON(w, "1000", "Invalid Request", 400, requestID)
		return
	}
}

// processPost - Parse URL for all the POST paths and call the controller action
/*
	POST  "/v1/users/change_email"
	POST  "/v1/users/change_password/{id}"
	POST  "/v1/users/getuserbyemail"
*/

func (uc *UserController) processPost(w http.ResponseWriter, r *http.Request, user *common.ContextData, requestID string, pathParts []string) {
	if (len(pathParts) == 3) && (pathParts[1] == "users") {
		if pathParts[2] == "change_email" {
			uc.ChangeEmail(w, r, user, requestID)
		} else if pathParts[2] == "getuserbyemail" {
			uc.GetUserByEmail(w, r, user, requestID)
		} else {
			common.RenderErrorJSON(w, "1000", "Invalid Request", 400, requestID)
			return
		}
	} else if (len(pathParts) == 4) && (pathParts[1] == "users") {
		if pathParts[2] == "change_password" {
			uc.ChangePassword(w, r, pathParts[3], user, requestID)
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
 PUT  "/v1/users/{id}"
*/

func (uc *UserController) processPut(w http.ResponseWriter, r *http.Request, user *common.ContextData, requestID string, pathParts []string) {

	if (len(pathParts) == 3) && (pathParts[1] == "users") {
		uc.UpdateUser(w, r, pathParts[2], user, requestID)
	} else {
		common.RenderErrorJSON(w, "1000", "Invalid Request", 400, requestID)
		return
	}

}

// processDelete - Parse URL for all the delete paths and call the controller action
/*
 DELETE  "/v1/users/{id}"
*/

func (uc *UserController) processDelete(w http.ResponseWriter, r *http.Request, user *common.ContextData, requestID string, pathParts []string) {

	if (len(pathParts) == 3) && (pathParts[1] == "users") {
		uc.DeleteUser(w, r, pathParts[2], user, requestID)
	} else {
		common.RenderErrorJSON(w, "1000", "Invalid Request", 400, requestID)
		return
	}

}

// GetUsers - Get Users
func (uc *UserController) GetUsers(w http.ResponseWriter, r *http.Request, limit string, cursor string, user *common.ContextData, requestID string) {
	AllowedRoles := []string{"co_admin"}
	err := common.CheckRoles(AllowedRoles, user.Roles)
	if err != nil {
		log.WithFields(log.Fields{
			"user":   user.Email,
			"reqid":  requestID,
			"msgnum": 1300,
		}).Error(err)
		common.RenderErrorJSON(w, "1300", "You are Not Authorised", 402, requestID)
		return
	}

	ctx := r.Context()

	select {
	case <-ctx.Done():
		common.RenderErrorJSON(w, "1002", "Client closed connection", 402, requestID)
		return
	default:
		users, err := uc.Service.GetUsers(ctx, limit, cursor, user.Email, requestID)
		if err != nil {
			log.WithFields(log.Fields{
				"user":   user.Email,
				"reqid":  requestID,
				"msgnum": 1301,
			}).Error(err)
			common.RenderErrorJSON(w, "1301", err.Error(), 402, requestID)
			return
		}
		common.RenderJSON(w, users)
	}
}

// GetUser - Get User Details
func (uc *UserController) GetUser(w http.ResponseWriter, r *http.Request, id string, user *common.ContextData, requestID string) {
	AllowedRoles := []string{"co_admin"}
	err := common.CheckRoles(AllowedRoles, user.Roles)
	if err != nil {
		log.WithFields(log.Fields{
			"user":   user.Email,
			"reqid":  requestID,
			"msgnum": 1302,
		}).Error(err)
		common.RenderErrorJSON(w, "1302", "You are Not Authorised", 402, requestID)
		return
	}

	ctx := r.Context()

	select {
	case <-ctx.Done():
		common.RenderErrorJSON(w, "1002", "Client closed connection", 402, requestID)
		return
	default:
		usr, err := uc.Service.GetUser(ctx, id, user.Email, requestID)
		if err != nil {
			log.WithFields(log.Fields{
				"user":   user.Email,
				"reqid":  requestID,
				"msgnum": 1303,
			}).Error(err)
			common.RenderErrorJSON(w, "1303", err.Error(), 400, requestID)
			return
		}

		common.RenderJSON(w, usr)
	}

}

// ChangeEmail - Changes Email
func (uc *UserController) ChangeEmail(w http.ResponseWriter, r *http.Request, user *common.ContextData, requestID string) {
	ctx := r.Context()

	select {
	case <-ctx.Done():
		common.RenderErrorJSON(w, "1002", "Client closed connection", 402, requestID)
		return
	default:
		form := userservices.ChangeEmailForm{}
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&form)
		if err != nil {
			log.WithFields(log.Fields{
				"user":   user.Email,
				"reqid":  requestID,
				"msgnum": 1304,
			}).Error(err)
			common.RenderErrorJSON(w, "1304", err.Error(), 402, requestID)
			return
		}
		err = uc.Service.ChangeEmail(ctx, &form, r.Host, user.Email, requestID)
		if err != nil {
			log.WithFields(log.Fields{
				"user":   user.Email,
				"reqid":  requestID,
				"msgnum": 1305,
			}).Error(err)
			common.RenderErrorJSON(w, "1305", err.Error(), 402, requestID)
			return
		}

		common.RenderJSON(w, "Your Email Changed successfully, Please Check your email and confirm your acoount")
	}
}

// ChangePassword - Changes Password
func (uc *UserController) ChangePassword(w http.ResponseWriter, r *http.Request, id string, user *common.ContextData, requestID string) {
	ctx := r.Context()

	select {
	case <-ctx.Done():
		common.RenderErrorJSON(w, "1002", "Client closed connection", 402, requestID)
		return
	default:
		form := userservices.PasswordForm{}
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&form)
		if err != nil {
			log.WithFields(log.Fields{
				"user":   user.Email,
				"reqid":  requestID,
				"msgnum": 1306,
			}).Error(err)
			common.RenderErrorJSON(w, "1306", err.Error(), 402, requestID)
			return
		}
		err = uc.Service.ChangePassword(ctx, &form, user.Email, requestID)
		if err != nil {
			log.WithFields(log.Fields{
				"user":   user.Email,
				"reqid":  requestID,
				"msgnum": 1307,
			}).Error(err)
			common.RenderErrorJSON(w, "1307", err.Error(), 402, requestID)
			return
		}

		common.RenderJSON(w, "Your Password Changed successfully")
	}
}

// GetUserByEmail - Get User By email
func (uc *UserController) GetUserByEmail(w http.ResponseWriter, r *http.Request, user *common.ContextData, requestID string) {
	ctx := r.Context()

	select {
	case <-ctx.Done():
		common.RenderErrorJSON(w, "1002", "Client closed connection", 402, requestID)
		return
	default:
		form := userservices.UserEmailForm{}
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&form)
		if err != nil {
			log.WithFields(log.Fields{
				"user":   user.Email,
				"reqid":  requestID,
				"msgnum": 1308,
			}).Error(err)
			common.RenderErrorJSON(w, "1308", err.Error(), 402, requestID)
			return
		}
		usr, err := uc.Service.GetUserByEmail(ctx, form.Email, user.Email, requestID)
		if err != nil {
			log.WithFields(log.Fields{
				"user":   user.Email,
				"reqid":  requestID,
				"msgnum": 1309,
			}).Error(err)
			common.RenderErrorJSON(w, "1309", err.Error(), 402, requestID)
			return
		}
		common.RenderJSON(w, usr)
	}
}

// UpdateUser - Update User
func (uc *UserController) UpdateUser(w http.ResponseWriter, r *http.Request, id string, user *common.ContextData, requestID string) {
	ctx := r.Context()

	select {
	case <-ctx.Done():
		common.RenderErrorJSON(w, "1002", "Client closed connection", 402, requestID)
		return
	default:
		form := userservices.User{}
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&form)
		if err != nil {
			log.WithFields(log.Fields{"user": user.Email, "reqid": requestID, "msgnum": 1310}).Error(err)
			common.RenderErrorJSON(w, "1310", err.Error(), 402, requestID)
			return
		}
		err = uc.Service.UpdateUser(ctx, id, &form, user.UserID, user.Email, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": user.Email, "reqid": requestID, "msgnum": 1311}).Error(err)
			common.RenderErrorJSON(w, "1311", err.Error(), 402, requestID)
			return
		}

		common.RenderJSON(w, "Updated Successfully")
	}
}

// DeleteUser - delete user
func (uc *UserController) DeleteUser(w http.ResponseWriter, r *http.Request, id string, user *common.ContextData, requestID string) {
	ctx := r.Context()

	select {
	case <-ctx.Done():
		common.RenderErrorJSON(w, "1002", "Client closed connection", 402, requestID)
		return
	default:
		err := uc.Service.DeleteUser(ctx, id, user.Email, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": user.Email, "reqid": requestID, "msgnum": 1312}).Error(err)
			common.RenderErrorJSON(w, "1312", err.Error(), 402, requestID)
			return
		}

		common.RenderJSON(w, "Deleted Successfully")
	}
}
