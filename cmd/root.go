// Package cmd provides the command-line interface for generating configurations using AI.
//
// Copyright © 2025 PATRICK HERMANN
package cmd

import (
	"os"

	"github.com/spf13/cobra"
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
}
