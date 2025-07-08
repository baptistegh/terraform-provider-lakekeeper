package profile

import (
	"encoding/json"
)

// GCSStorageSettings represents the storage settings for a warehouse
// where data are stored on Google Cloud Storage.
type GCSStorageSettings struct {
	// Name of the GCS bucket
	Bucket string `json:"bucket"`
	// Subpath in the bucket to use.
	KeyPrefix *string `json:"key-prefix,omitempty"`
}

type GCSStorageSettingsOptions func(*GCSStorageSettings) error

func (sp *GCSStorageSettings) GetStorageFamily() StorageFamily {
	return StorageFamilyADLS
}

// NewGCSStorageSettings creates a new GCS storage profile considering
// the options given.
func NewGCSStorageSettings(bucket string, opts ...GCSStorageSettingsOptions) (*GCSStorageSettings, error) {
	// Default configuration
	profile := GCSStorageSettings{
		Bucket: bucket,
	}

	// Apply options
	for _, v := range opts {
		if err := v(&profile); err != nil {
			return nil, err
		}
	}

	return &profile, nil
}

func WithGCSKeyPrefix(prefix string) GCSStorageSettingsOptions {
	return func(sp *GCSStorageSettings) error {
		sp.KeyPrefix = &prefix
		return nil
	}
}

func (s *GCSStorageSettings) AsProfile() StorageProfile {
	return StorageProfile{s}
}

func (s GCSStorageSettings) MarshalJSON() ([]byte, error) {
	type Alias GCSStorageSettings
	aux := struct {
		Type string `json:"type"`
		Alias
	}{
		Type:  string(StorageFamilyADLS),
		Alias: Alias(s),
	}
	return json.Marshal(aux)
}
