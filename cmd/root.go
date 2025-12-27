// Package cmd provides the command-line interface for generating configurations using AI.
//
// Copyright © 2025 PATRICK HERMANN
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/stuttgart-things/k2n/internal/menu"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "k2n",
	Short: "Generate Kubernetes and infrastructure configurations using AI",
	Long: `k2n is a powerful CLI tool that leverages AI to intelligently generate
Kubernetes and infrastructure configurations. It transforms natural language
descriptions into properly formatted configuration files, reducing manual
effort and ensuring consistency across your infrastructure.

Use k2n to generate Kubernetes manifests, Helm values, Crossplane compositions,
KCL modules, and other infrastructure-as-code artifacts with ease.`,
	Run: func(cmd *cobra.Command, args []string) {
		// If no subcommand or arguments are provided, launch the interactive menu
		if len(args) == 0 && !cmd.Flags().Changed("toggle") {
			if err := menu.ShowInteractiveMenu(cmd); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
			return
		}

		// Otherwise, show help
		cmd.Help()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.k2n.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// Disable default completion command for cleaner interface
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}
