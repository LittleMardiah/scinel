package repository

import (
    "context"
    "fmt"

    "github.com/google/uuid"
    "github.com/lib/pq"
    "github.com/LittleMardiah/scinel/internal/database"
    "github.com/LittleMardiah/scinel/internal/model"
)

type apiKeyRepository struct {
    db *database.DB
}

func NewAPIKeyRepository(db *database.DB) APIKeyRepository {
    return &apiKeyRepository{db: db}
}

// GetByUser retrieves all API keys for a specific user
func (r *apiKeyRepository) GetByUser(ctx context.Context, userID uuid.UUID) ([]model.APIKey, error) {
    query := `
        SELECT id, user_id, key_prefix, key_hash, name, description,
               permissions, is_active, last_used_at, expires_at,
               created_at, updated_at
        FROM api_keys
        WHERE user_id = $1
        ORDER BY created_at DESC
    `
    rows, err := r.db.QueryContext(ctx, query, userID)
    if err != nil {
        return nil, fmt.Errorf("failed to query API keys: %w", err)
    }
    defer rows.Close()

    var keys []model.APIKey
    for rows.Next() {
        var k model.APIKey
        err := rows.Scan(
            &k.ID, &k.UserID, &k.KeyPrefix, &k.KeyHash, &k.Name, &k.Description,
            pq.Array(&k.Permissions), &k.IsActive, &k.LastUsedAt, &k.ExpiresAt,
            &k.CreatedAt, &k.UpdatedAt,
        )
        if err != nil {
            return nil, fmt.Errorf("failed to scan API key: %w", err)
        }
        keys = append(keys, k)
    }
    if err = rows.Err(); err != nil {
        return nil, fmt.Errorf("rows iteration error: %w", err)
    }
    return keys, nil
}

// Create inserts a new API key
func (r *apiKeyRepository) Create(ctx context.Context, key *model.APIKey) (*model.APIKey, error) {
    query := `
        INSERT INTO api_keys (id, user_id, key_prefix, key_hash, name, description,
                              permissions, is_active, expires_at, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
        RETURNING id, user_id, key_prefix, key_hash, name, description,
                  permissions, is_active, last_used_at, expires_at,
                  created_at, updated_at
    `
    row := r.db.QueryRowContext(ctx, query,
        key.ID, key.UserID, key.KeyPrefix, key.KeyHash, key.Name, key.Description,
        pq.Array(key.Permissions), key.IsActive, key.ExpiresAt, key.CreatedAt, key.UpdatedAt,
    )
    err := row.Scan(
        &key.ID, &key.UserID, &key.KeyPrefix, &key.KeyHash, &key.Name, &key.Description,
        pq.Array(&key.Permissions), &key.IsActive, &key.LastUsedAt, &key.ExpiresAt,
        &key.CreatedAt, &key.UpdatedAt,
    )
    if err != nil {
        return nil, fmt.Errorf("failed to create API key: %w", err)
    }
    return key, nil
}

// GetByID retrieves API key by ID (with ownership check)
func (r *apiKeyRepository) GetByID(ctx context.Context, keyID, userID uuid.UUID) (*model.APIKey, error) {
    query := `
        SELECT id, user_id, key_prefix, key_hash, name, description,
               permissions, is_active, last_used_at, expires_at,
               created_at, updated_at
        FROM api_keys
        WHERE id = $1 AND user_id = $2
    `
    var k model.APIKey
    err := r.db.QueryRowContext(ctx, query, keyID, userID).Scan(
        &k.ID, &k.UserID, &k.KeyPrefix, &k.KeyHash, &k.Name, &k.Description,
        pq.Array(&k.Permissions), &k.IsActive, &k.LastUsedAt, &k.ExpiresAt,
        &k.CreatedAt, &k.UpdatedAt,
    )
    if err != nil {
        return nil, fmt.Errorf("failed to get API key: %w", err)
    }
    return &k, nil
}

// GetByPrefix retrieves API key by prefix (for validation)
func (r *apiKeyRepository) GetByPrefix(ctx context.Context, prefix string) (*model.APIKey, error) {
    query := `
        SELECT id, user_id, key_prefix, key_hash, name, description,
               permissions, is_active, last_used_at, expires_at,
               created_at, updated_at
        FROM api_keys
        WHERE key_prefix = $1
    `
    var k model.APIKey
    err := r.db.QueryRowContext(ctx, query, prefix).Scan(
        &k.ID, &k.UserID, &k.KeyPrefix, &k.KeyHash, &k.Name, &k.Description,
        pq.Array(&k.Permissions), &k.IsActive, &k.LastUsedAt, &k.ExpiresAt,
        &k.CreatedAt, &k.UpdatedAt,
    )
    if err != nil {
        return nil, fmt.Errorf("failed to get API key by prefix: %w", err)
    }
    return &k, nil
}

// UpdateLastUsed updates the last_used_at timestamp
func (r *apiKeyRepository) UpdateLastUsed(ctx context.Context, keyID uuid.UUID) error {
    query := `UPDATE api_keys SET last_used_at = NOW() WHERE id = $1`
    _, err := r.db.ExecContext(ctx, query, keyID)
    return err
}

// Delete revokes an API key (soft delete by setting is_active = false)
func (r *apiKeyRepository) Delete(ctx context.Context, keyID, userID uuid.UUID) error {
    query := `UPDATE api_keys SET is_active = false, updated_at = NOW() WHERE id = $1 AND user_id = $2`
    result, err := r.db.ExecContext(ctx, query, keyID, userID)
    if err != nil {
        return fmt.Errorf("failed to revoke API key: %w", err)
    }
    rows, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to get rows affected: %w", err)
    }
    if rows == 0 {
        return fmt.Errorf("API key not found or not owned by user")
    }
    return nil
}

// Rotate updates the key_hash for a key (rotate secret)
func (r *apiKeyRepository) Rotate(ctx context.Context, keyID, userID uuid.UUID, newKeyHash string) error {
    query := `UPDATE api_keys SET key_hash = $1, updated_at = NOW() WHERE id = $2 AND user_id = $3`
    result, err := r.db.ExecContext(ctx, query, newKeyHash, keyID, userID)
    if err != nil {
        return fmt.Errorf("failed to rotate API key: %w", err)
    }
    rows, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to get rows affected: %w", err)
    }
    if rows == 0 {
        return fmt.Errorf("API key not found or not owned by user")
    }
    return nil
}

// Update updates an existing API key
func (r *apiKeyRepository) Update(ctx context.Context, key *model.APIKey) error {
    query := `
        UPDATE api_keys
        SET name = $1, description = $2, permissions = $3,
            is_active = $4, expires_at = $5, updated_at = NOW()
        WHERE id = $6 AND user_id = $7
    `
    result, err := r.db.ExecContext(ctx, query,
        key.Name, key.Description, pq.Array(key.Permissions),
        key.IsActive, key.ExpiresAt,
        key.ID, key.UserID,
    )
    if err != nil {
        return fmt.Errorf("failed to update API key: %w", err)
    }
    rows, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to get rows affected: %w", err)
    }
    if rows == 0 {
        return fmt.Errorf("API key not found or not owned by user")
    }
    return nil
}
