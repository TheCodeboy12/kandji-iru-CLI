package kandji

import "encoding/json"

// Device matches the official API response for GET /api/v1/devices and
// GET /api/v1/devices/{device_id}. See Iru Endpoint Management API Postman collection.
type Device struct {
	DeviceID                   string   `json:"device_id"`
	DeviceName                 string   `json:"device_name"`
	Model                      string   `json:"model,omitempty"`
	SerialNumber               string   `json:"serial_number,omitempty"`
	Platform                   string   `json:"platform,omitempty"`
	OSVersion                  string   `json:"os_version,omitempty"`
	SupplementalBuildVersion   string   `json:"supplemental_build_version,omitempty"`
	SupplementalOSVersionExtra string   `json:"supplemental_os_version_extra,omitempty"`
	LastCheckIn                string   `json:"last_check_in,omitempty"`
	User                       *User    `json:"user,omitempty"`
	AssetTag                   string   `json:"asset_tag,omitempty"`
	BlueprintID                string   `json:"blueprint_id,omitempty"`
	MDMEnabled                 bool     `json:"mdm_enabled,omitempty"`
	AgentInstalled             bool     `json:"agent_installed,omitempty"`
	IsMissing                  bool     `json:"is_missing,omitempty"`
	IsRemoved                  bool     `json:"is_removed,omitempty"`
	AgentVersion               string   `json:"agent_version,omitempty"`
	FirstEnrollment            string   `json:"first_enrollment,omitempty"`
	LastEnrollment             string   `json:"last_enrollment,omitempty"`
	BlueprintName              string   `json:"blueprint_name,omitempty"`
	LostModeStatus             string   `json:"lost_mode_status,omitempty"`
	Tags                       []string `json:"tags,omitempty"`
}

// deviceRaw is used for decoding so "user" can be string or object.
type deviceRaw struct {
	DeviceID                   string   `json:"device_id"`
	DeviceName                 string   `json:"device_name"`
	Model                      string   `json:"model,omitempty"`
	SerialNumber               string   `json:"serial_number,omitempty"`
	Platform                   string   `json:"platform,omitempty"`
	OSVersion                  string   `json:"os_version,omitempty"`
	SupplementalBuildVersion   string   `json:"supplemental_build_version,omitempty"`
	SupplementalOSVersionExtra string   `json:"supplemental_os_version_extra,omitempty"`
	LastCheckIn                string   `json:"last_check_in,omitempty"`
	User                       userFlex `json:"user,omitempty"`
	AssetTag                   string   `json:"asset_tag,omitempty"`
	BlueprintID                string   `json:"blueprint_id,omitempty"`
	MDMEnabled                 bool     `json:"mdm_enabled,omitempty"`
	AgentInstalled             bool     `json:"agent_installed,omitempty"`
	IsMissing                  bool     `json:"is_missing,omitempty"`
	IsRemoved                  bool     `json:"is_removed,omitempty"`
	AgentVersion               string   `json:"agent_version,omitempty"`
	FirstEnrollment            string   `json:"first_enrollment,omitempty"`
	LastEnrollment             string   `json:"last_enrollment,omitempty"`
	BlueprintName              string   `json:"blueprint_name,omitempty"`
	LostModeStatus             string   `json:"lost_mode_status,omitempty"`
	Tags                       []string `json:"tags,omitempty"`
}

func (d *Device) UnmarshalJSON(data []byte) error {
	var raw deviceRaw
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	d.DeviceID = raw.DeviceID
	d.DeviceName = raw.DeviceName
	d.Model = raw.Model
	d.SerialNumber = raw.SerialNumber
	d.Platform = raw.Platform
	d.OSVersion = raw.OSVersion
	d.SupplementalBuildVersion = raw.SupplementalBuildVersion
	d.SupplementalOSVersionExtra = raw.SupplementalOSVersionExtra
	d.LastCheckIn = raw.LastCheckIn
	d.User = raw.User.User
	d.AssetTag = raw.AssetTag
	d.BlueprintID = raw.BlueprintID
	d.MDMEnabled = raw.MDMEnabled
	d.AgentInstalled = raw.AgentInstalled
	d.IsMissing = raw.IsMissing
	d.IsRemoved = raw.IsRemoved
	d.AgentVersion = raw.AgentVersion
	d.FirstEnrollment = raw.FirstEnrollment
	d.LastEnrollment = raw.LastEnrollment
	d.BlueprintName = raw.BlueprintName
	d.LostModeStatus = raw.LostModeStatus
	d.Tags = raw.Tags
	return nil
}

// userFlex accepts user as either a string or an object (API can return both).
type userFlex struct {
	*User
}

func (u *userFlex) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || string(data) == "null" {
		u.User = nil
		return nil
	}
	if data[0] == '"' {
		var s string
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}
		if s == "" {
			u.User = nil
		} else {
			u.User = &User{Name: s}
		}
		return nil
	}
	u.User = &User{}
	return json.Unmarshal(data, u.User)
}

// User is the nested user object on Device (official API: email, name, id, is_archived).
type User struct {
	Email      string `json:"email,omitempty"`
	Name       string `json:"name,omitempty"`
	ID         string `json:"id,omitempty"`
	IsArchived bool   `json:"is_archived,omitempty"`
}

// DeviceSummary is a row for table output.
type DeviceSummary struct {
	DeviceID     string
	DeviceName   string
	SerialNumber string
	Platform     string
	UserEmail    string
	UserName     string
}

// Summary returns a DeviceSummary for table output.
func (d *Device) Summary() DeviceSummary {
	var s DeviceSummary
	s.DeviceID = d.DeviceID
	s.DeviceName = d.DeviceName
	s.SerialNumber = d.SerialNumber
	s.Platform = d.Platform
	if d.User != nil {
		s.UserEmail = d.User.Email
		s.UserName = d.User.Name
	}
	return s
}
