package configservice

import (
	"github.com/ernie-mlg/ErniePJT-main-api-go/model"
)

// An interface representing something that contains a Config, such as the app.App struct
type ConfigService interface {
	Config() *model.Config
	AddConfigListener(func(old, current *model.Config)) string
	RemoveConfigListener(string)
}
