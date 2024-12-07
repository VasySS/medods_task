-- +goose Up
-- +goose StatementBegin
ALTER TABLE user_sessions 
    ADD COLUMN used BOOLEAN NOT NULL DEFAULT false;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE user_sessions
    DROP COLUMN used;
-- +goose StatementEnd
