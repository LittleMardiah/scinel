package model

import (
    "database/sql/driver"
    "encoding/json"
    "time"

    "github.com/google/uuid"
)

// Transaction represents an inventory in/out transaction
type Transaction struct {
    ID               uuid.UUID       `db:"id"`
    WarehouseID      uuid.UUID       `db:"warehouse_id"`
    SKUID            uuid.UUID       `db:"sku_id"`
    SKUCode          string          `db:"sku_code"` // Denormalized for queries
    Type             TransactionType `db:"type"`
    Quantity         int32           `db:"quantity"` // Positive or negative
    Reason           string          `db:"reason"`   // e.g., "Supplier delivery #PO-2026-001"
    PreviousQuantity int32           `db:"previous_quantity"`
    NewQuantity      int32           `db:"new_quantity"`
    CreatedBy        uuid.UUID       `db:"created_by"`
    Details          json.RawMessage `db:"details"` // JSONB: Additional metadata
    CreatedAt        time.Time       `db:"created_at"`
}

// TransactionType represents transaction direction
type TransactionType string

const (
    TransactionIn  TransactionType = "in"  // Inventory increase (purchase, return, correction)
    TransactionOut TransactionType = "out" // Inventory decrease (sale, damage, transfer)
)

// Scan implements sql.Scanner interface for TransactionType
func (t *TransactionType) Scan(value interface{}) error {
    *t = TransactionType(value.(string))
    return nil
}

// Value implements driver.Valuer interface for TransactionType
func (t TransactionType) Value() (driver.Value, error) {
    return string(t), nil
}

// TransactionFilter for querying transactions with filters
type TransactionFilter struct {
    WarehouseID *uuid.UUID
    SKUID       *uuid.UUID
    Type        *TransactionType
    DateFrom    *time.Time
    DateTo      *time.Time
    CreatedBy   *uuid.UUID
    Reason      *string
}