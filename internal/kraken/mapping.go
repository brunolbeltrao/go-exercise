package kraken

import "go-exercise/internal/models"

// SupportedPairs is the whitelist of pairs accepted by the API.
var SupportedPairs = []models.Pair{
	models.PairBTCUSD,
	models.PairBTCCHF,
	models.PairBTCEUR,
}

// appToKraken maps API pairs to Kraken pair symbols for querying.
// Note: Kraken commonly uses XBT for Bitcoin (not BTC).
var appToKraken = map[models.Pair]string{
	models.PairBTCUSD: "XBTUSD",
	models.PairBTCCHF: "XBTCHF",
	models.PairBTCEUR: "XBTEUR",
}

// krakenAltKeys lists possible response keys for a requested Kraken pair,
// as Kraken often returns alternative pair names in the response.
var krakenAltKeys = map[string][]string{
	"XBTUSD": {"XBTUSD", "XXBTZUSD"},
	"XBTCHF": {"XBTCHF", "XXBTZCHF"},
	"XBTEUR": {"XBTEUR", "XXBTZEUR"},
}

// IsSupported returns true if the given pair is supported.
func IsSupported(p models.Pair) bool {
	_, ok := appToKraken[p]
	return ok
}

// MapToKraken returns the Kraken symbol for a given API pair and a set of possible response keys.
func MapToKraken(p models.Pair) (symbol string, possibleResponseKeys []string, ok bool) {
	s, ok := appToKraken[p]
	if !ok {
		return "", nil, false
	}
	return s, krakenAltKeys[s], true
}

// AllKrakenSymbols returns the list of Kraken symbols for the supported pairs, in deterministic order.
func AllKrakenSymbols() []string {
	return []string{"XBTUSD", "XBTCHF", "XBTEUR"}
}
