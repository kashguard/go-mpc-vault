-- +migrate Up
CREATE TABLE IF NOT EXISTS address_book (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid (),
    organization_id uuid REFERENCES organizations (id) ON DELETE CASCADE,
    chain_id varchar(50) REFERENCES chains (id),
    address varchar(255) NOT NULL,
    name varchar(100) NOT NULL,
    is_whitelisted boolean DEFAULT TRUE,
    created_by uuid REFERENCES users (id),
    created_at timestamptz DEFAULT NOW(),
    updated_at timestamptz DEFAULT NOW(),
    UNIQUE (organization_id, chain_id, address)
);

CREATE TABLE IF NOT EXISTS spending_limits (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid (),
    vault_id uuid REFERENCES vaults (id) ON DELETE CASCADE,
    asset_id uuid REFERENCES assets (id), -- DEFAULT NULL means USD total limit
    amount DECIMAL(36, 18) NOT NULL,
    window_seconds int NOT NULL, -- 3600(1h), 86400(24h)
    action varchar(20) DEFAULT 'REJECT', -- 'REJECT', 'REQUIRE_ADMIN'
    created_at timestamptz DEFAULT NOW(),
    updated_at timestamptz DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_address_book_chain_id ON address_book (chain_id);

CREATE INDEX IF NOT EXISTS idx_address_book_created_by ON address_book (created_by);

CREATE INDEX IF NOT EXISTS idx_spending_limits_asset_id ON spending_limits (asset_id);

CREATE INDEX IF NOT EXISTS idx_spending_limits_vault_id ON spending_limits (vault_id);

-- +migrate Down
DROP TABLE IF EXISTS spending_limits;

DROP TABLE IF EXISTS address_book;

