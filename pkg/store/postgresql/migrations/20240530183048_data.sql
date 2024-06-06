-- +goose Up
-- +goose StatementBegin
CREATE TABLE if not exists data (
    data_id serial primary key,
    encrypt_data bytea not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS data;
-- +goose StatementEnd
