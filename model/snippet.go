package model

import (
	"encoding/json"
	"io"
	"time"

	"github.com/topoface/snippet-challenge/mlog"
)

// Snippet structure
type Snippet struct {
	URL       string    `json:"url"`
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
	Name      string `json:"name" validate:"blank:false;required"`
	ExpiresIn uint64 `json:"expires_in"`
	Body      string `json:"snippet" validate:"blank:false;required"`
}
