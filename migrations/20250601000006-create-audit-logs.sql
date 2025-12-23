-- +migrate Up
CREATE TABLE IF NOT EXISTS audit_logs (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid (),
    organization_id uuid REFERENCES organizations (id) ON DELETE CASCADE,
    user_id uuid REFERENCES users (id),
    action varchar(50) NOT NULL, -- 'LOGIN', 'CREATE_VAULT', 'APPROVE_TX', 'MODIFY_POLICY'
    resource_type varchar(50), -- 'vault', 'wallet', 'user'
    resource_id varchar(255),
    ip_address varchar(45),
    user_agent text,
    details jsonb,
    created_at timestamptz DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_audit_logs_organization_id ON audit_logs (organization_id);

CREATE INDEX IF NOT EXISTS idx_audit_logs_user_id ON audit_logs (user_id);

-- +migrate Down
DROP TABLE IF EXISTS audit_logs;

