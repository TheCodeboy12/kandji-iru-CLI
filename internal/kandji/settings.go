package kandji

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// GetLicensing calls GET /api/v1/settings/licensing. Returns raw JSON (structure may vary).
func (c *Client) GetLicensing(ctx context.Context) ([]byte, error) {
	path := apiPathPrefix + "/settings/licensing"
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
		return nil, fmt.Errorf("get licensing: %s: %s", resp.Status, string(body))
	}
	return body, nil
}

// ParseLicensing parses raw response for optional table use.
func ParseLicensing(raw []byte) (map[string]interface{}, error) {
	var m map[string]interface{}
	err := json.Unmarshal(raw, &m)
	return m, err
}
