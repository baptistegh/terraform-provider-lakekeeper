package sdk

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	managementv1 "github.com/lakekeeper/go-lakekeeper/pkg/apis/management/v1"
	"github.com/lakekeeper/go-lakekeeper/pkg/storage/credential"
)

type (
	StorageCredsModel interface {
		AsSDK() (managementv1.StorageCredential, error)
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

	AzureSystemIdentityCredsModel struct {
		Enabled types.Bool `tfsdk:"enabled"`
	}
	GCPSystemIdentityCredsModel struct {
		Enabled types.Bool `tfsdk:"enabled"`
	}
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

func (m *S3AccessKeyCredsModel) AsSDK() (managementv1.StorageCredential, error) {
	if m.AccessKeyID.IsNull() || m.AccessKeyID.IsUnknown() {
		return managementv1.StorageCredential{}, errors.New("access_key_id is required")
	}

	if m.SecretAccessKey.IsNull() || m.SecretAccessKey.IsUnknown() {
		return managementv1.StorageCredential{}, errors.New("secret_access_key is required")
	}

	if m.ExternalID.ValueString() != "" {
		return credential.NewS3AccessKeyWithExternalID(
			m.AccessKeyID.ValueString(),
			m.SecretAccessKey.ValueString(),
			m.ExternalID.ValueString(),
		), nil
	}

	return credential.NewS3AccessKey(
		m.AccessKeyID.ValueString(),
		m.SecretAccessKey.ValueString(),
	), nil
}

func (m *S3AccessKeyCredsModel) IsEmpty() bool {
	return m == nil
}

func (m *AWSSystemIdentityCredsModel) AsSDK() (managementv1.StorageCredential, error) {
	if m.ExternalID.IsNull() || m.ExternalID.IsUnknown() {
		return managementv1.StorageCredential{}, errors.New("external_id is required")
	}

	return credential.NewS3AwsSystemIdentityWithExternalID(m.ExternalID.ValueString()), nil
}

func (m *AWSSystemIdentityCredsModel) IsEmpty() bool {
	return m == nil
}

func (m *CloudflareR2CredsModel) AsSDK() (managementv1.StorageCredential, error) {
	if m.AccessKeyID.IsNull() || m.AccessKeyID.IsUnknown() {
		return managementv1.StorageCredential{}, errors.New("access_key_id is required")
	}

	if m.SecretAccessKey.IsNull() || m.SecretAccessKey.IsUnknown() {
		return managementv1.StorageCredential{}, errors.New("secret_access_key is required")
	}

	if m.AccountID.IsNull() || m.AccountID.IsUnknown() {
		return managementv1.StorageCredential{}, errors.New("account_id is required")
	}

	if m.Token.IsNull() || m.Token.IsUnknown() {
		return managementv1.StorageCredential{}, errors.New("token is required")
	}

	return credential.NewS3CloudflareR2(
		m.AccessKeyID.ValueString(),
		m.SecretAccessKey.ValueString(),
		m.AccountID.ValueString(),
		m.Token.ValueString(),
	), nil
}

func (m *CloudflareR2CredsModel) IsEmpty() bool {
	return m == nil
}

func (m *AZClientCredentialsCredsModel) AsSDK() (managementv1.StorageCredential, error) {
	if m.ClientID.IsNull() || m.ClientID.IsUnknown() {
		return managementv1.StorageCredential{}, errors.New("client_id is required")
	}

	if m.ClientSecret.IsNull() || m.ClientSecret.IsUnknown() {
		return managementv1.StorageCredential{}, errors.New("client_secret is required")
	}

	if m.TenantID.IsNull() || m.TenantID.IsUnknown() {
		return managementv1.StorageCredential{}, errors.New("tenant_id is required")
	}

	return credential.NewAZClientCredentials(
		m.TenantID.ValueString(),
		m.ClientID.ValueString(),
		m.ClientSecret.ValueString(),
	), nil
}

func (m *AZClientCredentialsCredsModel) IsEmpty() bool {
	return m == nil
}

func (m *AZSharedAccessKeyCredsModel) AsSDK() (managementv1.StorageCredential, error) {
	if m.Key.IsNull() || m.Key.IsUnknown() {
		return managementv1.StorageCredential{}, errors.New("key is required")
	}

	return credential.NewAZSharedAccessKey(m.Key.ValueString()), nil

}

func (m *AZSharedAccessKeyCredsModel) IsEmpty() bool {
	return m == nil
}

func (m *AzureSystemIdentityCredsModel) AsSDK() (managementv1.StorageCredential, error) {
	return credential.NewAZManagedIdentity(), nil
}

func (m *AzureSystemIdentityCredsModel) IsEmpty() bool {
	return m == nil
}

func (m *GCSServiceAccountKeyCredsModel) AsSDK() (managementv1.StorageCredential, error) {
	if m.Key.IsNull() || m.Key.IsUnknown() {
		return managementv1.StorageCredential{}, errors.New("key is required")
	}

	// Keys issued before Google added universe_domain lack it, while the
	// Lakekeeper API requires it. Default it to the public Google Cloud universe.
	fields := map[string]any{}
	if err := json.Unmarshal([]byte(m.Key.ValueString()), &fields); err != nil {
		return managementv1.StorageCredential{}, err
	}
	if _, ok := fields["universe_domain"]; !ok {
		fields["universe_domain"] = defaultGCSUniverseDomain
	}

	raw, err := json.Marshal(fields)
	if err != nil {
		return managementv1.StorageCredential{}, err
	}

	creds := managementv1.GcsServiceKey{}
	if err := json.Unmarshal(raw, &creds); err != nil {
		return managementv1.StorageCredential{}, err
	}

	return credential.NewGCSServiceAccountKey(creds), nil
}

func (m *GCSServiceAccountKeyCredsModel) IsEmpty() bool {
	return m == nil
}

func (m *GCPSystemIdentityCredsModel) AsSDK() (managementv1.StorageCredential, error) {
	return credential.NewGCSSystemIdentity(), nil
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
