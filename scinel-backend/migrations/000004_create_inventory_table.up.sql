CREATE TABLE IF NOT EXISTS inventory (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    warehouse_id UUID NOT NULL REFERENCES warehouses(id) ON DELETE CASCADE,
    sku_id UUID NOT NULL REFERENCES skus(id) ON DELETE CASCADE,
    quantity INT NOT NULL DEFAULT 0,
    reorder_point INT DEFAULT 100,
    max_stock_level INT DEFAULT 5000,
    last_restock_date TIMESTAMP WITH TIME ZONE,
    last_transaction_id UUID,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX idx_inventory_warehouse_sku ON inventory(warehouse_id, sku_id);
CREATE INDEX idx_inventory_quantity_low ON inventory(quantity) WHERE quantity <= reorder_point;
CREATE INDEX idx_inventory_warehouse ON inventory(warehouse_id);
CREATE INDEX idx_inventory_sku ON inventory(sku_id);