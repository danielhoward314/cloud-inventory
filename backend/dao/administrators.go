package dao

type Administrator struct {
	ID               string `json:"id"`
	Email            string `json:"email"`
	DisplayName      string `json:"display_name"`
	PasswordHashType string `json:"password_hash_type"`
	PasswordHash     string `json:"password_hash"`
	OrganizationID   string `json:"organization_id"`
}

type Administrators interface {
	Create(administrator *Administrator) (string, error)
	// Read(id string) (*Administrator, error)
	// Update(*Administrator) (*Administrator, error)
	// Delete(id string) (*Administrator, error)
}
