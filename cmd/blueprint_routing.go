package cmd

import (
	"fmt"
	"net/url"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"kandji-iru-cli/internal/kandji"
)

var blueprintRoutingCmd = &cobra.Command{
	Use:   "blueprint-routing",
	Short: "Blueprint routing (enrollment code)",
	Long:  `Get or update blueprint routing and view activity.`,
}

var blueprintRoutingGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get blueprint routing",
	Long:  `GET /api/v1/blueprint-routing/. Use -o json for raw output.`,
	RunE:  runBlueprintRoutingGet,
}

var blueprintRoutingActivityCmd = &cobra.Command{
	Use:   "activity",
	Short: "Get blueprint routing activity",
	Long:  `GET /api/v1/blueprint-routing/activity. Use -o json for raw output.`,
	RunE:  runBlueprintRoutingActivity,
}

func init() {
	rootCmd.AddCommand(blueprintRoutingCmd)
	blueprintRoutingCmd.AddCommand(blueprintRoutingGetCmd)
	blueprintRoutingCmd.AddCommand(blueprintRoutingActivityCmd)
	blueprintRoutingActivityCmd.Flags().Int("limit", 0, "Limit results")
	blueprintRoutingActivityCmd.Flags().String("params", "", "Extra query params as JSON (e.g. {\"limit\":100}). Merges with flags.")
	_ = viper.BindPFlag("blueprint_routing_activity_limit", blueprintRoutingActivityCmd.Flags().Lookup("limit"))
	_ = viper.BindPFlag("blueprint_routing_activity_params", blueprintRoutingActivityCmd.Flags().Lookup("params"))
}

func runBlueprintRoutingGet(cmd *cobra.Command, args []string) error {
	client := kandji.New(viper.GetString("resolved_base_url"), viper.GetString("token"))
	body, err := client.GetBlueprintRouting(cmd.Context())
	if err != nil {
		return err
	}
	os.Stdout.Write(body)
	if len(body) > 0 && body[len(body)-1] != '\n' {
		fmt.Println()
	}
	return nil
}

func runBlueprintRoutingActivity(cmd *cobra.Command, args []string) error {
	client := kandji.New(viper.GetString("resolved_base_url"), viper.GetString("token"))
	q := url.Values{}
	if limit := viper.GetInt("blueprint_routing_activity_limit"); limit > 0 {
		q.Set("limit", strconv.Itoa(limit))
	}
	for k, v := range parseExtraParams(viper.GetString("blueprint_routing_activity_params")) {
		if v != "" {
			q.Set(k, v)
		}
	}
	body, err := client.GetBlueprintRoutingActivity(cmd.Context(), q)
	if err != nil {
		return err
	}
	os.Stdout.Write(body)
	if len(body) > 0 && body[len(body)-1] != '\n' {
		fmt.Println()
	}
	return nil
}
