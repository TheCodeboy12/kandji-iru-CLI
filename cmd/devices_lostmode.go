package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"kandji-iru-cli/internal/kandji"
)

var devicesLostModeGetCmd = &cobra.Command{
	Use:   "get [device_id]",
	Short: "Get device lost mode details",
	Long:  `GET /api/v1/devices/{id}/details/lostmode. Use -o json for raw output.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runDevicesLostModeGet,
}

var devicesLostModeCancelCmd = &cobra.Command{
	Use:   "cancel [device_id]",
	Short: "Cancel lost mode",
	Long:  `DELETE /api/v1/devices/{id}/details/lostmode to disable lost mode.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runDevicesLostModeCancel,
}

var devicesLostModeCmd = &cobra.Command{
	Use:   "lostmode",
	Short: "Device lost mode",
}

func init() {
	devicesCmd.AddCommand(devicesLostModeCmd)
	devicesLostModeCmd.AddCommand(devicesLostModeGetCmd)
	devicesLostModeCmd.AddCommand(devicesLostModeCancelCmd)
}

func runDevicesLostModeGet(cmd *cobra.Command, args []string) error {
	client := kandji.New(viper.GetString("resolved_base_url"), viper.GetString("token"))
	body, err := client.GetDeviceLostMode(cmd.Context(), args[0])
	if err != nil {
		return err
	}
	os.Stdout.Write(body)
	if len(body) > 0 && body[len(body)-1] != '\n' {
		fmt.Println()
	}
	return nil
}

func runDevicesLostModeCancel(cmd *cobra.Command, args []string) error {
	client := kandji.New(viper.GetString("resolved_base_url"), viper.GetString("token"))
	if err := client.CancelDeviceLostMode(cmd.Context(), args[0]); err != nil {
		return err
	}
	fmt.Fprintln(os.Stdout, "Lost mode cancelled.")
	return nil
}
