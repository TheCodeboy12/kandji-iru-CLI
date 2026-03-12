package cmd

import (
	"fmt"
	"net/url"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"kandji-iru-cli/internal/kandji"
)

// deviceInfoSubcommands: activity, apps, library-items, parameters, status
var devicesActivityCmd = &cobra.Command{
	Use:   "activity [device_id]",
	Short: "Get device activity",
	Long:  `GET /api/v1/devices/{id}/activity. Use -o json for raw output.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runDeviceInfo("activity"),
}

var devicesAppsCmd = &cobra.Command{
	Use:   "apps [device_id]",
	Short: "Get device apps",
	Long:  `GET /api/v1/devices/{id}/apps. Use -o json for raw output.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runDeviceInfo("apps"),
}

var devicesLibraryItemsCmd = &cobra.Command{
	Use:   "library-items [device_id]",
	Short: "Get device library items",
	Long:  `GET /api/v1/devices/{id}/library-items. Use -o json for raw output.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runDeviceInfo("library-items"),
}

var devicesParametersCmd = &cobra.Command{
	Use:   "parameters [device_id]",
	Short: "Get device parameters",
	Long:  `GET /api/v1/devices/{id}/parameters. Use -o json for raw output.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runDeviceInfo("parameters"),
}

var devicesStatusCmd = &cobra.Command{
	Use:   "status [device_id]",
	Short: "Get device status",
	Long:  `GET /api/v1/devices/{id}/status. Use -o json for raw output.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runDeviceInfo("status"),
}

func init() {
	devicesCmd.AddCommand(devicesActivityCmd)
	devicesCmd.AddCommand(devicesAppsCmd)
	devicesCmd.AddCommand(devicesLibraryItemsCmd)
	devicesCmd.AddCommand(devicesParametersCmd)
	devicesCmd.AddCommand(devicesStatusCmd)

	devicesActivityCmd.Flags().Int("limit", 0, "Limit number of results")
	devicesActivityCmd.Flags().String("params", "", "Extra query params as JSON (e.g. {\"limit\":50}). Merges with flags.")
	_ = viper.BindPFlag("devices_activity_limit", devicesActivityCmd.Flags().Lookup("limit"))
	_ = viper.BindPFlag("devices_activity_params", devicesActivityCmd.Flags().Lookup("params"))
}

func runDeviceInfo(subpath string) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		client := kandji.New(viper.GetString("resolved_base_url"), viper.GetString("token"))
		q := url.Values{}
		if subpath == "activity" {
			if l := viper.GetInt("devices_activity_limit"); l > 0 {
				q.Set("limit", fmt.Sprintf("%d", l))
			}
			for k, v := range parseExtraParams(viper.GetString("devices_activity_params")) {
				if v != "" {
					q.Set(k, v)
				}
			}
		}
		body, err := client.GetDeviceJSON(cmd.Context(), args[0], subpath, q)
		if err != nil {
			return err
		}
		os.Stdout.Write(body)
		if len(body) > 0 && body[len(body)-1] != '\n' {
			fmt.Println()
		}
		return nil
	}
}
