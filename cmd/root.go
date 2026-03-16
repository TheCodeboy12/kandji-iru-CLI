/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"

	"kandji-iru-cli/internal/keyring"
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
		return validateConfig(cmd)
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

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/kandji-iru-cli/config.yaml)")
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
		viper.SetConfigFile(filepath.Join(home, ".config", "kandji-iru-cli", "config.yaml"))
	}

	viper.SetEnvPrefix("KANDJI")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))

	// So config file keys "base-url" and "token" are found when we GetString("base_url") / GetString("token").
	viper.RegisterAlias("base_url", "base-url")

	if err := viper.ReadInConfig(); err != nil {
		// Config file or directory not created yet (e.g. before running "init") — not fatal.
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return
		}
		if errors.Is(err, os.ErrNotExist) {
			return
		}
		var pathErr *os.PathError
		if errors.As(err, &pathErr) && errors.Is(pathErr.Err, os.ErrNotExist) {
			return
		}
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func validateConfig(cmd *cobra.Command) error {
	// Skip validation for init, completion, and token so they work without credentials.
	if len(os.Args) >= 2 {
		switch os.Args[1] {
		case "init", "completion", "token":
			return nil
		}
	}
	token := resolveToken(cmd)
	if token == "" {
		configPath := viper.ConfigFileUsed()
		if configPath == "" {
			home, _ := os.UserHomeDir()
			if home != "" {
				configPath = filepath.Join(home, ".config", "kandji-iru-cli", "config.yaml")
			} else {
				configPath = "~/.config/kandji-iru-cli/config.yaml"
			}
		}
		return fmt.Errorf("API token not configured\n"+
			"For best security, store your token in the system keyring:\n"+
			"  kandji-iru-cli init --keyring   (interactive setup)\n"+
			"  kandji-iru-cli token store     (store from --token, KANDJI_TOKEN, or stdin)\n"+
			"Otherwise add token to your config file or set KANDJI_TOKEN.\n"+
			"Config file: %s", configPath)
	}
	viper.Set("token", token)

	baseURL, err := resolveBaseURL(cmd)
	if err != nil {
		return err
	}
	viper.Set("resolved_base_url", strings.TrimSuffix(baseURL, "/"))
	return nil
}

// resolveToken returns the API token in order: --token flag (if set) > keyring > config file > env.
func resolveToken(cmd *cobra.Command) string {
	// 1. Explicit --token flag
	if f := cmd.Root().PersistentFlags().Lookup("token"); f != nil && f.Changed {
		if t := viper.GetString("token"); t != "" {
			return t
		}
	}
	// 2. System keyring (preferred over config/env)
	if kr, err := keyring.GetToken(); err == nil && kr != "" {
		return kr
	}
	// 3. Config file (before env so file takes precedence)
	if t := getTokenFromConfigFile(); t != "" {
		return t
	}
	// 4. Environment
	return os.Getenv("KANDJI_TOKEN")
}

// resolveBaseURL returns the API base URL in order: flag > keyring > config file > env; or builds from subdomain.
func resolveBaseURL(cmd *cobra.Command) (string, error) {
	// 1. Explicit --base-url flag
	if f := cmd.Root().PersistentFlags().Lookup("base-url"); f != nil && f.Changed {
		if u := strings.TrimSpace(viper.GetString("base_url")); u != "" {
			return u, nil
		}
	}
	// 2. Keyring
	if u, err := keyring.GetBaseURL(); err == nil && strings.TrimSpace(u) != "" {
		return strings.TrimSuffix(strings.TrimSpace(u), "/"), nil
	}
	// 3. Config file
	if u := getBaseURLFromConfigFile(); u != "" {
		return strings.TrimSuffix(strings.TrimSpace(u), "/"), nil
	}
	// 4. Environment
	if u := os.Getenv("KANDJI_BASE_URL"); strings.TrimSpace(u) != "" {
		return strings.TrimSuffix(strings.TrimSpace(u), "/"), nil
	}
	// No base URL — try subdomain to build URL
	subdomain := resolveSubdomain(cmd)
	if subdomain != "" {
		return fmt.Sprintf(defaultBaseURLTemplate, subdomain), nil
	}
	return "", fmt.Errorf("API base URL not configured\n"+
		"Set KANDJI_BASE_URL or KANDJI_SUBDOMAIN, use --base-url or --subdomain, or store them in the keyring (kandji-iru-cli init --keyring)")
}

// resolveSubdomain returns subdomain in order: flag > keyring > config file > env.
func resolveSubdomain(cmd *cobra.Command) string {
	if f := cmd.Root().PersistentFlags().Lookup("subdomain"); f != nil && f.Changed {
		if s := strings.TrimSpace(viper.GetString("subdomain")); s != "" {
			return s
		}
	}
	if s, err := keyring.GetSubdomain(); err == nil && strings.TrimSpace(s) != "" {
		return strings.TrimSpace(s)
	}
	if s := getSubdomainFromConfigFile(); s != "" {
		return strings.TrimSpace(s)
	}
	return strings.TrimSpace(os.Getenv("KANDJI_SUBDOMAIN"))
}

// getTokenFromConfigFile returns the token key from the config file only (ignores env).
func getTokenFromConfigFile() string {
	path := viper.ConfigFileUsed()
	if path == "" {
		return ""
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	var cfg struct {
		Token    string `yaml:"token"`
		BaseURL  string `yaml:"base-url"`
		Subdomain string `yaml:"subdomain"`
	}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return ""
	}
	return strings.TrimSpace(cfg.Token)
}

// getBaseURLFromConfigFile returns the base-url key from the config file only.
func getBaseURLFromConfigFile() string {
	path := viper.ConfigFileUsed()
	if path == "" {
		return ""
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	var cfg struct {
		BaseURL string `yaml:"base-url"`
	}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return ""
	}
	return strings.TrimSpace(cfg.BaseURL)
}

// getSubdomainFromConfigFile returns the subdomain key from the config file only.
func getSubdomainFromConfigFile() string {
	path := viper.ConfigFileUsed()
	if path == "" {
		return ""
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	var cfg struct {
		Subdomain string `yaml:"subdomain"`
	}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return ""
	}
	return strings.TrimSpace(cfg.Subdomain)
}
