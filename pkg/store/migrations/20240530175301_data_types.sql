-- +goose Up
-- +goose StatementBegin
CREATE TABLE if not exists data_types (
    data_type SERIAL PRIMARY KEY,
    type text
);
INSERT INTO data_types (type)
VALUES ('credentials'), ('credit_card_data'), ('file');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
