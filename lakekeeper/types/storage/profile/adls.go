package profile

import "encoding/json"

// ADLSStorageSettings represents the storage settings for a warehouse
// where data are stored on Azure Data Lake Storage.
type ADLSStorageSettings struct {
	// Name of the azure storage account.
	AccountName string `json:"account-name"`
	// Name of the adls filesystem, in blobstorage also known as container.
	Filesystem string `json:"filesystem"`
	// Allow alternative protocols such as wasbs:// in locations.
	// This is disabled by default. We do not recommend to use this setting
	// except for migration of old tables via the register endpoint.
	AllowAlternativeProtocols *bool `json:"allow-alternative-protocols,omitempty"`
	// The authority host to use for authentication.
	// Default: https://login.microsoftonline.com.
	AuthorityHost *string `json:"authority-host,omitempty"`
	// The host to use for the storage account. Default: dfs.core.windows.net.
	Host *string `json:"host,omitempty"`
	// Subpath in the filesystem to use.
	KeyPrefix *string `json:"key-prefix,omitempty"`
	// The validity of the sas token in seconds. Default: 3600.
	SASTokenValiditySeconds *int64 `json:"sas-token-validity-seconds,omitempty"`
}

func (sp *ADLSStorageSettings) GetStorageFamily() StorageFamily {
	return StorageFamilyADLS
}

type ADLSStorageSettingsOptions func(*ADLSStorageSettings) error

// NewADLSStorageSettings creates a new ADLS storage profile considering
// the options given.
func NewADLSStorageSettings(accountName, fs string, opts ...ADLSStorageSettingsOptions) (*ADLSStorageSettings, error) {
	// Default configuration
	profile := ADLSStorageSettings{
		AccountName: accountName,
		Filesystem:  fs,
	}

	// Apply options
	for _, v := range opts {
		if err := v(&profile); err != nil {
			return nil, err
		}
	}

	return &profile, nil
}

func WithADLSAlternativeProtocols() ADLSStorageSettingsOptions {
	return func(sp *ADLSStorageSettings) error {
		activated := true
		sp.AllowAlternativeProtocols = &activated
		return nil
	}
}

func WithAuthorityHost(host string) ADLSStorageSettingsOptions {
	return func(sp *ADLSStorageSettings) error {
		sp.AuthorityHost = &host
		return nil
	}
}

func WithADLSKeyPrefix(prefix string) ADLSStorageSettingsOptions {
	return func(sp *ADLSStorageSettings) error {
		sp.KeyPrefix = &prefix
		return nil
	}
}

func WithSASTokenValiditySeconds(seconds int64) ADLSStorageSettingsOptions {
	return func(sp *ADLSStorageSettings) error {
		sp.SASTokenValiditySeconds = &seconds
		return nil
	}
}

func WithHost(host string) ADLSStorageSettingsOptions {
	return func(sp *ADLSStorageSettings) error {
		sp.Host = &host
		return nil
	}
}

func (s *ADLSStorageSettings) AsProfile() StorageProfile {
	return StorageProfile{s}
}

func (s ADLSStorageSettings) MarshalJSON() ([]byte, error) {
	type Alias ADLSStorageSettings
	aux := struct {
		Type string `json:"type"`
		Alias
	}{
		Type:  string(StorageFamilyADLS),
		Alias: Alias(s),
	}
	return json.Marshal(aux)
}
