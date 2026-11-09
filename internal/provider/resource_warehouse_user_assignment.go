package provider

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	managementv1 "github.com/lakekeeper/go-lakekeeper/pkg/apis/management/v1"
	lakekeeper "github.com/lakekeeper/go-lakekeeper/pkg/client"
	"github.com/lakekeeper/go-lakekeeper/pkg/permissions"
)

// Ensure provider defined types fully satisfy framework interfaces
var (
	_ resource.Resource                = &lakekeeperWarehouseUserAssignmentResource{}
	_ resource.ResourceWithConfigure   = &lakekeeperWarehouseUserAssignmentResource{}
	_ resource.ResourceWithImportState = &lakekeeperWarehouseUserAssignmentResource{}
)

func init() {
	registerResource(NewLakekeeperWarehouseUserAssignment)
}

// NewLakekeeperWarehouseAssignment is a helper function to simplify the provider implementation.
func NewLakekeeperWarehouseUserAssignment() resource.Resource {
	return &lakekeeperWarehouseUserAssignmentResource{}
}

func (r *lakekeeperWarehouseUserAssignmentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_warehouse_user_assignment"
}

// lakekeeperWarehouseUserAssignmentResource defines the resource implementation.
type lakekeeperWarehouseUserAssignmentResource struct {
	client *lakekeeper.Client
}

// lakekeeperWarehouseUserAssignmentResourceModel describes the resource data model.
type lakekeeperWarehouseUserAssignmentResourceModel struct {
	ID          types.String `tfsdk:"id"` // form: <warehouse_id>:<role_id>
	WarehouseID types.String `tfsdk:"warehouse_id"`
	UserID      types.String `tfsdk:"user_id"`
	Assignments types.Set    `tfsdk:"assignments"`
}

func (r *lakekeeperWarehouseUserAssignmentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: fmt.Sprintf(`The ` + "`lakekeeper_warehouse_user_assignment`" + ` resource allows to manage the lifecycle of a user assignement to a warehouse.

**Upstream API**: [Lakekeeper REST API docs](https://docs.lakekeeper.io/docs/nightly/api/management/#tag/permissions/operation/update_warehouse_assignments)`),
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The internal ID of this resource. In the form: `{{warehouse_id}}/{{user_id}}`",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"warehouse_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the warehouse.",
				Required:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile("^[^/]+$"), "must be a warehouse UUID and NOT include the project UUID"),
				},
			},
			"user_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the user to assign to this warehouse.",
				Required:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"assignments": schema.SetAttribute{
				ElementType: types.StringType,
				MarkdownDescription: "List of assignments the user has on this warehouse. values can be " +
					enumValuesList(managementv1.AllowedWarehouseRelationEnumValues),
				Required: true,
				Validators: []validator.Set{setvalidator.ValueStringsAre(
					stringvalidator.OneOf(enumStrings(managementv1.AllowedWarehouseRelationEnumValues)...)),
				},
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *lakekeeperWarehouseUserAssignmentResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = resourceClient(req, resp)
}

// Create creates a new upstream resources and adds it into the Terraform state.
func (r *lakekeeperWarehouseUserAssignmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state, plan lakekeeperWarehouseUserAssignmentResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var assignments []string
	resp.Diagnostics.Append(plan.Assignments.ElementsAs(ctx, &assignments, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := warehouseAssignmentScope(r.client, plan.WarehouseID.ValueString()).apply(ctx, permissions.PrincipalUser, plan.UserID.ValueString(), assignments, nil); err != nil {
		resp.Diagnostics.AddError("Lakekeeper API error occurred", fmt.Sprintf("Unable to write warehouse assignment, %s", err.Error()))
		return
	}

	state.ID = types.StringValue(fmt.Sprintf("%s/%s", plan.WarehouseID.ValueString(), plan.UserID.ValueString()))
	state.WarehouseID = plan.WarehouseID
	state.UserID = plan.UserID
	state.Assignments = plan.Assignments

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *lakekeeperWarehouseUserAssignmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state lakekeeperWarehouseUserAssignmentResourceModel

	// Read Terraform state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	objectID, principalID := splitInternalID(state.ID)

	relations, err := warehouseAssignmentScope(r.client, objectID).relations(ctx, permissions.PrincipalUser, principalID)
	if err != nil {
		resp.Diagnostics.AddError("Lakekeeper API error occurred", fmt.Sprintf("Unable to read warehouse assignments, %s", err.Error()))
		return
	}

	newAssignments, diags := types.SetValueFrom(ctx, types.StringType, relations)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.WarehouseID = types.StringValue(objectID)
	state.UserID = types.StringValue(principalID)
	state.Assignments = newAssignments
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update updates the resource in-place.
func (r *lakekeeperWarehouseUserAssignmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state lakekeeperWarehouseUserAssignmentResourceModel

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

	objectID, principalID := splitInternalID(state.ID)

	if err := warehouseAssignmentScope(r.client, objectID).apply(ctx, permissions.PrincipalUser, principalID, stringValues(added), stringValues(removed)); err != nil {
		resp.Diagnostics.AddError("Lakekeeper API error occurred", fmt.Sprintf("Unable to update warehouse assignment, %v", err.Error()))
		return
	}

	state.UserID = types.StringValue(principalID)
	state.Assignments = plan.Assignments

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Delete removes the resource.
func (r *lakekeeperWarehouseUserAssignmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state lakekeeperWarehouseUserAssignmentResourceModel

	// Read Terraform state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var assignments []string
	resp.Diagnostics.Append(state.Assignments.ElementsAs(ctx, &assignments, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := warehouseAssignmentScope(r.client, state.WarehouseID.ValueString()).apply(ctx, permissions.PrincipalUser, state.UserID.ValueString(), nil, assignments); err != nil {
		resp.Diagnostics.AddError("Lakekeeper API error occurred", fmt.Sprintf("Unable to delete warehouse assignment, %v", err.Error()))
	}

	resp.State.RemoveResource(ctx)
}

// ImportState imports the resource into the Terraform state.
func (r *lakekeeperWarehouseUserAssignmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Expected format: "warehouse_id/user_id"
	parts := strings.Split(req.ID, "/")
	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			"Invalid imp&ort ID format",
			"Expected format: warehouse_id/user_id",
		)
		return
	}

	resp.State.SetAttribute(ctx, path.Root("warehouse_id"), parts[0])
	resp.State.SetAttribute(ctx, path.Root("user_id"), parts[1])

	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
