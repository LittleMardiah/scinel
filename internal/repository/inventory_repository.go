package repository

import (
    "context"

    "github.com/google/uuid"
    "github.com/LittleMardiah/scinel/internal/model"
)

// InventoryRepository defines inventory data access operations
type InventoryRepository interface {
    // GetBySKU retrieves inventory for a specific warehouse-SKU combination
    // Returns: nil if not found (not an error - let caller decide action)
    GetBySKU(ctx context.Context, warehouseID, skuID uuid.UUID) (*model.Inventory, error)

    // GetByWarehouse retrieves all SKU inventory in a warehouse
    GetByWarehouse(ctx context.Context, warehouseID uuid.UUID, limit, offset int) ([]model.Inventory, int64, error)

    // GetAll retrieves all inventory across all warehouses
    GetAll(ctx context.Context, limit, offset int) ([]model.Inventory, int64, error)

    // Create creates a new inventory record
    Create(ctx context.Context, inventory *model.Inventory) (*model.Inventory, error)

    // UpdateQuantity atomically updates inventory quantity
    // Returns: Updated inventory or error
    // Note: Uses pessimistic locking (SELECT ... FOR UPDATE) for consistency
    UpdateQuantity(ctx context.Context, warehouseID, skuID uuid.UUID, delta int32) (*model.Inventory, error)

    // LockForUpdate locks inventory record for update in a transaction
    // Used for: Pessimistic locking in UPDATE inventory flow
    LockForUpdate(ctx context.Context, warehouseID, skuID uuid.UUID) (*model.Inventory, error)

    // GetBySKUCached retrieves inventory with Redis cache fallback (30s TTL)
    // On cache miss: Query DB, then cache result
    // On cache error: Log but continue (cache is optional)
    GetBySKUCached(ctx context.Context, warehouseID, skuID uuid.UUID) (*model.Inventory, error)

    // InvalidateCache removes inventory from cache after update
    InvalidateCache(ctx context.Context, warehouseID, skuID uuid.UUID) error

    // CreateSnapshot records point-in-time inventory level
    CreateSnapshot(ctx context.Context, snapshot *model.InventorySnapshot) error
}