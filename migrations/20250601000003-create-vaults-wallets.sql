-- +migrate Up
CREATE TABLE IF NOT EXISTS vaults (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid (),
    organization_id uuid REFERENCES organizations (id) ON DELETE CASCADE,
    name varchar(100) NOT NULL,
    threshold int NOT NULL DEFAULT 0,
    created_at timestamptz DEFAULT NOW(),
    updated_at timestamptz DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS vault_keys (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid (),
    vault_id uuid REFERENCES vaults (id) ON DELETE CASCADE,
    key_id varchar(255) NOT NULL, -- MPC Infra KeyID
    algorithm varchar(50) NOT NULL, -- 'ECDSA', 'EdDSA'
    curve varchar(50) NOT NULL, -- 'secp256k1', 'ed25519'
    public_key_hex text NOT NULL, -- Aggregated Public Key
    created_at timestamptz DEFAULT NOW(),
    updated_at timestamptz DEFAULT NOW(),
    UNIQUE (vault_id, algorithm)
);

CREATE TABLE IF NOT EXISTS wallets (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid (),
    vault_id uuid REFERENCES vaults (id) ON DELETE CASCADE,
    chain_id varchar(50) REFERENCES chains (id),
    key_id varchar(255) NOT NULL, -- vault_keys.key_id
    address varchar(255) NOT NULL,
    derive_path varchar(255) NOT NULL, -- e.g. "m/44'/60'/0'/0/0"
    derive_index int NOT NULL,
    created_at timestamptz DEFAULT NOW(),
    updated_at timestamptz DEFAULT NOW(),
    UNIQUE (vault_id, chain_id, derive_index)
);

CREATE TABLE IF NOT EXISTS wallet_balances (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid (),
    wallet_id uuid REFERENCES wallets (id) ON DELETE CASCADE,
    asset_id uuid REFERENCES assets (id),
    balance DECIMAL(36, 18) DEFAULT 0,
    raw_balance varchar(100), -- Wei/Lamports
    updated_at timestamptz DEFAULT NOW(),
    UNIQUE (wallet_id, asset_id)
);

CREATE INDEX IF NOT EXISTS idx_vaults_organization_id ON vaults (organization_id);

CREATE INDEX IF NOT EXISTS idx_wallets_chain_id ON wallets (chain_id);

CREATE INDEX IF NOT EXISTS idx_wallet_balances_asset_id ON wallet_balances (asset_id);

-- +migrate Down
DROP TABLE IF EXISTS wallet_balances;

DROP TABLE IF EXISTS wallets;

DROP TABLE IF EXISTS vault_keys;

DROP TABLE IF EXISTS vaults;

