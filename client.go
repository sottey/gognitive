package gognitive

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/sottey/gognitive/internal"
)

// Client wraps the API key and HTTP client for Limitless API access.
type Client struct {
	APIKey     string
	HTTPClient *http.Client
}

// NewClient creates a new API client with a default timeout.
func NewClient(apiKey string) *Client {
	return &Client{
		APIKey:     apiKey,
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
	}
}

// ListLifelogs retrieves a list of lifelogs with optional filters.
// Supports pagination via cursor.

func (c *Client) ListLifelogs(limit int, cursor, date, start, end, timezone string) ([]Lifelog, string, error) {
	baseURL := "https://api.limitless.ai/v1/lifelogs"
	params := url.Values{}
	if limit > 0 {
		params.Set("limit", fmt.Sprintf("%d", limit))
	}
	if cursor != "" {
		params.Set("cursor", cursor)
	}
	// ... add other params as needed ...

	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, "", err
	}

	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, "", fmt.Errorf("API request failed: %s", resp.Status)
	}

	var result struct {
		Items      []Lifelog `json:"items"`
		NextCursor string    `json:"nextCursor"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, "", err
	}

	return result.Items, result.NextCursor, nil
}

/*
func (c *Client) ListLifelogs(limit int, cursor, date, start, end, timezone string) ([]Lifelog, string, error) {
	baseURL := "https://api.limitless.ai/v1/lifelogs"
	params := url.Values{}
	if limit > 0 {
		params.Set("limit", fmt.Sprintf("%d", limit))
	}
	if cursor != "" {
		params.Set("cursor", cursor)
	}
	if date != "" {
		params.Set("date", date)
	}
	if start != "" {
		params.Set("start", start)
	}
	if end != "" {
		params.Set("end", end)
	}
	if timezone != "" {
		params.Set("timezone", timezone)
	}

	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, "", err
	}

	var result struct {
		Data struct {
			Lifelogs []Lifelog `json:"lifelogs"`
		} `json:"data"`
		Meta struct {
			Lifelogs struct {
				NextCursor string `json:"nextCursor"`
			} `json:"lifelogs"`
		} `json:"meta"`
	}

	if err := internal.DoRequest(c.HTTPClient, req, c.APIKey, &result); err != nil {
		return nil, "", err
	}

	return result.Data.Lifelogs, result.Meta.Lifelogs.NextCursor, nil
}
*/

// GetLifelog retrieves a single lifelog by ID.
func (c *Client) GetLifelog(id string) (*Lifelog, error) {
	reqURL := fmt.Sprintf("https://api.limitless.ai/v1/lifelogs?id=%s", url.QueryEscape(id))
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Data struct {
			Lifelogs []Lifelog `json:"lifelogs"`
		} `json:"data"`
	}

	if err := internal.DoRequest(c.HTTPClient, req, c.APIKey, &result); err != nil {
		return nil, err
	}

	if len(result.Data.Lifelogs) == 0 {
		return nil, fmt.Errorf("no lifelog found for ID %s", id)
	}

	return &result.Data.Lifelogs[0], nil
}

// GetEnrichedLifelog retrieves a lifelog and enriches it with generated tags
func (c *Client) GetEnrichedLifelog(id string) (*EnrichedLifelog, error) {
	log, err := c.GetLifelog(id)
	if err != nil {
		return nil, err
	}

	tags := internal.GenerateTags(log.Markdown)

	return &EnrichedLifelog{
		ID:        log.ID,
		Title:     log.Title,
		StartTime: log.StartTime,
		EndTime:   log.EndTime,
		Markdown:  log.Markdown,
		Contents:  log.Contents,
		Tags:      tags,
	}, nil
}
