package dao

type Session struct {
	OrganizationID  string `json:"organization_id"`
	AdministratorID string `json:"administrator_id"`
}

// SessionDatastore defines the interface for session operations in a key-value datastore
type SessionDatastore interface {
	Create(session *Session) (string, error)
	Read(token string) (*Session, error)
	Delete(token string) error
}
