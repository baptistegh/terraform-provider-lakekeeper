package provider

import (
	"context"
	"fmt"

	managementv1 "github.com/baptistegh/go-lakekeeper/pkg/apis/management/v1"
	lakekeeper "github.com/baptistegh/go-lakekeeper/pkg/client"
	"github.com/baptistegh/go-lakekeeper/pkg/core"
	tftypes "github.com/baptistegh/terraform-provider-lakekeeper/internal/provider/types"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Ensure provider defined types fully satisfy framework interfaces
var (
	_ resource.Resource              = &lakekeeperWarehouseResource{}
	_ resource.ResourceWithConfigure = &lakekeeperWarehouseResource{}

	// Lakekeeper API does not give storage credentials on GET
	// we can't activate the import on warehouse for now
	// _ resource.ResourceWithImportState = &lakekeeperWarehouseResource{}
)

func init() {
	registerResource(NewLakekeeperWarehouseResource)
}

// NewLakekeeperWarehouseResource is a helper function to simplify the provider implementation.
func NewLakekeeperWarehouseResource() resource.Resource {
	return &lakekeeperWarehouseResource{}
}

func (r *lakekeeperWarehouseResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_warehouse"
}

// lakekeeperWarehouseResource defines the resource implementation.
type lakekeeperWarehouseResource struct {
	client *lakekeeper.Client
}

// lakekeeperWarehouseResourceModel describes the resource data model.
type lakekeeperWarehouseResourceModel struct {
	ID                types.String                    `tfsdk:"id"` // form: project_id:warehouse_id (internal ID)
	WarehouseID       types.String                    `tfsdk:"warehouse_id"`
	Name              types.String                    `tfsdk:"name"`
	ProjectID         types.String                    `tfsdk:"project_id"` // Optional, if not provided, the default project will be used.
	Protected         types.Bool                      `tfsdk:"protected"`
	Active            types.Bool                      `tfsdk:"active"`
	StorageProfile    *tftypes.StorageProfileModel    `tfsdk:"storage_profile"`
	DeleteProfile     *tftypes.DeleteProfileModel     `tfsdk:"delete_profile"`
	StorageCredential *tftypes.StorageCredentialModel `tfsdk:"storage_credential"`
}

func (r *lakekeeperWarehouseResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: fmt.Sprintf(`The ` + "`lakekeeper_warehouse`" + ` resource allows to manage the lifecycle of a lakekeeper warehouse.

**Upstream API**: [Lakekeeper REST API docs](https://docs.lakekeeper.io/docs/nightly/api/management/#tag/warehouse)`),
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The internale ID the warehouse. In the form: <project_id>:<warehouse_id>",
				Computed:            true,
			},
			"warehouse_id": schema.StringAttribute{
				MarkdownDescription: "The ID the warehouse.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the warehouse.",
				Required:            true,
				Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"project_id": schema.StringAttribute{
				MarkdownDescription: "The project ID to which the warehouse belongs. If not provided, the default project will be used.",
				Required:            true,
				Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"protected": schema.BoolAttribute{
				MarkdownDescription: "Whether the warehouse is protected from being deleted.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"active": schema.BoolAttribute{
				MarkdownDescription: "Whether the warehouse is active.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"storage_profile":    tftypes.StorageProfileResourceSchema(),
			"delete_profile":     tftypes.DeleteProfileResourceSchema(),
			"storage_credential": tftypes.StorageCredentialSchema(),
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *lakekeeperWarehouseResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	resourceData := req.ProviderData.(*LakekeeperResourceData)
	r.client = resourceData.Client
}

// Create creates a new upstream resources and adds it into the Terraform state.
func (r *lakekeeperWarehouseResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state *lakekeeperWarehouseResourceModel
	var plan types.Object

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(plan.As(ctx, &state, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    true,
		UnhandledUnknownAsEmpty: true,
	})...)
	if resp.Diagnostics.HasError() {
		return
	}

	opts, err := state.ToWarehouseCreateRequest()
	if err != nil {
		resp.Diagnostics.AddError("Error decoding state to model", fmt.Sprintf("Incorrect Warehouse creation request, %v", err))
		return
	}

	w, _, err := r.client.WarehouseV1(state.ProjectID.ValueString()).Create(opts, core.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError("Lakekeeper API error occurred",
			fmt.Sprintf("Unable to create warehouse: %v", err))
		return
	}

	warehouse, _, err := r.client.WarehouseV1(state.ProjectID.ValueString()).Get(w.ID, core.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError("Lakekeeper API error occurred",
			fmt.Sprintf("Unable to read warehouse: %v", err))
		return
	}

	diags := state.RefreshFromSettings(warehouse)
	if diags.HasError() {
		resp.Diagnostics = append(resp.Diagnostics, diags...)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *lakekeeperWarehouseResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state lakekeeperWarehouseResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	projectID, warehouseID := splitInternalID(state.ID)
	warehouse, _, err := r.client.WarehouseV1(projectID).Get(warehouseID, core.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError("Lakekeeper API error occurred", fmt.Sprintf("Unable to read warehouse %s in project %s, %s", warehouseID, projectID, err.Error()))
		return
	}

	diags := state.RefreshFromSettings(warehouse)
	if diags.HasError() {
		resp.Diagnostics = append(resp.Diagnostics, diags...)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Updates updates the resource in-place.
func (r *lakekeeperWarehouseResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state lakekeeperWarehouseResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	projectID, warehouseID := splitInternalID(state.ID)

	// Deactivate the warehouse if the active field is set to false
	if !plan.Active.IsNull() && !plan.Active.ValueBool() {
		if _, err := r.client.WarehouseV1(projectID).Deactivate(warehouseID, core.WithContext(ctx)); err != nil {
			resp.Diagnostics.AddError("Lakekeeper API error occurred", fmt.Sprintf("Unable to deactivate warehouse %s in project %s, %s", warehouseID, projectID, err.Error()))
			return
		}
	}

	// Rename the warehouse if the name field is different
	if !plan.Name.IsNull() && plan.Name.ValueString() != state.Name.ValueString() {
		if _, err := r.client.WarehouseV1(projectID).Rename(warehouseID, &managementv1.RenameWarehouseOptions{
			NewName: plan.Name.ValueString(),
		}, core.WithContext(ctx)); err != nil {
			resp.Diagnostics.AddError("Lakekeeper API error occurred", fmt.Sprintf("Unable to rename warehouse %s in project %s, %s", warehouseID, projectID, err.Error()))
			return
		}
	}

	// Set warehouse protection if the protected field is different
	if plan.Protected.ValueBool() != state.Protected.ValueBool() {
		if _, _, err := r.client.WarehouseV1(projectID).SetProtection(warehouseID, plan.Protected.ValueBool()); err != nil {
			resp.Diagnostics.AddError("Lakekeeper API error occurred", fmt.Sprintf("Unable to set protection for warehouse %s in project %s, %s", warehouseID, projectID, err.Error()))
			return
		}
	}

	opts, err := plan.ToWarehouseCreateRequest()
	if err != nil {
		resp.Diagnostics.AddError("Error decoding plan to model", fmt.Sprintf("Incorrect Warehouse update request, %v", err))
		return
	}

	// Update the delete profile
	if !plan.DeleteProfile.Type.Equal(state.DeleteProfile.Type) || !plan.DeleteProfile.ExpirationSeconds.Equal(state.DeleteProfile.ExpirationSeconds) {
		if _, err := r.client.WarehouseV1(projectID).UpdateDeleteProfile(warehouseID, &managementv1.UpdateDeleteProfileOptions{
			DeleteProfile: *opts.DeleteProfile,
		}, core.WithContext(ctx)); err != nil {
			resp.Diagnostics.AddError("Lakekeeper API error occurred", fmt.Sprintf("Unable to update delete profile for warehouse %s in project %s, %s", warehouseID, projectID, err.Error()))
			return
		}
	}

	// Update the storage profile and its storage credential
	if _, err := r.client.WarehouseV1(projectID).UpdateStorageProfile(warehouseID, &managementv1.UpdateStorageProfileOptions{
		StorageProfile:    opts.StorageProfile,
		StorageCredential: &opts.StorageCredential,
	}, core.WithContext(ctx)); err != nil {
		resp.Diagnostics.AddError("Lakekeeper API error occurred", fmt.Sprintf("Unable to update storage profile for warehouse %s in project %s, %s", warehouseID, projectID, err.Error()))
		return
	}

	// Refresh the state with the updated warehouse settings
	warehouse, _, err := r.client.WarehouseV1(projectID).Get(warehouseID, core.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError("Lakekeeper API error occurred", fmt.Sprintf("Unable to read warehouse %s in project %s, %s", warehouseID, projectID, err.Error()))
		return
	}

	diags := state.RefreshFromSettings(warehouse)
	if diags.HasError() {
		resp.Diagnostics = append(resp.Diagnostics, diags...)
		return
	}

	// Lakekeeper API does not return storage credentials on GET, so we need to set it manually
	state.StorageCredential = plan.StorageCredential

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Deletes removes the resource.
func (r *lakekeeperWarehouseResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state lakekeeperWarehouseResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	projectID, warehouseID := splitInternalID(state.ID)

	opts := managementv1.DeleteWarehouseOptions{
		Force: core.Ptr(true),
	}

	if _, err := r.client.WarehouseV1(projectID).Delete(warehouseID, &opts, core.WithContext(ctx)); err != nil {
		resp.Diagnostics.AddError("Lakekeeper API error occurred", fmt.Sprintf("Unable to delete warehouse %s in project %s, %s", warehouseID, projectID, err.Error()))
		return
	}

	resp.State.RemoveResource(ctx)
}
