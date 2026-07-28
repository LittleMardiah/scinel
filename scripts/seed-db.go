package main

import (
"context"
"database/sql"
"fmt"
"log"
"os"
"time"

_ "github.com/lib/pq"
)

func main() {
// Load database URL from environment
databaseURL := os.Getenv("DATABASE_URL")
if databaseURL == "" {
log.Fatal("DATABASE_URL not set")
}

// Connect to database
db, err := sql.Open("postgres", databaseURL)
if err != nil {
log.Fatalf("Failed to connect to database: %v", err)
}
defer db.Close()

ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

if err = seedData(ctx, db); err != nil {
log.Fatalf("Failed to seed data: %v", err)
}

fmt.Println("✅ Data seeding complete")
}

func seedData(ctx context.Context, db *sql.DB) error {
// Seed warehouses
warehouseQuery := `
INSERT INTO warehouses (code, name, location, capacity)
VALUES ('WH-JKT', 'Jakarta Main', 'Jakarta', 10000),
       ('WH-SBY', 'Surabaya Branch', 'Surabaya', 5000),
       ('WH-BDG', 'Bandung Branch', 'Bandung', 3000)
ON CONFLICT (code) DO NOTHING
`
if _, err := db.ExecContext(ctx, warehouseQuery); err != nil {
return fmt.Errorf("failed to seed warehouses: %w", err)
}
fmt.Println("✅ Warehouses seeded")

// Seed SKUs
skuQuery := `
INSERT INTO skus (sku_code, name, category, unit_price)
VALUES 
('SKU-NIKE-001', 'Nike Shoes Model A', 'Shoes', 1500000),
('SKU-NIKE-002', 'Nike Shirt Model B', 'Apparel', 300000),
('SKU-ADIDAS-001', 'Adidas Shoes Model C', 'Shoes', 1200000),
('SKU-ADIDAS-002', 'Adidas Hat', 'Accessories', 150000),
('SKU-PUMA-001', 'Puma Shoes Model D', 'Shoes', 900000)
ON CONFLICT (sku_code) DO NOTHING
`
if _, err := db.ExecContext(ctx, skuQuery); err != nil {
return fmt.Errorf("failed to seed SKUs: %w", err)
}
fmt.Println("✅ SKUs seeded")

// Seed inventory (5 SKUs × 3 warehouses = 15 records)
inventoryQuery := `
INSERT INTO inventory (warehouse_id, sku_id, quantity, reorder_point, max_stock_level)
SELECT w.id, s.id, FLOOR(RANDOM() * 1000)::INT, 100, 5000
FROM warehouses w, skus s
ON CONFLICT (warehouse_id, sku_id) DO NOTHING
`
if _, err := db.ExecContext(ctx, inventoryQuery); err != nil {
return fmt.Errorf("failed to seed inventory: %w", err)
}
fmt.Println("✅ Inventory seeded")

fmt.Println("✅ All seed data inserted successfully")
return nil
}
