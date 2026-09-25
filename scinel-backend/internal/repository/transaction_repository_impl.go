package repository

import (
    "context"
    "fmt"
    "strings"

    "github.com/google/uuid"
    "github.com/LittleMardiah/scinel/internal/database"
    "github.com/LittleMardiah/scinel/internal/model"
)

type transactionRepository struct {
    db *database.DB
}

func NewTransactionRepository(db *database.DB) TransactionRepository {
    return &transactionRepository{db: db}
}

const transactionColumns = `
    t.id, t.warehouse_id, t.sku_id, s.sku_code,
    t.type, t.quantity, t.reason,
    t.previous_quantity, t.new_quantity, t.created_by, t.created_at
`

const transactionSelect = `
    SELECT ` + transactionColumns + `
    FROM transactions t
    JOIN skus s ON t.sku_id = s.id
`

func (r *transactionRepository) Create(ctx context.Context, txn *model.Transaction) (*model.Transaction, error) {
    // Note: the `details` column does not exist in the current schema, so it is
    // intentionally omitted from the INSERT. Use pessimistic locking via the
    // caller/service by wrapping in WithTx if update-on-create is needed.
    query := `
        INSERT INTO transactions (id, warehouse_id, sku_id, type, quantity, reason,
                                  previous_quantity, new_quantity, created_by, created_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
        RETURNING id, warehouse_id, sku_id, type, quantity, reason,
                  previous_quantity, new_quantity, created_by, created_at
    `
    var created model.Transaction
    err := r.db.QueryRowContext(ctx, query,
        txn.ID, txn.WarehouseID, txn.SKUID, txn.Type, txn.Quantity, txn.Reason,
        txn.PreviousQuantity, txn.NewQuantity, txn.CreatedBy, txn.CreatedAt,
    ).Scan(
        &created.ID, &created.WarehouseID, &created.SKUID, &created.Type, &created.Quantity,
        &created.Reason, &created.PreviousQuantity, &created.NewQuantity, &created.CreatedBy,
        &created.CreatedAt,
    )
    if err != nil {
        return nil, fmt.Errorf("failed to create transaction: %w", err)
    }
    return &created, nil
}

func (r *transactionRepository) GetByID(ctx context.Context, txnID uuid.UUID) (*model.Transaction, error) {
    query := transactionSelect + `
        WHERE t.id = $1
    `
    var txn model.Transaction
    err := r.db.QueryRowContext(ctx, query, txnID).Scan(
        &txn.ID, &txn.WarehouseID, &txn.SKUID, &txn.SKUCode,
        &txn.Type, &txn.Quantity, &txn.Reason,
        &txn.PreviousQuantity, &txn.NewQuantity, &txn.CreatedBy, &txn.CreatedAt,
    )
    if err != nil {
        return nil, fmt.Errorf("failed to get transaction by ID: %w", err)
    }
    return &txn, nil
}

func (r *transactionRepository) Query(ctx context.Context, filter model.TransactionFilter, limit, offset int) ([]model.Transaction, int64, error) {
    var conds []string
    var args []interface{}

    if filter.WarehouseID != nil {
        args = append(args, *filter.WarehouseID)
        conds = append(conds, fmt.Sprintf("t.warehouse_id = $%d", len(args)))
    }
    if filter.SKUID != nil {
        args = append(args, *filter.SKUID)
        conds = append(conds, fmt.Sprintf("t.sku_id = $%d", len(args)))
    }
    if filter.Type != nil {
        args = append(args, *filter.Type)
        conds = append(conds, fmt.Sprintf("t.type = $%d", len(args)))
    }
    if filter.CreatedBy != nil {
        args = append(args, *filter.CreatedBy)
        conds = append(conds, fmt.Sprintf("t.created_by = $%d", len(args)))
    }
    if filter.Reason != nil {
        args = append(args, "%"+*filter.Reason+"%")
        conds = append(conds, fmt.Sprintf("t.reason ILIKE $%d", len(args)))
    }
    if filter.DateFrom != nil {
        args = append(args, *filter.DateFrom)
        conds = append(conds, fmt.Sprintf("t.created_at >= $%d", len(args)))
    }
    if filter.DateTo != nil {
        args = append(args, *filter.DateTo)
        conds = append(conds, fmt.Sprintf("t.created_at <= $%d", len(args)))
    }

    where := ""
    if len(conds) > 0 {
        where = " WHERE " + strings.Join(conds, " AND ")
    }

    listQuery := transactionSelect + where + `
        ORDER BY t.created_at DESC
        LIMIT ` + argPlaceholder(len(args)+1) + ` OFFSET ` + argPlaceholder(len(args)+2)
    allArgs := append(append([]interface{}{}, args...), limit, offset)

    rows, err := r.db.QueryContext(ctx, listQuery, allArgs...)
    if err != nil {
        return nil, 0, fmt.Errorf("failed to query transactions: %w", err)
    }
    defer rows.Close()

    var items []model.Transaction
    for rows.Next() {
        var txn model.Transaction
        if err := rows.Scan(
            &txn.ID, &txn.WarehouseID, &txn.SKUID, &txn.SKUCode,
            &txn.Type, &txn.Quantity, &txn.Reason,
            &txn.PreviousQuantity, &txn.NewQuantity, &txn.CreatedBy, &txn.CreatedAt,
        ); err != nil {
            return nil, 0, fmt.Errorf("failed to scan transaction: %w", err)
        }
        items = append(items, txn)
    }
    if err = rows.Err(); err != nil {
        return nil, 0, fmt.Errorf("rows iteration error: %w", err)
    }

    var total int64
    countQuery := `SELECT COUNT(*) FROM transactions t` + where
    err = r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
    if err != nil {
        return nil, 0, fmt.Errorf("failed to count transactions: %w", err)
    }

    return items, total, nil
}

func (r *transactionRepository) GetByWarehouse(ctx context.Context, warehouseID uuid.UUID, limit, offset int) ([]model.Transaction, int64, error) {
    query := transactionSelect + `
        WHERE t.warehouse_id = $1
        ORDER BY t.created_at DESC
        LIMIT $2 OFFSET $3
    `
    rows, err := r.db.QueryContext(ctx, query, warehouseID, limit, offset)
    if err != nil {
        return nil, 0, fmt.Errorf("failed to query transactions by warehouse: %w", err)
    }
    defer rows.Close()

    var items []model.Transaction
    for rows.Next() {
        var txn model.Transaction
        if err := rows.Scan(
            &txn.ID, &txn.WarehouseID, &txn.SKUID, &txn.SKUCode,
            &txn.Type, &txn.Quantity, &txn.Reason,
            &txn.PreviousQuantity, &txn.NewQuantity, &txn.CreatedBy, &txn.CreatedAt,
        ); err != nil {
            return nil, 0, fmt.Errorf("failed to scan transaction: %w", err)
        }
        items = append(items, txn)
    }
    if err = rows.Err(); err != nil {
        return nil, 0, fmt.Errorf("rows iteration error: %w", err)
    }

    var total int64
    err = r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM transactions WHERE warehouse_id = $1`, warehouseID).Scan(&total)
    if err != nil {
        return nil, 0, fmt.Errorf("failed to count transactions by warehouse: %w", err)
    }

    return items, total, nil
}

func (r *transactionRepository) GetBySKU(ctx context.Context, skuID uuid.UUID, limit, offset int) ([]model.Transaction, int64, error) {
    query := transactionSelect + `
        WHERE t.sku_id = $1
        ORDER BY t.created_at DESC
        LIMIT $2 OFFSET $3
    `
    rows, err := r.db.QueryContext(ctx, query, skuID, limit, offset)
    if err != nil {
        return nil, 0, fmt.Errorf("failed to query transactions by SKU: %w", err)
    }
    defer rows.Close()

    var items []model.Transaction
    for rows.Next() {
        var txn model.Transaction
        if err := rows.Scan(
            &txn.ID, &txn.WarehouseID, &txn.SKUID, &txn.SKUCode,
            &txn.Type, &txn.Quantity, &txn.Reason,
            &txn.PreviousQuantity, &txn.NewQuantity, &txn.CreatedBy, &txn.CreatedAt,
        ); err != nil {
            return nil, 0, fmt.Errorf("failed to scan transaction: %w", err)
        }
        items = append(items, txn)
    }
    if err = rows.Err(); err != nil {
        return nil, 0, fmt.Errorf("rows iteration error: %w", err)
    }

    var total int64
    err = r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM transactions WHERE sku_id = $1`, skuID).Scan(&total)
    if err != nil {
        return nil, 0, fmt.Errorf("failed to count transactions by SKU: %w", err)
    }

    return items, total, nil
}

// argPlaceholder builds a numbered placeholder (e.g. 1 -> "$1").
func argPlaceholder(n int) string {
    return fmt.Sprintf("$%d", n)
}
