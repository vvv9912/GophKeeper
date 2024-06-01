-- +goose Up
-- +goose StatementBegin
CREATE TABLE if not exists users_data (
  user_data_id SERIAL PRIMARY KEY,
  data_id integer, --unique
  user_id integer not null ,
  data_type integer not null ,
  name text not null,
  description text not null,
  hash text not null,
  created_at timestamp default NOW(),
  update_at timestamp default NOW(),
  is_Deleted boolean default false
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users_data;
-- +goose StatementEnd
