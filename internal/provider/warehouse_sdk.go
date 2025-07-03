package provider

import (
	"fmt"

	tftypes "github.com/baptistegh/terraform-provider-lakekeeper/internal/provider/types"
	"github.com/baptistegh/terraform-provider-lakekeeper/lakekeeper"
	"github.com/baptistegh/terraform-provider-lakekeeper/lakekeeper/storage"
	"github.com/baptistegh/terraform-provider-lakekeeper/lakekeeper/storage/credential"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (m *lakekeeperWarehouseResourceModel) ToWarehouseCreateRequest() (*lakekeeper.WarehouseCreateOptions, error) {
	if !m.Active.ValueBool() {
		return nil, fmt.Errorf("could not create a warehouse with inactive status")
	}
	req := &lakekeeper.WarehouseCreateOptions{
		Name:      m.Name.ValueString(),
		Protected: m.Protected.ValueBool(),
		Status:    "active",
	}

	if m.DeleteProfile != nil {
		deleteProfile, err := m.DeleteProfileSettings()
		if err != nil {
			return nil, err
		}
		req.DeleteProfile = deleteProfile
	}

	if m.StorageProfile != nil {
		storageProfile, err := m.StorageProfileSettings()
		if err != nil {
			return nil, err
		}
		req.StorageProfile = storageProfile
	}

	if m.StorageCredential != nil {
		storageCredential, err := m.StorageCredentialSettings()
		if err != nil {
			return nil, err
		}
		req.StorageCredential = storageCredential
	}

	return req, nil
}

func (m *lakekeeperWarehouseResourceModel) RefreshFromSettings(diags *diag.Diagnostics, w *lakekeeper.Warehouse) {
	m.ID = types.StringValue(w.ProjectID + ":" + w.ID)
	m.WarehouseID = types.StringValue(w.ID)
	m.ProjectID = types.StringValue(w.ProjectID)
	m.Protected = types.BoolValue(w.Protected)
	m.Active = types.BoolValue(w.IsActive())
	m.Name = types.StringValue(w.Name)

	m.refreshStorageProfile(diags, w.StorageProfileWrapper)
	m.refreshStorageCredential(w.StorageCredentialWrapper)
	m.refreshDeleteProfile(diags, w.DeleteProfileWrapper)
}

func (m *lakekeeperWarehouseResourceModel) DeleteProfileSettings() (lakekeeper.DeleteProfile, error) {
	if m.DeleteProfile == nil {
		return nil, nil
	}

	switch m.DeleteProfile.Type.ValueString() {
	case "soft":
		return &lakekeeper.SoftDeleteProfile{
			Type:           "soft",
			ExpiredSeconds: m.DeleteProfile.ExpirationSeconds.ValueInt32(),
		}, nil
	case "hard":
		return &lakekeeper.HardDeleteProfile{
			Type: "hard",
		}, nil
	default:
		return nil, fmt.Errorf("incorrect delete profile definition, type must be [soft,hard]")
	}
}

func (m *lakekeeperWarehouseResourceModel) StorageProfileSettings() (storage.StorageProfile, error) {
	if m.StorageProfile == nil {
		return nil, nil
	}
	switch m.StorageProfile.Type.ValueString() {
	case "s3":
		profile := storage.NewStorageProfileS3(
			m.StorageProfile.Bucket.ValueString(),
			m.StorageProfile.Region.ValueString(),
			m.StorageProfile.STSEnabled.ValueBool(),
		)
		if !m.StorageProfile.Endpoint.IsNull() && !m.StorageProfile.Endpoint.IsUnknown() {
			profile.Endpoint = m.StorageProfile.Endpoint.ValueStringPointer()
		} else {
			profile.Endpoint = nil
		}

		if !m.StorageProfile.PathStyleAccess.IsNull() && !m.StorageProfile.PathStyleAccess.IsUnknown() {
			profile.PathStyleAccess = m.StorageProfile.PathStyleAccess.ValueBoolPointer()
		} else {
			profile.PathStyleAccess = nil
		}

		if !m.StorageProfile.KeyPrefix.IsNull() && !m.StorageProfile.KeyPrefix.IsUnknown() {
			profile.KeyPrefix = m.StorageProfile.KeyPrefix.ValueStringPointer()
		} else {
			profile.KeyPrefix = nil
		}
		return profile, nil
	case "adls":
		return storage.NewStorageProfileADLS(
			m.StorageProfile.AccountName.ValueString(),
			m.StorageProfile.Filesystem.ValueString(),
		), nil
	case "gcs":
		return storage.NewStorageProfileGCS(
			m.StorageProfile.Bucket.ValueString(),
		), nil
	}
	return nil, fmt.Errorf("invalid storage profile definitions, type must be [s3,gcs,adls]")
}

func (m *lakekeeperWarehouseResourceModel) StorageCredentialSettings() (credential.StorageCredential, error) {
	if m.StorageCredential == nil {
		return nil, fmt.Errorf("invalid storage credential definitions, must be defined")
	}
	storageType := m.StorageCredential.Type.ValueString()
	switch storageType {
	case "s3_access_key":
		return credential.NewS3CredentialAccessKey(
			m.StorageCredential.AccessKeyID.ValueString(),
			m.StorageCredential.SecretAccessKey.ValueString(),
			m.StorageCredential.ExternalID.ValueString(),
		), nil
	default:
		return nil, fmt.Errorf("incorrect storage credential definition, type must be one of %v, got %s", tftypes.ValidStorageCredentialTypes, storageType)
	}
}

func (m *lakekeeperWarehouseResourceModel) refreshDeleteProfile(diags *diag.Diagnostics, d *lakekeeper.DeleteProfileWrapper) {
	if d == nil {
		m.DeleteProfile = nil
		return
	}
	switch profile := d.DeleteProfile.(type) {
	case lakekeeper.SoftDeleteProfile:
		m.DeleteProfile = &tftypes.DeleteProfileModel{
			Type:              types.StringValue("soft"),
			ExpirationSeconds: types.Int32Value(profile.ExpiredSeconds),
		}
	case lakekeeper.HardDeleteProfile:
		m.DeleteProfile = &tftypes.DeleteProfileModel{
			Type: types.StringValue("hard"),
		}
	default:
		diags.AddError("Error converting model to state", fmt.Sprintf("Incorrect delete profile type: %T, valid: [soft,hard]", d))
	}
}

func (m *lakekeeperWarehouseResourceModel) refreshStorageProfile(diags *diag.Diagnostics, d *storage.StorageProfileWrapper) {
	if d == nil {
		diags.AddError("Error converting model to state", "Incorrect storage profile: must be defined")
		m.StorageProfile = nil
		return
	}
	if m.StorageProfile == nil {
		m.StorageProfile = &tftypes.StorageProfileModel{}
	}
	m.StorageProfile.RefreshFromSettings(d.StorageProfile)
}

func (m *lakekeeperWarehouseResourceModel) refreshStorageCredential(d *credential.StorageCredentialWrapper) {
	if d == nil {
		return
	}
	if m.StorageCredential == nil {
		m.StorageCredential = &tftypes.StorageCredentialModel{}
	}
	m.StorageCredential.RefreshFromSettings(d.StorageCredential)
}
