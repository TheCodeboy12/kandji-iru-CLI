package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"kandji-iru-cli/internal/kandji"
)

var blueprintsGetCmd = &cobra.Command{
	Use:   "get [blueprint_id]",
	Short: "Get blueprint by ID",
	Long:  `Fetch a single blueprint by ID. Use -o json to print raw JSON.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runBlueprintsGet,
}

func init() {
	blueprintsCmd.AddCommand(blueprintsGetCmd)
}

func runBlueprintsGet(cmd *cobra.Command, args []string) error {
	blueprintID := args[0]
	baseURL := viper.GetString("resolved_base_url")
	token := viper.GetString("token")
	client := kandji.New(baseURL, token)

	bp, err := client.GetBlueprint(cmd.Context(), blueprintID)
	if err != nil {
		return fmt.Errorf("get blueprint: %w", err)
	}

	switch outputFormat() {
	case "json":
		return writeJSON(bp)
	default:
		writeBlueprintsTable([]kandji.Blueprint{*bp})
		return nil
	}
}
