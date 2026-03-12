package cmd

import (
	"github.com/spf13/cobra"
)

var auditCmd = &cobra.Command{
	Use:   "audit",
	Short: "Audit log and activity",
	Long:  `Query audit log events from the Activity module.`,
}

func init() {
	rootCmd.AddCommand(auditCmd)
}
