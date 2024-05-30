-- +goose Up
-- +goose StatementBegin
CREATE TABLE if not exists data (
    data_id serial primary key,
    encrypt_data bytea
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
