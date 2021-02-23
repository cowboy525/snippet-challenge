package app

import (
	"context"
	"net/http"

	"github.com/topoface/snippet-challenge/mlog"
	"github.com/topoface/snippet-challenge/model"
	"github.com/topoface/snippet-challenge/store"
)

// Iface : app interface
type Iface interface {
	Config() *model.Config
	Log() *mlog.Logger
	Handle404(w http.ResponseWriter, r *http.Request)
	Path() string
	SetContext(c context.Context)
	SetPath(s string)
	SetServer(srv *Server)
	Srv() *Server
	Store() *store.Store
}
