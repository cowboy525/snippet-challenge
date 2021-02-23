package web

import (
	"net/http"
	"strings"
)

// InitStatic : serve static files
func (w *Web) InitStatic() {
	staticDir := *w.ConfigService.Config().FileSettings.Directory
	staticDir = strings.TrimPrefix(staticDir, ".")
	w.MainRouter.
		PathPrefix(staticDir).
		Handler(http.StripPrefix(staticDir, http.FileServer(http.Dir("."+staticDir))))
}
