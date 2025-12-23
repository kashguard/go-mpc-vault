-- +migrate Up
CREATE TABLE IF NOT EXISTS chains (
    id varchar(50) PRIMARY KEY, -- e.g. 'ETH', 'BTC', 'SOL_MAINNET'
    name varchar(100) NOT NULL, -- e.g. 'Ethereum Mainnet'
    type VARCHAR(20) NOT NULL, -- 'EVM', 'UTXO', 'SOLANA'
    chain_id varchar(50), -- Chain ID, e.g. '1', 'solana-mainnet'
    algorithm varchar(50) NOT NULL, -- 'ECDSA', 'EdDSA', 'Schnorr'
    curve varchar(50) NOT NULL, -- 'secp256k1', 'ed25519'
    currency_symbol varchar(20) NOT NULL, -- 'ETH', 'BTC', 'SOL'
    rpc_url text,
    explorer_url text,
    icon_url text,
    is_testnet boolean DEFAULT FALSE,
    created_at timestamptz DEFAULT NOW(),
    updated_at timestamptz DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS assets (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid (),
    chain_id varchar(50) REFERENCES chains (id),
    symbol varchar(20) NOT NULL, -- 'USDT', 'USDC', 'ETH'
    name varchar(100) NOT NULL, -- 'Tether USD'
    type VARCHAR(20) NOT NULL, -- 'NATIVE', 'ERC20', 'TRC20', 'SPL'
    contract_address varchar(255), -- Token contract address (DEFAULT NULL for Native)
    decimals int NOT NULL DEFAULT 0,
    icon_url text,
    is_active boolean DEFAULT TRUE,
    created_at timestamptz DEFAULT NOW(),
    updated_at timestamptz DEFAULT NOW(),
    UNIQUE (chain_id, contract_address)
);

-- +migrate Down
DROP TABLE IF EXISTS assets;

DROP TABLE IF EXISTS chains;

