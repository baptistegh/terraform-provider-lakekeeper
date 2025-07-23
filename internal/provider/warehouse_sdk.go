package provider

import (
	"errors"
	"fmt"

	tftypes "github.com/baptistegh/terraform-provider-lakekeeper/internal/provider/types"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	managementv1 "github.com/baptistegh/go-lakekeeper/pkg/apis/management/v1"
	"github.com/baptistegh/go-lakekeeper/pkg/apis/management/v1/storage/credential"
	"github.com/baptistegh/go-lakekeeper/pkg/apis/management/v1/storage/profile"
)

func (m *lakekeeperWarehouseResourceModel) ToWarehouseCreateRequest() (*managementv1.CreateWarehouseOptions, error) {
	req := managementv1.CreateWarehouseOptions{
		Name: m.Name.ValueString(),
	}

	if m.DeleteProfile != nil {
		settings, err := m.DeleteProfileSettings()
		if err != nil {
			return nil, err
		}
		if settings == nil {
			return nil, errors.New("delete profile is empty")
		}

		req.DeleteProfile = settings.AsProfile()
	}

	if m.StorageProfile != nil {
		settings, err := m.StorageSettings()
		if err != nil {
			return nil, err
		}

		if settings == nil {
			return nil, errors.New("storage profile is empty")
		}

		req.StorageProfile = settings.AsProfile()
	}

	if m.StorageCredential != nil {
		settings, err := m.StorageCredentialSettings()
		if err != nil {
			return nil, err
		}
		req.StorageCredential = settings.AsCredential()
	}

	return &req, nil
}

// TODO: refactor RefreshFromSettings on datasource and resource
// because these functions are almost identical

func (m *lakekeeperWarehouseResourceModel) RefreshFromSettings(w *managementv1.Warehouse) diag.Diagnostics {
	m.ID = types.StringValue(w.ProjectID + "/" + w.ID)
	m.WarehouseID = types.StringValue(w.ID)
	m.ProjectID = types.StringValue(w.ProjectID)
	m.Protected = types.BoolValue(w.Protected)
	m.Active = types.BoolValue(w.IsActive())
	m.Name = types.StringValue(w.Name)

	diags := diag.Diagnostics{}
	const errorMessage = "Error refreshing warehouse state"

	if w.StorageProfile.StorageSettings == nil {
		m.StorageProfile = nil
		diags.AddError(errorMessage, "Storage profile must be defined")
	} else {
		m.StorageProfile = &tftypes.StorageProfileModel{}
		settings := w.StorageProfile.StorageSettings

		switch sp := settings.(type) {
		case *profile.ADLSStorageSettings:
			m.StorageProfile.Type = types.StringValue(string(sp.GetStorageFamily()))
			m.StorageProfile.AccountName = types.StringValue(sp.AccountName)
			m.StorageProfile.AllowAlternativeProtocols = types.BoolPointerValue(sp.AllowAlternativeProtocols)
			m.StorageProfile.AuthorityHost = types.StringPointerValue(sp.AuthorityHost)
			m.StorageProfile.Filesystem = types.StringValue(sp.Filesystem)
			m.StorageProfile.Host = types.StringPointerValue(sp.Host)
			m.StorageProfile.KeyPrefix = types.StringPointerValue(sp.KeyPrefix)
			m.StorageProfile.SASTokenValiditySeconds = types.Int64PointerValue(sp.SASTokenValiditySeconds)
		case *profile.GCSStorageSettings:
			m.StorageProfile.Type = types.StringValue(string(sp.GetStorageFamily()))
			m.StorageProfile.Bucket = types.StringValue(sp.Bucket)
			m.StorageProfile.KeyPrefix = types.StringPointerValue(sp.KeyPrefix)
		case *profile.S3StorageSettings:
			m.StorageProfile.Type = types.StringValue(string(sp.GetStorageFamily()))
			m.StorageProfile.AllowAlternativeProtocols = types.BoolPointerValue(sp.AllowAlternativeProtocols)
			m.StorageProfile.AssumeRoleARN = types.StringPointerValue(sp.AssumeRoleARN)
			m.StorageProfile.AWSKMSKeyARN = types.StringPointerValue(sp.AWSKMSKeyARN)
			m.StorageProfile.Bucket = types.StringValue(sp.Bucket)
			m.StorageProfile.Endpoint = types.StringPointerValue(sp.Endpoint)
			m.StorageProfile.Flavor = types.StringValue(string(*sp.Flavor))
			m.StorageProfile.KeyPrefix = types.StringPointerValue(sp.KeyPrefix)
			m.StorageProfile.PathStyleAccess = types.BoolPointerValue(sp.PathStyleAccess)
			m.StorageProfile.PushS3DeleteDisabled = types.BoolPointerValue(sp.PushS3DeleteDisabled)
			m.StorageProfile.Region = types.StringValue(sp.Region)
			m.StorageProfile.RemoteSigningURLStyle = types.StringValue(string(*sp.RemoteSigningURLStyle))
			m.StorageProfile.STSEnabled = types.BoolValue(sp.STSEnabled)
			m.StorageProfile.STSRoleARN = types.StringPointerValue(sp.STSRoleARN)
			m.StorageProfile.STSTokenValiditySeconds = types.Int64PointerValue(sp.STSTokenValiditySeconds)
		default:
			diags.AddError(errorMessage, fmt.Sprintf("Incorrect storage profile type: %T, valid: [s3,adls,gcs]", sp))
		}
	}

	if w.DeleteProfile == nil || w.DeleteProfile.DeleteProfileSettings == nil {
		m.DeleteProfile = nil
	} else {
		switch deleteProfile := w.DeleteProfile.DeleteProfileSettings.(type) {
		case *profile.TabularDeleteProfileSoft:
			m.DeleteProfile = &tftypes.DeleteProfileModel{
				Type:              types.StringValue("soft"),
				ExpirationSeconds: types.Int32Value(deleteProfile.ExpirationSeconds),
			}
		case *profile.TabularDeleteProfileHard:
			m.DeleteProfile = &tftypes.DeleteProfileModel{
				Type: types.StringValue("hard"),
			}
		default:
			diags.AddError(errorMessage, fmt.Sprintf("Incorrect delete profile type: %T, valid: [soft,hard]", deleteProfile))
		}
	}

	// Lakekeeper API does not give storage-credential on GET, it can't be refreshed
	if m.StorageCredential == nil {
		diags.AddError(errorMessage, "Storage credential must be defined")
	}

	return diags
}

func (m *lakekeeperWarehouseDataSourceModel) RefreshFromSettings(w *managementv1.Warehouse) diag.Diagnostics {
	m.ID = types.StringValue(w.ProjectID + "/" + w.ID)
	m.WarehouseID = types.StringValue(w.ID)
	m.ProjectID = types.StringValue(w.ProjectID)
	m.Protected = types.BoolValue(w.Protected)
	m.Active = types.BoolValue(w.IsActive())
	m.Name = types.StringValue(w.Name)

	diags := diag.Diagnostics{}
	const errorMessage = "Error refreshing warehouse state"

	if w.StorageProfile.StorageSettings == nil {
		m.StorageProfile = nil
		diags.AddError(errorMessage, "Storage profile must be defined")
	} else {
		m.StorageProfile = &tftypes.StorageProfileModel{}
		storageProfile := w.StorageProfile.StorageSettings

		switch sp := storageProfile.(type) {
		case *profile.ADLSStorageSettings:
			m.StorageProfile.Type = types.StringValue(string(sp.GetStorageFamily()))
			m.StorageProfile.AccountName = types.StringValue(sp.AccountName)
			m.StorageProfile.AllowAlternativeProtocols = types.BoolPointerValue(sp.AllowAlternativeProtocols)
			m.StorageProfile.AuthorityHost = types.StringPointerValue(sp.AuthorityHost)
			m.StorageProfile.Filesystem = types.StringValue(sp.Filesystem)
			m.StorageProfile.Host = types.StringPointerValue(sp.Host)
			m.StorageProfile.KeyPrefix = types.StringPointerValue(sp.KeyPrefix)
			m.StorageProfile.SASTokenValiditySeconds = types.Int64PointerValue(sp.SASTokenValiditySeconds)
		case *profile.GCSStorageSettings:
			m.StorageProfile.Type = types.StringValue(string(sp.GetStorageFamily()))
			m.StorageProfile.Bucket = types.StringValue(sp.Bucket)
			m.StorageProfile.KeyPrefix = types.StringPointerValue(sp.KeyPrefix)
		case *profile.S3StorageSettings:
			m.StorageProfile.Type = types.StringValue(string(sp.GetStorageFamily()))
			m.StorageProfile.AllowAlternativeProtocols = types.BoolPointerValue(sp.AllowAlternativeProtocols)
			m.StorageProfile.AssumeRoleARN = types.StringPointerValue(sp.AssumeRoleARN)
			m.StorageProfile.AWSKMSKeyARN = types.StringPointerValue(sp.AWSKMSKeyARN)
			m.StorageProfile.Bucket = types.StringValue(sp.Bucket)
			m.StorageProfile.Endpoint = types.StringPointerValue(sp.Endpoint)
			m.StorageProfile.Flavor = types.StringValue(string(*sp.Flavor))
			m.StorageProfile.KeyPrefix = types.StringPointerValue(sp.KeyPrefix)
			m.StorageProfile.PathStyleAccess = types.BoolPointerValue(sp.PathStyleAccess)
			m.StorageProfile.PushS3DeleteDisabled = types.BoolPointerValue(sp.PushS3DeleteDisabled)
			m.StorageProfile.Region = types.StringValue(sp.Region)
			m.StorageProfile.RemoteSigningURLStyle = types.StringValue(string(*sp.RemoteSigningURLStyle))
			m.StorageProfile.STSEnabled = types.BoolValue(sp.STSEnabled)
			m.StorageProfile.STSRoleARN = types.StringPointerValue(sp.STSRoleARN)
			m.StorageProfile.STSTokenValiditySeconds = types.Int64PointerValue(sp.STSTokenValiditySeconds)
		default:
			diags.AddError(errorMessage, fmt.Sprintf("Incorrect storage profile type: %T, valid: [s3,adls,gcs]", sp))
		}
	}

	if w.DeleteProfile == nil || w.DeleteProfile.DeleteProfileSettings == nil {
		m.DeleteProfile = nil
	} else {
		m.DeleteProfile = &tftypes.DeleteProfileModel{}
		switch deleteProfile := w.DeleteProfile.DeleteProfileSettings.(type) {
		case *profile.TabularDeleteProfileSoft:
			m.DeleteProfile = &tftypes.DeleteProfileModel{
				Type:              types.StringValue("soft"),
				ExpirationSeconds: types.Int32Value(deleteProfile.ExpirationSeconds),
			}
		case *profile.TabularDeleteProfileHard:
			m.DeleteProfile = &tftypes.DeleteProfileModel{
				Type: types.StringValue("hard"),
			}
		default:
			diags.AddError(errorMessage, fmt.Sprintf("Incorrect delete profile type: %T, valid: [soft,hard]", deleteProfile))
		}
	}

	return diags
}

func (m *lakekeeperWarehouseResourceModel) DeleteProfileSettings() (profile.DeleteProfileSettings, error) {
	if m.DeleteProfile == nil {
		return nil, nil
	}

	if m.DeleteProfile.Type.IsNull() || m.DeleteProfile.Type.IsUnknown() {
		return profile.NewTabularDeleteProfileHard(), nil
	}

	switch m.DeleteProfile.Type.ValueString() {
	case "soft":
		return profile.NewTabularDeleteProfileSoft(m.DeleteProfile.ExpirationSeconds.ValueInt32()), nil
	case "hard":
		return profile.NewTabularDeleteProfileHard(), nil
	default:
		return nil, fmt.Errorf("incorrect delete profile definition, type must be [soft,hard]")
	}
}

func (m *lakekeeperWarehouseResourceModel) StorageSettings() (profile.StorageSettings, error) {
	if m.StorageProfile == nil {
		return nil, nil
	}
	switch m.StorageProfile.Type.ValueString() {
	case "s3":
		opts := []profile.S3StorageSettingsOptions{}

		if m.StorageProfile.STSEnabled.ValueBool() {
			opts = append(opts, profile.WithSTSEnabled())
		}

		if m.StorageProfile.AllowAlternativeProtocols.ValueBool() {
			opts = append(opts, profile.WithS3AlternativeProtocols())
		}

		if !m.StorageProfile.AssumeRoleARN.IsNull() && !m.StorageProfile.AssumeRoleARN.IsUnknown() {
			opts = append(opts, profile.WithAssumeRoleARN(m.StorageProfile.AssumeRoleARN.ValueString()))
		}

		if !m.StorageProfile.AWSKMSKeyARN.IsNull() && !m.StorageProfile.AWSKMSKeyARN.IsUnknown() {
			opts = append(opts, profile.WithAWSKMSKeyARN(m.StorageProfile.AWSKMSKeyARN.ValueString()))
		}

		if !m.StorageProfile.Endpoint.IsNull() && !m.StorageProfile.Endpoint.IsUnknown() {
			opts = append(opts, profile.WithEndpoint(m.StorageProfile.Endpoint.ValueString()))
		}

		if !m.StorageProfile.Flavor.IsNull() && !m.StorageProfile.Flavor.IsUnknown() {
			flavor := profile.S3Flavor(m.StorageProfile.Flavor.ValueString())
			opts = append(opts, profile.WithFlavor(flavor))
		}

		if !m.StorageProfile.KeyPrefix.IsNull() && !m.StorageProfile.KeyPrefix.IsUnknown() {
			opts = append(opts, profile.WithS3KeyPrefix(m.StorageProfile.KeyPrefix.ValueString()))
		}

		if m.StorageProfile.PathStyleAccess.ValueBool() {
			opts = append(opts, profile.WithPathStyleAccess())
		}

		if !m.StorageProfile.PushS3DeleteDisabled.IsNull() && !m.StorageProfile.PushS3DeleteDisabled.IsUnknown() {
			opts = append(opts, profile.WithPushS3DeleteDisabled(m.StorageProfile.PushS3DeleteDisabled.ValueBool()))
		}

		if !m.StorageProfile.RemoteSigningURLStyle.IsNull() && !m.StorageProfile.RemoteSigningURLStyle.IsUnknown() {
			style := profile.RemoteSigningURLStyle(m.StorageProfile.RemoteSigningURLStyle.ValueString())
			opts = append(opts, profile.WithRemoteSigningURLStyle(style))
		}

		if !m.StorageProfile.STSRoleARN.IsNull() && !m.StorageProfile.STSRoleARN.IsUnknown() {
			opts = append(opts, profile.WithSTSRoleARN(m.StorageProfile.STSRoleARN.ValueString()))
		}

		if !m.StorageProfile.STSTokenValiditySeconds.IsNull() && !m.StorageProfile.STSTokenValiditySeconds.IsUnknown() {
			opts = append(opts, profile.WithSTSTokenValiditySeconds(m.StorageProfile.STSTokenValiditySeconds.ValueInt64()))
		}

		profile := profile.NewS3StorageSettings(
			m.StorageProfile.Bucket.ValueString(),
			m.StorageProfile.Region.ValueString(),
			opts...,
		)

		return profile, nil
	case "adls":
		opts := []profile.ADLSStorageSettingsOptions{}

		if m.StorageProfile.AllowAlternativeProtocols.ValueBool() {
			opts = append(opts, profile.WithADLSAlternativeProtocols())
		}

		if !m.StorageProfile.AuthorityHost.IsNull() && !m.StorageProfile.AuthorityHost.IsUnknown() {
			opts = append(opts, profile.WithAuthorityHost(m.StorageProfile.AuthorityHost.ValueString()))
		}

		if !m.StorageProfile.Host.IsNull() && !m.StorageProfile.Host.IsUnknown() {
			opts = append(opts, profile.WithHost(m.StorageProfile.Host.ValueString()))
		}

		if !m.StorageProfile.KeyPrefix.IsNull() && !m.StorageProfile.KeyPrefix.IsUnknown() {
			opts = append(opts, profile.WithADLSKeyPrefix(m.StorageProfile.KeyPrefix.ValueString()))
		}

		if !m.StorageProfile.SASTokenValiditySeconds.IsNull() && !m.StorageProfile.SASTokenValiditySeconds.IsUnknown() {
			opts = append(opts, profile.WithSASTokenValiditySeconds(m.StorageProfile.SASTokenValiditySeconds.ValueInt64()))
		}

		profile := profile.NewADLSStorageSettings(
			m.StorageProfile.AccountName.ValueString(),
			m.StorageProfile.Filesystem.ValueString(),
			opts...,
		)

		return profile, nil
	case "gcs":
		opts := []profile.GCSStorageSettingsOptions{}

		if !m.StorageProfile.KeyPrefix.IsNull() && !m.StorageProfile.KeyPrefix.IsUnknown() {
			opts = append(opts, profile.WithGCSKeyPrefix(m.StorageProfile.KeyPrefix.ValueString()))
		}

		profile := profile.NewGCSStorageSettings(
			m.StorageProfile.Bucket.ValueString(),
			opts...,
		)

		return profile, nil

	default:
		return nil, fmt.Errorf("invalid storage profile definitions, type must be [s3,gcs,adls]")
	}

}

func (m *lakekeeperWarehouseResourceModel) StorageCredentialSettings() (credential.CredentialSettings, error) {
	if m.StorageCredential == nil {
		return nil, fmt.Errorf("invalid storage credential definitions, must be defined")
	}
	storageType := m.StorageCredential.Type.ValueString()
	switch storageType {
	case "s3_access_key":
		opts := []credential.S3CredentialAccessKeyOptions{}
		if !m.StorageCredential.ExternalID.IsNull() && !m.StorageCredential.ExternalID.IsUnknown() {
			opts = append(opts, credential.WithExternalID(m.StorageCredential.ExternalID.ValueString()))
		}
		creds := credential.NewS3CredentialAccessKey(
			m.StorageCredential.AccessKeyID.ValueString(),
			m.StorageCredential.SecretAccessKey.ValueString(),
			opts...,
		)

		return creds, nil
	case "s3_aws_system_identity":
		return credential.NewS3CredentialSystemIdentity(m.StorageCredential.AccountID.ValueString()), nil
	case "s3_cloudflare_r2":
		return credential.NewCloudflareR2Credential(
			m.StorageCredential.AccessKeyID.ValueString(),
			m.StorageCredential.SecretAccessKey.ValueString(),
			m.StorageCredential.AccountID.ValueString(),
			m.StorageCredential.Token.ValueString(),
		), nil
	case "az_client_credentials":
		return credential.NewAZCredentialClientCredentials(
			m.StorageCredential.ClientID.ValueString(),
			m.StorageCredential.ClientSecret.ValueString(),
			m.StorageCredential.TenantID.ValueString(),
		), nil
	case "az_shared_access_key":
		return credential.NewAZCredentialSharedAccessKey(
			m.StorageCredential.AZKey.ValueString(),
		), nil
	case "az_azure_system_identity":
		return credential.NewAZCredentialManagedIdentity(), nil
	case "gcs_service_account_key":
		return credential.NewGCSCredentialServiceAccountKey(
			credential.GCSServiceKey{
				AuthProviderX509CertURL: m.StorageCredential.Key.AuthProviderX509CertURL.ValueString(),
				AuthURI:                 m.StorageCredential.Key.AuthURI.ValueString(),
				ClientEmail:             m.StorageCredential.Key.ClientEmail.ValueString(),
				ClientID:                m.StorageCredential.Key.ClientID.ValueString(),
				ClientX509CertURL:       m.StorageCredential.Key.ClientX509CertURL.ValueString(),
				PrivateKey:              m.StorageCredential.Key.PrivateKey.ValueString(),
				PrivateKeyID:            m.StorageCredential.Key.PrivateKeyID.ValueString(),
				ProjectID:               m.StorageCredential.Key.ProjectID.ValueString(),
				TokenURI:                m.StorageCredential.Key.TokenURI.ValueString(),
				Type:                    m.StorageCredential.Key.Type.ValueString(),
				UniverseDomain:          m.StorageCredential.Key.UniverseDomain.ValueString(),
			},
		), nil
	case "gcs_gcp_system_identity":
		return credential.NewGCSCredentialSystemIdentity(), nil
	default:
		return nil, fmt.Errorf("incorrect storage credential definition, type must be one of %v, got %s", tftypes.ValidStorageCredentialTypes, storageType)
	}
}
