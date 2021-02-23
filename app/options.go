package app

import (
	"github.com/ernie-mlg/ErniePJT-main-api-go/config"
	"github.com/ernie-mlg/ErniePJT-main-api-go/mlog"
	"github.com/pkg/errors"
)

type Option func(s *Server) error

// Config applies the given config dsn, whether a path to config.json or a database connection string.
func Config(dsn string, watch bool) Option {
	return func(s *Server) error {
		configStore, err := config.NewStore(dsn, watch)
		if err != nil {
			return errors.Wrap(err, "failed to apply Config option")
		}

		s.configStore = configStore
		return nil
	}
}

// ConfigStore applies the given config store, typically to replace the traditional sources with a memory store for testing.
func ConfigStore(configStore config.Store) Option {
	return func(s *Server) error {
		s.configStore = configStore

		return nil
	}
}

func SetLogger(logger *mlog.Logger) Option {
	return func(s *Server) error {
		s.Log = logger
		return nil
	}
}

type AppOption func(a *App)
type AppOptionCreator func() []AppOption

func ServerConnector(s *Server) AppOption {
	return func(a *App) {
		a.srv = s
		if *s.Config().ServiceSettings.AtomicRequest {
			a.store = s.Store.CopyStore()
		}
	}
}
