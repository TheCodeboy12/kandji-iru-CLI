// Package keyring provides OS keyring storage for the Kandji API token and base URL.
package keyring

import (
	"github.com/zalando/go-keyring"
)

const (
	// Service is the keyring service name for kandji-iru-cli.
	Service = "kandji-iru-cli"
	// User is the keyring account name used for the API token.
	User = "api-token"
	// UserBaseURL is the keyring account name for the API base URL.
	UserBaseURL = "base-url"
	// UserSubdomain is the keyring account name for the tenant subdomain.
	UserSubdomain = "subdomain"
)

// GetToken returns the Kandji API token from the system keyring, or empty string if not found or on error.
func GetToken() (string, error) {
	return keyring.Get(Service, User)
}

// SetToken stores the Kandji API token in the system keyring.
func SetToken(token string) error {
	return keyring.Set(Service, User, token)
}

// DeleteToken removes the Kandji API token from the system keyring.
func DeleteToken() error {
	return keyring.Delete(Service, User)
}

// GetBaseURL returns the Kandji API base URL from the system keyring, or empty string if not found or on error.
func GetBaseURL() (string, error) {
	return keyring.Get(Service, UserBaseURL)
}

// SetBaseURL stores the Kandji API base URL in the system keyring.
func SetBaseURL(baseURL string) error {
	return keyring.Set(Service, UserBaseURL, baseURL)
}

// DeleteBaseURL removes the base URL from the system keyring.
func DeleteBaseURL() error {
	return keyring.Delete(Service, UserBaseURL)
}

// GetSubdomain returns the tenant subdomain from the system keyring, or empty string if not found or on error.
func GetSubdomain() (string, error) {
	return keyring.Get(Service, UserSubdomain)
}

// SetSubdomain stores the tenant subdomain in the system keyring.
func SetSubdomain(subdomain string) error {
	return keyring.Set(Service, UserSubdomain, subdomain)
}

// DeleteSubdomain removes the subdomain from the system keyring.
func DeleteSubdomain() error {
	return keyring.Delete(Service, UserSubdomain)
}
