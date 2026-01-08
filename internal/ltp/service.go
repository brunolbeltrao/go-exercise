package ltp

import (
	"context"
	"fmt"

	"go-exercise/internal/cache"
	"go-exercise/internal/kraken"
	"go-exercise/internal/models"
)

type Service interface {
	GetLTP(ctx context.Context, pairs []models.Pair) ([]models.LTPEntry, error)
}

type service struct {
	cache        *cache.MemoryCache
	krakenClient kraken.Client
}

func NewService(c *cache.MemoryCache, k kraken.Client) Service {
	return &service{
		cache:        c,
		krakenClient: k,
	}
}

func (s *service) GetLTP(ctx context.Context, pairs []models.Pair) ([]models.LTPEntry, error) {
	// Build list of pairs to fetch from cache or Kraken
	type pairState struct {
		value float64
		hit   bool
	}
	state := make(map[models.Pair]*pairState, len(pairs))
	missing := make([]models.Pair, 0, len(pairs))

	for _, p := range pairs {
		if _, ok := state[p]; ok {
			continue // de-dup
		}
		if v, ok := s.cache.Get(p); ok {
			state[p] = &pairState{value: v, hit: true}
			continue
		}
		state[p] = &pairState{}
		missing = append(missing, p)
	}

	// Fetch missing pairs from Kraken (batch by kraken symbol)
	if len(missing) > 0 {
		krakenSymbolsSet := make(map[string]struct{}, len(missing))
		// map from krakenSymbol -> []models.Pair (app pairs that use it)
		symbolToPairs := make(map[string][]models.Pair, len(missing))
		for _, p := range missing {
			sym, _, ok := kraken.MapToKraken(p)
			if !ok {
				return nil, fmt.Errorf("unsupported pair: %s", p)
			}
			krakenSymbolsSet[sym] = struct{}{}
			symbolToPairs[sym] = append(symbolToPairs[sym], p)
		}
		krakenSymbols := make([]string, 0, len(krakenSymbolsSet))
		for s := range krakenSymbolsSet {
			krakenSymbols = append(krakenSymbols, s)
		}
		prices, err := s.krakenClient.GetLastTradedPrices(ctx, krakenSymbols)
		if err != nil {
			return nil, err
		}
		for sym, pairsForSym := range symbolToPairs {
			price, ok := prices[sym]
			if !ok {
				return nil, fmt.Errorf("kraken: missing price for symbol %s", sym)
			}
			for _, p := range pairsForSym {
				state[p].value = price
				state[p].hit = true
				s.cache.Set(p, price)
			}
		}
	}

	// Build response list in the order requested (including duplicates requested)
	out := make([]models.LTPEntry, 0, len(pairs))
	for _, p := range pairs {
		ps := state[p]
		if ps == nil || !ps.hit {
			return nil, fmt.Errorf("failed to obtain price for %s", p)
		}
		out = append(out, models.LTPEntry{
			Pair:   string(p),
			Amount: ps.value,
		})
	}
	return out, nil
}
