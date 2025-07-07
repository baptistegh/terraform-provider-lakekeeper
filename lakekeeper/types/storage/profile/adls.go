package profile

import "encoding/json"

// StorageProfileADLS represents the storage settings for a warehouse
// where data are stored on Azure Data Lake Storage.
type StorageProfileADLS struct {
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

func (sp *StorageProfileADLS) GetStorageProfileType() StorageFamily {
	return StorageFamilyADLS
}

type StorageProfileADLSOptions func(*StorageProfileADLS) error

// NewStorageProfileADLS creates a new ADLS storage profile considering
// the options given.
func NewStorageProfileADLS(accountName, fs string, opts ...StorageProfileADLSOptions) (*StorageProfileADLS, error) {
	// Default configuration
	profile := StorageProfileADLS{
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

func WithADLSAlternativeProtocols() StorageProfileADLSOptions {
	return func(sp *StorageProfileADLS) error {
		activated := true
		sp.AllowAlternativeProtocols = &activated
		return nil
	}
}

func WithAuthorityHost(host string) StorageProfileADLSOptions {
	return func(sp *StorageProfileADLS) error {
		sp.AuthorityHost = &host
		return nil
	}
}

func WithADLSKeyPrefix(prefix string) StorageProfileADLSOptions {
	return func(sp *StorageProfileADLS) error {
		sp.KeyPrefix = &prefix
		return nil
	}
}

func WithSASTokenValiditySeconds(seconds int64) StorageProfileADLSOptions {
	return func(sp *StorageProfileADLS) error {
		sp.SASTokenValiditySeconds = &seconds
		return nil
	}
}

func WithHost(host string) StorageProfileADLSOptions {
	return func(sp *StorageProfileADLS) error {
		sp.Host = &host
		return nil
	}
}

func (s *StorageProfileADLS) AsProfile() *StorageProfile {
	return &StorageProfile{s}
}

func (s StorageProfileADLS) MarshalJSON() ([]byte, error) {
	type Alias StorageProfileADLS
	aux := struct {
		Type string `json:"type"`
		Alias
	}{
		Type:  string(StorageFamilyADLS),
		Alias: Alias(s),
	}
	return json.Marshal(aux)
}
