package model

import (
    "time"

    "github.com/google/uuid"
)

// APIKey represents an API key for programmatic access
type APIKey struct {
    ID          uuid.UUID  `db:"id"`
    UserID      uuid.UUID  `db:"user_id"`
    KeyPrefix   string     `db:"key_prefix"`   // e.g., "sk_test_1a2b3c"
    KeyHash     string     `db:"key_hash"`     // bcrypt hash of full key
    Name        string     `db:"name"`
    Description *string    `db:"description"`
    Permissions []string   `db:"permissions"`  // e.g., ["inventory:read", "forecast:write"]
    IsActive    bool       `db:"is_active"`
    LastUsedAt  *time.Time `db:"last_used_at"`
    ExpiresAt   *time.Time `db:"expires_at"`
    CreatedAt   time.Time  `db:"created_at"`
    UpdatedAt   time.Time  `db:"updated_at"`
}