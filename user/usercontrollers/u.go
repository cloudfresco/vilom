package usercontrollers

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/cloudfresco/vilom/common"
	"github.com/cloudfresco/vilom/user/userservices"
)

// UController - create u controller
type UController struct {
	Service *userservices.UserService
}

// NewUController - create u handler
func NewUController(s *userservices.UserService) *UController {
	return &UController{s}
}

func (uc *UController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	requestID := common.GetRequestID()
	pathParts, _, err := common.ParseURL(r.URL.String())
	if err != nil {
		common.RenderErrorJSON(w, "1000", "Invalid Request", 400, requestID)
		return
	}

	switch r.Method {
	case http.MethodGet:

		/*
					   GET /v1/u/confirmation/:token
			       GET /v1/u/change_email/:token

		*/

		if (len(pathParts) == 4) && (pathParts[1] == "u") && (pathParts[2] == "confirmation") {
			uc.ConfirmEmail(w, r, pathParts[3], requestID)
		} else if (len(pathParts) == 4) && (pathParts[1] == "u") && (pathParts[2] == "change_email") {
			uc.ConfirmChangeEmail(w, r, pathParts[3], requestID)
		} else {
			common.RenderErrorJSON(w, "1000", "Invalid Request", 400, requestID)
			return
		}

	case http.MethodPost:

		/*
			    POST /v1/u/login
			    POST /v1/u/create
				  POST /v1/u/forgot_password
				  POST /v1/u/reset_password/:token
		*/

		if (len(pathParts) == 3) && (pathParts[1] == "u") && (pathParts[2] == "login") {
			uc.Login(w, r, requestID)
		} else if (len(pathParts) == 3) && (pathParts[1] == "u") && (pathParts[2] == "create") {
			uc.Create(w, r, requestID)
		} else if (len(pathParts) == 3) && (pathParts[1] == "u") && (pathParts[2] == "forgot_password") {
			uc.ForgotPassword(w, r, requestID)
		} else if (len(pathParts) == 4) && (pathParts[1] == "u") && (pathParts[2] == "reset_password") {
			uc.ConfirmForgotPassword(w, r, pathParts[3], requestID)
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

// Login - User logins
func (uc *UController) Login(w http.ResponseWriter, r *http.Request, requestID string) {
	ctx := r.Context()

	select {
	case <-ctx.Done():
		common.RenderErrorJSON(w, "1002", "Client closed connection", 402, requestID)
		return
	default:
		form := userservices.LoginForm{}
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&form)
		if err != nil {
			log.WithFields(log.Fields{
				"reqid":  requestID,
				"msgnum": 1100,
			}).Error(err)
			common.RenderErrorJSON(w, "1100", err.Error(), 402, requestID)
			return
		}
		user, err := uc.Service.Login(ctx, &form, requestID)
		if err != nil {
			log.WithFields(log.Fields{
				"reqid":  requestID,
				"msgnum": 1101,
			}).Error(err)
			common.RenderErrorJSON(w, "1101", err.Error(), 402, requestID)
			return
		}
		common.RenderJSON(w, user)
	}
}

// ConfirmEmail - Confirmation of email
func (uc *UController) ConfirmEmail(w http.ResponseWriter, r *http.Request, id string, requestID string) {
	ctx := r.Context()

	select {
	case <-ctx.Done():
		common.RenderErrorJSON(w, "1002", "Client closed connection", 402, requestID)
		return
	default:
		err := uc.Service.ConfirmEmail(ctx, id, requestID)
		if err != nil {
			log.WithFields(log.Fields{
				"reqid":  requestID,
				"msgnum": 1102,
			}).Error(err)
			common.RenderErrorJSON(w, "1102", err.Error(), 402, requestID)
			return
		}

		common.RenderJSON(w, "Your Account confirmed successfully")
	}
}

// Create - Create User
func (uc *UController) Create(w http.ResponseWriter, r *http.Request, requestID string) {
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
			log.WithFields(log.Fields{
				"reqid":  requestID,
				"msgnum": 1103,
			}).Error(err)
			common.RenderErrorJSON(w, "1103", err.Error(), 402, requestID)
			return
		}
		user, err := uc.Service.Create(ctx, &form, r.Host, requestID)
		if err != nil {
			log.WithFields(log.Fields{
				"reqid":  requestID,
				"msgnum": 1104,
			}).Error(err)
			common.RenderErrorJSON(w, "1104", err.Error(), 402, requestID)
			return
		}
		common.RenderJSON(w, user)
	}
}

// ForgotPassword - Send Link to reset password
func (uc *UController) ForgotPassword(w http.ResponseWriter, r *http.Request, requestID string) {
	ctx := r.Context()

	select {
	case <-ctx.Done():
		common.RenderErrorJSON(w, "1002", "Client closed connection", 402, requestID)
		return
	default:
		form := userservices.ForgotPasswordForm{}
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&form)
		if err != nil {
			log.WithFields(log.Fields{
				"reqid":  requestID,
				"msgnum": 1105,
			}).Error(err)
			common.RenderErrorJSON(w, "1105", err.Error(), 402, requestID)
			return
		}
		err = uc.Service.ForgotPassword(ctx, &form, r.Host, requestID)
		if err != nil {
			log.WithFields(log.Fields{
				"reqid":  requestID,
				"msgnum": 1106,
			}).Error(err)
			common.RenderErrorJSON(w, "1106", err.Error(), 402, requestID)
			return
		}

		common.RenderJSON(w, "Please Check your email and get token to reset password")
	}
}

// ConfirmForgotPassword - Reset password
func (uc *UController) ConfirmForgotPassword(w http.ResponseWriter, r *http.Request, id string, requestID string) {
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
				"reqid":  requestID,
				"msgnum": 1107,
			}).Error(err)
			common.RenderErrorJSON(w, "1107", err.Error(), 402, requestID)
			return
		}
		err = uc.Service.ConfirmForgotPassword(ctx, &form, id, requestID)
		if err != nil {
			log.WithFields(log.Fields{
				"reqid":  requestID,
				"msgnum": 1108,
			}).Error(err)
			common.RenderErrorJSON(w, "1108", err.Error(), 402, requestID)
			return
		}

		common.RenderJSON(w, "Your Password Changed successfully")
	}
}

// ConfirmChangeEmail - Confirm Change Email
func (uc *UController) ConfirmChangeEmail(w http.ResponseWriter, r *http.Request, id string, requestID string) {
	ctx := r.Context()

	select {
	case <-ctx.Done():
		common.RenderErrorJSON(w, "1002", "Client closed connection", 402, requestID)
		return
	default:
		err := uc.Service.ConfirmChangeEmail(ctx, id, requestID)
		if err != nil {
			log.WithFields(log.Fields{
				"reqid":  requestID,
				"msgnum": 1109,
			}).Error(err)
			common.RenderErrorJSON(w, "1109", err.Error(), 402, requestID)
			return
		}

		common.RenderJSON(w, "Your Account confirmed successfully")
	}
}
