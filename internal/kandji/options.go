package kandji

import (
	"net/url"
	"strconv"
)

// ListDeviceOptions holds query parameters for GET /api/v1/devices/.
// ExtraParams are merged into the query string; use any API query param names (e.g. serial_number).
type ListDeviceOptions struct {
	Limit        int
	Offset       int
	DeviceID     string
	DeviceName   string
	SerialNumber string
	MacAddress   string
	UserName     string
	UserEmail    string
	Platform     string
	BlueprintID  string
	// ExtraParams are arbitrary query params (JSON key → value). Merged after built-in flags; same key overrides.
	ExtraParams map[string]string
}

// QueryValues returns url.Values for the list devices request.
func (o ListDeviceOptions) QueryValues() url.Values {
	v := url.Values{}
	if o.Limit > 0 {
		v.Set("limit", strconv.Itoa(o.Limit))
	}
	if o.Offset > 0 {
		v.Set("offset", strconv.Itoa(o.Offset))
	}
	if o.DeviceID != "" {
		v.Set("device_id", o.DeviceID)
	}
	if o.DeviceName != "" {
		v.Set("device_name", o.DeviceName)
	}
	if o.SerialNumber != "" {
		v.Set("serial_number", o.SerialNumber)
	}
	if o.MacAddress != "" {
		v.Set("mac_address", o.MacAddress)
	}
	if o.UserName != "" {
		v.Set("user_name", o.UserName)
	}
	if o.UserEmail != "" {
		v.Set("user_email", o.UserEmail)
	}
	if o.Platform != "" {
		v.Set("platform", o.Platform)
	}
	if o.BlueprintID != "" {
		v.Set("blueprint_id", o.BlueprintID)
	}
	for k, val := range o.ExtraParams {
		if val != "" {
			v.Set(k, val)
		}
	}
	return v
}
