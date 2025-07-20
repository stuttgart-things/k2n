/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/stuttgart-things/k2n/internal"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "k2n",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
	internal.PrintBanner()

	envData := map[string]string{
		"GIT-REPO":        "",
		"VAULT_ADDR":      os.Getenv("VAULT_ADDR"),
		"VAULT_NAMESPACE": os.Getenv("VAULT_NAMESPACE"),
		"VAULT_ROLE_ID":   os.Getenv("VAULT_ROLE_ID"),
		"VAULT_SECRET_ID": os.Getenv("VAULT_SECRET_ID"),
		"VAULT_TOKEN":     os.Getenv("VAULT_TOKEN"),
	}

	internal.PrintEnvTable(envData)

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
