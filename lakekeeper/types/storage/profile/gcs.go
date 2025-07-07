package profile

import (
	"encoding/json"
)

// StorageProfileGCS represents the storage settings for a warehouse
// where data are stored on Google Cloud Storage.
type StorageProfileGCS struct {
	// Name of the GCS bucket
	Bucket string `json:"bucket"`
	// Subpath in the bucket to use.
	KeyPrefix *string `json:"key-prefix,omitempty"`
}

type StorageProfileGCSOptions func(*StorageProfileGCS) error

func (sp *StorageProfileGCS) GetStorageProfileType() StorageFamily {
	return StorageFamilyADLS
}

// NewStorageProfileGCS creates a new GCS storage profile considering
// the options given.
func NewStorageProfileGCS(bucket string, opts ...StorageProfileGCSOptions) (*StorageProfileGCS, error) {
	// Default configuration
	profile := StorageProfileGCS{
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

func WithGCSKeyPrefix(prefix string) StorageProfileGCSOptions {
	return func(sp *StorageProfileGCS) error {
		sp.KeyPrefix = &prefix
		return nil
	}
}

func (s *StorageProfileGCS) AsProfile() *StorageProfile {
	return &StorageProfile{s}
}

func (s StorageProfileGCS) MarshalJSON() ([]byte, error) {
	type Alias StorageProfileGCS
	aux := struct {
		Type string `json:"type"`
		Alias
	}{
		Type:  string(StorageFamilyADLS),
		Alias: Alias(s),
	}
	return json.Marshal(aux)
}
