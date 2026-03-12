package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"kandji-iru-cli/internal/kandji"
)

var tagsDeleteCmd = &cobra.Command{
	Use:   "delete [tag_id]",
	Short: "Delete a tag",
	Args:  cobra.ExactArgs(1),
	RunE:  runTagsDelete,
}

func init() {
	tagsCmd.AddCommand(tagsDeleteCmd)
}

func runTagsDelete(cmd *cobra.Command, args []string) error {
	client := kandji.New(viper.GetString("resolved_base_url"), viper.GetString("token"))
	if err := client.DeleteTag(cmd.Context(), args[0]); err != nil {
		return err
	}
	fmt.Fprintln(os.Stdout, "Deleted.")
	return nil
}
