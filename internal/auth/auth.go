package auth

import (
	"errors"
	"net/http"
)

// getAPIKeyFromHeader gets the API Key from the header
// Example: X-API-KEY: <API_KEY>
func GetAPIKeyFromHeader(header http.Header) (string, error) {
	apiKey := header.Get("X-API-KEY")
	if apiKey == "" {
		return "", errors.New("no API key provided")
	}
	return apiKey, nil
}
