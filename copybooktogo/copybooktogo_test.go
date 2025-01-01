package copybooktogo

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	// Create a temporary file for testing
	tmpFile, err := os.CreateTemp("", "copybook-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	tests := map[string]struct {
		copybookPath   string
		packageName    string
		expectedConfig *Config
		assertError    assert.ErrorAssertionFunc
	}{
		"ValidConfig_ReturnsConfig": {
			copybookPath: tmpFile.Name(),
			packageName:  "validpackage",
			expectedConfig: &Config{
				CopybookPath: tmpFile.Name(),
				PackageName:  "validpackage",
			},
			assertError: assert.NoError,
		},
		"InvalidFilePath_ReturnsError": {
			copybookPath:   "/nonexistent/path",
			packageName:    "validpackage",
			expectedConfig: nil,
			assertError:    assert.Error,
		},
		"InvalidPackageName_ReturnsError": {
			copybookPath:   tmpFile.Name(),
			packageName:    "invalid-package",
			expectedConfig: nil,
			assertError:    assert.Error,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			cfg, err := NewConfig(tt.copybookPath, tt.packageName)
			tt.assertError(t, err)
			assert.Equal(t, tt.expectedConfig, cfg)
		})
	}
}
