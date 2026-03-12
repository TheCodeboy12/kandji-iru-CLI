package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"kandji-iru-cli/internal/kandji"
)

var devicesSecretsCmd = &cobra.Command{
	Use:   "secrets [device_id] [type]",
	Short: "Get device secret",
	Long:  `Get a device secret. Type: bypasscode, filevaultkey, unlockpin, recoverypassword.`,
	Args:  cobra.ExactArgs(2),
	RunE:  runDevicesSecrets,
}

func init() {
	devicesCmd.AddCommand(devicesSecretsCmd)
}

func runDevicesSecrets(cmd *cobra.Command, args []string) error {
	client := kandji.New(viper.GetString("resolved_base_url"), viper.GetString("token"))
	body, err := client.GetDeviceSecret(cmd.Context(), args[0], args[1])
	if err != nil {
		return err
	}
	os.Stdout.Write(body)
	if len(body) > 0 && body[len(body)-1] != '\n' {
		fmt.Println()
	}
	return nil
}
