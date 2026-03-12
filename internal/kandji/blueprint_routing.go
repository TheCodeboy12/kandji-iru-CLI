package kandji

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// GetBlueprintRouting calls GET /api/v1/blueprint-routing/.
func (c *Client) GetBlueprintRouting(ctx context.Context) ([]byte, error) {
	path := apiPathPrefix + "/blueprint-routing/"
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("get blueprint routing: %w", err)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("get blueprint routing: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("get blueprint routing read: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get blueprint routing: %s: %s", resp.Status, string(body))
	}
	return body, nil
}

// UpdateBlueprintRoutingRequest matches PATCH body (enrollment_code).
type UpdateBlueprintRoutingRequest struct {
	EnrollmentCode *struct {
		IsActive bool   `json:"is_active"`
		Code     string `json:"code"`
	} `json:"enrollment_code,omitempty"`
}

// UpdateBlueprintRouting calls PATCH /api/v1/blueprint-routing/.
func (c *Client) UpdateBlueprintRouting(ctx context.Context, payload UpdateBlueprintRoutingRequest) ([]byte, error) {
	path := apiPathPrefix + "/blueprint-routing/"
	raw, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := c.newRequest(ctx, http.MethodPatch, path, bytes.NewReader(raw))
	if err != nil {
		return nil, fmt.Errorf("update blueprint routing: %w", err)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("update blueprint routing: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("update blueprint routing read: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("update blueprint routing: %s: %s", resp.Status, string(body))
	}
	return body, nil
}

// GetBlueprintRoutingActivity calls GET /api/v1/blueprint-routing/activity.
func (c *Client) GetBlueprintRoutingActivity(ctx context.Context, limit int) ([]byte, error) {
	path := apiPathPrefix + "/blueprint-routing/activity"
	v := url.Values{}
	if limit > 0 {
		v.Set("limit", fmt.Sprintf("%d", limit))
	}
	if q := v.Encode(); q != "" {
		path += "?" + q
	}
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("get blueprint routing activity: %w", err)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("get blueprint routing activity: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("get blueprint routing activity read: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get blueprint routing activity: %s: %s", resp.Status, string(body))
	}
	return body, nil
}
