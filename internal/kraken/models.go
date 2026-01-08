package kraken

// APIResponse models the Kraken Ticker API response.
type APIResponse struct {
	Error  []string               `json:"error"`
	Result map[string]TickerEntry `json:"result"`
}

// TickerEntry captures the relevant parts of the ticker payload.
// Field "c" is "last trade closed": [ <price>, <lot volume> ]
type TickerEntry struct {
	C []string `json:"c"`
}
