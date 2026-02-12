package pgdb

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	TableUser               = "users"
	TableUserBalanceHistory = "user_balance_history"
)

var migrations = map[string]string{
	TableUser:               "../../schema/user.up.sql",
	TableUserBalanceHistory: "../../schema/user_balance_history.up.sql",
}

// NOTE: simplified migration mechanism intended for demo / test task.
func RunMigrations(ctx context.Context, pgdb *pgxpool.Pool, tables []string) error {
	for _, table := range tables {
		schemaFile, ok := migrations[table]
		if !ok {
			return fmt.Errorf("no migration found for table: %s", table)
		}

		schema, err := os.ReadFile(schemaFile)
		if err != nil {
			return fmt.Errorf("failed to read schema file %s: %w", schemaFile, err)
		}

		_, err = pgdb.Exec(ctx, string(schema))
		if err != nil {
			return fmt.Errorf("failed to execute migration for table %s: %w", table, err)
		}
	}

	return nil
}
