-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS  tasks_edges (
    id BIGSERIAL PRIMARY KEY,
    pipeline_id bigserial NOT NULL REFERENCES pipelines(id) ON DELETE CASCADE,
    from_node BIGSERIAL NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    to_node BIGSERIAL NOT NULL REFERENCES tasks(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE tasks_edges;
-- +goose StatementEnd
