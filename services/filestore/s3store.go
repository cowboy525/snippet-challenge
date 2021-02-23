package filestore

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudfront/sign"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/ernie-mlg/ErniePJT-main-api-go/mlog"
	"github.com/ernie-mlg/ErniePJT-main-api-go/model"
)

type S3FileBackend struct {
	endpoint    string
	accessKey   string
	secretKey   string
	secure      bool
	region      string
	bucket      string
	signKeyId   string
	signKeyPath string
}

func (b *S3FileBackend) FileExists(path string) (bool, *model.AppError) {
	bucket := aws.String(b.bucket)
	keyname := aws.String(path)

	sess := session.Must(session.NewSession())

	creds := credentials.NewStaticCredentials(b.accessKey, b.secretKey, "")

	svc := s3.New(sess, &aws.Config{
		Region:      aws.String(b.region),
		Credentials: creds,
	})

	params := &s3.HeadObjectInput{
		Bucket: bucket,
		Key:    keyname,
	}
	_, err := svc.HeadObject(params)
	if err != nil {
		return false, model.NewAppError("RemoveFile", "services.file.remove_file.s3", nil, err.Error(), http.StatusInternalServerError)
	}
	return true, nil
}

func (b *S3FileBackend) ReadFile(path string) ([]byte, *model.AppError) {
	bucket := aws.String(b.bucket)
	keyname := aws.String(path)

	sess := session.Must(session.NewSession())

	creds := credentials.NewStaticCredentials(b.accessKey, b.secretKey, "")

	svc := s3.New(sess, &aws.Config{
		Region:      aws.String(b.region),
		Credentials: creds,
	})

	params := &s3.GetObjectInput{
		Bucket: bucket,
		Key:    keyname,
	}
	resp, err := svc.GetObject(params)

	if err != nil {
		return nil, model.NewAppError("ReadFile", "services.file.read_file.s3", nil, err.Error(), http.StatusInternalServerError)
	}
	defer resp.Body.Close()

	s3objectBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, model.NewAppError("ReadFile", "services.file.read_file.s3", nil, err.Error(), http.StatusInternalServerError)
	}

	return s3objectBytes, nil
}

func (b *S3FileBackend) WriteFile(fr io.ReadSeeker, size int64, path string) (int64, *model.AppError) {
	bucket := aws.String(b.bucket)
	keyname := aws.String(path)

	sess := session.Must(session.NewSession())

	creds := credentials.NewStaticCredentials(b.accessKey, b.secretKey, "")

	svc := s3.New(sess, &aws.Config{
		Region:      aws.String(b.region),
		Credentials: creds,
	})

	var contentType string
	if ext := filepath.Ext(path); model.IsFileExtImage(ext) {
		contentType = model.GetImageMimeType(ext)
	} else {
		contentType = "binary/octet-stream"
	}

	params := &s3.PutObjectInput{
		Bucket:        bucket,
		Key:           keyname,
		ACL:           aws.String("bucket-owner-full-control"),
		Body:          fr,
		ContentLength: aws.Int64(size),
		ContentType:   &contentType,
	}
	_, err := svc.PutObject(params)

	if err != nil {
		return 0, model.NewAppError("WriteFile", "services.file.write_file.s3", nil, err.Error(), http.StatusInternalServerError)
	}
	return size, nil
}

func (b *S3FileBackend) RemoveFile(path string) *model.AppError {
	bucket := aws.String(b.bucket)
	keyname := aws.String(path)

	sess := session.Must(session.NewSession())

	creds := credentials.NewStaticCredentials(b.accessKey, b.secretKey, "")

	svc := s3.New(sess, &aws.Config{
		Region:      aws.String(b.region),
		Credentials: creds,
	})

	params := &s3.DeleteObjectInput{
		Bucket: bucket,
		Key:    keyname,
	}
	_, err := svc.DeleteObject(params)
	if err != nil {
		return model.NewAppError("RemoveFile", "services.file.remove_file.s3", nil, err.Error(), http.StatusInternalServerError)
	}
	return nil
}

func removeObjects(svc *s3.S3, bucket *string, contents []*s3.Object) {
	for _, obj := range contents {
		obj := obj

		params := &s3.DeleteObjectInput{
			Bucket: bucket,
			Key:    obj.Key,
		}
		svc.DeleteObject(params)
	}
}

func (b *S3FileBackend) RemoveDirectory(path string) *model.AppError {
	bucket := aws.String(b.bucket)
	keyname := aws.String(path)

	sess := session.Must(session.NewSession())

	creds := credentials.NewStaticCredentials(b.accessKey, b.secretKey, "")

	svc := s3.New(sess, &aws.Config{
		Region:      aws.String(b.region),
		Credentials: creds,
	})

	results, err := svc.ListObjects(&s3.ListObjectsInput{
		Bucket: bucket,
		Prefix: keyname,
	})
	if err != nil {
		return model.NewAppError("RemoveFile", "services.file.remove_directory.s3", nil, err.Error(), http.StatusInternalServerError)
	}
	go removeObjects(svc, bucket, results.Contents)

	return nil
}

func (b *S3FileBackend) GetSignedFileURL(path string, expire time.Time) (*string, *model.AppError) {
	tmp := strings.Split(path, "/")
	n := len(tmp)
	tmp[n-1] = url.QueryEscape(tmp[n-1])
	path = strings.Join(tmp, "/")

	privateKey, err := sign.LoadPEMPrivKeyFile(b.signKeyPath)
	if err != nil {
		mlog.Error("Couldn't load PEM file: ", mlog.Err(err))
		return nil, model.NewAppError("GetSignedFileURL", "services.file.sign_url.pem_load", nil, err.Error(), http.StatusInternalServerError)
	}
	signer := sign.NewURLSigner(b.signKeyId, privateKey)
	signedURL, err := signer.Sign("https://"+filepath.Join(b.endpoint, path), expire)
	if err != nil {
		mlog.Error("Couldn't get signed url: ", mlog.Err(err))
		return nil, model.NewAppError("GetSignedFileURL", "services.file.sign_url.sign", nil, err.Error(), http.StatusInternalServerError)
	}
	return &signedURL, nil
}
