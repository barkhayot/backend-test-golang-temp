package skinport

import (
	"context"
	"fmt"
	"time"
)

type Config struct {
	fetchInterval time.Duration
	fetchOnLaunch bool
}

type Option func(*Config)

func WithFetchInterval(d int) Option {
	return func(c *Config) {
		c.fetchInterval = time.Duration(d) * time.Minute
	}
}

func WithFetchOnLaunch(fetchOnLaunch bool) Option {
	return func(c *Config) {
		c.fetchOnLaunch = fetchOnLaunch
	}
}

type Service struct {
	cfg   *Config
	state *State
}

func NewService(opts ...Option) *Service {
	cfg := &Config{
		fetchInterval: defaultFetchInterval,
	}

	for _, opt := range opts {
		opt(cfg)
	}

	return &Service{
		cfg:   cfg,
		state: newState(),
	}
}

func (s *Service) Run(ctx context.Context) {
	ticker := time.NewTicker(s.cfg.fetchInterval)
	defer ticker.Stop()

	if s.cfg.fetchOnLaunch {
		if err := s.runOnce(ctx); err != nil {
			fmt.Printf("initial fetch failed: %v\n", err)
		}
	}

	for {
		select {
		case <-ticker.C:
			if err := s.runOnce(ctx); err != nil {
				fmt.Printf("fetch failed: %v\n", err)
			}
		case <-ctx.Done():
			return
		}
	}
}

func (s *Service) GetItems() []Item {
	return s.state.getAll()
}
