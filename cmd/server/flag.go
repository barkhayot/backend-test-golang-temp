package main

import (
	"fmt"
	"os"
	"strconv"
)

type config struct {
	FetchInterval int // in minutes
	FetchOnLaunch bool
	RunMigrations bool
	AppPort       int
	Pgdb          string
}

func setDefaultConfig() config {
	return config{
		FetchInterval: 5,
		FetchOnLaunch: true,
		RunMigrations: true,
		AppPort:       8080,
		Pgdb:          "postgres://user:password@localhost:5432/postgres?sslmode=disable",
	}
}

func load() (config, error) {
	cfg := setDefaultConfig()

	if v := os.Getenv("FETCH_INTERVAL_IN_MINUTE"); v != "" {
		i, err := strconv.Atoi(v)
		if err != nil {
			return cfg, err
		}
		cfg.FetchInterval = i

	}

	if v := os.Getenv("FETCH_ON_LAUNCH"); v != "" {
		b, err := strconv.ParseBool(v)
		if err != nil {
			return cfg, err
		}
		cfg.FetchOnLaunch = b
	}

	if v := os.Getenv("APP_PORT"); v != "" {
		p, err := strconv.Atoi(v)
		if err != nil {
			return cfg, err
		}

		cfg.AppPort = p
	}

	if cfg.AppPort <= 0 || cfg.AppPort > 65535 {
		return cfg, fmt.Errorf("invalid APP_PORT: %d", cfg.AppPort)
	}

	if v := os.Getenv("PGDB"); v != "" {
		cfg.Pgdb = v
	}

	if v := os.Getenv("RUN_MIGRATIONS"); v != "" {
		b, err := strconv.ParseBool(v)
		if err != nil {
			return cfg, err
		}
		cfg.RunMigrations = b
	}

	return cfg, nil
}
