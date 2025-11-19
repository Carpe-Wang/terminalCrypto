package keyring

import (
	"fmt"

	"github.com/zalando/go-keyring"
)

const serviceName = "terminalcrypto"

// Credentials holds API credentials
type Credentials struct {
	APIKey    string
	APISecret string
}

// StoreCredentials stores API credentials securely in the system keyring
func StoreCredentials(exchange, apiKey, apiSecret string) error {
	// Store API key
	if err := keyring.Set(serviceName+"-"+exchange, "api-key", apiKey); err != nil {
		return fmt.Errorf("failed to store API key: %w", err)
	}

	// Store API secret
	if err := keyring.Set(serviceName+"-"+exchange, "api-secret", apiSecret); err != nil {
		return fmt.Errorf("failed to store API secret: %w", err)
	}

	return nil
}

// GetCredentials retrieves API credentials from the system keyring
func GetCredentials(exchange string) (*Credentials, error) {
	// Get API key
	apiKey, err := keyring.Get(serviceName+"-"+exchange, "api-key")
	if err != nil {
		return nil, fmt.Errorf("failed to get API key for %s: %w (have you run 'setup %s'?)", exchange, err, exchange)
	}

	// Get API secret
	apiSecret, err := keyring.Get(serviceName+"-"+exchange, "api-secret")
	if err != nil {
		return nil, fmt.Errorf("failed to get API secret for %s: %w", exchange, err)
	}

	return &Credentials{
		APIKey:    apiKey,
		APISecret: apiSecret,
	}, nil
}

// DeleteCredentials removes API credentials from the system keyring
func DeleteCredentials(exchange string) error {
	// Delete API key
	if err := keyring.Delete(serviceName+"-"+exchange, "api-key"); err != nil {
		return fmt.Errorf("failed to delete API key: %w", err)
	}

	// Delete API secret
	if err := keyring.Delete(serviceName+"-"+exchange, "api-secret"); err != nil {
		return fmt.Errorf("failed to delete API secret: %w", err)
	}

	return nil
}

// HasCredentials checks if credentials exist for an exchange
func HasCredentials(exchange string) bool {
	_, err := GetCredentials(exchange)
	return err == nil
}
