package skinport

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/barkhayot/request/pkg/request"
)

const (
	// https://docs.skinport.com/items
	baseURL      = "https://api.skinport.com/v1/items"
	fetchTimeout = 10 * time.Second // can be decreased if needed

	// rate limit: 8 requests per 5 minutes
	// should be > 37.5s to respect rate limit
	// it can be descreased from 1 minute to 40s safely
	defaultFetchInterval = 1 * time.Minute
)

var (
	tradableParam   = "1"
	untradableParam = "0"
)

type Item struct {
	MarketHashName     string   `json:"market_hash_name"`
	Version            *string  `json:"version"`
	Currency           string   `json:"currency"`
	SuggestedPrice     float64  `json:"suggested_price"`
	ItemPage           string   `json:"item_page"`
	MarketPage         string   `json:"market_page"`
	Quantity           int      `json:"quantity"`
	CreatedAt          int64    `json:"created_at"`
	UpdatedAt          int64    `json:"updated_at"`
	MinTradablePrice   *float64 `json:"min_tradable_price"`
	MaxTradablePrice   *float64 `json:"max_tradable_price"`
	MinUntradablePrice *float64 `json:"min_untradable_price"`
	MaxUntradablePrice *float64 `json:"max_untradable_price"`
}

type raw struct {
	MarketHashName string   `json:"market_hash_name"`
	Version        *string  `json:"version"`
	Currency       string   `json:"currency"`
	SuggestedPrice float64  `json:"suggested_price"`
	ItemPage       string   `json:"item_page"`
	MarketPage     string   `json:"market_page"`
	MinPrice       *float64 `json:"min_price"`
	MaxPrice       *float64 `json:"max_price"`
	MeanPrice      *float64 `json:"mean_price"`
	MedianPrice    *float64 `json:"median_price"`
	Quantity       int      `json:"quantity"`
	CreatedAt      int64    `json:"created_at"`
	UpdatedAt      int64    `json:"updated_at"`
}

func (s *Service) runOnce(parent context.Context) error {
	ctx, cancel := context.WithTimeout(parent, fetchTimeout)
	defer cancel()

	// NOTE: if any of the fetches fail, we don't update the state
	tradable, err := fetch(ctx, tradableParam)
	if err != nil {
		return fmt.Errorf("tradable fetch: %w", err)
	}

	untradable, err := fetch(ctx, untradableParam)
	if err != nil {
		return fmt.Errorf("untradable fetch: %w", err)
	}

	items := postProcess(tradable, untradable)
	s.state.set(items)

	return nil
}

func fetch(ctx context.Context, tradable string) ([]raw, error) {
	if tradable != tradableParam && tradable != untradableParam {
		return nil, errors.New("invalid tradable param")
	}

	headers := http.Header{}
	headers.Add("Accept-Encoding", "br")

	params := url.Values{}
	params.Add("tradable", tradable)

	resp, err := request.Request[[]raw](
		ctx,
		request.WithMethod("GET"),
		request.WithEndpoint(baseURL),
		request.WithHeaders(headers),
		request.WithQueryParams(params),
		request.WithTimeout(fetchTimeout),
	)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func postProcess(tradable, untradable []raw) []Item {
	m := make(map[string]*Item)

	for _, r := range tradable {
		it := getOrCreate(m, r)
		it.MinTradablePrice = r.MinPrice
	}

	for _, r := range untradable {
		it := getOrCreate(m, r)
		it.MinUntradablePrice = r.MinPrice
	}

	out := make([]Item, 0, len(m))
	for _, v := range m {
		out = append(out, *v)
	}

	return out
}

func getOrCreate(m map[string]*Item, r raw) *Item {
	if it, ok := m[r.MarketHashName]; ok {
		return it
	}

	it := &Item{
		MarketHashName: r.MarketHashName,
		Currency:       r.Currency,
	}
	m[r.MarketHashName] = it
	return it
}
