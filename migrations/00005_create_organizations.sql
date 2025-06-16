-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS organizations (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);

CREATE TABLE IF NOT EXISTS organization_members (
    id BIGSERIAL PRIMARY KEY,
    organization_id BIGSERIAL NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    user_id BIGSERIAL NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role_id BIGSERIAL NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now(),
    UNIQUE(organization_id, user_id)
);

CREATE TABLE IF NOT EXISTS roles (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    description TEXT,
    permissions JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);

-- Insert default roles
INSERT INTO roles (name, description, permissions) VALUES
    ('admin', 'Organization administrator with full access', '{"*": "*"}'::jsonb),
    ('developer', 'Regular organization member', '{"read": "*", "write": ["pipeline", "task"]}'::jsonb),
    ('viewer', 'Read-only access to organization resources', '{"read": "*"}'::jsonb);



-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS organization_roles;
DROP TABLE IF EXISTS roles;
DROP TABLE IF EXISTS organization_members;
DROP TABLE IF EXISTS organizations;
-- +goose StatementEnd 