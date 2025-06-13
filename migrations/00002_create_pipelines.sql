-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS pipelines (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGSERIAL NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255),
    status VARCHAR(15),
    -- version control
    version INT,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE pipelines;
-- +goose StatementEnd
