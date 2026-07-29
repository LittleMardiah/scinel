package repository

import (
    "context"

    "github.com/google/uuid"
    "github.com/LittleMardiah/scinel/internal/model"
)

// APIKeyRepository defines API key data access operations
type APIKeyRepository interface {
    // Create inserts a new API key
    Create(ctx context.Context, key *model.APIKey) (*model.APIKey, error)

    // GetByID retrieves API key by ID (with ownership check)
    GetByID(ctx context.Context, keyID, userID uuid.UUID) (*model.APIKey, error)

    // GetByPrefix retrieves API key by prefix (for validation)
    GetByPrefix(ctx context.Context, prefix string) (*model.APIKey, error)

    // GetByUser retrieves all API keys for a user
    GetByUser(ctx context.Context, userID uuid.UUID) ([]model.APIKey, error)

    // Update updates API key (name, permissions, is_active, last_used)
    Update(ctx context.Context, key *model.APIKey) error

    // Delete revokes API key (soft delete or hard delete)
    Delete(ctx context.Context, keyID, userID uuid.UUID) error
}