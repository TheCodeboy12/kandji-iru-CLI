package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"kandji-iru-cli/internal/kandji"
)

var usersDeleteCmd = &cobra.Command{
	Use:   "delete [user_id]",
	Short: "Delete a directory user",
	Long:  `DELETE /api/v1/users/{user_id}. Removes the user record from the directory integration.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runUsersDelete,
}

func init() {
	usersCmd.AddCommand(usersDeleteCmd)
}

func runUsersDelete(cmd *cobra.Command, args []string) error {
	client := kandji.New(viper.GetString("resolved_base_url"), viper.GetString("token"))
	if err := client.DeleteUser(cmd.Context(), args[0]); err != nil {
		return err
	}
	fmt.Fprintln(os.Stdout, "User deleted.")
	return nil
}
