package utils_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ernie-mlg/ErniePJT-main-api-go/model"
	"github.com/ernie-mlg/ErniePJT-main-api-go/utils"
)

func TestGetSubpathFromConfig(t *testing.T) {
	testCases := []struct {
		Description     string
		SiteURL         *string
		ExpectedError   bool
		ExpectedSubpath string
	}{
		{
			"empty SiteURL",
			sToP(""),
			false,
			"/",
		},
		{
			"invalid SiteURL",
			sToP("cache_object:foo/bar"),
			true,
			"",
		},
		{
			"nil SiteURL",
			nil,
			false,
			"/",
		},
		{
			"no trailing slash",
			sToP("http://localhost:8065"),
			false,
			"/",
		},
		{
			"trailing slash",
			sToP("http://localhost:8065/"),
			false,
			"/",
		},
		{
			"subpath, no trailing slash",
			sToP("http://localhost:8065/subpath"),
			false,
			"/subpath",
		},
		{
			"trailing slash",
			sToP("http://localhost:8065/subpath/"),
			false,
			"/subpath",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Description, func(t *testing.T) {
			config := &model.Config{
				ServiceSettings: model.ServiceSettings{
					SiteURL: testCase.SiteURL,
				},
			}

			subpath, err := utils.GetSubpathFromConfig(config)
			if testCase.ExpectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, testCase.ExpectedSubpath, subpath)
		})
	}
}

func sToP(s string) *string {
	return &s
}
