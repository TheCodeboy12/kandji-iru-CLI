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

// DeviceNote matches one note in GET/POST/PATCH device notes (Postman).
type DeviceNote struct {
	ID        string `json:"id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

// ListDeviceNotes calls GET /api/v1/devices/{device_id}/notes.
func (c *Client) ListDeviceNotes(ctx context.Context, deviceID string) ([]DeviceNote, error) {
	path := apiPathPrefix + "/devices/" + url.PathEscape(deviceID) + "/notes"
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("list device notes: %w", err)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("list device notes: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("list device notes read: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("list device notes: %s: %s", resp.Status, string(body))
	}
	var notes []DeviceNote
	if err := json.Unmarshal(body, &notes); err != nil {
		return nil, fmt.Errorf("list device notes decode: %w", err)
	}
	return notes, nil
}

// ListDeviceNotesRaw returns the raw response body from GET /api/v1/devices/{id}/notes.
func (c *Client) ListDeviceNotesRaw(ctx context.Context, deviceID string) ([]byte, error) {
	path := apiPathPrefix + "/devices/" + url.PathEscape(deviceID) + "/notes"
	return c.GetRaw(ctx, path)
}

// GetDeviceNoteRaw returns the raw response body from GET /api/v1/devices/{id}/notes/{note_id}.
func (c *Client) GetDeviceNoteRaw(ctx context.Context, deviceID, noteID string) ([]byte, error) {
	path := apiPathPrefix + "/devices/" + url.PathEscape(deviceID) + "/notes/" + url.PathEscape(noteID)
	return c.GetRaw(ctx, path)
}

// GetDeviceNote calls GET /api/v1/devices/{device_id}/notes/{note_id}.
func (c *Client) GetDeviceNote(ctx context.Context, deviceID, noteID string) (*DeviceNote, error) {
	path := apiPathPrefix + "/devices/" + url.PathEscape(deviceID) + "/notes/" + url.PathEscape(noteID)
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("get device note: %w", err)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("get device note: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("note not found: %s", noteID)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("get device note read: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get device note: %s: %s", resp.Status, string(body))
	}
	var note DeviceNote
	if err := json.Unmarshal(body, &note); err != nil {
		return nil, fmt.Errorf("get device note decode: %w", err)
	}
	return &note, nil
}

// CreateDeviceNoteRequest is the body for POST device notes.
type CreateDeviceNoteRequest struct {
	Content string `json:"content"`
}

// CreateDeviceNote calls POST /api/v1/devices/{device_id}/notes.
func (c *Client) CreateDeviceNote(ctx context.Context, deviceID, content string) (*DeviceNote, error) {
	path := apiPathPrefix + "/devices/" + url.PathEscape(deviceID) + "/notes"
	body := CreateDeviceNoteRequest{Content: content}
	raw, _ := json.Marshal(body)
	req, err := c.newRequest(ctx, http.MethodPost, path, bytes.NewReader(raw))
	if err != nil {
		return nil, fmt.Errorf("create device note: %w", err)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("create device note: %w", err)
	}
	defer resp.Body.Close()
	out, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("create device note read: %w", err)
	}
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("create device note: %s: %s", resp.Status, string(out))
	}
	var note DeviceNote
	if err := json.Unmarshal(out, &note); err != nil {
		return nil, fmt.Errorf("create device note decode: %w", err)
	}
	return &note, nil
}

// UpdateDeviceNote calls PATCH /api/v1/devices/{device_id}/notes/{note_id}.
func (c *Client) UpdateDeviceNote(ctx context.Context, deviceID, noteID, content string) (*DeviceNote, error) {
	path := apiPathPrefix + "/devices/" + url.PathEscape(deviceID) + "/notes/" + url.PathEscape(noteID)
	body := CreateDeviceNoteRequest{Content: content}
	raw, _ := json.Marshal(body)
	req, err := c.newRequest(ctx, http.MethodPatch, path, bytes.NewReader(raw))
	if err != nil {
		return nil, fmt.Errorf("update device note: %w", err)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("update device note: %w", err)
	}
	defer resp.Body.Close()
	out, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("update device note read: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("update device note: %s: %s", resp.Status, string(out))
	}
	var note DeviceNote
	if err := json.Unmarshal(out, &note); err != nil {
		return nil, fmt.Errorf("update device note decode: %w", err)
	}
	return &note, nil
}

// DeleteDeviceNote calls DELETE /api/v1/devices/{device_id}/notes/{note_id}.
func (c *Client) DeleteDeviceNote(ctx context.Context, deviceID, noteID string) error {
	path := apiPathPrefix + "/devices/" + url.PathEscape(deviceID) + "/notes/" + url.PathEscape(noteID)
	req, err := c.newRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return fmt.Errorf("delete device note: %w", err)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("delete device note: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNoContent || resp.StatusCode == http.StatusOK {
		return nil
	}
	body, _ := io.ReadAll(resp.Body)
	return fmt.Errorf("delete device note: %s: %s", resp.Status, string(body))
}

// GetDeviceJSON calls GET on a device subresource and returns raw JSON (activity, apps, library-items, parameters, status).
func (c *Client) GetDeviceJSON(ctx context.Context, deviceID, subpath string, query url.Values) ([]byte, error) {
	path := apiPathPrefix + "/devices/" + url.PathEscape(deviceID) + "/" + url.PathEscape(subpath)
	if q := query.Encode(); q != "" {
		path += "?" + q
	}
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

// GetDeviceLostMode calls GET /api/v1/devices/{device_id}/details/lostmode.
func (c *Client) GetDeviceLostMode(ctx context.Context, deviceID string) ([]byte, error) {
	path := apiPathPrefix + "/devices/" + url.PathEscape(deviceID) + "/details/lostmode"
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

// CancelDeviceLostMode calls DELETE /api/v1/devices/{device_id}/details/lostmode.
func (c *Client) CancelDeviceLostMode(ctx context.Context, deviceID string) error {
	path := apiPathPrefix + "/devices/" + url.PathEscape(deviceID) + "/details/lostmode"
	req, err := c.newRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}
	body, _ := io.ReadAll(resp.Body)
	return fmt.Errorf("%s: %s", resp.Status, string(body))
}

// UpdateDeviceRequest is the body for PATCH /api/v1/devices/{device_id}.
type UpdateDeviceRequest struct {
	User        interface{} `json:"user,omitempty"`         // UUID string or null to clear
	AssetTag    interface{} `json:"asset_tag,omitempty"`     // string or null to clear
	BlueprintID string      `json:"blueprint_id,omitempty"`
	Tags        []string    `json:"tags,omitempty"` // empty to clear all
}

// UpdateDevice calls PATCH /api/v1/devices/{device_id}.
func (c *Client) UpdateDevice(ctx context.Context, deviceID string, payload UpdateDeviceRequest) ([]byte, error) {
	path := apiPathPrefix + "/devices/" + url.PathEscape(deviceID)
	raw, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := c.newRequest(ctx, http.MethodPatch, path, bytes.NewReader(raw))
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
		return nil, fmt.Errorf("update device: %s: %s", resp.Status, string(body))
	}
	return body, nil
}

// GetDeviceSecret calls GET /api/v1/devices/{device_id}/secrets/{secret_type} (bypasscode, filevaultkey, unlockpin, recoverypassword).
func (c *Client) GetDeviceSecret(ctx context.Context, deviceID, secretType string) ([]byte, error) {
	path := apiPathPrefix + "/devices/" + url.PathEscape(deviceID) + "/secrets/" + url.PathEscape(secretType)
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
		return nil, fmt.Errorf("get device secret %s: %s: %s", secretType, resp.Status, string(body))
	}
	return body, nil
}
