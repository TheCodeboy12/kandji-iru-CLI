package kandji

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const apiPathPrefix = "/api/v1"

// Client calls the Kandji API with a base URL and bearer token.
type Client struct {
	baseURL    string
	token      string
	httpClient *http.Client
}

// New returns a Client. baseURL should not have a trailing slash.
func New(baseURL, token string) *Client {
	return &Client{
		baseURL: baseURL,
		token:   token,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// ListDevices calls GET /api/v1/devices with optional query params.
// Response is a JSON array of Device (official API; see Postman collection).
func (c *Client) ListDevices(ctx context.Context, opts ListDeviceOptions) ([]Device, error) {
	path := apiPathPrefix + "/devices"
	if q := opts.QueryValues().Encode(); q != "" {
		path += "?" + q
	}
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("list devices: %w", err)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("list devices: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("list devices read: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("list devices: %s: %s", resp.Status, string(body))
	}
	var devices []Device
	if err := json.Unmarshal(body, &devices); err != nil {
		return nil, fmt.Errorf("list devices decode: %w", err)
	}
	return devices, nil
}

// GetDevice calls GET /api/v1/devices/{device_id}.
func (c *Client) GetDevice(ctx context.Context, deviceID string) (*Device, error) {
	path := apiPathPrefix + "/devices/" + url.PathEscape(deviceID)
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("get device: %w", err)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("get device: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("device not found: %s", deviceID)
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("get device: %s: %s", resp.Status, string(body))
	}
	var out Device
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("get device decode: %w", err)
	}
	return &out, nil
}

// GetDeviceDetails calls GET /api/v1/devices/{device_id}/details (full device details).
func (c *Client) GetDeviceDetails(ctx context.Context, deviceID string) ([]byte, error) {
	path := apiPathPrefix + "/devices/" + url.PathEscape(deviceID) + "/details"
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("get device details: %w", err)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("get device details: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("get device details read: %w", err)
	}
	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("device not found: %s", deviceID)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get device details: %s: %s", resp.Status, string(body))
	}
	return body, nil
}

// DeviceAction sends POST /api/v1/devices/{device_id}/action/{action} with optional body.
// action is e.g. "lock", "restart", "blankpush", "dailycheckin", "renewmdmprofile", "setname", etc.
func (c *Client) DeviceAction(ctx context.Context, deviceID, action string, body io.Reader) (int, []byte, error) {
	path := apiPathPrefix + "/devices/" + url.PathEscape(deviceID) + "/action/" + url.PathEscape(action)
	req, err := c.newRequest(ctx, http.MethodPost, path, body)
	if err != nil {
		return 0, nil, fmt.Errorf("device action %s: %w", action, err)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, nil, fmt.Errorf("device action %s: %w", action, err)
	}
	defer resp.Body.Close()
	out, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return resp.StatusCode, out, nil
	}
	return resp.StatusCode, out, fmt.Errorf("device action %s: %s: %s", action, resp.Status, string(out))
}

func (c *Client) newRequest(ctx context.Context, method, path string, body io.Reader) (*http.Request, error) {
	u := c.baseURL + path
	req, err := http.NewRequestWithContext(ctx, method, u, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	return req, nil
}

// GetRaw performs GET and returns the response body unchanged (for -o raw).
func (c *Client) GetRaw(ctx context.Context, path string) ([]byte, error) {
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %s", resp.Status, string(body))
	}
	return body, nil
}

// ListDevicesRaw returns the raw response body from GET /api/v1/devices.
func (c *Client) ListDevicesRaw(ctx context.Context, opts ListDeviceOptions) ([]byte, error) {
	path := apiPathPrefix + "/devices"
	if q := opts.QueryValues().Encode(); q != "" {
		path += "?" + q
	}
	return c.GetRaw(ctx, path)
}

// GetDeviceRaw returns the raw response body from GET /api/v1/devices/{id}.
func (c *Client) GetDeviceRaw(ctx context.Context, deviceID string) ([]byte, error) {
	return c.GetRaw(ctx, apiPathPrefix+"/devices/"+url.PathEscape(deviceID))
}
