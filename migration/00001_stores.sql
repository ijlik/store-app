-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS stores (
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    name VARCHAR(50) NOT NULL,
    url VARCHAR(100) NOT NULL UNIQUE,
    address TEXT NOT NULL,
    phone VARCHAR(15) NOT NULL,
    operational_time_start INT NOT NULL,
    operational_time_end INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NULL,
    PRIMARY KEY (id)
);

-- +goose Down
DROP TABLE IF EXISTS stores;
