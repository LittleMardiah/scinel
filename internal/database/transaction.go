package database

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

// TxFunc is a function that runs within a transaction
type TxFunc func(context.Context, *sql.Tx) error

// WithTx executes a function within a transaction with deadlock retry logic
func (db *DB) WithTx(ctx context.Context, fn TxFunc) error {
	const maxRetries = 3
	var lastErr error

	for attempt := 0; attempt < maxRetries; attempt++ {
		tx, err := db.BeginTx(ctx, &sql.TxOptions{
			Isolation: sql.LevelReadCommitted,
		})
		if err != nil {
			return fmt.Errorf("failed to begin transaction: %w", err)
		}

		if err = fn(ctx, tx); err != nil {
			tx.Rollback()
			if isDeadlock(err) && attempt < maxRetries-1 {
				lastErr = err
				backoff := time.Duration((1 << uint(attempt)) * 100) * time.Millisecond
				time.Sleep(backoff)
				continue
			}
			return fmt.Errorf("transaction error (attempt %d): %w", attempt+1, err)
		}

		if err = tx.Commit(); err != nil {
			if isDeadlock(err) && attempt < maxRetries-1 {
				lastErr = err
				backoff := time.Duration((1 << uint(attempt)) * 100) * time.Millisecond
				time.Sleep(backoff)
				continue
			}
			return fmt.Errorf("failed to commit transaction: %w", err)
		}
		return nil
	}

	return fmt.Errorf("transaction failed after %d retries (last error: %w)", maxRetries, lastErr)
}

// isDeadlock checks if error is a deadlock error
func isDeadlock(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	return strings.Contains(errStr, "deadlock detected") ||
		strings.Contains(errStr, "40P01") ||
		strings.HasSuffix(errStr, "1")
}
