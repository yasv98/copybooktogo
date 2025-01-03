package main

import (
	"os"

	"github.com/yasv98/copybooktogo/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
