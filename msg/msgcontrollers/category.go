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

// CategoryController - Create Category Controller
type CategoryController struct {
	Service  msgservices.CategoryServiceIntf
	Serviceu userservices.UserServiceIntf
}

// NewCategoryController - Create Category Handler
func NewCategoryController(s msgservices.CategoryServiceIntf, su userservices.UserServiceIntf) *CategoryController {
	return &CategoryController{
		Service:  s,
		Serviceu: su,
	}
}

// ServeHTTP - parse url and call controller action
func (cc *CategoryController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
 GET  "/v1/categories/"
 GET  "/v1/categories/{id}"
 GET  "/v1/categories/topcats"
 GET  "/v1/categories/{id}/chdn"
 GET  "/v1/categories/{id}/getparent"
*/

func (cc *CategoryController) processGet(w http.ResponseWriter, r *http.Request, user *common.ContextData, requestID string, pathParts []string, queryString url.Values) {

	if (len(pathParts) == 2) && (pathParts[1] == "categories") {
		limit := queryString.Get("limit")
		cursor := queryString.Get("cursor")
		cc.GetCategories(w, r, limit, cursor, user, requestID)
	} else if len(pathParts) == 3 {
		if pathParts[2] == "topcats" {
			cc.GetTopLevelCategories(w, r, user, requestID)
		} else if pathParts[1] == "categories" {
			cc.GetCategoryWithTopics(w, r, pathParts[2], user, requestID)
		} else {
			common.RenderErrorJSON(w, "1000", "Invalid Request", 400, requestID)
			return
		}
	} else if (len(pathParts) == 4) && (pathParts[1] == "categories") {
		if pathParts[3] == "chdn" {
			cc.GetChildCategories(w, r, pathParts[2], user, requestID)
		} else if pathParts[3] == "getparent" {
			cc.GetParentCategory(w, r, pathParts[2], user, requestID)
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
 POST  "/v1/categories/create/"
 POST  "/v1/categories/chdcreate/"
*/
func (cc *CategoryController) processPost(w http.ResponseWriter, r *http.Request, user *common.ContextData, requestID string, pathParts []string) {
	if (len(pathParts) == 3) && (pathParts[1] == "categories") {
		if pathParts[2] == "create" {
			cc.CreateCategory(w, r, user, requestID)
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
 PUT  "/v1/categories/{id}"
*/

func (cc *CategoryController) processPut(w http.ResponseWriter, r *http.Request, user *common.ContextData, requestID string, pathParts []string) {

	if (len(pathParts) == 3) && (pathParts[1] == "categories") {
		cc.UpdateCategory(w, r, pathParts[2], user, requestID)
	} else {
		common.RenderErrorJSON(w, "1000", "Invalid Request", 400, requestID)
		return
	}

}

// processDelete - Parse URL for all the delete paths and call the controller action
/*
 DELETE  "/v1/categories/{id}"
*/

func (cc *CategoryController) processDelete(w http.ResponseWriter, r *http.Request, user *common.ContextData, requestID string, pathParts []string) {

	if (len(pathParts) == 3) && (pathParts[1] == "categories") {
		cc.DeleteCategory(w, r, pathParts[2], user, requestID)
	} else {
		common.RenderErrorJSON(w, "1000", "Invalid Request", 400, requestID)
		return
	}

}

// GetCategories - used to view all categories
func (cc *CategoryController) GetCategories(w http.ResponseWriter, r *http.Request, limit string, cursor string, user *common.ContextData, requestID string) {
	ctx := r.Context()

	select {
	case <-ctx.Done():
		common.RenderErrorJSON(w, "1002", "Client closed connection", 402, requestID)
		return
	default:
		categories, err := cc.Service.GetCategories(ctx, limit, cursor, user.Email, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": user.Email, "reqid": requestID, "msgnum": 4000}).Error(err)
			common.RenderErrorJSON(w, "4000", err.Error(), 402, requestID)
			return
		}

		common.RenderJSON(w, categories)
	}
}

// GetCategoryWithTopics - used to view category
func (cc *CategoryController) GetCategoryWithTopics(w http.ResponseWriter, r *http.Request, id string, user *common.ContextData, requestID string) {
	ctx := r.Context()

	select {
	case <-ctx.Done():
		common.RenderErrorJSON(w, "1002", "Client closed connection", 402, requestID)
		return
	default:
		category, err := cc.Service.GetCategoryWithTopics(ctx, id, user.Email, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": user.Email, "reqid": requestID, "msgnum": 4001}).Error(err)
			common.RenderErrorJSON(w, "4001", err.Error(), 402, requestID)
			return
		}

		common.RenderJSON(w, category)
	}
}

// CreateCategory - used to Create Category
func (cc *CategoryController) CreateCategory(w http.ResponseWriter, r *http.Request, user *common.ContextData, requestID string) {
	ctx := r.Context()

	select {
	case <-ctx.Done():
		common.RenderErrorJSON(w, "1002", "Client closed connection", 402, requestID)
		return
	default:
		form := msgservices.Category{}
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&form)
		if err != nil {
			log.WithFields(log.Fields{"user": user.Email, "reqid": requestID, "msgnum": 4002}).Error(err)
			common.RenderErrorJSON(w, "4002", err.Error(), 402, requestID)
			return
		}
		v := common.NewValidator()
		v.IsStrLenBetMinMax("Category Name", form.CategoryName, msgservices.CategoryNameLenMin, msgservices.CategoryNameLenMax)
		v.IsStrLenBetMinMax("Category Description", form.CategoryDesc, msgservices.CategoryDescLenMin, msgservices.CategoryDescLenMax)
		if v.IsValid() {
			common.RenderErrorJSON(w, "4012", v.Error(), 402, requestID)
			return
		}
		cat, err := cc.Service.CreateCategory(ctx, &form, user.UserID, user.Email, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": user.Email, "reqid": requestID, "msgnum": 4003}).Error(err)
			common.RenderErrorJSON(w, "4003", err.Error(), 402, requestID)
			return
		}

		common.RenderJSON(w, cat)
	}
}

// CreateChild - used to Create SubCategory
func (cc *CategoryController) CreateChild(w http.ResponseWriter, r *http.Request, user *common.ContextData, requestID string) {
	ctx := r.Context()

	select {
	case <-ctx.Done():
		common.RenderErrorJSON(w, "1002", "Client closed connection", 402, requestID)
		return
	default:
		form := msgservices.Category{}
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&form)
		if err != nil {
			log.WithFields(log.Fields{"user": user.Email, "reqid": requestID, "msgnum": 4004}).Error(err)
			common.RenderErrorJSON(w, "4004", err.Error(), 402, requestID)
			return
		}
		cat, err := cc.Service.CreateChild(ctx, &form, user.UserID, user.Email, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": user.Email, "reqid": requestID, "msgnum": 4005}).Error(err)
			common.RenderErrorJSON(w, "4005", err.Error(), 402, requestID)
			return
		}

		common.RenderJSON(w, cat)
	}
}

// GetTopLevelCategories - Get all top level categories
func (cc *CategoryController) GetTopLevelCategories(w http.ResponseWriter, r *http.Request, user *common.ContextData, requestID string) {
	ctx := r.Context()

	select {
	case <-ctx.Done():
		common.RenderErrorJSON(w, "1002", "Client closed connection", 402, requestID)
		return
	default:
		categories, err := cc.Service.GetTopLevelCategories(ctx, user.Email, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": user.Email, "reqid": requestID, "msgnum": 4006}).Error(err)
			common.RenderErrorJSON(w, "4006", err.Error(), 402, requestID)
			return
		}

		common.RenderJSON(w, categories)
	}
}

// GetChildCategories - Get children of category
func (cc *CategoryController) GetChildCategories(w http.ResponseWriter, r *http.Request, id string, user *common.ContextData, requestID string) {
	ctx := r.Context()

	select {
	case <-ctx.Done():
		common.RenderErrorJSON(w, "1002", "Client closed connection", 402, requestID)
		return
	default:
		categories, err := cc.Service.GetChildCategories(ctx, id, user.Email, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": user.Email, "reqid": requestID, "msgnum": 4007}).Error(err)
			common.RenderErrorJSON(w, "4007", err.Error(), 402, requestID)
			return
		}

		common.RenderJSON(w, categories)
	}
}

// GetParentCategory - Get parent category
func (cc *CategoryController) GetParentCategory(w http.ResponseWriter, r *http.Request, id string, user *common.ContextData, requestID string) {
	ctx := r.Context()

	select {
	case <-ctx.Done():
		common.RenderErrorJSON(w, "1002", "Client closed connection", 402, requestID)
		return
	default:
		category, err := cc.Service.GetParentCategory(ctx, id, user.Email, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": user.Email, "reqid": requestID, "msgnum": 4008}).Error(err)
			common.RenderErrorJSON(w, "4008", err.Error(), 402, requestID)
			return
		}

		common.RenderJSON(w, category)
	}
}

// UpdateCategory - Update Category
func (cc *CategoryController) UpdateCategory(w http.ResponseWriter, r *http.Request, id string, user *common.ContextData, requestID string) {
	ctx := r.Context()

	select {
	case <-ctx.Done():
		common.RenderErrorJSON(w, "1002", "Client closed connection", 402, requestID)
		return
	default:
		form := msgservices.Category{}
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&form)
		if err != nil {
			log.WithFields(log.Fields{"user": user.Email, "reqid": requestID, "msgnum": 4009}).Error(err)
			common.RenderErrorJSON(w, "4009", err.Error(), 402, requestID)
			return
		}
		err = cc.Service.UpdateCategory(ctx, id, &form, user.UserID, user.Email, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": user.Email, "reqid": requestID, "msgnum": 4010}).Error(err)
			common.RenderErrorJSON(w, "4010", err.Error(), 402, requestID)
			return
		}

		common.RenderJSON(w, "Updated Successfully")
	}
}

// DeleteCategory - delete category
func (cc *CategoryController) DeleteCategory(w http.ResponseWriter, r *http.Request, id string, user *common.ContextData, requestID string) {
	ctx := r.Context()

	select {
	case <-ctx.Done():
		common.RenderErrorJSON(w, "1002", "Client closed connection", 402, requestID)
		return
	default:
		err := cc.Service.DeleteCategory(ctx, id, user.Email, requestID)
		if err != nil {
			log.WithFields(log.Fields{"user": user.Email, "reqid": requestID, "msgnum": 4011}).Error(err)
			common.RenderErrorJSON(w, "4011", err.Error(), 402, requestID)
			return
		}

		common.RenderJSON(w, "Deleted Successfully")
	}
}
