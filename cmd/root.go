package cmd

import (
	"github.com/carlmjohnson/versioninfo"
	"github.com/spf13/cobra"

	"github.com/yasv98/copybooktogo/copybooktogo"
)

var (
	rootCmd = &cobra.Command{
		Use:   "copybooktogo",
		Short: "Convert COBOL copybooks to Go structs",
		Long: `copybooktogo is a CLI tool that converts COBOL copybooks to Go struct definitions.
It handles the normalization, parsing and generation of equivalent Go code.`,
		RunE:         run,
		SilenceUsage: true,
	}

	// Flag variables.
	copybookPath  string
	packageName   string
	typeOverrides map[string]string
	outputPath    string
)

// Execute runs the root command.
func Execute() error {
	rootCmd.Version = versioninfo.Version
	return rootCmd.Execute()
}

func init() {
	rootCmd.Flags().StringVarP(&copybookPath, "copybook", "c", "", "Path to the copybook file (required)")
	rootCmd.Flags().StringVarP(&packageName, "package", "p", "main", "Package name for generated Go code")
	rootCmd.Flags().StringToStringVarP(&typeOverrides, "typeOverrides", "t", nil,
		"Custom overrides that map PIC types to configured Go types in from=to format (e.g., unsigned=int,decimal=custom.Type)")
	rootCmd.Flags().StringVarP(&outputPath, "output", "o", "", "Path to the output file or directory")

	_ = rootCmd.MarkFlagRequired("copybook")
}

func run(_ *cobra.Command, _ []string) error {
	cfg, err := copybooktogo.NewConfig(copybookPath, packageName, outputPath, typeOverrides)
	if err != nil {
		return err
	}
	return copybooktogo.Process(cfg)
}
