-- +goose Up
-- +goose StatementBegin
CREATE TABLE if not exists data_types (
    data_type SERIAL PRIMARY KEY,
    type text not null
);
INSERT INTO data_types (type)
VALUES ('credentials'), ('credit_card_data'), ('binary_file'), ('text_file');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS data_types;
-- +goose StatementEnd
