package skinport

import "sync"

type State struct {
	items []Item

	// assuming that reads are more frequent than writes
	mu sync.RWMutex
}

func newState() *State {
	return &State{
		items: make([]Item, 0),
	}
}

func (s *State) GetAll() []Item {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make([]Item, len(s.items))
	copy(out, s.items)
	return out
}

func (s *State) set(items []Item) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items = items
}
