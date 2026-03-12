package kandji

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// AuditEvent matches one entry in GET /api/v1/audit/events (Postman collection).
type AuditEvent struct {
	ID              string                 `json:"id"`
	Action          string                 `json:"action"`
	OccurredAt      string                 `json:"occurred_at"`
	ActorID         string                 `json:"actor_id"`
	ActorType       string                 `json:"actor_type"`
	TargetID        string                 `json:"target_id"`
	TargetType      string                 `json:"target_type"`
	TargetComponent string                 `json:"target_component"`
	NewState        map[string]interface{} `json:"new_state,omitempty"`
	Metadata        map[string]interface{} `json:"metadata,omitempty"`
}

// AuditEventsResponse is the response for GET /api/v1/audit/events.
type AuditEventsResponse struct {
	Results  []AuditEvent `json:"results"`
	Previous *string      `json:"previous"`
	Next     *string      `json:"next"`
}

// ListAuditEventsOptions holds query params for audit events.
type ListAuditEventsOptions struct {
	Limit     int
	SortBy    string
	StartDate string
	EndDate   string
	Cursor    string
}

// QueryValues returns url.Values for the audit events request.
func (o ListAuditEventsOptions) QueryValues() url.Values {
	v := url.Values{}
	if o.Limit > 0 {
		v.Set("limit", fmt.Sprintf("%d", o.Limit))
	}
	if o.SortBy != "" {
		v.Set("sort_by", o.SortBy)
	}
	if o.StartDate != "" {
		v.Set("start_date", o.StartDate)
	}
	if o.EndDate != "" {
		v.Set("end_date", o.EndDate)
	}
	if o.Cursor != "" {
		v.Set("cursor", o.Cursor)
	}
	return v
}

// ListAuditEvents calls GET /api/v1/audit/events.
func (c *Client) ListAuditEvents(ctx context.Context, opts ListAuditEventsOptions) (*AuditEventsResponse, error) {
	path := apiPathPrefix + "/audit/events"
	if q := opts.QueryValues().Encode(); q != "" {
		path += "?" + q
	}
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("list audit events: %w", err)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("list audit events: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("list audit events read: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("list audit events: %s: %s", resp.Status, string(body))
	}
	var out AuditEventsResponse
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, fmt.Errorf("list audit events decode: %w", err)
	}
	return &out, nil
}

// ListAuditEventsRaw returns the raw response body from GET /api/v1/audit/events.
func (c *Client) ListAuditEventsRaw(ctx context.Context, opts ListAuditEventsOptions) ([]byte, error) {
	path := apiPathPrefix + "/audit/events"
	if q := opts.QueryValues().Encode(); q != "" {
		path += "?" + q
	}
	return c.GetRaw(ctx, path)
}
