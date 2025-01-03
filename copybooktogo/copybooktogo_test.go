package copybooktogo

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"

	"github.com/yasv98/copybooktogo/parse"
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
		typeOverrides  map[string]string
		expectedConfig *Config
		assertError    assert.ErrorAssertionFunc
	}{
		"ValidConfig_ReturnsConfig": {
			copybookPath:  tmpFile.Name(),
			packageName:   "validpackage",
			typeOverrides: nil,
			expectedConfig: &Config{
				CopybookPath:  tmpFile.Name(),
				PackageName:   "validpackage",
				TypeOverrides: map[parse.PicType]string{},
			},
			assertError: assert.NoError,
		},
		"ValidConfigWithOverrides_ReturnsConfigWithParsedOverrides": {
			copybookPath: tmpFile.Name(),
			packageName:  "validpackage",
			typeOverrides: map[string]string{
				"unsigned": "int",
				"decimal":  "string",
			},
			expectedConfig: &Config{
				CopybookPath: tmpFile.Name(),
				PackageName:  "validpackage",
				TypeOverrides: map[parse.PicType]string{
					parse.Unsigned: "int",
					parse.Decimal:  "string",
				},
			},
			assertError: assert.NoError,
		},
		"InvalidFilePath_ReturnsError": {
			copybookPath:   "/nonexistent/path",
			packageName:    "validpackage",
			typeOverrides:  nil,
			expectedConfig: nil,
			assertError:    assert.Error,
		},
		"InvalidPackageName_ReturnsError": {
			copybookPath:   tmpFile.Name(),
			packageName:    "invalid-package",
			typeOverrides:  nil,
			expectedConfig: nil,
			assertError:    assert.Error,
		},
		"InvalidTypeOverrides_ReturnsError": {
			copybookPath: tmpFile.Name(),
			packageName:  "validpackage",
			typeOverrides: map[string]string{
				"invalid": "int",
			},
			expectedConfig: nil,
			assertError:    assert.Error,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			cfg, err := NewConfig(tt.copybookPath, tt.packageName, "", tt.typeOverrides)
			tt.assertError(t, err)
			assert.True(t, cmp.Equal(tt.expectedConfig, cfg, cmpopts.IgnoreFields(Config{}, "OutputPath")))
		})
	}
}

func Test_determineOutputPath(t *testing.T) {
	tests := map[string]struct {
		outputPath         string
		copybookPath       string
		expectedOutputPath string
	}{
		"EmptyOutputPath_ReturnsGeneratedGoFileName": {
			outputPath:         "",
			copybookPath:       "/path/to/copybook.cpy",
			expectedOutputPath: "/path/to/copybook.generated.go",
		},
		"OutputPathWithGoExtension_ReturnsOutputPath": {
			outputPath:         "/different/path/to/output.go",
			copybookPath:       "/path/to/copybook.cpy",
			expectedOutputPath: "/different/path/to/output.go",
		},
		"OutputPathToDirectory_ReturnsOutputPath": {
			outputPath:         "/different/path/to/output",
			copybookPath:       "/path/to/copybook.cpy",
			expectedOutputPath: "/different/path/to/output/copybook.generated.go",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tt.expectedOutputPath, determineOutputPath(tt.outputPath, tt.copybookPath))
		})
	}
}
