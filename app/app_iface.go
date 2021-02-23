package app

import (
	"context"
	"io"
	"net/http"

	"github.com/topoface/snippet-challenge/model"
	"github.com/topoface/snippet-challenge/services/filestore"
	"github.com/topoface/snippet-challenge/store"
	"github.com/topoface/snippet-challenge/store/sqlstore/pagination"
)

// Iface : app interface
type Iface interface {
	AcceptLanguage() string
	Config() *model.Config
	CreateSpace(space *model.Space) (*model.Space, *model.AppError)
	CreateSpaceWithOwner(spaceRequest *model.SpaceRequest, avatar *model.SpaceRequestAvatar, ownerAvatar *model.SpaceRequestAvatar) (map[string]interface{}, *model.AppError)
	FileBackend() (filestore.FileBackend, *model.AppError)
	GetSpace(spaceID uint64) (*model.Space, *model.AppError)
	GetSpaces(options model.GetSpacesOptions) (*pagination.Paginator, *model.AppError)
	Handle404(w http.ResponseWriter, r *http.Request)
	Path() string
	ReadFile(path string) ([]byte, *model.AppError)
	RemoveFile(path *string) *model.AppError
	SetContext(c context.Context)
	SetPath(s string)
	SetServer(srv *Server)
	Srv() *Server
	Store() store.Store
	UpdateSpace(space *model.Space) (*model.Space, *model.AppError)
	WriteFile(fr io.ReadSeeker, size int64, path string) (int64, *model.AppError)
}
