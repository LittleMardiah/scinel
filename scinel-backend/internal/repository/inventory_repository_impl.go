package repository

import (
    "context"
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
    // TODO: Implement
    return nil, fmt.Errorf("GetBySKU not implemented yet")
}

// Create creates a new inventory record
func (r *inventoryRepository) Create(ctx context.Context, inventory *model.Inventory) (*model.Inventory, error) {
    // TODO: Implement
    return nil, fmt.Errorf("Create not implemented yet")
}

// UpdateQuantity atomically updates inventory quantity
func (r *inventoryRepository) UpdateQuantity(ctx context.Context, warehouseID, skuID uuid.UUID, delta int32) (*model.Inventory, error) {
    // TODO: Implement
    return nil, fmt.Errorf("UpdateQuantity not implemented yet")
}

// LockForUpdate locks inventory record for update in a transaction
func (r *inventoryRepository) LockForUpdate(ctx context.Context, warehouseID, skuID uuid.UUID) (*model.Inventory, error) {
    // TODO: Implement
    return nil, fmt.Errorf("LockForUpdate not implemented yet")
}

// GetBySKUCached retrieves inventory with Redis cache fallback (30s TTL)
func (r *inventoryRepository) GetBySKUCached(ctx context.Context, warehouseID, skuID uuid.UUID) (*model.Inventory, error) {
    // TODO: Implement
    return nil, fmt.Errorf("GetBySKUCached not implemented yet")
}

// InvalidateCache removes inventory from cache after update
func (r *inventoryRepository) InvalidateCache(ctx context.Context, warehouseID, skuID uuid.UUID) error {
    // TODO: Implement
    return fmt.Errorf("InvalidateCache not implemented yet")
}

// CreateSnapshot records point-in-time inventory level
func (r *inventoryRepository) CreateSnapshot(ctx context.Context, snapshot *model.InventorySnapshot) error {
    // TODO: Implement
    return fmt.Errorf("CreateSnapshot not implemented yet")
}