package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"kandji-iru-cli/internal/kandji"
)

var tagsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List tags",
	Long:  `List tags. Use --search to filter. Use -o json for raw output.`,
	RunE:  runTagsList,
}

func init() {
	tagsCmd.AddCommand(tagsListCmd)
	tagsListCmd.Flags().String("search", "", "Search filter")
	tagsListCmd.Flags().String("params", "", "Extra query params as JSON (e.g. {\"search\":\"eng\"}). Merges with/overrides flags.")
	_ = viper.BindPFlag("tags_list_search", tagsListCmd.Flags().Lookup("search"))
	_ = viper.BindPFlag("tags_list_params", tagsListCmd.Flags().Lookup("params"))
}

func runTagsList(cmd *cobra.Command, args []string) error {
	client := kandji.New(viper.GetString("resolved_base_url"), viper.GetString("token"))
	opts := kandji.ListTagsOptions{
		Search:      viper.GetString("tags_list_search"),
		ExtraParams: parseExtraParams(viper.GetString("tags_list_params")),
	}
	if outputFormat() == "raw" {
		body, err := client.ListTagsRaw(cmd.Context(), opts)
		if err != nil {
			return fmt.Errorf("list tags: %w", err)
		}
		_, _ = os.Stdout.Write(body)
		return nil
	}
	tags, err := client.ListTags(cmd.Context(), opts)
	if err != nil {
		return fmt.Errorf("list tags: %w", err)
	}
	if outputFormat() == "json" {
		return writeJSON(tags)
	}
	if len(tags) == 0 {
		fmt.Fprintln(os.Stdout, "No tags.")
		return nil
	}
	for _, t := range tags {
		fmt.Fprintf(os.Stdout, "%s\t%s\n", t.ID, t.Name)
	}
	return nil
}
