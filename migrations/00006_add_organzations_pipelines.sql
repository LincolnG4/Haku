-- +goose Up
-- +goose StatementBegin
ALTER TABLE pipelines 
ADD COLUMN organization_id INTEGER REFERENCES organizations(id);

ALTER TABLE pipelines DROP COLUMN user_id;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE pipelines DROP COLUMN organization_id;
ALTER TABLE pipelines 
ADD COLUMN user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE;
-- +goose StatementEnd
