package provider

import (
	"context"
	"fmt"
	"strings"

	managementv1 "github.com/baptistegh/go-lakekeeper/pkg/apis/management/v1"
	lakekeeper "github.com/baptistegh/go-lakekeeper/pkg/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces
var (
	_ resource.Resource                = &lakekeeperRoleResource{}
	_ resource.ResourceWithConfigure   = &lakekeeperRoleResource{}
	_ resource.ResourceWithImportState = &lakekeeperRoleResource{}
)

type LakekeeperRoleResourceModel struct {
	ID          types.String `tfsdk:"id"`
	RoleID      types.String `tfsdk:"role_id"`
	ProjectID   types.String `tfsdk:"project_id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`

	CreatedAt types.String `tfsdk:"created_at"`
	UpdatedAt types.String `tfsdk:"updated_at"`
}

func init() {
	registerResource(NewLakekeeperRoleResource)
}

// NewLakekeeperRoleResource is a helper function to simplify the provider implementation.
func NewLakekeeperRoleResource() resource.Resource {
	return &lakekeeperRoleResource{}
}

func (r *lakekeeperRoleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_role"
}

// lakekeeperRoleResource defines the resource implementation.
type lakekeeperRoleResource struct {
	client *lakekeeper.Client
}

func (r *lakekeeperRoleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: fmt.Sprintf(`The ` + "`lakekeeper_role`" + ` resource allows to manage the lifecycle of a lakekeeper role.

**Upstream API**: [Lakekeeper REST API docs](https://docs.lakekeeper.io/docs/nightly/api/management/#tag/role)`),

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of the role. in the form `{{project_id}}/{{role_id}}`",
				Computed:            true,
			},
			"role_id": schema.StringAttribute{
				MarkdownDescription: `The internal ID of the role.`,
				Computed:            true,
			},
			"project_id": schema.StringAttribute{
				MarkdownDescription: `The ID of the project the role belongs to.`,
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the role.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the role.",
				Optional:            true,
			},
			"created_at": schema.StringAttribute{
				MarkdownDescription: "When the role has been created.",
				Computed:            true,
			},
			"updated_at": schema.StringAttribute{
				MarkdownDescription: "When the role has last been modified.",
				Computed:            true,
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *lakekeeperRoleResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	resourceData := req.ProviderData.(*LakekeeperResourceData)
	r.client = resourceData.Client
}

// Create creates a new upstream resources and adds it into the Terraform state.
func (r *lakekeeperRoleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state LakekeeperRoleResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	opts := managementv1.CreateRoleOptions{
		Name:        state.Name.ValueString(),
		Description: state.Description.ValueStringPointer(),
	}

	role, _, err := r.client.RoleV1(state.ProjectID.ValueString()).Create(ctx, &opts)
	if err != nil {
		resp.Diagnostics.AddError("Lakekeeper API error occurred", fmt.Sprintf("Unable to create role, %v", err))
		return
	}

	state.RoleID = types.StringValue(role.ID)
	state.ProjectID = types.StringValue(role.ProjectID)

	state.ID = types.StringValue(fmt.Sprintf("%s/%s", role.ProjectID, role.ID))

	state.CreatedAt = types.StringValue(role.CreatedAt)
	state.UpdatedAt = types.StringPointerValue(role.UpdatedAt)
	state.Description = types.StringPointerValue(role.Description)

	// Log the creation of the resource
	tflog.Debug(ctx, "created a role", map[string]any{
		"name": state.Name.ValueString(), "id": state.ID.ValueString(),
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *lakekeeperRoleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state LakekeeperRoleResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	projectID, roleID := splitInternalID(state.ID)

	role, _, err := r.client.RoleV1(projectID).Get(ctx, roleID)
	if err != nil {
		resp.Diagnostics.AddError("Lakekeeper API error occurred", fmt.Sprintf("Unable to read role %s in project %s, %v", roleID, projectID, err))
		return
	}

	state.RoleID = types.StringValue(role.ID)
	state.ProjectID = types.StringValue(role.ProjectID)

	state.ID = types.StringValue(fmt.Sprintf("%s/%s", role.ProjectID, role.ID))

	state.Name = types.StringValue(role.Name)
	state.CreatedAt = types.StringValue(role.CreatedAt)

	state.Description = types.StringPointerValue(role.Description)
	state.UpdatedAt = types.StringPointerValue(role.UpdatedAt)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Updates updates the resource in-place.
func (r *lakekeeperRoleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state LakekeeperRoleResourceModel
	var plan LakekeeperRoleResourceModel

	// Read Terraform state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if state.ID.IsNull() || state.ID.IsUnknown() || state.ID.ValueString() == "" {
		resp.Diagnostics.AddError("Incorrect resource definition", "Resource was requested to perform an in-place upgrade with a null ID.")
		return
	}

	projectID, roleID := splitInternalID(state.ID)

	opts := managementv1.UpdateRoleOptions{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueStringPointer(),
	}

	role, _, err := r.client.RoleV1(projectID).Update(ctx, roleID, &opts)
	if err != nil {
		resp.Diagnostics.AddError("Lakekeeper API error occurred", fmt.Sprintf("Unable to update role %s in project %s, %v", roleID, projectID, err))
		return
	}

	state.RoleID = types.StringValue(role.ID)
	state.ProjectID = types.StringValue(role.ProjectID)
	state.Name = types.StringValue(role.Name)
	state.ID = types.StringValue(fmt.Sprintf("%s/%s", role.ProjectID, role.ID))

	state.CreatedAt = types.StringValue(role.CreatedAt)
	state.UpdatedAt = types.StringPointerValue(role.UpdatedAt)
	state.Description = types.StringPointerValue(role.Description)

	// Log the creation of the resource
	tflog.Debug(ctx, "created a role", map[string]any{
		"name": state.Name.ValueString(), "id": state.ID.ValueString(),
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Deletes removes the resource.
func (r *lakekeeperRoleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state LakekeeperRoleResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	projectID, roleID := splitInternalID(state.ID)

	_, err := r.client.RoleV1(projectID).Delete(ctx, roleID)
	if err != nil {
		resp.Diagnostics.AddError("Lakekeeper API error occurred", fmt.Sprintf("Unable to delete role %s in project %s, %v", roleID, projectID, err))
		return
	}

	resp.State.RemoveResource(ctx)
}

// ImportState imports the resource into the Terraform state.
func (r *lakekeeperRoleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Expected format: "project_id/role_id"
	parts := strings.Split(req.ID, "/")
	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			"Invalid import ID format",
			"Expected format: project_id/role_id",
		)
		return
	}

	resp.State.SetAttribute(ctx, path.Root("project_id"), parts[0])
	resp.State.SetAttribute(ctx, path.Root("role_id"), parts[1])

	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
