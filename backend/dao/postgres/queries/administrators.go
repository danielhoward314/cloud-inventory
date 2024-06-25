package queries

const AdministratorsInsert = `INSERT INTO administrators (email, display_name, organization_id, password_hash_type, password_hash)
VALUES ($1, $2, $3, $4, $5)
RETURNING id`

const AdministratorsSelect = `SELECT id, email, display_name, password_hash_type, password_hash, organization_id, verified FROM administrators where id = $1`
const AdministratorsSelectByEmail = `SELECT id, email, display_name, password_hash_type, password_hash, organization_id, verified FROM administrators where email = $1`

const AdministratorsUpdate = `UPDATE administrators
SET
    email = $1,
    display_name = $2,
    password_hash_type = $3,
    password_hash = $4,
    organization_id = $5,
    verified = $6,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $7
`
