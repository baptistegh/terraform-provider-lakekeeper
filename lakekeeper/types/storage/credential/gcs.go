package credential

import "encoding/json"

type GCSSCredentialSettings interface {
	GetGCSCredentialType() GCSCredentialType

	CredentialSettings
}

type GCSCredentialType string

const (
	ServiceAccountKey GCSCredentialType = "service-account-key"
	GCPSystemIdentity GCSCredentialType = "gcp-system-identity"
)

// verify implementations
var (
	_ GCSSCredentialSettings = (*GCSCredentialServiceAccountKey)(nil)
	_ GCSSCredentialSettings = (*GCSCredentialSystemIdentity)(nil)

	_ CredentialSettings = (*GCSCredentialServiceAccountKey)(nil)
	_ CredentialSettings = (*GCSCredentialSystemIdentity)(nil)
)

type GCSCredentialServiceAccountKey struct {
	Key GCSServiceKey `json:"key"`
}

type GCSServiceKey struct {
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

func NewGCSCredentialServiceAccountKey(key GCSServiceKey) *GCSCredentialServiceAccountKey {
	return &GCSCredentialServiceAccountKey{}
}

type GCSCredentialSystemIdentity struct{}

func NewGCSCredentialSystemIdentity() *GCSCredentialSystemIdentity {
	return &GCSCredentialSystemIdentity{}
}

func (*GCSCredentialServiceAccountKey) GetCredentialFamily() CredentialFamily {
	return GCSCredentialFamily
}
func (*GCSCredentialServiceAccountKey) GetGCSCredentialType() GCSCredentialType {
	return ServiceAccountKey
}
func (c *GCSCredentialServiceAccountKey) AsCredential() StorageCredential {
	return StorageCredential{Settings: c}
}
func (s GCSCredentialServiceAccountKey) MarshalJSON() ([]byte, error) {
	type Alias GCSCredentialServiceAccountKey
	aux := struct {
		Type           string `json:"type"`
		CredentialType string `json:"credential-type"`
		Alias
	}{
		Type:           string(s.GetCredentialFamily()),
		CredentialType: string(s.GetGCSCredentialType()),
		Alias:          Alias(s),
	}
	return json.Marshal(aux)
}

func (*GCSCredentialSystemIdentity) GetCredentialFamily() CredentialFamily {
	return GCSCredentialFamily
}
func (*GCSCredentialSystemIdentity) GetGCSCredentialType() GCSCredentialType {
	return GCPSystemIdentity
}
func (c *GCSCredentialSystemIdentity) AsCredential() StorageCredential {
	return StorageCredential{Settings: c}
}
func (s GCSCredentialSystemIdentity) MarshalJSON() ([]byte, error) {
	aux := struct {
		Type           string `json:"type"`
		CredentialType string `json:"credential-type"`
	}{
		Type:           string(s.GetCredentialFamily()),
		CredentialType: string(s.GetGCSCredentialType()),
	}
	return json.Marshal(aux)
}
