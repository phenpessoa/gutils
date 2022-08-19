package zaputils

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCustomLoggers(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Log("skipping TestNewCustomLoggers: this test can only be run on unix systems")
		return
	}

	_, _, _, err := NewCustomLoggers("test")
	require.Nil(t, err)

	for _, tc := range []struct {
		testName   string
		prettyName string
		mainName   string
		sysName    string
		shouldErr  bool
	}{
		{"valid", "test", "test", "test", false},
		{"invalid_pretty", "", "test", "test", true},
		{"invalid_main", "test", "", "test", true},
		{"invalid_sys", "test", "test", "", true},
	} {
		t.Run(tc.testName, func(t *testing.T) {
			_, _, _, err := newCustomLoggers(tc.prettyName, tc.mainName, tc.sysName)
			if tc.shouldErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}

}
