package storage

import (
	"encoding/json"
	"fmt"
)

type StorageCredentialWrapper struct {
	StorageCredential StorageCredential
}

type StorageCredential interface {
	GetStorageCredentialType() string
}

func (w *StorageCredentialWrapper) UnmarshalJSON(data []byte) error {
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	t, _ := raw["type"].(string)
	c, _ := raw["credential-type"].(string)
	switch t {
	case "s3":
		switch c {
		case "access-key":
			var sc S3CredentialAccessKey
			if err := json.Unmarshal(data, &sc); err != nil {
				return err
			}
			w.StorageCredential = sc
		case "aws-system-identity":
			var sc S3CredentialAccessKey
			if err := json.Unmarshal(data, &sc); err != nil {
				return err
			}
			w.StorageCredential = sc
		case "cloudflare-r2":
			var sc CloudflareR2Credential
			if err := json.Unmarshal(data, &sc); err != nil {
				return err
			}
			w.StorageCredential = sc
		default:
			return fmt.Errorf("unknown credential-type for s3: %s", t)
		}
	case "az":
		switch c {
		case "client_credentials":
			var sc AZCredentialClientCredentials
			if err := json.Unmarshal(data, &sc); err != nil {
				return err
			}
			w.StorageCredential = sc
		case "shared-access-key":
			var sc AZCredentialSharedAccessKey
			if err := json.Unmarshal(data, &sc); err != nil {
				return err
			}
			w.StorageCredential = sc
		case "azure-system-identity":
			var sc AZCredentialManagedIdentity
			if err := json.Unmarshal(data, &sc); err != nil {
				return err
			}
			w.StorageCredential = sc
		default:
			return fmt.Errorf("unknown credential-type for az: %s", t)
		}
	case "gcs":
		switch c {
		case "service-account-key":
			var sc GCSCredentialServiceAccountKey
			if err := json.Unmarshal(data, &sc); err != nil {
				return err
			}
			w.StorageCredential = sc
		case "gcp-system-identity":
			var sc GCSCredentialSystemIdentity
			if err := json.Unmarshal(data, &sc); err != nil {
				return err
			}
			w.StorageCredential = sc
		default:
			return fmt.Errorf("unknown credential-type for gcs: %s", t)
		}
	default:
		return fmt.Errorf("unknown storage-profile type: %s", t)
	}
	return nil
}

func (w StorageCredentialWrapper) MarshalJSON() ([]byte, error) {
	if w.StorageCredential == nil {
		return []byte("null"), nil
	}
	return json.Marshal(w.StorageCredential)
}
