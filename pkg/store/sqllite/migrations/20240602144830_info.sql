-- +goose Up
-- +goose StatementBegin
CREATE TABLE if not exists info (
        info_id INTEGER PRIMARY KEY AUTOINCREMENT,
        last_time_update timestamp not null,
        jwt_Token text
);
-- +goose StatementEnd
INSERT INTO info (info_id, last_time_update, jwt_Token)
VALUES (1, '2022-06-01 00:00:00', '');
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS info;
-- +goose StatementEnd
