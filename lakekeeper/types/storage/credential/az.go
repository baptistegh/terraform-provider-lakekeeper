package credential

import "encoding/json"

type AZCredentialSettings interface {
	GetAZCredentialType() AZSCredentialType

	CredentialSettings
}

type AZSCredentialType string

const (
	ClientCredentials   = "client-credentials"
	SharedAccessKey     = "shared-access-key"
	AzureSystemIdentity = "azure-system-identity"
)

// verify implementations
var (
	_ AZCredentialSettings = (*AZCredentialClientCredentials)(nil)
	_ AZCredentialSettings = (*AZCredentialSharedAccessKey)(nil)
	_ AZCredentialSettings = (*AZCredentialManagedIdentity)(nil)

	_ CredentialSettings = (*AZCredentialClientCredentials)(nil)
	_ CredentialSettings = (*AZCredentialSharedAccessKey)(nil)
	_ CredentialSettings = (*AZCredentialManagedIdentity)(nil)
)

type AZCredentialClientCredentials struct {
	ClientID     string `json:"client-id"`
	ClientSecret string `json:"client-secret"`
	TenantID     string `json:"tenant-id"`
}

func NewAZCredentialClientCredentials(clientID, clientSecret, tenantID string) *AZCredentialClientCredentials {
	return &AZCredentialClientCredentials{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TenantID:     tenantID,
	}
}

type AZCredentialSharedAccessKey struct {
	Key string `json:"key"`
}

func NewAZCredentialSharedAccessKey(key string) *AZCredentialSharedAccessKey {
	return &AZCredentialSharedAccessKey{
		Key: key,
	}
}

type AZCredentialManagedIdentity struct{}

func NewAZCredentialManagedIdentity() *AZCredentialManagedIdentity {
	return &AZCredentialManagedIdentity{}
}

func (*AZCredentialClientCredentials) GetCredentialFamily() CredentialFamily {
	return AZCredentialFamily
}
func (*AZCredentialClientCredentials) GetAZCredentialType() AZSCredentialType {
	return ClientCredentials
}
func (c *AZCredentialClientCredentials) AsCredential() StorageCredential {
	return StorageCredential{Settings: c}
}
func (s AZCredentialClientCredentials) MarshalJSON() ([]byte, error) {
	type Alias AZCredentialClientCredentials
	aux := struct {
		Type           string `json:"type"`
		CredentialType string `json:"credential-type"`
		Alias
	}{
		Type:           string(s.GetCredentialFamily()),
		CredentialType: string(s.GetAZCredentialType()),
		Alias:          Alias(s),
	}
	return json.Marshal(aux)
}

func (*AZCredentialSharedAccessKey) GetCredentialFamily() CredentialFamily {
	return AZCredentialFamily
}
func (*AZCredentialSharedAccessKey) GetAZCredentialType() AZSCredentialType {
	return SharedAccessKey
}
func (c *AZCredentialSharedAccessKey) AsCredential() StorageCredential {
	return StorageCredential{Settings: c}
}
func (s AZCredentialSharedAccessKey) MarshalJSON() ([]byte, error) {
	type Alias AZCredentialSharedAccessKey
	aux := struct {
		Type           string `json:"type"`
		CredentialType string `json:"credential-type"`
		Alias
	}{
		Type:           string(s.GetCredentialFamily()),
		CredentialType: string(s.GetAZCredentialType()),
		Alias:          Alias(s),
	}
	return json.Marshal(aux)
}

func (*AZCredentialManagedIdentity) GetCredentialFamily() CredentialFamily {
	return AZCredentialFamily
}
func (*AZCredentialManagedIdentity) GetAZCredentialType() AZSCredentialType {
	return AzureSystemIdentity
}
func (c *AZCredentialManagedIdentity) AsCredential() StorageCredential {
	return StorageCredential{Settings: c}
}
func (s AZCredentialManagedIdentity) MarshalJSON() ([]byte, error) {
	aux := struct {
		Type           string `json:"type"`
		CredentialType string `json:"credential-type"`
	}{
		Type:           string(s.GetCredentialFamily()),
		CredentialType: string(s.GetAZCredentialType()),
	}
	return json.Marshal(aux)
}
