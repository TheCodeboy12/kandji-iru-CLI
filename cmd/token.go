package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"kandji-iru-cli/internal/keyring"
)

var tokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Manage API token in system keyring",
	Long: `Store or remove the Kandji API token in the system keyring (macOS Keychain,
Linux Secret Service, Windows Credential Manager). When the token is stored here,
you do not need to set KANDJI_TOKEN, --token, or token in the config file.`,
}

var tokenStoreBaseURL, tokenStoreSubdomain string

var tokenStoreCmd = &cobra.Command{
	Use:   "store",
	Short: "Store API token and/or base URL in system keyring",
	Long: `Store the Kandji API token and/or base URL (or subdomain) in the system keyring.

Token: from --token, KANDJI_TOKEN, or stdin. Optional if you only want to store --base-url or --subdomain.

Base URL / subdomain: use --base-url or --subdomain so the CLI works without a config file.`,
	RunE: runTokenStore,
}

var tokenDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Remove API token from system keyring",
	Long:  `Remove the Kandji API token from the system keyring. This does not remove the token from your config file or environment.`,
	RunE:  runTokenDelete,
}

func init() {
	rootCmd.AddCommand(tokenCmd)
	tokenCmd.AddCommand(tokenStoreCmd)
	tokenCmd.AddCommand(tokenDeleteCmd)
	tokenStoreCmd.Flags().StringVar(&tokenStoreBaseURL, "base-url", "", "Store this base URL in keyring (e.g. https://your-tenant.api.kandji.io)")
	tokenStoreCmd.Flags().StringVar(&tokenStoreSubdomain, "subdomain", "", "Store this subdomain in keyring (base URL will be built from it)")
}

func runTokenStore(cmd *cobra.Command, args []string) error {
	stored := false
	token := viper.GetString("token")
	if token == "" && tokenStoreBaseURL == "" && tokenStoreSubdomain == "" {
		// Read from stdin (e.g. echo "token" | kandji-iru-cli token store)
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			token = strings.TrimSpace(scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			return fmt.Errorf("reading token from stdin: %w", err)
		}
	}
	if strings.TrimSpace(token) != "" {
		if err := keyring.SetToken(token); err != nil {
			return fmt.Errorf("storing token in keyring: %w", err)
		}
		fmt.Fprintln(os.Stderr, "Token stored in system keyring.")
		stored = true
	}
	if strings.TrimSpace(tokenStoreBaseURL) != "" {
		if err := keyring.SetBaseURL(strings.TrimSpace(tokenStoreBaseURL)); err != nil {
			return fmt.Errorf("storing base URL in keyring: %w", err)
		}
		fmt.Fprintln(os.Stderr, "Base URL stored in system keyring.")
		stored = true
	}
	if strings.TrimSpace(tokenStoreSubdomain) != "" {
		if err := keyring.SetSubdomain(strings.TrimSpace(tokenStoreSubdomain)); err != nil {
			return fmt.Errorf("storing subdomain in keyring: %w", err)
		}
		fmt.Fprintln(os.Stderr, "Subdomain stored in system keyring.")
		stored = true
	}
	if !stored {
		return fmt.Errorf("nothing to store: provide token (--token, KANDJI_TOKEN, or stdin), --base-url, or --subdomain")
	}
	return nil
}

func runTokenDelete(cmd *cobra.Command, args []string) error {
	if err := keyring.DeleteToken(); err != nil {
		return fmt.Errorf("deleting token from keyring: %w", err)
	}
	fmt.Fprintln(os.Stderr, "Token removed from system keyring.")
	return nil
}
