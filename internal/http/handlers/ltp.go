package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"go-exercise/internal/kraken"
	"go-exercise/internal/ltp"
	"go-exercise/internal/models"
)

type LTPHandler struct {
	svc ltp.Service
}

func NewLTPHandler(s ltp.Service) *LTPHandler {
	return &LTPHandler{svc: s}
}

func (h *LTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Extract pairs from query
	pairs := parsePairs(r)
	if len(pairs) == 0 {
		// default to all supported
		for _, p := range kraken.SupportedPairs {
			pairs = append(pairs, p)
		}
	}
	// Validate
	for _, p := range pairs {
		if !kraken.IsSupported(p) {
			http.Error(w, "unsupported pair: "+string(p), http.StatusBadRequest)
			return
		}
	}

	entries, err := h.svc.GetLTP(r.Context(), pairs)
	if err != nil {
		writeJSON(w, http.StatusBadGateway, models.LTPResponse{LTP: []models.LTPEntry{}})
		return
	}
	writeJSON(w, http.StatusOK, models.LTPResponse{LTP: entries})
}

func parsePairs(r *http.Request) []models.Pair {
	var out []models.Pair
	// repeated ?pair=
	for _, v := range r.URL.Query()["pair"] {
		v = strings.TrimSpace(v)
		if v == "" {
			continue
		}
		out = append(out, models.Pair(v))
	}
	// csv ?pairs=
	if csv := strings.TrimSpace(r.URL.Query().Get("pairs")); csv != "" {
		for _, v := range strings.Split(csv, ",") {
			v = strings.TrimSpace(v)
			if v == "" {
				continue
			}
			out = append(out, models.Pair(v))
		}
	}
	return out
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// urlError helps classify upstream errors; placeholder for future expansion
type urlError struct{}
