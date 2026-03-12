package kandji

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// DirectoryUser is a user from the directory (GET /api/v1/users). Named to avoid clash with Device.User.
type DirectoryUser struct {
	Active           bool   `json:"active"`
	Archived         bool   `json:"archived"`
	CreatedAt        string `json:"created_at,omitempty"`
	Department       string `json:"department,omitempty"`
	DeprecatedUserID string `json:"deprecated_user_id,omitempty"`
	Email            string `json:"email,omitempty"`
	ID               string `json:"id"`
	Integration      *struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		UUID string `json:"uuid"`
		Type string `json:"type"`
	} `json:"integration,omitempty"`
	JobTitle    string `json:"job_title,omitempty"`
	Name        string `json:"name,omitempty"`
	DeviceCount int    `json:"device_count,omitempty"`
	UpdatedAt   string `json:"updated_at,omitempty"`
}

// DirectoryUserListResponse is the response for GET /api/v1/users.
type DirectoryUserListResponse struct {
	Next     *string         `json:"next"`
	Previous *string         `json:"previous"`
	Results  []DirectoryUser `json:"results"`
}

// ListDirectoryUsersOptions holds query params for list users.
type ListDirectoryUsersOptions struct {
	Email         string
	ID            string
	IntegrationID string
	Archived      string
	Cursor        string
}

// QueryValues returns url.Values for list users.
func (o ListDirectoryUsersOptions) QueryValues() url.Values {
	v := url.Values{}
	if o.Email != "" {
		v.Set("email", o.Email)
	}
	if o.ID != "" {
		v.Set("id", o.ID)
	}
	if o.IntegrationID != "" {
		v.Set("integration_id", o.IntegrationID)
	}
	if o.Archived != "" {
		v.Set("archived", o.Archived)
	}
	if o.Cursor != "" {
		v.Set("cursor", o.Cursor)
	}
	return v
}

// ListUsers calls GET /api/v1/users.
func (c *Client) ListUsers(ctx context.Context, opts ListDirectoryUsersOptions) (*DirectoryUserListResponse, error) {
	path := apiPathPrefix + "/users"
	if q := opts.QueryValues().Encode(); q != "" {
		path += "?" + q
	}
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("list users read: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("list users: %s: %s", resp.Status, string(body))
	}
	var out DirectoryUserListResponse
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, fmt.Errorf("list users decode: %w", err)
	}
	return &out, nil
}

// ListUsersRaw returns the raw response body from GET /api/v1/users.
func (c *Client) ListUsersRaw(ctx context.Context, opts ListDirectoryUsersOptions) ([]byte, error) {
	path := apiPathPrefix + "/users"
	if q := opts.QueryValues().Encode(); q != "" {
		path += "?" + q
	}
	return c.GetRaw(ctx, path)
}

// GetUserRaw returns the raw response body from GET /api/v1/users/{user_id}.
func (c *Client) GetUserRaw(ctx context.Context, userID string) ([]byte, error) {
	return c.GetRaw(ctx, apiPathPrefix+"/users/"+url.PathEscape(userID))
}

// GetUser calls GET /api/v1/users/{user_id}.
func (c *Client) GetUser(ctx context.Context, userID string) (*DirectoryUser, error) {
	path := apiPathPrefix + "/users/" + url.PathEscape(userID)
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("user not found: %s", userID)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("get user read: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get user: %s: %s", resp.Status, string(body))
	}
	var out DirectoryUser
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, fmt.Errorf("get user decode: %w", err)
	}
	return &out, nil
}

// DeleteUser calls DELETE /api/v1/users/{user_id}.
func (c *Client) DeleteUser(ctx context.Context, userID string) error {
	path := apiPathPrefix + "/users/" + url.PathEscape(userID)
	req, err := c.newRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNoContent || resp.StatusCode == http.StatusOK {
		return nil
	}
	body, _ := io.ReadAll(resp.Body)
	return fmt.Errorf("delete user: %s: %s", resp.Status, string(body))
}
