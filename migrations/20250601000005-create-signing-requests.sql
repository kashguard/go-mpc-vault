-- +migrate Up
CREATE TABLE IF NOT EXISTS signing_requests (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid (),
    vault_id uuid REFERENCES vaults (id) ON DELETE CASCADE,
    wallet_id uuid REFERENCES wallets (id),
    initiator_id uuid REFERENCES users (id),
    tx_data text NOT NULL, -- Raw tx data (Hex)
    tx_hash varchar(255), -- Tx Hash
    amount DECIMAL(36, 18),
    to_address varchar(255),
    note text,
    status varchar(50) DEFAULT 'pending', -- pending, approved, signing, signed, rejected, failed
    mpc_session_id varchar(255), -- MPC Session ID
    signature text, -- Final Signature
    created_at timestamptz DEFAULT NOW(),
    updated_at timestamptz DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS approvals (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid (),
    request_id uuid REFERENCES signing_requests (id) ON DELETE CASCADE,
    user_id uuid REFERENCES users (id),
    action varchar(50) NOT NULL, -- 'approve', 'reject'
    comment text,
    created_at timestamptz DEFAULT NOW(),
    updated_at timestamptz DEFAULT NOW(),
    UNIQUE (request_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_signing_requests_initiator_id ON signing_requests (initiator_id);

CREATE INDEX IF NOT EXISTS idx_signing_requests_vault_id ON signing_requests (vault_id);

CREATE INDEX IF NOT EXISTS idx_signing_requests_wallet_id ON signing_requests (wallet_id);

CREATE INDEX IF NOT EXISTS idx_approvals_user_id ON approvals (user_id);

-- +migrate Down
DROP TABLE IF EXISTS approvals;

DROP TABLE IF EXISTS signing_requests;

