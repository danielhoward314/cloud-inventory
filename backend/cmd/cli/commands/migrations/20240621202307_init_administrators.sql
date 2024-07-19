-- +goose Up
-- +goose StatementBegin
-- Create an enum that signals what hashing algorithm should be used for the password_hash
CREATE TYPE password_hash_type AS ENUM (
    'BCRYPT'
);

CREATE TYPE authorization_role AS ENUM (
    'PRIMARY_ADMIN',
    'SECONDARY_ADMIN'
);

-- Create the administrators table
CREATE TABLE administrators (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    display_name VARCHAR(255),
    password_hash_type password_hash_type NOT NULL,
    authorization_role authorization_role NOT NULL,
    password_hash TEXT NOT NULL,
    verified BOOLEAN NOT NULL DEFAULT false,
    organization_id UUID NOT NULL,
    CONSTRAINT fk_organization
        FOREIGN KEY(organization_id) 
        REFERENCES organizations(id)
        ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE administrators;
-- +goose StatementEnd
