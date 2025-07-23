package provider

import (
	"context"
	"fmt"
	"strings"

	permissionv1 "github.com/baptistegh/go-lakekeeper/pkg/apis/management/v1/permission"
	lakekeeper "github.com/baptistegh/go-lakekeeper/pkg/client"
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
	_ resource.Resource                = &lakekeeperWarehouseRoleAssignmentResource{}
	_ resource.ResourceWithConfigure   = &lakekeeperWarehouseRoleAssignmentResource{}
	_ resource.ResourceWithImportState = &lakekeeperWarehouseRoleAssignmentResource{}
)

func init() {
	registerResource(NewLakekeeperWarehouseRoleAssignment)
}

// NewLakekeeperWarehouseAssignment is a helper function to simplify the provider implementation.
func NewLakekeeperWarehouseRoleAssignment() resource.Resource {
	return &lakekeeperWarehouseRoleAssignmentResource{}
}

func (r *lakekeeperWarehouseRoleAssignmentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_warehouse_role_assignment"
}

// lakekeeperWarehouseRoleAssignmentResource defines the resource implementation.
type lakekeeperWarehouseRoleAssignmentResource struct {
	client *lakekeeper.Client
}

// lakekeeperWarehouseRoleAssignmentResourceModel describes the resource data model.
type lakekeeperWarehouseRoleAssignmentResourceModel struct {
	ID          types.String `tfsdk:"id"` // form: role_id:warehouse_id (internal ID)
	WarehouseID types.String `tfsdk:"warehouse_id"`
	RoleID      types.String `tfsdk:"role_id"`
	Assignments types.Set    `tfsdk:"assignments"`
}

func (r *lakekeeperWarehouseRoleAssignmentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: fmt.Sprintf(`The ` + "`lakekeeper_warehouse_role_assignment`" + ` resource allows to manage the lifecycle of a role assignement to a warehouse.

**Upstream API**: [Lakekeeper REST API docs](https://docs.lakekeeper.io/docs/nightly/api/management/#tag/permissions/operation/update_warehouse_assignments)`),
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The internal ID of this resource. In the form: <warehouse_id>:<role_id>",
				Computed:            true,
			},
			"warehouse_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the warehouse.",
				Required:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"role_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the role to assign to the role.",
				Required:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"assignments": schema.SetAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "List of role assignments for this warehouse. values can be `ownership`, `pass_grants_admin`, `manage_grants_admin`, `describe_warehouse`, `select_warehouse`, `create_warehouse`, `modify_warehouse`",
				Required:            true,
				Validators: []validator.Set{setvalidator.ValueStringsAre(
					stringvalidator.OneOf(
						string(permissionv1.OwnershipWarehouseAssignment),
						string(permissionv1.PassGrantsAdminWarehouseAssignment),
						string(permissionv1.ManageGrantsAdminWarehouseAssignment),
						string(permissionv1.DescribeWarehouseAssignment),
						string(permissionv1.SelectWarehouseAssignment),
						string(permissionv1.CreateWarehouseAssignment),
						string(permissionv1.ModifyWarehouseAssignment),
					)),
				},
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *lakekeeperWarehouseRoleAssignmentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	resourceData := req.ProviderData.(*LakekeeperResourceData)
	r.client = resourceData.Client
}

// Create creates a new upstream resources and adds it into the Terraform state.
func (r *lakekeeperWarehouseRoleAssignmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state, plan lakekeeperWarehouseRoleAssignmentResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var opts []*permissionv1.WarehouseAssignment

	for _, v := range plan.Assignments.Elements() {
		s, ok := v.(types.String)
		if !ok {
			resp.Diagnostics.AddError("Error converting model to resource", fmt.Sprintf("Unable to read assignment %s", v.String()))
			return
		}

		opts = append(opts, &permissionv1.WarehouseAssignment{
			Assignee: permissionv1.UserOrRole{
				Value: plan.RoleID.ValueString(),
				Type:  permissionv1.RoleType,
			},
			Assignment: permissionv1.WarehouseAssignmentType(s.ValueString()),
		})
	}

	_, err := r.client.PermissionV1().WarehousePermission().Update(ctx, plan.WarehouseID.ValueString(), &permissionv1.UpdateWarehousePermissionsOptions{
		Writes: opts,
	})
	if err != nil {
		resp.Diagnostics.AddError("Lakekeeper API error occurred", fmt.Sprintf("Unable to write warehouse assignment, %s", err.Error()))
		return
	}

	state.ID = types.StringValue(fmt.Sprintf("%s/%s", plan.WarehouseID.ValueString(), plan.RoleID.ValueString()))

	state.WarehouseID = plan.WarehouseID
	state.RoleID = plan.RoleID
	state.Assignments = plan.Assignments

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *lakekeeperWarehouseRoleAssignmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state lakekeeperWarehouseRoleAssignmentResourceModel

	// Read Terraform state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	warehouseID, roleID := splitInternalID(state.ID)

	assignments, _, err := r.client.PermissionV1().WarehousePermission().GetAssignments(ctx, warehouseID, nil)
	if err != nil {
		resp.Diagnostics.AddError("Lakekeeper API error occurred", fmt.Sprintf("Unable to read role assignments, %s", err.Error()))
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

	state.WarehouseID = types.StringValue(warehouseID)
	state.RoleID = types.StringValue(roleID)

	state.Assignments = newAssignments
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Updates updates the resource in-place.
func (r *lakekeeperWarehouseRoleAssignmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state lakekeeperWarehouseRoleAssignmentResourceModel

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

	var writes []*permissionv1.WarehouseAssignment
	var deletes []*permissionv1.WarehouseAssignment

	for _, v := range added {
		writes = append(writes, &permissionv1.WarehouseAssignment{
			Assignee: permissionv1.UserOrRole{
				Type:  permissionv1.RoleType,
				Value: plan.RoleID.ValueString(),
			},
			Assignment: permissionv1.WarehouseAssignmentType(v.ValueString()),
		})
	}

	for _, v := range removed {
		deletes = append(deletes, &permissionv1.WarehouseAssignment{
			Assignee: permissionv1.UserOrRole{
				Type:  permissionv1.RoleType,
				Value: plan.RoleID.ValueString(),
			},
			Assignment: permissionv1.WarehouseAssignmentType(v.ValueString()),
		})
	}

	warehouseID, roleID := splitInternalID(state.ID)

	if _, err := r.client.PermissionV1().WarehousePermission().Update(ctx, warehouseID, &permissionv1.UpdateWarehousePermissionsOptions{
		Writes:  writes,
		Deletes: deletes,
	}); err != nil {
		resp.Diagnostics.AddError("Lakekeeper API error occurred", fmt.Sprintf("Unable to update role assignment, %v", err.Error()))
		return
	}

	state.RoleID = types.StringValue(roleID)
	state.Assignments = plan.Assignments

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Deletes removes the resource.
func (r *lakekeeperWarehouseRoleAssignmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state lakekeeperWarehouseRoleAssignmentResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var assignments []types.String
	var deletes []*permissionv1.WarehouseAssignment

	resp.Diagnostics.Append(state.Assignments.ElementsAs(ctx, &assignments, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	for _, v := range assignments {
		deletes = append(deletes, &permissionv1.WarehouseAssignment{
			Assignee: permissionv1.UserOrRole{
				Value: state.RoleID.ValueString(),
				Type:  permissionv1.RoleType,
			},
			Assignment: permissionv1.WarehouseAssignmentType(v.ValueString()),
		})
	}

	if _, err := r.client.PermissionV1().WarehousePermission().Update(ctx, state.WarehouseID.ValueString(), &permissionv1.UpdateWarehousePermissionsOptions{
		Deletes: deletes,
	}); err != nil {
		resp.Diagnostics.AddError("Lakekeeper API error occurred", fmt.Sprintf("Unable to delete role assignment, %v", err.Error()))
	}

	resp.State.RemoveResource(ctx)
}

// ImportState imports the resource into the Terraform state.
func (r *lakekeeperWarehouseRoleAssignmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Expected format: "warehouse_id/role_id"
	parts := strings.Split(req.ID, "/")
	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			"Invalid import ID format",
			"Expected format: warehouse_id/role_id",
		)
		return
	}

	resp.State.SetAttribute(ctx, path.Root("warehouse_id"), parts[0])
	resp.State.SetAttribute(ctx, path.Root("role_id"), parts[1])

	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
