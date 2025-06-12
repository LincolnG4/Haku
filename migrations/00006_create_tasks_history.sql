-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS tasks_history (
    id BIGSERIAL PRIMARY KEY,
    pipeline_id BIGSERIAL NOT NULL,
    user_id bigserial NOT NULL,
    name VARCHAR(255),
    status VARCHAR(255),
    version INT,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now(),
    recorded_at TIMESTAMP DEFAULT now(), -- when this snapshot was taken
    action VARCHAR(20), -- 'create', 'update', 'delete'
    
    FOREIGN KEY (pipeline_id) REFERENCES pipelines(id)
);

CREATE OR REPLACE FUNCTION log_tasks_history()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO tasks_history (
        pipeline_id, user_id, name, status, version,
        created_at, updated_at, recorded_at, action
    )
    VALUES (
        OLD.id, OLD.user_id, OLD.name, OLD.status, OLD.version,
        OLD.created_at, OLD.updated_at, now(), TG_OP
    );
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_tasks_history
AFTER UPDATE OR DELETE ON tasks
FOR EACH ROW EXECUTE FUNCTION log_tasks_history();


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE tasks_history;
-- +goose StatementEnd
