package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"kandji-iru-cli/internal/kandji"
)

// Postman device actions: lock, restart, blankpush, dailycheckin, renewmdmprofile,
// setname, clearpasscode, erase, deleteuser, reinstallagent, remotedesktop, etc.
var devicesActionCmd = &cobra.Command{
	Use:   "action [device_id] [action]",
	Short: "Run a device action",
	Long: `Send an action to a device. Action is the path segment (e.g. lock, restart, blankpush, dailycheckin).
Use --body to send JSON (e.g. for setname: {"DeviceName":"NewName"}).
Examples:
  kandji-iru-cli devices action <device_id> lock
  kandji-iru-cli devices action <device_id> blankpush
  kandji-iru-cli devices action <device_id> setname --body '{"DeviceName":"My Mac"}'`,
	Args: cobra.ExactArgs(2),
	RunE: runDevicesAction,
}

var devicesActionBody string

func init() {
	devicesCmd.AddCommand(devicesActionCmd)
	devicesActionCmd.Flags().StringVar(&devicesActionBody, "body", "", "JSON body for actions that require it (e.g. setname, erase)")
	_ = viper.BindPFlag("devices_action_body", devicesActionCmd.Flags().Lookup("body"))
}

func runDevicesAction(cmd *cobra.Command, args []string) error {
	deviceID := args[0]
	action := strings.TrimSpace(args[1])
	if action == "" {
		return fmt.Errorf("action cannot be empty")
	}

	baseURL := viper.GetString("resolved_base_url")
	token := viper.GetString("token")
	client := kandji.New(baseURL, token)

	var body io.Reader
	if b := viper.GetString("devices_action_body"); b != "" {
		body = bytes.NewReader([]byte(b))
	}

	code, out, err := client.DeviceAction(cmd.Context(), deviceID, action, body)
	if err != nil {
		return err
	}
	if len(out) > 0 {
		os.Stdout.Write(out)
		if out[len(out)-1] != '\n' {
			fmt.Println()
		}
	} else if code == 200 || code == 204 {
		fmt.Fprintln(os.Stdout, "OK")
	}
	return nil
}
