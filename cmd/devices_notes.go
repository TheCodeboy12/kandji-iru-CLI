package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"kandji-iru-cli/internal/kandji"
)

var devicesNotesListCmd = &cobra.Command{
	Use:   "list [device_id]",
	Short: "List device notes",
	Long:  `List all notes for a device.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runDevicesNotesList,
}

var devicesNoteGetCmd = &cobra.Command{
	Use:   "get [device_id] [note_id]",
	Short: "Get a device note",
	Args:  cobra.ExactArgs(2),
	RunE:  runDevicesNoteGet,
}

var devicesNoteCreateCmd = &cobra.Command{
	Use:   "create [device_id]",
	Short: "Create a device note",
	Long:  `Create a note for a device. Use --content for the note body.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runDevicesNoteCreate,
}

var devicesNoteUpdateCmd = &cobra.Command{
	Use:   "update [device_id] [note_id]",
	Short: "Update a device note",
	Long:  `Update a note. Use --content for the new body.`,
	Args:  cobra.ExactArgs(2),
	RunE:  runDevicesNoteUpdate,
}

var devicesNoteDeleteCmd = &cobra.Command{
	Use:   "delete [device_id] [note_id]",
	Short: "Delete a device note",
	Args:  cobra.ExactArgs(2),
	RunE:  runDevicesNoteDelete,
}

var devicesNotesCmd = &cobra.Command{
	Use:   "notes",
	Short: "Manage device notes",
}

func init() {
	devicesCmd.AddCommand(devicesNotesCmd)
	devicesNotesCmd.AddCommand(devicesNotesListCmd)
	devicesNotesCmd.AddCommand(devicesNoteGetCmd)
	devicesNotesCmd.AddCommand(devicesNoteCreateCmd)
	devicesNotesCmd.AddCommand(devicesNoteUpdateCmd)
	devicesNotesCmd.AddCommand(devicesNoteDeleteCmd)

	devicesNoteCreateCmd.Flags().String("content", "", "Note content")
	_ = viper.BindPFlag("devices_note_content", devicesNoteCreateCmd.Flags().Lookup("content"))
	devicesNoteUpdateCmd.Flags().String("content", "", "Note content")
	_ = viper.BindPFlag("devices_note_update_content", devicesNoteUpdateCmd.Flags().Lookup("content"))
}

func runDevicesNotesList(cmd *cobra.Command, args []string) error {
	client := kandji.New(viper.GetString("resolved_base_url"), viper.GetString("token"))
	notes, err := client.ListDeviceNotes(cmd.Context(), args[0])
	if err != nil {
		return fmt.Errorf("list device notes: %w", err)
	}
	if outputFormat() == "json" {
		return writeJSON(notes)
	}
	if len(notes) == 0 {
		fmt.Fprintln(os.Stdout, "No notes.")
		return nil
	}
	for _, n := range notes {
		fmt.Fprintf(os.Stdout, "%s\t%s\t%s\n", n.ID, n.CreatedAt, n.Content)
	}
	return nil
}

func runDevicesNoteGet(cmd *cobra.Command, args []string) error {
	client := kandji.New(viper.GetString("resolved_base_url"), viper.GetString("token"))
	note, err := client.GetDeviceNote(cmd.Context(), args[0], args[1])
	if err != nil {
		return err
	}
	if outputFormat() == "json" {
		return writeJSON(note)
	}
	fmt.Fprintf(os.Stdout, "ID: %s\nCreated: %s\nUpdated: %s\nContent: %s\n", note.ID, note.CreatedAt, note.UpdatedAt, note.Content)
	return nil
}

func runDevicesNoteCreate(cmd *cobra.Command, args []string) error {
	content := viper.GetString("devices_note_content")
	if content == "" {
		return fmt.Errorf("--content is required")
	}
	client := kandji.New(viper.GetString("resolved_base_url"), viper.GetString("token"))
	note, err := client.CreateDeviceNote(cmd.Context(), args[0], content)
	if err != nil {
		return err
	}
	if outputFormat() == "json" {
		return writeJSON(note)
	}
	fmt.Fprintf(os.Stdout, "Created note %s\n", note.ID)
	return nil
}

func runDevicesNoteUpdate(cmd *cobra.Command, args []string) error {
	content := viper.GetString("devices_note_update_content")
	if content == "" {
		return fmt.Errorf("--content is required")
	}
	client := kandji.New(viper.GetString("resolved_base_url"), viper.GetString("token"))
	note, err := client.UpdateDeviceNote(cmd.Context(), args[0], args[1], content)
	if err != nil {
		return err
	}
	if outputFormat() == "json" {
		return writeJSON(note)
	}
	fmt.Fprintf(os.Stdout, "Updated note %s\n", note.ID)
	return nil
}

func runDevicesNoteDelete(cmd *cobra.Command, args []string) error {
	client := kandji.New(viper.GetString("resolved_base_url"), viper.GetString("token"))
	if err := client.DeleteDeviceNote(cmd.Context(), args[0], args[1]); err != nil {
		return err
	}
	fmt.Fprintln(os.Stdout, "Deleted.")
	return nil
}
