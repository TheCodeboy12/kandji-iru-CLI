package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"kandji-iru-cli/internal/kandji"
)

var devicesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List devices",
	Long: `List devices enrolled in Kandji. Optional filters match the API query parameters.
Use -o json to print raw JSON instead of a table.`,
	RunE: runDevicesList,
}

func init() {
	devicesCmd.AddCommand(devicesListCmd)

	devicesListCmd.Flags().Int("limit", 0, "Max number of devices to return (API default if not set)")
	devicesListCmd.Flags().Int("offset", 0, "Offset for pagination")
	devicesListCmd.Flags().String("device-id", "", "Filter by device ID")
	devicesListCmd.Flags().String("device-name", "", "Filter by device name")
	devicesListCmd.Flags().String("serial-number", "", "Filter by serial number")
	devicesListCmd.Flags().String("mac-address", "", "Filter by MAC address")
	devicesListCmd.Flags().String("user-name", "", "Filter by user name")
	devicesListCmd.Flags().String("user-email", "", "Filter by user email")
	devicesListCmd.Flags().String("platform", "", "Filter by platform (e.g. Mac, iPhone)")
	devicesListCmd.Flags().String("blueprint-id", "", "Filter by blueprint ID")

	_ = viper.BindPFlag("devices_list_limit", devicesListCmd.Flags().Lookup("limit"))
	_ = viper.BindPFlag("devices_list_offset", devicesListCmd.Flags().Lookup("offset"))
	_ = viper.BindPFlag("devices_list_device_id", devicesListCmd.Flags().Lookup("device-id"))
	_ = viper.BindPFlag("devices_list_device_name", devicesListCmd.Flags().Lookup("device-name"))
	_ = viper.BindPFlag("devices_list_serial_number", devicesListCmd.Flags().Lookup("serial-number"))
	_ = viper.BindPFlag("devices_list_mac_address", devicesListCmd.Flags().Lookup("mac-address"))
	_ = viper.BindPFlag("devices_list_user_name", devicesListCmd.Flags().Lookup("user-name"))
	_ = viper.BindPFlag("devices_list_user_email", devicesListCmd.Flags().Lookup("user-email"))
	_ = viper.BindPFlag("devices_list_platform", devicesListCmd.Flags().Lookup("platform"))
	_ = viper.BindPFlag("devices_list_blueprint_id", devicesListCmd.Flags().Lookup("blueprint-id"))
}

func runDevicesList(cmd *cobra.Command, args []string) error {
	baseURL := viper.GetString("resolved_base_url")
	token := viper.GetString("token")
	client := kandji.New(baseURL, token)

	opts := kandji.ListDeviceOptions{
		Limit:        viper.GetInt("devices_list_limit"),
		Offset:       viper.GetInt("devices_list_offset"),
		DeviceID:     viper.GetString("devices_list_device_id"),
		DeviceName:   viper.GetString("devices_list_device_name"),
		SerialNumber: viper.GetString("devices_list_serial_number"),
		MacAddress:   viper.GetString("devices_list_mac_address"),
		UserName:     viper.GetString("devices_list_user_name"),
		UserEmail:    viper.GetString("devices_list_user_email"),
		Platform:     viper.GetString("devices_list_platform"),
		BlueprintID:  viper.GetString("devices_list_blueprint_id"),
	}

	devices, err := client.ListDevices(cmd.Context(), opts)
	if err != nil {
		return fmt.Errorf("list devices: %w", err)
	}

	switch outputFormat() {
	case "json":
		return writeDevicesJSON(devices)
	default:
		writeDevicesTable(devices)
		return nil
	}
}
