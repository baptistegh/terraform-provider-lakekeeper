package provider

import (
	"context"
	"fmt"

	permissionv1 "github.com/baptistegh/go-lakekeeper/pkg/apis/management/v1/permission"
	lakekeeper "github.com/baptistegh/go-lakekeeper/pkg/client"
	"github.com/baptistegh/go-lakekeeper/pkg/core"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces
var (
	_ resource.Resource                = &lakekeeperServerRoleAssignmentResource{}
	_ resource.ResourceWithConfigure   = &lakekeeperServerRoleAssignmentResource{}
	_ resource.ResourceWithImportState = &lakekeeperServerRoleAssignmentResource{}
)

func init() {
	registerResource(NewLakekeeperServerRoleAssignment)
}

// NewLakekeeperServerAssignment is a helper function to simplify the provider implementation.
func NewLakekeeperServerRoleAssignment() resource.Resource {
	return &lakekeeperServerRoleAssignmentResource{}
}

func (r *lakekeeperServerRoleAssignmentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_server_role_assignment"
}

// lakekeeperServerRoleAssignmentResource defines the resource implementation.
type lakekeeperServerRoleAssignmentResource struct {
	client *lakekeeper.Client
}

// lakekeeperServerRoleAssignmentResourceModel describes the resource data model.
type lakekeeperServerRoleAssignmentResourceModel struct {
	ID          types.String `tfsdk:"id"` // form: project_id:role_id (internal ID)
	RoleID      types.String `tfsdk:"role_id"`
	Assignments types.Set    `tfsdk:"assignments"`
}

func (r *lakekeeperServerRoleAssignmentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: fmt.Sprintf(`The ` + "`lakekeeper_server_role_assignment`" + ` resource allows to manage the lifecycle of a role assignement to the server.

**Upstream API**: [Lakekeeper REST API docs](https://docs.lakekeeper.io/docs/nightly/api/management/#tag/permissions/operation/update_server_assignments)`),
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The internal ID of this resource. In the form: <role_id>",
				Computed:            true,
			},
			"role_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the role.",
				Required:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"assignments": schema.SetAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "List of server assignments for this role. values can be `operator` or `admin`",
				Required:            true,
				Validators: []validator.Set{setvalidator.ValueStringsAre(
					stringvalidator.OneOf(string(permissionv1.AdminServerAssignment), string(permissionv1.OperatorServerAssignment)),
				)},
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *lakekeeperServerRoleAssignmentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	resourceData := req.ProviderData.(*LakekeeperResourceData)
	r.client = resourceData.Client
}

// Create creates a new upstream resources and adds it into the Terraform state.
func (r *lakekeeperServerRoleAssignmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state, plan lakekeeperServerRoleAssignmentResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var opts []*permissionv1.ServerAssignment

	for _, v := range plan.Assignments.Elements() {
		s, ok := v.(types.String)
		if !ok {
			resp.Diagnostics.AddError("Error converting model to resource", fmt.Sprintf("Unable to read assignment %s", v.String()))
			return
		}

		opts = append(opts, &permissionv1.ServerAssignment{
			Assignee: permissionv1.UserOrRole{
				Value: plan.RoleID.ValueString(),
				Type:  permissionv1.RoleType,
			},
			Assignment: permissionv1.ServerAssignmentType(s.ValueString()),
		})
	}

	_, err := r.client.PermissionV1().ServerPermission().Update(&permissionv1.UpdateServerPermissionsOptions{
		Writes: opts,
	}, core.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError("Lakekeeper API error occurred", fmt.Sprintf("Unable to write server assignment, %s", err.Error()))
		return
	}

	state.ID = plan.RoleID
	state.RoleID = plan.RoleID
	state.Assignments = plan.Assignments

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *lakekeeperServerRoleAssignmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state lakekeeperServerRoleAssignmentResourceModel

	// Read Terraform state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	roleID := state.ID.ValueString()
	state.RoleID = types.StringValue(roleID)

	assignments, _, err := r.client.PermissionV1().ServerPermission().GetAssignments(nil, core.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError("Lakekeeper API error occurred", fmt.Sprintf("Unable to read server assignments, %s", err.Error()))
		return
	}

	var elems []attr.Value
	for _, v := range assignments.Assignments {
		if v.Assignee.Value == roleID && v.Assignee.Type == permissionv1.RoleType {
			elems = append(elems, types.StringValue(string(v.Assignment)))
		}
	}

	newAssignments, diags := types.SetValue(types.StringType, elems)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	state.Assignments = newAssignments
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Updates updates the resource in-place.
func (r *lakekeeperServerRoleAssignmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state lakekeeperServerRoleAssignmentResourceModel

	// Read Terraform state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var planAssignments []types.String
	resp.Diagnostics.Append(plan.Assignments.ElementsAs(ctx, &planAssignments, false)...)

	var stateAssignments []types.String
	resp.Diagnostics.Append(state.Assignments.ElementsAs(ctx, &stateAssignments, false)...)

	if resp.Diagnostics.HasError() {
		return
	}

	added, removed := DiffTypedStrings(stateAssignments, planAssignments)

	var writes []*permissionv1.ServerAssignment
	var deletes []*permissionv1.ServerAssignment

	for _, v := range added {
		writes = append(writes, &permissionv1.ServerAssignment{
			Assignee: permissionv1.UserOrRole{
				Type:  permissionv1.RoleType,
				Value: plan.RoleID.ValueString(),
			},
			Assignment: permissionv1.ServerAssignmentType(v.ValueString()),
		})
	}

	for _, v := range removed {
		deletes = append(deletes, &permissionv1.ServerAssignment{
			Assignee: permissionv1.UserOrRole{
				Type:  permissionv1.RoleType,
				Value: plan.RoleID.ValueString(),
			},
			Assignment: permissionv1.ServerAssignmentType(v.ValueString()),
		})
	}

	if _, err := r.client.PermissionV1().ServerPermission().Update(&permissionv1.UpdateServerPermissionsOptions{
		Writes:  writes,
		Deletes: deletes,
	}, core.WithContext(ctx)); err != nil {
		resp.Diagnostics.AddError("Lakekeeper API error occurred", fmt.Sprintf("Unable to update server assignment, %v", err.Error()))
		return
	}

	state.Assignments = plan.Assignments

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Deletes removes the resource.
func (r *lakekeeperServerRoleAssignmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state lakekeeperServerRoleAssignmentResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var assignments []types.String
	var deletes []*permissionv1.ServerAssignment

	resp.Diagnostics.Append(state.Assignments.ElementsAs(ctx, &assignments, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	for _, v := range assignments {
		deletes = append(deletes, &permissionv1.ServerAssignment{
			Assignee: permissionv1.UserOrRole{
				Value: state.RoleID.ValueString(),
				Type:  permissionv1.RoleType,
			},
			Assignment: permissionv1.ServerAssignmentType(v.ValueString()),
		})
	}

	if _, err := r.client.PermissionV1().ServerPermission().Update(&permissionv1.UpdateServerPermissionsOptions{
		Deletes: deletes,
	}, core.WithContext(ctx)); err != nil {
		resp.Diagnostics.AddError("Lakekeeper API error occurred", fmt.Sprintf("Unable to delete server assignment, %v", err.Error()))
	}

	resp.State.RemoveResource(ctx)
}

// ImportState imports the resource into the Terraform state.
func (r *lakekeeperServerRoleAssignmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
