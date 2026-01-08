package kraken

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// Client defines the behavior needed by the service to fetch LTPs.
type Client interface {
	// GetLastTradedPrices returns a map of requested Kraken symbols to their last traded price.
	GetLastTradedPrices(ctx context.Context, krakenSymbols []string) (map[string]float64, error)
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type RealClient struct {
	BaseURL    string
	HTTPClient HTTPClient
}

func NewRealClient(httpTimeout time.Duration) *RealClient {
	return &RealClient{
		BaseURL: "https://api.kraken.com",
		HTTPClient: &http.Client{
			Timeout: httpTimeout,
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
				DialContext: (&net.Dialer{
					Timeout:   5 * time.Second,
					KeepAlive: 30 * time.Second,
				}).DialContext,
				MaxIdleConns:        100,
				IdleConnTimeout:     90 * time.Second,
				TLSHandshakeTimeout: 5 * time.Second,
			},
		},
	}
}

func (c *RealClient) GetLastTradedPrices(ctx context.Context, krakenSymbols []string) (map[string]float64, error) {
	if len(krakenSymbols) == 0 {
		return map[string]float64{}, nil
	}

	// Build query string, requesting all symbols at once
	joined := strings.Join(krakenSymbols, ",")
	u, err := url.Parse(c.BaseURL)
	if err != nil {
		return nil, err
	}
	u.Path = "/0/public/Ticker"
	q := u.Query()
	q.Set("pair", joined)
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return nil, fmt.Errorf("kraken: status %d: %s", resp.StatusCode, string(body))
	}

	var parsed APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return nil, err
	}
	if len(parsed.Error) > 0 {
		return nil, errors.New(strings.Join(parsed.Error, "; "))
	}

	// Map requested symbol -> price using alt keys to account for Kraken's response naming
	out := make(map[string]float64, len(krakenSymbols))
	for _, sym := range krakenSymbols {
		keys := krakenAltKeys[sym]
		var priceStr string
		for _, k := range keys {
			if entry, ok := parsed.Result[k]; ok {
				if len(entry.C) > 0 {
					priceStr = entry.C[0]
					break
				}
			}
		}
		if priceStr == "" {
			// try direct
			if entry, ok := parsed.Result[sym]; ok && len(entry.C) > 0 {
				priceStr = entry.C[0]
			}
		}
		if priceStr == "" {
			return nil, fmt.Errorf("kraken: missing price for %s", sym)
		}
		// Parse as float64
		val, err := parseFloat(priceStr)
		if err != nil {
			return nil, fmt.Errorf("kraken: invalid price for %s: %w", sym, err)
		}
		out[sym] = val
	}
	return out, nil
}

// parseFloat handles Kraken's price string format safely.
func parseFloat(s string) (float64, error) {
	// Avoid locale issues; Kraken returns dot-decimal.
	return strconvParseFloat(s)
}

// tiny indirection for easy testing
var strconvParseFloat = func(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}
