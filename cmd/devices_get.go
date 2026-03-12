package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"kandji-iru-cli/internal/kandji"
)

var devicesGetCmd = &cobra.Command{
	Use:   "get [device_id]",
	Short: "Get device details",
	Long: `Fetch full details for a single device by device_id (UUID).
Use -o json to print raw JSON instead of a table.`,
	Args: cobra.ExactArgs(1),
	RunE: runDevicesGet,
}

func init() {
	devicesCmd.AddCommand(devicesGetCmd)
}

func runDevicesGet(cmd *cobra.Command, args []string) error {
	deviceID := args[0]
	baseURL := viper.GetString("resolved_base_url")
	token := viper.GetString("token")
	client := kandji.New(baseURL, token)

	device, err := client.GetDevice(cmd.Context(), deviceID)
	if err != nil {
		return fmt.Errorf("get device: %w", err)
	}

	switch outputFormat() {
	case "json":
		return writeDeviceJSON(device)
	default:
		writeDevicesTable([]kandji.Device{*device})
		return nil
	}
}
