package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"kandji-iru-cli/internal/kandji"
)

var tagsUpdateCmd = &cobra.Command{
	Use:   "update [tag_id]",
	Short: "Update a tag",
	Long:  `Update a tag name. Use --name for the new name.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runTagsUpdate,
}

func init() {
	tagsCmd.AddCommand(tagsUpdateCmd)
	tagsUpdateCmd.Flags().String("name", "", "New tag name")
	_ = viper.BindPFlag("tags_update_name", tagsUpdateCmd.Flags().Lookup("name"))
}

func runTagsUpdate(cmd *cobra.Command, args []string) error {
	name := viper.GetString("tags_update_name")
	if name == "" {
		return fmt.Errorf("--name is required")
	}
	client := kandji.New(viper.GetString("resolved_base_url"), viper.GetString("token"))
	tag, err := client.UpdateTag(cmd.Context(), args[0], name)
	if err != nil {
		return err
	}
	if outputFormat() == "json" {
		return writeJSON(tag)
	}
	fmt.Fprintf(os.Stdout, "Updated tag %s\n", tag.ID)
	return nil
}
