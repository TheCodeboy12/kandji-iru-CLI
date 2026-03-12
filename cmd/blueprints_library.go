package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"kandji-iru-cli/internal/kandji"
)

var blueprintsLibraryItemsCmd = &cobra.Command{
	Use:   "library-items [blueprint_id]",
	Short: "List library items assigned to a blueprint",
	Long:  `GET /api/v1/blueprints/{id}/list-library-items. Use -o json for raw output.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runBlueprintsLibraryItems,
}

var blueprintsTemplatesCmd = &cobra.Command{
	Use:   "templates",
	Short: "Get blueprint templates",
	Long:  `GET /api/v1/blueprints/templates/. Use -o json for raw output.`,
	RunE:  runBlueprintsTemplates,
}

func init() {
	blueprintsCmd.AddCommand(blueprintsLibraryItemsCmd)
	blueprintsCmd.AddCommand(blueprintsTemplatesCmd)
}

func runBlueprintsLibraryItems(cmd *cobra.Command, args []string) error {
	client := kandji.New(viper.GetString("resolved_base_url"), viper.GetString("token"))
	body, err := client.ListBlueprintLibraryItems(cmd.Context(), args[0])
	if err != nil {
		return err
	}
	os.Stdout.Write(body)
	if len(body) > 0 && body[len(body)-1] != '\n' {
		fmt.Println()
	}
	return nil
}

func runBlueprintsTemplates(cmd *cobra.Command, args []string) error {
	client := kandji.New(viper.GetString("resolved_base_url"), viper.GetString("token"))
	body, err := client.GetBlueprintTemplates(cmd.Context())
	if err != nil {
		return err
	}
	os.Stdout.Write(body)
	if len(body) > 0 && body[len(body)-1] != '\n' {
		fmt.Println()
	}
	return nil
}
