package storage

import (
	"encoding/json"
	"fmt"
)

var ValidStorageProfileTypes = []string{"adls", "s3", "gcs"}

// validate implementations
var (
	_ StorageProfile = &ADLSStorageSettings{}
	_ StorageProfile = &GCSStorageSettings{}
	_ StorageProfile = &StorageProfileS3{}
)

type StorageProfile interface {
	GetStorageType() string
}

type ADLSStorageSettings struct {
	Type                      string  `json:"type"`
	AccountName               string  `json:"account-name"`
	AllowAlternativeProtocols bool    `json:"allow-alternative-protocols,omitempty"`
	AuthorityHost             *string `json:"authority-host,omitempty"`
	Filesystem                string  `json:"filesystem"`
	Host                      *string `json:"host,omitempty"`
	KeyPrefix                 *string `json:"key-prefix,omitempty"`
	SASTokenValiditySeconds   *int64  `json:"sas-token-validity-seconds,omitempty"`
}

func (ADLSStorageSettings) GetStorageType() string {
	return "adls"
}

func NewADLSStorageSettings(accountName, fileSystem string) *ADLSStorageSettings {
	return &ADLSStorageSettings{
		Type:                      "adls",
		AccountName:               accountName,
		Filesystem:                fileSystem,
		AllowAlternativeProtocols: false,
	}
}

type StorageProfileS3 struct {
	Type                      string  `json:"type"`
	AllowAlternativeProtocols bool    `json:"allow-alternative-protocols,omitempty"`
	AssumeRoleARN             *string `json:"assume-role-arn,omitempty"`
	AWSKMSKeyARN              *string `json:"aws-kms-key-arn,omitempty"`
	Bucket                    string  `json:"bucket"`
	Endpoint                  *string `json:"endpoint,omitempty"`
	Flavor                    *string `json:"flavor,omitempty"`
	KeyPrefix                 *string `json:"key-prefix,omitempty"`
	PathStyleAccess           *bool   `json:"path-style-access,omitempty"`
	PushS3DeleteDisabled      *bool   `json:"push-s3-delete-disabled,omitempty"`
	Region                    string  `json:"region"`
	RemoteSigningURLStyle     *string `json:"remote-signing-url-style,omitempty"`
	STSEnabled                bool    `json:"sts-enabled"`
	STSRoleARN                *string `json:"sts-role-arn,omitempty"`
	STSTokenValiditySeconds   *int64  `json:"sts-token-validity-seconds,omitempty"`
}

func (StorageProfileS3) GetStorageType() string {
	return "s3"
}

func NewStorageProfileS3(bucket, region string, stsEnabled bool) *StorageProfileS3 {
	return &StorageProfileS3{
		Type:       "s3",
		Bucket:     bucket,
		Region:     region,
		STSEnabled: stsEnabled,
	}
}

type GCSStorageSettings struct {
	Type      string  `json:"type"`
	Bucket    string  `json:"bucket"`
	KeyPrefix *string `json:"key-prefix,omitempty"`
}

func (s GCSStorageSettings) GetStorageType() string {
	return "gcs"
}

func NewGCSStorageSettings(bucket string) *GCSStorageSettings {
	return &GCSStorageSettings{
		Type:   "gcs",
		Bucket: bucket,
	}
}

type StorageProfileWrapper struct {
	StorageProfile StorageProfile
}

func (w *StorageProfileWrapper) UnmarshalJSON(data []byte) error {
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	t, _ := raw["type"].(string)
	switch t {
	case "adls":
		var s ADLSStorageSettings
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}
		w.StorageProfile = s
	case "s3":
		var s StorageProfileS3
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}
		w.StorageProfile = s
	case "gcs":
		var s GCSStorageSettings
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}
		w.StorageProfile = s
	default:
		return fmt.Errorf("unknown storage-profile type: %s", t)
	}
	return nil
}

func (w StorageProfileWrapper) MarshalJSON() ([]byte, error) {
	if w.StorageProfile == nil {
		return []byte("null"), nil
	}
	return json.Marshal(w.StorageProfile)
}
