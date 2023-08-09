-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS products (
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    store_id uuid NOT NULL,
    name VARCHAR(50) NOT NULL,
    url VARCHAR(100) NOT NULL UNIQUE,
    price FLOAT NOT NULL,
    description TEXT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (store_id) REFERENCES stores (id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE IF EXISTS products;
