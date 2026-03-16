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

// configTemplatePlain is used when init --no-keyring: token and base-url in file (less secure).
const configTemplatePlain = `# Kandji API configuration
# Edit this file and add your token and base URL (or subdomain).
# Get your API token: Kandji tenant > Settings > Access
#
# For better security, use keyring instead: kandji-iru-cli init (default)
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

// configTemplateKeyring is the default: no token in file (stored in keyring).
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
	Short: "Create config and store credentials (default: keyring)",
	Long: `Initialize the CLI. By default uses the system keyring (most secure): token is stored
in the keyring, and you can set base URL or subdomain in the config file when prompted.

Use --no-keyring only if you need to store the API token in the config file (less secure);
then edit the file to add token and base-url (or subdomain).

One-line init (non-interactive): pass --api-key and --base-url (or --subdomain) to configure
without prompts, e.g. kandji-iru-cli init --api-key=YOUR_TOKEN --base-url=https://tenant.api.kandji.io`,
	RunE: runInit,
}

var initForce bool
var initNoKeyring bool
var initApiKey string

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().BoolVar(&initForce, "force", false, "Overwrite config file if it already exists")
	initCmd.Flags().BoolVar(&initNoKeyring, "no-keyring", false, "Store token in config file instead of keyring (less secure)")
	initCmd.Flags().StringVar(&initApiKey, "api-key", "", "API token (with --base-url or --subdomain for one-line init)")
}

// initToken returns the token for init: init --api-key flag, or global --token / KANDJI_TOKEN.
func initToken() string {
	if s := strings.TrimSpace(initApiKey); s != "" {
		return s
	}
	return strings.TrimSpace(viper.GetString("token"))
}

// initBaseURL returns the base URL for init from global --base-url or KANDJI_BASE_URL.
// We read the flag directly so user-passed --base-url is not overwritten by config file.
func initBaseURL() string {
	if f := rootCmd.PersistentFlags().Lookup("base-url"); f != nil && f.Changed {
		return strings.TrimSpace(f.Value.String())
	}
	return strings.TrimSpace(viper.GetString("base_url"))
}

// initSubdomain returns the subdomain for init from global --subdomain or KANDJI_SUBDOMAIN.
func initSubdomain() string {
	if f := rootCmd.PersistentFlags().Lookup("subdomain"); f != nil && f.Changed {
		return strings.TrimSpace(f.Value.String())
	}
	return strings.TrimSpace(viper.GetString("subdomain"))
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
		if initNoKeyring {
			fmt.Fprintln(os.Stderr, "Edit it to update your token and base-url (or subdomain). Use --force to overwrite with the default template.")
		} else {
			fmt.Fprintln(os.Stderr, "Use --force to overwrite. To only store a token in keyring, use: kandji-iru-cli token store")
		}
		return nil
	}

	useKeyring := !initNoKeyring
	template := configTemplateKeyring
	if initNoKeyring {
		template = configTemplatePlain
	}
	if err := os.WriteFile(path, []byte(template), 0600); err != nil {
		return fmt.Errorf("write config file: %w", err)
	}

	if exists {
		fmt.Fprintf(os.Stderr, "Overwrote config file: %s\n", path)
	} else {
		fmt.Fprintf(os.Stderr, "Created config file: %s\n", path)
	}

	if useKeyring {
		token := initToken()
		baseURL := initBaseURL()
		subdomain := initSubdomain()
		hasURL := baseURL != "" || subdomain != ""

		// One-line init: token + base URL (or subdomain) provided → no prompts
		if token != "" && hasURL {
			if err := keyring.SetToken(token); err != nil {
				return fmt.Errorf("storing token in keyring: %w", err)
			}
			val := baseURL
			if val == "" {
				val = subdomain
			}
			if err := writeKeyringConfigWithURL(path, val); err != nil {
				return fmt.Errorf("writing config file: %w", err)
			}
			fmt.Fprintln(os.Stderr, "Initialized. Token stored in keyring; base URL in config file.")
			return nil
		}

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
		if hasURL {
			val := baseURL
			if val == "" {
				val = subdomain
			}
			if err := writeKeyringConfigWithURL(path, val); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: could not update config file with URL: %v\n", err)
			} else {
				fmt.Fprintln(os.Stderr, "Base URL written to config file.")
			}
		} else {
			fmt.Fprint(os.Stderr, "Base URL or subdomain for config file (optional; e.g. https://tenant.api.kandji.io or tenant): ")
			scanner := bufio.NewScanner(os.Stdin)
			if scanner.Scan() {
				val := strings.TrimSpace(scanner.Text())
				if val != "" {
					if err := writeKeyringConfigWithURL(path, val); err != nil {
						fmt.Fprintf(os.Stderr, "Warning: could not update config file with URL: %v\n", err)
					} else if strings.HasPrefix(val, "http://") || strings.HasPrefix(val, "https://") {
						fmt.Fprintln(os.Stderr, "Base URL written to config file.")
					} else {
						fmt.Fprintln(os.Stderr, "Subdomain written to config file.")
					}
				}
			}
			_ = scanner.Err()
		}
		fmt.Fprintln(os.Stderr, "Done. Token is in keyring; base URL is in the config file (edit it anytime).")
		return nil
	}

	fmt.Fprintln(os.Stderr, "Edit the config file and add your token and base-url (or subdomain), then run any CLI command.")
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

// writeKeyringConfigWithURL rewrites the config file at path with base-url or subdomain set.
// val is either a full URL (https://...) or a subdomain name.
func writeKeyringConfigWithURL(path, val string) error {
	const header = `# Kandji API configuration (token stored in system keyring)
# Edit base-url or subdomain below. Token is read from the keyring.
#
# US tenants:
#   Set base-url: https://YOUR_TENANT.api.kandji.io
#   Or set subdomain: YOUR_TENANT (URL will be built for you)
#
# EU tenants:
#   Set base-url: https://YOUR_TENANT.api.eu.kandji.io

`
	// Escape for YAML double-quoted string (backslash and quote)
	escape := func(s string) string {
		s = strings.ReplaceAll(s, "\\", "\\\\")
		return strings.ReplaceAll(s, "\"", "\\\"")
	}
	var body string
	if strings.HasPrefix(val, "http://") || strings.HasPrefix(val, "https://") {
		body = header + "base-url: \"" + escape(val) + "\"\n# subdomain: \"\"\n"
	} else {
		body = header + "base-url: \"\"\nsubdomain: \"" + escape(val) + "\"\n"
	}
	return os.WriteFile(path, []byte(body), 0600)
}
