package provider

import (
	"errors"
	"fmt"

	"github.com/baptistegh/terraform-provider-lakekeeper/internal/provider/sdk"
	"github.com/baptistegh/terraform-provider-lakekeeper/internal/provider/sdk/deprecated"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	managementv1 "github.com/baptistegh/go-lakekeeper/pkg/apis/management/v1"
	"github.com/baptistegh/go-lakekeeper/pkg/apis/management/v1/storage/profile"
)

func (m *lakekeeperWarehouseResourceModel) toWarehouseCreateRequest() (*managementv1.CreateWarehouseOptions, error) {
	req := managementv1.CreateWarehouseOptions{
		Name: m.Name.ValueString(),
	}

	s, err := sdk.OnlyOneStorageProfile(m.S3StorageProfile, m.ADLSStorageProfile, m.GCSStorageProfile)
	if err != nil {
		return nil, err
	}

	storage, err := s.AsSDK()
	if err != nil {
		return nil, err
	}
	req.StorageProfile = storage.AsProfile()

	creds, err := s.CredentialAsSDK()
	if err != nil {
		return nil, err
	}
	req.StorageCredential = creds.AsCredential()

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

	return &req, nil
}

// TODO: refactor RefreshFromSettings on datasource and resource
// because these functions are almost identical

func (m *lakekeeperWarehouseResourceModel) RefreshFromSettings(w *managementv1.Warehouse, plan *lakekeeperWarehouseResourceModel) diag.Diagnostics {
	m.ID = types.StringValue(w.ProjectID + "/" + w.ID)
	m.WarehouseID = types.StringValue(w.ID)
	m.ProjectID = types.StringValue(w.ProjectID)
	m.Protected = types.BoolValue(w.Protected)
	m.Active = types.BoolValue(w.IsActive())
	m.Name = types.StringValue(w.Name)

	diags := diag.Diagnostics{}
	const errorMessage = "Error refreshing warehouse state"

	if w.StorageProfile.StorageSettings == nil {
		diags.AddError(errorMessage, "Storage profile must be defined")
		return diags
	}

	var oldProfile sdk.StorageProfileModel
	if plan != nil {
		s, err := sdk.OnlyOneStorageProfile(plan.S3StorageProfile, plan.GCSStorageProfile, plan.ADLSStorageProfile)
		if err != nil {
			diags.AddError(errorMessage, err.Error())
			return diags
		}
		oldProfile = s
	} else {
		s, err := sdk.OnlyOneStorageProfile(m.S3StorageProfile, m.GCSStorageProfile, m.ADLSStorageProfile)
		if err != nil {
			diags.AddError(errorMessage, err.Error())
			return diags
		}
		oldProfile = s
	}

	creds, err := oldProfile.GetCredentials()
	if err != nil {
		diags.AddError(errorMessage, err.Error())
		return diags
	}

	storageProfile, err := sdk.StorageProfileModelFromSDK(w.StorageProfile)
	if err != nil {
		diags.AddError(errorMessage, err.Error())
		return diags
	}

	if err := storageProfile.AddCreds(creds); err != nil {
		diags.AddError(errorMessage, err.Error())
		return diags
	}

	switch sp := storageProfile.(type) {
	case *sdk.S3StorageProfileModel:
		m.S3StorageProfile = sp
		m.ADLSStorageProfile = nil
		m.GCSStorageProfile = nil
	case *sdk.ADLSStorageProfileModel:
		m.ADLSStorageProfile = sp
		m.S3StorageProfile = nil
		m.GCSStorageProfile = nil
	case *sdk.GCSStorageProfileModel:
		m.GCSStorageProfile = sp
		m.ADLSStorageProfile = nil
		m.S3StorageProfile = nil
	default:
		diags.AddError(errorMessage, "Incorrect storage profile type")
	}

	if w.DeleteProfile == nil || w.DeleteProfile.DeleteProfileSettings == nil {
		m.DeleteProfile = nil
	} else {
		switch deleteProfile := w.DeleteProfile.DeleteProfileSettings.(type) {
		case *profile.TabularDeleteProfileSoft:
			m.DeleteProfile = &sdk.DeleteProfileModel{
				Type:              types.StringValue("soft"),
				ExpirationSeconds: types.Int32Value(deleteProfile.ExpirationSeconds),
			}
		case *profile.TabularDeleteProfileHard:
			m.DeleteProfile = &sdk.DeleteProfileModel{
				Type: types.StringValue("hard"),
			}
		default:
			diags.AddError(errorMessage, fmt.Sprintf("Incorrect delete profile type: %T, valid: [soft,hard]", deleteProfile))
		}
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
		m.StorageProfile = &deprecated.StorageProfileModel{}
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
		m.DeleteProfile = &sdk.DeleteProfileModel{}
		switch deleteProfile := w.DeleteProfile.DeleteProfileSettings.(type) {
		case *profile.TabularDeleteProfileSoft:
			m.DeleteProfile = &sdk.DeleteProfileModel{
				Type:              types.StringValue("soft"),
				ExpirationSeconds: types.Int32Value(deleteProfile.ExpirationSeconds),
			}
		case *profile.TabularDeleteProfileHard:
			m.DeleteProfile = &sdk.DeleteProfileModel{
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
