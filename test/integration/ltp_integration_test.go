package integration

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"go-exercise/internal/cache"
	httpapi "go-exercise/internal/http"
	"go-exercise/internal/http/handlers"
	"go-exercise/internal/kraken"
	"go-exercise/internal/ltp"
	"go-exercise/internal/models"
)

type mockKraken struct {
	symbolToPrice map[string]float64
}

func (m *mockKraken) GetLastTradedPrices(_ context.Context, krakenSymbols []string) (map[string]float64, error) {
	out := make(map[string]float64, len(krakenSymbols))
	for _, s := range krakenSymbols {
		if v, ok := m.symbolToPrice[s]; ok {
			out[s] = v
		}
	}
	return out, nil
}

func TestLTP_DefaultAllPairs(t *testing.T) {
	c := cache.NewMemoryCache(60 * time.Second)
	k := &mockKraken{symbolToPrice: map[string]float64{
		"XBTUSD": 52000.12,
		"XBTEUR": 50000.12,
		"XBTCHF": 49000.12,
	}}
	s := ltp.NewService(c, k)
	ltpHandler := handlers.NewLTPHandler(s)
	server := httptest.NewServer(httpapi.NewRouter(ltpHandler))
	defer server.Close()

	resp, err := http.Get(server.URL + "/api/v1/ltp")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status: %d", resp.StatusCode)
	}
	var body models.LTPResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if len(body.LTP) != 3 {
		t.Fatalf("expected 3 items, got %d", len(body.LTP))
	}
	// Basic sanity: amounts > 0
	for _, e := range body.LTP {
		if e.Amount <= 0 {
			t.Fatalf("expected positive amount for %s", e.Pair)
		}
	}
}

func TestLTP_SelectPairs(t *testing.T) {
	c := cache.NewMemoryCache(60 * time.Second)
	k := &mockKraken{symbolToPrice: map[string]float64{
		"XBTUSD": 52000.12,
		"XBTEUR": 50000.12,
		"XBTCHF": 49000.12,
	}}
	s := ltp.NewService(c, k)
	ltpHandler := handlers.NewLTPHandler(s)
	server := httptest.NewServer(httpapi.NewRouter(ltpHandler))
	defer server.Close()

	resp, err := http.Get(server.URL + "/api/v1/ltp?pairs=BTC/USD,BTC/EUR")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status: %d", resp.StatusCode)
	}
	var body models.LTPResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if len(body.LTP) != 2 {
		t.Fatalf("expected 2 items, got %d", len(body.LTP))
	}
	pairs := map[string]bool{}
	for _, e := range body.LTP {
		pairs[e.Pair] = true
	}
	if !pairs[string(models.PairBTCUSD)] || !pairs[string(models.PairBTCEUR)] {
		t.Fatalf("unexpected pairs in response: %+v", pairs)
	}
}

func TestLTP_LiveKrakenOptional(t *testing.T) {
	if os.Getenv("TEST_LIVE_KRAKEN") != "1" {
		t.Skip("live test disabled")
	}
	c := cache.NewMemoryCache(60 * time.Second)
	k := kraken.NewRealClient(5 * time.Second)
	s := ltp.NewService(c, k)
	ltpHandler := handlers.NewLTPHandler(s)
	server := httptest.NewServer(httpapi.NewRouter(ltpHandler))
	defer server.Close()

	resp, err := http.Get(server.URL + "/api/v1/ltp?pair=BTC/USD")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status: %d", resp.StatusCode)
	}
	var body models.LTPResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if len(body.LTP) != 1 {
		t.Fatalf("expected 1 item, got %d", len(body.LTP))
	}
	if body.LTP[0].Amount <= 0 {
		t.Fatalf("expected positive amount")
	}
}
