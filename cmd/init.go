package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

const configTemplate = `# Kandji API configuration
# Edit this file and add your token and base URL (or subdomain).
# Get your API token: Kandji tenant > Settings > Access
#
# US tenants:
#   Set base-url: https://YOUR_TENANT.api.kandji.io
#   Or set subdomain: YOUR_TENANT (URL will be built for you)
#
# EU tenants:
#   Set base-url: https://YOUR_TENANT.api.eu.kandji.io

token: ""
base-url: ""
# subdomain: ""
`

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create the config file",
	Long: `Create the config file so you can edit it with your API token and base URL.
Default path: ~/.kandji.yaml (or use --config to specify a path).
After running init, edit the file and add your token and base-url (or subdomain).`,
	RunE: runInit,
}

var initForce bool

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().BoolVar(&initForce, "force", false, "Overwrite config file if it already exists")
}

func runInit(cmd *cobra.Command, args []string) error {
	path := configPath()
	if path == "" {
		return fmt.Errorf("could not determine config file path (set --config or ensure HOME is set)")
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("create config directory: %w", err)
	}

	if _, err := os.Stat(path); err == nil && !initForce {
		return fmt.Errorf("config file already exists: %s\nUse --force to overwrite", path)
	}

	if err := os.WriteFile(path, []byte(configTemplate), 0600); err != nil {
		return fmt.Errorf("write config file: %w", err)
	}

	fmt.Fprintf(os.Stderr, "Created config file: %s\n", path)
	fmt.Fprintln(os.Stderr, "Edit it and add your token and base-url (or subdomain), then run any CLI command.")
	return nil
}

// configPath returns the path where the config file is or would be written.
func configPath() string {
	if cfgFile != "" {
		return cfgFile
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".kandji.yaml")
}
