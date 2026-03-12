package cmd

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// parseExtraParams parses a JSON object into map[string]string for query params.
// Empty or invalid JSON returns nil. Used by any list/GET command that supports --params.
// Accepts standard JSON (double quotes). If that fails, tries single-quoted style so
// --params "{'archived':'true'}" works (converted to {"archived":"true"}).
func parseExtraParams(raw string) map[string]string {
	if raw == "" {
		return nil
	}
	raw = strings.TrimSpace(raw)
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(raw), &m); err != nil {
		// Common mistake: single-quoted JSON (invalid). Try converting ' to ".
		try := strings.ReplaceAll(raw, "'", "\"")
		if err2 := json.Unmarshal([]byte(try), &m); err2 != nil {
			return nil
		}
	}
	out := make(map[string]string, len(m))
	for k, v := range m {
		if v == nil {
			continue
		}
		switch t := v.(type) {
		case string:
			out[k] = t
		case float64:
			out[k] = strconv.FormatFloat(t, 'f', -1, 64)
		case bool:
			out[k] = strconv.FormatBool(t)
		default:
			out[k] = fmt.Sprint(t)
		}
	}
	return out
}
