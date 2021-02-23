package web

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/topoface/snippet-challenge/app"
	"github.com/topoface/snippet-challenge/mlog"
	"github.com/topoface/snippet-challenge/services/configservice"
)

// Web structure
type Web struct {
	GetGlobalAppOptions app.AppOptionCreator
	ConfigService       configservice.ConfigService
	MainRouter          *mux.Router
}

// New : creat new web instance
func New(config configservice.ConfigService, globalOptions app.AppOptionCreator, root *mux.Router) *Web {
	mlog.Debug("Initializing web routes")

	web := &Web{
		GetGlobalAppOptions: globalOptions,
		ConfigService:       config,
		MainRouter:          root,
	}

	web.InitStatic()

	return web
}

// ReturnStatusNoContent : return NoContent status
func ReturnStatusNoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}
