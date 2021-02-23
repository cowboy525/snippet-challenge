package api

import (
	"net/http"
)

// InitAPIRoot : init root api router
func (api *API) InitAPIRoot() {
	api.BaseRoutes.APIRoot.Handle("/", api.APIHandler(healthCheck)).Methods("GET")
}

func healthCheck(c *Context, w http.ResponseWriter, r *http.Request) {
	ReturnStatusOK(w)
}
