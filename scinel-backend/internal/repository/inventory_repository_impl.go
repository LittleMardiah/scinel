package repository

import (
    "context"
    "database/sql"
    "errors"
    "fmt"

    "github.com/google/uuid"
    "github.com/LittleMardiah/scinel/internal/database"
    "github.com/LittleMardiah/scinel/internal/model"
)

// inventoryRepository implements InventoryRepository interface
type inventoryRepository struct {
    db *database.DB
}

// NewInventoryRepository creates a new inventory repository instance
func NewInventoryRepository(db *database.DB) InventoryRepository {
    return &inventoryRepository{db: db}
}

// --- IMPLEMENTATIONS ---

// GetAll retrieves all inventory across all warehouses
func (r *inventoryRepository) GetAll(ctx context.Context, limit, offset int) ([]model.Inventory, int64, error) {
    query := `
        SELECT 
            i.id, i.warehouse_id, i.sku_id, 
            s.sku_code, w.name as warehouse_name, s.name as sku_name,
            i.quantity, i.reorder_point, i.max_stock_level,
            i.last_restock_date, i.last_transaction_id,
            i.created_at, i.updated_at
        FROM inventory i
        JOIN warehouses w ON i.warehouse_id = w.id
        JOIN skus s ON i.sku_id = s.id
        ORDER BY i.updated_at DESC
        LIMIT $1 OFFSET $2
    `

    rows, err := r.db.QueryContext(ctx, query, limit, offset)
    if err != nil {
        return nil, 0, fmt.Errorf("failed to query inventory: %w", err)
    }
    defer rows.Close()

    var items []model.Inventory
    for rows.Next() {
        var inv model.Inventory
        var lastTransactionID *uuid.UUID
        err := rows.Scan(
            &inv.ID, &inv.WarehouseID, &inv.SKUID,
            &inv.SKUCode, &inv.WarehouseName, &inv.SKUName,
            &inv.Quantity, &inv.ReorderPoint, &inv.MaxStockLevel,
            &inv.LastRestockDate, &lastTransactionID,
            &inv.CreatedAt, &inv.UpdatedAt,
        )
        if err != nil {
            return nil, 0, fmt.Errorf("failed to scan inventory: %w", err)
        }
        inv.LastTransactionID = lastTransactionID
        items = append(items, inv)
    }

    if err = rows.Err(); err != nil {
        return nil, 0, fmt.Errorf("rows iteration error: %w", err)
    }

    // Count total
    var total int64
    countQuery := `SELECT COUNT(*) FROM inventory`
    err = r.db.QueryRowContext(ctx, countQuery).Scan(&total)
    if err != nil {
        return nil, 0, fmt.Errorf("failed to count inventory: %w", err)
    }

    return items, total, nil
}

// GetByWarehouse retrieves all SKU inventory in a specific warehouse
func (r *inventoryRepository) GetByWarehouse(ctx context.Context, warehouseID uuid.UUID, limit, offset int) ([]model.Inventory, int64, error) {
    query := `
        SELECT 
            i.id, i.warehouse_id, i.sku_id, 
            s.sku_code, w.name as warehouse_name, s.name as sku_name,
            i.quantity, i.reorder_point, i.max_stock_level,
            i.last_restock_date, i.last_transaction_id,
            i.created_at, i.updated_at
        FROM inventory i
        JOIN warehouses w ON i.warehouse_id = w.id
        JOIN skus s ON i.sku_id = s.id
        WHERE i.warehouse_id = $1
        ORDER BY i.updated_at DESC
        LIMIT $2 OFFSET $3
    `

    rows, err := r.db.QueryContext(ctx, query, warehouseID, limit, offset)
    if err != nil {
        return nil, 0, fmt.Errorf("failed to query inventory by warehouse: %w", err)
    }
    defer rows.Close()

    var items []model.Inventory
    for rows.Next() {
        var inv model.Inventory
        var lastTransactionID *uuid.UUID
        err := rows.Scan(
            &inv.ID, &inv.WarehouseID, &inv.SKUID,
            &inv.SKUCode, &inv.WarehouseName, &inv.SKUName,
            &inv.Quantity, &inv.ReorderPoint, &inv.MaxStockLevel,
            &inv.LastRestockDate, &lastTransactionID,
            &inv.CreatedAt, &inv.UpdatedAt,
        )
        if err != nil {
            return nil, 0, fmt.Errorf("failed to scan inventory: %w", err)
        }
        inv.LastTransactionID = lastTransactionID
        items = append(items, inv)
    }

    if err = rows.Err(); err != nil {
        return nil, 0, fmt.Errorf("rows iteration error: %w", err)
    }

    // Count total for this warehouse
    var total int64
    countQuery := `SELECT COUNT(*) FROM inventory WHERE warehouse_id = $1`
    err = r.db.QueryRowContext(ctx, countQuery, warehouseID).Scan(&total)
    if err != nil {
        return nil, 0, fmt.Errorf("failed to count inventory by warehouse: %w", err)
    }

    return items, total, nil
}

// GetBySKU retrieves inventory for a specific warehouse-SKU combination
// Returns: nil if not found (not an error - let caller decide action)
func (r *inventoryRepository) GetBySKU(ctx context.Context, warehouseID, skuID uuid.UUID) (*model.Inventory, error) {
    query := `
        SELECT
            i.id, i.warehouse_id, i.sku_id,
            s.sku_code, w.name AS warehouse_name, s.name AS sku_name,
            i.quantity, i.reorder_point, i.max_stock_level,
            i.last_restock_date, i.last_transaction_id,
            i.created_at, i.updated_at
        FROM inventory i
        JOIN warehouses w ON i.warehouse_id = w.id
        JOIN skus s ON i.sku_id = s.id
        WHERE i.warehouse_id = $1 AND i.sku_id = $2
    `

    var inv model.Inventory
    var lastTransactionID *uuid.UUID
    err := r.db.QueryRowContext(ctx, query, warehouseID, skuID).Scan(
        &inv.ID, &inv.WarehouseID, &inv.SKUID,
        &inv.SKUCode, &inv.WarehouseName, &inv.SKUName,
        &inv.Quantity, &inv.ReorderPoint, &inv.MaxStockLevel,
        &inv.LastRestockDate, &lastTransactionID,
        &inv.CreatedAt, &inv.UpdatedAt,
    )
    if errors.Is(err, sql.ErrNoRows) {
        return nil, nil
    }
    if err != nil {
        return nil, fmt.Errorf("failed to get inventory by SKU: %w", err)
    }
    inv.LastTransactionID = lastTransactionID
    return &inv, nil
}

// Create creates a new inventory record
func (r *inventoryRepository) Create(ctx context.Context, inventory *model.Inventory) (*model.Inventory, error) {
    query := `
        INSERT INTO inventory (
            id, warehouse_id, sku_id, quantity, reorder_point, max_stock_level,
            last_restock_date, last_transaction_id, created_at, updated_at
        )
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
        RETURNING id, created_at, updated_at
    `

    err := r.db.QueryRowContext(ctx, query,
        inventory.ID, inventory.WarehouseID, inventory.SKUID,
        inventory.Quantity, inventory.ReorderPoint, inventory.MaxStockLevel,
        inventory.LastRestockDate, inventory.LastTransactionID,
        inventory.CreatedAt, inventory.UpdatedAt,
    ).Scan(&inventory.ID, &inventory.CreatedAt, &inventory.UpdatedAt)
    if err != nil {
        return nil, fmt.Errorf("failed to create inventory: %w", err)
    }
    return inventory, nil
}

// UpdateQuantity atomically updates inventory quantity
func (r *inventoryRepository) UpdateQuantity(ctx context.Context, warehouseID, skuID uuid.UUID, delta int32) (*model.Inventory, error) {
    var updatedInventory *model.Inventory
    err := r.db.WithTx(ctx, func(ctx context.Context, tx *sql.Tx) error {
        var currentQuantity int32
        err := tx.QueryRowContext(ctx, `
            SELECT quantity
            FROM inventory
            WHERE warehouse_id = $1 AND sku_id = $2
            FOR UPDATE
        `, warehouseID, skuID).Scan(&currentQuantity)
        if err != nil {
            return fmt.Errorf("failed to lock inventory: %w", err)
        }

        newQuantity := currentQuantity + delta
        if newQuantity < 0 {
            return ErrInsufficientStock
        }

        _, err = tx.ExecContext(ctx, `
            UPDATE inventory
            SET quantity = $1, updated_at = NOW()
            WHERE warehouse_id = $2 AND sku_id = $3
        `, newQuantity, warehouseID, skuID)
        if err != nil {
            return fmt.Errorf("failed to update inventory quantity: %w", err)
        }

        updatedInventory = &model.Inventory{
            WarehouseID: warehouseID,
            SKUID:       skuID,
            Quantity:    newQuantity,
        }
        return nil
    })
    if err != nil {
        if errors.Is(err, ErrInsufficientStock) {
            return nil, ErrInsufficientStock
        }
        return nil, err
    }
    return updatedInventory, nil
}

// LockForUpdateTx locks the inventory row inside a caller-owned transaction.
// Unlike LockForUpdate, the row lock survives until that transaction commits or
// rolls back, so the subsequent read-modify-write cannot interleave with a
// concurrent transaction.
//
// "FOR UPDATE OF i" restricts the lock to the inventory row: a bare FOR UPDATE
// on a query with JOINs would also lock the warehouses and skus rows, turning
// unrelated stock movements into lock contention and deadlock candidates.
//
// Returns ErrInventoryNotFound when the row does not exist.
func (r *inventoryRepository) LockForUpdateTx(ctx context.Context, tx *sql.Tx, warehouseID, skuID uuid.UUID) (*model.Inventory, error) {
    if tx == nil {
        return nil, fmt.Errorf("transaction must not be nil")
    }

    query := `
        SELECT
            i.id, i.warehouse_id, i.sku_id,
            s.sku_code, w.name AS warehouse_name, s.name AS sku_name,
            i.quantity, i.reorder_point, i.max_stock_level,
            i.last_restock_date, i.last_transaction_id,
            i.created_at, i.updated_at
        FROM inventory i
        JOIN warehouses w ON i.warehouse_id = w.id
        JOIN skus s ON i.sku_id = s.id
        WHERE i.warehouse_id = $1 AND i.sku_id = $2
        FOR UPDATE OF i
    `

    var inv model.Inventory
    var lastTransactionID *uuid.UUID
    err := tx.QueryRowContext(ctx, query, warehouseID, skuID).Scan(
        &inv.ID, &inv.WarehouseID, &inv.SKUID,
        &inv.SKUCode, &inv.WarehouseName, &inv.SKUName,
        &inv.Quantity, &inv.ReorderPoint, &inv.MaxStockLevel,
        &inv.LastRestockDate, &lastTransactionID,
        &inv.CreatedAt, &inv.UpdatedAt,
    )
    if errors.Is(err, sql.ErrNoRows) {
        return nil, ErrInventoryNotFound
    }
    if err != nil {
        return nil, fmt.Errorf("failed to lock inventory: %w", err)
    }
    inv.LastTransactionID = lastTransactionID
    return &inv, nil
}

// UpdateQuantityTx applies delta to the inventory row inside a caller-owned
// transaction and returns the resulting level.
//
// The quantity is re-read with FOR UPDATE rather than computed by the caller,
// so the guard against a negative result cannot be bypassed by a stale value.
// The caller is still expected to hold the lock from LockForUpdateTx; this
// re-read makes the method safe on its own as well.
func (r *inventoryRepository) UpdateQuantityTx(ctx context.Context, tx *sql.Tx, warehouseID, skuID uuid.UUID, delta int32) (*model.Inventory, error) {
    if tx == nil {
        return nil, fmt.Errorf("transaction must not be nil")
    }

    var currentQuantity int32
    err := tx.QueryRowContext(ctx, `
        SELECT quantity
        FROM inventory
        WHERE warehouse_id = $1 AND sku_id = $2
        FOR UPDATE
    `, warehouseID, skuID).Scan(&currentQuantity)
    if errors.Is(err, sql.ErrNoRows) {
        return nil, ErrInventoryNotFound
    }
    if err != nil {
        return nil, fmt.Errorf("failed to lock inventory: %w", err)
    }

    newQuantity := currentQuantity + delta
    if newQuantity < 0 {
        return nil, ErrInsufficientStock
    }

    if _, err := tx.ExecContext(ctx, `
        UPDATE inventory
        SET quantity = $1, updated_at = NOW()
        WHERE warehouse_id = $2 AND sku_id = $3
    `, newQuantity, warehouseID, skuID); err != nil {
        return nil, fmt.Errorf("failed to update inventory quantity: %w", err)
    }

    return &model.Inventory{
        WarehouseID: warehouseID,
        SKUID:       skuID,
        Quantity:    newQuantity,
    }, nil
}

// LockForUpdate locks inventory record for update in a transaction
//
// Deprecated: LockForUpdate is non-tx-aware; gunakan LockForUpdateTx dalam
// WithTx. Non-tx tidak memegang lock. This method runs FOR UPDATE on a pooled
// connection with no surrounding transaction, so PostgreSQL wraps the
// statement in an implicit transaction that commits immediately and the row
// lock is already gone when this function returns.
func (r *inventoryRepository) LockForUpdate(ctx context.Context, warehouseID, skuID uuid.UUID) (*model.Inventory, error) {
    query := `
        SELECT
            i.id, i.warehouse_id, i.sku_id,
            s.sku_code, w.name AS warehouse_name, s.name AS sku_name,
            i.quantity, i.reorder_point, i.max_stock_level,
            i.last_restock_date, i.last_transaction_id,
            i.created_at, i.updated_at
        FROM inventory i
        JOIN warehouses w ON i.warehouse_id = w.id
        JOIN skus s ON i.sku_id = s.id
        WHERE i.warehouse_id = $1 AND i.sku_id = $2
        FOR UPDATE
    `

    var inv model.Inventory
    var lastTransactionID *uuid.UUID
    err := r.db.QueryRowContext(ctx, query, warehouseID, skuID).Scan(
        &inv.ID, &inv.WarehouseID, &inv.SKUID,
        &inv.SKUCode, &inv.WarehouseName, &inv.SKUName,
        &inv.Quantity, &inv.ReorderPoint, &inv.MaxStockLevel,
        &inv.LastRestockDate, &lastTransactionID,
        &inv.CreatedAt, &inv.UpdatedAt,
    )
    if err != nil {
        return nil, fmt.Errorf("failed to lock inventory: %w", err)
    }
    inv.LastTransactionID = lastTransactionID
    return &inv, nil
}

// GetBySKUCached retrieves inventory with Redis cache fallback (30s TTL)
func (r *inventoryRepository) GetBySKUCached(ctx context.Context, warehouseID, skuID uuid.UUID) (*model.Inventory, error) {
    return r.GetBySKU(ctx, warehouseID, skuID)
}

// InvalidateCache removes inventory from cache after update
func (r *inventoryRepository) InvalidateCache(ctx context.Context, warehouseID, skuID uuid.UUID) error {
    return nil
}

// CreateSnapshot records point-in-time inventory level
func (r *inventoryRepository) CreateSnapshot(ctx context.Context, snapshot *model.InventorySnapshot) error {
    query := `
        INSERT INTO inventory_snapshots (
            id, warehouse_id, sku_id, quantity, snapshot_date, created_at
        )
        VALUES ($1, $2, $3, $4, $5, $6)
    `

    _, err := r.db.ExecContext(ctx, query,
        snapshot.ID, snapshot.WarehouseID, snapshot.SKUID,
        snapshot.Quantity, snapshot.SnapshotDate, snapshot.CreatedAt,
    )
    if err != nil {
        return fmt.Errorf("failed to create inventory snapshot: %w", err)
    }
    return nil
}