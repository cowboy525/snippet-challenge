package api

import (
	"net/http"
)

func (api *API) InitSnippets() {
	// api.BaseRoutes.Snippets.Handle("/confirm/", api.APIHandler(confirmSpace)).Methods("POST")
}

func confirmSpace(c *Context, w http.ResponseWriter, r *http.Request) {
}
