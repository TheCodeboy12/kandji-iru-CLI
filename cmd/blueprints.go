package cmd

import (
	"github.com/spf13/cobra"
)

var blueprintsCmd = &cobra.Command{
	Use:   "blueprints",
	Short: "Manage blueprints",
	Long:  `List and inspect blueprints.`,
}

func init() {
	rootCmd.AddCommand(blueprintsCmd)
}
