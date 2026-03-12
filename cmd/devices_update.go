package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"kandji-iru-cli/internal/kandji"
)

var devicesUpdateCmd = &cobra.Command{
	Use:   "update [device_id]",
	Short: "Update device attributes",
	Long:  `PATCH device: user, asset_tag, blueprint_id, tags. Use flags or --body for JSON.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runDevicesUpdate,
}

func init() {
	devicesCmd.AddCommand(devicesUpdateCmd)
	devicesUpdateCmd.Flags().String("user", "", "Assigned user ID (UUID); use empty to clear")
	devicesUpdateCmd.Flags().String("asset-tag", "", "Asset tag; use --clear-asset-tag to clear")
	devicesUpdateCmd.Flags().Bool("clear-asset-tag", false, "Clear asset tag")
	devicesUpdateCmd.Flags().String("blueprint-id", "", "Blueprint ID")
	devicesUpdateCmd.Flags().StringSlice("tags", nil, "Tags (comma-separated); use --clear-tags to clear all")
	devicesUpdateCmd.Flags().Bool("clear-tags", false, "Clear all tags")
	devicesUpdateCmd.Flags().String("body", "", "Raw JSON body instead of flags")
	_ = viper.BindPFlag("devices_update_user", devicesUpdateCmd.Flags().Lookup("user"))
	_ = viper.BindPFlag("devices_update_asset_tag", devicesUpdateCmd.Flags().Lookup("asset-tag"))
	_ = viper.BindPFlag("devices_update_clear_asset_tag", devicesUpdateCmd.Flags().Lookup("clear-asset-tag"))
	_ = viper.BindPFlag("devices_update_blueprint_id", devicesUpdateCmd.Flags().Lookup("blueprint-id"))
	_ = viper.BindPFlag("devices_update_tags", devicesUpdateCmd.Flags().Lookup("tags"))
	_ = viper.BindPFlag("devices_update_clear_tags", devicesUpdateCmd.Flags().Lookup("clear-tags"))
	_ = viper.BindPFlag("devices_update_body", devicesUpdateCmd.Flags().Lookup("body"))
}

func runDevicesUpdate(cmd *cobra.Command, args []string) error {
	client := kandji.New(viper.GetString("resolved_base_url"), viper.GetString("token"))
	deviceID := args[0]

	if b := viper.GetString("devices_update_body"); b != "" {
		var payload kandji.UpdateDeviceRequest
		if err := json.Unmarshal([]byte(b), &payload); err != nil {
			return fmt.Errorf("invalid --body JSON: %w", err)
		}
		out, err := client.UpdateDevice(cmd.Context(), deviceID, payload)
		if err != nil {
			return err
		}
		os.Stdout.Write(out)
		if len(out) > 0 && out[len(out)-1] != '\n' {
			fmt.Println()
		}
		return nil
	}

	var payload kandji.UpdateDeviceRequest
	if u := viper.GetString("devices_update_user"); u != "" {
		payload.User = u
	}
	if viper.GetBool("devices_update_clear_asset_tag") {
		payload.AssetTag = nil
	} else if t := viper.GetString("devices_update_asset_tag"); t != "" {
		payload.AssetTag = t
	}
	if b := viper.GetString("devices_update_blueprint_id"); b != "" {
		payload.BlueprintID = b
	}
	if viper.GetBool("devices_update_clear_tags") {
		payload.Tags = []string{}
	} else if tags := viper.GetStringSlice("devices_update_tags"); len(tags) > 0 {
		payload.Tags = tags
	}

	out, err := client.UpdateDevice(cmd.Context(), deviceID, payload)
	if err != nil {
		return err
	}
	os.Stdout.Write(out)
	if len(out) > 0 && out[len(out)-1] != '\n' {
		fmt.Println()
	}
	return nil
}
