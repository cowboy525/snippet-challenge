package model

import (
	"encoding/json"
	"io"
	"math"
	"net"
	"net/http"
	"net/url"
	"strconv"

	"github.com/ernie-mlg/ErniePJT-main-api-go/mlog"
)

const (
	AWS_ACCESS_KEY_ID          = "xxxxx"
	AWS_SECRET_ACCESS_KEY      = "xxxxx"
	TEST_AWS_ACCESS_KEY_ID     = "xxxxx"
	TEST_AWS_SECRET_ACCESS_KEY = "xxxxx"

	AWS_SNS_REGION_NAME = "us-west-2"

	AWS_S3_REGION_NAME           = "ap-northeast-1"
	AWS_STORAGE_BUCKET_NAME      = "file.ernie.works"
	ENABLED_S3_TEST              = 0
	TEST_AWS_STORAGE_BUCKET_NAME = "file-test.ernie.works"

	AWS_CLOUDFRONT_SIGN_KEY_ID           = "xxxxx"
	AWS_CLOUDFRONT_SIGN_PRIVATE_KEY_PATH = ".keys/cloudfront.pem"
	AWS_S3_CUSTOM_DOMAIN                 = "file.ernie.works"
	TEST_AWS_S3_CUSTOM_DOMAIN            = "file-test.ernie.works"

	ELASTICSEARCH_CONNECTION_URL = "http://localhost:9200"

	CONN_SECURITY_NONE     = ""
	CONN_SECURITY_PLAIN    = "PLAIN"
	CONN_SECURITY_TLS      = "TLS"
	CONN_SECURITY_STARTTLS = "STARTTLS"

	IMAGE_DRIVER_LOCAL = "local"
	IMAGE_DRIVER_S3    = "amazons3"

	FILE_SETTINGS_DEFAULT_DIRECTORY = "./data/"

	DATABASE_DRIVER_SQLITE   = "sqlite3"
	DATABASE_DRIVER_MYSQL    = "mysql"
	DATABASE_DRIVER_POSTGRES = "postgres"

	SQL_SETTINGS_DEFAULT_DATA_SOURCE = "admin:admin@tcp(localhost:3306)/pjmtool_go?charset=utf8mb4,utf8&readTimeout=30s&writeTimeout=30s&parseTime=true"

	EMAIL_SETTINGS_DEFAULT_FEEDBACK_ORGANIZATION = ""

	SERVICE_SETTINGS_DEFAULT_SITE_URL           = "http://localhost:13000"
	SERVICE_SETTINGS_DEFAULT_LISTEN_AND_ADDRESS = ":13000"

	FAKE_SETTING = "********************************"
)

// ServiceSettings structure
type ServiceSettings struct {
	SiteURL                *string `restricted:"true"`
	ListenAddress          *string `restricted:"true"`
	ConnectionSecurity     *string `restricted:"true"`
	EnableDeveloper        *bool   `restricted:"true"`
	SessionCacheInMinutes  *int    `restricted:"true"`
	SessionLengthWebInDays *int    `restricted:"true"`
	AtomicRequest          *bool   `restricted:"true"`
}

// SetDefaults sets default service settings
func (s *ServiceSettings) SetDefaults() {
	if s.SiteURL == nil {
		if s.EnableDeveloper != nil && *s.EnableDeveloper {
			s.SiteURL = NewString(SERVICE_SETTINGS_DEFAULT_SITE_URL)
		} else {
			s.SiteURL = NewString("")
		}
	}

	if s.ListenAddress == nil {
		s.ListenAddress = NewString(SERVICE_SETTINGS_DEFAULT_LISTEN_AND_ADDRESS)
	}

	if s.ConnectionSecurity == nil {
		s.ConnectionSecurity = NewString("")
	}

	if s.EnableDeveloper == nil {
		s.EnableDeveloper = NewBool(false)
	}

	if s.SessionCacheInMinutes == nil {
		s.SessionCacheInMinutes = NewInt(10)
	}

	if s.SessionLengthWebInDays == nil {
		s.SessionLengthWebInDays = NewInt(180)
	}

	if s.AtomicRequest == nil {
		s.AtomicRequest = NewBool(false)
	}
}

func (s *ServiceSettings) isValid() *AppError {
	if !(*s.ConnectionSecurity == CONN_SECURITY_NONE || *s.ConnectionSecurity == CONN_SECURITY_TLS) {
		return NewAppError("Config.IsValid", "model.config.is_valid.webserver_security.app_error", nil, "", http.StatusBadRequest)
	}

	if len(*s.SiteURL) != 0 {
		if _, err := url.ParseRequestURI(*s.SiteURL); err != nil {
			return NewAppError("Config.IsValid", "model.config.is_valid.site_url.app_error", nil, "", http.StatusBadRequest)
		}
	}

	host, port, _ := net.SplitHostPort(*s.ListenAddress)
	var isValidHost bool
	if host == "" {
		isValidHost = true
	} else {
		isValidHost = (net.ParseIP(host) != nil) || IsDomainName(host)
	}
	portInt, err := strconv.Atoi(port)
	if err != nil || !isValidHost || portInt < 0 || portInt > math.MaxUint16 {
		return NewAppError("Config.IsValid", "model.config.is_valid.listen_address.app_error", nil, "", http.StatusBadRequest)
	}

	return nil
}

// ConfigFunc : config func
type ConfigFunc func() *Config

// Config structure
type Config struct {
	FileSettings    FileSettings
	ServiceSettings ServiceSettings
	LogSettings     LogSettings
}

// Clone creates clone of config
func (o *Config) Clone() *Config {
	var ret Config
	if err := json.Unmarshal([]byte(o.ToJSON()), &ret); err != nil {
		panic(err)
	}
	return &ret
}

// ToJSON convert a Config to a json string
func (o *Config) ToJSON() string {
	b, _ := json.Marshal(o)
	return string(b)
}

// ConfigFromJSON will decode the input and return a Config
func ConfigFromJSON(data io.Reader) *Config {
	var o *Config
	if err := json.NewDecoder(data).Decode(&o); err != nil {
		mlog.Error(err.Error())
		return nil
	}
	return o
}

// SetDefaults sets default config settings
func (o *Config) SetDefaults() {
	o.FileSettings.SetDefaults()
	o.ServiceSettings.SetDefaults()
	o.FileSettings.SetDefaults()
	o.LogSettings.SetDefaults()
}

// IsValid check if config is valid
func (o *Config) IsValid() *AppError {
	if err := o.ServiceSettings.isValid(); err != nil {
		return err
	}

	return nil
}

// LogSettings structure
type LogSettings struct {
	EnableConsole          *bool   `restricted:"true"`
	ConsoleLevel           *string `restricted:"true"`
	ConsoleJSON            *bool   `restricted:"true"`
	EnableFile             *bool   `restricted:"true"`
	FileLevel              *string `restricted:"true"`
	FileJSON               *bool   `restricted:"true"`
	FileLocation           *string `restricted:"true"`
	EnableWebhookDebugging *bool   `restricted:"true"`
	EnableDiagnostics      *bool   `restricted:"true"`
}

// SetDefaults sets default log settings
func (s *LogSettings) SetDefaults() {
	if s.EnableConsole == nil {
		s.EnableConsole = NewBool(true)
	}

	if s.ConsoleLevel == nil {
		s.ConsoleLevel = NewString("DEBUG")
	}

	if s.EnableFile == nil {
		s.EnableFile = NewBool(true)
	}

	if s.FileLevel == nil {
		s.FileLevel = NewString("INFO")
	}

	if s.FileLocation == nil {
		s.FileLocation = NewString("")
	}

	if s.EnableWebhookDebugging == nil {
		s.EnableWebhookDebugging = NewBool(true)
	}

	if s.EnableDiagnostics == nil {
		s.EnableDiagnostics = NewBool(true)
	}

	if s.ConsoleJSON == nil {
		s.ConsoleJSON = NewBool(true)
	}

	if s.FileJSON == nil {
		s.FileJSON = NewBool(true)
	}
}

// FileSettings structure
type FileSettings struct {
	MaxFileSize *int64
	DriverName  *string `restricted:"true"`
	Directory   *string `restricted:"true"`
}

// SetDefaults sets default file settings
func (s *FileSettings) SetDefaults() {
	if s.MaxFileSize == nil {
		s.MaxFileSize = NewInt64(52428800) // 50MB
	}

	if s.DriverName == nil {
		s.DriverName = NewString(IMAGE_DRIVER_LOCAL)
	}

	if s.Directory == nil {
		s.Directory = NewString(FILE_SETTINGS_DEFAULT_DIRECTORY)
	}
}
