package credential

import (
	"encoding/json"
	"fmt"
)

type CredentialSettings interface {
	GetCredentialFamily() CredentialFamily
	AsCredential() StorageCredential

	json.Marshaler
}

type StorageCredential struct {
	Settings CredentialSettings
}

type CredentialFamily string

const (
	S3CredentialFamily  CredentialFamily = "s3"
	GCSCredentialFamily CredentialFamily = "gcs"
	AZCredentialFamily  CredentialFamily = "az"
)

func (sc *StorageCredential) UnmarshalJSON(data []byte) error {
	var peek struct {
		Type           string `json:"type"`
		CredentialType string `json:"credential-type"`
	}
	if err := json.Unmarshal(data, &peek); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}

	switch fmt.Sprintf("%s:%s", peek.Type, peek.CredentialType) {
	case fmt.Sprintf("%s:%s", S3CredentialFamily, peek.CredentialType):
		var cfg S3CredentialAccessKey
		if err := json.Unmarshal(data, &cfg); err != nil {
			return err
		}
		sc.Settings = &cfg
	case fmt.Sprintf("%s:%s", S3CredentialFamily, peek.CredentialType):
		var cfg S3CredentialSystemIdentity
		if err := json.Unmarshal(data, &cfg); err != nil {
			return err
		}
		sc.Settings = &cfg
	case fmt.Sprintf("%s:%s", S3CredentialFamily, peek.CredentialType):
		var cfg CloudflareR2Credential
		if err := json.Unmarshal(data, &cfg); err != nil {
			return err
		}
		sc.Settings = &cfg
	case fmt.Sprintf("%s:%s", GCSCredentialFamily, ServiceAccountKey):
		var cfg GCSCredentialServiceAccountKey
		if err := json.Unmarshal(data, &cfg); err != nil {
			return err
		}
		sc.Settings = &cfg
	case fmt.Sprintf("%s:%s", GCSCredentialFamily, GCPSystemIdentity):
		var cfg GCSCredentialSystemIdentity
		if err := json.Unmarshal(data, &cfg); err != nil {
			return err
		}
		sc.Settings = &cfg
	case fmt.Sprintf("%s:%s", AZCredentialFamily, ClientCredentials):
		var cfg AZCredentialClientCredentials
		if err := json.Unmarshal(data, &cfg); err != nil {
			return err
		}
		sc.Settings = &cfg
	case fmt.Sprintf("%s:%s", AZCredentialFamily, SharedAccessKey):
		var cfg AZCredentialSharedAccessKey
		if err := json.Unmarshal(data, &cfg); err != nil {
			return err
		}
		sc.Settings = &cfg
	case fmt.Sprintf("%s:%s", AZCredentialFamily, AzureSystemIdentity):
		var cfg AZCredentialManagedIdentity
		if err := json.Unmarshal(data, &cfg); err != nil {
			return err
		}
		sc.Settings = &cfg
	default:
		return fmt.Errorf("unsupported storage credential type: %s / %s", peek.Type, peek.CredentialType)
	}
	return nil
}

func (sc StorageCredential) MarshalJSON() ([]byte, error) {
	return json.Marshal(sc.Settings)
}

// Type-safe helpers
func (sc StorageCredential) AsS3() (S3SCredentialSettings, bool) {
	cfg, ok := sc.Settings.(S3SCredentialSettings)
	return cfg, ok
}

func (sc StorageCredential) AsAZ() (AZCredentialSettings, bool) {
	cfg, ok := sc.Settings.(AZCredentialSettings)
	return cfg, ok
}

func (sc StorageCredential) AsGCS() (GCSSCredentialSettings, bool) {
	cfg, ok := sc.Settings.(GCSSCredentialSettings)
	return cfg, ok
}
