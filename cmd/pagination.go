package cmd

import (
	"fmt"
	"net/url"
	"os"
	"strings"
)

// printPaginationHint writes a one-line hint to stderr when there is a next or previous page.
// nextURL and prevURL are the API response next/previous values (full URL or cursor).
// For cursor-based APIs (audit, users): pass the response next/prev; we extract cursor and show --cursor=.
// For offset-based (blueprints, devices): we show --offset= and --limit=.
func printPaginationHint(mode string, nextURL, prevURL *string, count int, limit, offset int) {
	if nextURL == nil && prevURL == nil && mode != "devices" {
		return
	}
	// For devices we only have count/limit/offset (no next URL); hint when we got a full page
	if mode == "devices" {
		if limit <= 0 || count < limit {
			return
		}
		nextOffset := offset + limit
		fmt.Fprintf(os.Stderr, "Next page: use --offset=%d --limit=%d\n", nextOffset, limit)
		return
	}

	var nextHint, prevHint string
	if nextURL != nil && *nextURL != "" {
		nextHint = extractPaginationFlags(mode, *nextURL)
	}
	if prevURL != nil && *prevURL != "" {
		prevHint = extractPaginationFlags(mode, *prevURL)
	}

	if nextHint != "" {
		fmt.Fprintf(os.Stderr, "Next page: %s\n", nextHint)
	}
	if prevHint != "" {
		fmt.Fprintf(os.Stderr, "Previous page: %s\n", prevHint)
	}
}

// extractPaginationFlags parses next/previous URL and returns a string like "--cursor=abc" or "--offset=300 --limit=100".
func extractPaginationFlags(mode, raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}
	// If it doesn't look like a URL, treat as cursor (audit/users).
	if !strings.HasPrefix(raw, "http://") && !strings.HasPrefix(raw, "https://") {
		if mode == "audit" || mode == "users" {
			return "use --cursor=" + raw
		}
		return ""
	}
	u, err := url.Parse(raw)
	if err != nil {
		return fmt.Sprintf("use --cursor=%s", raw)
	}
	q := u.Query()

	switch mode {
	case "audit", "users":
		if c := q.Get("cursor"); c != "" {
			return fmt.Sprintf("use --cursor=%s", c)
		}
		return ""
	case "blueprints":
		var parts []string
		if o := q.Get("offset"); o != "" {
			parts = append(parts, fmt.Sprintf("--offset=%s", o))
		}
		if l := q.Get("limit"); l != "" {
			parts = append(parts, fmt.Sprintf("--limit=%s", l))
		}
		if len(parts) > 0 {
			return "use " + strings.Join(parts, " ")
		}
		return ""
	default:
		return ""
	}
}
