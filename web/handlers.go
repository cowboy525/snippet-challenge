package web

import (
	"net/http"
	"reflect"
	"runtime"
	"strings"

	"github.com/topoface/snippet-challenge/app"
	"github.com/topoface/snippet-challenge/mlog"
	"github.com/topoface/snippet-challenge/model"
	"github.com/topoface/snippet-challenge/utils"
)

// GetHandlerName : get handler name from handler func
func GetHandlerName(h func(*Context, http.ResponseWriter, *http.Request)) string {
	handlerName := runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name()
	pos := strings.LastIndex(handlerName, ".")
	if pos != -1 && len(handlerName) > pos {
		handlerName = handlerName[pos+1:]
	}
	return handlerName
}

// Handler : web handler structure
type Handler struct {
	GetGlobalAppOptions app.AppOptionCreator
	HandleFunc          func(*Context, http.ResponseWriter, *http.Request)
	HandlerName         string
	RequireSession      bool
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mlog.Debug("Received HTTP request", mlog.String("method", r.Method), mlog.String("url", r.URL.Path))

	c := &Context{}
	c.App = app.New(
		h.GetGlobalAppOptions()...,
	)

	c.App.SetPath(r.URL.Path)
	c.Params = ParamsFromRequest(r)
	c.Log = c.App.Log()

	subpath, _ := utils.GetSubpathFromConfig(c.App.Config())
	siteURLHeader := app.GetProtocol(r) + "://" + r.Host + subpath
	c.SetSiteURLHeader(siteURLHeader)

	// All api response bodies will be JSON formatted by default
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodGet {
		w.Header().Set("Expires", "0")
	}

	// Authentication
	c.Log = c.App.Log().With(
		mlog.String("path", c.App.Path()),
		mlog.String("method", r.Method),
	)

	// process requests
	if c.Err == nil {
		h.HandleFunc(c, w, r)
	}

	// Handle errors that have occurred
	if c.Err != nil {
		var isInfo bool
		isInfo = false
		if detail, ok := c.Err.Errors[0]["detail"]; ok {
			if errArr, ok := detail.([]*model.Error); ok && errArr[0].ID == "api.context.session_expired" {
				isInfo = true
			}
		}
		if isInfo {
			c.LogInfo(c.Err)
		} else {
			c.LogError(c.Err)
		}

		c.Err.Where = r.URL.Path

		// Block out detailed error when not in developer mode
		if !*c.App.Config().ServiceSettings.EnableDeveloper {
			c.Err.DetailedError = ""
		}

		w.WriteHeader(c.Err.StatusCode)
		w.Write([]byte(c.Err.ToJSON()))
		return
	}
}
