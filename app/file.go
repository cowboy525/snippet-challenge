package app

import (
	"io"
	"time"

	"github.com/rwcarlsen/goexif/exif"
	"github.com/topoface/snippet-challenge/model"
	"github.com/topoface/snippet-challenge/services/filestore"
)

const (
	/*
	  EXIF Image Orientations
	  1        2       3      4         5            6           7          8

	  888888  888888      88  88      8888888888  88                  88  8888888888
	  88          88      88  88      88  88      88  88          88  88      88  88
	  8888      8888    8888  8888    88          8888888888  8888888888          88
	  88          88      88  88
	  88          88  888888  888888
	*/
	Upright            = 1
	UprightMirrored    = 2
	UpsideDown         = 3
	UpsideDownMirrored = 4
	RotatedCWMirrored  = 5
	RotatedCCW         = 6
	RotatedCCWMirrored = 7
	RotatedCW          = 8

	MaxImageSize         = 6048 * 4032 // 24 megapixels, roughly 36MB as a raw image
	ImageThumbnailWidth  = 120
	ImageThumbnailHeight = 100
	ImageThumbnailRatio  = float64(ImageThumbnailHeight) / float64(ImageThumbnailWidth)
	ImagePreviewWidth    = 1920
)

// FileBackend create and return new file backend
func (a *App) FileBackend() (filestore.FileBackend, *model.AppError) {
	return a.Srv().FileBackend()
}

// ReadFile read file from path
func (a *App) ReadFile(path string) ([]byte, *model.AppError) {
	backend, err := a.FileBackend()
	if err != nil {
		return nil, err
	}

	return backend.ReadFile(path)
}

// WriteFile write file to path
func (a *App) WriteFile(fr io.ReadSeeker, size int64, path string) (int64, *model.AppError) {
	backend, err := a.FileBackend()
	if err != nil {
		return 0, err
	}

	return backend.WriteFile(fr, size, path)
}

// RemoveFile remove file from path
func (a *App) RemoveFile(path *string) *model.AppError {
	if path == nil {
		return nil
	}

	backend, err := a.FileBackend()
	if err != nil {
		return err
	}

	return backend.RemoveFile(*path)
}

// RemoveDirectory remove files in the directory
func (a *App) RemoveDirectory(path string) *model.AppError {

	backend, err := a.FileBackend()
	if err != nil {
		return err
	}

	return backend.RemoveDirectory(path)
}

// GetSignedFileURL return signed url of file path
func (a *App) GetSignedFileURL(path *string, expire time.Time) (*string, *model.AppError) {
	if path == nil || len(*path) == 0 {
		return nil, nil
	}

	backend, err := a.FileBackend()
	if err != nil {
		return nil, err
	}

	return backend.GetSignedFileURL(*path, expire)
}

func getImageOrientation(input io.Reader) (int, error) {
	exifData, err := exif.Decode(input)
	if err != nil {
		return Upright, err
	}

	tag, err := exifData.Get("Orientation")
	if err != nil {
		return Upright, err
	}

	orientation, err := tag.Int(0)
	if err != nil {
		return Upright, err
	}

	return orientation, nil
}
