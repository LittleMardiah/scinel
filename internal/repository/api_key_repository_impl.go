package repository

import (
    "context"
    "errors"

    "github.com/google/uuid"
    "github.com/LittleMardiah/scinel/internal/database"
    "github.com/LittleMardiah/scinel/internal/model"
)

type apiKeyRepository struct {
    db *database.DB
}

// NewAPIKeyRepository creates a new API key repository instance
func NewAPIKeyRepository(db *database.DB) APIKeyRepository {
    return &apiKeyRepository{db: db}
}

// --- PLACEHOLDER IMPLEMENTATIONS ---
func (r *apiKeyRepository) Create(ctx context.Context, key *model.APIKey) (*model.APIKey, error) {
    return nil, errors.New("Create API key not implemented yet")
}

func (r *apiKeyRepository) GetByID(ctx context.Context, keyID, userID uuid.UUID) (*model.APIKey, error) {
    return nil, errors.New("GetByID not implemented yet")
}

func (r *apiKeyRepository) GetByPrefix(ctx context.Context, prefix string) (*model.APIKey, error) {
    return nil, errors.New("GetByPrefix not implemented yet")
}

func (r *apiKeyRepository) GetByUser(ctx context.Context, userID uuid.UUID) ([]model.APIKey, error) {
    return nil, errors.New("GetByUser not implemented yet")
}

func (r *apiKeyRepository) Update(ctx context.Context, key *model.APIKey) error {
    return errors.New("Update not implemented yet")
}

func (r *apiKeyRepository) Delete(ctx context.Context, keyID, userID uuid.UUID) error {
    return errors.New("Delete not implemented yet")
}