package profile

import (
	"encoding/json"
	"fmt"
)

type StorageProfileSettings interface {
	GetStorageProfileType() StorageFamily
	AsProfile() *StorageProfile

	json.Marshaler
}

type StorageProfile struct {
	StorageProfile StorageProfileSettings
}

type StorageFamily string

const (
	StorageFamilyADLS StorageFamily = "adls"
	StorageFamilyGCS  StorageFamily = "gcs"
	StorageFamilyS3   StorageFamily = "s3"
)

// Check the implementation
var (
	_ StorageProfileSettings = (*StorageProfileADLS)(nil)
	_ StorageProfileSettings = (*StorageProfileGCS)(nil)
	_ StorageProfileSettings = (*StorageProfileS3)(nil)
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
		var cfg StorageProfileS3
		if err := json.Unmarshal(data, &cfg); err != nil {
			return err
		}
		sc.StorageProfile = &cfg
	case "adls":
		var cfg StorageProfileADLS
		if err := json.Unmarshal(data, &cfg); err != nil {
			return err
		}
		sc.StorageProfile = &cfg
	case "gcs":
		var cfg StorageProfileGCS
		if err := json.Unmarshal(data, &cfg); err != nil {
			return err
		}
		sc.StorageProfile = &cfg
	default:
		return fmt.Errorf("unsupported storage type: %s", peek.Type)
	}
	return nil
}

func (sc StorageProfile) MarshalJSON() ([]byte, error) {
	return json.Marshal(sc.StorageProfile)
}

// Type-safe helpers
func (sc StorageProfile) AsS3() (*StorageProfileS3, bool) {
	cfg, ok := sc.StorageProfile.(*StorageProfileS3)
	return cfg, ok
}

func (sc StorageProfile) AsADLS() (*StorageProfileADLS, bool) {
	cfg, ok := sc.StorageProfile.(*StorageProfileADLS)
	return cfg, ok
}

func (sc StorageProfile) AsGCS() (*StorageProfileGCS, bool) {
	cfg, ok := sc.StorageProfile.(*StorageProfileGCS)
	return cfg, ok
}
