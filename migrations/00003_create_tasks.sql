-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS tasks (
    id UUID PRIMARY KEY,
    pipeline_id bigserial NOT NULL REFERENCES pipelines(id) ON DELETE CASCADE,
    title VARCHAR NOT NULL,
    description TEXT,
    ui_display INT,
    ----
    type VARCHAR(255) NOT NULL,
    config JSONB NOT NULL,
    status VARCHAR NOT NULL,
    error TEXT,
    -- version control
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()

);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE tasks;

-- +goose StatementEnd
