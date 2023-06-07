-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS orders (
    order_id UUID NOT NULL PRIMARY KEY,
    courier_id UUID NULL REFERENCES couriers(courier_id),
    weight FLOAT NOT NULL,
    regions INT NOT NULL,
    delivery_hours TEXT[] NOT NULL,
    cost INT NOT NULL,
    completed_time TIMESTAMP NOT NULL,
    distribution_date TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS orders;
-- +goose StatementEnd
