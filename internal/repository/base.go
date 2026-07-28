package repository

import (
"context"
"database/sql"

"github.com/LittleMardiah/scinel/internal/database"
)

// Repository is the base interface for all repositories
type Repository interface {
// Lifecycle methods
Close() error
Ping(ctx context.Context) error
WithTx(ctx context.Context, fn database.TxFunc) error
}

// BaseRepository provides common repository functionality
type BaseRepository struct {
db *database.DB
}

// NewBaseRepository creates a new base repository
func NewBaseRepository(db *database.DB) *BaseRepository {
return &BaseRepository{db: db}
}

// Close closes the database connection
func (r *BaseRepository) Close() error {
return r.db.Close()
}

// Ping tests the database connection
func (r *BaseRepository) Ping(ctx context.Context) error {
return r.db.PingContext(ctx)
}

// WithTx executes a function within a transaction
func (r *BaseRepository) WithTx(ctx context.Context, fn database.TxFunc) error {
return r.db.WithTx(ctx, fn)
}

// Helper method to get the database connection
func (r *BaseRepository) getDB() *sql.DB {
return r.db.DB
}
