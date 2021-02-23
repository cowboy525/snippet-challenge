package web

import (
	"net/http"

	"github.com/topoface/snippet-challenge/app"
	"github.com/topoface/snippet-challenge/mlog"
	"github.com/topoface/snippet-challenge/model"
)

// Context structure
type Context struct {
	App           app.Iface
	Log           *mlog.Logger
	Err           *model.AppError
	siteURLHeader string
}

func (c *Context) LogError(err *model.AppError) {
	// Filter out 404s, endless reconnects and browser compatibility errors
	var isDebug bool
	isDebug = false
	if err.StatusCode == http.StatusNotFound {
		isDebug = true
	}
	if detail, ok := c.Err.Errors[0]["detail"]; ok {
		if errArr, ok := detail.([]*model.Error); ok && errArr[0].ID == "web.check_browser_compatibility.app_error" {
			isDebug = true
		}
	}
	if isDebug {
		c.LogDebug(err)
	} else {
		c.Log.Error(
			err.SystemMessage(),
			mlog.String("err_where", err.Where),
			mlog.Int("http_code", err.StatusCode),
			mlog.String("err_details", err.DetailedError),
		)
	}
}

func (c *Context) LogInfo(err *model.AppError) {
	// Filter out 401s
	if err.StatusCode == http.StatusUnauthorized {
		c.LogDebug(err)
	} else {
		c.Log.Info(
			err.SystemMessage(),
			mlog.String("err_where", err.Where),
			mlog.Int("http_code", err.StatusCode),
			mlog.String("err_details", err.DetailedError),
		)
	}
}

func (c *Context) LogDebug(err *model.AppError) {
	c.Log.Debug(
		err.SystemMessage(),
		mlog.String("err_where", err.Where),
		mlog.Int("http_code", err.StatusCode),
		mlog.String("err_details", err.DetailedError),
	)
}
