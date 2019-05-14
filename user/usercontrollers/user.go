package usercontrollers

import (
	"encoding/json"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/palantir/stacktrace"

	"github.com/cloudfresco/vilom/common"
	"github.com/cloudfresco/vilom/user/userservices"
)

// UsersController - used for
type UsersController struct {
	Service *userservices.UserService
}

// NewUsersController - Used to create a users handler
func NewUsersController(s *userservices.UserService) *UsersController {
	return &UsersController{s}
}

func (uc *UsersController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
			     GET  "/v1/users/"
				   GET  "/v1/users/{id}"
		*/

		if (len(pathParts) == 2) && (pathParts[1] == "users") {
			limit := queryString.Get("limit")
			cursor := queryString.Get("cursor")
			uc.Index(w, r, limit, cursor, user, requestID)
		} else if (len(pathParts) == 3) && (pathParts[1] == "users") {
			uc.Show(w, r, pathParts[2], user, requestID)
		} else {
			common.RenderErrorJSON(w, "1000", "Invalid Request", 400, requestID)
			return
		}

	case http.MethodPost:

		/*
						     POST  "/v1/users/change_email"
								 POST  "/v1/users/change_password/{id}"
			           POST  "/v1/users/getuserbyemail"
		*/

		if (len(pathParts) == 3) && (pathParts[1] == "users") && (pathParts[2] == "change_email") {
			uc.ChangeEmail(w, r, user, requestID)
		} else if (len(pathParts) == 4) && (pathParts[1] == "users") && (pathParts[2] == "change_password") {
			uc.ChangePassword(w, r, pathParts[3], user, requestID)
		} else if (len(pathParts) == 3) && (pathParts[1] == "users") && (pathParts[2] == "getuserbyemail") {
			uc.Getuserbyemail(w, r, user, requestID)
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

// Index - Get Users
func (uc *UsersController) Index(w http.ResponseWriter, r *http.Request, limit string, cursor string, user *common.ContextData, requestID string) {
	AllowedRoles := []string{"co_admin"}
	err := common.CheckRoles(AllowedRoles, user.Roles)
	if err != nil {
		log.Error(stacktrace.Propagate(err, ""))
		common.RenderErrorJSON(w, "1150", "You are Not Authorised", 402, requestID)
		return
	}

	ctx := r.Context()

	select {
	case <-ctx.Done():
		common.RenderErrorJSON(w, "1002", "Client closed connection", 402, requestID)
		return
	default:
		users, err := uc.Service.GetUsers(limit, cursor)
		if err != nil {
			log.Error(stacktrace.Propagate(err, ""))
			common.RenderErrorJSON(w, "1151", err.Error(), 402, requestID)
			return
		}
		common.RenderJSON(w, users)
	}
}

// Show - Get User Details
func (uc *UsersController) Show(w http.ResponseWriter, r *http.Request, id string, user *common.ContextData, requestID string) {
	AllowedRoles := []string{"co_admin"}
	err := common.CheckRoles(AllowedRoles, user.Roles)
	if err != nil {
		log.Error(stacktrace.Propagate(err, ""))
		common.RenderErrorJSON(w, "1152", "You are Not Authorised", 402, requestID)
		return
	}

	ctx := r.Context()

	select {
	case <-ctx.Done():
		common.RenderErrorJSON(w, "1002", "Client closed connection", 402, requestID)
		return
	default:
		usr, err := uc.Service.GetUser(id)
		if err != nil {
			log.Error(stacktrace.Propagate(err, ""))
			common.RenderErrorJSON(w, "1153", err.Error(), 400, requestID)
			return
		}

		common.RenderJSON(w, usr)
	}

}

// ChangeEmail - Changes Email
func (uc *UsersController) ChangeEmail(w http.ResponseWriter, r *http.Request, user *common.ContextData, requestID string) {
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
			log.Error(stacktrace.Propagate(err, ""))
			common.RenderErrorJSON(w, "1154", err.Error(), 402, requestID)
			return
		}
		err = uc.Service.ChangeEmail(&form, r.Host)
		if err != nil {
			log.Error(stacktrace.Propagate(err, ""))
			common.RenderErrorJSON(w, "1155", err.Error(), 402, requestID)
			return
		}

		common.RenderJSON(w, "Your Email Changed successfully, Please Check your email and confirm your acoount")
	}
}

// ChangePassword - Changes Password
func (uc *UsersController) ChangePassword(w http.ResponseWriter, r *http.Request, id string, user *common.ContextData, requestID string) {
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
			log.Error(stacktrace.Propagate(err, ""))
			common.RenderErrorJSON(w, "1156", err.Error(), 402, requestID)
			return
		}
		err = uc.Service.ChangePassword(&form)
		if err != nil {
			log.Error(stacktrace.Propagate(err, ""))
			common.RenderErrorJSON(w, "1157", err.Error(), 402, requestID)
			return
		}

		common.RenderJSON(w, "Your Password Changed successfully")
	}
}

// Getuserbyemail - Get User By email
func (uc *UsersController) Getuserbyemail(w http.ResponseWriter, r *http.Request, user *common.ContextData, requestID string) {
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
			log.Error(stacktrace.Propagate(err, ""))
			common.RenderErrorJSON(w, "1158", err.Error(), 402, requestID)
			return
		}
		usr, err := uc.Service.GetUserByEmail(form.Email)
		if err != nil {
			log.Error(stacktrace.Propagate(err, ""))
			common.RenderErrorJSON(w, "1159", err.Error(), 402, requestID)
			return
		}
		common.RenderJSON(w, usr)
	}
}
