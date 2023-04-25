package config

import (
	"fmt"
	"github.com/stretchr/testify/assert"

	"testing"
)

func TestGetFilePath(t *testing.T) {
	assert.Equal(t, "etc/config/local.env", GetFilePath("local"))
}

// TestReadConfig reads actual config
func TestReadConfig(t *testing.T) {
	testCases := []struct {
		name, filePath string
		err            error
	}{
		{
			name:     "invalid_file",
			filePath: "etc/config/local.env",
			err:      fmt.Errorf("error loading etc/config/local.env file"),
		},

		{
			name:     "local",
			filePath: "../../../etc/config/local.env",
			err:      nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Read(tc.filePath)

			assert.Equal(t, tc.err, err)

			if tc.err == nil {
				assert.NotZero(t, Global.JwtSecretAccessToken)
			}

			ResetGlobalConfig()
		})
	}
}
