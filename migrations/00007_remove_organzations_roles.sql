-- +goose Up
-- +goose StatementBegin

DROP TABLE IF EXISTS organization_rolesl;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS organization_roles (
    id BIGSERIAL PRIMARY KEY,
    organization_id BIGSERIAL NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    role_id BIGSERIAL NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now(),
    UNIQUE(organization_id, role_id)
);
-- +goose StatementEnd
