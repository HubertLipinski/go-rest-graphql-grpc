-- +goose Up
-- +goose StatementBegin
CREATE TABLE projects
(
    id          SERIAL PRIMARY KEY,
    name        varchar(255) NOT NULL,
    description TEXT,
    owner_id    INTEGER      NOT NULL REFERENCES users (id),
    created_at  TIMESTAMP DEFAULT now(),
    updated_at  TIMESTAMP DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE projects;
-- +goose StatementEnd
