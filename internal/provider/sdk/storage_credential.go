package sdk

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/baptistegh/go-lakekeeper/pkg/apis/management/v1/storage/credential"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type (
	StorageCredsModel interface {
		AsSDK() (credential.CredentialSettings, error)
		IsEmpty() bool
	}

	S3AccessKeyCredsModel struct {
		AccessKeyID     types.String `tfsdk:"access_key_id"`
		SecretAccessKey types.String `tfsdk:"secret_access_key"`
		ExternalID      types.String `tfsdk:"external_id"`
	}

	AWSSystemIdentityCredsModel struct {
		ExternalID types.String `tfsdk:"external_id"`
	}

	CloudflareR2CredsModel struct {
		AccessKeyID     types.String `tfsdk:"access_key_id"`
		SecretAccessKey types.String `tfsdk:"secret_access_key"`
		AccountID       types.String `tfsdk:"account_id"`
		Token           types.String `tfsdk:"token"`
	}

	AZClientCredentialsCredsModel struct {
		ClientID     types.String `tfsdk:"client_id"`
		ClientSecret types.String `tfsdk:"client_secret"`
		TenantID     types.String `tfsdk:"tenant_id"`
	}

	AZSharedAccessKeyCredsModel struct {
		Key types.String `tfsdk:"key"`
	}

	GCSServiceAccountKeyCredsModel struct {
		Key types.String `tfsdk:"key"`
	}

	AzureSystemIdentityCredsModel struct{}
	GCPSystemIdentityCredsModel   struct{}
)

var (
	_ StorageCredsModel = (*S3AccessKeyCredsModel)(nil)
	_ StorageCredsModel = (*CloudflareR2CredsModel)(nil)
	_ StorageCredsModel = (*AWSSystemIdentityCredsModel)(nil)

	_ StorageCredsModel = (*AZClientCredentialsCredsModel)(nil)
	_ StorageCredsModel = (*AZSharedAccessKeyCredsModel)(nil)
	_ StorageCredsModel = (*AzureSystemIdentityCredsModel)(nil)

	_ StorageCredsModel = (*GCPSystemIdentityCredsModel)(nil)
	_ StorageCredsModel = (*GCSServiceAccountKeyCredsModel)(nil)
)

func (m *S3AccessKeyCredsModel) AsSDK() (credential.CredentialSettings, error) {
	if m.AccessKeyID.IsNull() || m.AccessKeyID.IsUnknown() {
		return nil, errors.New("access_key_id is required")
	}

	if m.SecretAccessKey.IsNull() || m.SecretAccessKey.IsUnknown() {
		return nil, errors.New("secret_access_key is required")
	}

	opts := []credential.S3CredentialAccessKeyOptions{}

	if m.ExternalID.ValueString() != "" {
		opts = append(opts, credential.WithExternalID(m.ExternalID.ValueString()))
	}

	creds := credential.NewS3CredentialAccessKey(
		m.AccessKeyID.ValueString(),
		m.SecretAccessKey.ValueString(),
		opts...,
	)

	return creds, nil
}

func (m *S3AccessKeyCredsModel) IsEmpty() bool {
	return m == nil
}

func (m *AWSSystemIdentityCredsModel) AsSDK() (credential.CredentialSettings, error) {
	if m.ExternalID.IsNull() || m.ExternalID.IsUnknown() {
		return nil, errors.New("external_id is required")
	}

	return credential.NewS3CredentialSystemIdentity(m.ExternalID.ValueString()), nil
}

func (m *AWSSystemIdentityCredsModel) IsEmpty() bool {
	return m == nil
}

func (m *CloudflareR2CredsModel) AsSDK() (credential.CredentialSettings, error) {
	if m.AccessKeyID.IsNull() || m.AccessKeyID.IsUnknown() {
		return nil, errors.New("access_key_id is required")
	}

	if m.SecretAccessKey.IsNull() || m.SecretAccessKey.IsUnknown() {
		return nil, errors.New("secret_access_key is required")
	}

	if m.AccountID.IsNull() || m.AccountID.IsUnknown() {
		return nil, errors.New("account_id is required")
	}

	if m.Token.IsNull() || m.Token.IsUnknown() {
		return nil, errors.New("token is required")
	}

	return credential.NewCloudflareR2Credential(
		m.AccessKeyID.ValueString(),
		m.SecretAccessKey.ValueString(),
		m.AccountID.ValueString(),
		m.Token.ValueString(),
	), nil
}

func (m *CloudflareR2CredsModel) IsEmpty() bool {
	return m == nil
}

func (m *AZClientCredentialsCredsModel) AsSDK() (credential.CredentialSettings, error) {
	if m.ClientID.IsNull() || m.ClientID.IsUnknown() {
		return nil, errors.New("client_id is required")
	}

	if m.ClientSecret.IsNull() || m.ClientSecret.IsUnknown() {
		return nil, errors.New("client_secret is required")
	}

	if m.TenantID.IsNull() || m.TenantID.IsUnknown() {
		return nil, errors.New("tenant_id is required")
	}

	return credential.NewAZCredentialClientCredentials(
		m.ClientID.ValueString(),
		m.ClientSecret.ValueString(),
		m.TenantID.ValueString(),
	), nil
}

func (m *AZClientCredentialsCredsModel) IsEmpty() bool {
	return m == nil
}

func (m *AZSharedAccessKeyCredsModel) AsSDK() (credential.CredentialSettings, error) {
	if m.Key.IsNull() || m.Key.IsUnknown() {
		return nil, errors.New("key is required")
	}

	return credential.NewAZCredentialSharedAccessKey(m.Key.ValueString()), nil

}

func (m *AZSharedAccessKeyCredsModel) IsEmpty() bool {
	return m == nil
}

func (m *AzureSystemIdentityCredsModel) AsSDK() (credential.CredentialSettings, error) {
	return credential.NewAZCredentialManagedIdentity(), nil
}

func (m *AzureSystemIdentityCredsModel) IsEmpty() bool {
	return m == nil
}

func (m *GCSServiceAccountKeyCredsModel) AsSDK() (credential.CredentialSettings, error) {
	if m.Key.IsNull() || m.Key.IsUnknown() {
		return nil, errors.New("key is required")
	}

	creds := credential.GCSServiceKey{}

	if err := json.Unmarshal([]byte(m.Key.ValueString()), &creds); err != nil {
		return nil, err
	}

	return credential.NewGCSCredentialServiceAccountKey(creds), nil
}

func (m *GCSServiceAccountKeyCredsModel) IsEmpty() bool {
	return m == nil
}

func (m *GCPSystemIdentityCredsModel) AsSDK() (credential.CredentialSettings, error) {
	return credential.NewGCSCredentialSystemIdentity(), nil
}

func (m *GCPSystemIdentityCredsModel) IsEmpty() bool {
	return m == nil
}

// OnlyOne returns the only non-nil argument, or an error if 0 or >1 are non-nil.
func OnlyOneStorageCredential(args ...StorageCredsModel) (StorageCredsModel, error) {
	count := 0
	var selected StorageCredsModel

	for _, arg := range args {
		if !arg.IsEmpty() {
			count++
			selected = arg
		}
	}

	if count != 1 {
		return nil, fmt.Errorf("you can set one and only one storage credential, got %d", count)
	}

	return selected, nil
}
