-- +migrate Up
CREATE TABLE IF NOT EXISTS user_credentials (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid (),
    user_id uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    credential_id text NOT NULL, -- WebAuthn Credential ID (Base64)
    public_key text NOT NULL, -- Public Key (COSE Key format)
    attestation_type varchar(50), -- e.g. "none", "direct"
    aaguid uuid, -- Authenticator Attestation GUID
    sign_count int DEFAULT 0, -- Sign count (anti-cloning)
    device_name varchar(100), -- e.g. "iPhone 15 Pro"
    last_used_at timestamptz,
    created_at timestamptz DEFAULT NOW(),
    updated_at timestamptz DEFAULT NOW(),
    UNIQUE (user_id, credential_id)
);

-- +migrate Down
DROP TABLE IF EXISTS user_credentials;

