package cmd

import (
	"github.com/spf13/cobra"
)

var usersCmd = &cobra.Command{
	Use:   "users",
	Short: "Directory users",
	Long:  `List and get users from user directory integrations.`,
}

func init() {
	rootCmd.AddCommand(usersCmd)
}
