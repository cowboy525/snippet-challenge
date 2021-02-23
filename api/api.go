package api

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/topoface/snippet-challenge/app"
	"github.com/topoface/snippet-challenge/model"
	"github.com/topoface/snippet-challenge/services/configservice"
	"github.com/topoface/snippet-challenge/web"
)

// Routes : define api routes
type Routes struct {
	Root    *mux.Router // ''
	APIRoot *mux.Router // ''

	Snippets *mux.Router // 'spaces'
}

// API structure
type API struct {
	ConfigService       configservice.ConfigService
	GetGlobalAppOptions app.AppOptionCreator
	BaseRoutes          *Routes
}

// Init : init api routes
func Init(configservice configservice.ConfigService, globalOptionsFunc app.AppOptionCreator, root *mux.Router) *API {
	api := &API{
		ConfigService:       configservice,
		GetGlobalAppOptions: globalOptionsFunc,
		BaseRoutes:          &Routes{},
	}

	api.BaseRoutes.Root = root
	api.BaseRoutes.APIRoot = root.PathPrefix(model.API_URL_SUFFIX).Subrouter()

	api.BaseRoutes.Snippets = api.BaseRoutes.APIRoot.PathPrefix("/snippets").Subrouter()

	api.InitSnippets()

	// root.Handle("/api/{anything:.*}", http.HandlerFunc(api.Handle404))

	return api
}

// Handle404 : handle requests to undefined endpoints
func (api *API) Handle404(w http.ResponseWriter, r *http.Request) {
	// web.Handle404(api.ConfigService, w, r)
}

// ReturnStatusNoContent : return NoContent status
var ReturnStatusNoContent = web.ReturnStatusNoContent
