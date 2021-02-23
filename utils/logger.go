package utils

import (
	"path/filepath"
	"strings"

	"github.com/ernie-mlg/ErniePJT-main-api-go/mlog"
	"github.com/ernie-mlg/ErniePJT-main-api-go/model"
	"github.com/ernie-mlg/ErniePJT-main-api-go/utils/fileutils"
)

const (
	LOG_ROTATE_SIZE = 10000
	LOG_FILENAME    = "erniepjt.log"
)

type fileLocationFunc func(string) string

func MloggerConfigFromLoggerConfig(s *model.LogSettings, getFileFunc fileLocationFunc) *mlog.LoggerConfiguration {
	return &mlog.LoggerConfiguration{
		EnableConsole: *s.EnableConsole,
		ConsoleJSON:   *s.ConsoleJSON,
		ConsoleLevel:  strings.ToLower(*s.ConsoleLevel),
		EnableFile:    *s.EnableFile,
		FileJSON:      *s.FileJSON,
		FileLevel:     strings.ToLower(*s.FileLevel),
		FileLocation:  getFileFunc(*s.FileLocation),
	}
}

func GetLogFileLocation(fileLocation string) string {
	if fileLocation == "" {
		fileLocation, _ = fileutils.FindDir("logs")
	}

	return filepath.Join(fileLocation, LOG_FILENAME)
}

// DON'T USE THIS Modify the level on the app logger
func DisableDebugLogForTest() {
	mlog.GloballyDisableDebugLogForTest()
}

// DON'T USE THIS Modify the level on the app logger
func EnableDebugLogForTest() {
	mlog.GloballyEnableDebugLogForTest()
}
