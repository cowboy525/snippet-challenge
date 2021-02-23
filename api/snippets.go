package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/topoface/snippet-challenge/binding"
	"github.com/topoface/snippet-challenge/model"
)

func (api *API) InitSnippets() {
	api.BaseRoutes.Snippets.Handle("", api.APIHandler(createSnippet)).Methods("POST")
	api.BaseRoutes.Snippets.Handle("/{name}", api.APIHandler(getSnippet)).Methods("GET")
}

func createSnippet(c *Context, w http.ResponseWriter, r *http.Request) {
	var snippetRequest model.SnippetRequest
	if _, err := binding.JSON.Bind(r, &snippetRequest); err != nil {
		c.Err = err
		return
	}

	snippet, err := c.App.CreateSnippet(&snippetRequest)
	if err != nil {
		c.Err = err
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(snippet.ToJSON()))
}

func getSnippet(c *Context, w http.ResponseWriter, r *http.Request) {
	props := mux.Vars(r)
	snippetName, ok := props["name"]
	if !ok {
		c.Err = model.ValidationError("", "No snippet name", nil, "")
		return
	}

	snippet, err := c.App.GetSnippet(snippetName)
	if err != nil {
		c.Err = err
		return
	}

	w.Write([]byte(snippet.ToJSON()))
}
