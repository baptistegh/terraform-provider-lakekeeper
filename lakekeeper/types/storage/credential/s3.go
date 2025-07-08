package credential

import "encoding/json"

type S3SCredentialSettings interface {
	GetS3CredentialType() S3CredentialType

	CredentialSettings
}

type S3CredentialType string

const (
	AccessKey         S3CredentialType = "access-key"
	AWSSystemIdentity S3CredentialType = "aws-system-identity"
	CloudflareR2      S3CredentialType = "cloudflare-r2"
)

// verify implementations
var (
	_ S3SCredentialSettings = (*S3CredentialAccessKey)(nil)
	_ S3SCredentialSettings = (*S3CredentialSystemIdentity)(nil)
	_ S3SCredentialSettings = (*CloudflareR2Credential)(nil)

	_ CredentialSettings = (*S3CredentialAccessKey)(nil)
	_ CredentialSettings = (*S3CredentialSystemIdentity)(nil)
	_ CredentialSettings = (*CloudflareR2Credential)(nil)
)

type S3CredentialAccessKey struct {
	// Access key ID used for IO operations of Lakekeeper
	AWSAccessKeyID string `json:"aws-access-key-id"`
	// Secret key associated with the access key ID.
	AWSSecretAccessKey string  `json:"aws-secret-access-key"`
	ExternalID         *string `json:"external-id"`
}

type S3CredentialAccessKeyOptions func(*S3CredentialAccessKey) error

func NewS3CredentialAccessKey(accessKey, secretKey string, options ...S3CredentialAccessKeyOptions) (*S3CredentialAccessKey, error) {
	s := S3CredentialAccessKey{
		AWSAccessKeyID:     accessKey,
		AWSSecretAccessKey: secretKey,
	}

	for _, v := range options {
		err := v(&s)
		if err != nil {
			return nil, err
		}
	}

	return &s, nil
}

type S3CredentialSystemIdentity struct {
	ExternalID string `json:"external-id"`
}

func NewS3CredentialSystemIdentity(externalID string) *S3CredentialSystemIdentity {
	return &S3CredentialSystemIdentity{
		ExternalID: externalID,
	}
}

type CloudflareR2Credential struct {
	// Access key ID used for IO operations of Lakekeeper
	AccessKeyID string `json:"access-key-id"`
	// Secret key associated with the access key ID.
	SecretAccessKey string `json:"secret-access-key"`
	// Cloudflare account ID, used to determine the temporary credentials
	// endpoint.
	AccountID string `json:"account=id"`
	// Token associated with the access key ID.
	// This is used to fetch downscoped temporary credentials
	// for vended credentials.
	Token string `json:"token"`
}

func NewCloudflareR2Credential(accessKey, secretKey, accountID, token string) *CloudflareR2Credential {
	return &CloudflareR2Credential{
		AccessKeyID:     accessKey,
		SecretAccessKey: secretKey,
		AccountID:       accountID,
		Token:           token,
	}
}

func WithExternalID(externalID string) S3CredentialAccessKeyOptions {
	return func(c *S3CredentialAccessKey) error {
		c.ExternalID = &externalID
		return nil
	}
}

func (*S3CredentialAccessKey) GetCredentialFamily() CredentialFamily {
	return S3CredentialFamily
}
func (*S3CredentialAccessKey) GetS3CredentialType() S3CredentialType {
	return AccessKey
}
func (c *S3CredentialAccessKey) AsCredential() StorageCredential {
	return StorageCredential{Settings: c}
}
func (s S3CredentialAccessKey) MarshalJSON() ([]byte, error) {
	type Alias S3CredentialAccessKey
	aux := struct {
		Type           string `json:"type"`
		CredentialType string `json:"credential-type"`
		Alias
	}{
		Type:           string(S3CredentialFamily),
		CredentialType: string(AccessKey),
		Alias:          Alias(s),
	}
	return json.Marshal(aux)
}

func (*S3CredentialSystemIdentity) GetCredentialFamily() CredentialFamily {
	return S3CredentialFamily
}
func (*S3CredentialSystemIdentity) GetS3CredentialType() S3CredentialType {
	return AWSSystemIdentity
}
func (c *S3CredentialSystemIdentity) AsCredential() StorageCredential {
	return StorageCredential{Settings: c}
}
func (s S3CredentialSystemIdentity) MarshalJSON() ([]byte, error) {
	type Alias S3CredentialSystemIdentity
	aux := struct {
		Type           string `json:"type"`
		CredentialType string `json:"credential-type"`
		Alias
	}{
		Type:           string(S3CredentialFamily),
		CredentialType: string(AWSSystemIdentity),
		Alias:          Alias(s),
	}
	return json.Marshal(aux)
}

func (*CloudflareR2Credential) GetCredentialFamily() CredentialFamily {
	return S3CredentialFamily
}
func (*CloudflareR2Credential) GetS3CredentialType() S3CredentialType {
	return CloudflareR2
}
func (c *CloudflareR2Credential) AsCredential() StorageCredential {
	return StorageCredential{Settings: c}
}
func (s CloudflareR2Credential) MarshalJSON() ([]byte, error) {
	type Alias CloudflareR2Credential
	aux := struct {
		Type           string `json:"type"`
		CredentialType string `json:"credential-type"`
		Alias
	}{
		Type:           string(S3CredentialFamily),
		CredentialType: string(CloudflareR2),
		Alias:          Alias(s),
	}
	return json.Marshal(aux)
}
