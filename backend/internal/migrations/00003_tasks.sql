-- +goose Up
-- +goose StatementBegin
CREATE TABLE tasks
(
    id          SERIAL PRIMARY KEY,
--     project_id  INTEGER      NOT NULL REFERENCES projects (id),
--     assigned_id INTEGER REFERENCES users (id),
    title       VARCHAR(255) NOT NULL,
    description TEXT,
    status      VARCHAR(255) NOT NULL default 'todo',
    priority    VARCHAR(255) NOT NULL default 'medium',
    due_date    TIMESTAMP,
    created_at  TIMESTAMP             DEFAULT now(),
    updated_at  TIMESTAMP             DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE tasks;
-- +goose StatementEnd
