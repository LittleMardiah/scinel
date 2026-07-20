package database

import (
"database/sql"
"fmt"
"time"

_ "github.com/lib/pq"
)

type DB struct {
*sql.DB
}

// ConnectionConfig holds database configuration
type ConnectionConfig struct {
URL                string
MinConnections     int
MaxConnections     int
ConnectTimeout     time.Duration
IdleTimeout        time.Duration
MaxLifetime        time.Duration
IsolationLevel     sql.IsolationLevel
}

// NewConnection creates a new database connection pool
func NewConnection(config ConnectionConfig) (*DB, error) {
// Set defaults if not provided
if config.MinConnections == 0 {
config.MinConnections = 10
}
if config.MaxConnections == 0 {
config.MaxConnections = 50
}
if config.ConnectTimeout == 0 {
config.ConnectTimeout = 5 * time.Second
}
if config.IdleTimeout == 0 {
config.IdleTimeout = 30 * time.Second
}
if config.MaxLifetime == 0 {
config.MaxLifetime = 5 * time.Minute
}

// Open database connection
sqlDB, err := sql.Open("postgres", config.URL)
if err != nil {
return nil, fmt.Errorf("failed to open database: %w", err)
}

// Set connection pool parameters
sqlDB.SetMaxOpenConns(config.MaxConnections)
sqlDB.SetMaxIdleConns(config.MinConnections)
sqlDB.SetConnMaxIdleTime(config.IdleTimeout)
sqlDB.SetConnMaxLifetime(config.MaxLifetime)

// Test connection
if err = sqlDB.Ping(); err != nil {
return nil, fmt.Errorf("failed to ping database: %w", err)
}

return &DB{sqlDB}, nil
}

// Close closes the database connection
func (db *DB) Close() error {
return db.DB.Close()
}
