package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"kandji-iru-cli/internal/kandji"
)

var tagsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a tag",
	Long:  `Create a tag. Use --name for the tag name.`,
	RunE:  runTagsCreate,
}

func init() {
	tagsCmd.AddCommand(tagsCreateCmd)
	tagsCreateCmd.Flags().String("name", "", "Tag name")
	_ = viper.BindPFlag("tags_create_name", tagsCreateCmd.Flags().Lookup("name"))
}

func runTagsCreate(cmd *cobra.Command, args []string) error {
	name := viper.GetString("tags_create_name")
	if name == "" {
		return fmt.Errorf("--name is required")
	}
	client := kandji.New(viper.GetString("resolved_base_url"), viper.GetString("token"))
	tag, err := client.CreateTag(cmd.Context(), name)
	if err != nil {
		return err
	}
	if outputFormat() == "json" {
		return writeJSON(tag)
	}
	fmt.Fprintf(os.Stdout, "Created tag %s (%s)\n", tag.Name, tag.ID)
	return nil
}
