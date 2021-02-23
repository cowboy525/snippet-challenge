package web

import (
	"net/http"
	"reflect"
	"runtime"
	"strings"

	"github.com/ernie-mlg/ErniePJT-main-api-go/app"
	"github.com/ernie-mlg/ErniePJT-main-api-go/mlog"
	"github.com/ernie-mlg/ErniePJT-main-api-go/model"
	"github.com/ernie-mlg/ErniePJT-main-api-go/utils"
	"github.com/ernie-mlg/ErniePJT-main-api-go/web/middleware"
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
	clientID := r.Header.Get(model.HEADER_CLIENT_ID)
	mlog.Debug("Received HTTP request", mlog.String("method", r.Method), mlog.String("url", r.URL.Path), mlog.String("client_id", clientID))

	c := &Context{}
	c.App = app.New(
		h.GetGlobalAppOptions()...,
	)

	if r.Method != http.MethodGet {
		c.App.Store().LockToMaster()
	}

	t, locale := utils.GetTranslationsAndLocale(w, r)
	c.App.SetT(t)
	c.App.SetLocale(locale)
	c.App.SetAcceptLanguage(r.Header.Get("Accept-Language"))
	c.App.SetClientID(clientID)
	c.App.SetUserAgent(r.UserAgent())
	c.App.SetClientSecretKey(r.Header.Get(model.HEADER_CLIENT_SECRETKEY))
	c.App.SetPath(r.URL.Path)
	c.Params = ParamsFromRequest(r)
	c.Log = c.App.Log()

	subpath, _ := utils.GetSubpathFromConfig(c.App.Config())
	siteURLHeader := app.GetProtocol(r) + "://" + r.Host + subpath
	c.SetSiteURLHeader(siteURLHeader)

	w.Header().Set(model.HEADER_CLIENT_ID, c.App.ClientID())

	// All api response bodies will be JSON formatted by default
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodGet {
		w.Header().Set("Expires", "0")
	}

	// Authentication
	if c.Err == nil && h.RequireSession {
		token, tokenLocation := app.ParseAuthTokenFromRequest(r)

		if len(clientID) == 0 {
			c.Err = model.AuthenticationFailedCustomError("ServeHTTP", "model.app_error.no_client_id", nil, "Client Id is empty")
		} else if len(token) == 0 {
			c.Err = model.AuthenticationFailedCustomError("ServeHTTP", "model.app_error.no_credentials", nil, "Token is empty")
		} else {
			session, err := c.App.GetSession(token)
			if err != nil {
				c.Err = err
				if err.StatusCode != http.StatusInternalServerError {
					c.RemoveSessionCookie(w, r)
				}
			} else if tokenLocation == app.TokenLocationQueryString {
				c.Err = model.AuthenticationFailedCustomError("ServeHTTP", "model.app_error.query_token", nil, "token="+token)
			} else {
				c.App.SetSession(session)
			}
		}
	}

	c.Log = c.App.Log().With(
		mlog.String("path", c.App.Path()),
		mlog.String("client_id", c.App.ClientID()),
		mlog.Uint64("user_id", c.App.Session().UserID),
		mlog.String("method", r.Method),
	)

	// process requests
	if c.Err == nil {
		if *c.App.Config().ServiceSettings.AtomicRequest {
			c.App.Store().Begin()
		}

		h.HandleFunc(c, w, r)

		if *c.App.Config().ServiceSettings.AtomicRequest {
			if c.Err == nil {
				c.App.Store().Commit()
			} else {
				c.App.Store().Rollback()
			}
		}
	}

	if r.MultipartForm != nil {
		r.MultipartForm.RemoveAll()
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

		c.Err.Translate(c.App.T)
		c.Err.ClientID = c.App.ClientID()
		c.Err.Where = r.URL.Path

		// Block out detailed error when not in developer mode
		if !*c.App.Config().ServiceSettings.EnableDeveloper {
			c.Err.DetailedError = ""
		}

		w.WriteHeader(c.Err.StatusCode)
		w.Write([]byte(c.Err.ToJSON()))
		return
	}

	middleware.Notification(c.App)
	if !c.App.Srv().IsTest() {
		middleware.SlackNotification(c.App)
		middleware.Synchronization(c.App)
		middleware.SendPushNotificaton(c.App)
	}
}
