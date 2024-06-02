-- +goose Up
-- +goose StatementBegin
CREATE TABLE if not exists users_data (
  user_data_id INTEGER PRIMARY KEY, --AUTOINCREMENT,
  data_id integer not null , --unique
  data_type integer not null ,
  name text not null,
  description text not null,
  hash text not null,
  created_at timestamp not null ,
  update_at timestamp not null ,
  is_Deleted boolean default false
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users_data;
-- +goose StatementEnd
