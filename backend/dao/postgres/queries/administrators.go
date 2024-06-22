package queries

const AdministratorsInsert = `INSERT INTO administrators (email, display_name, organization_id)
VALUES ($1, $2, $3)
RETURNING id`
