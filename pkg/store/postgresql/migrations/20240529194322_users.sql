-- +goose Up
-- +goose StatementBegin
CREATE TABLE if not exists users (
    user_id SERIAL PRIMARY KEY,
    login varchar(20) not null,
    password varchar(32) not null, --из расчета хэш функции
    created_at timestamp default NOW(),
    unique(login)
    );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
