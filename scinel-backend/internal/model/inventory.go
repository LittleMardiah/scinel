package model

import (
    "time"

    "github.com/google/uuid"
)

// Inventory represents a stock level for a SKU in a warehouse
type Inventory struct {
    ID                uuid.UUID  `db:"id"`
    WarehouseID       uuid.UUID  `db:"warehouse_id"`
    SKUID             uuid.UUID  `db:"sku_id"`
    SKUCode           string     `db:"sku_code"`           // e.g., "SKU-NIKE-001"
    WarehouseName     string     `db:"-"`                  // NOT from DB, filled by JOIN
    SKUName           string     `db:"-"`                  // NOT from DB, filled by JOIN
    Quantity          int32      `db:"quantity"`
    ReorderPoint      int32      `db:"reorder_point"`      // Trigger for low-stock alerts
    MaxStockLevel     int32      `db:"max_stock_level"`    // Safety cap
    LastRestockDate   *time.Time `db:"last_restock_date"`
    LastTransactionID *uuid.UUID `db:"last_transaction_id"`
    CreatedAt         time.Time  `db:"created_at"`
    UpdatedAt         time.Time  `db:"updated_at"`
}

// InventorySnapshot represents a point-in-time stock level (for historical tracking)
type InventorySnapshot struct {
    ID           uuid.UUID `db:"id"`
    WarehouseID  uuid.UUID `db:"warehouse_id"`
    SKUID        uuid.UUID `db:"sku_id"`
    Quantity     int32     `db:"quantity"`
    SnapshotDate time.Time `db:"snapshot_date"`
    CreatedAt    time.Time `db:"created_at"`
}

// InventoryStatus represents inventory status classification
type InventoryStatus string

const (
    StatusNormal   InventoryStatus = "normal"    // Quantity > reorder_point
    StatusLowStock InventoryStatus = "low_stock" // Quantity <= reorder_point
    StatusCritical InventoryStatus = "critical"  // Quantity == 0 or negative
)

// GetStatus returns inventory status based on quantity
func (i *Inventory) GetStatus() InventoryStatus {
    if i.Quantity <= 0 {
        return StatusCritical
    }
    if i.Quantity <= i.ReorderPoint {
        return StatusLowStock
    }
    return StatusNormal
}