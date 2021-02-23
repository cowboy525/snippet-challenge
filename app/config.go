package app

import (
	"net/http"
	"runtime/debug"

	"github.com/pkg/errors"

	"github.com/topoface/snippet-challenge/mlog"
	"github.com/topoface/snippet-challenge/model"
)

func (s *Server) Config() *model.Config {
	return s.configStore.Get()
}

func (a *App) Config() *model.Config {
	return a.Srv().Config()
}

func (s *Server) EnvironmentConfig() map[string]interface{} {
	return s.configStore.GetEnvironmentOverrides()
}

func (a *App) EnvironmentConfig() map[string]interface{} {
	return a.Srv().EnvironmentConfig()
}

func (s *Server) UpdateConfig(f func(*model.Config)) {
	old := s.Config()
	updated := old.Clone()
	f(updated)
	if _, err := s.configStore.Set(updated); err != nil {
		mlog.Error("Failed to update config", mlog.Err(err))
	}
}

func (a *App) UpdateConfig(f func(*model.Config)) {
	a.Srv().UpdateConfig(f)
}

func (s *Server) ReloadConfig() error {
	debug.FreeOSMemory()
	if err := s.configStore.Load(); err != nil {
		return err
	}
	return nil
}

func (a *App) ReloadConfig() error {
	return a.Srv().ReloadConfig()
}

// Registers a function with a given listener to be called when the config is reloaded and may have changed. The function
// will be called with two arguments: the old config and the new config. AddConfigListener returns a unique ID
// for the listener that can later be used to remove it.
func (s *Server) AddConfigListener(listener func(*model.Config, *model.Config)) string {
	return s.configStore.AddListener(listener)
}

func (a *App) AddConfigListener(listener func(*model.Config, *model.Config)) string {
	return a.Srv().AddConfigListener(listener)
}

// Removes a listener function by the unique ID returned when AddConfigListener was called
func (s *Server) RemoveConfigListener(id string) {
	s.configStore.RemoveListener(id)
}

func (a *App) RemoveConfigListener(id string) {
	a.Srv().RemoveConfigListener(id)
}

func (a *App) GetSiteURL() string {
	return *a.Config().ServiceSettings.SiteURL
}

// GetConfigFile proxies access to the given configuration file to the underlying config store.
func (a *App) GetConfigFile(name string) ([]byte, error) {
	data, err := a.Srv().configStore.GetFile(name)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get config file %s", name)
	}

	return data, nil
}

// GetSanitizedConfig gets the configuration for a system admin without any secrets.
func (a *App) GetSanitizedConfig() *model.Config {
	cfg := a.Config().Clone()

	return cfg
}

// GetEnvironmentConfig returns a map of configuration keys whose values have been overridden by an environment variable.
func (a *App) GetEnvironmentConfig() map[string]interface{} {
	return a.EnvironmentConfig()
}

// SaveConfig replaces the active configuration, optionally notifying cluster peers.
func (a *App) SaveConfig(newCfg *model.Config) *model.AppError {
	_, err := a.Srv().configStore.Set(newCfg)
	if err != nil {
		return model.NewAppError("saveConfig", "app.save_config.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return nil
}
