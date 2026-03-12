/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	defaultBaseURLTemplate = "https://%s.clients.us-1.kandji.io"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kandji-iru-cli",
	Short: "CLI for the Kandji API",
	Long: `A CLI to manage Kandji devices and other resources via the Kandji API.
Requires KANDJI_TOKEN and KANDJI_BASE_URL (or KANDJI_SUBDOMAIN).`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return validateConfig()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.kandji.yaml)")
	rootCmd.PersistentFlags().String("base-url", "", "Kandji API base URL (e.g. https://your-tenant.clients.us-1.kandji.io)")
	rootCmd.PersistentFlags().String("subdomain", "", "Kandji tenant subdomain (used to build base URL if --base-url is not set)")
	rootCmd.PersistentFlags().String("token", "", "Kandji API bearer token")
	rootCmd.PersistentFlags().StringP("output", "o", "table", "Output format: table, json, or raw (exact API response bytes)")
	rootCmd.PersistentFlags().Bool("raw", false, "Emit raw API response body (same as -o raw)")
	rootCmd.PersistentFlags().Bool("verbose", false, "Verbose output")

	_ = viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))
	_ = viper.BindPFlag("base_url", rootCmd.PersistentFlags().Lookup("base-url"))
	_ = viper.BindPFlag("subdomain", rootCmd.PersistentFlags().Lookup("subdomain"))
	_ = viper.BindPFlag("token", rootCmd.PersistentFlags().Lookup("token"))
	_ = viper.BindPFlag("output", rootCmd.PersistentFlags().Lookup("output"))
	_ = viper.BindPFlag("raw", rootCmd.PersistentFlags().Lookup("raw"))
	_ = viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName(".kandji")
	}

	viper.SetEnvPrefix("KANDJI")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}

func validateConfig() error {
	// Skip validation for init so the config file can be created without credentials.
	if len(os.Args) >= 2 && os.Args[1] == "init" {
		return nil
	}
	token := viper.GetString("token")
	if token == "" {
		return fmt.Errorf("API token not configured\nSet KANDJI_TOKEN or use --token or add to config file (~/.kandji.yaml)")
	}

	baseURL := viper.GetString("base_url")
	if baseURL == "" {
		subdomain := viper.GetString("subdomain")
		if subdomain == "" {
			return fmt.Errorf("API base URL not configured\nSet KANDJI_BASE_URL or KANDJI_SUBDOMAIN, or use --base-url or --subdomain")
		}
		// Subdomain is the tenant name, e.g. "acme" -> https://acme.clients.us-1.kandji.io
		baseURL = fmt.Sprintf(defaultBaseURLTemplate, strings.TrimSpace(subdomain))
	}
	viper.Set("resolved_base_url", strings.TrimSuffix(baseURL, "/"))
	return nil
}
