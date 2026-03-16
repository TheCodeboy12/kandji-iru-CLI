package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"kandji-iru-cli/internal/keyring"
)

const configTemplate = `# Kandji API configuration
# Edit this file and add your token and base URL (or subdomain).
# Get your API token: Kandji tenant > Settings > Access
#
# For best security, store the token in the system keyring instead:
#   kandji-iru-cli init --keyring   (interactive) or kandji-iru-cli token store
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

// configTemplateKeyring is used when init --keyring: no token in file (stored in keyring).
const configTemplateKeyring = `# Kandji API configuration (token stored in system keyring)
# Edit base-url or subdomain below. Token is read from the keyring.
#
# US tenants:
#   Set base-url: https://YOUR_TENANT.api.kandji.io
#   Or set subdomain: YOUR_TENANT (URL will be built for you)
#
# EU tenants:
#   Set base-url: https://YOUR_TENANT.api.eu.kandji.io

base-url: ""
# subdomain: ""
`

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create the config file (optionally store token in keyring)",
	Long: `Create the config file and optionally store your API token in the system keyring.

Without --keyring: creates a config file with empty token and base-url; edit the file to add them.

With --keyring: creates the config file (no token in file) and stores the API token in the
system keyring (macOS Keychain, Linux Secret Service, Windows Credential Manager). Token can
be provided via --token, KANDJI_TOKEN, or stdin (e.g. echo "YOUR_TOKEN" | kandji-iru-cli init --keyring).
Then edit the config file to add base-url or subdomain only.`,
	RunE: runInit,
}

var initForce bool
var initKeyring bool

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().BoolVar(&initForce, "force", false, "Overwrite config file if it already exists")
	initCmd.Flags().BoolVar(&initKeyring, "keyring", false, "Store API token in system keyring (no token in config file)")
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

	exists := false
	if _, err := os.Stat(path); err == nil {
		exists = true
	}

	if exists && !initForce {
		fmt.Fprintf(os.Stderr, "Config file already exists: %s\n", path)
		if initKeyring {
			fmt.Fprintln(os.Stderr, "Use --force to overwrite. To only store a token in keyring, use: kandji-iru-cli token store")
		} else {
			fmt.Fprintln(os.Stderr, "Edit it to update your token and base-url (or subdomain). Use --force to overwrite with the default template.")
		}
		return nil
	}

	template := configTemplate
	if initKeyring {
		template = configTemplateKeyring
	}
	if err := os.WriteFile(path, []byte(template), 0600); err != nil {
		return fmt.Errorf("write config file: %w", err)
	}

	if exists {
		fmt.Fprintf(os.Stderr, "Overwrote config file: %s\n", path)
	} else {
		fmt.Fprintf(os.Stderr, "Created config file: %s\n", path)
	}

	if initKeyring {
		token := viper.GetString("token")
		if token == "" {
			fmt.Fprint(os.Stderr, "API token (will be stored in system keyring, or pipe from stdin): ")
			scanner := bufio.NewScanner(os.Stdin)
			if scanner.Scan() {
				token = strings.TrimSpace(scanner.Text())
			}
			_ = scanner.Err()
		}
		if token == "" {
			fmt.Fprintln(os.Stderr, "No token provided. Store it later with: kandji-iru-cli token store")
			fmt.Fprintln(os.Stderr, "Edit the config file to add base-url (or subdomain), then run any CLI command.")
			return nil
		}
		if err := keyring.SetToken(token); err != nil {
			return fmt.Errorf("storing token in keyring: %w", err)
		}
		fmt.Fprintln(os.Stderr, "Token stored in system keyring.")
		// Optionally store base URL or subdomain so CLI works without config file
		fmt.Fprint(os.Stderr, "Base URL or subdomain (optional, stored in keyring; e.g. https://tenant.api.kandji.io or tenant): ")
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			val := strings.TrimSpace(scanner.Text())
			if val != "" {
				if strings.HasPrefix(val, "http://") || strings.HasPrefix(val, "https://") {
					_ = keyring.SetBaseURL(val)
					fmt.Fprintln(os.Stderr, "Base URL stored in system keyring.")
				} else {
					_ = keyring.SetSubdomain(val)
					fmt.Fprintln(os.Stderr, "Subdomain stored in system keyring.")
				}
			}
		}
		_ = scanner.Err()
		fmt.Fprintln(os.Stderr, "You can run CLI commands without a config file; or edit the config file for base-url/subdomain only.")
		return nil
	}

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
	return filepath.Join(home, ".config", "kandji-iru-cli", "config.yaml")
}
