package provider

import (
	"errors"
	"fmt"

	"github.com/baptistegh/terraform-provider-lakekeeper/internal/provider/sdk"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	managementv1 "github.com/lakekeeper/go-lakekeeper/pkg/apis/management/v1"
)

func (m *lakekeeperWarehouseResourceModel) toWarehouseCreateRequest() (*managementv1.CreateWarehouseRequest, error) {
	req := managementv1.CreateWarehouseRequest{
		WarehouseName: m.Name.ValueString(),
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
	req.StorageProfile = storage

	creds, err := s.CredentialAsSDK()
	if err != nil {
		return nil, err
	}
	req.SetStorageCredential(creds)

	if m.DeleteProfile != nil {
		settings, err := m.DeleteProfileSettings()
		if err != nil {
			return nil, err
		}
		if settings == nil {
			return nil, errors.New("delete profile is empty")
		}

		req.DeleteProfile = settings
	}

	return &req, nil
}

// TODO: refactor RefreshFromSettings on datasource and resource
// because these functions are almost identical

func (m *lakekeeperWarehouseResourceModel) RefreshFromSettings(w *managementv1.GetWarehouseResponse, plan *lakekeeperWarehouseResourceModel) diag.Diagnostics {
	m.ID = types.StringValue(w.ProjectId + "/" + w.Id)
	m.WarehouseID = types.StringValue(w.Id)
	m.ProjectID = types.StringValue(w.ProjectId)
	m.Protected = types.BoolValue(w.Protected)
	m.Active = types.BoolValue(w.Status == managementv1.WAREHOUSESTATUS_ACTIVE)
	m.Name = types.StringValue(w.Name)

	diags := diag.Diagnostics{}
	const errorMessage = "Error refreshing warehouse state"

	var oldProfile sdk.StorageProfileModel
	if plan != nil && plan.StorageProfile != nil {
		s, err := sdk.OnlyOneStorageProfile(plan.StorageProfile.S3StorageProfile, plan.StorageProfile.GCSStorageProfile, plan.StorageProfile.ADLSStorageProfile)
		if err != nil {
			diags.AddError(errorMessage, err.Error())
			return diags
		}
		oldProfile = s
	} else if m.StorageProfile != nil {
		s, err := sdk.OnlyOneStorageProfile(m.StorageProfile.S3StorageProfile, m.StorageProfile.GCSStorageProfile, m.StorageProfile.ADLSStorageProfile)
		if err != nil {
			diags.AddError(errorMessage, err.Error())
			return diags
		}
		oldProfile = s
	} else {
		oldProfile = nil
	}

	storageProfile, err := sdk.StorageProfileModelFromSDK(w.StorageProfile)
	if err != nil {
		diags.AddError(errorMessage, err.Error())
		return diags
	}

	if oldProfile != nil {
		creds, err := oldProfile.GetCredentials()
		if err != nil {
			diags.AddError(errorMessage, err.Error())
			return diags
		}

		if err := storageProfile.AddCreds(creds); err != nil {
			diags.AddError(errorMessage, err.Error())
			return diags
		}
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

	m.DeleteProfile = deleteProfileModelFromSDK(w.DeleteProfile)

	return diags
}

func (m *lakekeeperWarehouseDataSourceModel) RefreshDataSourceFromSettings(w *managementv1.GetWarehouseResponse) diag.Diagnostics {
	m.ID = types.StringValue(w.ProjectId + "/" + w.Id)
	m.WarehouseID = types.StringValue(w.Id)
	m.ProjectID = types.StringValue(w.ProjectId)
	m.Protected = types.BoolValue(w.Protected)
	m.Active = types.BoolValue(w.Status == managementv1.WAREHOUSESTATUS_ACTIVE)
	m.Name = types.StringValue(w.Name)

	diags := diag.Diagnostics{}
	const errorMessage = "Error refreshing warehouse state"

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

	m.DeleteProfile = deleteProfileModelFromSDK(w.DeleteProfile)

	return diags
}

func (m *lakekeeperWarehouseResourceModel) DeleteProfileSettings() (*managementv1.TabularDeleteProfile, error) {
	if m.DeleteProfile == nil {
		return nil, nil
	}

	hard := managementv1.TabularDeleteProfileHardAsTabularDeleteProfile(managementv1.NewTabularDeleteProfileHard(sdk.DeleteProfileTypeHard))

	if m.DeleteProfile.Type.IsNull() || m.DeleteProfile.Type.IsUnknown() {
		return &hard, nil
	}

	switch m.DeleteProfile.Type.ValueString() {
	case sdk.DeleteProfileTypeSoft:
		soft := managementv1.TabularDeleteProfileSoftAsTabularDeleteProfile(managementv1.NewTabularDeleteProfileSoft(
			int64(m.DeleteProfile.ExpirationSeconds.ValueInt32()),
			sdk.DeleteProfileTypeSoft,
		))
		return &soft, nil
	case sdk.DeleteProfileTypeHard:
		return &hard, nil
	default:
		return nil, fmt.Errorf("incorrect delete profile definition, type must be [soft,hard]")
	}
}

// deleteProfileModelFromSDK converts a tabular delete profile, or returns nil when none is set.
func deleteProfileModelFromSDK(p managementv1.TabularDeleteProfile) *sdk.DeleteProfileModel {
	switch {
	case p.TabularDeleteProfileSoft != nil:
		return &sdk.DeleteProfileModel{
			Type:              types.StringValue(sdk.DeleteProfileTypeSoft),
			ExpirationSeconds: types.Int32Value(int32(p.TabularDeleteProfileSoft.ExpirationSeconds)),
		}
	case p.TabularDeleteProfileHard != nil:
		return &sdk.DeleteProfileModel{
			Type: types.StringValue(sdk.DeleteProfileTypeHard),
		}
	default:
		return nil
	}
}
