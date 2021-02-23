package web

import (
	"net/http"

	"github.com/ernie-mlg/ErniePJT-main-api-go/app"
	"github.com/ernie-mlg/ErniePJT-main-api-go/mlog"
	"github.com/ernie-mlg/ErniePJT-main-api-go/model"
	"github.com/ernie-mlg/ErniePJT-main-api-go/services/configservice"
	"github.com/gorilla/mux"
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

// ReturnStatusOK : return OK status
func ReturnStatusOK(w http.ResponseWriter) {
	m := make(map[string]string)
	m[model.STATUS] = model.STATUS_OK
	w.Write([]byte(model.MapToJSON(m)))
}

// ReturnStatusNoContent : return NoContent status
func ReturnStatusNoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}
