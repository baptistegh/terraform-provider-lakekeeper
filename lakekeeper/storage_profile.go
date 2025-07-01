package lakekeeper

import (
	"encoding/json"
	"fmt"
)

type StorageProfile interface {
	IsStorageProfile()
}

type StorageProfileADLS struct {
	Type                      string `json:"type"`
	AccountName               string `json:"account-name"`
	AllowAlternativeProtocols bool   `json:"allow-alternative-protocols,omitempty"`
	AuthorityHost             string `json:"authority-host,omitempty"`
	Filesystem                string `json:"filesystem"`
	Host                      string `json:"host,omitempty"`
	KeyPrefix                 string `json:"key-prefix,omitempty"`
	SASTokenValiditySeconds   int    `json:"sas-token-validity-seconds,omitempty"`
}

func (StorageProfileADLS) IsStorageProfile() {}

type StorageProfileS3 struct {
	Type                      string `json:"type"`
	AllowAlternativeProtocols bool   `json:"allow-alternative-protocols,omitempty"`
	AssumeRoleARN             string `json:"assume-role-arn,omitempty"`
	AWSKMSKeyARN              string `json:"aws-kms-key-arn,omitempty"`
	Bucket                    string `json:"bucket"`
	Endpoint                  string `json:"endpoint,omitempty"`
	Flavor                    string `json:"flavor,omitempty"`
	KeyPrefix                 string `json:"key-prefix,omitempty"`
	PathStyleAccess           bool   `json:"path-style-access,omitempty"`
	PushS3DeleteDisabled      bool   `json:"push-s3-delete-disabled,omitempty"`
	Region                    string `json:"region"`
	RemoteSigningURLStyle     string `json:"remote-signing-url-style,omitempty"`
	STSEnabled                bool   `json:"sts-enabled"`
	STSRoleARN                string `json:"sts-role-arn,omitempty"`
	STSTokenValiditySeconds   int    `json:"sts-token-validity-seconds,omitempty"`
}

func (StorageProfileS3) IsStorageProfile() {}

type StorageProfileGCS struct {
	Type      string `json:"type"`
	Bucket    string `json:"bucket"`
	KeyPrefix string `json:"key-prefix,omitempty"`
}

func (s StorageProfileGCS) GetType() string {
	return "gcs"
}

func (StorageProfileGCS) IsStorageProfile() {}

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
		var s StorageProfileADLS
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
		var s StorageProfileGCS
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
