package skinport

import (
	"fmt"
	"sync"
	"testing"
)

func TestNewState_Empty(t *testing.T) {
	s := newState()

	items := s.getAll()
	if len(items) != 0 {
		t.Fatalf("expected empty state, got %d items", len(items))
	}
}

func TestState_SetAndgetAll(t *testing.T) {
	s := newState()

	items := []Item{
		{MarketHashName: "AK-47"},
		{MarketHashName: "M4A1-S"},
	}

	s.set(items)

	got := s.getAll()

	if len(got) != len(items) {
		t.Fatalf("expected %d items, got %d", len(items), len(got))
	}

	if got[0].MarketHashName != items[0].MarketHashName {
		t.Errorf("unexpected item value")
	}
}

func TestState_getAllReturnsCopy(t *testing.T) {
	s := newState()

	items := []Item{
		{MarketHashName: "AK-47"},
	}
	s.set(items)

	got := s.getAll()

	// mutate returned slice
	got[0].MarketHashName = "HACKED"

	again := s.getAll()
	if again[0].MarketHashName == "HACKED" {
		t.Fatalf("internal state was mutated via getAll")
	}
}

func TestState_ConcurrentReads(t *testing.T) {
	s := newState()
	s.set([]Item{
		{MarketHashName: "AK-47"},
	})

	var wg sync.WaitGroup

	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = s.getAll()
		}()
	}

	wg.Wait()
}

func TestState_ReadWhileWrite(t *testing.T) {
	s := newState()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 100; i++ {
			s.set([]Item{
				{MarketHashName: fmt.Sprintf("item-%d", i)},
			})
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 100; i++ {
			_ = s.getAll()
		}
	}()

	wg.Wait()
}
