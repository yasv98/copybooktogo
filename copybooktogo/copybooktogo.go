// Package copybooktogo provides functionality for processing COBOL copybooks and generating equivalent Go struct
// definitions.
package copybooktogo

import (
	"fmt"
	"go/token"
	"os"
	"path"
	"path/filepath"
	"strings"

	"copybooktogo/generate"
	"copybooktogo/normalise"
	"copybooktogo/parse"
)

// Process reads a COBOL copybook file and generates Go struct definitions.
func Process(cfg *Config) error {
	copybook, err := os.ReadFile(cfg.CopybookPath)
	if err != nil {
		return fmt.Errorf("reading copybook file: %w", err)
	}

	normalisedCopybook, err := normalise.Format(copybook)
	if err != nil {
		return fmt.Errorf("normalizing copybook: %w", err)
	}

	ast, err := parse.BuildAST(normalisedCopybook)
	if err != nil {
		return fmt.Errorf("parsing copybook: %w", err)
	}

	data, err := generate.ToGoStructsData(ast, getCopybookName(cfg.CopybookPath), cfg.PackageName, cfg.TypeOverrides)
	if err != nil {
		return fmt.Errorf("generating Go structs: %w", err)
	}

	if err := os.WriteFile(cfg.OutputPath, data, 0o600); err != nil {
		return fmt.Errorf("writing output file: %w", err)
	}

	fmt.Printf("Successfully generated Go structs in: %s\n", cfg.OutputPath)
	return nil
}

// Config holds the configuration for the copybooktogo tool.
type Config struct {
	CopybookPath  string
	PackageName   string
	TypeOverrides map[parse.PicType]string
	OutputPath    string
}

// NewConfig creates new Config and validates it.
func NewConfig(copybookPath, packageName, outputPath string, typeOverrides map[string]string) (*Config, error) {
	if _, err := os.Stat(copybookPath); err != nil {
		return nil, fmt.Errorf("copybook file path error: %w", err)
	}

	if !token.IsIdentifier(packageName) {
		return nil, fmt.Errorf("package name %q is not a valid Go identifier", packageName)
	}

	overrides, err := convertTypeOverrides(typeOverrides)
	if err != nil {
		return nil, err
	}

	outputPath = determineOutputPath(outputPath, copybookPath)

	cfg := &Config{
		CopybookPath:  copybookPath,
		PackageName:   packageName,
		TypeOverrides: overrides,
		OutputPath:    outputPath,
	}
	return cfg, nil
}

func convertTypeOverrides(typeOverrides map[string]string) (map[parse.PicType]string, error) {
	overrides := make(map[parse.PicType]string)
	for fromPicType, toGoType := range typeOverrides {
		picType, err := parse.PicTypeFromString(fromPicType)
		if err != nil {
			return nil, err
		}
		overrides[picType] = toGoType
	}
	return overrides, nil
}

func getCopybookName(filePath string) string {
	fileNameWithExtension := filepath.Base(filePath)
	return strings.TrimSuffix(fileNameWithExtension, filepath.Ext(fileNameWithExtension))
}

func createGoFileName(filePath string) string {
	return filepath.Join(path.Dir(filePath), strings.ToLower(getCopybookName(filePath))+".generated.go")
}

func determineOutputPath(outputPath, copybookPath string) string {
	if outputPath == "" {
		return createGoFileName(copybookPath)
	}
	if filepath.Ext(outputPath) == ".go" {
		return outputPath
	}

	return filepath.Join(outputPath, strings.ToLower(getCopybookName(copybookPath))+".generated.go")
}
