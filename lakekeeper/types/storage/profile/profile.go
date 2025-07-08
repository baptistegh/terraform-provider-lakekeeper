package profile

import (
	"encoding/json"
	"fmt"
)

type StorageSettings interface {
	GetStorageFamily() StorageFamily
	AsProfile() StorageProfile

	json.Marshaler
}

type StorageProfile struct {
	StorageSettings StorageSettings
}

type StorageFamily string

const (
	StorageFamilyADLS StorageFamily = "adls"
	StorageFamilyGCS  StorageFamily = "gcs"
	StorageFamilyS3   StorageFamily = "s3"
)

// Check the implementation
var (
	_ StorageSettings = (*ADLSStorageSettings)(nil)
	_ StorageSettings = (*GCSStorageSettings)(nil)
	_ StorageSettings = (*S3StorageSettings)(nil)
)

func (sc *StorageProfile) UnmarshalJSON(data []byte) error {
	var peek struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &peek); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}

	switch peek.Type {
	case "s3":
		var cfg S3StorageSettings
		if err := json.Unmarshal(data, &cfg); err != nil {
			return err
		}
		sc.StorageSettings = &cfg
	case "adls":
		var cfg ADLSStorageSettings
		if err := json.Unmarshal(data, &cfg); err != nil {
			return err
		}
		sc.StorageSettings = &cfg
	case "gcs":
		var cfg GCSStorageSettings
		if err := json.Unmarshal(data, &cfg); err != nil {
			return err
		}
		sc.StorageSettings = &cfg
	default:
		return fmt.Errorf("unsupported storage type: %s", peek.Type)
	}
	return nil
}

func (sc StorageProfile) MarshalJSON() ([]byte, error) {
	return json.Marshal(sc.StorageSettings)
}

// Type-safe helpers
func (sc StorageProfile) AsS3() (*S3StorageSettings, bool) {
	cfg, ok := sc.StorageSettings.(*S3StorageSettings)
	return cfg, ok
}

func (sc StorageProfile) AsADLS() (*ADLSStorageSettings, bool) {
	cfg, ok := sc.StorageSettings.(*ADLSStorageSettings)
	return cfg, ok
}

func (sc StorageProfile) AsGCS() (*GCSStorageSettings, bool) {
	cfg, ok := sc.StorageSettings.(*GCSStorageSettings)
	return cfg, ok
}
