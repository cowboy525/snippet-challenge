package app

import (
	"context"
	"net/http"

	"github.com/ernie-mlg/ErniePJT-main-api-go/mlog"
	"github.com/ernie-mlg/ErniePJT-main-api-go/store"
)

type App struct {
	srv *Server

	store store.Store

	path    string
	context context.Context
}

func New(options ...AppOption) *App {
	app := &App{}

	for _, option := range options {
		option(app)
	}

	return app
}

// DO NOT CALL THIS.
// This is to avoid having to change all the code in cmd/commands/* for now
// shutdown should be called directly on the server
func (a *App) Shutdown() {
	a.Srv().Shutdown()
	a.srv = nil
}

func (a *App) Handle404(w http.ResponseWriter, r *http.Request) {
	mlog.Debug("not found handler triggered", mlog.String("path", r.URL.Path), mlog.Int("code", 404))

	http.NotFound(w, r)
}

func (a *App) Srv() *Server {
	return a.srv
}

func (a *App) Store() store.Store {
	if *a.Config().ServiceSettings.AtomicRequest {
		return a.store
	}
	return a.Srv().Store
}

func (a *App) Path() string {
	return a.path
}

func (a *App) SetPath(s string) {
	a.path = s
}
func (a *App) SetContext(c context.Context) {
	a.context = c
}
func (a *App) SetServer(srv *Server) {
	a.srv = srv
}
