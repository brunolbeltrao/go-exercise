package models

type Pair string

const (
	PairBTCUSD Pair = "BTC/USD"
	PairBTCCHF Pair = "BTC/CHF"
	PairBTCEUR Pair = "BTC/EUR"
)

// LTPEntry represents a single pair + amount in the response.
type LTPEntry struct {
	Pair   string  `json:"pair"`
	Amount float64 `json:"amount"`
}

// LTPResponse is the top-level API response payload.
type LTPResponse struct {
	LTP []LTPEntry `json:"ltp"`
}
