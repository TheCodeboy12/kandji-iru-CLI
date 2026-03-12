package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"kandji-iru-cli/internal/kandji"
)

var usersGetCmd = &cobra.Command{
	Use:   "get [user_id]",
	Short: "Get directory user by ID",
	Long:  `Fetch a single directory user by UUID. Use -o json to print raw JSON.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runUsersGet,
}

func init() {
	usersCmd.AddCommand(usersGetCmd)
}

func runUsersGet(cmd *cobra.Command, args []string) error {
	userID := args[0]
	baseURL := viper.GetString("resolved_base_url")
	token := viper.GetString("token")
	client := kandji.New(baseURL, token)

	u, err := client.GetUser(cmd.Context(), userID)
	if err != nil {
		return fmt.Errorf("get user: %w", err)
	}

	switch outputFormat() {
	case "json":
		return writeJSON(u)
	default:
		writeUsersTable([]kandji.DirectoryUser{*u})
		return nil
	}
}
