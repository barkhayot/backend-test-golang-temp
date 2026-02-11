package skinport

import "testing"

func TestFetchTradable(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	resp, err := fetch(t.Context(), tradableParam)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(resp) == 0 {
		t.Fatal("expected non-empty response")
	}

	t.Logf("fetched %d tradable items", len(resp))
}

func TestFetchUntradable(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	resp, err := fetch(t.Context(), untradableParam)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(resp) == 0 {
		t.Fatal("expected non-empty response")
	}

	t.Logf("fetched %d untradable items", len(resp))
}

func TestFetchWithInvalidParam(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	_, err := fetch(t.Context(), "invalid_param")
	if err == nil {
		t.Fatal("expected error for invalid tradable param, got nil")
	}

	t.Logf("received expected error: %v", err)
}

func TestGetOrCreate_NewItem(t *testing.T) {
	m := make(map[string]*Item)

	r := raw{
		MarketHashName: "AK-47 | Redline",
		Currency:       "USD",
	}

	it := getOrCreate(m, r)

	if it == nil {
		t.Fatal("expected item, got nil")
	}

	if it.MarketHashName != r.MarketHashName {
		t.Errorf("market_hash_name mismatch")
	}

	if it.Currency != r.Currency {
		t.Errorf("currency mismatch")
	}

	if len(m) != 1 {
		t.Errorf("expected map size 1, got %d", len(m))
	}
}

func TestGetOrCreate_ExistingItem(t *testing.T) {
	m := make(map[string]*Item)

	r := raw{
		MarketHashName: "AK-47 | Redline",
		Currency:       "USD",
	}

	first := getOrCreate(m, r)
	second := getOrCreate(m, r)

	if first != second {
		t.Errorf("expected same pointer, got different items")
	}
}

func TestPostProcess_MergesTradableAndUntradable(t *testing.T) {
	tradablePrice := 10.0
	untradablePrice := 7.5

	tradable := []raw{
		{
			MarketHashName: "AK-47 | Redline",
			Currency:       "USD",
			MinPrice:       &tradablePrice,
		},
	}

	untradable := []raw{
		{
			MarketHashName: "AK-47 | Redline",
			Currency:       "USD",
			MinPrice:       &untradablePrice,
		},
	}

	items := postProcess(tradable, untradable)

	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}

	it := items[0]

	if it.MinTradablePrice == nil || *it.MinTradablePrice != tradablePrice {
		t.Errorf("invalid MinTradablePrice")
	}

	if it.MinUntradablePrice == nil || *it.MinUntradablePrice != untradablePrice {
		t.Errorf("invalid MinUntradablePrice")
	}
}

func TestPostProcess_ItemOnlyTradable(t *testing.T) {
	price := 12.0

	tradable := []raw{
		{
			MarketHashName: "M4A1-S | Basilisk",
			Currency:       "USD",
			MinPrice:       &price,
		},
	}

	items := postProcess(tradable, nil)

	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}

	it := items[0]

	if it.MinTradablePrice == nil {
		t.Errorf("expected MinTradablePrice to be set")
	}

	if it.MinUntradablePrice != nil {
		t.Errorf("expected MinUntradablePrice to be nil")
	}
}
