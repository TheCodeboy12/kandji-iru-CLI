package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"kandji-iru-cli/internal/kandji"
)

var usersListCmd = &cobra.Command{
	Use:   "list",
	Short: "List directory users",
	Long:  `List users from directory integrations. Use -o json to print raw JSON.`,
	RunE:  runUsersList,
}

func init() {
	usersCmd.AddCommand(usersListCmd)

	usersListCmd.Flags().String("email", "", "Filter by email (containing)")
	usersListCmd.Flags().String("id", "", "Search by user UUID")
	usersListCmd.Flags().String("integration-id", "", "Filter by integration UUID")
	usersListCmd.Flags().String("archived", "", "Filter archived: true or false")
	usersListCmd.Flags().String("cursor", "", "Pagination cursor")
	usersListCmd.Flags().String("params", "", "Extra query params as JSON (e.g. {\"email\":\"@company.com\"}). Merges with/overrides flags.")

	_ = viper.BindPFlag("users_list_email", usersListCmd.Flags().Lookup("email"))
	_ = viper.BindPFlag("users_list_id", usersListCmd.Flags().Lookup("id"))
	_ = viper.BindPFlag("users_list_integration_id", usersListCmd.Flags().Lookup("integration-id"))
	_ = viper.BindPFlag("users_list_archived", usersListCmd.Flags().Lookup("archived"))
	_ = viper.BindPFlag("users_list_cursor", usersListCmd.Flags().Lookup("cursor"))
	_ = viper.BindPFlag("users_list_params", usersListCmd.Flags().Lookup("params"))
}

func runUsersList(cmd *cobra.Command, args []string) error {
	baseURL := viper.GetString("resolved_base_url")
	token := viper.GetString("token")
	client := kandji.New(baseURL, token)

	opts := kandji.ListDirectoryUsersOptions{
		Email:         viper.GetString("users_list_email"),
		ID:            viper.GetString("users_list_id"),
		IntegrationID: viper.GetString("users_list_integration_id"),
		Archived:      viper.GetString("users_list_archived"),
		Cursor:        viper.GetString("users_list_cursor"),
		ExtraParams:   parseExtraParams(viper.GetString("users_list_params")),
	}

	switch outputFormat() {
	case "raw":
		body, err := client.ListUsersRaw(cmd.Context(), opts)
		if err != nil {
			return fmt.Errorf("list users: %w", err)
		}
		_, _ = os.Stdout.Write(body)
		return nil
	}

	resp, err := client.ListUsers(cmd.Context(), opts)
	if err != nil {
		return fmt.Errorf("list users: %w", err)
	}

	switch outputFormat() {
	case "json":
		return writeJSON(resp)
	default:
		writeUsersTable(resp.Results)
		printPaginationHint("users", resp.Next, resp.Previous, len(resp.Results), 0, 0)
		return nil
	}
}
