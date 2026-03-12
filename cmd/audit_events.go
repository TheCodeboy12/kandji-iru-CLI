package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"kandji-iru-cli/internal/kandji"
)

var auditEventsCmd = &cobra.Command{
	Use:   "events",
	Short: "List audit events",
	Long: `List audit log events (blueprint/library item changes, device lifecycle, API token actions, etc.).
Use -o json to print raw JSON.`,
	RunE: runAuditEvents,
}

func init() {
	auditCmd.AddCommand(auditEventsCmd)

	auditEventsCmd.Flags().Int("limit", 500, "Max records per request (API max 500)")
	auditEventsCmd.Flags().String("sort-by", "-occurred_at", "Sort by occurred_at or id; prefix with - for descending")
	auditEventsCmd.Flags().String("start-date", "", "Filter by start date (ISO or YYYY-MM-DD)")
	auditEventsCmd.Flags().String("end-date", "", "Filter by end date (ISO or YYYY-MM-DD)")
	auditEventsCmd.Flags().String("cursor", "", "Pagination cursor")
	auditEventsCmd.Flags().String("params", "", "Extra query params as JSON (e.g. {\"limit\":100}). Merges with/overrides flags.")

	_ = viper.BindPFlag("audit_events_limit", auditEventsCmd.Flags().Lookup("limit"))
	_ = viper.BindPFlag("audit_events_sort_by", auditEventsCmd.Flags().Lookup("sort-by"))
	_ = viper.BindPFlag("audit_events_start_date", auditEventsCmd.Flags().Lookup("start-date"))
	_ = viper.BindPFlag("audit_events_end_date", auditEventsCmd.Flags().Lookup("end-date"))
	_ = viper.BindPFlag("audit_events_cursor", auditEventsCmd.Flags().Lookup("cursor"))
	_ = viper.BindPFlag("audit_events_params", auditEventsCmd.Flags().Lookup("params"))
}

func runAuditEvents(cmd *cobra.Command, args []string) error {
	baseURL := viper.GetString("resolved_base_url")
	token := viper.GetString("token")
	client := kandji.New(baseURL, token)

	opts := kandji.ListAuditEventsOptions{
		Limit:       viper.GetInt("audit_events_limit"),
		SortBy:      viper.GetString("audit_events_sort_by"),
		StartDate:   viper.GetString("audit_events_start_date"),
		EndDate:     viper.GetString("audit_events_end_date"),
		Cursor:      viper.GetString("audit_events_cursor"),
		ExtraParams: parseExtraParams(viper.GetString("audit_events_params")),
	}
	if opts.Limit <= 0 {
		opts.Limit = 500
	}
	if opts.SortBy == "" {
		opts.SortBy = "-occurred_at"
	}

	switch outputFormat() {
	case "raw":
		body, err := client.ListAuditEventsRaw(cmd.Context(), opts)
		if err != nil {
			return fmt.Errorf("audit events: %w", err)
		}
		_, _ = os.Stdout.Write(body)
		return nil
	}

	resp, err := client.ListAuditEvents(cmd.Context(), opts)
	if err != nil {
		return fmt.Errorf("audit events: %w", err)
	}

	switch outputFormat() {
	case "json":
		return writeJSON(resp)
	default:
		writeAuditEventsTable(resp.Results)
		printPaginationHint("audit", resp.Next, resp.Previous, len(resp.Results), opts.Limit, 0)
		return nil
	}
}
