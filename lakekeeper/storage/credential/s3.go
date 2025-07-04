package credential

type S3Credential interface {
	GetS3CredentialType() string
}

type S3CredentialAccessKey struct {
	Type               string `json:"type"`
	CredentialType     string `json:"credential-type"`
	AWSAccessKeyID     string `json:"aws-access-key-id"`
	AWSSecretAccessKey string `json:"aws-secret-access-key"`
	ExternalID         string `json:"external-id,omitempty"`
}

type S3CredentialSystemIdentity struct {
	Type           string `json:"type"`
	CredentialType string `json:"credential-type"`
	ExternalID     string `json:"external-id,omitempty"`
}

type CloudflareR2Credential struct {
	Type            string `json:"type"`
	CredentialType  string `json:"credential-type"`
	AccessKeyID     string `json:"access-key-id"`
	SecretAccessKey string `json:"secret-access-key"`
	Token           string `json:"token"`
	AccountID       string `json:"account-id"`
}

// validate implementations
var (
	_ S3Credential      = &S3CredentialAccessKey{}
	_ StorageCredential = &S3CredentialAccessKey{}
	_ S3Credential      = &S3CredentialSystemIdentity{}
	_ StorageCredential = &S3CredentialSystemIdentity{}
	_ S3Credential      = &CloudflareR2Credential{}
	_ StorageCredential = &CloudflareR2Credential{}
)

func NewS3CredentialAccessKey(accessKeyID, secretAccessKey, externalID string) S3CredentialAccessKey {
	return S3CredentialAccessKey{
		Type:               "s3",
		CredentialType:     "access-key",
		AWSAccessKeyID:     accessKeyID,
		AWSSecretAccessKey: secretAccessKey,
		ExternalID:         externalID,
	}
}

func (S3CredentialAccessKey) GetStorageCredentialType() string {
	return "s3"
}

func (S3CredentialAccessKey) GetS3CredentialType() string {
	return "access-key"
}

func NewS3CredentialSystemIdentity(externalID string) S3CredentialSystemIdentity {
	return S3CredentialSystemIdentity{
		Type:           "s3",
		CredentialType: "aws-system-identity",
		ExternalID:     externalID,
	}
}

func (S3CredentialSystemIdentity) GetStorageCredentialType() string {
	return "s3"
}

func (S3CredentialSystemIdentity) GetS3CredentialType() string {
	return "aws-system-identity"
}

func NewCloudflareR2Credential(accessKeyID, secretAccessKey, token, accountID string) CloudflareR2Credential {
	return CloudflareR2Credential{
		Type:            "s3",
		CredentialType:  "cloudflare-r2",
		AccessKeyID:     accessKeyID,
		SecretAccessKey: secretAccessKey,
		Token:           token,
		AccountID:       accountID,
	}
}

func (CloudflareR2Credential) GetStorageCredentialType() string {
	return "s3"
}

func (CloudflareR2Credential) GetS3CredentialType() string {
	return "cloudflare-r2"
}
