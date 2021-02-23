package binding

import (
	"net/http"

	"github.com/topoface/snippet-challenge/model"
	"github.com/topoface/snippet-challenge/validator"
)

type Binding interface {
	Name() string
	Bind(*http.Request, interface{}) (map[string]interface{}, *model.AppError)
}

// These implement the Binding interface and can be used to bind the data
// present in the request to struct instances.
var (
	JSON = jsonBinding{}
)

func Default() Binding {
	return JSON
}

func validate(obj interface{}, data interface{}) *model.AppError {
	return validator.Validate(obj, data)
}
