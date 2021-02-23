package config

import (
	"strings"

	"github.com/topoface/snippet-challenge/model"
	"github.com/topoface/snippet-challenge/utils"
)

// desanitize replaces fake settings with their actual values.
func desanitize(actual, target *model.Config) {

	if *target.EmailSettings.SMTPPassword == model.FAKE_SETTING {
		target.EmailSettings.SMTPPassword = actual.EmailSettings.SMTPPassword
	}

	if *target.SQLSettings.DataSource == model.FAKE_SETTING {
		*target.SQLSettings.DataSource = *actual.SQLSettings.DataSource
	}

	target.SQLSettings.DataSourceReplicas = make([]string, len(actual.SQLSettings.DataSourceReplicas))
	for i := range target.SQLSettings.DataSourceReplicas {
		target.SQLSettings.DataSourceReplicas[i] = actual.SQLSettings.DataSourceReplicas[i]
	}
}

// fixConfig patches invalid or missing data in the configuration, returning true if changed.
func fixConfig(cfg *model.Config) bool {
	changed := false

	// Ensure SiteURL has no trailing slash.
	if strings.HasSuffix(*cfg.ServiceSettings.SiteURL, "/") {
		*cfg.ServiceSettings.SiteURL = strings.TrimRight(*cfg.ServiceSettings.SiteURL, "/")
		changed = true
	}

	return changed
}

// Merge merges two configs together. The receiver's values are overwritten with the patch's
// values except when the patch's values are nil.
func Merge(cfg *model.Config, patch *model.Config, mergeConfig *utils.MergeConfig) (*model.Config, error) {
	ret, err := utils.Merge(cfg, patch, mergeConfig)
	if err != nil {
		return nil, err
	}

	retCfg := ret.(model.Config)
	return &retCfg, nil
}

// stripPassword remove the password from a given DSN
// func stripPassword(dsn, schema string) string {
// 	prefix := schema + "://"
// 	dsn = strings.TrimPrefix(dsn, prefix)

// 	i := strings.Index(dsn, ":")
// 	j := strings.LastIndex(dsn, "@")

// 	// Return error if no @ sign is found
// 	if j < 0 {
// 		return "(omitted due to error parsing the DSN)"
// 	}

// 	// Return back the input if no password is found
// 	if i < 0 || i > j {
// 		return prefix + dsn
// 	}

// 	return prefix + dsn[:i+1] + dsn[j:]
// }
