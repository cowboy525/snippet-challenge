package api

import (
	"net/http"

	"github.com/ernie-mlg/ErniePJT-main-api-go/web"
)

// Context type
type Context = web.Context

// APIHandler provides a handler for API endpoints which do not require the user to be logged in order for access to be
// granted.
func (api *API) APIHandler(h func(*Context, http.ResponseWriter, *http.Request)) http.Handler {
	handler := &web.Handler{
		GetGlobalAppOptions: api.GetGlobalAppOptions,
		HandleFunc:          h,
		HandlerName:         web.GetHandlerName(h),
		RequireSession:      false,
	}
	return handler
}
