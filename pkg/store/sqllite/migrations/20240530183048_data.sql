-- +goose Up
-- +goose StatementBegin
CREATE TABLE if not exists data (
    data_id INTEGER PRIMARY KEY AUTOINCREMENT,
    encrypt_data bytea not null,
    meta_data bytea default '{}'

);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS data;
-- +goose StatementEnd
