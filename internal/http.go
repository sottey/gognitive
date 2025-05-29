package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// DoRequest handles common request logic including headers, timeouts, decoding, and error reporting.
func DoRequest(client *http.Client, req *http.Request, apiKey string, result interface{}) error {
	if client == nil {
		client = &http.Client{Timeout: 10 * time.Second}
	}
	req.Header.Set("X-API-Key", apiKey)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API request failed: %s", resp.Status)
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	return nil
}
