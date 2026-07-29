package repository

import (
    "context"

    "github.com/google/uuid"
    "github.com/LittleMardiah/scinel/internal/model"
)

// TransactionRepository defines transaction data access operations
type TransactionRepository interface {
    // Create inserts a new transaction record
    Create(ctx context.Context, txn *model.Transaction) (*model.Transaction, error)

    // GetByID retrieves single transaction
    GetByID(ctx context.Context, txnID uuid.UUID) (*model.Transaction, error)

    // Query retrieves transactions with filters and pagination
    Query(ctx context.Context, filter model.TransactionFilter, limit, offset int) ([]model.Transaction, int64, error)

    // GetByWarehouse retrieves all transactions in a warehouse
    GetByWarehouse(ctx context.Context, warehouseID uuid.UUID, limit, offset int) ([]model.Transaction, int64, error)

    // GetBySKU retrieves all transactions for a SKU
    GetBySKU(ctx context.Context, skuID uuid.UUID, limit, offset int) ([]model.Transaction, int64, error)
}