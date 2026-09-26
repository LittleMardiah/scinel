package repository

import (
    "context"
    "database/sql"
    "errors"

    "github.com/google/uuid"
    "github.com/LittleMardiah/scinel/internal/model"
)

// Sentinel errors returned by the inventory repository. They are exported so
// the handler layer can map them onto HTTP status codes after the error has
// travelled through database.WithTx (which wraps with %w).
var (
    // ErrInventoryNotFound means no inventory row exists for the
    // warehouse/SKU pair.
    ErrInventoryNotFound = errors.New("inventory record not found")
    // ErrInsufficientStock means the requested delta would drive the quantity
    // below zero.
    ErrInsufficientStock = errors.New("insufficient stock")
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
    //
    // Deprecated: LockForUpdate is non-tx-aware; gunakan LockForUpdateTx dalam
    // WithTx. Non-tx tidak memegang lock. Outside an explicit transaction
    // PostgreSQL runs the statement in its own implicit transaction, so the row
    // lock is released the instant the SELECT returns and provides no mutual
    // exclusion at all.
    LockForUpdate(ctx context.Context, warehouseID, skuID uuid.UUID) (*model.Inventory, error)

    // LockForUpdateTx locks the inventory row inside a caller-owned
    // transaction. The lock is held until that transaction commits or rolls
    // back, which is what makes read-modify-write sequences safe.
    // Returns: ErrInventoryNotFound when the row does not exist.
    LockForUpdateTx(ctx context.Context, tx *sql.Tx, warehouseID, skuID uuid.UUID) (*model.Inventory, error)

    // UpdateQuantityTx applies delta to the inventory row inside a
    // caller-owned transaction and returns the updated level.
    // The caller MUST hold the row lock (LockForUpdateTx) first, otherwise the
    // read-modify-write is not atomic.
    // Returns: ErrInsufficientStock when the result would be negative.
    UpdateQuantityTx(ctx context.Context, tx *sql.Tx, warehouseID, skuID uuid.UUID, delta int32) (*model.Inventory, error)

    // GetBySKUCached retrieves inventory with Redis cache fallback (30s TTL)
    // On cache miss: Query DB, then cache result
    // On cache error: Log but continue (cache is optional)
    GetBySKUCached(ctx context.Context, warehouseID, skuID uuid.UUID) (*model.Inventory, error)

    // InvalidateCache removes inventory from cache after update
    InvalidateCache(ctx context.Context, warehouseID, skuID uuid.UUID) error

    // CreateSnapshot records point-in-time inventory level
    CreateSnapshot(ctx context.Context, snapshot *model.InventorySnapshot) error
}
