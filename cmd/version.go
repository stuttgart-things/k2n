// Package cmd provides the command-line interface for generating configurations using AI.
//
// Copyright © 2025 PATRICK HERMANN
package cmd

import (
	"fmt"

	goVersion "go.hein.dev/go-version"

	"github.com/spf13/cobra"
	"github.com/stuttgart-things/k2n/internal"
)

var (
	shortened  = false
	version    = "unset"
	commit     = "unknown"
	date       = "unknown"
	output     = "yaml"
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "version will output the current build information",
		Long:  "version will output the current build information",

		Run: func(_ *cobra.Command, _ []string) {

			internal.PrintBanner()

			resp := goVersion.FuncWithOutput(shortened, version, commit, date, output)
			fmt.Print(resp)
		},
	}
)

func init() {
	versionCmd.Flags().BoolVarP(&shortened, "short", "s", false, "Print just the version number.")
	versionCmd.Flags().StringVarP(&output, "output", "o", "yaml", "Output format. One of 'yaml' or 'json'.")
	rootCmd.AddCommand(versionCmd)
}
