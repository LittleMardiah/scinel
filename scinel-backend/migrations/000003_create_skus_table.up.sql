CREATE TABLE IF NOT EXISTS skus (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sku_code VARCHAR(100) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    category VARCHAR(100),
    unit_price DECIMAL(10, 2),
    status VARCHAR(50) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'discontinued', 'coming_soon')),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_skus_code ON skus(sku_code);
CREATE INDEX idx_skus_category ON skus(category);
CREATE INDEX idx_skus_status ON skus(status);