package kandji

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// Blueprint matches one item from GET /api/v1/blueprints (Postman collection).
type Blueprint struct {
	ID              string                 `json:"id"`
	Name            string                 `json:"name"`
	Icon            string                 `json:"icon,omitempty"`
	Color           string                 `json:"color,omitempty"`
	Description     string                 `json:"description,omitempty"`
	Params          map[string]interface{} `json:"params,omitempty"`
	ComputersCount  int                    `json:"computers_count,omitempty"`
	EnrollmentCode *struct {
		Code     string `json:"code"`
		IsActive bool   `json:"is_active"`
	} `json:"enrollment_code,omitempty"`
	Type string `json:"type,omitempty"`
}

// BlueprintListResponse is the response for GET /api/v1/blueprints.
type BlueprintListResponse struct {
	Count    int         `json:"count"`
	Next     *string     `json:"next"`
	Previous *string     `json:"previous"`
	Results  []Blueprint `json:"results"`
}

// ListBlueprintsOptions holds query params for list blueprints.
type ListBlueprintsOptions struct {
	ID     string
	IDIn   string
	Name   string
	Limit  int
	Offset int
}

// QueryValues returns url.Values for list blueprints.
func (o ListBlueprintsOptions) QueryValues() url.Values {
	v := url.Values{}
	if o.ID != "" {
		v.Set("id", o.ID)
	}
	if o.IDIn != "" {
		v.Set("id__in", o.IDIn)
	}
	if o.Name != "" {
		v.Set("name", o.Name)
	}
	if o.Limit > 0 {
		v.Set("limit", fmt.Sprintf("%d", o.Limit))
	}
	if o.Offset > 0 {
		v.Set("offset", fmt.Sprintf("%d", o.Offset))
	}
	return v
}

// ListBlueprints calls GET /api/v1/blueprints.
func (c *Client) ListBlueprints(ctx context.Context, opts ListBlueprintsOptions) (*BlueprintListResponse, error) {
	path := apiPathPrefix + "/blueprints"
	if q := opts.QueryValues().Encode(); q != "" {
		path += "?" + q
	}
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("list blueprints: %w", err)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("list blueprints: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("list blueprints read: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("list blueprints: %s: %s", resp.Status, string(body))
	}
	var out BlueprintListResponse
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, fmt.Errorf("list blueprints decode: %w", err)
	}
	return &out, nil
}

// GetBlueprint calls GET /api/v1/blueprints/{blueprint_id}.
func (c *Client) GetBlueprint(ctx context.Context, blueprintID string) (*Blueprint, error) {
	path := apiPathPrefix + "/blueprints/" + url.PathEscape(blueprintID)
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("get blueprint: %w", err)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("get blueprint: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("blueprint not found: %s", blueprintID)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("get blueprint read: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get blueprint: %s: %s", resp.Status, string(body))
	}
	var out Blueprint
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, fmt.Errorf("get blueprint decode: %w", err)
	}
	return &out, nil
}

// ListBlueprintLibraryItems calls GET /api/v1/blueprints/{blueprint_id}/list-library-items.
// Returns raw JSON (array or object per API).
func (c *Client) ListBlueprintLibraryItems(ctx context.Context, blueprintID string) ([]byte, error) {
	path := apiPathPrefix + "/blueprints/" + url.PathEscape(blueprintID) + "/list-library-items"
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("list blueprint library items: %w", err)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("list blueprint library items: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("list blueprint library items read: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("list blueprint library items: %s: %s", resp.Status, string(body))
	}
	return body, nil
}

// GetBlueprintTemplates calls GET /api/v1/blueprints/templates/.
// Returns raw JSON.
func (c *Client) GetBlueprintTemplates(ctx context.Context) ([]byte, error) {
	path := apiPathPrefix + "/blueprints/templates/"
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("get blueprint templates: %w", err)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("get blueprint templates: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("get blueprint templates read: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get blueprint templates: %s: %s", resp.Status, string(body))
	}
	return body, nil
}
