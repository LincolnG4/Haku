-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS pipelines_history (
    id BIGSERIAL PRIMARY KEY,
    pipeline_id BIGSERIAL NOT NULL REFERENCES pipelines(id) ON DELETE CASCADE,
    user_id BIGSERIAL NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255),
    status VARCHAR(255),
    version INT,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now(),
    action VARCHAR(20) -- 'create', 'update', 'delete'
);


CREATE OR REPLACE FUNCTION log_pipeline_history()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO pipelines_history (
        pipeline_id, user_id, name, status, version,
        created_at, updated_at, action
    )
    VALUES (
        OLD.id, OLD.user_id, OLD.name, OLD.status, OLD.version,
        OLD.created_at, now(), TG_OP
    );
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_pipeline_history
AFTER UPDATE OR DELETE ON pipelines
FOR EACH ROW EXECUTE FUNCTION log_pipeline_history();


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE pipelines_history;
-- +goose StatementEnd
