package provider

import (
	"context"
	"fmt"

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
	_ resource.Resource                = &lakekeeperProjectUserAssignmentResource{}
	_ resource.ResourceWithConfigure   = &lakekeeperProjectUserAssignmentResource{}
	_ resource.ResourceWithImportState = &lakekeeperProjectUserAssignmentResource{}
)

func init() {
	registerResource(NewLakekeeperProjectUserAssignment)
}

// NewLakekeeperProjectAssignment is a helper function to simplify the provider implementation.
func NewLakekeeperProjectUserAssignment() resource.Resource {
	return &lakekeeperProjectUserAssignmentResource{}
}

func (r *lakekeeperProjectUserAssignmentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project_user_assignment"
}

// lakekeeperProjectUserAssignmentResource defines the resource implementation.
type lakekeeperProjectUserAssignmentResource struct {
	client *lakekeeper.Client
}

// lakekeeperProjectUserAssignmentResourceModel describes the resource data model.
type lakekeeperProjectUserAssignmentResourceModel struct {
	ID          types.String `tfsdk:"id"` // form: project_id:user_id (internal ID)
	ProjectID   types.String `tfsdk:"project_id"`
	UserID      types.String `tfsdk:"user_id"`
	Assignments types.Set    `tfsdk:"assignments"`
}

func (r *lakekeeperProjectUserAssignmentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: fmt.Sprintf(`The ` + "`lakekeeper_project_user_assignment`" + ` resource allows to manage the lifecycle of a user assignement to a project.

**Upstream API**: [Lakekeeper REST API docs](https://docs.lakekeeper.io/docs/nightly/api/management/#tag/permissions/operation/update_project_assignments)`),
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The internal ID of this resource. In the form: <project_id>:<user_id>",
				Computed:            true,
			},
			"project_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the project.",
				Required:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"user_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the user.",
				Required:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"assignments": schema.SetAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "List of project assignments for this user. values can be `project_admin` `security_admin` `data_admin` `role_creator` `describe` `select` `create` `modify`",
				Required:            true,
				Validators: []validator.Set{setvalidator.ValueStringsAre(
					stringvalidator.OneOf(
						string(permissionv1.AdminProjectAssignment),
						string(permissionv1.SecurityAdminProjectAssignment),
						string(permissionv1.DataAdminProjectAssignment),
						string(permissionv1.RoleCreatorProjectAssignment),
						string(permissionv1.DescribeProjectAssignment),
						string(permissionv1.SelectProjectAssignment),
						string(permissionv1.CreateProjectAssignment),
						string(permissionv1.ModifyProjectAssignment),
					),
				)},
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *lakekeeperProjectUserAssignmentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	resourceData := req.ProviderData.(*LakekeeperResourceData)
	r.client = resourceData.Client
}

// Create creates a new upstream resources and adds it into the Terraform state.
func (r *lakekeeperProjectUserAssignmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state, plan lakekeeperProjectUserAssignmentResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	projectID := plan.ProjectID.ValueString()
	userID := plan.UserID.ValueString()

	id := fmt.Sprintf("%s:%s", projectID, userID)

	var opts []*permissionv1.ProjectAssignment

	for _, v := range plan.Assignments.Elements() {
		s, ok := v.(types.String)
		if !ok {
			resp.Diagnostics.AddError("Error converting model to resource", fmt.Sprintf("Unable to read assignment %s", v.String()))
			return
		}

		opts = append(opts, &permissionv1.ProjectAssignment{
			Assignee: permissionv1.UserOrRole{
				Value: plan.UserID.ValueString(),
				Type:  permissionv1.UserType,
			},
			Assignment: permissionv1.ProjectAssignmentType(s.ValueString()),
		})
	}

	_, err := r.client.PermissionV1().ProjectPermission().Update(ctx, projectID, &permissionv1.UpdateProjectPermissionsOptions{
		Writes: opts,
	})
	if err != nil {
		resp.Diagnostics.AddError("Lakekeeper API error occurred", fmt.Sprintf("Unable to write project assignment, %s", err.Error()))
		return
	}

	state.ID = types.StringValue(id)
	state.ProjectID = plan.ProjectID
	state.UserID = plan.UserID
	state.Assignments = plan.Assignments

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *lakekeeperProjectUserAssignmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state lakekeeperProjectUserAssignmentResourceModel

	// Read Terraform state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	projectID, userID := splitInternalID(state.ID)

	assignments, _, err := r.client.PermissionV1().ProjectPermission().GetAssignments(ctx, projectID, nil)
	if err != nil {
		resp.Diagnostics.AddError("Lakekeeper API error occurred", fmt.Sprintf("Unable to read project assignments, %s", err.Error()))
		return
	}

	var elems []attr.Value
	for _, v := range assignments.Assignments {
		if v.Assignee.Value == userID && v.Assignee.Type == permissionv1.UserType {
			elems = append(elems, types.StringValue(string(v.Assignment)))
		}
	}

	newAssignments, diags := types.SetValue(types.StringType, elems)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	state.UserID = types.StringValue(userID)
	state.ProjectID = types.StringValue(projectID)
	state.Assignments = newAssignments

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Updates updates the resource in-place.
func (r *lakekeeperProjectUserAssignmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state lakekeeperProjectUserAssignmentResourceModel

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

	var writes []*permissionv1.ProjectAssignment
	var deletes []*permissionv1.ProjectAssignment

	for _, v := range added {
		writes = append(writes, &permissionv1.ProjectAssignment{
			Assignee: permissionv1.UserOrRole{
				Type:  permissionv1.UserType,
				Value: plan.UserID.ValueString(),
			},
			Assignment: permissionv1.ProjectAssignmentType(v.ValueString()),
		})
	}

	for _, v := range removed {
		deletes = append(deletes, &permissionv1.ProjectAssignment{
			Assignee: permissionv1.UserOrRole{
				Type:  permissionv1.UserType,
				Value: plan.UserID.ValueString(),
			},
			Assignment: permissionv1.ProjectAssignmentType(v.ValueString()),
		})
	}

	if _, err := r.client.PermissionV1().ProjectPermission().Update(ctx, plan.ProjectID.ValueString(), &permissionv1.UpdateProjectPermissionsOptions{
		Writes:  writes,
		Deletes: deletes,
	}); err != nil {
		resp.Diagnostics.AddError("Lakekeeper API error occurred", fmt.Sprintf("Unable to update project assignment, %v", err.Error()))
		return
	}

	state.Assignments = plan.Assignments

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Deletes removes the resource.
func (r *lakekeeperProjectUserAssignmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state lakekeeperProjectUserAssignmentResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var assignments []types.String
	var deletes []*permissionv1.ProjectAssignment

	resp.Diagnostics.Append(state.Assignments.ElementsAs(ctx, &assignments, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	for _, v := range assignments {
		deletes = append(deletes, &permissionv1.ProjectAssignment{
			Assignee: permissionv1.UserOrRole{
				Value: state.UserID.ValueString(),
				Type:  permissionv1.UserType,
			},
			Assignment: permissionv1.ProjectAssignmentType(v.ValueString()),
		})
	}

	if _, err := r.client.PermissionV1().ProjectPermission().Update(ctx, state.ProjectID.ValueString(), &permissionv1.UpdateProjectPermissionsOptions{
		Deletes: deletes,
	}); err != nil {
		resp.Diagnostics.AddError("Lakekeeper API error occurred", fmt.Sprintf("Unable to delete project assignment, %v", err.Error()))
	}

	resp.State.RemoveResource(ctx)
}

// ImportState imports the resource into the Terraform state.
func (r *lakekeeperProjectUserAssignmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
