package provider

import (
	"errors"
	"fmt"

	"github.com/baptistegh/terraform-provider-lakekeeper/internal/provider/sdk"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	managementv1 "github.com/baptistegh/go-lakekeeper/pkg/apis/management/v1"
	"github.com/baptistegh/go-lakekeeper/pkg/apis/management/v1/storage/profile"
)

func (m *lakekeeperWarehouseResourceModel) toWarehouseCreateRequest() (*managementv1.CreateWarehouseOptions, error) {
	req := managementv1.CreateWarehouseOptions{
		Name: m.Name.ValueString(),
	}

	if m.StorageProfile == nil {
		return nil, errors.New("storage profile is required")
	}

	s, err := sdk.OnlyOneStorageProfile(m.StorageProfile.S3StorageProfile, m.StorageProfile.ADLSStorageProfile, m.StorageProfile.GCSStorageProfile)
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
	if plan != nil && plan.StorageProfile != nil {
		s, err := sdk.OnlyOneStorageProfile(plan.StorageProfile.S3StorageProfile, plan.StorageProfile.GCSStorageProfile, plan.StorageProfile.ADLSStorageProfile)
		if err != nil {
			diags.AddError(errorMessage, err.Error())
			return diags
		}
		oldProfile = s
	} else {
		s, err := sdk.OnlyOneStorageProfile(m.StorageProfile.S3StorageProfile, m.StorageProfile.GCSStorageProfile, m.StorageProfile.ADLSStorageProfile)
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

	m.StorageProfile = &storageProfileWrapper{}

	switch sp := storageProfile.(type) {
	case *sdk.S3StorageProfileModel:
		m.StorageProfile.S3StorageProfile = sp
	case *sdk.ADLSStorageProfileModel:
		m.StorageProfile.ADLSStorageProfile = sp
	case *sdk.GCSStorageProfileModel:
		m.StorageProfile.GCSStorageProfile = sp
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

func (m *lakekeeperWarehouseDataSourceModel) RefreshDataSourceFromSettings(w *managementv1.Warehouse) diag.Diagnostics {
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

	storageProfile, err := sdk.StorageProfileModelFromSDK(w.StorageProfile)
	if err != nil {
		diags.AddError(errorMessage, err.Error())
		return diags
	}

	m.StorageProfile = &storageProfileDataSourceWrapper{}

	switch sp := storageProfile.(type) {
	case *sdk.S3StorageProfileModel:
		m.StorageProfile.S3StorageProfile = &sdk.S3StorageProfileDataSourceModel{}
		if err := m.StorageProfile.S3StorageProfile.FromResourceModel(sp); err != nil {
			diags.AddError(errorMessage, err.Error())
		}
	case *sdk.ADLSStorageProfileModel:
		m.StorageProfile.ADLSStorageProfile = &sdk.ADLSStorageProfileDataSourceModel{}
		if err := m.StorageProfile.ADLSStorageProfile.FromResourceModel(sp); err != nil {
			diags.AddError(errorMessage, err.Error())
		}
	case *sdk.GCSStorageProfileModel:
		m.StorageProfile.GCSStorageProfile = &sdk.GCSStorageProfileDataSourceModel{}
		if err := m.StorageProfile.GCSStorageProfile.FromResourceModel(sp); err != nil {
			diags.AddError(errorMessage, err.Error())
		}
	default:
		diags.AddError(errorMessage, fmt.Sprintf("Incorrect storage profile type, got %T", storageProfile))
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
