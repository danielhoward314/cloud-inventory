-- +goose Up
-- +goose StatementBegin
-- Create the necessary extension for generating UUIDs
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create an enum for the billing plans
CREATE TYPE billing_plan_type AS ENUM (
    'FREE',
    'PREMIUM'
);

-- Create the organizations table
CREATE TABLE organizations (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    primary_administrator_email VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    billing_plan_type billing_plan_type NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE organizations;
-- +goose StatementEnd
