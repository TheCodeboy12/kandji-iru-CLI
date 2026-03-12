package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"kandji-iru-cli/internal/kandji"
)

var blueprintsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List blueprints",
	Long:  `List blueprints. Use -o json to print raw JSON.`,
	RunE:  runBlueprintsList,
}

func init() {
	blueprintsCmd.AddCommand(blueprintsListCmd)

	blueprintsListCmd.Flags().String("id", "", "Look up a specific blueprint by ID")
	blueprintsListCmd.Flags().String("id-in", "", "Comma-separated blueprint IDs")
	blueprintsListCmd.Flags().String("name", "", "Filter by name (containing)")
	blueprintsListCmd.Flags().Int("limit", 0, "Results per page")
	blueprintsListCmd.Flags().Int("offset", 0, "Starting index")
	blueprintsListCmd.Flags().String("params", "", "Extra query params as JSON (e.g. {\"name\":\"Engineering\"}). Merges with/overrides flags.")

	_ = viper.BindPFlag("blueprints_list_id", blueprintsListCmd.Flags().Lookup("id"))
	_ = viper.BindPFlag("blueprints_list_id_in", blueprintsListCmd.Flags().Lookup("id-in"))
	_ = viper.BindPFlag("blueprints_list_name", blueprintsListCmd.Flags().Lookup("name"))
	_ = viper.BindPFlag("blueprints_list_limit", blueprintsListCmd.Flags().Lookup("limit"))
	_ = viper.BindPFlag("blueprints_list_offset", blueprintsListCmd.Flags().Lookup("offset"))
	_ = viper.BindPFlag("blueprints_list_params", blueprintsListCmd.Flags().Lookup("params"))
}

func runBlueprintsList(cmd *cobra.Command, args []string) error {
	baseURL := viper.GetString("resolved_base_url")
	token := viper.GetString("token")
	client := kandji.New(baseURL, token)

	opts := kandji.ListBlueprintsOptions{
		ID:          viper.GetString("blueprints_list_id"),
		IDIn:        viper.GetString("blueprints_list_id_in"),
		Name:        viper.GetString("blueprints_list_name"),
		Limit:       viper.GetInt("blueprints_list_limit"),
		Offset:      viper.GetInt("blueprints_list_offset"),
		ExtraParams: parseExtraParams(viper.GetString("blueprints_list_params")),
	}

	switch outputFormat() {
	case "raw":
		body, err := client.ListBlueprintsRaw(cmd.Context(), opts)
		if err != nil {
			return fmt.Errorf("list blueprints: %w", err)
		}
		_, _ = os.Stdout.Write(body)
		return nil
	}

	resp, err := client.ListBlueprints(cmd.Context(), opts)
	if err != nil {
		return fmt.Errorf("list blueprints: %w", err)
	}

	switch outputFormat() {
	case "json":
		return writeJSON(resp)
	default:
		writeBlueprintsTable(resp.Results)
		printPaginationHint("blueprints", resp.Next, resp.Previous, len(resp.Results), opts.Limit, opts.Offset)
		return nil
	}
}
