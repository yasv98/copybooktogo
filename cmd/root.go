package cmd

import (
	"github.com/carlmjohnson/versioninfo"
	"github.com/spf13/cobra"

	"copybooktogo/copybooktogo"
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
	copybookPath string
	packageName  string
)

// Execute runs the root command.
func Execute() error {
	rootCmd.Version = versioninfo.Version
	return rootCmd.Execute()
}

func init() {
	rootCmd.Flags().StringVarP(&copybookPath, "copybook", "c", "", "Path to the copybook file (required)")
	rootCmd.Flags().StringVarP(&packageName, "package", "p", "main", "Package name for generated Go code")

	_ = rootCmd.MarkFlagRequired("copybook")
}

func run(_ *cobra.Command, _ []string) error {
	cfg, err := copybooktogo.NewConfig(copybookPath, packageName)
	if err != nil {
		return err
	}
	return copybooktogo.Process(cfg)
}
