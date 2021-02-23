package model

import (
	"encoding/json"
	"io"
	"time"

	"github.com/ernie-mlg/ErniePJT-main-api-go/mlog"
)

// Snippet structure
type Snippet struct {
	URL       string    `json:"url" validate:"blank:false;required"`
	Name      string    `json:"name"`
	ExpiresAt time.Time `json:"expires_at"`
	Body      string    `json:"snippet"`
}

func SnippetFromJSON(data io.Reader) *Snippet {
	var o *Snippet
	if err := json.NewDecoder(data).Decode(&o); err != nil {
		mlog.Error(err.Error())
		return nil
	}
	return o
}

func (o *Snippet) ToJSON() string {
	b, _ := json.Marshal(o)
	return string(b)
}

// SnippetRequest structure
type SnippetRequest struct {
	Name      string `json:"name"`
	ExpiresIn uint64 `json:"expires_in"`
	Body      string `json:"snippet"`
}
