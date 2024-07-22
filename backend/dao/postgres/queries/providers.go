package queries

const ProvidersSelectByOrganizationID = `SELECT
	id,
	external_identifier,
	display_name,
	provider_name,
	metadata,
	organization_id
FROM providers where organization_id = $1`
