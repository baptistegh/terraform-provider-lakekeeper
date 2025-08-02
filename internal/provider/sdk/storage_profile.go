package sdk

import (
	"errors"
	"fmt"

	"github.com/baptistegh/go-lakekeeper/pkg/apis/management/v1/storage/credential"
	"github.com/baptistegh/go-lakekeeper/pkg/apis/management/v1/storage/profile"
	"github.com/baptistegh/go-lakekeeper/pkg/core"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type (
	StorageProfileModel interface {
		StorageFamily() profile.StorageFamily
		AsSDK() (profile.StorageSettings, error)
		CredentialAsSDK() (credential.CredentialSettings, error)
		AddCreds(StorageCredsModel) error
		IsEmpty() bool
		GetCredentials() (StorageCredsModel, error)
	}

	S3StorageProfileModel struct {
		Bucket                    types.String                 `tfsdk:"bucket"`
		Region                    types.String                 `tfsdk:"region"`
		KeyPrefix                 types.String                 `tfsdk:"key_prefix"`
		AllowAlternativeProtocols types.Bool                   `tfsdk:"allow_alternative_protocols"`
		AssumeRoleARN             types.String                 `tfsdk:"assume_role_arn"`
		AWSKMSKeyARN              types.String                 `tfsdk:"aws_kms_key_arn"`
		Endpoint                  types.String                 `tfsdk:"endpoint"`
		Flavor                    types.String                 `tfsdk:"flavor"`
		PathStyleAccess           types.Bool                   `tfsdk:"path_style_access"`
		PushS3DeleteDisabled      types.Bool                   `tfsdk:"push_s3_delete_disabled"`
		RemoteSigningURLStyle     types.String                 `tfsdk:"remote_signing_url_style"`
		STSEnabled                types.Bool                   `tfsdk:"sts_enabled"`
		STSRoleARN                types.String                 `tfsdk:"sts_role_arn"`
		STSTokenValiditySeconds   types.Int64                  `tfsdk:"sts_token_validity_seconds"`
		AccessKey                 *S3AccessKeyCredsModel       `tfsdk:"access_key"`
		CloudflareR2              *CloudflareR2CredsModel      `tfsdk:"cloudflare_r2"`
		AWSSystemIdentity         *AWSSystemIdentityCredsModel `tfsdk:"aws_system_identity"`
	}

	ADLSStorageProfileModel struct {
		AccountName               types.String                   `tfsdk:"account_name"`
		Filesystem                types.String                   `tfsdk:"filesystem"`
		Host                      types.String                   `tfsdk:"host"`
		AuthorityHost             types.String                   `tfsdk:"authority_host"`
		KeyPrefix                 types.String                   `tfsdk:"key_prefix"`
		SASTokenValiditySeconds   types.Int64                    `tfsdk:"sas_token_validity_seconds"`
		AllowAlternativeProtocols types.Bool                     `tfsdk:"allow_alternative_protocols"`
		SharedAccessKey           *AZSharedAccessKeyCredsModel   `tfsdk:"shared_access_key"`
		ClientCredentials         *AZClientCredentialsCredsModel `tfsdk:"client_credentials"`
		SystemIdentity            *AzureSystemIdentityCredsModel `tfsdk:"azure_system_identity"`
	}

	GCSStorageProfileModel struct {
		Bucket            types.String                    `tfsdk:"bucket"`
		KeyPrefix         types.String                    `tfsdk:"key_prefix"`
		ServiceAccountKey *GCSServiceAccountKeyCredsModel `tfsdk:"service_account_key"`
		SystemIdentity    *GCPSystemIdentityCredsModel    `tfsdk:"gcp_system_identity"`
	}
)

var (
	_ StorageProfileModel = (*S3StorageProfileModel)(nil)
	_ StorageProfileModel = (*ADLSStorageProfileModel)(nil)
	_ StorageProfileModel = (*GCSStorageProfileModel)(nil)
)

func (m *S3StorageProfileModel) AsSDK() (profile.StorageSettings, error) {
	if m.Bucket.IsNull() || m.Bucket.IsUnknown() {
		return nil, errors.New("bucket is required")
	}

	if m.Region.IsNull() || m.Region.IsUnknown() {
		return nil, errors.New("region is required")
	}

	opt := []profile.S3StorageSettingsOptions{}

	if m.KeyPrefix.ValueString() != "" {
		opt = append(opt, profile.WithS3KeyPrefix(m.KeyPrefix.ValueString()))
	}

	if m.AllowAlternativeProtocols.ValueBool() {
		opt = append(opt, profile.WithS3AlternativeProtocols())
	}

	if m.AssumeRoleARN.ValueString() != "" {
		opt = append(opt, profile.WithAssumeRoleARN(m.AssumeRoleARN.ValueString()))
	}

	if m.AWSKMSKeyARN.ValueString() != "" {
		opt = append(opt, profile.WithAWSKMSKeyARN(m.AWSKMSKeyARN.ValueString()))
	}

	if m.Endpoint.ValueString() != "" {
		opt = append(opt, profile.WithEndpoint(m.Endpoint.ValueString()))
	}

	if m.Flavor.ValueString() != "" {
		opt = append(opt, profile.WithFlavor(profile.S3Flavor(m.Flavor.ValueString())))
	}

	if m.PathStyleAccess.ValueBool() {
		opt = append(opt, profile.WithPathStyleAccess())
	}

	if m.PushS3DeleteDisabled.ValueBool() {
		opt = append(opt, profile.WithPushS3DeleteDisabled(m.PushS3DeleteDisabled.ValueBool()))
	}

	if m.RemoteSigningURLStyle.ValueString() != "" {
		opt = append(opt, profile.WithRemoteSigningURLStyle(profile.RemoteSigningURLStyle(m.RemoteSigningURLStyle.ValueString())))
	}

	if m.STSEnabled.ValueBool() {
		opt = append(opt, profile.WithSTSEnabled())
	}

	if m.STSRoleARN.ValueString() != "" {
		opt = append(opt, profile.WithSTSRoleARN(m.STSRoleARN.ValueString()))
	}

	if m.STSTokenValiditySeconds.ValueInt64() != 0 {
		opt = append(opt, profile.WithSTSTokenValiditySeconds(m.STSTokenValiditySeconds.ValueInt64()))
	}

	sp := profile.NewS3StorageSettings(
		m.Bucket.ValueString(),
		m.Region.ValueString(),
		opt...,
	)

	return sp, nil
}

func (m *S3StorageProfileModel) CredentialAsSDK() (credential.CredentialSettings, error) {
	storage, err := OnlyOneStorageCredential(m.AccessKey, m.CloudflareR2, m.AWSSystemIdentity)
	if err != nil {
		return nil, err
	}

	return storage.AsSDK()
}

func (m *S3StorageProfileModel) StorageFamily() profile.StorageFamily {
	return profile.StorageFamilyS3
}

func (m *S3StorageProfileModel) IsEmpty() bool {
	return m == nil
}

func (m *S3StorageProfileModel) AddCreds(c StorageCredsModel) error {
	switch v := c.(type) {
	case *S3AccessKeyCredsModel:
		m.AccessKey = v
		m.AWSSystemIdentity = nil
		m.CloudflareR2 = nil
	case *AWSSystemIdentityCredsModel:
		m.AccessKey = nil
		m.AWSSystemIdentity = v
		m.CloudflareR2 = nil
	case *CloudflareR2CredsModel:
		m.AccessKey = nil
		m.AWSSystemIdentity = nil
		m.CloudflareR2 = v
	default:
		return fmt.Errorf("incorrect storage credential type %T", v)
	}
	return nil
}

func (m *S3StorageProfileModel) GetCredentials() (StorageCredsModel, error) {
	return OnlyOneStorageCredential(m.AccessKey, m.AWSSystemIdentity, m.CloudflareR2)
}

func (m *ADLSStorageProfileModel) AsSDK() (profile.StorageSettings, error) {
	if m.AccountName.IsNull() || m.AccountName.IsUnknown() {
		return nil, errors.New("account_name is required")
	}

	if m.Filesystem.IsNull() || m.Filesystem.IsUnknown() {
		return nil, errors.New("filesystem is required")
	}

	opt := []profile.ADLSStorageSettingsOptions{}

	if m.AuthorityHost.ValueString() != "" {
		opt = append(opt, profile.WithAuthorityHost(m.AuthorityHost.ValueString()))
	}

	if m.Host.ValueString() != "" {
		opt = append(opt, profile.WithHost(m.Host.ValueString()))
	}

	if m.KeyPrefix.ValueString() != "" {
		opt = append(opt, profile.WithADLSKeyPrefix(m.KeyPrefix.ValueString()))
	}

	if m.SASTokenValiditySeconds.ValueInt64() != 0 {
		opt = append(opt, profile.WithSASTokenValiditySeconds(m.SASTokenValiditySeconds.ValueInt64()))
	}

	if m.AllowAlternativeProtocols.ValueBool() {
		opt = append(opt, profile.WithADLSAlternativeProtocols())
	}

	sp := profile.NewADLSStorageSettings(
		m.AccountName.ValueString(),
		m.Filesystem.ValueString(),
		opt...,
	)

	return sp, nil
}

func (m *ADLSStorageProfileModel) CredentialAsSDK() (credential.CredentialSettings, error) {
	storage, err := OnlyOneStorageCredential(m.SharedAccessKey, m.ClientCredentials, m.SystemIdentity)
	if err != nil {
		return nil, err
	}

	return storage.AsSDK()
}

func (m *GCSStorageProfileModel) StorageFamily() profile.StorageFamily {
	return profile.StorageFamilyGCS
}

func (m *GCSStorageProfileModel) IsEmpty() bool {
	return m == nil
}

func (m *GCSStorageProfileModel) AddCreds(c StorageCredsModel) error {
	switch v := c.(type) {
	case *GCSServiceAccountKeyCredsModel:
		m.ServiceAccountKey = v
		m.SystemIdentity = nil
	case *GCPSystemIdentityCredsModel:
		m.ServiceAccountKey = nil
		m.SystemIdentity = v
	default:
		return fmt.Errorf("incorrect storage credential type %T", v)
	}
	return nil
}

func (m *GCSStorageProfileModel) GetCredentials() (StorageCredsModel, error) {
	return OnlyOneStorageCredential(m.ServiceAccountKey, m.SystemIdentity)
}

func (m *GCSStorageProfileModel) AsSDK() (profile.StorageSettings, error) {
	if m.Bucket.IsNull() || m.Bucket.IsUnknown() {
		return nil, errors.New("bucket is required")
	}

	opt := []profile.GCSStorageSettingsOptions{}

	if m.KeyPrefix.ValueString() != "" {
		opt = append(opt, profile.WithGCSKeyPrefix(m.KeyPrefix.ValueString()))
	}

	sp := profile.NewGCSStorageSettings(
		m.Bucket.ValueString(),
		opt...,
	)

	return sp, nil
}

func (m *GCSStorageProfileModel) CredentialAsSDK() (credential.CredentialSettings, error) {
	creds, err := OnlyOneStorageCredential(m.ServiceAccountKey, m.SystemIdentity)
	if err != nil {
		return nil, err
	}

	return creds.AsSDK()
}

func (m *ADLSStorageProfileModel) StorageFamily() profile.StorageFamily {
	return profile.StorageFamilyADLS
}

func (m *ADLSStorageProfileModel) IsEmpty() bool {
	return m == nil
}

func (m *ADLSStorageProfileModel) AddCreds(c StorageCredsModel) error {
	switch v := c.(type) {
	case *AZClientCredentialsCredsModel:
		m.ClientCredentials = v
		m.SharedAccessKey = nil
		m.SystemIdentity = nil
	case *AZSharedAccessKeyCredsModel:
		m.ClientCredentials = nil
		m.SharedAccessKey = v
		m.SystemIdentity = nil
	case *AzureSystemIdentityCredsModel:
		m.ClientCredentials = nil
		m.SharedAccessKey = nil
		m.SystemIdentity = v
	default:
		return fmt.Errorf("incorrect storage credential type %T", v)
	}
	return nil
}

func (m *ADLSStorageProfileModel) GetCredentials() (StorageCredsModel, error) {
	return OnlyOneStorageCredential(m.ClientCredentials, m.SharedAccessKey, m.SystemIdentity)
}

func StorageProfileModelFromSDK(sp profile.StorageProfile) (StorageProfileModel, error) {
	if sp.StorageSettings == nil {
		return nil, errors.New("storage profile is empty")
	}

	switch sp.StorageSettings.GetStorageFamily() {
	case profile.StorageFamilyS3:
		cfg, ok := sp.AsS3()
		if !ok {
			return nil, errors.New("invalid storage profile")
		}
		sp := &S3StorageProfileModel{
			Bucket:                    types.StringValue(cfg.Bucket),
			Region:                    types.StringValue(cfg.Region),
			STSEnabled:                types.BoolValue(cfg.STSEnabled),
			STSRoleARN:                types.StringPointerValue(cfg.STSRoleARN),
			STSTokenValiditySeconds:   types.Int64PointerValue(cfg.STSTokenValiditySeconds),
			PushS3DeleteDisabled:      types.BoolPointerValue(cfg.PushS3DeleteDisabled),
			KeyPrefix:                 types.StringPointerValue(cfg.KeyPrefix),
			AllowAlternativeProtocols: types.BoolPointerValue(cfg.AllowAlternativeProtocols),
			AssumeRoleARN:             types.StringPointerValue(cfg.AssumeRoleARN),
			AWSKMSKeyARN:              types.StringPointerValue(cfg.AWSKMSKeyARN),
			Endpoint:                  types.StringPointerValue(cfg.Endpoint),
			PathStyleAccess:           types.BoolPointerValue(cfg.PathStyleAccess),
		}

		var f *string
		if cfg.Flavor != nil {
			f = core.Ptr(string(*cfg.Flavor))
		}

		var style *string
		if cfg.RemoteSigningURLStyle != nil {
			style = core.Ptr(string(*cfg.RemoteSigningURLStyle))
		}

		sp.Flavor = types.StringPointerValue(f)
		sp.RemoteSigningURLStyle = types.StringPointerValue(style)

		return sp, nil
	case profile.StorageFamilyADLS:
		cfg, ok := sp.AsADLS()
		if !ok {
			return nil, errors.New("invalid storage profile")
		}

		return &ADLSStorageProfileModel{
			AccountName:               types.StringValue(cfg.AccountName),
			Filesystem:                types.StringValue(cfg.Filesystem),
			Host:                      types.StringPointerValue(cfg.Host),
			AuthorityHost:             types.StringPointerValue(cfg.AuthorityHost),
			KeyPrefix:                 types.StringPointerValue(cfg.KeyPrefix),
			SASTokenValiditySeconds:   types.Int64PointerValue(cfg.SASTokenValiditySeconds),
			AllowAlternativeProtocols: types.BoolPointerValue(cfg.AllowAlternativeProtocols),
		}, nil
	case profile.StorageFamilyGCS:
		cfg, ok := sp.AsGCS()
		if !ok {
			return nil, errors.New("invalid storage profile")
		}

		return &GCSStorageProfileModel{
			Bucket:    types.StringValue(cfg.Bucket),
			KeyPrefix: types.StringPointerValue(cfg.KeyPrefix),
		}, nil
	}

	return nil, errors.New("unsupported storage profile")
}

// OnlyOne returns the only non-nil argument, or an error if 0 or >1 are non-nil.
func OnlyOneStorageProfile(args ...StorageProfileModel) (StorageProfileModel, error) {
	count := 0
	var selected StorageProfileModel

	for _, arg := range args {
		if !arg.IsEmpty() {
			count++
			selected = arg
		}
	}

	if count != 1 {
		return nil, fmt.Errorf("you can set one and only one storage profile, got %d", count)
	}

	return selected, nil
}
