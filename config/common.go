package config

import (
	"bytes"
	"io"
	"strings"
	"sync"

	"github.com/pkg/errors"
	"github.com/topoface/snippet-challenge/model"
)

// commonStore enables code sharing between different backing implementations
type commonStore struct {
	emitter

	configLock             sync.RWMutex
	config                 *model.Config
	configWithoutOverrides *model.Config
	environmentOverrides   map[string]interface{}
}

// Get fetches the current, cached configuration.
func (cs *commonStore) Get() *model.Config {
	cs.configLock.RLock()
	defer cs.configLock.RUnlock()

	return cs.config
}

// GetEnvironmentOverrides fetches the configuration fields overridden by environment variables.
func (cs *commonStore) GetEnvironmentOverrides() map[string]interface{} {
	cs.configLock.RLock()
	defer cs.configLock.RUnlock()

	return cs.environmentOverrides
}

// set replaces the current configuration in its entirety, and updates the backing store
// using the persist function argument.
//
// This function assumes no lock has been acquired, as it acquires a write lock itself.
func (cs *commonStore) set(newCfg *model.Config, allowEnvironmentOverrides bool, validate func(*model.Config) error, persist func(*model.Config) error) (*model.Config, error) {
	cs.configLock.Lock()
	var unlockOnce sync.Once
	defer unlockOnce.Do(cs.configLock.Unlock)

	oldCfg := cs.config

	// To both clone and re-apply the environment variable overrides we marshal and then
	// unmarshal the config again.
	var err error
	newCfg, _, err = unmarshalConfig(strings.NewReader(newCfg.ToJSON()), allowEnvironmentOverrides)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal config with env overrides")
	}

	newCfg.SetDefaults()

	// Sometimes the config is received with "fake" data in sensitive fields. Apply the real
	// data from the existing config as necessary.
	desanitize(oldCfg, newCfg)

	if validate != nil {
		if err := validate(newCfg); err != nil {
			return nil, errors.Wrap(err, "new configuration is invalid")
		}
	}

	if err := persist(cs.RemoveEnvironmentOverrides(newCfg)); err != nil {
		return nil, errors.Wrap(err, "failed to persist")
	}

	cs.config = newCfg

	unlockOnce.Do(cs.configLock.Unlock)

	// Notify listeners synchronously. Ideally, this would be asynchronous, but existing code
	// assumes this and there would be increased complexity to avoid racing updates.
	cs.invokeConfigListeners(oldCfg, newCfg)

	return oldCfg, nil
}

// load updates the current configuration from the given io.ReadCloser.
//
// This function assumes no lock has been acquired, as it acquires a write lock itself.
func (cs *commonStore) load(f io.ReadCloser, needsSave bool, validate func(*model.Config) error, persist func(*model.Config) error) error {
	// Duplicate f so that we can read a configuration without applying environment overrides
	f2 := new(bytes.Buffer)
	tee := io.TeeReader(f, f2)

	allowEnvironmentOverrides := true
	loadedCfg, environmentOverrides, err := unmarshalConfig(tee, allowEnvironmentOverrides)
	if err != nil {
		return errors.Wrapf(err, "failed to unmarshal config with env overrides")
	}

	// Keep track of the original values that the Environment settings overrode
	loadedCfgWithoutEnvOverrides, _, err := unmarshalConfig(f2, false)
	if err != nil {
		return errors.Wrapf(err, "failed to unmarshal config without env overrides")
	}

	loadedCfg.SetDefaults()
	loadedCfgWithoutEnvOverrides.SetDefaults()

	if validate != nil {
		if err = validate(loadedCfg); err != nil {
			return errors.Wrap(err, "invalid config")
		}
	}

	if changed := fixConfig(loadedCfg); changed {
		needsSave = true
	}

	cs.configLock.Lock()
	var unlockOnce sync.Once
	defer unlockOnce.Do(cs.configLock.Unlock)

	if needsSave && persist != nil {
		cfgWithoutEnvOverrides := removeEnvOverrides(loadedCfg, loadedCfgWithoutEnvOverrides, environmentOverrides)
		if err = persist(cfgWithoutEnvOverrides); err != nil {
			return errors.Wrap(err, "failed to persist required changes after load")
		}
	}

	cs.config = loadedCfg
	cs.configWithoutOverrides = loadedCfgWithoutEnvOverrides
	cs.environmentOverrides = environmentOverrides

	unlockOnce.Do(cs.configLock.Unlock)

	return nil
}

// validate checks if the given configuration is valid
func (cs *commonStore) validate(cfg *model.Config) error {
	if err := cfg.IsValid(); err != nil {
		return err
	}

	return nil
}

// RemoveEnvironmentOverrides returns a new config without the given environment overrides.
func (cs *commonStore) RemoveEnvironmentOverrides(cfg *model.Config) *model.Config {
	return removeEnvOverrides(cfg, cs.configWithoutOverrides, cs.environmentOverrides)
}
