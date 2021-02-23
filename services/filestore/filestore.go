package filestore

import (
	"io"
	"time"

	"github.com/topoface/snippet-challenge/model"
)

type ReadCloseSeeker interface {
	io.ReadCloser
	io.Seeker
}

type FileBackend interface {
	FileExists(path string) (bool, *model.AppError)
	ReadFile(path string) ([]byte, *model.AppError)
	WriteFile(fr io.ReadSeeker, size int64, path string) (int64, *model.AppError)
	RemoveFile(path string) *model.AppError
	RemoveDirectory(path string) *model.AppError
	GetSignedFileURL(path string, expire time.Time) (*string, *model.AppError)
}

func NewFileBackend(config *model.Config) (FileBackend, *model.AppError) {
	return &LocalFileBackend{
		baseUrl:   *config.ServiceSettings.SiteURL,
		directory: *config.FileSettings.Directory,
	}, nil
}
