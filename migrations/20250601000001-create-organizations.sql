-- +migrate Up
CREATE TABLE IF NOT EXISTS organizations (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid (),
    name varchar(100) NOT NULL,
    owner_id uuid NOT NULL REFERENCES users (id),
    created_at timestamptz DEFAULT NOW(),
    updated_at timestamptz DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS organization_members (
    organization_id uuid REFERENCES organizations (id) ON DELETE CASCADE,
    user_id uuid REFERENCES users (id) ON DELETE CASCADE,
    role VARCHAR(50) NOT NULL, -- 'admin', 'operator', 'auditor'
    created_at timestamptz DEFAULT NOW(),
    PRIMARY KEY (organization_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_organizations_owner_id ON organizations (owner_id);

CREATE INDEX IF NOT EXISTS idx_organization_members_user_id ON organization_members (user_id);

-- +migrate Down
DROP TABLE IF EXISTS organization_members;

DROP TABLE IF EXISTS organizations;

