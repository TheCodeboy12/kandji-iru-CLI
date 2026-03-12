package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/viper"
	"kandji-iru-cli/internal/kandji"
)

func outputFormat() string {
	return viper.GetString("output")
}

func writeDevicesTable(devices []kandji.Device) {
	if len(devices) == 0 {
		// Avoid printing a header-only table; give a clear message
		fmt.Fprintln(os.Stdout, "No devices.")
		return
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.Header("Device ID", "Device Name", "Serial Number", "Platform", "User")
	for _, d := range devices {
		s := d.Summary()
		user := s.UserEmail
		if user == "" {
			user = s.UserName
		}
		table.Append([]string{s.DeviceID, s.DeviceName, s.SerialNumber, s.Platform, user})
	}
	table.Render()
}

func writeDevicesJSON(devices []kandji.Device) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(devices)
}

func writeDeviceJSON(d *kandji.Device) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(d)
}
