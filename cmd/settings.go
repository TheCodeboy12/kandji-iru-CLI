package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"kandji-iru-cli/internal/kandji"
)

var settingsCmd = &cobra.Command{
	Use:   "settings",
	Short: "Settings and licensing",
}

var settingsLicensingCmd = &cobra.Command{
	Use:   "licensing",
	Short: "Get licensing information",
	Long:  `GET /api/v1/settings/licensing. Use -o json for raw output.`,
	RunE:  runSettingsLicensing,
}

func init() {
	rootCmd.AddCommand(settingsCmd)
	settingsCmd.AddCommand(settingsLicensingCmd)
}

func runSettingsLicensing(cmd *cobra.Command, args []string) error {
	client := kandji.New(viper.GetString("resolved_base_url"), viper.GetString("token"))
	body, err := client.GetLicensing(cmd.Context())
	if err != nil {
		return err
	}
	os.Stdout.Write(body)
	if len(body) > 0 && body[len(body)-1] != '\n' {
		fmt.Println()
	}
	return nil
}
