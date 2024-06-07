-- +goose Up
-- +goose StatementBegin
CREATE TABLE if not exists info (
        info_id INTEGER PRIMARY KEY AUTOINCREMENT,
        last_time_update timestamp not null,
        jwt_Token text
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS info;
-- +goose StatementEnd
