-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS couriers (
    courier_id UUID NOT NULL PRIMARY KEY,
    courier_type CHAR(4) NOT NULL,
    regions INTEGER[] NOT NULL,
    working_hours TEXT[] NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS couriers;
-- +goose StatementEnd
