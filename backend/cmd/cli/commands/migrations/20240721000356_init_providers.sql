-- +goose Up
-- +goose StatementBegin
-- Create the provider_name enum if it doesn't exist
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'provider_name') THEN
        CREATE TYPE provider_name AS ENUM (
            'AWS',
            'GCP',
            'AZURE'
        );
    END IF;
END$$;

CREATE TABLE IF NOT EXISTS providers (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    external_identifier TEXT NOT NULL,
    display_name VARCHAR(255),
    provider_name provider_name NOT NULL,
    metadata JSONB,
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
DROP TABLE IF EXISTS providers;
DROP TYPE IF EXISTS provider_name;
-- +goose StatementEnd
