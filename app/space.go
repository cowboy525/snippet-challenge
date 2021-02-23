package app

import (
	"net/http"

	"github.com/ernie-mlg/ErniePJT-main-api-go/model"
)

// CreateSpace create a new space
func (a *App) CreateSpace(space *model.Space) (*model.Space, *model.AppError) {
	rspace, err := a.Store().Space().Create(space)
	if err != nil {
		err.Where = "CreateSpace"
		err.StatusCode = http.StatusBadRequest
		return nil, err
	}

	return rspace, nil
}
