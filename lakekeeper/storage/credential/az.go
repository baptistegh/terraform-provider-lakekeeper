package credential

type AZCredential interface {
	GetAZCredentialType() string
}

type AZCredentialClientCredentials struct {
	Type           string `json:"type"`
	CredentialType string `json:"credential-type"`
	ClientID       string `json:"client_id"`
	ClientSecret   string `json:"client_secret"`
	TenantID       string `json:"tenant_id"`
}

type AZCredentialManagedIdentity struct {
	Type           string `json:"type"`
	CredentialType string `json:"credential-type"`
	Key            string `json:"key"`
}

type AZCredentialSharedAccessKey struct {
	Type           string `json:"type"`
	CredentialType string `json:"credential-type"`
	Key            string `json:"key"`
}

// validate implementations
var (
	_ AZCredential      = &AZCredentialClientCredentials{}
	_ StorageCredential = &AZCredentialClientCredentials{}
	_ AZCredential      = &AZCredentialManagedIdentity{}
	_ StorageCredential = &AZCredentialManagedIdentity{}
	_ AZCredential      = &AZCredentialSharedAccessKey{}
	_ StorageCredential = &AZCredentialSharedAccessKey{}
)

func (AZCredentialClientCredentials) GetStorageCredentialType() string {
	return "az"
}

func (AZCredentialClientCredentials) GetAZCredentialType() string {
	return "client_credentials"
}

func (AZCredentialSharedAccessKey) GetStorageCredentialType() string {
	return "az"
}

func (AZCredentialSharedAccessKey) GetAZCredentialType() string {
	return "shared-access-key"
}

func (AZCredentialManagedIdentity) GetStorageCredentialType() string {
	return "az"
}

func (AZCredentialManagedIdentity) GetAZCredentialType() string {
	return "azure-system-identity"
}
