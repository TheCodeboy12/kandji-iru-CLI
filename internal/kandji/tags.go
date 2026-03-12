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

// Tag matches one tag from the tags API (Postman).
type Tag struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// ListTagsOptions holds query params for list tags.
type ListTagsOptions struct {
	Search string
}

// QueryValues returns url.Values for list tags.
func (o ListTagsOptions) QueryValues() url.Values {
	v := url.Values{}
	if o.Search != "" {
		v.Set("search", o.Search)
	}
	return v
}

// ListTags calls GET /api/v1/tags.
func (c *Client) ListTags(ctx context.Context, opts ListTagsOptions) ([]Tag, error) {
	path := apiPathPrefix + "/tags"
	if q := opts.QueryValues().Encode(); q != "" {
		path += "?" + q
	}
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("list tags: %w", err)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("list tags: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("list tags read: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("list tags: %s: %s", resp.Status, string(body))
	}
	var tags []Tag
	if err := json.Unmarshal(body, &tags); err != nil {
		return nil, fmt.Errorf("list tags decode: %w", err)
	}
	return tags, nil
}

// GetTag calls GET /api/v1/tags/{tag_id}.
func (c *Client) GetTag(ctx context.Context, tagID string) (*Tag, error) {
	path := apiPathPrefix + "/tags/" + url.PathEscape(tagID)
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("get tag: %w", err)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("get tag: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("tag not found: %s", tagID)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("get tag read: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get tag: %s: %s", resp.Status, string(body))
	}
	var tag Tag
	if err := json.Unmarshal(body, &tag); err != nil {
		return nil, fmt.Errorf("get tag decode: %w", err)
	}
	return &tag, nil
}

// CreateTagRequest is the body for POST /api/v1/tags.
type CreateTagRequest struct {
	Name string `json:"name"`
}

// CreateTag calls POST /api/v1/tags.
func (c *Client) CreateTag(ctx context.Context, name string) (*Tag, error) {
	path := apiPathPrefix + "/tags"
	body := CreateTagRequest{Name: name}
	raw, _ := json.Marshal(body)
	req, err := c.newRequest(ctx, http.MethodPost, path, bytes.NewReader(raw))
	if err != nil {
		return nil, fmt.Errorf("create tag: %w", err)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("create tag: %w", err)
	}
	defer resp.Body.Close()
	out, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("create tag read: %w", err)
	}
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("create tag: %s: %s", resp.Status, string(out))
	}
	var tag Tag
	if err := json.Unmarshal(out, &tag); err != nil {
		return nil, fmt.Errorf("create tag decode: %w", err)
	}
	return &tag, nil
}

// UpdateTagRequest is the body for PATCH /api/v1/tags/{tag_id}.
type UpdateTagRequest struct {
	Name string `json:"name"`
}

// UpdateTag calls PATCH /api/v1/tags/{tag_id}.
func (c *Client) UpdateTag(ctx context.Context, tagID, name string) (*Tag, error) {
	path := apiPathPrefix + "/tags/" + url.PathEscape(tagID)
	body := UpdateTagRequest{Name: name}
	raw, _ := json.Marshal(body)
	req, err := c.newRequest(ctx, http.MethodPatch, path, bytes.NewReader(raw))
	if err != nil {
		return nil, fmt.Errorf("update tag: %w", err)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("update tag: %w", err)
	}
	defer resp.Body.Close()
	out, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("update tag read: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("update tag: %s: %s", resp.Status, string(out))
	}
	var tag Tag
	if err := json.Unmarshal(out, &tag); err != nil {
		return nil, fmt.Errorf("update tag decode: %w", err)
	}
	return &tag, nil
}

// DeleteTag calls DELETE /api/v1/tags/{tag_id}.
func (c *Client) DeleteTag(ctx context.Context, tagID string) error {
	path := apiPathPrefix + "/tags/" + url.PathEscape(tagID)
	req, err := c.newRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return fmt.Errorf("delete tag: %w", err)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("delete tag: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNoContent || resp.StatusCode == http.StatusOK {
		return nil
	}
	body, _ := io.ReadAll(resp.Body)
	return fmt.Errorf("delete tag: %s: %s", resp.Status, string(body))
}
