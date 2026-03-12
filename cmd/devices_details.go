package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"kandji-iru-cli/internal/kandji"
)

var devicesDetailsCmd = &cobra.Command{
	Use:   "details [device_id]",
	Short: "Get full device details",
	Long:  `Fetch full device details (general, hardware, network, etc.) for a device.
Use -o json to print raw JSON.`,
	Args: cobra.ExactArgs(1),
	RunE: runDevicesDetails,
}

func init() {
	devicesCmd.AddCommand(devicesDetailsCmd)
}

func runDevicesDetails(cmd *cobra.Command, args []string) error {
	deviceID := args[0]
	baseURL := viper.GetString("resolved_base_url")
	token := viper.GetString("token")
	client := kandji.New(baseURL, token)

	body, err := client.GetDeviceDetails(cmd.Context(), deviceID)
	if err != nil {
		return fmt.Errorf("device details: %w", err)
	}

	switch outputFormat() {
	case "json", "raw":
		_, err = os.Stdout.Write(body)
		if err != nil {
			return err
		}
		fmt.Println()
		return nil
	default:
		_, err = os.Stdout.Write(body)
		if err != nil {
			return err
		}
		fmt.Println()
		return nil
	}
}
