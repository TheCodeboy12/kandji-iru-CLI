package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"kandji-iru-cli/internal/kandji"
)

var tagsGetCmd = &cobra.Command{
	Use:   "get [tag_id]",
	Short: "Get a tag",
	Args:  cobra.ExactArgs(1),
	RunE:  runTagsGet,
}

func init() {
	tagsCmd.AddCommand(tagsGetCmd)
}

func runTagsGet(cmd *cobra.Command, args []string) error {
	client := kandji.New(viper.GetString("resolved_base_url"), viper.GetString("token"))
	if outputFormat() == "raw" {
		body, err := client.GetTagRaw(cmd.Context(), args[0])
		if err != nil {
			return err
		}
		_, _ = os.Stdout.Write(body)
		return nil
	}
	tag, err := client.GetTag(cmd.Context(), args[0])
	if err != nil {
		return err
	}
	if outputFormat() == "json" {
		return writeJSON(tag)
	}
	fmt.Fprintf(os.Stdout, "ID: %s\nName: %s\n", tag.ID, tag.Name)
	return nil
}
