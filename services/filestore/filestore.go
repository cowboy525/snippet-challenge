package filestore

import (
	"io"
	"net/http"
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
	switch *config.FileSettings.DriverName {
	case model.IMAGE_DRIVER_S3:
		return &S3FileBackend{
			endpoint:    *config.AwsSettings.AwsS3CustomDomain,
			accessKey:   *config.AwsSettings.AwsAccessKeyID,
			secretKey:   *config.AwsSettings.AwsSecretAccessKey,
			secure:      config.AwsSettings.AwsS3SSL == nil || *config.AwsSettings.AwsS3SSL,
			region:      *config.AwsSettings.AwsS3RegionName,
			bucket:      *config.AwsSettings.AwsStorageBucketName,
			signKeyId:   *config.AwsSettings.AwsCloudFrontSignKeyID,
			signKeyPath: *config.AwsSettings.AwsCloudFrontSignPrivateKeyPath,
		}, nil
	case model.IMAGE_DRIVER_LOCAL:
		return &LocalFileBackend{
			baseUrl:   *config.ServiceSettings.SiteURL,
			directory: *config.FileSettings.Directory,
		}, nil
	}
	return nil, model.NewAppError("NewFileBackend", "services.file.no_driver", nil, "", http.StatusInternalServerError)
}
