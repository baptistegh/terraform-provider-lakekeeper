package credential

type GCSCredential interface {
	GetGCSCredentialType() string
}

type GCSCredentialServiceAccountKey struct {
	Type           string `json:"type"`
	CredentialType string `json:"credential-type"`

	*GCSKey `json:"key"`
}

type GCSKey struct {
	AuthProviderX509CertURL string `json:"auth_provider_x509_cert_url"`
	AuthURI                 string `json:"auth_uri"`
	ClientEmail             string `json:"client_email"`
	ClientID                string `json:"client_id"`
	ClientX509CertURL       string `json:"client_x509_cert_url"`
	PrivateKey              string `json:"private_key"`
	PrivateKeyID            string `json:"private_key_id"`
	ProjectID               string `json:"project_id"`
	TokenURI                string `json:"token_uri"`
	Type                    string `json:"type"`
	UniverseDomain          string `json:"universe_domain"`
}

type GCSCredentialSystemIdentity struct {
	Type           string `json:"type"`
	CredentialType string `json:"credential-type"`
}

// validate implementations
var (
	_ GCSCredential     = &GCSCredentialServiceAccountKey{}
	_ StorageCredential = &GCSCredentialServiceAccountKey{}
	_ GCSCredential     = &GCSCredentialSystemIdentity{}
	_ GCSCredential     = &GCSCredentialSystemIdentity{}
)

func (GCSCredentialServiceAccountKey) GetStorageCredentialType() string {
	return "gcs"
}

func (GCSCredentialServiceAccountKey) GetGCSCredentialType() string {
	return "service-account-key"
}

func NewGCSCredentialServiceACcountKey(key *GCSKey) GCSCredentialServiceAccountKey {
	return GCSCredentialServiceAccountKey{
		Type:           "gcs",
		CredentialType: "service-account-key",
		GCSKey:         key,
	}
}

func (GCSCredentialSystemIdentity) GetStorageCredentialType() string {
	return "gcs"
}

func (GCSCredentialSystemIdentity) GetGCSCredentialType() string {
	return "gcp-system-identity"
}

func NewGCSCredentialSystemIdentity() GCSCredentialSystemIdentity {
	return GCSCredentialSystemIdentity{
		Type:           "gcs",
		CredentialType: "gcp-system-identity",
	}
}
