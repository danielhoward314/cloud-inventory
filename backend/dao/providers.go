package dao

type Provider struct {
	ID                 string `json:"id"`
	ExternalIdentifier string `json:"external_identifier"`
	Name               string `json:"name"`
	ProviderName       string `json:"provider_name"`
	Metadata           []byte `json:"metadata"`
	OrganizationID     string `json:"organization_id"`
}

type Providers interface {
	List(organizationID string) ([]*Provider, error)
}
