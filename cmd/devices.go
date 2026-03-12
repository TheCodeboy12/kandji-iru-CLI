package cmd

import (
	"github.com/spf13/cobra"
)

var devicesCmd = &cobra.Command{
	Use:   "devices",
	Short: "Manage Kandji devices",
	Long:  `List and inspect devices enrolled in Kandji.`,
}

func init() {
	rootCmd.AddCommand(devicesCmd)
}
