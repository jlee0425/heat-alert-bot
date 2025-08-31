package wbgt

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

const (
	HeatAlertDownloadURL = "https://www.wbgt.env.go.jp/alert/dl"
	tokyoTimezone        = "Asia/Tokyo"
)

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL:    baseURL,
		httpClient: &http.Client{},
	}
}

func GetTokyoTime() time.Time {
	// Get current UTC time
	utcTime := time.Now().UTC()

	// Load Tokyo timezone
	timezone, err := time.LoadLocation(tokyoTimezone)
	if err != nil {
		// Fallback to fixed offset if timezone loading fails
		timezone = time.FixedZone("JST", 9*60*60) // UTC+9
	}

	// Convert to Tokyo local time
	return utcTime.In(timezone)
}

func GetAlertEndpoint() string {
	// target string format: "2025/alert_20250831_05.csv"
	// Format: YYYY/alert_YYYYMMDD_05.csv
	// https://pkg.go.dev/time#example-Time.Format
	tokyoTime := GetTokyoTime()
	year := tokyoTime.Format("2006")
	dateStr := tokyoTime.Format("20060102")

	return fmt.Sprintf("%s/alert_%s_05.csv", year, dateStr)
}

func (c *Client) FetchCSVData(ctx context.Context, endpoint string) ([]byte, error) {
	baseURL, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}

	fullURL, err := url.JoinPath(baseURL.String(), endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to construct URL: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	fmt.Println("[GET] fetching today's alert data from", fullURL)
	req.Header.Set("Accept", "text/csv")
	req.Header.Set("User-Agent", "heat-alert-bot/1.0")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}

	return body, nil
}
