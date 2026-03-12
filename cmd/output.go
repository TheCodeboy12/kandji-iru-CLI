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

// writeJSON writes v as indented JSON to stdout (for any list/get response).
func writeJSON(v interface{}) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
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
	return writeJSON(d)
}

func writeAuditEventsTable(events []kandji.AuditEvent) {
	if len(events) == 0 {
		fmt.Fprintln(os.Stdout, "No audit events.")
		return
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.Header("ID", "Occurred At", "Action", "Actor", "Target Type", "Target ID")
	for _, e := range events {
		table.Append([]string{e.ID, e.OccurredAt, e.Action, e.ActorID, e.TargetType, e.TargetID})
	}
	table.Render()
}

func writeBlueprintsTable(bps []kandji.Blueprint) {
	if len(bps) == 0 {
		fmt.Fprintln(os.Stdout, "No blueprints.")
		return
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.Header("ID", "Name", "Type", "Computers")
	for _, b := range bps {
		table.Append([]string{b.ID, b.Name, b.Type, fmt.Sprintf("%d", b.ComputersCount)})
	}
	table.Render()
}

func writeUsersTable(users []kandji.DirectoryUser) {
	if len(users) == 0 {
		fmt.Fprintln(os.Stdout, "No users.")
		return
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.Header("ID", "Name", "Email", "Device Count", "Archived")
	for _, u := range users {
		archived := "false"
		if u.Archived {
			archived = "true"
		}
		table.Append([]string{u.ID, u.Name, u.Email, fmt.Sprintf("%d", u.DeviceCount), archived})
	}
	table.Render()
}
