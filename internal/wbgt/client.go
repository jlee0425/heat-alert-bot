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
	TokyoTimezone        = "Asia/Tokyo"
)

func getTokyoTime() time.Time {
	utcTime := time.Now().UTC()
	timezone, err := time.LoadLocation(TokyoTimezone)
	if err != nil {
		timezone = time.FixedZone("JST", 9*60*60) // UTC+9 fallback
	}
	return utcTime.In(timezone)
}

func GetAlertEndpoint() string {
	tokyoTime := getTokyoTime()
	year := tokyoTime.Format("2006")
	dateStr := tokyoTime.Format("20060102")
	return fmt.Sprintf("%s/alert_%s_05.csv", year, dateStr)
}

type HeaderRoundTripper struct {
	base    http.RoundTripper
	headers map[string]string
}

func (h *HeaderRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	// Clone the request to avoid modifying the original
	newReq := req.Clone(req.Context())

	// Add default headers
	for key, value := range h.headers {
		if newReq.Header.Get(key) == "" { // Only set if not already present
			newReq.Header.Set(key, value)
		}
	}

	return h.base.RoundTrip(newReq)
}

func NewClient(baseURL string) *Client {
	// Default headers for all requests
	defaultHeaders := map[string]string{
		"Accept":          "text/csv",
		"Accept-Language": "ja,en-US",
		"Accept-Charset":  "UTF-8",
		"User-Agent":      "heat-alert-bot/1.0",
	}

	transport := &HeaderRoundTripper{
		base:    http.DefaultTransport,
		headers: defaultHeaders,
	}

	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Transport: transport,
			Timeout:   30 * time.Second,
		},
	}
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

	fmt.Printf("[GET] fetching today's alert data from\n > %s\n", fullURL)

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
