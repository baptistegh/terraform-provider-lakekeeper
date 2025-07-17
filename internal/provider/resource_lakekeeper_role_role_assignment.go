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
	_ resource.Resource                = &lakekeeperRoleRoleAssignmentResource{}
	_ resource.ResourceWithConfigure   = &lakekeeperRoleRoleAssignmentResource{}
	_ resource.ResourceWithImportState = &lakekeeperRoleRoleAssignmentResource{}
)

func init() {
	registerResource(NewLakekeeperRoleRoleAssignment)
}

// NewLakekeeperRoleAssignment is a helper function to simplify the provider implementation.
func NewLakekeeperRoleRoleAssignment() resource.Resource {
	return &lakekeeperRoleRoleAssignmentResource{}
}

func (r *lakekeeperRoleRoleAssignmentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_role_role_assignment"
}

// lakekeeperRoleRoleAssignmentResource defines the resource implementation.
type lakekeeperRoleRoleAssignmentResource struct {
	client *lakekeeper.Client
}

// lakekeeperRoleRoleAssignmentResourceModel describes the resource data model.
type lakekeeperRoleRoleAssignmentResourceModel struct {
	ID          types.String `tfsdk:"id"` // form: role_id:assignee_id (internal ID)
	RoleID      types.String `tfsdk:"role_id"`
	AssigneeID  types.String `tfsdk:"assignee_id"`
	Assignments types.Set    `tfsdk:"assignments"`
}

func (r *lakekeeperRoleRoleAssignmentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: fmt.Sprintf(`The ` + "`lakekeeper_role_role_assignment`" + ` resource allows to manage the lifecycle of a role assignement to a role.

**Upstream API**: [Lakekeeper REST API docs](https://docs.lakekeeper.io/docs/nightly/api/management/#tag/permissions/operation/update_role_assignments)`),
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The internal ID of this resource. In the form: <role_id>:<assignee_id>",
				Computed:            true,
			},
			"role_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the role.",
				Required:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"assignee_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the role to assign to the role.",
				Required:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"assignments": schema.SetAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "List of role assignments for this role. values can be `ownership` or `assignee`",
				Required:            true,
				Validators: []validator.Set{setvalidator.ValueStringsAre(
					stringvalidator.OneOf(string(permissionv1.OwnershipRoleAssignment), string(permissionv1.AssigneeRoleAssignment)),
				)},
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *lakekeeperRoleRoleAssignmentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	resourceData := req.ProviderData.(*LakekeeperResourceData)
	r.client = resourceData.Client
}

// Create creates a new upstream resources and adds it into the Terraform state.
func (r *lakekeeperRoleRoleAssignmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state, plan lakekeeperRoleRoleAssignmentResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var opts []*permissionv1.RoleAssignment

	for _, v := range plan.Assignments.Elements() {
		s, ok := v.(types.String)
		if !ok {
			resp.Diagnostics.AddError("Error converting model to resource", fmt.Sprintf("Unable to read assignment %s", v.String()))
			return
		}

		opts = append(opts, &permissionv1.RoleAssignment{
			Assignee: permissionv1.UserOrRole{
				Value: plan.AssigneeID.ValueString(),
				Type:  permissionv1.RoleType,
			},
			Assignment: permissionv1.RoleAssignmentType(s.ValueString()),
		})
	}

	_, err := r.client.PermissionV1().RolePermission().Update(plan.RoleID.ValueString(), &permissionv1.UpdateRolePermissionsOptions{
		Writes: opts,
	}, core.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError("Lakekeeper API error occurred", fmt.Sprintf("Unable to write role assignment, %s", err.Error()))
		return
	}

	state.ID = types.StringValue(fmt.Sprintf("%s:%s", plan.RoleID.ValueString(), plan.AssigneeID.ValueString()))

	state.RoleID = plan.RoleID
	state.AssigneeID = plan.AssigneeID
	state.Assignments = plan.Assignments

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *lakekeeperRoleRoleAssignmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state lakekeeperRoleRoleAssignmentResourceModel

	// Read Terraform state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	roleID, assigneeID := splitInternalID(state.ID)

	assignments, _, err := r.client.PermissionV1().RolePermission().GetAssignments(roleID, nil, core.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError("Lakekeeper API error occurred", fmt.Sprintf("Unable to read role assignments, %s", err.Error()))
		return
	}

	var elems []attr.Value
	for _, v := range assignments.Assignments {
		if v.Assignee.Value == assigneeID && v.Assignee.Type == permissionv1.RoleType {
			elems = append(elems, types.StringValue(string(v.Assignment)))
		}
	}

	newAssignments, diags := types.SetValue(types.StringType, elems)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	state.RoleID = types.StringValue(roleID)
	state.AssigneeID = types.StringValue(assigneeID)

	state.Assignments = newAssignments
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Updates updates the resource in-place.
func (r *lakekeeperRoleRoleAssignmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state lakekeeperRoleRoleAssignmentResourceModel

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

	var writes []*permissionv1.RoleAssignment
	var deletes []*permissionv1.RoleAssignment

	for _, v := range added {
		writes = append(writes, &permissionv1.RoleAssignment{
			Assignee: permissionv1.UserOrRole{
				Type:  permissionv1.RoleType,
				Value: plan.AssigneeID.ValueString(),
			},
			Assignment: permissionv1.RoleAssignmentType(v.ValueString()),
		})
	}

	for _, v := range removed {
		deletes = append(deletes, &permissionv1.RoleAssignment{
			Assignee: permissionv1.UserOrRole{
				Type:  permissionv1.RoleType,
				Value: plan.AssigneeID.ValueString(),
			},
			Assignment: permissionv1.RoleAssignmentType(v.ValueString()),
		})
	}

	roleID, assigneeID := splitInternalID(state.ID)

	if _, err := r.client.PermissionV1().RolePermission().Update(roleID, &permissionv1.UpdateRolePermissionsOptions{
		Writes:  writes,
		Deletes: deletes,
	}, core.WithContext(ctx)); err != nil {
		resp.Diagnostics.AddError("Lakekeeper API error occurred", fmt.Sprintf("Unable to update role assignment, %v", err.Error()))
		return
	}

	state.AssigneeID = types.StringValue(assigneeID)
	state.Assignments = plan.Assignments

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Deletes removes the resource.
func (r *lakekeeperRoleRoleAssignmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state lakekeeperRoleRoleAssignmentResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var assignments []types.String
	var deletes []*permissionv1.RoleAssignment

	resp.Diagnostics.Append(state.Assignments.ElementsAs(ctx, &assignments, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	for _, v := range assignments {
		deletes = append(deletes, &permissionv1.RoleAssignment{
			Assignee: permissionv1.UserOrRole{
				Value: state.AssigneeID.ValueString(),
				Type:  permissionv1.RoleType,
			},
			Assignment: permissionv1.RoleAssignmentType(v.ValueString()),
		})
	}

	if _, err := r.client.PermissionV1().RolePermission().Update(state.RoleID.ValueString(), &permissionv1.UpdateRolePermissionsOptions{
		Deletes: deletes,
	}, core.WithContext(ctx)); err != nil {
		resp.Diagnostics.AddError("Lakekeeper API error occurred", fmt.Sprintf("Unable to delete role assignment, %v", err.Error()))
	}

	resp.State.RemoveResource(ctx)
}

// ImportState imports the resource into the Terraform state.
func (r *lakekeeperRoleRoleAssignmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
