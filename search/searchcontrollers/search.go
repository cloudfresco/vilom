package searchcontrollers

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/cloudfresco/vilom/common"
	"github.com/cloudfresco/vilom/search/searchservices"
	"github.com/cloudfresco/vilom/user/userservices"
)

/* error message range: 7000-7299 */

// SearchController - Create Search Controller
type SearchController struct {
	Service  searchservices.SearchServiceIntf
	Serviceu userservices.UserServiceIntf
}

// NewSearchController - Create Search Handler
func NewSearchController(s searchservices.SearchServiceIntf, su userservices.UserServiceIntf) *SearchController {
	return &SearchController{
		Service:  s,
		Serviceu: su,
	}
}

func (sc *SearchController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user, requestID, err := sc.Serviceu.GetAuthUserDetails(r)
	if err != nil {
		common.RenderErrorJSON(w, "1001", err.Error(), 401, requestID)
		return
	}
	var pathParts []string

	path := r.URL.Path
	pathParts = common.GetPathParts(path)

	switch r.Method {
	case http.MethodGet:
	case http.MethodPost:
		/*
		   GET  "/v1/search/"
		*/
		if (len(pathParts) == 2) && (pathParts[1] == "search") {
			sc.Search(w, r, user, requestID)
		}
	case http.MethodPut:
	case http.MethodDelete:
	default:
		common.RenderErrorJSON(w, "1000", "Invalid Request", 400, requestID)
		return
	}

}

// Search - Search Channels
func (sc *SearchController) Search(w http.ResponseWriter, r *http.Request, user *common.ContextData, requestID string) {
	ctx := r.Context()

	select {
	case <-ctx.Done():
		common.RenderErrorJSON(w, "1002", "Client closed connection", 402, requestID)
		return
	default:
		form := searchservices.BleveForm{}
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&form)
		if err != nil {
			log.WithFields(log.Fields{
				"user":   user.Email,
				"reqid":  requestID,
				"msgnum": 7000,
			}).Error(err)
			common.RenderErrorJSON(w, "7000", err.Error(), 402, requestID)
			return
		}
		SearchResults, err := sc.Service.Search(&form, user.Email, requestID)
		if err != nil {
			log.WithFields(log.Fields{
				"user":   user.Email,
				"reqid":  requestID,
				"msgnum": 7001,
			}).Error(err)
			common.RenderErrorJSON(w, "7001", err.Error(), 402, requestID)
			return
		}

		common.RenderJSON(w, SearchResults.Hits)
	}
}
