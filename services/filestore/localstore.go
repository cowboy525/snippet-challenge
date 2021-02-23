package filestore

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/topoface/snippet-challenge/model"
)

type LocalFileBackend struct {
	baseUrl   string
	directory string
}

func (b *LocalFileBackend) ReadFile(path string) ([]byte, *model.AppError) {
	f, err := ioutil.ReadFile(filepath.Join(b.directory, path))
	if err != nil {
		return nil, model.NewAppError("ReadFile", "services.file.read_file.reading_local", nil, err.Error(), http.StatusInternalServerError)
	}
	return f, nil
}

func (b *LocalFileBackend) FileExists(path string) (bool, *model.AppError) {
	_, err := os.Stat(filepath.Join(b.directory, path))

	if os.IsNotExist(err) {
		return false, nil
	}

	if err != nil {
		return false, model.NewAppError("FileExists", "services.file.file_exists.exists_local", nil, err.Error(), http.StatusInternalServerError)
	}
	return true, nil
}

func (b *LocalFileBackend) WriteFile(fr io.ReadSeeker, size int64, path string) (int64, *model.AppError) {
	if exists, _ := b.FileExists(path); exists {
		err := b.RemoveFile(path)
		if err != nil {
			return 0, err
		}
	}
	return writeFileLocally(fr, filepath.Join(b.directory, path))
}

func writeFileLocally(fr io.ReadSeeker, path string) (int64, *model.AppError) {
	if err := os.MkdirAll(filepath.Dir(path), 0750); err != nil {
		directory, _ := filepath.Abs(filepath.Dir(path))
		return 0, model.NewAppError("WriteFile", "services.file.write_file_locally.create_dir", nil, "directory="+directory+", err="+err.Error(), http.StatusInternalServerError)
	}
	fw, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return 0, model.NewAppError("WriteFile", "services.file.write_file_locally.writing", nil, err.Error(), http.StatusInternalServerError)
	}
	defer fw.Close()
	written, err := io.Copy(fw, fr)
	if err != nil {
		return written, model.NewAppError("WriteFile", "services.file.write_file_locally.writing", nil, err.Error(), http.StatusInternalServerError)
	}
	return written, nil
}

func (b *LocalFileBackend) RemoveFile(path string) *model.AppError {
	if err := os.Remove(filepath.Join(b.directory, path)); err != nil {
		return model.NewAppError("RemoveFile", "services.file.remove_file.local", nil, err.Error(), http.StatusInternalServerError)
	}
	return nil
}

func (b *LocalFileBackend) RemoveDirectory(path string) *model.AppError {
	if err := os.RemoveAll(filepath.Join(b.directory, path)); err != nil {
		return model.NewAppError("RemoveDirectory", "services.file.remove_directory.local", nil, err.Error(), http.StatusInternalServerError)
	}
	return nil
}

func (b *LocalFileBackend) GetSignedFileURL(path string, expire time.Time) (*string, *model.AppError) {
	tmp := strings.Split(path, "/")
	n := len(tmp)
	tmp[n-1] = url.QueryEscape(tmp[n-1])
	path = strings.Join(tmp, "/")

	u, _ := url.Parse(b.baseUrl)
	signedURL := u.Scheme + "://" + filepath.Join(u.Host, filepath.Join(b.directory, path))
	signedURL = strings.ReplaceAll(signedURL, "\\", "/")
	return &signedURL, nil
}
